import { useState, useEffect } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { ShieldAlert, ShieldX, Activity, CheckCircle2, ShieldCheck, Trash2 } from 'lucide-react';
import { ResponsiveContainer, AreaChart, Area, XAxis, YAxis, Tooltip, CartesianGrid, ReferenceArea } from 'recharts';
import { Sheet, SheetContent, SheetDescription, SheetHeader, SheetTitle } from '@/components/ui/sheet';

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

export default function SecurityPage() {
  const [logs, setLogs] = useState<any[]>([]);
  const [wafSummary, setWafSummary] = useState({ total: 0, blocked: 0, allowed: 0, threatLevel: 'Low' });
  const [threatData, setThreatData] = useState<any[]>([]);

  const [rules, setRules] = useState<any[]>([]);
  const [ipList, setIpList] = useState<any[]>([]);
  const [targetIP, setTargetIP] = useState("");
  const [targetReason, setTargetReason] = useState("");
  const [selectedLog, setSelectedLog] = useState<any | null>(null);

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
      await Promise.all(
        uniqueIPs.map(ip =>
          fetch('http://localhost:8080/api/ip', {
            method: 'POST',
            body: JSON.stringify({ ip, action, reason: `Bulk ${action} from Logs` }),
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

    try {
      await fetch('http://localhost:8080/api/stats?t=' + Date.now(), {
        method: 'DELETE',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ ids: selectedIds }),
      });
      // Optimistic local update: remove deleted logs from local state immediately
      setLogs(prev => prev.filter(log => {
        const logId = log.id || `${log.timestamp}-${log.ip}`;
        return !selectedLogIds[logId];
      }));
      setSelectedLogIds({});
      // Fetch latest state from server to sync
      fetchLogs();
    } catch (e) {
      console.error("Bulk delete failed", e);
    }
  };

  const handleBulkIPActionFromManager = async (action: 'ban' | 'whitelist' | 'remove') => {
    const ips = Object.keys(selectedIPs).filter(ip => selectedIPs[ip]);
    if (ips.length === 0) return;

    try {
      await Promise.all(
        ips.map(ip =>
          fetch('http://localhost:8080/api/ip', {
            method: 'POST',
            body: JSON.stringify({ ip, action, reason: action === 'remove' ? '' : 'Bulk Action' }),
          })
        )
      );
      // Optimistic local UI updates
      if (action === 'remove') {
        setIpList(prev => prev.filter(entry => !selectedIPs[entry.ip]));
      } else {
        setIpList(prev => prev.map(entry => {
          if (selectedIPs[entry.ip]) {
            return { ...entry, status: action === 'ban' ? 'Banned' : 'Whitelisted', reason: 'Bulk Action' };
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
      await fetch('http://localhost:8080/api/rules', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
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
      const res = await fetch('http://localhost:8080/api/rules?t=' + Date.now());
      setRules(await res.json());
    } catch (e) {
      console.error(e);
    }
  };

  const fetchIpList = async () => {
    try {
      const res = await fetch('http://localhost:8080/api/ip?t=' + Date.now());
      setIpList((await res.json()) || []);
    } catch (e) {
      console.error(e);
    }
  };

  const fetchLogs = async () => {
    try {
      const res = await fetch('http://localhost:8080/api/stats?t=' + Date.now());
      const data = await res.json();
      const rawLogs = data.logs || [];
      const sorted = [...rawLogs].sort((a, b) => parseSafeDate(b.timestamp).getTime() - parseSafeDate(a.timestamp).getTime());
      setLogs(sorted);
      
      const total = rawLogs.length;
      const blocked = rawLogs.filter((l: any) => l.status === 'Blocked').length;
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

    setThreatData(Object.values(dataMap).sort((a, b) => a.time - b.time));
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

  const handleIPSubmit = async (action: string, ipOverride?: string, reasonOverride?: string) => {
    const ip = ipOverride || targetIP;
    const reason = reasonOverride !== undefined ? reasonOverride : targetReason;
    if (!ip) return;

    await fetch('http://localhost:8080/api/ip', {
      method: 'POST',
      body: JSON.stringify({ ip, action, reason }),
    });

    if (!ipOverride) {
      setTargetIP(""); 
      setTargetReason(""); 
    }
    fetchIpList(); 
  };

  const toggleRule = async (ruleId: string, currentStatus: boolean) => {
    await fetch('http://localhost:8080/api/rules', {
      method: 'PUT',
      body: JSON.stringify({ id: ruleId, enabled: !currentStatus }),
    });
    fetchRules();
  };

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
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h2 className="text-2xl font-bold tracking-tight text-slate-900">Security Monitoring</h2>
          <p className="text-slate-500">Real-time threat analysis and WAF status.</p>
        </div>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card className="bg-slate-900 text-white border-slate-800">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium text-slate-300">Attack Attempts</CardTitle>
            <ShieldAlert className="h-4 w-4 text-red-400" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-red-400">{wafSummary.blocked.toLocaleString()}</div>
            <p className="text-xs text-slate-400">Total requests blocked by rules</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Safe Traffic</CardTitle>
            <CheckCircle2 className="h-4 w-4 text-emerald-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-emerald-600">{wafSummary.allowed.toLocaleString()}</div>
            <p className="text-xs text-slate-500">Allowed requests</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Anomaly Score</CardTitle>
            <Activity className="h-4 w-4 text-orange-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-orange-600">
              {((wafSummary.blocked / wafSummary.total) * 100).toFixed(1)}%
            </div>
            <p className="text-xs text-slate-500">Traffic flagged as malicious</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Active Incidents</CardTitle>
            <ShieldX className="h-4 w-4 text-rose-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-rose-600">{ipList.filter(ip => ip.status === 'Banned').length}</div>
            <p className="text-xs text-slate-500">Currently banned IPs</p>
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader className="flex flex-col sm:flex-row sm:items-center justify-between">
          <div>
            <CardTitle>Threat Landscape (Real Data)</CardTitle>
            <CardDescription>Volume of allowed vs blocked requests over time.</CardDescription>
          </div>
          <div className="flex items-center gap-3 mt-4 sm:mt-0">
            <select 
              className="bg-white border border-slate-300 rounded-md px-3 py-1.5 text-sm"
              value={rangeType}
              onChange={(e) => {
                const val = e.target.value as any;
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
              <option value="Today">Today (00:00 - 24:00)</option>
              <option value="Custom">Custom Range</option>
            </select>
            {rangeType === 'Custom' && (
              <div className="flex items-center gap-2">
                <Input type="datetime-local" className="h-8 text-xs w-[170px]" value={tempStart} onChange={e => setTempStart(e.target.value)} />
                <span className="text-slate-400 text-sm">to</span>
                <Input type="datetime-local" className="h-8 text-xs w-[170px]" value={tempEnd} onChange={e => setTempEnd(e.target.value)} />
                <Button size="sm" className="h-8 bg-slate-900 hover:bg-slate-800 text-white font-medium px-3" onClick={handleApplyCustomRange}>
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
                  <stop offset="5%" stopColor="#ef4444" stopOpacity={0.3}/>
                  <stop offset="95%" stopColor="#ef4444" stopOpacity={0}/>
                </linearGradient>
                <linearGradient id="colorAllowed" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#10b981" stopOpacity={0.1}/>
                  <stop offset="95%" stopColor="#10b981" stopOpacity={0}/>
                </linearGradient>
              </defs>
              <XAxis 
                dataKey="time" 
                type="number"
                scale="time"
                domain={[activeStartTime, activeEndTime]}
                axisLine={false} 
                tickLine={false} 
                tick={{ fontSize: 12 }} 
                dy={10} 
                tickFormatter={(unixTime) => {
                  const d = new Date(unixTime);
                  return `${String(d.getHours()).padStart(2,'0')}:${String(d.getMinutes()).padStart(2,'0')}`;
                }}
              />
              <YAxis axisLine={false} tickLine={false} tick={{ fontSize: 12 }} />
              <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="#e2e8f0" />
              <Tooltip 
                contentStyle={{ borderRadius: '8px', border: '1px solid #e2e8f0', boxShadow: '0 4px 12px rgba(0,0,0,0.1)' }}
                labelFormatter={(unixTime: any) => {
                  const d = new Date(unixTime);
                  return `${String(d.getHours()).padStart(2,'0')}:${String(d.getMinutes()).padStart(2,'0')}`;
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

      {/* IP Blacklist Manager */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2"><ShieldCheck className="text-emerald-500" /> IP Blacklist Manager</CardTitle>
          <CardDescription>Control access and logging behavior for specific IPs.</CardDescription>
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="flex flex-col md:flex-row gap-4 items-end bg-slate-50 p-4 rounded-lg border border-slate-100">
            <div className="space-y-2 flex-1 w-full">
              <label className="text-xs font-medium text-slate-700">Target IP Address</label>
              <Input
                placeholder="e.g. 192.168.1.5 or ::1"
                className="bg-white border-slate-300 font-mono"
                value={targetIP}
                onChange={e => setTargetIP(e.target.value)}
              />
            </div>
            <div className="space-y-2 flex-1 w-full">
              <label className="text-xs font-medium text-slate-700">Reason (Optional)</label>
              <Input
                placeholder="e.g. Spamming API"
                className="bg-white border-slate-300"
                value={targetReason}
                onChange={e => setTargetReason(e.target.value)}
              />
            </div>
            <div className="flex gap-2 w-full md:w-auto">
              <Button variant="destructive" className="flex-1 md:flex-none" onClick={() => handleIPSubmit('ban')}>Ban</Button>
              <Button className="bg-emerald-600 hover:bg-emerald-700 flex-1 md:flex-none" onClick={() => handleIPSubmit('whitelist')}>Whitelist</Button>
            </div>
          </div>

          {Object.keys(selectedIPs).filter(ip => selectedIPs[ip]).length > 0 && (
            <div className="flex items-center justify-between bg-slate-50 border border-slate-200 rounded-lg p-3 mb-4 animate-in fade-in slide-in-from-top-2 duration-200">
              <span className="text-xs font-semibold text-slate-700">
                {Object.keys(selectedIPs).filter(ip => selectedIPs[ip]).length} IP(s) selected
              </span>
              <div className="flex gap-2">
                <Button 
                  size="sm" 
                  variant="destructive" 
                  className="h-8 text-xs font-medium" 
                  onClick={() => handleBulkIPActionFromManager('ban')}
                >
                  Bulk Ban
                </Button>
                <Button 
                  size="sm" 
                  className="h-8 text-xs bg-emerald-600 hover:bg-emerald-700 text-white font-medium" 
                  onClick={() => handleBulkIPActionFromManager('whitelist')}
                >
                  Bulk Whitelist
                </Button>
                <Button 
                  size="sm" 
                  variant="outline" 
                  className="h-8 text-xs border-slate-300 hover:bg-slate-100 text-slate-700 font-medium" 
                  onClick={() => handleBulkIPActionFromManager('remove')}
                >
                  Bulk Remove / Forget
                </Button>
              </div>
            </div>
          )}

          <div className="overflow-x-auto border rounded-md">
            <Table>
              <TableHeader>
                <TableRow className="bg-slate-50">
                  <TableHead className="w-[40px]">
                    <input 
                      type="checkbox" 
                      className="rounded border-slate-300 text-slate-900 focus:ring-slate-500 h-4 w-4 cursor-pointer"
                      checked={ipList.length > 0 && ipList.every(entry => selectedIPs[entry.ip])}
                      onChange={e => {
                        const checked = e.target.checked;
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
                {ipList.map((entry: any) => (
                  <TableRow key={entry.ip}>
                    <TableCell>
                      <input 
                        type="checkbox" 
                        className="rounded border-slate-300 text-slate-900 focus:ring-slate-500 h-4 w-4 cursor-pointer"
                        checked={!!selectedIPs[entry.ip]}
                        onChange={e => {
                          setSelectedIPs(prev => ({
                            ...prev,
                            [entry.ip]: e.target.checked
                          }));
                        }}
                      />
                    </TableCell>
                    <TableCell className="font-mono text-slate-700">{entry.ip}</TableCell>
                    <TableCell>
                      {entry.status === 'Banned' && <Badge variant="destructive">Banned</Badge>}
                      {entry.status === 'Whitelisted' && <Badge className="bg-emerald-600">Whitelisted</Badge>}
                      {entry.status === 'Ignored' && <Badge variant="secondary">Muted (No Log)</Badge>}
                    </TableCell>
                    <TableCell className="text-slate-500 text-sm">{entry.reason || "-"}</TableCell>
                    <TableCell className="text-right space-x-2">
                      <Button size="sm" variant="ghost" className="h-8 w-8 p-0 text-slate-500 hover:text-rose-600" onClick={() => handleIPSubmit('remove', entry.ip)}>
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
        </CardContent>
      </Card>

      <div className="grid gap-6 lg:grid-cols-2">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between">
            <div>
              <CardTitle>Active Rules Registry</CardTitle>
              <CardDescription>Dynamically toggle security rules from waf-rules.json.</CardDescription>
            </div>
            <Button size="sm" onClick={() => setIsAddRuleOpen(true)}>+ Add Rule</Button>
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
                    <TableCell className="text-right">
                       <Button 
                          size="sm" 
                          variant={rule.enabled ? "outline" : "default"}
                          className={!rule.enabled ? "bg-emerald-600 hover:bg-emerald-700" : ""}
                          onClick={() => toggleRule(rule.id, rule.enabled)}
                       >
                         {rule.enabled ? "Disable" : "Enable"}
                       </Button>
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

        <Card>
          <CardHeader className="flex flex-col sm:flex-row sm:items-center justify-between space-y-2 sm:space-y-0">
            <div>
              <CardTitle>Security Logs</CardTitle>
              <CardDescription>Recent suspicious activities and blocks.</CardDescription>
            </div>
            <div className="flex items-center gap-2">
              <select
                className="bg-white border border-slate-300 rounded-md px-2 py-1 text-xs font-medium text-slate-700 focus:outline-none focus:ring-1 focus:ring-slate-400"
                value={statusFilter}
                onChange={(e) => setStatusFilter(e.target.value)}
              >
                <option value="All">All Status</option>
                <option value="Passed">Passed</option>
                <option value="Blocked">Blocked</option>
              </select>
              <select
                className="bg-white border border-slate-300 rounded-md px-2 py-1 text-xs font-medium text-slate-700 focus:outline-none focus:ring-1 focus:ring-slate-400"
                value={rowsPerPage}
                onChange={(e) => setRowsPerPage(Number(e.target.value))}
              >
                <option value={5}>5 rows</option>
                <option value={10}>10 rows</option>
                <option value={20}>20 rows</option>
                <option value={50}>50 rows</option>
              </select>
            </div>
          </CardHeader>
          <CardContent>
            {Object.keys(selectedLogIds).filter(id => selectedLogIds[id]).length > 0 && (
              <div className="flex items-center justify-between bg-slate-50 border border-slate-200 rounded-lg p-3 mb-4 animate-in fade-in slide-in-from-top-2 duration-200">
                <span className="text-xs font-semibold text-slate-700">
                  {Object.keys(selectedLogIds).filter(id => selectedLogIds[id]).length} log(s) selected
                </span>
                <div className="flex gap-2">
                  <Button 
                    size="sm" 
                    variant="destructive" 
                    className="h-8 text-xs font-medium" 
                    onClick={() => handleBulkIPAction('ban')}
                  >
                    Bulk Ban IPs
                  </Button>
                  <Button 
                    size="sm" 
                    className="h-8 text-xs bg-emerald-600 hover:bg-emerald-700 text-white font-medium" 
                    onClick={() => handleBulkIPAction('whitelist')}
                  >
                    Bulk Whitelist IPs
                  </Button>
                  <Button 
                    size="sm" 
                    variant="outline" 
                    className="h-8 text-xs border-slate-300 hover:bg-slate-100 text-slate-700 font-medium" 
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
                    <input 
                      type="checkbox" 
                      className="rounded border-slate-300 text-slate-900 focus:ring-slate-500 h-4 w-4 cursor-pointer"
                      checked={displayLogs.length > 0 && displayLogs.every(l => {
                        const logId = l.id || `${l.timestamp}-${l.ip}`;
                        return selectedLogIds[logId];
                      })}
                      onChange={e => handleSelectAll(e.target.checked)}
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
                        <input 
                          type="checkbox" 
                          className="rounded border-slate-300 text-slate-900 focus:ring-slate-500 h-4 w-4 cursor-pointer"
                          checked={!!selectedLogIds[logId]}
                          onChange={e => handleSelectLog(logId, e.target.checked)}
                        />
                      </TableCell>
                      <TableCell className="text-xs whitespace-nowrap">
                        {parseSafeDate(log.timestamp).toLocaleTimeString()}
                      </TableCell>
                      <TableCell className="font-mono text-xs">{log.ip}</TableCell>
                      <TableCell 
                        className="text-xs truncate max-w-[150px] cursor-pointer hover:bg-slate-50 transition-colors text-blue-600 underline-offset-4 hover:underline" 
                        title="Click to view full payload"
                        onClick={() => setSelectedLog(log)}
                      >
                        <span className="font-bold mr-1 text-slate-800">{log.method}</span>
                        {log.url}
                      </TableCell>
                      <TableCell className="text-xs">
                        {log.rule_id ? <Badge variant="secondary" className="text-[10px]">{log.rule_id}</Badge> : '-'}
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
        <SheetContent className="sm:max-w-[500px]">
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
              </div>

              <div>
                <div className="font-semibold text-slate-800 mb-2 text-sm">Full Request URL / Payload:</div>
                <div className="bg-slate-950 text-emerald-400 p-4 rounded-md text-xs font-mono break-all max-h-[300px] overflow-y-auto border border-slate-800 shadow-inner">
                  {selectedLog.url}
                </div>
              </div>
              
              {selectedLog.details && (
                <div>
                  <div className="font-semibold text-slate-800 mb-2 text-sm">Evaluation Details:</div>
                  <div className="bg-slate-100 text-slate-800 p-3 rounded-md text-xs font-mono border border-slate-200">
                    {selectedLog.details}
                  </div>
                </div>
              )}
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
              <label className="text-sm font-medium mb-1 block text-slate-700">Description</label>
              <Input placeholder="e.g. Detect malicious payload" value={newRuleDesc} onChange={(e) => setNewRuleDesc(e.target.value)} />
            </div>
            <div>
              <label className="text-sm font-medium mb-1 block text-slate-700">Regex Pattern</label>
              <Input placeholder="e.g. (?i)(select|drop|union)" value={newRulePattern} onChange={(e) => setNewRulePattern(e.target.value)} />
              <p className="text-[10px] text-slate-500 mt-1">Must be a valid Go regex pattern.</p>
            </div>
            <div>
              <label className="text-sm font-medium mb-1 block text-slate-700">Tags (comma separated)</label>
              <Input placeholder="e.g. sqli, xss" value={newRuleTags} onChange={(e) => setNewRuleTags(e.target.value)} />
            </div>
            <div>
              <label className="text-sm font-medium mb-1 block text-slate-700">Impact Score (1-10)</label>
              <Input type="number" min="1" max="10" value={newRuleImpact} onChange={(e) => setNewRuleImpact(e.target.value)} />
            </div>
            <Button className="w-full mt-4" onClick={handleAddRule}>Save Rule</Button>
          </div>
        </SheetContent>
      </Sheet>
    </div>
  );
}
