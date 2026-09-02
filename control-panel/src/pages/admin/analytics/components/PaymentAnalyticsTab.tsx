import { CreditCard, CheckCircle2, Clock, AlertOctagon, RotateCcw } from 'lucide-react';
import { ResponsiveContainer, BarChart, Bar, XAxis, YAxis, Tooltip, Cell, LabelList } from 'recharts';
import type { PaymentMetricsResponse } from '../../../../models/Analytics';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '../../../../components/ui/table';
import { Skeleton } from '../../../../components/ui/skeleton';

interface Props {
  data: PaymentMetricsResponse | null;
  loading: boolean;
}

export default function PaymentAnalyticsTab({ data, loading }: Props) {
  const formatIDR = (val?: number | null) => {
    if (val === undefined || val === null || isNaN(val)) return 'Rp 0';
    return `Rp ${val.toLocaleString('id-ID')}`;
  };

  if (loading) {
    return (
      <div className="space-y-6">
        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4 border-b border-border/60 pb-8">
          {Array.from({ length: 4 }).map((_, i) => (
            <div key={i} className="space-y-2">
              <Skeleton className="h-4 w-28 bg-muted animate-pulse" />
              <Skeleton className="h-8 w-36 bg-muted animate-pulse" />
              <Skeleton className="h-3 w-40 bg-muted animate-pulse" />
            </div>
          ))}
        </div>
        <Skeleton className="h-64 rounded-2xl bg-muted animate-pulse" />
      </div>
    );
  }

  const summary = data?.summary;
  const breakdown = Array.isArray(data?.breakdown) ? data.breakdown : [];

  const chartData = breakdown.map(item => ({
    name: item.method_name || 'Unknown',
    amount: typeof item.amount === 'number' ? item.amount : 0,
    count: typeof item.count === 'number' ? item.count : 0,
    rate: typeof item.success_rate === 'number' ? item.success_rate : 0,
  }));

  const colors = ['hsl(142.4, 71.8%, 29.2%)', '#3b82f6', '#8b5cf6', '#f59e0b', '#ec4899'];

  return (
    <div className="space-y-8 animate-in fade-in slide-in-from-left-4 duration-300">
      {/* Payment KPI Grid - SecurityPage Borderless Style */}
      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4 border-b border-border/60 pb-8">
        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <span className="text-sm font-bold font-display text-foreground">Total Paid</span>
            <CheckCircle2 className="h-4 w-4 text-emerald-600" />
          </div>
          <div className="text-3xl font-bold tracking-tight text-emerald-600 dark:text-emerald-400">{formatIDR(summary?.total_paid)}</div>
          <p className="text-xs text-muted-foreground">
            Success Rate: <span className="font-semibold text-emerald-600">{((summary?.payment_success_rate || 0) * 100).toFixed(1)}%</span>
          </p>
        </div>

        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <span className="text-sm font-bold font-display text-foreground">Pending Payment</span>
            <Clock className="h-4 w-4 text-amber-500" />
          </div>
          <div className="text-3xl font-bold tracking-tight text-amber-600 dark:text-amber-400">{formatIDR(summary?.total_pending)}</div>
          <p className="text-xs text-muted-foreground">
            Avg time to pay: {summary?.avg_time_to_pay ? `${(summary.avg_time_to_pay / 60).toFixed(1)} mins` : 'N/A'}
          </p>
        </div>

        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <span className="text-sm font-bold font-display text-foreground">Expired Payments</span>
            <AlertOctagon className="h-4 w-4 text-muted-foreground" />
          </div>
          <div className="text-3xl font-bold tracking-tight text-foreground">{formatIDR(summary?.total_expired)}</div>
          <p className="text-xs text-muted-foreground">Unpaid expired checkouts</p>
        </div>

        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <span className="text-sm font-bold font-display text-foreground">Refunded Amount</span>
            <RotateCcw className="h-4 w-4 text-rose-500" />
          </div>
          <div className="text-3xl font-bold tracking-tight text-rose-600 dark:text-rose-400">{formatIDR(summary?.total_refunded)}</div>
          <p className="text-xs text-rose-500 font-medium">Total customer refunds processed</p>
        </div>
      </div>

      {/* Payment Method Chart & Table */}
      <div className="grid gap-8 md:grid-cols-7">
        {/* Payment Chart */}
        <div className="md:col-span-4 space-y-4">
          <h3 className="font-bold font-display text-lg text-foreground flex items-center gap-2">
            <CreditCard className="h-5 w-5 text-primary" />
            Revenue by Payment Method
          </h3>
          <div className="h-[260px] w-full">
            {chartData.length === 0 ? (
              <div className="h-full flex items-center justify-center text-xs text-muted-foreground">
                No payment breakdown data available.
              </div>
            ) : (
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={chartData} margin={{ top: 20, right: 10, left: -10, bottom: 5 }}>
                  <XAxis dataKey="name" axisLine={false} tickLine={false} tick={{ fontSize: 10, fill: 'var(--muted-foreground)' }} />
                  <YAxis
                    axisLine={false}
                    tickLine={false}
                    tickFormatter={(val) => `Rp ${(val / 1000000).toFixed(0)}M`}
                    tick={{ fontSize: 10, fill: 'var(--muted-foreground)' }}
                  />
                  <Tooltip
                    content={({ active, payload }) => {
                      if (active && payload && payload.length) {
                        const d = payload[0].payload;
                        return (
                          <div className="bg-popover text-popover-foreground border border-border rounded-xl p-3 shadow-md text-xs font-sans">
                            <p className="font-semibold font-display mb-1">{d.name}</p>
                            <p className="text-primary font-semibold">Volume: {formatIDR(d.amount)}</p>
                            <p className="text-muted-foreground">Transactions: {d.count}</p>
                            <p className="text-muted-foreground">Success Rate: {(d.rate * 100).toFixed(1)}%</p>
                          </div>
                        );
                      }
                      return null;
                    }}
                  />
                  <Bar dataKey="amount" radius={[6, 6, 0, 0]}>
                    {chartData.map((_, idx) => (
                      <Cell key={`cell-${idx}`} fill={colors[idx % colors.length]} />
                    ))}
                    <LabelList
                      dataKey="amount"
                      position="top"
                      formatter={(val: any) => {
                        const num = Number(val);
                        if (!num) return '';
                        return num >= 1000000 ? `${(num / 1000000).toFixed(1)}M` : `${(num / 1000).toFixed(0)}k`;
                      }}
                      style={{ fontSize: 10, fontWeight: 700, fill: 'var(--foreground)' }}
                    />
                  </Bar>
                </BarChart>
              </ResponsiveContainer>
            )}
          </div>
        </div>

        {/* Payment Table */}
        <div className="md:col-span-3 space-y-4">
          <h3 className="font-bold font-display text-base text-foreground">Method Details</h3>
          {breakdown.length === 0 ? (
            <p className="text-xs text-muted-foreground">No payment records found.</p>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Method</TableHead>
                  <TableHead className="text-right">Txns</TableHead>
                  <TableHead className="text-right">Success</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {breakdown.map((item) => (
                  <TableRow key={item.method_id}>
                    <TableCell className="font-medium text-xs text-foreground">{item.method_name}</TableCell>
                    <TableCell className="text-right text-xs text-muted-foreground">{item.count}</TableCell>
                    <TableCell className="text-right text-xs font-semibold text-emerald-600">
                      {((item.success_rate || 0) * 100).toFixed(0)}%
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </div>
      </div>
    </div>
  );
}
