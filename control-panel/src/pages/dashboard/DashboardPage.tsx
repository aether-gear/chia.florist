import { useState } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Shield, ShieldAlert, ShieldCheck, Zap, Activity } from 'lucide-react';
import { ResponsiveContainer, BarChart, Bar, XAxis, YAxis, Tooltip, CartesianGrid } from 'recharts';
import { salesData, topSellingProducts } from '@/data/mockSales';
import { getWafSummary, getRecentLogs } from '@/data/wafData';

export default function DashboardPage() {
  const wafSummary = getWafSummary();
  const [logCount, setLogCount] = useState(5);
  const [logStatus, setLogStatus] = useState('All');

  const allRecentLogs = getRecentLogs(500);
  const recentLogs = allRecentLogs.filter(log => {
    if (logStatus === 'All') return true;
    return log.status === logStatus;
  }).slice(0, logCount);

  return (
    <div className="space-y-10 animate-in fade-in duration-300">
      
      {/* WAF Summary Cards */}
      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
        <Card className="border-0 shadow-none bg-zinc-50/40 dark:bg-slate-900/40">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-bold font-display text-foreground">Total Requests</CardTitle>
            <Activity className="h-4 w-4 text-muted-foreground/60" />
          </CardHeader>
          <CardContent>
            <div className="text-3xl font-bold tracking-tight text-foreground">{wafSummary.total.toLocaleString()}</div>
            <p className="text-xs text-muted-foreground mt-1">Since last system restart</p>
          </CardContent>
        </Card>
        <Card className="border-0 shadow-none bg-zinc-50/40 dark:bg-slate-900/40">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-bold font-display text-foreground">Blocked Threats</CardTitle>
            <ShieldAlert className="h-4 w-4 text-rose-500" />
          </CardHeader>
          <CardContent>
            <div className="text-3xl font-bold tracking-tight text-rose-600 dark:text-rose-400">{wafSummary.blocked.toLocaleString()}</div>
            <p className="text-xs text-rose-500 font-medium mt-1">Action Required</p>
          </CardContent>
        </Card>
        <Card className="border-0 shadow-none bg-zinc-50/40 dark:bg-slate-900/40">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-bold font-display text-foreground">Threat Level</CardTitle>
            <Shield className="h-4 w-4 text-orange-500" />
          </CardHeader>
          <CardContent>
            <div className="text-3xl font-bold tracking-tight text-orange-600 dark:text-orange-400">{wafSummary.threatLevel}</div>
            <p className="text-xs text-muted-foreground mt-1">Based on recent anomalies</p>
          </CardContent>
        </Card>
        <Card className="border-0 shadow-none bg-zinc-50/40 dark:bg-slate-900/40">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-bold font-display text-foreground">Active Rules</CardTitle>
            <ShieldCheck className="h-4 w-4 text-primary" />
          </CardHeader>
          <CardContent>
            <div className="text-3xl font-bold tracking-tight text-primary">{wafSummary.activeRules}</div>
            <p className="text-xs text-muted-foreground mt-1">WAF policies enforced</p>
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-6 md:grid-cols-7 lg:grid-cols-7">
        
        {/* Sales Graph */}
        <Card className="md:col-span-4 border-0 shadow-none bg-zinc-50/40 dark:bg-slate-900/40">
          <CardHeader>
            <CardTitle className="font-bold font-display tracking-tight text-lg">Sales Overview</CardTitle>
            <CardDescription className="text-muted-foreground text-sm">Weekly revenue and orders performance.</CardDescription>
          </CardHeader>
          <CardContent className="h-[300px]">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={salesData} margin={{ top: 10, right: 10, left: -20, bottom: 0 }}>
                <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="rgba(0,0,0,0.04)" />
                <XAxis dataKey="name" axisLine={false} tickLine={false} tick={{ fontSize: 12, fill: "currentColor" }} dy={10} />
                <YAxis axisLine={false} tickLine={false} tick={{ fontSize: 12, fill: "currentColor" }} tickFormatter={(val) => `$${val}`} />
                <Tooltip cursor={{ fill: 'rgba(0,0,0,0.02)' }} contentStyle={{ borderRadius: '12px', border: 'none', boxShadow: '0 4px 20px rgba(0,0,0,0.06)' }} />
                <Bar dataKey="revenue" fill="hsl(var(--primary))" radius={[6, 6, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>

        {/* AI Insights */}
        <Card className="md:col-span-3 border-0 shadow-none bg-zinc-50/40 dark:bg-slate-900/40">
          <CardHeader>
            <CardTitle className="flex items-center font-bold font-display tracking-tight text-lg text-foreground">
              <Zap className="w-5 h-5 mr-2 text-amber-500" fill="currentColor" />
              AI Insights
            </CardTitle>
            <CardDescription className="text-muted-foreground text-sm">Smart business and security analysis.</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="p-4 bg-amber-500/5 border border-amber-500/10 rounded-2xl">
              <h4 className="text-sm font-bold font-display text-amber-800 dark:text-amber-300 mb-1">Unusual Traffic Spike</h4>
              <p className="text-xs text-amber-700/80 dark:text-amber-300/70 leading-relaxed">
                We detected a 400% increase in traffic originating from unknown IPs trying to access `/catalog`. WAF successfully blocked 98% of these malicious requests.
              </p>
            </div>
            <div className="p-4 bg-primary/5 border border-primary/10 rounded-2xl">
              <h4 className="text-sm font-bold font-display text-primary mb-1">Top Selling Product</h4>
              <p className="text-xs text-primary/80 dark:text-primary/70 leading-relaxed">
                "{topSellingProducts[0].name}" is performing exceptionally well this week with {topSellingProducts[0].sales} sales, contributing to 65% of total revenue.
              </p>
            </div>
            <div className="p-4 bg-primary/5 border border-primary/10 rounded-2xl">
              <h4 className="text-sm font-bold font-display text-primary mb-1">Conversion Suggestion</h4>
              <p className="text-xs text-muted-foreground leading-relaxed">
                Consider offering a discount on "Enterprise Setup". Abandoned cart rates are high for this tier.
              </p>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* WAF Logs Table */}
      <Card className="border-0 shadow-none bg-zinc-50/40 dark:bg-slate-900/40">
        <CardHeader className="flex flex-row items-center justify-between">
          <div>
            <CardTitle className="font-bold font-display tracking-tight text-lg">Recent WAF Logs</CardTitle>
            <CardDescription className="text-muted-foreground text-sm">Latest traffic evaluated by the Web Application Firewall.</CardDescription>
          </div>
          <div className="flex gap-2">
            <select
              className="h-9 w-28 rounded-xl border border-border bg-background px-3 py-1 text-sm shadow-sm transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring text-foreground"
              value={logStatus}
              onChange={(e) => setLogStatus(e.target.value)}
            >
              <option value="All">All Status</option>
              <option value="Allowed">Allowed</option>
              <option value="Blocked">Blocked</option>
            </select>
            <select
              className="h-9 w-20 rounded-xl border border-border bg-background px-3 py-1 text-sm shadow-sm transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring text-foreground"
              value={logCount}
              onChange={(e) => setLogCount(Number(e.target.value))}
            >
              <option value={5}>5 Rows</option>
              <option value={10}>10 Rows</option>
              <option value={20}>20 Rows</option>
              <option value={50}>50 Rows</option>
            </select>
          </div>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Time</TableHead>
                <TableHead>IP Address</TableHead>
                <TableHead>Method & Path</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Details</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {recentLogs.map((log) => (
                <TableRow key={log.id}>
                  <TableCell className="text-xs whitespace-nowrap text-muted-foreground">
                    {new Date(log.timestamp).toLocaleTimeString()}
                  </TableCell>
                  <TableCell className="font-mono text-xs text-foreground">{log.ip}</TableCell>
                  <TableCell>
                    <span className="font-bold text-foreground mr-2">{log.method}</span>
                    <span className="text-muted-foreground text-sm truncate max-w-[200px] inline-block align-bottom">{log.url}</span>
                  </TableCell>
                  <TableCell>
                    {log.status === 'Blocked' ? (
                      <Badge variant="destructive" className="rounded-lg">Blocked</Badge>
                    ) : (
                      <Badge variant="secondary" className="bg-primary/15 text-primary hover:bg-primary/15 border-0 rounded-lg">Allowed</Badge>
                    )}
                  </TableCell>
                  <TableCell className="text-sm text-muted-foreground truncate max-w-[250px]">{log.details}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
}
