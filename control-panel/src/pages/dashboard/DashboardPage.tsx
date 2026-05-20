import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Shield, ShieldAlert, ShieldCheck, Zap, Activity } from 'lucide-react';
import { ResponsiveContainer, BarChart, Bar, XAxis, YAxis, Tooltip, CartesianGrid } from 'recharts';
import { salesData, topSellingProducts } from '@/data/mockSales';
import { getWafSummary, getRecentLogs } from '@/data/wafData';

export default function DashboardPage() {
  const wafSummary = getWafSummary();
  const recentLogs = getRecentLogs(5);

  return (
    <div className="space-y-6">
      
      {/* WAF Summary Cards */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Requests</CardTitle>
            <Activity className="h-4 w-4 text-slate-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{wafSummary.total.toLocaleString()}</div>
            <p className="text-xs text-slate-500">Since last system restart</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Blocked Threats</CardTitle>
            <ShieldAlert className="h-4 w-4 text-red-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-red-600">{wafSummary.blocked.toLocaleString()}</div>
            <p className="text-xs text-red-500 font-medium">Action Required</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Threat Level</CardTitle>
            <Shield className="h-4 w-4 text-orange-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-orange-600">{wafSummary.threatLevel}</div>
            <p className="text-xs text-slate-500">Based on recent anomalies</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Active Rules</CardTitle>
            <ShieldCheck className="h-4 w-4 text-emerald-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-emerald-600">{wafSummary.activeRules}</div>
            <p className="text-xs text-slate-500">WAF policies enforced</p>
          </CardContent>
        </Card>
      </div>

      <div className="grid gap-4 md:grid-cols-7 lg:grid-cols-7">
        
        {/* Sales Graph */}
        <Card className="md:col-span-4">
          <CardHeader>
            <CardTitle>Sales Overview</CardTitle>
            <CardDescription>Weekly revenue and orders performance.</CardDescription>
          </CardHeader>
          <CardContent className="h-[300px]">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={salesData} margin={{ top: 10, right: 10, left: -20, bottom: 0 }}>
                <CartesianGrid strokeDasharray="3 3" vertical={false} />
                <XAxis dataKey="name" axisLine={false} tickLine={false} tick={{ fontSize: 12 }} dy={10} />
                <YAxis axisLine={false} tickLine={false} tick={{ fontSize: 12 }} tickFormatter={(val) => `$${val}`} />
                <Tooltip cursor={{ fill: 'rgba(0,0,0,0.05)' }} contentStyle={{ borderRadius: '8px', border: 'none', boxShadow: '0 4px 12px rgba(0,0,0,0.1)' }} />
                <Bar dataKey="revenue" fill="#4f46e5" radius={[4, 4, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          </CardContent>
        </Card>

        {/* AI Insights */}
        <Card className="md:col-span-3">
          <CardHeader>
            <CardTitle className="flex items-center">
              <Zap className="w-5 h-5 mr-2 text-amber-500" fill="currentColor" />
              AI Insights
            </CardTitle>
            <CardDescription>Smart business and security analysis.</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="p-3 bg-amber-50 border border-amber-100 rounded-lg">
              <h4 className="text-sm font-semibold text-amber-900 mb-1">Unusual Traffic Spike</h4>
              <p className="text-xs text-amber-800">
                We detected a 400% increase in traffic originating from unknown IPs trying to access `/catalog`. WAF successfully blocked 98% of these malicious requests.
              </p>
            </div>
            <div className="p-3 bg-emerald-50 border border-emerald-100 rounded-lg">
              <h4 className="text-sm font-semibold text-emerald-900 mb-1">Top Selling Product</h4>
              <p className="text-xs text-emerald-800">
                "{topSellingProducts[0].name}" is performing exceptionally well this week with {topSellingProducts[0].sales} sales, contributing to 65% of total revenue.
              </p>
            </div>
            <div className="p-3 bg-indigo-50 border border-indigo-100 rounded-lg">
              <h4 className="text-sm font-semibold text-indigo-900 mb-1">Conversion Suggestion</h4>
              <p className="text-xs text-indigo-800">
                Consider offering a discount on "Enterprise Setup". Abandoned cart rates are high for this tier.
              </p>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* WAF Logs Table */}
      <Card>
        <CardHeader>
          <CardTitle>Recent WAF Logs</CardTitle>
          <CardDescription>Latest traffic evaluated by the Web Application Firewall.</CardDescription>
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
                  <TableCell className="text-xs whitespace-nowrap">
                    {new Date(log.timestamp).toLocaleTimeString()}
                  </TableCell>
                  <TableCell className="font-mono text-xs">{log.ip}</TableCell>
                  <TableCell>
                    <span className="font-bold text-slate-700 mr-2">{log.method}</span>
                    <span className="text-slate-500 text-sm truncate max-w-[200px] inline-block align-bottom">{log.url}</span>
                  </TableCell>
                  <TableCell>
                    {log.status === 'Blocked' ? (
                      <Badge variant="destructive">Blocked</Badge>
                    ) : (
                      <Badge variant="secondary" className="bg-emerald-100 text-emerald-800 hover:bg-emerald-100">Allowed</Badge>
                    )}
                  </TableCell>
                  <TableCell className="text-sm text-slate-600 truncate max-w-[250px]">{log.details}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </div>
  );
}
