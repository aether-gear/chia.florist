import { useState, useEffect, useRef, useMemo } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { ShieldAlert, ShieldX, Activity, CheckCircle2, ShieldCheck, Trash2, Eye, FolderOpen, ListFilter } from 'lucide-react';
import { ResponsiveContainer, AreaChart, Area, XAxis, YAxis, Tooltip, CartesianGrid, ReferenceArea } from 'recharts';
import { Sheet, SheetContent, SheetDescription, SheetHeader, SheetTitle } from '@/components/ui/sheet';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription } from '@/components/ui/dialog';
import { fetchApi } from '../../lib/api';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Switch } from '@/components/ui/switch';
import { Checkbox } from '@/components/ui/checkbox';

const isPublicIp = (ip: string): boolean => {
  if (!ip) return false;
  const cleanIp = ip.trim().replace(/^::ffff:/, '');
  if (cleanIp === '127.0.0.1' || cleanIp === '::1' || cleanIp === 'localhost') return false;
  if (cleanIp.startsWith('10.')) return false;
  if (cleanIp.startsWith('192.168.')) return false;
  if (cleanIp.startsWith('169.254.')) return false;
  const parts = cleanIp.split('.').map(Number);
  if (parts.length === 4) {
    if (parts[0] === 172 && parts[1] >= 16 && parts[1] <= 31) return false;
  }
  return true;
};

interface GeoPoint {
  ip: string;
  lat: number;
  lng: number;
  city: string;
  country: string;
  isBlocked: boolean;
}

interface ThreatDataPoint {
  time: number;
  blocked: number;
  allowed: number;
}

const LEAFLET_MAP_HTML = `
<!DOCTYPE html>
<html>
<head>
  <link rel="stylesheet" href="https://unpkg.com/leaflet@1.9.4/dist/leaflet.css" />
  <script src="https://unpkg.com/leaflet@1.9.4/dist/leaflet.js"></script>
  <style>
    html, body, #map {
      margin: 0; padding: 0; width: 100%; height: 100%;
      background: #0f172a;
    }
    .leaflet-popup-content-wrapper {
      background: rgba(15, 23, 42, 0.95) !important;
      color: #f8fafc !important;
      border: 1px solid #334155 !important;
      backdrop-filter: blur(4px);
    }
    .leaflet-popup-tip {
      background: rgba(15, 23, 42, 0.95) !important;
      border: 1px solid #334155 !important;
    }
    body.light .leaflet-popup-content-wrapper {
      background: rgba(255, 255, 255, 0.95) !important;
      color: #0f172a !important;
      border: 1px solid #e2e8f0 !important;
    }
    body.light .leaflet-popup-tip {
      background: rgba(255, 255, 255, 0.95) !important;
      border: 1px solid #e2e8f0 !important;
    }
  </style>
</head>
<body>
  <div id="map"></div>
  <script>
    let map = L.map('map', {
      zoomControl: false,
      attributionControl: false,
      minZoom: 2,
      maxBounds: [[-90, -180], [90, 180]],
      maxBoundsViscosity: 1.0
    }).setView([20, 0], 2);

    let currentTileLayer = null;
    let activeTheme = null;

    function setTileLayer(isDark) {
      if (activeTheme === isDark) return;
      activeTheme = isDark;
      if (currentTileLayer) {
        map.removeLayer(currentTileLayer);
      }
      if (isDark) {
        document.body.classList.remove('light');
        currentTileLayer = L.tileLayer('https://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}{r}.png', {
          attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors &copy; <a href="https://carto.com/attributions">CARTO</a>'
        }).addTo(map);
      } else {
        document.body.classList.add('light');
        currentTileLayer = L.tileLayer('https://{s}.basemaps.cartocdn.com/rastertiles/voyager/{z}/{x}/{y}{r}.png', {
          attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors &copy; <a href="https://carto.com/attributions">CARTO</a>'
        }).addTo(map);
      }
    }

    setTileLayer(true);

    let markers = [];
    let hasFittedBounds = false;
    let lastPointsStr = "";

    window.addEventListener('message', function(event) {
      const data = event.data;
      if (data.type === 'setTheme') {
        setTileLayer(data.isDark);
      } else if (data.type === 'setPoints') {
        const points = data.points || [];
        const pointsStr = JSON.stringify(points);
        if (lastPointsStr === pointsStr) {
          return; // Skip redrawing if points list has not changed
        }
        lastPointsStr = pointsStr;

        markers.forEach(m => map.removeLayer(m));
        markers = [];

        points.forEach(pt => {
          if (!pt.lat || !pt.lng) return;
          
          let marker = L.circleMarker([pt.lat, pt.lng], {
            radius: 6,
            fillColor: pt.isBlocked ? '#ef4444' : '#10b981',
            color: '#ffffff',
            weight: 1.5,
            opacity: 1,
            fillOpacity: 0.8
          }).addTo(map);

          let tooltipHtml = '<div style="font-family: sans-serif; font-size: 11px; line-height: 1.4; padding: 2px;">' +
            '<strong style="color: ' + (pt.isBlocked ? '#ef4444' : '#10b981') + '; font-size: 12px;">' + pt.ip + '</strong><br/>' +
            '<strong>Location:</strong> ' + (pt.city ? pt.city + ', ' : '') + (pt.country || 'Unknown') + '<br/>' +
            '<strong>Status:</strong> ' + (pt.isBlocked ? 'Blocked Threat' : 'Allowed Request') +
            '</div>';

          marker.bindTooltip(tooltipHtml, {
            permanent: false,
            direction: 'top',
            className: 'custom-tooltip'
          });

          markers.push(marker);
        });

        if (markers.length > 0 && !hasFittedBounds) {
          let group = new L.featureGroup(markers);
          map.fitBounds(group.getBounds().pad(0.15));
          hasFittedBounds = true;
        }
      }
    });
  </script>
</body>
</html>
`;


const formatDateTimeLocal = (timestamp: number) => {
  const d = new Date(timestamp);
  const pad = (n: number) => String(n).padStart(2, '0');
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`;
};

const parseSafeDate = (dateStr: any): Date => {
  if (!dateStr) return new Date();
  if (dateStr instanceof Date) return dateStr;
  const match = String(dateStr).match(/^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3})\d+([+-]\d{2}:\d{2}|Z)$/);
  if (match) {
    return new Date(match[1] + match[2]);
  }
  return new Date(dateStr);
};

const getIPCategory = (reason: string): string => {
  if (!reason) return 'Manual / Other';
  const r = reason.toLowerCase();
  if (r.includes('rate limit') || r.includes('brute force') || r.includes('spam') || r.includes('brute-force')) {
    return 'Rate Limiting';
  }
  if (r.includes('sqli') || r.includes('sql') || r.includes('union') || r.includes('select')) {
    return 'SQL Injection';
  }
  if (r.includes('xss') || r.includes('script')) {
    return 'XSS Attacks';
  }
  if (r.includes('lfi') || r.includes('path') || r.includes('passwd') || r.includes('win.ini') || r.includes('traversal')) {
    return 'Path Traversal (LFI)';
  }
  if (r.includes('rce') || r.includes('command') || r.includes('shellshock') || r.includes('jndi') || r.includes('ls -la') || r.includes('decode')) {
    return 'RCE & Command Injection';
  }
  if (r.includes('scanner') || r.includes('sqlmap') || r.includes('nikto') || r.includes('user-agent')) {
    return 'Malicious User-Agent';
  }
  if (r.includes('probe') || r.includes('admin') || r.includes('config') || r.includes('setup')) {
    return 'Recon & Probing';
  }
  return 'Manual / Other';
};

export default function SecurityPage() {
  const [isDark, setIsDark] = useState<boolean>(() => document.documentElement.classList.contains('dark'));

  useEffect(() => {
    const observer = new MutationObserver(() => {
      setIsDark(document.documentElement.classList.contains('dark'));
    });
    observer.observe(document.documentElement, { attributes: true, attributeFilter: ['class'] });
    return () => observer.disconnect();
  }, []);

  const [logs, setLogs] = useState<any[]>([]);
  const [wafSummary, setWafSummary] = useState({ total: 0, blocked: 0, allowed: 0, threatLevel: 'Low' });
  const [threatData, setThreatData] = useState<any[]>([]);
  const [isConfirmClearLogsOpen, setIsConfirmClearLogsOpen] = useState(false);

  const [rules, setRules] = useState<any[]>([]);
  const [ipList, setIpList] = useState<any[]>([]);
  const [targetIP, setTargetIP] = useState("");
  const [targetReason, setTargetReason] = useState("");
  const [selectedLog, setSelectedLog] = useState<any | null>(null);
  const [ipRowsPerPage, setIpRowsPerPage] = useState<number>(5);
  const [ipCurrentPage, setIpCurrentPage] = useState<number>(1);
  const [selectedDetailIP, setSelectedDetailIP] = useState<string | null>(null);
  const [isCategoryModalOpen, setIsCategoryModalOpen] = useState(false);
  const [selectedAnalysisCategory, setSelectedAnalysisCategory] = useState("Rate Limiting");
  const [selectedRuleForDetail, setSelectedRuleForDetail] = useState<any | null>(null);
  const [isEditingRule, setIsEditingRule] = useState(false);
  const [editRuleDesc, setEditRuleDesc] = useState("");
  const [editRulePattern, setEditRulePattern] = useState("");
  const [editRuleTags, setEditRuleTags] = useState("");
  const [editRuleImpact, setEditRuleImpact] = useState("5");

  // VirusTotal API Key Configuration (Paste your API key here)
  const VIRUSTOTAL_API_KEY = "5e28d0fd12d4b881c0f32993e0d44e51997fbb16bf02cb9908294c5f833f9cc7"; // PASTE YOUR VIRUSTOTAL API KEY HERE

  // IP2Location API Key Configuration (Paste your API key here)
  const IP2LOCATION_API_KEY = "863EE843FB581979DB85BE72BE0CFD14";

  // VirusTotal Reputation Check State
  const [vtLoading, setVtLoading] = useState<boolean>(false);
  const [vtResult, setVtResult] = useState<any>(null);
  const [vtError, setVtError] = useState<string | null>(null);

  const [geoPoints, setGeoPoints] = useState<GeoPoint[]>([]);
  const iframeRef = useRef<HTMLIFrameElement>(null);

  // States for editing IP details from the popup profile
  const [editIPStatus, setEditIPStatus] = useState<string>('ban');
  const [editIPReason, setEditIPReason] = useState<string>('');

  const prevSelectedDetailIPRef = useRef<string | null>(null);

  useEffect(() => {
    if (selectedDetailIP && selectedDetailIP !== prevSelectedDetailIPRef.current) {
      const entry = ipList.find(e => e.ip === selectedDetailIP);
      if (entry) {
        let mappedStatus = 'ban';
        const s = String(entry.status).toLowerCase();
        if (s === 'ignored' || s === 'ignore') {
          mappedStatus = 'ignore';
        } else if (s === 'whitelisted' || s === 'whitelist') {
          mappedStatus = 'whitelist';
        } else if (s === 'banned_muted') {
          mappedStatus = 'banned_muted';
        } else if (s === 'whitelisted_muted') {
          mappedStatus = 'whitelisted_muted';
        }
        setTimeout(() => {
          setEditIPStatus(mappedStatus);
          setEditIPReason(entry.reason || '');
        }, 0);
      } else {
        setTimeout(() => {
          setEditIPStatus('ban');
          setEditIPReason('');
        }, 0);
      }
    }
    prevSelectedDetailIPRef.current = selectedDetailIP;
  }, [selectedDetailIP, ipList]);

  useEffect(() => {
    setTimeout(() => {
      setVtResult(null);
      setVtError(null);
    }, 0);
  }, [selectedLog]);

  // Real-time IP Geolocation Loader Effect
  useEffect(() => {
    if (logs.length === 0) return;

    const publicIps = Array.from(new Set(logs.map(l => l.ip).filter(isPublicIp)));
    if (publicIps.length === 0) {
      setTimeout(() => {
        setGeoPoints([]);
      }, 0);
      return;
    }

    const cacheKey = "waf_geo_cache";
    let geoCache: Record<string, { lat: number; lng: number; city: string; country: string }> = {};
    try {
      const cached = localStorage.getItem(cacheKey);
      if (cached) geoCache = JSON.parse(cached);
    } catch (err) {
      console.error("Failed to parse geo cache", err);
    }

    // Identify which IPs we need to resolve
    const toFetch = publicIps.filter(ip => !geoCache[ip]);

    // Construct the initial list from cache
    const pointsList = publicIps
      .map(ip => {
        const cached = geoCache[ip];
        if (!cached || !cached.lat || !cached.lng) return null;
        return {
          ip,
          lat: cached.lat,
          lng: cached.lng,
          city: cached.city,
          country: cached.country,
          isBlocked: logs.some(l => l.ip === ip && l.status === 'Blocked')
        };
      })
      .filter(Boolean) as GeoPoint[];

    setTimeout(() => {
      setGeoPoints(pointsList);
    }, 0);

    if (toFetch.length === 0) return;

    const fetchGeoData = async () => {
      let updated = false;
      for (const ip of toFetch) {
        try {
          const res = await fetch(`http://localhost:8080/api/geo/${ip}?key=${IP2LOCATION_API_KEY}`);
          if (!res.ok) continue;
          const data = await res.json();
          if (data.latitude && data.longitude) {
            geoCache[ip] = {
              lat: data.latitude,
              lng: data.longitude,
              city: data.city_name || "",
              country: data.country_name || ""
            };
            updated = true;
          }
        } catch (err) {
          console.error(`Failed to geolocate ${ip}`, err);
        }
      }

      if (updated) {
        localStorage.setItem(cacheKey, JSON.stringify(geoCache));
        const freshPoints = publicIps
          .map(ip => {
            const cached = geoCache[ip];
            if (!cached || !cached.lat || !cached.lng) return null;
            return {
              ip,
              lat: cached.lat,
              lng: cached.lng,
              city: cached.city,
              country: cached.country,
              isBlocked: logs.some(l => l.ip === ip && l.status === 'Blocked')
            };
          })
          .filter(Boolean) as GeoPoint[];
        setTimeout(() => {
          setGeoPoints(freshPoints);
        }, 0);
      }
    };

    fetchGeoData();
  }, [logs]);

  // Post messages to Leaflet iframe
  useEffect(() => {
    if (iframeRef.current && iframeRef.current.contentWindow) {
      iframeRef.current.contentWindow.postMessage({
        type: 'setTheme',
        isDark
      }, '*');

      iframeRef.current.contentWindow.postMessage({
        type: 'setPoints',
        points: geoPoints
      }, '*');
    }
  }, [geoPoints, isDark]);

  // Time Range & Brushing State
  const [rangeType, setRangeType] = useState<"Today" | "Custom">("Today");
  const [customStart, setCustomStart] = useState("");
  const [customEnd, setCustomEnd] = useState("");
  const [tempStart, setTempStart] = useState("");
  const [tempEnd, setTempEnd] = useState("");
  const [selectedLogIds, setSelectedLogIds] = useState<Record<string, boolean>>({});
  const [selectedIPs, setSelectedIPs] = useState<Record<string, boolean>>({});
  const [refAreaLeft, setRefAreaLeft] = useState<number | null>(null);
  const [refAreaRight, setRefAreaRight] = useState<number | null>(null);

  const handleApplyCustomRange = () => {
    setCustomStart(tempStart);
    setCustomEnd(tempEnd);
  };

  const handleSelectAll = (checked: boolean) => {
    const nextSelected: Record<string, boolean> = {};
    if (checked) {
      displayLogs.forEach(log => {
        const logId = log.id || `${log.timestamp}-${log.ip}`;
        nextSelected[logId] = true;
      });
    }
    setSelectedLogIds(nextSelected);
  };

  const handleSelectLog = (logId: string, checked: boolean) => {
    setSelectedLogIds(prev => ({
      ...prev,
      [logId]: checked
    }));
  };

  const handleBulkIPAction = async (action: string) => {
    const selectedIds = Object.keys(selectedLogIds).filter(id => selectedLogIds[id]);
    if (selectedIds.length === 0) return;

    // Retrieve unique IPs of the selected logs
    const selectedLogs = logs.filter(l => {
      const logId = l.id || `${l.timestamp}-${l.ip}`;
      return selectedLogIds[logId];
    });
    const uniqueIPs = Array.from(new Set(selectedLogs.map(l => l.ip).filter(Boolean)));

    if (uniqueIPs.length === 0) return;

    try {
      const mappedAction = action === 'remove' || action === 'reset' ? 'reset' : action;
      await Promise.all(
        uniqueIPs.map(ip =>
          fetchApi('/api/ip', {
            method: 'POST',
            body: JSON.stringify({ ip, action: mappedAction, reason: `Bulk ${action} from Logs` }),
          })
        )
      );
      setSelectedLogIds({});
      fetchIpList();
    } catch (e) {
      console.error("Bulk IP action failed", e);
    }
  };

  const handleBulkDeleteLogs = async () => {
    const selectedIds = Object.keys(selectedLogIds).filter(id => selectedLogIds[id]);
    if (selectedIds.length === 0) return;

    const realIds = selectedIds.filter(id => id.length === 36);

    try {
      if (realIds.length > 0) {
        await fetchApi(`/api/stats?ids=${realIds.join(',')}`, {
          method: 'DELETE',
        });
      }

      setLogs(prev => prev.filter(log => {
        const logId = log.id || `${log.timestamp}-${log.ip}`;
        return !selectedLogIds[logId];
      }));
      setSelectedLogIds({});
      fetchLogs();
    } catch (e) {
      console.error("Failed to bulk delete audit logs", e);
    }
  };

  const handleClearAllLogs = () => {
    setIsConfirmClearLogsOpen(true);
  };

  const executeClearAllLogs = async () => {
    try {
      await fetchApi('/api/stats?all=true', {
        method: 'DELETE',
      });
      setLogs([]);
      setSelectedLogIds({});
      setWafSummary({ total: 0, blocked: 0, allowed: 0, threatLevel: 'Low' });
      fetchLogs();
    } catch (e) {
      console.error("Failed to clear all audit logs", e);
    }
  };
  const handleCheckReputation = async (ip: string) => {
    setVtLoading(true);
    setVtError(null);
    setVtResult(null);
    try {
      const headers: Record<string, string> = {
        "Content-Type": "application/json",
      };
      if (VIRUSTOTAL_API_KEY) {
        headers["X-VT-Key"] = VIRUSTOTAL_API_KEY;
      }

      const vtPromise = fetch(`http://localhost:8080/api/analyze/${ip}`, {
        method: "GET",
        headers,
      }).then(async r => {
        if (!r.ok) {
          const errText = await r.text();
          throw new Error(errText || `VirusTotal returned ${r.status}`);
        }
        return r.json();
      });

      const ip2locPromise = fetch(`http://localhost:8080/api/geo/${ip}?key=${IP2LOCATION_API_KEY}`)
        .then(async r => {
          if (!r.ok) {
            const errText = await r.text();
            throw new Error(errText || `IP2Location returned ${r.status}`);
          }
          return r.json();
        });

      const [vtData, ip2locData] = await Promise.all([vtPromise, ip2locPromise]);

      // Extract stats and info from standard VirusTotal v3 payload
      const attrs = vtData?.data?.attributes;
      if (!attrs) {
        throw new Error("Invalid VirusTotal API response format");
      }

      const stats = attrs.last_analysis_stats || { harmless: 0, malicious: 0, suspicious: 0, undetected: 0 };
      setVtResult({
        stats,
        asOwner: ip2locData.as || attrs.as_owner || "Unknown ISP / AS Owner",
        country: ip2locData.country_name || attrs.country || "UNK",
        countryCode: ip2locData.country_code || "UNK",
        regionName: ip2locData.region_name || "",
        cityName: ip2locData.city_name || "",
        latitude: ip2locData.latitude || 0,
        longitude: ip2locData.longitude || 0,
        network: attrs.network || "N/A",
        isProxy: ip2locData.is_proxy || false,
      });
    } catch (e: any) {
      console.error("IP reputation or location scan failed", e);
      setVtError(e.message || "An unknown error occurred during scanning");
    } finally {
      setVtLoading(false);
    }
  };

  const handleBulkIPActionFromManager = async (action: 'ban' | 'whitelist' | 'ignore' | 'remove') => {
    const ips = Object.keys(selectedIPs).filter(ip => selectedIPs[ip]);
    if (ips.length === 0) return;

    try {
      const mappedAction = action === 'remove' ? 'reset' : action;
      await Promise.all(
        ips.map(ip =>
          fetchApi('/api/ip', {
            method: 'POST',
            body: JSON.stringify({ ip, action: mappedAction, reason: action === 'remove' ? '' : 'Bulk Action' }),
          })
        )
      );
      // Optimistic local UI updates
      if (action === 'remove') {
        setIpList(prev => prev.filter(entry => !selectedIPs[entry.ip]));
      } else {
        setIpList(prev => prev.map(entry => {
          if (selectedIPs[entry.ip]) {
            return {
              ...entry,
              status: action === 'ban' ? 'Banned' : action === 'whitelist' ? 'Whitelisted' : 'Ignored',
              reason: 'Bulk Action'
            };
          }
          return entry;
        }));
      }
      setSelectedIPs({});
      fetchIpList();
    } catch (e) {
      console.error("Bulk manager IP action failed", e);
    }
  };

  const handleClearAllIPs = async () => {
    if (!window.confirm("Are you sure you want to forget and reset ALL configured IPs?")) return;
    try {
      await Promise.all(
        ipList.map(entry =>
          fetchApi('/api/ip', {
            method: 'POST',
            body: JSON.stringify({ ip: entry.ip, action: 'reset', reason: '' }),
          })
        )
      );
      setIpList([]);
      setSelectedIPs({});
      fetchIpList();
    } catch (e) {
      console.error("Failed to clear all IPs", e);
    }
  };

  // Security Logs filter and pagination
  const [rowsPerPage, setRowsPerPage] = useState<number>(10);
  const [statusFilter, setStatusFilter] = useState<string>("All");

  // New Rule Form State
  const [isAddRuleOpen, setIsAddRuleOpen] = useState(false);
  const [newRuleDesc, setNewRuleDesc] = useState("");
  const [newRulePattern, setNewRulePattern] = useState("");
  const [newRuleTags, setNewRuleTags] = useState("");
  const [newRuleImpact, setNewRuleImpact] = useState("5");

  const handleAddRule = async () => {
    if (!newRuleDesc || !newRulePattern) return;
    try {
      await fetchApi('/api/rules', {
        method: 'POST',
        body: JSON.stringify({
          description: newRuleDesc,
          pattern: newRulePattern,
          tags: newRuleTags.split(',').map(t => t.trim()).filter(Boolean),
          impact: newRuleImpact,
          enabled: true
        })
      });
      setIsAddRuleOpen(false);
      setNewRuleDesc("");
      setNewRulePattern("");
      setNewRuleTags("");
      fetchRules();
    } catch (e) {
      console.error("Failed to add rule", e);
    }
  };

  const fetchRules = async () => {
    try {
      const data = await fetchApi('/api/rules?t=' + Date.now());
      setRules(data?.rules || []);
    } catch (e) {
      console.error(e);
    }
  };

  const fetchIpList = async () => {
    try {
      const data = await fetchApi('/api/ip?t=' + Date.now());
      setIpList(data?.entries || []);
    } catch (e) {
      console.error(e);
    }
  };

  const fetchLogs = async () => {
    try {
      const data = await fetchApi('/api/stats?limit=1000&t=' + Date.now());
      const rawLogs = data?.audit_logs || [];
      const mappedLogs = rawLogs.map((log: any) => {
        const isBlocked = log.outcome === 'blocked' || log.action === 'request_blocked';
        return {
          id: log.id,
          timestamp: log.created_at,
          ip: log.client_ip || '127.0.0.1',
          method: log.metadata?.method || 'GET',
          url: log.metadata?.path || '/',
          status: isBlocked ? 'Blocked' : 'Allowed',
          ruleId: log.metadata?.rule_id || '-',
          reason: log.metadata?.reason || '-',
          payload: log.metadata?.payload || '-',
          userAgent: log.metadata?.user_agent || '-'
        };
      });

      const sorted = [...mappedLogs].sort((a, b) => parseSafeDate(b.timestamp).getTime() - parseSafeDate(a.timestamp).getTime());
      setLogs(sorted);

      const total = mappedLogs.length;
      const blocked = mappedLogs.filter((l: any) => l.status === 'Blocked').length;
      setWafSummary({
        total,
        blocked,
        allowed: total - blocked,
        threatLevel: blocked > 50 ? 'High' : (blocked > 10 ? 'Medium' : 'Low')
      });
    } catch (e) {
      console.error(e);
    }
  };

  useEffect(() => {
    fetchRules();
    fetchIpList();
    fetchLogs();

    const interval = setInterval(() => {
      fetchRules();
      fetchIpList();
      fetchLogs();
    }, 3000); // Poll every 3 seconds for live updates
    return () => clearInterval(interval);
  }, []);

  useEffect(() => {
    const dataMap: Record<number, { time: number; blocked: number; allowed: number }> = {};
    let startTime = 0;
    let endTime = Number.MAX_SAFE_INTEGER;
    let step = 1000 * 60 * 15; // Default to 15-minute intervals

    if (rangeType === 'Today') {
      const now = new Date();
      startTime = new Date(now.getFullYear(), now.getMonth(), now.getDate()).getTime();
      endTime = startTime + 24 * 60 * 60 * 1000;
      step = 1000 * 60 * 10; // 10-minute intervals for today (144 points)
    } else if (rangeType === 'Custom') {
      const now = new Date();
      startTime = customStart ? new Date(customStart).getTime() : now.getTime() - 24 * 60 * 60 * 1000;
      endTime = customEnd ? new Date(customEnd).getTime() : now.getTime();

      const diff = endTime - startTime;
      if (diff > 0) {
        // Keep number of points bounded for performance
        step = Math.max(1000 * 60 * 5, Math.floor(diff / 150));
      }
    }

    // Pre-populate dataMap with empty intervals to fill the entire spectrum on the X-axis
    if (startTime > 0 && endTime < Number.MAX_SAFE_INTEGER) {
      for (let t = startTime; t <= endTime; t += step) {
        dataMap[t] = { time: t, blocked: 0, allowed: 0 };
      }
    }

    logs.forEach(log => {
      const ts = parseSafeDate(log.timestamp).getTime();
      if (ts >= startTime && ts <= endTime) {
        // Round to nearest step
        const roundedTs = Math.floor(ts / step) * step;
        if (!dataMap[roundedTs]) {
          dataMap[roundedTs] = { time: roundedTs, blocked: 0, allowed: 0 };
        }
        if (log.status === 'Blocked') {
          dataMap[roundedTs].blocked += 1;
        } else {
          dataMap[roundedTs].allowed += 1;
        }
      }
    });

    const sortedThreatData = Object.values(dataMap).sort((a: ThreatDataPoint, b: ThreatDataPoint) => a.time - b.time);
    setTimeout(() => {
      setThreatData(sortedThreatData);
    }, 0);
  }, [logs, rangeType, customStart, customEnd]);

  const displayLogs = logs.filter(log => {
    const ts = parseSafeDate(log.timestamp).getTime();
    let startTime = 0;
    let endTime = Number.MAX_SAFE_INTEGER;
    if (rangeType === 'Today') {
      const now = new Date();
      startTime = new Date(now.getFullYear(), now.getMonth(), now.getDate()).getTime();
      endTime = startTime + 24 * 60 * 60 * 1000;
    } else if (rangeType === 'Custom') {
      startTime = customStart ? new Date(customStart).getTime() : 0;
      endTime = customEnd ? new Date(customEnd).getTime() : Number.MAX_SAFE_INTEGER;
    }
    const matchesTime = ts >= startTime && ts <= endTime;

    if (!matchesTime) return false;

    if (statusFilter !== "All") {
      if (statusFilter === "Blocked") {
        return log.status === "Blocked";
      } else if (statusFilter === "Passed") {
        return log.status === "Allowed";
      }
    }

    return true;
  });

  const displayedIps = ipList.slice((ipCurrentPage - 1) * ipRowsPerPage, ipCurrentPage * ipRowsPerPage);

  const handleIPSubmit = async (action: string, ipOverride?: string, reasonOverride?: string) => {
    const ip = ipOverride || targetIP;
    const reason = reasonOverride !== undefined ? reasonOverride : targetReason;
    if (!ip) return;

    const mappedAction = action === 'remove' || action === 'reset' ? 'reset' : action;
    await fetchApi('/api/ip', {
      method: 'POST',
      body: JSON.stringify({ ip, action: mappedAction, reason }),
    });

    if (!ipOverride) {
      setTargetIP("");
      setTargetReason("");
    }
    fetchIpList();
  };

  const toggleRule = async (ruleId: string, currentStatus: boolean) => {
    await fetchApi(`/api/rules/${ruleId}`, {
      method: 'PUT',
      body: JSON.stringify({ enabled: !currentStatus }),
    });
    fetchRules();
  };

  const handleDeleteRule = async (ruleId: string) => {
    if (!window.confirm("Are you sure you want to delete this WAF rule?")) return;
    try {
      await fetchApi(`/api/rules/${ruleId}`, {
        method: 'DELETE',
      });
      fetchRules();
    } catch (e) {
      console.error("Failed to delete rule", e);
    }
  };

  const handleUpdateRule = async () => {
    if (!selectedRuleForDetail) return;
    try {
      await fetchApi(`/api/rules/${selectedRuleForDetail.id}`, {
        method: 'PUT',
        body: JSON.stringify({
          description: editRuleDesc,
          pattern: editRulePattern,
          tags: editRuleTags.split(',').map(t => t.trim()).filter(Boolean),
          impact: editRuleImpact,
        })
      });
      setIsEditingRule(false);
      const updated = {
        ...selectedRuleForDetail,
        description: editRuleDesc,
        pattern: editRulePattern,
        tags: editRuleTags.split(',').map(t => t.trim()).filter(Boolean),
        impact: editRuleImpact,
      };
      setSelectedRuleForDetail(updated);
      fetchRules();
    } catch (e) {
      console.error("Failed to update rule", e);
    }
  };

  useEffect(() => {
    if (iframeRef.current) {
      iframeRef.current.srcdoc = LEAFLET_MAP_HTML;
    }
  }, []);

  const mapIframe = useMemo(() => {
    return (
      <iframe
        ref={iframeRef}
        title="IP Geolocation Map"
        width="100%"
        height="100%"
        frameBorder="0"
        scrolling="no"
        marginHeight={0}
        marginWidth={0}
        onLoad={() => {
          if (iframeRef.current && iframeRef.current.contentWindow) {
            iframeRef.current.contentWindow.postMessage({
              type: 'setTheme',
              isDark
            }, '*');
            iframeRef.current.contentWindow.postMessage({
              type: 'setPoints',
              points: geoPoints
            }, '*');
          }
        }}
        className="w-full h-full"
      />
    );
  }, []);

  const now = new Date();
  let activeStartTime = 0;
  let activeEndTime = Number.MAX_SAFE_INTEGER;
  if (rangeType === 'Today') {
    activeStartTime = new Date(now.getFullYear(), now.getMonth(), now.getDate()).getTime();
    activeEndTime = activeStartTime + 24 * 60 * 60 * 1000;
  } else if (rangeType === 'Custom') {
    activeStartTime = customStart ? new Date(customStart).getTime() : now.getTime() - 24 * 60 * 60 * 1000;
    activeEndTime = customEnd ? new Date(customEnd).getTime() : now.getTime();
  }

  return (
    <div className="space-y-10 animate-in fade-in duration-300">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h2 className="text-3xl font-bold font-display tracking-tight text-foreground">Security Monitoring</h2>
          <p className="text-muted-foreground text-sm">Real-time threat analysis and WAF status.</p>
        </div>
      </div>

      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
        <Card className="border-0 shadow-none bg-zinc-50/40 dark:bg-slate-900/40">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-bold font-display text-foreground">Attack Attempts</CardTitle>
            <ShieldAlert className="h-4 w-4 text-rose-500" />
          </CardHeader>
          <CardContent>
            <div className="text-3xl font-bold tracking-tight text-rose-600 dark:text-rose-400">{wafSummary.blocked.toLocaleString()}</div>
            <p className="text-xs text-rose-500 font-medium mt-1">Total requests blocked by rules</p>
          </CardContent>
        </Card>
        <Card className="border-0 shadow-none bg-zinc-50/40 dark:bg-slate-900/40">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-bold font-display text-foreground">Safe Traffic</CardTitle>
            <CheckCircle2 className="h-4 w-4 text-primary" />
          </CardHeader>
          <CardContent>
            <div className="text-3xl font-bold tracking-tight text-primary">{wafSummary.allowed.toLocaleString()}</div>
            <p className="text-xs text-muted-foreground mt-1">Allowed requests</p>
          </CardContent>
        </Card>
        <Card className="border-0 shadow-none bg-zinc-50/40 dark:bg-slate-900/40">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-bold font-display text-foreground">Anomaly Score</CardTitle>
            <Activity className="h-4 w-4 text-orange-500" />
          </CardHeader>
          <CardContent>
            <div className="text-3xl font-bold tracking-tight text-orange-600 dark:text-orange-400">
              {((wafSummary.blocked / wafSummary.total) * 100).toFixed(1)}%
            </div>
            <p className="text-xs text-muted-foreground mt-1">Traffic flagged as malicious</p>
          </CardContent>
        </Card>
        <Card className="border-0 shadow-none bg-zinc-50/40 dark:bg-slate-900/40">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-bold font-display text-foreground">Active Incidents</CardTitle>
            <ShieldX className="h-4 w-4 text-rose-500" />
          </CardHeader>
          <CardContent>
            <div className="text-3xl font-bold tracking-tight text-rose-600 dark:text-rose-400">{ipList.filter(ip => ip.status === 'Banned').length}</div>
            <p className="text-xs text-muted-foreground mt-1">Currently banned IPs</p>
          </CardContent>
        </Card>
      </div>

      <Card className="border-0 shadow-none bg-zinc-50/40 dark:bg-slate-900/40">
        <CardHeader className="flex flex-col sm:flex-row sm:items-center justify-between">
          <div>
            <CardTitle className="font-bold font-display tracking-tight text-lg text-foreground">Threat Landscape</CardTitle>
            <CardDescription className="text-muted-foreground text-sm">Volume of allowed vs blocked requests over time.</CardDescription>
          </div>
          <div className="flex items-center gap-3 mt-4 sm:mt-0">
            <Select
              value={rangeType}
              onValueChange={(val: any) => {
                setRangeType(val);
                if (val === 'Today') {
                  setCustomStart("");
                  setCustomEnd("");
                  setTempStart("");
                  setTempEnd("");
                } else {
                  const now = new Date();
                  const startStr = formatDateTimeLocal(now.getTime() - 24 * 60 * 60 * 1000);
                  const endStr = formatDateTimeLocal(now.getTime());
                  setTempStart(startStr);
                  setTempEnd(endStr);
                  setCustomStart(startStr);
                  setCustomEnd(endStr);
                }
              }}
            >
              <SelectTrigger className="w-[180px] h-9 bg-background text-foreground border border-input">
                <SelectValue placeholder="Select Range" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="Today">Today (00:00 - 24:00)</SelectItem>
                <SelectItem value="Custom">Custom Range</SelectItem>
              </SelectContent>
            </Select>
            {rangeType === 'Custom' && (
              <div className="flex items-center gap-2">
                <Input type="datetime-local" className="h-8 text-xs w-[170px]" value={tempStart} onChange={e => setTempStart(e.target.value)} />
                <span className="text-slate-400 text-sm">to</span>
                <Input type="datetime-local" className="h-8 text-xs w-[170px]" value={tempEnd} onChange={e => setTempEnd(e.target.value)} />
                <Button size="sm" className="h-8 px-3" onClick={handleApplyCustomRange}>
                  Apply
                </Button>
              </div>
            )}
            {rangeType === 'Custom' && (
              <Button size="sm" variant="outline" onClick={() => { setRangeType("Today"); setCustomStart(""); setCustomEnd(""); setTempStart(""); setTempEnd(""); }}>
                Reset Zoom
              </Button>
            )}
          </div>
        </CardHeader>
        <CardContent className="h-[350px]">
          <ResponsiveContainer width="100%" height="100%">
            <AreaChart
              data={threatData}
              margin={{ top: 10, right: 30, left: 0, bottom: 0 }}
              onMouseDown={(e) => { if (e && e.activeLabel) setRefAreaLeft(Number(e.activeLabel)); }}
              onMouseMove={(e) => { if (refAreaLeft !== null && e && e.activeLabel) setRefAreaRight(Number(e.activeLabel)); }}
              onMouseUp={() => {
                if (refAreaLeft !== null && refAreaRight !== null) {
                  let [start, end] = [refAreaLeft, refAreaRight];
                  if (start > end) [start, end] = [end, start];
                  if (end - start > 1000 * 60) { // Zoom only if dragging at least 1 minute
                    setRangeType("Custom");
                    const startStr = formatDateTimeLocal(start);
                    const endStr = formatDateTimeLocal(end);
                    setTempStart(startStr);
                    setTempEnd(endStr);
                    setCustomStart(startStr);
                    setCustomEnd(endStr);
                  }
                }
                setRefAreaLeft(null);
                setRefAreaRight(null);
              }}
              style={{ userSelect: 'none' }}
            >
              <defs>
                <linearGradient id="colorBlocked" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#ef4444" stopOpacity={0.3} />
                  <stop offset="95%" stopColor="#ef4444" stopOpacity={0} />
                </linearGradient>
                <linearGradient id="colorAllowed" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#10b981" stopOpacity={0.1} />
                  <stop offset="95%" stopColor="#10b981" stopOpacity={0} />
                </linearGradient>
              </defs>
              <XAxis
                dataKey="time"
                type="number"
                scale="time"
                domain={[activeStartTime, activeEndTime]}
                axisLine={false}
                tickLine={false}
                tick={{ fill: isDark ? "#94a3b8" : "#64748b", fontSize: 12 }}
                dy={10}
                tickFormatter={(unixTime) => {
                  const d = new Date(unixTime);
                  return `${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`;
                }}
              />
              <YAxis
                axisLine={false}
                tickLine={false}
                tick={{ fill: isDark ? "#94a3b8" : "#64748b", fontSize: 12 }}
                allowDecimals={false}
                padding={{ top: 30 }}
              />
              <CartesianGrid strokeDasharray="3 3" vertical={false} stroke={isDark ? "#334155" : "#e2e8f0"} />
              <Tooltip
                contentStyle={{
                  borderRadius: '8px',
                  border: `1px solid ${isDark ? '#334155' : '#e2e8f0'}`,
                  boxShadow: '0 4px 12px rgba(0,0,0,0.1)',
                  backgroundColor: isDark ? '#0f172a' : '#ffffff',
                  color: isDark ? '#f1f5f9' : '#0f172a'
                }}
                labelFormatter={(unixTime: any) => {
                  const d = new Date(unixTime);
                  return `${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`;
                }}
              />
              <Area type="monotone" dataKey="allowed" name="Allowed Traffic" stroke="#10b981" fillOpacity={1} fill="url(#colorAllowed)" />
              <Area type="monotone" dataKey="blocked" name="Blocked Threats" stroke="#ef4444" fillOpacity={1} fill="url(#colorBlocked)" />

              {refAreaLeft !== null && refAreaRight !== null ? (
                <ReferenceArea x1={refAreaLeft} x2={refAreaRight} strokeOpacity={0.3} fill="#8884d8" />
              ) : null}
            </AreaChart>
          </ResponsiveContainer>
        </CardContent>
      </Card>

      {/* Real-time Threat Geolocation Map */}
      <Card className="border-0 shadow-none bg-zinc-50/40 dark:bg-slate-900/40">
        <CardHeader>
          <CardTitle className="font-bold font-display tracking-tight text-lg text-foreground">
            Real-time Threat Geolocation Map
          </CardTitle>
          <CardDescription className="text-muted-foreground text-sm">
            Pinpoints the physical location of all audited public IP addresses in real-time. Hover over points for threat details.
          </CardDescription>
        </CardHeader>
        <CardContent className="pt-0">
          <div className="h-[400px] flex items-center justify-center bg-background border border-border rounded-2xl overflow-hidden p-0 relative">
            {mapIframe}
            {geoPoints.length === 0 && (
              <div className="absolute inset-0 flex flex-col items-center justify-center bg-background/90 text-center p-6 text-muted-foreground space-y-2 z-20">
                <div className="text-3xl">🌐</div>
                <div className="font-medium text-foreground">No Geolocation Data Available</div>
                <div className="text-xs max-w-sm mx-auto text-muted-foreground">
                  No public IP addresses have been logged yet. Local or private traffic will not appear on the threat map.
                </div>
              </div>
            )}
          </div>
        </CardContent>
      </Card>

      {/* IP Blacklist Manager */}
      <Card className="border-0 shadow-none bg-zinc-50/40 dark:bg-slate-900/40">
        <CardHeader className="flex flex-row items-center justify-between space-y-0">
          <div>
            <CardTitle className="flex items-center gap-2 font-bold font-display tracking-tight text-lg text-foreground"><ShieldCheck className="text-primary" /> IP Blacklist Manager</CardTitle>
            <CardDescription className="text-muted-foreground text-sm">Control access and logging behavior for specific IPs.</CardDescription>
          </div>
          <div className="flex gap-2">
            <Button
              size="sm"
              variant="outline"
              className="h-9 gap-2 border-border text-foreground hover:bg-muted"
              onClick={handleClearAllIPs}
            >
              Reset / Forget All IPs
            </Button>
            <Button
              size="sm"
              variant="outline"
              className="h-9 gap-2 border-primary/20 text-primary hover:bg-primary/5"
              onClick={() => setIsCategoryModalOpen(true)}
            >
              <FolderOpen className="h-4 w-4" /> View Category Analysis
            </Button>
          </div>
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="flex flex-col md:flex-row gap-4 items-end bg-background p-5 rounded-2xl border border-border/80">
            <div className="space-y-2 flex-1 w-full">
              <label className="text-xs font-bold text-muted-foreground uppercase tracking-wider">Target IP Address</label>
              <Input
                placeholder="e.g. 192.168.1.5 or ::1"
                className="bg-background text-foreground border-border rounded-xl font-mono"
                value={targetIP}
                onChange={e => setTargetIP(e.target.value)}
              />
            </div>
            <div className="space-y-2 flex-1 w-full">
              <label className="text-xs font-bold text-muted-foreground uppercase tracking-wider">Reason (Optional)</label>
              <Input
                placeholder="e.g. Spamming API"
                className="bg-background text-foreground border-border rounded-xl"
                value={targetReason}
                onChange={e => setTargetReason(e.target.value)}
              />
            </div>
            <div className="flex gap-2 w-full md:w-auto">
              <Button variant="destructive" className="flex-1 md:flex-none rounded-xl" onClick={() => handleIPSubmit('ban')}>Ban</Button>
              <Button className="flex-1 md:flex-none bg-primary hover:bg-primary/90 text-primary-foreground rounded-xl" onClick={() => handleIPSubmit('whitelist')}>Whitelist</Button>
              <Button variant="secondary" className="flex-1 md:flex-none rounded-xl" onClick={() => handleIPSubmit('ignore')}>Mute</Button>
            </div>
          </div>

          {Object.keys(selectedIPs).filter(ip => selectedIPs[ip]).length > 0 && (
            <div className="flex items-center justify-between bg-muted/40 border border-border/80 rounded-2xl p-4 mb-4 animate-in fade-in slide-in-from-top-2 duration-200">
              <span className="text-xs font-semibold text-muted-foreground">
                {Object.keys(selectedIPs).filter(ip => selectedIPs[ip]).length} IP(s) selected
              </span>
              <div className="flex gap-2">
                <Button
                  size="sm"
                  variant="destructive"
                  className="h-8 text-xs font-semibold rounded-xl"
                  onClick={() => handleBulkIPActionFromManager('ban')}
                >
                  Bulk Ban
                </Button>
                <Button
                  size="sm"
                  className="h-8 text-xs font-semibold bg-primary hover:bg-primary/90 text-primary-foreground rounded-xl"
                  onClick={() => handleBulkIPActionFromManager('whitelist')}
                >
                  Bulk Whitelist
                </Button>
                <Button
                  size="sm"
                  variant="secondary"
                  className="h-8 text-xs"
                  onClick={() => handleBulkIPActionFromManager('ignore')}
                >
                  Bulk Mute
                </Button>
                <Button
                  size="sm"
                  variant="outline"
                  className="h-8 text-xs border-slate-300 dark:border-slate-800 hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-700 dark:text-slate-200 font-medium"
                  onClick={() => handleBulkIPActionFromManager('remove')}
                >
                  Bulk Remove / Forget
                </Button>
              </div>
            </div>
          )}

          <div className="overflow-x-auto border border-border rounded-2xl">
            <Table>
              <TableHeader>
                <TableRow className="bg-muted/50">
                  <TableHead className="w-[40px]">
                    <Checkbox
                      checked={ipList.length > 0 && ipList.every(entry => selectedIPs[entry.ip])}
                      onCheckedChange={(checked) => {
                        const updated: Record<string, boolean> = {};
                        if (checked) {
                          ipList.forEach(entry => {
                            updated[entry.ip] = true;
                          });
                        }
                        setSelectedIPs(updated);
                      }}
                    />
                  </TableHead>
                  <TableHead>IP Address</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Reason</TableHead>
                  <TableHead className="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {ipList.length === 0 && (
                  <TableRow>
                    <TableCell colSpan={5} className="text-center text-slate-500 py-8">No IPs configured yet.</TableCell>
                  </TableRow>
                )}
                {displayedIps.map((entry: any) => (
                  <TableRow key={entry.ip}>
                    <TableCell>
                      <Checkbox
                        checked={!!selectedIPs[entry.ip]}
                        onCheckedChange={(checked) => {
                          setSelectedIPs(prev => ({
                            ...prev,
                            [entry.ip]: !!checked
                          }));
                        }}
                      />
                    </TableCell>
                    <TableCell
                      className="font-mono text-slate-700 hover:text-emerald-600 hover:underline cursor-pointer"
                      onClick={() => setSelectedDetailIP(entry.ip)}
                    >
                      {entry.ip}
                    </TableCell>
                    <TableCell>
                      {(entry.status === 'Banned' || entry.status === 'banned') && <Badge variant="destructive">Banned</Badge>}
                      {(entry.status === 'Whitelisted' || entry.status === 'whitelisted') && <Badge className="bg-emerald-600">Whitelisted</Badge>}
                      {(entry.status === 'Ignored' || entry.status === 'ignored') && <Badge variant="secondary" className="bg-slate-100 dark:bg-slate-800 text-slate-800 dark:text-slate-200 border-0">(No Log)</Badge>}
                      {(entry.status === 'banned_muted') && (
                        <div className="flex items-center gap-1.5">
                          <Badge variant="destructive">Banned</Badge>
                          <Badge variant="secondary" className="bg-slate-100 dark:bg-slate-800 text-slate-800 dark:text-slate-200 border-0">Muted</Badge>
                        </div>
                      )}
                      {(entry.status === 'whitelisted_muted') && (
                        <div className="flex items-center gap-1.5">
                          <Badge className="bg-emerald-600">Whitelisted</Badge>
                          <Badge variant="secondary" className="bg-slate-100 dark:bg-slate-800 text-slate-800 dark:text-slate-200 border-0">Muted</Badge>
                        </div>
                      )}
                    </TableCell>
                    <TableCell className="text-slate-500 text-sm">{entry.reason || "-"}</TableCell>
                    <TableCell className="text-right space-x-2">
                      <Button
                        size="sm"
                        variant="ghost"
                        className="h-8 w-8 p-0 text-slate-500 hover:text-slate-900"
                        onClick={() => setSelectedDetailIP(entry.ip)}
                        title="View IP Details"
                      >
                        <Eye className="h-4 w-4" />
                      </Button>
                      <Button size="sm" variant="ghost" className="h-8 w-8 p-0 text-slate-500 hover:text-rose-600" onClick={() => handleIPSubmit('remove', entry.ip)}>
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>

          {/* Pagination Controls */}
          {ipList.length > 0 && (
            <div className="flex items-center justify-between border-t pt-4">
              <div className="flex items-center gap-2">
                <span className="text-xs text-slate-500">Rows per page:</span>
                <Select
                  value={String(ipRowsPerPage)}
                  onValueChange={(val) => {
                    setIpRowsPerPage(Number(val));
                    setIpCurrentPage(1);
                  }}
                >
                  <SelectTrigger className="w-[85px] h-8 text-xs bg-background text-foreground border border-input">
                    <SelectValue placeholder="5 rows" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="5">5 rows</SelectItem>
                    <SelectItem value="10">10 rows</SelectItem>
                    <SelectItem value="20">20 rows</SelectItem>
                    <SelectItem value="50">50 rows</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="flex items-center gap-2">
                <span className="text-xs text-slate-500">
                  Page {ipCurrentPage} of {Math.ceil(ipList.length / ipRowsPerPage) || 1}
                </span>
                <div className="flex gap-1">
                  <Button
                    size="sm"
                    variant="outline"
                    className="h-8 text-xs px-2"
                    disabled={ipCurrentPage === 1}
                    onClick={() => setIpCurrentPage(prev => Math.max(1, prev - 1))}
                  >
                    Previous
                  </Button>
                  <Button
                    size="sm"
                    variant="outline"
                    className="h-8 text-xs px-2"
                    disabled={ipCurrentPage >= Math.ceil(ipList.length / ipRowsPerPage)}
                    onClick={() => setIpCurrentPage(prev => Math.min(Math.ceil(ipList.length / ipRowsPerPage), prev + 1))}
                  >
                    Next
                  </Button>
                </div>
              </div>
            </div>
          )}
        </CardContent>
      </Card>

      <div className="grid gap-6 lg:grid-cols-2">
        <Card className="border-0 shadow-none bg-zinc-50/40 dark:bg-slate-900/40">
          <CardHeader className="flex flex-row items-center justify-between">
            <div>
              <CardTitle className="font-bold font-display tracking-tight text-lg text-foreground">Active Rules Registry</CardTitle>
              <CardDescription className="text-muted-foreground text-sm">Dynamically toggle security rules from waf-rules.json.</CardDescription>
            </div>
            <Button size="sm" className="bg-primary hover:bg-primary/90 text-primary-foreground rounded-xl" onClick={() => setIsAddRuleOpen(true)}>+ Add Rule</Button>
          </CardHeader>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>ID</TableHead>
                  <TableHead>Description</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="text-right">Action</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {rules.map((rule: any) => (
                  <TableRow key={rule.id}>
                    <TableCell className="font-medium text-xs">{rule.id}</TableCell>
                    <TableCell className="text-sm">{rule.description || rule.id}</TableCell>
                    <TableCell>
                      {rule.enabled ? (
                        <Badge variant="secondary" className="bg-emerald-100 text-emerald-800">Active</Badge>
                      ) : (
                        <Badge variant="outline" className="text-slate-400">Disabled</Badge>
                      )}
                    </TableCell>
                    <TableCell className="text-right whitespace-nowrap">
                      <div className="flex gap-1.5 justify-end">
                        <Button
                          size="sm"
                          variant="outline"
                          className="h-8 text-xs font-medium border-slate-300 dark:border-slate-800 hover:bg-slate-100 dark:hover:bg-slate-850 text-slate-700 dark:text-slate-200"
                          onClick={() => setSelectedRuleForDetail(rule)}
                        >
                          View
                        </Button>
                        <div className="flex items-center justify-center min-w-[56px]">
                          <Switch
                            checked={rule.enabled}
                            onCheckedChange={() => toggleRule(rule.id, rule.enabled)}
                          />
                        </div>
                        <Button
                          size="sm"
                          variant="destructive"
                          className="h-8 text-xs font-medium"
                          onClick={() => handleDeleteRule(rule.id)}
                        >
                          Delete
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
                {rules.length === 0 && (
                  <TableRow>
                    <TableCell colSpan={4} className="text-center text-slate-500 py-6">
                      No rules loaded or backend is down.
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </CardContent>
        </Card>

        <Card className="border-0 shadow-none bg-zinc-50/40 dark:bg-slate-900/40">
          <CardHeader className="flex flex-col sm:flex-row sm:items-center justify-between space-y-2 sm:space-y-0">
            <div>
              <CardTitle className="font-bold font-display tracking-tight text-lg text-foreground">Security Logs</CardTitle>
              <CardDescription className="text-muted-foreground text-sm">Recent suspicious activities and blocks.</CardDescription>
            </div>
            <div className="flex items-center gap-2">
              <Button
                size="sm"
                variant="destructive"
                className="h-8 text-xs font-semibold rounded-xl bg-destructive hover:bg-destructive/90 text-destructive-foreground"
                onClick={handleClearAllLogs}
              >
                Clear All Logs
              </Button>
              <Select
                value={statusFilter}
                onValueChange={(val) => setStatusFilter(val)}
              >
                <SelectTrigger className="w-[110px] h-8 text-xs bg-background text-foreground border border-border rounded-xl">
                  <SelectValue placeholder="All Status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="All">All Status</SelectItem>
                  <SelectItem value="Passed">Passed</SelectItem>
                  <SelectItem value="Blocked">Blocked</SelectItem>
                </SelectContent>
              </Select>

              <Select
                value={String(rowsPerPage)}
                onValueChange={(val) => setRowsPerPage(Number(val))}
              >
                <SelectTrigger className="w-[85px] h-8 text-xs bg-background text-foreground border border-border rounded-xl">
                  <SelectValue placeholder="5 rows" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="5">5 rows</SelectItem>
                  <SelectItem value="10">10 rows</SelectItem>
                  <SelectItem value="20">20 rows</SelectItem>
                  <SelectItem value="50">50 rows</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </CardHeader>
          <CardContent>
            {Object.keys(selectedLogIds).filter(id => selectedLogIds[id]).length > 0 && (
              <div className="flex items-center justify-between bg-muted/40 border border-border/80 rounded-2xl p-4 mb-4 animate-in fade-in slide-in-from-top-2 duration-200">
                <span className="text-xs font-semibold text-muted-foreground">
                  {Object.keys(selectedLogIds).filter(id => selectedLogIds[id]).length} log(s) selected
                </span>
                <div className="flex gap-2">
                  <Button
                    size="sm"
                    variant="destructive"
                    className="h-8 text-xs font-semibold rounded-xl"
                    onClick={() => handleBulkIPAction('ban')}
                  >
                    Bulk Ban IPs
                  </Button>
                  <Button
                    size="sm"
                    className="h-8 text-xs font-semibold bg-primary hover:bg-primary/90 text-primary-foreground rounded-xl"
                    onClick={() => handleBulkIPAction('whitelist')}
                  >
                    Bulk Whitelist IPs
                  </Button>
                  <Button
                    size="sm"
                    variant="outline"
                    className="h-8 text-xs"
                    onClick={handleBulkDeleteLogs}
                  >
                    Delete Logs
                  </Button>
                </div>
              </div>
            )}
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className="w-[40px]">
                    <Checkbox
                      checked={displayLogs.length > 0 && displayLogs.every(l => {
                        const logId = l.id || `${l.timestamp}-${l.ip}`;
                        return selectedLogIds[logId];
                      })}
                      onCheckedChange={(checked) => handleSelectAll(!!checked)}
                    />
                  </TableHead>
                  <TableHead>Time</TableHead>
                  <TableHead>IP</TableHead>
                  <TableHead>Payload</TableHead>
                  <TableHead>Rule</TableHead>
                  <TableHead>Action</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {displayLogs.slice(0, rowsPerPage).map((log) => {
                  const logId = log.id || `${log.timestamp}-${log.ip}`;
                  return (
                    <TableRow key={logId}>
                      <TableCell>
                        <Checkbox
                          checked={!!selectedLogIds[logId]}
                          onCheckedChange={(checked) => handleSelectLog(logId, !!checked)}
                        />
                      </TableCell>
                      <TableCell className="text-xs whitespace-nowrap">
                        {parseSafeDate(log.timestamp).toLocaleTimeString()}
                      </TableCell>
                      <TableCell className="font-mono text-xs">{log.ip}</TableCell>
                      <TableCell
                        className="text-xs truncate max-w-[150px] cursor-pointer hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors text-blue-600 underline-offset-4 hover:underline"
                        title="Click to view full payload"
                        onClick={() => setSelectedLog(log)}
                      >
                        <span className="font-bold mr-1 text-slate-800 dark:text-slate-200">{log.method}</span>
                        {log.url}
                      </TableCell>
                      <TableCell className="text-xs">
                        {log.ruleId ? <Badge variant="secondary" className="text-[10px]">{log.ruleId}</Badge> : '-'}
                      </TableCell>
                      <TableCell>
                        {log.status === 'Blocked' ? (
                          <Badge variant="destructive" className="text-[10px]">Blocked</Badge>
                        ) : (
                          <Badge variant="outline" className="text-[10px]">Allowed</Badge>
                        )}
                      </TableCell>
                    </TableRow>
                  );
                })}
                {displayLogs.length === 0 && (
                  <TableRow>
                    <TableCell colSpan={6} className="text-center text-slate-500 py-6">
                      No logs found matching filters.
                    </TableCell>
                  </TableRow>
                )}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      </div>

      {/* Payload Details Sheet */}
      <Sheet open={!!selectedLog} onOpenChange={(open) => !open && setSelectedLog(null)}>
        <SheetContent className="sm:max-w-[500px] overflow-y-auto">
          <SheetHeader className="mb-6">
            <SheetTitle>Payload Details</SheetTitle>
            <SheetDescription>
              Detailed view of the HTTP request and WAF evaluation.
            </SheetDescription>
          </SheetHeader>
          {selectedLog && (
            <div className="space-y-6">
              <div className="grid grid-cols-2 gap-y-3 text-sm">
                <div className="font-semibold text-slate-500">Timestamp</div>
                <div>{parseSafeDate(selectedLog.timestamp).toLocaleString()}</div>

                <div className="font-semibold text-slate-500">IP Address</div>
                <div className="font-mono">{selectedLog.ip}</div>

                <div className="font-semibold text-slate-500">Status</div>
                <div>
                  {selectedLog.status === 'Blocked' ? (
                    <Badge variant="destructive">Blocked</Badge>
                  ) : (
                    <Badge variant="outline">Allowed</Badge>
                  )}
                </div>

                <div className="font-semibold text-slate-500">Rule ID</div>
                <div>{selectedLog.rule_id ? <Badge variant="secondary">{selectedLog.rule_id}</Badge> : '-'}</div>

                <div className="font-semibold text-slate-500">Method</div>
                <div className="font-bold">{selectedLog.method}</div>
              </div>              <div>
                <div className="font-bold font-display text-foreground mb-2 text-sm">Full Request URL / Payload:</div>
                <div className="bg-slate-950 text-emerald-400 p-4 rounded-xl text-xs font-mono break-all max-h-[220px] overflow-y-auto border border-border shadow-inner">
                  {selectedLog.url}
                </div>
              </div>

              {selectedLog.details && (
                <div>
                  <div className="font-bold font-display text-foreground mb-2 text-sm">Evaluation Details:</div>
                  <div className="bg-muted/50 text-foreground p-3.5 rounded-xl text-xs font-mono border border-border">
                    {selectedLog.details}
                  </div>
                </div>
              )}

              <hr className="border-border/85" />

              <div className="space-y-4">
                <div className="flex items-center justify-between">
                  <div className="font-bold font-display text-foreground text-sm">IP Threat Intel & Location</div>
                  {vtResult && (
                    <Badge variant={vtResult.stats.malicious > 0 ? "destructive" : "secondary"} className={vtResult.stats.malicious > 0 ? "bg-red-100 text-red-800" : "bg-primary/10 text-primary border-0"}>
                      {vtResult.stats.malicious > 0 ? `${vtResult.stats.malicious} Malicious` : "Clean / Harmless"}
                    </Badge>
                  )}
                </div>

                {/* Status indicator: Key set or Demo fallback */}
                {!VIRUSTOTAL_API_KEY && (
                  <p className="text-[10px] text-muted-foreground bg-muted/30 p-3 rounded-xl border border-border">
                    ℹ️ Running in demo mode. Configure your API key in <code>SecurityPage.tsx</code> to use your own limits.
                  </p>
                )}

                {!vtResult && !vtLoading && (
                  <Button
                    size="sm"
                    className="w-full h-9 bg-primary hover:bg-primary/90 text-primary-foreground font-semibold rounded-xl flex items-center justify-center gap-1.5"
                    onClick={() => handleCheckReputation(selectedLog.ip)}
                  >
                    <Activity className="h-4 w-4" /> Check Reputation & Location
                  </Button>
                )}

                {vtLoading && (
                  <div className="flex items-center justify-center p-6 bg-muted/20 border border-border rounded-xl">
                    <div className="flex flex-col items-center gap-2 text-muted-foreground text-xs">
                      <div className="animate-spin rounded-full h-5 w-5 border-2 border-primary/20 border-t-primary"></div>
                      Fetching IP intelligence & coordinates...
                    </div>
                  </div>
                )}

                {vtError && (
                  <div className="p-3 bg-red-50 dark:bg-red-950/30 text-red-700 dark:text-red-400 text-xs rounded-xl border border-red-200 dark:border-red-900/50 font-medium">
                    ⚠️ Error: {vtError}
                  </div>
                )}

                {vtResult && (
                  <div className="bg-muted/30 border border-border rounded-xl p-4 space-y-3 text-xs animate-in fade-in duration-200">
                    <div className="grid grid-cols-2 gap-y-2 border-b border-border/80 pb-2.5">
                      <div className="text-muted-foreground font-medium">ISP / AS Owner</div>
                      <div className="font-semibold text-slate-850 dark:text-slate-100 break-words">{vtResult.asOwner}</div>

                      <div className="text-slate-500 font-medium">Location</div>
                      <div className="font-semibold text-slate-850 dark:text-slate-100">
                        {vtResult.cityName && `${vtResult.cityName}, `}
                        {vtResult.regionName && `${vtResult.regionName}, `}
                        {vtResult.country} ({vtResult.countryCode})
                      </div>

                      <div className="text-slate-500 font-medium">Network (CIDR)</div>
                      <div className="font-mono text-slate-850 dark:text-slate-200">{vtResult.network}</div>

                      <div className="text-slate-500 font-medium">Proxy/VPN Status</div>
                      <div>
                        {vtResult.isProxy ? (
                          <Badge variant="destructive" className="bg-red-100 text-red-800 text-[10px] font-semibold">VPN / Proxy Detected</Badge>
                        ) : (
                          <Badge variant="outline" className="text-emerald-700 bg-emerald-50 border-emerald-200 text-[10px] font-semibold">Residential / Direct</Badge>
                        )}
                      </div>
                    </div>

                    <div className="space-y-1.5">
                      <div className="text-slate-500 font-medium pb-0.5">VirusTotal Engine Votes:</div>
                      <div className="grid grid-cols-4 gap-2 text-center text-[10px]">
                        <div className="bg-red-50 dark:bg-red-950/20 text-red-800 dark:text-red-400 rounded p-1.5 border border-red-100 dark:border-red-900/30 font-semibold">
                          <span className="block text-xs font-bold text-red-600 dark:text-red-400">{vtResult.stats.malicious}</span>
                          Malicious
                        </div>
                        <div className="bg-orange-50 dark:bg-orange-950/20 text-orange-800 dark:text-orange-400 rounded p-1.5 border border-orange-100 dark:border-orange-900/30 font-semibold">
                          <span className="block text-xs font-bold text-orange-600 dark:text-orange-400">{vtResult.stats.suspicious}</span>
                          Suspicious
                        </div>
                        <div className="bg-emerald-50 dark:bg-emerald-950/20 text-emerald-800 dark:text-emerald-400 rounded p-1.5 border border-emerald-100 dark:border-emerald-900/30 font-semibold">
                          <span className="block text-xs font-bold text-emerald-600 dark:text-emerald-400">{vtResult.stats.harmless}</span>
                          Harmless
                        </div>
                        <div className="bg-slate-100 dark:bg-slate-800 text-slate-800 dark:text-slate-200 rounded p-1.5 border border-slate-200 dark:border-slate-700 font-semibold">
                          <span className="block text-xs font-bold text-slate-600 dark:text-slate-400">{vtResult.stats.undetected}</span>
                          Undetected
                        </div>
                      </div>
                    </div>
                  </div>
                )}
              </div>
            </div>
          )}
        </SheetContent>
      </Sheet>

      {/* Add Rule Sheet */}
      <Sheet open={isAddRuleOpen} onOpenChange={setIsAddRuleOpen}>
        <SheetContent>
          <SheetHeader className="mb-6">
            <SheetTitle>Add New WAF Rule</SheetTitle>
            <SheetDescription>
              Create a new regex pattern rule to intercept malicious traffic.
            </SheetDescription>
          </SheetHeader>
          <div className="space-y-4">
            <div>
              <label className="text-sm font-medium mb-1 block text-slate-700 dark:text-slate-300">Description</label>
              <Input placeholder="e.g. Detect malicious payload" value={newRuleDesc} onChange={(e) => setNewRuleDesc(e.target.value)} />
            </div>
            <div>
              <label className="text-sm font-medium mb-1 block text-slate-700 dark:text-slate-300">Regex Pattern</label>
              <Input placeholder="e.g. (?i)(select|drop|union)" value={newRulePattern} onChange={(e) => setNewRulePattern(e.target.value)} />
              <p className="text-[10px] text-slate-500 mt-1">Must be a valid Go regex pattern.</p>
            </div>
            <div>
              <label className="text-sm font-medium mb-1 block text-slate-700 dark:text-slate-300">Tags (comma separated)</label>
              <Input placeholder="e.g. sqli, xss" value={newRuleTags} onChange={(e) => setNewRuleTags(e.target.value)} />
            </div>
            <div>
              <label className="text-sm font-medium mb-1 block text-slate-700 dark:text-slate-300">Impact Score (1-10)</label>
              <Input type="number" min="1" max="10" value={newRuleImpact} onChange={(e) => setNewRuleImpact(e.target.value)} />
            </div>
            <Button className="w-full mt-4" onClick={handleAddRule}>Save Rule</Button>
          </div>
        </SheetContent>
      </Sheet>

      {/* Rule Details Dialog */}
      <Dialog open={!!selectedRuleForDetail} onOpenChange={(open) => { if (!open) { setSelectedRuleForDetail(null); setIsEditingRule(false); } }}>
        <DialogContent className="sm:max-w-[600px]">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2 text-xl font-bold">
              <ShieldAlert className="text-emerald-500 h-5 w-5" /> WAF Rule Profile
            </DialogTitle>
            <DialogDescription>
              {isEditingRule ? "Modify rule settings and signature pattern." : "Technical specifications and signature pattern for this rule."}
            </DialogDescription>
          </DialogHeader>

          {selectedRuleForDetail && (
            isEditingRule ? (
              <div className="space-y-4">
                <div>
                  <label className="text-sm font-semibold mb-1 block text-slate-700">Description</label>
                  <Input
                    value={editRuleDesc}
                    onChange={(e) => setEditRuleDesc(e.target.value)}
                    placeholder="e.g. Detect SQL comment bypass"
                  />
                </div>
                <div>
                  <label className="text-sm font-semibold mb-1 block text-slate-700">Regex Pattern</label>
                  <Input
                    value={editRulePattern}
                    onChange={(e) => setEditRulePattern(e.target.value)}
                    placeholder="e.g. (?i)(select|union)"
                    className="font-mono text-xs"
                  />
                  <p className="text-[10px] text-slate-500 mt-1">Must be a valid Go regex pattern.</p>
                </div>
                <div>
                  <label className="text-sm font-semibold mb-1 block text-slate-700">Tags (comma separated)</label>
                  <Input
                    value={editRuleTags}
                    onChange={(e) => setEditRuleTags(e.target.value)}
                    placeholder="e.g. sqli, bypass"
                  />
                </div>
                <div>
                  <label className="text-sm font-semibold mb-1 block text-slate-700">Impact Score (1-10)</label>
                  <Input
                    type="number"
                    min="1"
                    max="10"
                    value={editRuleImpact}
                    onChange={(e) => setEditRuleImpact(e.target.value)}
                  />
                </div>
                <div className="flex justify-end gap-2 mt-6">
                  <Button onClick={handleUpdateRule}>Save Changes</Button>
                  <Button variant="outline" onClick={() => setIsEditingRule(false)}>Cancel</Button>
                </div>
              </div>
            ) : (
              <div className="space-y-6">
                <div className="grid grid-cols-2 gap-4 bg-slate-50 dark:bg-slate-900/40 p-4 rounded-lg border border-slate-100 dark:border-slate-800">
                  <div>
                    <div className="text-xs font-semibold text-slate-400 uppercase tracking-wider">Rule ID</div>
                    <div className="mt-1 font-mono text-sm text-slate-800 dark:text-slate-200">{selectedRuleForDetail.id}</div>
                  </div>
                  <div>
                    <div className="text-xs font-semibold text-slate-400 uppercase tracking-wider">Status</div>
                    <div className="mt-1">
                      {selectedRuleForDetail.enabled ? (
                        <Badge variant="secondary" className="bg-emerald-100 dark:bg-emerald-950/30 text-emerald-800 dark:text-emerald-450 border-0">Active</Badge>
                      ) : (
                        <Badge variant="outline" className="text-slate-400 border-slate-300 dark:border-slate-800">Disabled</Badge>
                      )}
                    </div>
                  </div>
                  <div>
                    <div className="text-xs font-semibold text-slate-400 uppercase tracking-wider">Impact Score</div>
                    <div className="mt-1 text-sm font-semibold text-slate-800 dark:text-slate-200">{selectedRuleForDetail.impact || '5'}/10</div>
                  </div>
                  <div>
                    <div className="text-xs font-semibold text-slate-400 uppercase tracking-wider">Tags</div>
                    <div className="mt-1 flex flex-wrap gap-1">
                      {selectedRuleForDetail.tags && selectedRuleForDetail.tags.length > 0 ? (
                        selectedRuleForDetail.tags.map((t: string) => (
                          <Badge key={t} variant="outline" className="text-[10px] bg-slate-100 dark:bg-slate-800 text-slate-700 dark:text-slate-300 border-transparent">{t}</Badge>
                        ))
                      ) : (
                        <span className="text-slate-400 text-xs">-</span>
                      )}
                    </div>
                  </div>
                </div>

                <div>
                  <div className="text-sm font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Description</div>
                  <div className="text-sm text-slate-800 dark:text-slate-200 bg-white dark:bg-slate-950 border border-slate-200 dark:border-slate-800 rounded p-3 leading-relaxed shadow-sm">
                    {selectedRuleForDetail.description}
                  </div>
                </div>

                <div>
                  <div className="text-sm font-semibold text-slate-700 dark:text-slate-300 mb-1.5">Regular Expression Pattern</div>
                  <div className="bg-slate-950 text-emerald-400 p-4 rounded-md text-xs font-mono break-all max-h-[150px] overflow-y-auto border border-slate-800 shadow-inner">
                    {selectedRuleForDetail.pattern}
                  </div>
                </div>

                <div className="flex items-center justify-between border-t border-slate-100 dark:border-slate-850 pt-4 mt-2">
                  <div className="flex items-center gap-2">
                    <span className="text-sm font-semibold text-slate-750 dark:text-slate-300">Status:</span>
                    <Switch
                      checked={selectedRuleForDetail.enabled}
                      onCheckedChange={() => {
                        toggleRule(selectedRuleForDetail.id, selectedRuleForDetail.enabled);
                        setSelectedRuleForDetail((prev: any) => prev ? { ...prev, enabled: !prev.enabled } : null);
                      }}
                    />
                    <span className="text-xs text-slate-500">
                      {selectedRuleForDetail.enabled ? "Active" : "Disabled"}
                    </span>
                  </div>
                  <div className="flex gap-2">
                    <Button
                      variant="outline"
                      className="border-slate-300 dark:border-slate-800 text-slate-700 dark:text-slate-200 hover:bg-slate-100 dark:hover:bg-slate-800"
                      onClick={() => {
                        setEditRuleDesc(selectedRuleForDetail.description || "");
                        setEditRulePattern(selectedRuleForDetail.pattern || "");
                        setEditRuleTags((selectedRuleForDetail.tags || []).join(", "));
                        setEditRuleImpact(selectedRuleForDetail.impact || "5");
                        setIsEditingRule(true);
                      }}
                    >
                      Edit Rule
                    </Button>
                    <Button
                      variant="destructive"
                      onClick={() => {
                        handleDeleteRule(selectedRuleForDetail.id);
                        setSelectedRuleForDetail(null);
                      }}
                    >
                      Delete
                    </Button>
                  </div>
                </div>
                <Button variant="secondary" className="w-full mt-2" onClick={() => setSelectedRuleForDetail(null)}>Close</Button>
              </div>
            )
          )}
        </DialogContent>
      </Dialog>

      {/* IP Details Modal */}
      <Dialog open={!!selectedDetailIP} onOpenChange={(open) => !open && setSelectedDetailIP(null)}>
        <DialogContent className="sm:max-w-[700px]">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2 text-xl font-bold">
              <ShieldAlert className="text-rose-500 h-5 w-5" /> IP Threat Profile
            </DialogTitle>
            <DialogDescription className="font-mono text-slate-500">
              Details for IP Address: {selectedDetailIP}
            </DialogDescription>
          </DialogHeader>

          {selectedDetailIP && (() => {
            // Loaded entry from list
            const ipLogs = logs.filter(l => l.ip === selectedDetailIP);
            const blockedLogs = ipLogs.filter(l => l.status === 'Blocked');
            const allowedLogs = ipLogs.filter(l => l.status === 'Allowed');

            return (
              <div className="space-y-6">
                <div className="grid grid-cols-2 gap-4 bg-slate-50 dark:bg-slate-900/40 p-4 rounded-lg border border-slate-100 dark:border-slate-800">
                  <div>
                    <div className="text-xs font-semibold text-slate-400 uppercase tracking-wider mb-1">Access Status</div>
                    <Select
                      value={editIPStatus}
                      onValueChange={(val) => setEditIPStatus(val)}
                    >
                      <SelectTrigger className="w-full h-9 bg-background text-foreground border border-input">
                        <SelectValue placeholder="Select Status" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="ban">Banned (Blocked)</SelectItem>
                        <SelectItem value="whitelist">Whitelisted (Allowed)</SelectItem>
                        <SelectItem value="ignore">Muted (No Log)</SelectItem>
                        <SelectItem value="banned_muted">Banned & Muted (Block, No Log)</SelectItem>
                        <SelectItem value="whitelisted_muted">Whitelisted & Muted (Allow, No Log)</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                  <div>
                    <div className="text-xs font-semibold text-slate-400 uppercase tracking-wider mb-1">Estimated Category</div>
                    <div className="h-8 flex items-center">
                      <Badge variant="outline" className="px-2 py-0.5 border-slate-300 dark:border-slate-700 text-slate-700 dark:text-slate-200 font-medium bg-background">
                        {getIPCategory(editIPReason)}
                      </Badge>
                    </div>
                  </div>
                  <div className="col-span-2">
                    <div className="text-xs font-semibold text-slate-400 uppercase tracking-wider mb-1">Action Trigger / Reason</div>
                    <Input
                      className="bg-background text-foreground border border-input rounded px-3 py-1.5 text-sm w-full"
                      value={editIPReason}
                      onChange={(e) => setEditIPReason(e.target.value)}
                      placeholder="e.g. Local developer testing / DDoS source"
                    />
                  </div>
                  <div>
                    <div className="text-xs font-semibold text-slate-400 uppercase tracking-wider">Total Threat Events</div>
                    <div className="mt-1 text-sm font-semibold text-rose-600">{blockedLogs.length} events blocked</div>
                  </div>
                  <div>
                    <div className="text-xs font-semibold text-slate-400 uppercase tracking-wider">Clean Requests Allowed</div>
                    <div className="mt-1 text-sm font-semibold text-emerald-600">{allowedLogs.length} requests</div>
                  </div>
                </div>

                <div>
                  <div className="text-sm font-bold text-slate-800 dark:text-slate-200 mb-3 flex items-center gap-1.5">
                    <Activity className="h-4 w-4 text-slate-500" /> Recent Activity Log ({ipLogs.length} records)
                  </div>
                  <div className="overflow-hidden border border-slate-200 dark:border-slate-800 rounded-md max-h-[250px] overflow-y-auto shadow-inner bg-white dark:bg-slate-950">
                    <Table>
                      <TableHeader className="bg-slate-50 dark:bg-slate-900/60">
                        <TableRow>
                          <TableHead className="text-xs py-2 h-9">Time</TableHead>
                          <TableHead className="text-xs py-2 h-9">Method & URL</TableHead>
                          <TableHead className="text-xs py-2 h-9">Rule Triggered</TableHead>
                          <TableHead className="text-xs py-2 h-9 text-right">Result</TableHead>
                        </TableRow>
                      </TableHeader>
                      <TableBody>
                        {ipLogs.length === 0 ? (
                          <TableRow>
                            <TableCell colSpan={4} className="text-center text-slate-400 py-6 text-xs">No recent traffic records found for this IP.</TableCell>
                          </TableRow>
                        ) : (
                          ipLogs.map((l, idx) => (
                            <TableRow key={idx} className="hover:bg-slate-50/50 dark:hover:bg-slate-800/30">
                              <TableCell className="text-[11px] whitespace-nowrap text-slate-500 py-1.5">
                                {parseSafeDate(l.timestamp).toLocaleTimeString()}
                              </TableCell>
                              <TableCell className="text-[11px] font-mono max-w-[280px] truncate text-slate-800 dark:text-slate-205 py-1.5" title={l.url}>
                                <span className="font-bold mr-1">{l.method}</span>{l.url}
                              </TableCell>
                              <TableCell className="text-[11px] py-1.5">
                                {l.rule_id ? <Badge variant="outline" className="text-[9px] px-1.5">{l.rule_id}</Badge> : <span className="text-slate-400">-</span>}
                              </TableCell>
                              <TableCell className="text-right py-1.5">
                                {l.status === 'Blocked' ? (
                                  <Badge variant="destructive" className="text-[9px] px-1.5">Blocked</Badge>
                                ) : (
                                  <Badge variant="outline" className="text-[9px] px-1.5 text-emerald-600 dark:text-emerald-400 border-emerald-200 dark:border-emerald-900/30 bg-emerald-50 dark:bg-emerald-950/20">Allowed</Badge>
                                )}
                              </TableCell>
                            </TableRow>
                          ))
                        )}
                      </TableBody>
                    </Table>
                  </div>
                </div>

                <div className="flex justify-between items-center border-t pt-4">
                  <Button
                    variant="outline"
                    size="sm"
                    className="h-9 border-rose-200 text-rose-600 hover:bg-rose-50 dark:hover:bg-rose-950/20 hover:text-rose-700 dark:text-rose-450 dark:border-rose-900/50 font-medium"
                    onClick={() => {
                      handleIPSubmit('remove', selectedDetailIP);
                      setSelectedDetailIP(null);
                    }}
                  >
                    Forget & Unban IP
                  </Button>
                  <div className="flex gap-2">
                    <Button
                      variant="outline"
                      size="sm"
                      className="h-9"
                      onClick={() => setSelectedDetailIP(null)}
                    >
                      Close
                    </Button>
                    <Button
                      size="sm"
                      className="h-9"
                      onClick={async () => {
                        await handleIPSubmit(editIPStatus, selectedDetailIP, editIPReason);
                        setSelectedDetailIP(null);
                      }}
                    >
                      Save Changes
                    </Button>
                  </div>
                </div>
              </div>
            );
          })()}
        </DialogContent>
      </Dialog>

      {/* Category Analysis Modal */}
      <Dialog open={isCategoryModalOpen} onOpenChange={setIsCategoryModalOpen}>
        <DialogContent className="sm:max-w-[850px] p-0 overflow-hidden">
          <div className="flex h-[550px] divide-x divide-slate-200 dark:divide-slate-800">
            {/* Sidebar Categories List */}
            <div className="w-[280px] bg-slate-50 dark:bg-slate-900 p-6 flex flex-col justify-between">
              <div className="space-y-6">
                <div>
                  <h3 className="text-base font-bold text-slate-900 dark:text-slate-100 flex items-center gap-2">
                    <ListFilter className="h-4 w-4 text-emerald-600" /> Categories
                  </h3>
                  <p className="text-xs text-slate-500 dark:text-slate-400 mt-1">Select category to audit banned IPs</p>
                </div>

                <div className="space-y-1.5">
                  {[
                    'Rate Limiting',
                    'SQL Injection',
                    'XSS Attacks',
                    'Path Traversal (LFI)',
                    'RCE & Command Injection',
                    'Malicious User-Agent',
                    'Recon & Probing',
                    'Manual / Other'
                  ].map((category) => {
                    const count = ipList.filter(e => getIPCategory(e.reason) === category).length;
                    const isActive = selectedAnalysisCategory === category;

                    return (
                      <button
                        key={category}
                        className={`w-full flex items-center justify-between px-3 py-2 rounded-lg text-xs font-semibold transition-all ${isActive
                          ? 'bg-emerald-600 text-white shadow-md shadow-emerald-600/10'
                          : 'hover:bg-slate-200/60 dark:hover:bg-slate-850/60 text-slate-700 dark:text-slate-350'
                          }`}
                        onClick={() => setSelectedAnalysisCategory(category)}
                      >
                        <span className="truncate mr-2">{category}</span>
                        <Badge
                          className={`text-[10px] px-1.5 py-0.5 rounded-full ${isActive
                            ? 'bg-white/20 text-white border-transparent'
                            : 'bg-slate-200 dark:bg-slate-800 text-slate-700 dark:text-slate-300'
                            }`}
                        >
                          {count}
                        </Badge>
                      </button>
                    );
                  })}
                </div>
              </div>

              <div className="text-[11px] text-slate-400 dark:text-slate-500 border-t dark:border-slate-800 pt-4">
                Powered by live Go WAF rules parser.
              </div>
            </div>

            {/* Right Pane List of IPs */}
            <div className="flex-1 p-6 flex flex-col justify-between bg-white dark:bg-slate-950">
              <div className="space-y-4 overflow-hidden flex flex-col h-full">
                <div className="border-b dark:border-slate-800 pb-3">
                  <div className="flex items-center gap-2">
                    <Badge className="bg-emerald-100 dark:bg-emerald-950/40 text-emerald-800 dark:text-emerald-400 border-transparent text-xs hover:bg-emerald-100 font-bold">
                      Category
                    </Badge>
                    <h2 className="text-lg font-bold text-slate-900 dark:text-slate-100">{selectedAnalysisCategory}</h2>
                  </div>
                  <p className="text-xs text-slate-500 dark:text-slate-400 mt-1">
                    IPs flagged by Go WAF matching signatures for this category.
                  </p>
                </div>

                {(() => {
                  const filteredIPs = ipList.filter(e => getIPCategory(e.reason) === selectedAnalysisCategory);

                  return (
                    <div className="flex-1 overflow-y-auto pr-1 min-h-0">
                      {filteredIPs.length === 0 ? (
                        <div className="flex flex-col items-center justify-center h-full py-12 text-center text-slate-400 space-y-2">
                          <CheckCircle2 className="h-10 w-10 text-slate-300 dark:text-slate-700" />
                          <p className="text-sm font-semibold text-slate-500 dark:text-slate-400">No IPs in this category</p>
                          <p className="text-xs text-slate-400 dark:text-slate-500">No threats detected matching these signatures currently.</p>
                        </div>
                      ) : (
                        <div className="space-y-2.5">
                          {filteredIPs.map((entry) => {
                            const ipLogsCount = logs.filter(l => l.ip === entry.ip).length;
                            return (
                              <div
                                key={entry.ip}
                                className="flex items-center justify-between p-3 border border-slate-100 dark:border-slate-800 rounded-lg hover:bg-slate-50 dark:hover:bg-slate-900 transition-colors shadow-sm bg-card"
                              >
                                <div className="space-y-1">
                                  <div className="flex items-center gap-2">
                                    <span
                                      className="font-mono text-sm font-bold text-slate-850 dark:text-slate-200 cursor-pointer hover:underline hover:text-emerald-600"
                                      onClick={() => {
                                        setSelectedDetailIP(entry.ip);
                                      }}
                                    >
                                      {entry.ip}
                                    </span>
                                    {entry.status === 'Banned' ? (
                                      <Badge variant="destructive" className="text-[9px] px-1.5 py-0">Banned</Badge>
                                    ) : (
                                      <Badge className="bg-emerald-600 text-white text-[9px] px-1.5 py-0">Whitelist</Badge>
                                    )}
                                  </div>
                                  <div className="text-[11px] text-slate-500 dark:text-slate-400 font-medium truncate max-w-[350px]" title={entry.reason}>
                                    Reason: {entry.reason || "-"}
                                  </div>
                                  <div className="text-[10px] text-slate-400 dark:text-slate-500">
                                    Activity log count: {ipLogsCount} records
                                  </div>
                                </div>
                                <div className="flex gap-2">
                                  <Button
                                    size="sm"
                                    variant="outline"
                                    className="h-8 text-xs font-semibold border-slate-200 dark:border-slate-800 text-slate-700 dark:text-slate-200 hover:bg-slate-100 dark:hover:bg-slate-900"
                                    onClick={() => setSelectedDetailIP(entry.ip)}
                                  >
                                    Details
                                  </Button>
                                  <Button
                                    size="sm"
                                    variant="ghost"
                                    className="h-8 text-xs text-rose-600 dark:text-rose-400 hover:text-rose-700 hover:bg-rose-50 dark:hover:bg-rose-950/20 font-semibold"
                                    onClick={() => handleIPSubmit('remove', entry.ip)}
                                  >
                                    Unban
                                  </Button>
                                </div>
                              </div>
                            );
                          })}
                        </div>
                      )}
                    </div>
                  );
                })()}
              </div>

              <div className="border-t dark:border-slate-800 pt-4 flex justify-end">
                <Button
                  className="bg-slate-900 dark:bg-slate-800 hover:bg-slate-800 dark:hover:bg-slate-700 text-white"
                  onClick={() => setIsCategoryModalOpen(false)}
                >
                  Close Analysis
                </Button>
              </div>
            </div>
          </div>
        </DialogContent>
      </Dialog>

      <Dialog open={isConfirmClearLogsOpen} onOpenChange={setIsConfirmClearLogsOpen}>
        <DialogContent className="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle className="text-slate-900 dark:text-white">Clear All Security Logs</DialogTitle>
            <DialogDescription className="mt-2 text-slate-500">
              Are you sure you want to permanently clear ALL security logs from the database? This action cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <div className="flex justify-end gap-2 mt-6">
            <Button
              variant="outline"
              onClick={() => setIsConfirmClearLogsOpen(false)}
            >
              Cancel
            </Button>
            <Button
              variant="destructive"
              onClick={() => {
                setIsConfirmClearLogsOpen(false);
                executeClearAllLogs();
              }}
            >
              Clear Logs
            </Button>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  );
}
