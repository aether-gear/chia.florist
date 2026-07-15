import { useState, useMemo } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Shield, ShieldAlert, ShieldCheck, Zap, Activity } from 'lucide-react';
import { ResponsiveContainer, BarChart, Bar, XAxis, YAxis, Tooltip, Cell } from 'recharts';
import { useProductStatsViewModel } from '../../viewmodels/useProductStatsViewModel';
import { Skeleton } from '../../components/ui/skeleton';
import { getWafSummary, getRecentLogs } from '@/data/wafData';

export default function DashboardPage() {
  const wafSummary = getWafSummary();
  const [logCount, setLogCount] = useState(5);
  const [logStatus, setLogStatus] = useState('All');
  const [timeWindow, setTimeWindow] = useState<'7d' | '30d' | '90d'>('30d');

  const { data: statsData, loading: statsLoading, error: statsError } = useProductStatsViewModel();

  const barChartData = useMemo(() => {
    if (!statsData?.stats) return [];
    const field =
      timeWindow === '7d'
        ? 'sales_velocity_7d'
        : timeWindow === '30d'
        ? 'sales_velocity_30d'
        : 'sales_velocity_90d';

    return [...statsData.stats]
      .sort((a, b) => b[field] - a[field])
      .slice(0, 6)
      .map((item) => ({
        name: item.name,
        sales: item[field],
      }));
  }, [statsData, timeWindow]);

  const topProduct = useMemo(() => {
    if (!statsData?.stats || statsData.stats.length === 0) return null;
    return [...statsData.stats].sort((a, b) => b.sales_velocity_7d - a.sales_velocity_7d)[0];
  }, [statsData]);

  const colors = [
    'hsl(var(--chart-1))',
    'hsl(var(--chart-2))',
    'hsl(var(--chart-3))',
    'hsl(var(--chart-4))',
    'hsl(var(--chart-5))',
    '#9ca3af',
  ];

  const BarTooltip = ({ active, payload }: any) => {
    if (active && payload && payload.length) {
      const data = payload[0].payload;
      return (
        <div className="bg-popover text-popover-foreground border border-border rounded-xl p-3 shadow-md text-xs font-sans">
          <p className="font-semibold font-display mb-1">{data.name}</p>
          <p className="text-muted-foreground">
            Units Sold: <span className="font-semibold text-primary">{data.sales}</span>
          </p>
        </div>
      );
    }
    return null;
  };

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
        
        {/* Sales Velocity Graph */}
        <Card className="md:col-span-4 border-0 shadow-none bg-zinc-50/40 dark:bg-slate-900/40 flex flex-col justify-between">
          <CardHeader className="pb-2">
            <div className="flex items-center justify-between gap-2 flex-wrap">
              <div>
                <CardTitle className="font-bold font-display tracking-tight text-lg text-foreground">
                  Product Sales Velocity
                </CardTitle>
                <CardDescription className="text-muted-foreground text-xs font-sans">
                  Top products sold by volume
                </CardDescription>
              </div>
              {/* Pill Toggles */}
              <div className="flex items-center gap-1 bg-muted p-0.5 rounded-lg text-xs">
                {(['7d', '30d', '90d'] as const).map((window) => (
                  <button
                    key={window}
                    onClick={() => setTimeWindow(window)}
                    className={`px-2 py-1 rounded-md transition-colors font-medium font-sans ${
                      timeWindow === window
                        ? 'bg-primary text-primary-foreground shadow-sm'
                        : 'text-muted-foreground hover:text-foreground'
                    }`}
                  >
                    {window}
                  </button>
                ))}
              </div>
            </div>
          </CardHeader>
          <CardContent className="h-[300px] pt-4">
            {statsLoading ? (
              <div className="h-full w-full flex flex-col justify-between py-4">
                {Array.from({ length: 5 }).map((_, i) => (
                  <div key={i} className="flex items-center gap-2">
                    <Skeleton className="h-4 w-20 animate-pulse bg-muted shrink-0" />
                    <Skeleton className="h-4 w-full rounded-md animate-pulse bg-muted" />
                  </div>
                ))}
              </div>
            ) : statsError ? (
              <div className="h-full w-full flex items-center justify-center text-xs text-destructive bg-destructive/5 rounded-xl border border-destructive/10 p-4">
                Failed to load sales velocity: {statsError}
              </div>
            ) : barChartData.length === 0 ? (
              <div className="h-full w-full flex items-center justify-center text-xs text-muted-foreground">
                No product sales statistics found.
              </div>
            ) : (
              <ResponsiveContainer width="100%" height="100%">
                <BarChart
                  data={barChartData}
                  layout="vertical"
                  margin={{ top: 5, right: 10, left: -10, bottom: 5 }}
                >
                  <XAxis type="number" hide />
                  <YAxis
                    dataKey="name"
                    type="category"
                    width={85}
                    axisLine={false}
                    tickLine={false}
                    tick={{ fontSize: 10, fill: 'var(--muted-foreground)' }}
                  />
                  <Tooltip content={<BarTooltip />} cursor={{ fill: 'rgba(21, 94, 55, 0.05)' }} />
                  <Bar dataKey="sales" radius={[0, 4, 4, 0]}>
                    {barChartData.map((_, index) => (
                      <Cell key={`cell-${index}`} fill={colors[index % 5]} />
                    ))}
                  </Bar>
                </BarChart>
              </ResponsiveContainer>
            )}
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
            {statsLoading ? (
              <div className="p-4 bg-primary/5 border border-primary/10 rounded-2xl space-y-2">
                <Skeleton className="h-4 w-32 animate-pulse bg-muted" />
                <Skeleton className="h-3 w-full animate-pulse bg-muted" />
              </div>
            ) : topProduct ? (
              <div className="p-4 bg-primary/5 border border-primary/10 rounded-2xl">
                <h4 className="text-sm font-bold font-display text-primary mb-1">Top Selling Product</h4>
                <p className="text-xs text-primary/80 dark:text-primary/70 leading-relaxed">
                  "{topProduct.name}" is performing exceptionally well this week with {topProduct.sales_velocity_7d} sales, contributing to {topProduct.revenue_contribution_percentage}% of total revenue.
                </p>
              </div>
            ) : (
              <div className="p-4 bg-primary/5 border border-primary/10 rounded-2xl">
                <h4 className="text-sm font-bold font-display text-primary mb-1">Top Selling Product</h4>
                <p className="text-xs text-muted-foreground leading-relaxed">
                  No sales velocity records available for this week.
                </p>
              </div>
            )}
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
