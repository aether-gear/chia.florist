import React, { useMemo } from 'react';
import {
  ResponsiveContainer,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip,
  CartesianGrid,
} from 'recharts';
import { ShieldAlert, Activity, Eye, Lock, Globe } from 'lucide-react';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '../../../components/ui/table';
import { Badge } from '../../../components/ui/badge';
import { Button } from '../../../components/ui/button';
import { Skeleton } from '../../../components/ui/skeleton';
import type { DashboardWafSummary, SecurityEventLog } from '../../../viewmodels/useDashboardViewModel';

interface DashboardCyberViewProps {
  wafSummary: DashboardWafSummary;
  securityLogs: SecurityEventLog[];
  loading: boolean;
  onInspectSecurityLog: (log: SecurityEventLog) => void;
}

export const DashboardCyberView: React.FC<DashboardCyberViewProps> = ({
  wafSummary,
  securityLogs,
  loading,
  onInspectSecurityLog,
}) => {
  // Compute time intervals for allowed vs blocked requests chart
  const chartData = useMemo(() => {
    if (!securityLogs || securityLogs.length === 0) return [];

    const buckets: Record<string, { time: string; allowed: number; blocked: number }> = {};

    // Group logs into hourly buckets
    securityLogs.forEach((log) => {
      try {
        const d = new Date(log.timestamp);
        const hourKey = `${String(d.getHours()).padStart(2, '0')}:00`;
        if (!buckets[hourKey]) {
          buckets[hourKey] = { time: hourKey, allowed: 0, blocked: 0 };
        }
        if (log.status === 'Blocked') {
          buckets[hourKey].blocked += 1;
        } else {
          buckets[hourKey].allowed += 1;
        }
      } catch {
        // Ignore date parsing error
      }
    });

    return Object.values(buckets).slice(-12);
  }, [securityLogs]);

  // Compute Attack Vector Categories
  const categoryStats = useMemo(() => {
    const counts = {
      sqli: 0,
      xss: 0,
      lfi: 0,
      rce: 0,
      rateLimit: 0,
      scanners: 0,
      other: 0,
    };

    securityLogs.forEach((log) => {
      if (log.status !== 'Blocked') return;
      const text = `${log.reason} ${log.url} ${log.ruleId}`.toLowerCase();
      if (text.includes('sql') || text.includes('union') || text.includes('select')) counts.sqli++;
      else if (text.includes('xss') || text.includes('script')) counts.xss++;
      else if (text.includes('passwd') || text.includes('traversal') || text.includes('..')) counts.lfi++;
      else if (text.includes('rce') || text.includes('exec') || text.includes('command')) counts.rce++;
      else if (text.includes('rate') || text.includes('limit') || text.includes('flood')) counts.rateLimit++;
      else if (text.includes('scan') || text.includes('nikto') || text.includes('sqlmap')) counts.scanners++;
      else counts.other++;
    });

    return counts;
  }, [securityLogs]);

  return (
    <div className="space-y-10 animate-in fade-in slide-in-from-left-4 duration-300">
      {/* 1. Cyber Posture & Active Policies Stats Row */}
      <div className="grid gap-4 grid-cols-2 sm:grid-cols-4 pb-6 border-b border-border/60">
        <div className="p-3.5 rounded-xl bg-muted/40 border border-border/40 space-y-1">
          <div className="flex items-center justify-between">
            <span className="text-[11px] text-muted-foreground font-sans uppercase tracking-wider font-semibold">Total Inspected</span>
            <Activity className="w-3.5 h-3.5 text-muted-foreground/70" />
          </div>
          <p className="text-xl font-bold font-display text-foreground">{wafSummary.total.toLocaleString()}</p>
          <p className="text-[11px] text-muted-foreground">Evaluated HTTP requests</p>
        </div>

        <div className="p-3.5 rounded-xl bg-muted/40 border border-border/40 space-y-1">
          <div className="flex items-center justify-between">
            <span className="text-[11px] text-muted-foreground font-sans uppercase tracking-wider font-semibold">Blocked Attacks</span>
            <ShieldAlert className="w-3.5 h-3.5 text-rose-500" />
          </div>
          <p className="text-xl font-bold font-display text-rose-600 dark:text-rose-400">{wafSummary.blocked}</p>
          <p className="text-[11px] text-rose-500 font-medium">Malicious vectors dropped</p>
        </div>

        <div className="p-3.5 rounded-xl bg-muted/40 border border-border/40 space-y-1">
          <div className="flex items-center justify-between">
            <span className="text-[11px] text-muted-foreground font-sans uppercase tracking-wider font-semibold">Active WAF Rules</span>
            <Lock className="w-3.5 h-3.5 text-primary" />
          </div>
          <p className="text-xl font-bold font-display text-foreground">{wafSummary.activeRules}</p>
          <p className="text-[11px] text-muted-foreground">Regex & signature policies</p>
        </div>

        <div className="p-3.5 rounded-xl bg-muted/40 border border-border/40 space-y-1">
          <div className="flex items-center justify-between">
            <span className="text-[11px] text-muted-foreground font-sans uppercase tracking-wider font-semibold">IP Enforcements</span>
            <Globe className="w-3.5 h-3.5 text-sky-500" />
          </div>
          <p className="text-xl font-bold font-display text-foreground">
            {wafSummary.bannedIps} <span className="text-xs font-normal text-muted-foreground">banned</span>
          </p>
          <p className="text-[11px] text-muted-foreground">{wafSummary.whitelistedIps} whitelisted</p>
        </div>
      </div>

      {/* 2. Traffic Area Chart & Attack Categories Grid */}
      <div className="grid gap-10 lg:grid-cols-7 pb-8 border-b border-border/60">
        {/* Left: Traffic Activity Chart */}
        <div className="lg:col-span-4 space-y-4">
          <div className="pb-2 border-b border-border/60 flex items-center justify-between">
            <div>
              <h3 className="font-bold font-display tracking-tight text-lg text-foreground">
                Traffic & Attack Interception Timeline
              </h3>
              <p className="text-muted-foreground text-xs font-sans">
                Real-time stream of legitimate traffic vs blocked threats.
              </p>
            </div>
          </div>

          <div className="h-[280px]">
            {loading ? (
              <div className="h-full flex flex-col justify-between py-4">
                <Skeleton className="h-4 w-full bg-muted" />
                <Skeleton className="h-36 w-full bg-muted rounded-xl" />
                <Skeleton className="h-4 w-full bg-muted" />
              </div>
            ) : chartData.length === 0 ? (
              <div className="h-full flex items-center justify-center text-xs text-muted-foreground bg-muted/10 rounded-xl border border-border/30">
                Awaiting real-time telemetry from WAF logger.
              </div>
            ) : (
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={chartData} margin={{ top: 10, right: 10, left: -20, bottom: 0 }}>
                  <CartesianGrid strokeDasharray="3 3" stroke="currentColor" className="text-border/40" />
                  <XAxis dataKey="time" tick={{ fontSize: 10, fill: 'var(--muted-foreground)' }} axisLine={false} tickLine={false} />
                  <YAxis tick={{ fontSize: 10, fill: 'var(--muted-foreground)' }} axisLine={false} tickLine={false} />
                  <Tooltip
                    content={({ active, payload }) => {
                      if (!active || !payload?.length) return null;
                      const data = payload[0].payload;
                      return (
                        <div className="bg-popover text-popover-foreground border border-border rounded-xl p-3 shadow-lg text-xs font-sans">
                          <p className="font-bold font-display mb-1">{data.time}</p>
                          <p className="text-muted-foreground">
                            Allowed: <span className="font-semibold text-primary">{data.allowed}</span>
                          </p>
                          <p className="text-muted-foreground">
                            Blocked: <span className="font-semibold text-rose-500">{data.blocked}</span>
                          </p>
                        </div>
                      );
                    }}
                  />
                  <Bar dataKey="allowed" name="Allowed" fill="hsl(var(--primary))" radius={[4, 4, 0, 0]} maxBarSize={28} />
                  <Bar dataKey="blocked" name="Blocked" fill="#ef4444" radius={[4, 4, 0, 0]} maxBarSize={28} />
                </BarChart>
              </ResponsiveContainer>
            )}
          </div>
        </div>

        {/* Right: Attack Vector Distribution */}
        <div className="lg:col-span-3 space-y-4">
          <div className="pb-2 border-b border-border/60">
            <h3 className="font-bold font-display tracking-tight text-lg text-foreground">
              Mitigated Threat Vectors
            </h3>
            <p className="text-muted-foreground text-xs font-sans">
              Categorized attack payloads intercepted by WAF filters.
            </p>
          </div>

          <div className="space-y-2.5">
            <div className="flex items-center justify-between p-3 rounded-xl bg-muted/30 border border-border/30">
              <span className="text-xs font-medium text-foreground">SQL Injection (SQLi)</span>
              <Badge variant="secondary" className="bg-rose-500/10 text-rose-600 dark:text-rose-400 border-rose-500/20 font-mono">
                {categoryStats.sqli} hits
              </Badge>
            </div>
            <div className="flex items-center justify-between p-3 rounded-xl bg-muted/30 border border-border/30">
              <span className="text-xs font-medium text-foreground">Cross-Site Scripting (XSS)</span>
              <Badge variant="secondary" className="bg-amber-500/10 text-amber-600 dark:text-amber-400 border-amber-500/20 font-mono">
                {categoryStats.xss} hits
              </Badge>
            </div>
            <div className="flex items-center justify-between p-3 rounded-xl bg-muted/30 border border-border/30">
              <span className="text-xs font-medium text-foreground">Path Traversal / LFI</span>
              <Badge variant="secondary" className="bg-indigo-500/10 text-indigo-600 dark:text-indigo-400 border-indigo-500/20 font-mono">
                {categoryStats.lfi} hits
              </Badge>
            </div>
            <div className="flex items-center justify-between p-3 rounded-xl bg-muted/30 border border-border/30">
              <span className="text-xs font-medium text-foreground">Rate Limit & Brute-Force</span>
              <Badge variant="secondary" className="bg-sky-500/10 text-sky-600 dark:text-sky-400 border-sky-500/20 font-mono">
                {categoryStats.rateLimit} hits
              </Badge>
            </div>
            <div className="flex items-center justify-between p-3 rounded-xl bg-muted/30 border border-border/30">
              <span className="text-xs font-medium text-foreground">Malicious Scanners & Probing</span>
              <Badge variant="secondary" className="bg-zinc-500/10 text-zinc-600 dark:text-zinc-400 border-zinc-500/20 font-mono">
                {categoryStats.scanners + categoryStats.other} hits
              </Badge>
            </div>
          </div>
        </div>
      </div>

      {/* 3. Live Security Audit Feed Table */}
      <div className="space-y-4">
        <div className="flex flex-row items-center justify-between pb-2 border-b border-border/60">
          <div>
            <h3 className="font-bold font-display tracking-tight text-lg text-foreground">
              Live Security Telemetry Feed
            </h3>
            <p className="text-muted-foreground text-xs font-sans">
              Evaluated incoming HTTP requests and WAF mitigation verdicts.
            </p>
          </div>
        </div>

        <div>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Time</TableHead>
                <TableHead>Client IP</TableHead>
                <TableHead>Method & Path</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Signature / Details</TableHead>
                <TableHead className="text-right">Action</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {loading ? (
                Array.from({ length: 5 }).map((_, i) => (
                  <TableRow key={i}>
                    <TableCell><Skeleton className="h-4 w-20 bg-muted" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-28 bg-muted" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-36 bg-muted" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-16 bg-muted" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-32 bg-muted" /></TableCell>
                    <TableCell className="text-right"><Skeleton className="h-7 w-16 bg-muted ml-auto" /></TableCell>
                  </TableRow>
                ))
              ) : securityLogs.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={6} className="text-center text-xs text-muted-foreground py-8">
                    No security events recorded.
                  </TableCell>
                </TableRow>
              ) : (
                securityLogs.slice(0, 10).map((log) => (
                  <TableRow key={log.id} className="hover:bg-muted/30">
                    <TableCell className="text-xs text-muted-foreground whitespace-nowrap">
                      {new Date(log.timestamp).toLocaleTimeString()}
                    </TableCell>
                    <TableCell className="font-mono text-xs font-semibold text-foreground">
                      {log.ip}
                    </TableCell>
                    <TableCell>
                      <span className="font-bold text-xs text-foreground mr-1.5">{log.method}</span>
                      <span className="text-muted-foreground text-xs font-mono truncate max-w-[200px] inline-block align-bottom">
                        {log.url}
                      </span>
                    </TableCell>
                    <TableCell>
                      {log.status === 'Blocked' ? (
                        <Badge variant="destructive" className="rounded-md uppercase">Blocked</Badge>
                      ) : (
                        <Badge variant="secondary" className="bg-primary/15 text-primary border-0 rounded-md uppercase">Allowed</Badge>
                      )}
                    </TableCell>
                    <TableCell className="text-xs text-muted-foreground truncate max-w-[220px]">
                      {log.reason !== '-' ? log.reason : log.ruleId !== '-' ? `Rule: ${log.ruleId}` : 'Legitimate Request'}
                    </TableCell>
                    <TableCell className="text-right">
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => onInspectSecurityLog(log)}
                        className="h-8 text-xs font-semibold text-primary hover:text-primary hover:bg-primary/10 rounded-lg gap-1"
                      >
                        <Eye className="w-3.5 h-3.5" />
                        Inspect
                      </Button>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </div>
      </div>
    </div>
  );
};

export default DashboardCyberView;
