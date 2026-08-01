import { ShoppingBag, XCircle, DollarSign, Truck } from 'lucide-react';
import { ResponsiveContainer, BarChart, Bar, XAxis, YAxis, Tooltip, Cell, LabelList } from 'recharts';
import type { OrderMetricsResponse } from '../../../../models/Analytics';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '../../../../components/ui/table';
import { Skeleton } from '../../../../components/ui/skeleton';

interface Props {
  data: OrderMetricsResponse | null;
  loading: boolean;
}

export default function OrderAnalyticsTab({ data, loading }: Props) {
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
  const topProducts = Array.isArray(data?.top_products) ? data.top_products : [];
  const topShops = Array.isArray(data?.top_shops) ? data.top_shops : [];
  const totalOrders = summary?.total_orders || 1;

  const statusChartData = [
    { name: 'Pending', count: summary?.pending_count || 0, color: '#f59e0b' },
    { name: 'Confirmed', count: summary?.confirmed_count || 0, color: '#3b82f6' },
    { name: 'Processing', count: summary?.processing_count || 0, color: '#a855f7' },
    { name: 'Shipped', count: summary?.shipped_count || 0, color: '#6366f1' },
    { name: 'Delivered', count: summary?.delivered_count || 0, color: 'hsl(142.4, 71.8%, 29.2%)' },
    { name: 'Cancelled', count: summary?.cancelled_count || 0, color: '#f43f5e' },
  ];

  return (
    <div className="space-y-8 animate-in fade-in duration-300">
      {/* Order KPI Cards - SecurityPage Borderless Style */}
      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4 border-b border-border/60 pb-8">
        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <span className="text-sm font-bold font-display text-foreground">Total Orders</span>
            <ShoppingBag className="h-4 w-4 text-primary" />
          </div>
          <div className="text-3xl font-bold tracking-tight text-foreground">{summary?.total_orders || 0}</div>
          <p className="text-xs text-muted-foreground">AOV: {formatIDR(summary?.aov)}</p>
        </div>

        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <span className="text-sm font-bold font-display text-foreground">Total GMV</span>
            <DollarSign className="h-4 w-4 text-primary" />
          </div>
          <div className="text-3xl font-bold tracking-tight text-primary">{formatIDR(summary?.total_gmv)}</div>
          <p className="text-xs text-muted-foreground">Net Revenue: {formatIDR(summary?.total_revenue)}</p>
        </div>

        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <span className="text-sm font-bold font-display text-foreground">Shipping Fees</span>
            <Truck className="h-4 w-4 text-primary" />
          </div>
          <div className="text-3xl font-bold tracking-tight text-foreground">{formatIDR(summary?.total_shipping_fee)}</div>
          <p className="text-xs text-muted-foreground">Collected shipping fees</p>
        </div>

        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <span className="text-sm font-bold font-display text-foreground">Cancellation Rate</span>
            <XCircle className="h-4 w-4 text-rose-500" />
          </div>
          <div className="text-3xl font-bold tracking-tight text-rose-600 dark:text-rose-400">
            {((summary?.cancellation_rate || 0) * 100).toFixed(1)}%
          </div>
          <p className="text-xs text-rose-500 font-medium">{summary?.cancelled_count || 0} orders cancelled</p>
        </div>
      </div>

      {/* Order Status Distribution (Mobile-Optimized Column Chart + Quick Data Cards) */}
      <div className="space-y-4">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2">
          <div>
            <h3 className="font-bold font-display text-lg text-foreground flex items-center gap-2">
              <ShoppingBag className="h-5 w-5 text-primary" />
              Order Lifecycle Status Breakdown
            </h3>
            <p className="text-xs text-muted-foreground">Distribution of active and completed order stages</p>
          </div>
        </div>

        {/* Recharts Bar Chart with Direct Value Labels */}
        <div className="h-[260px] w-full pt-2">
          <ResponsiveContainer width="100%" height="100%">
            <BarChart data={statusChartData} margin={{ top: 25, right: 10, left: -10, bottom: 5 }}>
              <XAxis
                dataKey="name"
                axisLine={false}
                tickLine={false}
                tick={{ fontSize: 10, fill: 'var(--muted-foreground)' }}
              />
              <YAxis
                axisLine={false}
                tickLine={false}
                allowDecimals={false}
                tick={{ fontSize: 10, fill: 'var(--muted-foreground)' }}
              />
              <Tooltip
                content={({ active, payload }) => {
                  if (active && payload && payload.length) {
                    const d = payload[0].payload;
                    const pct = Math.min(100, Math.round((d.count / (totalOrders || 1)) * 100));
                    return (
                      <div className="bg-popover text-popover-foreground border border-border rounded-xl p-3 shadow-md text-xs font-sans">
                        <p className="font-semibold font-display mb-1">{d.name}</p>
                        <p className="text-primary font-semibold">Orders: {d.count}</p>
                        <p className="text-muted-foreground">Share: {pct}%</p>
                      </div>
                    );
                  }
                  return null;
                }}
              />
              <Bar dataKey="count" radius={[6, 6, 0, 0]}>
                {statusChartData.map((entry, idx) => (
                  <Cell key={`cell-${idx}`} fill={entry.color} />
                ))}
                <LabelList
                  dataKey="count"
                  position="top"
                  formatter={(val: any) => {
                    const num = Number(val);
                    if (!num) return '';
                    const pct = Math.min(100, Math.round((num / (totalOrders || 1)) * 100));
                    return `${num} (${pct}%)`;
                  }}
                  style={{ fontSize: 10, fontWeight: 700, fill: 'var(--foreground)' }}
                />
              </Bar>
            </BarChart>
          </ResponsiveContainer>
        </div>
      </div>

      {/* Top Products & Top Shops Grid */}
      <div className="grid gap-8 md:grid-cols-2">
        {/* Top Products */}
        <div className="space-y-4">
          <h3 className="font-bold font-display text-base text-foreground flex items-center gap-2">
            <ShoppingBag className="h-4 w-4 text-primary" />
            Top Selling Products by Revenue
          </h3>
          {topProducts.length === 0 ? (
            <p className="text-xs text-muted-foreground">No top product records found.</p>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Product</TableHead>
                  <TableHead className="text-right">Qty</TableHead>
                  <TableHead className="text-right">Revenue</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {topProducts.map((p) => (
                  <TableRow key={p.product_id}>
                    <TableCell className="font-medium text-xs text-foreground">{p.product_name}</TableCell>
                    <TableCell className="text-right text-xs text-muted-foreground">{p.quantity}</TableCell>
                    <TableCell className="text-right text-xs font-semibold text-primary">{formatIDR(p.revenue)}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </div>

        {/* Top Shops */}
        <div className="space-y-4">
          <h3 className="font-bold font-display text-base text-foreground flex items-center gap-2">
            <ShoppingBag className="h-4 w-4 text-primary" />
            Top Performing Shop Branches
          </h3>
          {topShops.length === 0 ? (
            <p className="text-xs text-muted-foreground">No shop revenue records found.</p>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Shop Branch</TableHead>
                  <TableHead className="text-right">Orders</TableHead>
                  <TableHead className="text-right">Revenue</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {topShops.map((s) => (
                  <TableRow key={s.shop_id}>
                    <TableCell className="font-medium text-xs text-foreground">{s.shop_name}</TableCell>
                    <TableCell className="text-right text-xs text-muted-foreground">{s.orders}</TableCell>
                    <TableCell className="text-right text-xs font-semibold text-primary">{formatIDR(s.revenue)}</TableCell>
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
