package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// --- Configuration ---
const (
	ServerPort  = ":8080"
	LogFile     = "waf-logs.json"
	RulesFile   = "waf-rules.json"
	BlockedFile = "waf-blocked.json" // Reusing this name for unified IP config
	FiltersFile = "waf-filters.json"
)

// --- Models ---
type WAFLog struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	IP        string    `json:"ip"`
	Method    string    `json:"method"`
	URL       string    `json:"url"`
	Status    string    `json:"status"`
	RuleID    string    `json:"rule_id,omitempty"`
	Details   string    `json:"details,omitempty"`
	Geo       string    `json:"geo,omitempty"`
}

type WAFRule struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Pattern     string   `json:"pattern"`
	Tags        []string `json:"tags,omitempty"`
	Impact      string   `json:"impact,omitempty"`
	Enabled     bool     `json:"enabled"`
}

type IPRecord struct {
	IP     string `json:"ip"`
	Reason string `json:"reason,omitempty"`
}

type IPConfig struct {
	BlockedIPs  []IPRecord `json:"blocked_ips"`
	Whitelisted []string   `json:"whitelisted_ips"`
	IgnoredIPs  []string   `json:"ignored_ips"`
}

type FilterConfig struct {
	Keywords        []string `json:"keywords"`
	WhitelistedURLs []string `json:"whitelisted_urls"`
}

// --- In-Memory Stores ---
var (
	wafLogs   []WAFLog
	logsMutex sync.RWMutex

	blockedIPs     = make(map[string]string)
	whitelistedIPs = make(map[string]bool)
	ignoredIPs     = make(map[string]bool) // New: Mute logs
	ipMutex        sync.RWMutex

	rules      []WAFRule
	rulesMutex sync.RWMutex

	filterConfig FilterConfig
	filterMutex  sync.RWMutex
)

// --- Initialization ---
func init() {
	loadRules()
	loadLogs()
	loadIPConfig()
	loadFilters()
}

func loadRules() {
	file, err := os.Open(RulesFile)
	if err != nil {
		log.Println("Rules file not found, loading defaults.")
		rules = []WAFRule{
			{ID: "1001", Description: "SQL Injection Detection (Basic)", Pattern: `(?i)(union\s+select|select\s+.*\s+from|delete\s+from|drop\s+table|update\s+set)`},
		}
		return
	}
	defer file.Close()
	byteValue, _ := io.ReadAll(file)
	json.Unmarshal(byteValue, &rules)
	fmt.Printf("[INIT] Loaded %d rules\n", len(rules))
}

func loadLogs() {
	file, err := os.Open(LogFile)
	if err != nil {
		return
	}
	defer file.Close()
	byteValue, _ := io.ReadAll(file)
	if err := json.Unmarshal(byteValue, &wafLogs); err != nil {
		fmt.Println("Error loading logs, starting fresh")
		wafLogs = []WAFLog{}
	} else {
		fmt.Printf("[INIT] Restored %d logs\n", len(wafLogs))
	}
}

func loadIPConfig() {
	file, err := os.Open(BlockedFile)
	if err != nil {
		fmt.Println("[INIT] No existing IP config found, starting fresh.")
		return
	}
	defer file.Close()

	var config IPConfig
	byteValue, _ := io.ReadAll(file)
	if err := json.Unmarshal(byteValue, &config); err != nil {
		fmt.Println("[INIT] Error parsing IP config.")
		return
	}

	ipMutex.Lock()
	defer ipMutex.Unlock()

	for _, record := range config.BlockedIPs {
		blockedIPs[record.IP] = record.Reason
	}
	for _, ip := range config.Whitelisted {
		whitelistedIPs[ip] = true
	}
	for _, ip := range config.IgnoredIPs {
		ignoredIPs[ip] = true
	}
	fmt.Printf("[INIT] Loaded IP Config: %d Banned, %d Whitelisted, %d Ignored\n", len(blockedIPs), len(whitelistedIPs), len(ignoredIPs))
}

func saveIPConfig() {
	// ipMutex is already locked by the caller (handleIPAction)
	// Do NOT lock here to avoid deadlock

	config := IPConfig{
		BlockedIPs:  []IPRecord{},
		Whitelisted: []string{},
		IgnoredIPs:  []string{},
	}

	for ip, reason := range blockedIPs {
		config.BlockedIPs = append(config.BlockedIPs, IPRecord{IP: ip, Reason: reason})
	}
	for ip := range whitelistedIPs {
		config.Whitelisted = append(config.Whitelisted, ip)
	}
	for ip := range ignoredIPs {
		config.IgnoredIPs = append(config.IgnoredIPs, ip)
	}

	data, _ := json.MarshalIndent(config, "", "  ")
	os.WriteFile(BlockedFile, data, 0644)
}

func saveRules() {
	rulesMutex.RLock()
	defer rulesMutex.RUnlock()
	data, _ := json.MarshalIndent(rules, "", "  ")
	os.WriteFile(RulesFile, data, 0644)
}

func loadFilters() {
	file, err := os.Open(FiltersFile)
	if err != nil {
		fmt.Println("[INIT] No existing filters config found, array will be empty.")
		filterConfig = FilterConfig{Keywords: []string{}, WhitelistedURLs: []string{}}
		return
	}
	defer file.Close()
	byteValue, _ := io.ReadAll(file)
	json.Unmarshal(byteValue, &filterConfig)
	fmt.Printf("[INIT] Loaded %d Keywords, %d URL Whitelists\n", len(filterConfig.Keywords), len(filterConfig.WhitelistedURLs))
}
func saveFilters() {
	data, _ := json.MarshalIndent(filterConfig, "", "  ")
	os.WriteFile(FiltersFile, data, 0644)
}

func appendStringIfMissing(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i { return slice }
	}
	return append(slice, i)
}
func removeString(slice []string, s string) []string {
	var result []string
	for _, v := range slice {
		if v != s { result = append(result, v) }
	}
	return result
}

// --- Rate Limiting ---
type IPRateLimit struct {
	Count        int
	LastReset    time.Time
	BlockedUntil time.Time
}

var (
	rateLimits = make(map[string]*IPRateLimit)
	rateMutex  sync.Mutex
)

func isRateLimited(ip string) bool {
	rateMutex.Lock()
	defer rateMutex.Unlock()

	now := time.Now()
	limiter, exists := rateLimits[ip]
	if !exists {
		rateLimits[ip] = &IPRateLimit{
			Count:     1,
			LastReset: now,
		}
		return false
	}

	// Check if already temporarily blocked
	if now.Before(limiter.BlockedUntil) {
		return true
	}

	// If 1 second has passed, reset count
	if now.Sub(limiter.LastReset) > time.Second {
		limiter.Count = 1
		limiter.LastReset = now
		return false
	}

	limiter.Count++
	if limiter.Count > 50 {
		// Temporary block for 10 seconds
		limiter.BlockedUntil = now.Add(10 * time.Second)
		return true
	}

	return false
}

func autoBanIP(ip, reason string) {
	ipMutex.Lock()
	defer ipMutex.Unlock()
	if _, exists := blockedIPs[ip]; !exists {
		blockedIPs[ip] = "Auto-Banned: " + reason
		saveIPConfig()
	}
}

// --- WAF Middleware ---
func WAFMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)

		ipMutex.RLock()
		isWhite := whitelistedIPs[ip]
		banReason, isBan := blockedIPs[ip]
		ipMutex.RUnlock()

		if isWhite {
			next.ServeHTTP(w, r)
			return
		}

		if isBan {
			reasonStr := "IP is manually banned"
			if banReason != "" {
				reasonStr += " (" + banReason + ")"
			}
			logRequest(r, "Blocked", reasonStr, "", "N/A")
			http.Error(w, "403 Forbidden - IP Banned manually", http.StatusForbidden)
			return
		}

		// Rate Limiting Check
		if isRateLimited(ip) {
			logRequest(r, "Blocked", "Rate Limit Exceeded (>50 req/sec)", "RATELIMIT", "Rate Limiter Block")
			fmt.Printf("[WAF] RATE LIMIT BLOCKED %s\n", ip)
			autoBanIP(ip, "Rate Limit Exceeded (>50 req/sec)")
			http.Error(w, "429 Too Many Requests - Rate Limit Exceeded", http.StatusTooManyRequests)
			return
		}

		filterMutex.RLock()
		localFilters := filterConfig
		filterMutex.RUnlock()

		// 1. URL Whitelist Check
		for _, wURL := range localFilters.WhitelistedURLs {
			if strings.HasPrefix(strings.ToLower(r.URL.Path), strings.ToLower(wURL)) {
				// Bypass WAF rules
				logRequest(r, "Allowed", "Bypassed via URL Whitelist: "+wURL, "", "N/A")
				next.ServeHTTP(w, r)
				return
			}
		}

		decodedQuery := r.URL.RawQuery
		// Mitigation: Double/Triple URL Encode Evasion
		for i := 0; i < 3; i++ {
			if unescaped, err := url.QueryUnescape(decodedQuery); err == nil && unescaped != decodedQuery {
				decodedQuery = unescaped
			} else {
				break
			}
		}

		// Mitigation: Null Byte Injection (%00 / \x00)
		decodedQuery = strings.ReplaceAll(decodedQuery, "\x00", "")

		payload := r.URL.Path + " " + decodedQuery
		lowerPayload := strings.ToLower(payload)

		// 2. Keyword Filtering Check
		for _, kw := range localFilters.Keywords {
			if strings.Contains(lowerPayload, strings.ToLower(kw)) {
				logRequest(r, "Blocked", "Blocked by Keyword: "+kw, "KW-BLOCK", "Keyword Exact Match")
				fmt.Printf("[WAF] KEYWORD BLOCKED '%s' from %s\n", kw, ip)
				autoBanIP(ip, "Malicious Keyword: "+kw)
				http.Error(w, "403 Forbidden - WAF Keyword Match", http.StatusForbidden)
				return
			}
		}

		rulesMutex.RLock()
		localRules := rules
		rulesMutex.RUnlock()

		for _, rule := range localRules {
			if !rule.Enabled {
				continue
			}
			match, _ := regexp.MatchString(rule.Pattern, payload)
			if match {
				logRequest(r, "Blocked", "Matched Rule: "+rule.Description, rule.ID, rule.Description)
				fmt.Printf("[WAF] BLOCKED %s from %s\n", rule.Description, ip)
				autoBanIP(ip, "Rule Violation: "+rule.ID)
				http.Error(w, "403 Forbidden - WAF Blocked Request", http.StatusForbidden)
				return
			}
		}

		// Log Allowed Requests for Visibility
		logRequest(r, "Allowed", "Passed WAF Checks", "", "N/A")
		// fmt.Printf("[WAF] ALLOWED request from %s\n", ip)

		next.ServeHTTP(w, r)
	})
}

// --- Helpers ---
func getClientIP(r *http.Request) string {
	// 1. Simulation Header (Highest Priority for Testing)
	if simIP := r.Header.Get("x-simulated-ip"); simIP != "" {
		return simIP
	}

	// 2. Load Balancer / Proxy Header
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return strings.Split(forwarded, ",")[0]
	}

	// 3. RemoteAddr (Handle IPv6 [::1]:port properly)
	ip := r.RemoteAddr
	if strings.Contains(ip, "[") && strings.Contains(ip, "]") {
		// Extract content between brackets for IPv6
		// format: [::1]:8080 or [::1]
		parts := strings.Split(ip, "]")
		if len(parts) > 0 {
			ip = strings.TrimPrefix(parts[0], "[")
		}
	} else if strings.Contains(ip, ":") {
		// IPv4:port
		parts := strings.Split(ip, ":")
		ip = parts[0]
	}

	return ip
}

func getGeoLocation(ip string) string {
	// Mock func for immediate logs, detailed check happens via API
	return "UNK"
}

var logIDCtr int64

func logRequest(r *http.Request, status, details, ruleID, desc string) {
	ip := getClientIP(r)

	// FILTER 1: Muted IPs (Configured by User)
	ipMutex.RLock()
	isIgnored := ignoredIPs[ip]
	ipMutex.RUnlock()
	if isIgnored {
		return
	}

	// FILTER 2: Backend internal calls (Noise Reduction)
	if strings.HasPrefix(r.URL.Path, "/api/stats") ||
		strings.HasPrefix(r.URL.Path, "/api/geo/") ||
		strings.HasPrefix(r.URL.Path, "/api/rules") ||
		strings.HasPrefix(r.URL.Path, "/api/ip") {
		return
	}

	fmt.Printf("[WAF-LOG] Processing log for IP: %s, URL: %s, Status: %s\n", ip, r.URL.Path, status)

	// Unique ID Generation: Timestamp + Atomic Counter
	// Prevents collision during high-concurrency simulation
	uniqueID := fmt.Sprintf("%d-%d", time.Now().UnixNano(), atomic.AddInt64(&logIDCtr, 1))

	entry := WAFLog{
		ID:        uniqueID,
		Timestamp: time.Now(),
		IP:        ip,
		Method:    r.Method,
		URL:       r.URL.String(),
		Status:    status,
		Details:   details,
		RuleID:    ruleID,
		Geo:       getGeoLocation(ip),
	}

	logsMutex.Lock()
	wafLogs = append(wafLogs, entry)
	logsMutex.Unlock()

	go func() {
		logsMutex.RLock()
		defer logsMutex.RUnlock()
		data, _ := json.MarshalIndent(wafLogs, "", "  ")
		os.WriteFile(LogFile, data, 0644)
	}()
}

// --- API Handlers ---
func handleStats(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	logsMutex.RLock()
	defer logsMutex.RUnlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total_requests": len(wafLogs),
		"logs":           wafLogs,
	})
}

func handleRules(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodGet {
		rulesMutex.RLock()
		defer rulesMutex.RUnlock()
		json.NewEncoder(w).Encode(rules)
		return
	}

	if r.Method == http.MethodPost {
		var newRule WAFRule
		if err := json.NewDecoder(r.Body).Decode(&newRule); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		rulesMutex.Lock()
		// Auto ID
		newRule.ID = fmt.Sprintf("%d", len(rules)+1000)
		newRule.Enabled = true // Default to enabled
		rules = append(rules, newRule)
		rulesMutex.Unlock()

		saveRules()
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodPut {
		var req struct {
			ID      string `json:"id"`
			Enabled bool   `json:"enabled"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		rulesMutex.Lock()
		found := false
		for i := range rules {
			if rules[i].ID == req.ID {
				rules[i].Enabled = req.Enabled
				found = true
				break
			}
		}
		rulesMutex.Unlock()

		if found {
			saveRules()
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, "Rule not found", http.StatusNotFound)
		}
		return
	}
}

func handleIPAction(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions {
		return
	}

	if r.Method == http.MethodGet {
		ipMutex.RLock()
		defer ipMutex.RUnlock()

		type IPEntry struct {
			IP     string `json:"ip"`
			Status string `json:"status"`
			Reason string `json:"reason,omitempty"`
		}
		list := []IPEntry{} // Initialize as empty slice for valid JSON []

		for ip, reason := range blockedIPs {
			list = append(list, IPEntry{IP: ip, Status: "Banned", Reason: reason})
		}
		for ip := range whitelistedIPs {
			list = append(list, IPEntry{IP: ip, Status: "Whitelisted"})
		}
		for ip := range ignoredIPs {
			list = append(list, IPEntry{IP: ip, Status: "Ignored"})
		}

		json.NewEncoder(w).Encode(list)
		return
	}

	if r.Method == http.MethodPost {
		var req struct {
			IP     string `json:"ip"`
			Action string `json:"action"`
			Reason string `json:"reason,omitempty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		ipMutex.Lock()
		defer ipMutex.Unlock()

		// Reset status
		fmt.Printf("[WAF-IP] Resetting IP: %s (Action: %s)\n", req.IP, req.Action)
		delete(blockedIPs, req.IP)
		delete(whitelistedIPs, req.IP)
		delete(ignoredIPs, req.IP)

		switch req.Action {
		case "ban":
			blockedIPs[req.IP] = req.Reason
		case "whitelist":
			whitelistedIPs[req.IP] = true
		case "ignore":
			ignoredIPs[req.IP] = true
		}

		saveIPConfig() // Persist changes immediately
		w.WriteHeader(http.StatusOK)
	}
}

func handleFiltersAPI(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	if r.Method == http.MethodOptions { return }

	if r.Method == http.MethodGet {
		filterMutex.RLock()
		defer filterMutex.RUnlock()
		json.NewEncoder(w).Encode(filterConfig)
		return
	}

	if r.Method == http.MethodPost {
		var req struct {
			Type   string `json:"type"`   // "keyword" or "url"
			Value  string `json:"value"`
			Action string `json:"action"` // "add" or "remove"
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		filterMutex.Lock()
		defer filterMutex.Unlock()

		if req.Type == "keyword" {
			if req.Action == "add" {
				filterConfig.Keywords = appendStringIfMissing(filterConfig.Keywords, req.Value)
				filterConfig.WhitelistedURLs = removeString(filterConfig.WhitelistedURLs, req.Value) // Prevent conflict
			} else if req.Action == "remove" {
				filterConfig.Keywords = removeString(filterConfig.Keywords, req.Value)
			}
		} else if req.Type == "url" {
			if req.Action == "add" {
				filterConfig.WhitelistedURLs = appendStringIfMissing(filterConfig.WhitelistedURLs, req.Value)
				filterConfig.Keywords = removeString(filterConfig.Keywords, req.Value) // Prevent conflict
			} else if req.Action == "remove" {
				filterConfig.WhitelistedURLs = removeString(filterConfig.WhitelistedURLs, req.Value)
			}
		}

		saveFilters()
		w.WriteHeader(http.StatusOK)
	}
}

// REAL Proxy to VirusTotal
func handleVirusTotal(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	ip := strings.TrimPrefix(r.URL.Path, "/api/analyze/")
	apiKey := r.Header.Get("X-VT-Key") // User must provide key

	if apiKey == "" {
		apiKey = "5e28d0fd12d4b881c0f32993e0d44e51997fbb16bf02cb9908294c5f833f9cc7" // Demo Key Fallback
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://www.virustotal.com/api/v3/ip_addresses/"+ip, nil)
	req.Header.Set("x-apikey", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "VT API Error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	io.Copy(w, resp.Body)
}

// REAL Proxy to Geo2IP
func handleGeoIP(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)
	ip := strings.TrimPrefix(r.URL.Path, "/api/geo/")
	apiKey := r.URL.Query().Get("key")

	if apiKey == "" {
		apiKey = "YOUR_FREE_KEY"
	}

	resp, err := http.Get(fmt.Sprintf("https://api.ip2geoapi.com/ip/%s?key=%s&format=json", ip, apiKey))
	if err != nil {
		http.Error(w, "Geo API Error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	io.Copy(w, resp.Body)
}

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-VT-Key")
}

func vulnerableHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	w.Write([]byte(fmt.Sprintf("<h1>Flower Shop Search (Go Version)</h1><p>Results for: %s</p>", query)))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", vulnerableHandler)
	mux.HandleFunc("/api/stats", handleStats)
	mux.HandleFunc("/api/rules", handleRules)
	mux.HandleFunc("/api/filters", handleFiltersAPI)
	mux.HandleFunc("/api/ip", handleIPAction)
	mux.HandleFunc("/api/analyze/", handleVirusTotal) // Now a Proxy
	mux.HandleFunc("/api/geo/", handleGeoIP)          // New Proxy

	fmt.Printf("Starting Pro WAF Server on %s...\n", ServerPort)
	if err := http.ListenAndServe(ServerPort, WAFMiddleware(mux)); err != nil {
		log.Fatal(err)
	}
}
