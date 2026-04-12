# WAF-Golang Project Context
*This document provides a comprehensive overview of the "WAF Sentinel" project for AI assistants lacking direct directory access.*

## 1. Project Overview
The "WAF Sentinel" is a custom, high-performance Web Application Firewall (WAF) and real-time visualization dashboard designed for a small-to-medium e-commerce platform (like a flower shop). It acts as an active defense shield that intercepts inbound HTTP traffic, detects malicious signatures, blocks threats, and visualizes the attack data.

## 2. Tech Stack Overview
### Backend (The WAF Middleware & API)
- **Language**: Go (Golang 1.20+)
- **Server/Router**: Native `net/http` ServeMux.
- **Persistence Layer**: JSON files (No external database like PostgreSQL/GORM).
- **Core Security Engine**: Go's native standard Regular Expression (`regexp`) engine.
- **Concurrency**: Goroutines and `sync.RWMutex` for lock-safe memory operations on high-throughput traffic logs.
- **External APIs**: Integrated with VirusTotal (Reputation Scanning) and GeoIP APIs.

### Frontend (The Security Dashboard)
- **Framework**: Next.js (React) using the App Router (`app/` directory).
- **Styling**: Tailwind CSS + Shadcn/UI components (Dark Mode premium aesthetic).
- **Data Visualization**: 
  - `Recharts` for SIEM-style traffic analysis and zero-filled time-series bucket charts (up to 500 bucket performance limits).
  - `Leaflet.js` (via react-leaflet) for dynamic geospatial threat mapping.

## 3. Directory Structure
```text
D:\Antigravity-Workspace\web-dummy\WAF-Golang\
│
├── main.go                 # Core Go Backend (WAF Middleware, APIs, Log Management)
├── waf-logs.json           # Stores historical traffic logs
├── waf-rules.json          # Stores regex-based WAF signatures
├── waf-blocked.json        # Stores Banned, Whitelisted, and Ignored IP configurations
├── simulate_attacks.js     # Node.js script to simulate botnets/traffic
├── CHANGELOG_WAF_GOLANG.md # Detailed technical history of the codebase
│
└── ui/                     # Next.js Frontend Dashboard
    ├── package.json
    ├── app/
    │   ├── layout.tsx      # Global layout (forced "dark" class)
    │   ├── page.tsx        # Main Dashboard UI (Charts, Map, Log Table, API calls)
    │   └── config/
    │       └── page.tsx    # Configuration UI (Rule management, IP Blacklisting)
    └── components/
        ├── TrafficChart.tsx # Recharts implementation (Zero-filling, responsive intervals)
        ├── AttackMap.tsx    # Leaflet implementation (Threat geographic plotting)
        └── ui/              # Shadcn/UI primitive components
```

## 4. Core Mechanics & Architecture

### A. Traffic Interception (WAF Logic)
Located in `WAFMiddleware(next http.Handler)` in `main.go`.
1. Extracts Client IP (handling IPv4 and IPv6 proxy formats safely).
2. Checks memory cache for Whitelist/Blacklist.
3. If not whitelisted, the HTTP request URL and Query are matched against expressions loaded from `waf-rules.json`.
4. If a match is found (e.g., SQL Injection), the WAF returns `403 Forbidden` and logs the attack.
5. If clean, the request is passed through the middleware.

### B. Persistent State Management
Data relies on high-speed JSON read/writes controlled by Mutex locks to prevent race conditions:
- `waf-rules.json`: Dynamic rules array.
- `waf-logs.json`: Append-only array of traffic events.
- `waf-blocked.json`: Three lists -> `blocked_ips`, `whitelisted_ips`, `ignored_ips` (Muted IPs).

### C. The Dashboard (Frontend App)
The `ui/app/page.tsx` serves as the primary hub:
- **State Initialization**: Contains an `isReady` synchronization guard to prevent browser `localStorage` race conditions regarding user-preferred chart intervals (e.g., "1h", "24h").
- **Real-Time Polling**: Fetches `/api/stats` continuously at adjustable intervals (1s, 2s, 5s).
- **Smart Filtering**: Extracts unique IPs from incoming logs and dynamically proxies VirusTotal and GeoIP checks.

## 5. Backend REST API Endpoints (`main.go`)
- `GET /` : Vulnerable target application endpoint (Simulating a Flower Shop).
- `GET /api/stats` : Returns total requests and the full array of WAF logs.
- `GET /api/rules` : Fetches the current WAF RegEx ruleset.
- `POST /api/rules` : Adds a new security rule to the configuration.
- `GET /api/ip` : Returns the combined status of all managed IPs (Banned, Whitelisted, Ignored).
- `POST /api/ip` : Action handler receiving `{ "ip": "...", "action": "ban|whitelist|ignore" }`.
- `GET /api/analyze/{ip}` : Backend Proxy handling VirusTotal requests to hide API keys from the client.
- `GET /api/geo/{ip}` : Backend Proxy handling GeoIP resolution.

## 6. How to Run Locally
The project requires running two persistent terminals:
1. **Backend**: `cd WAF-Golang` -> `go run main.go` (Runs on `localhost:8080`)
2. **Frontend**: `cd WAF-Golang/ui` -> `npm run dev` (Runs on `localhost:3000`)

## 7. Known Nuances / Constraints
- **Performance Cap**: Recharts module has an adaptive bucket limit (~500 buckets) to prevent browser crashing when viewing large 30-day time ranges.
- **Frontend Race Conditions**: The `page.tsx` heavily relies on caching API responses in `useRef` (e.g., `geoCache`) to avoid exhausting rate limits on free Geo/VirusTotal APIs.
- **API Keys**: Keys are stored/injected via the Backend Proxy (`main.go`) to prevent leakage in the frontend source.

*This concludes the context setup for external AI processing.*
