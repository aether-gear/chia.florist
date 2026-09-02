import React, { useMemo } from 'react';
import {
  ResponsiveContainer,
  AreaChart,
  Area,
  XAxis,
  YAxis,
  Tooltip,
  CartesianGrid,
} from 'recharts';
import { ShoppingBag, Eye, CheckCircle2, Clock, Truck, XCircle, Warehouse } from 'lucide-react';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '../../../components/ui/table';
import { Badge } from '../../../components/ui/badge';
import { Button } from '../../../components/ui/button';
import { Skeleton } from '../../../components/ui/skeleton';
import type { OrderMetricsResponse } from '../../../models/Analytics';
import type { Order } from '../../../models/Order';
import type { ProductStat } from '../../../models/Product';

interface DashboardEcommerceViewProps {
  orderData: OrderMetricsResponse | null;
  productStats: ProductStat[];
  recentOrders: Order[];
  loading: boolean;
  onInspectOrder: (order: Order) => void;
}

export const DashboardEcommerceView: React.FC<DashboardEcommerceViewProps> = ({
  orderData,
  productStats,
  recentOrders,
  loading,
  onInspectOrder,
}) => {
  const formatCurrency = (val: number) => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      maximumFractionDigits: 0,
    }).format(val);
  };

  const chartData = useMemo(() => {
    if (!orderData?.time_series || orderData.time_series.length === 0) return [];
    return orderData.time_series.map((pt) => ({
      date: pt.date ? pt.date.slice(5) : '',
      gmv: pt.gmv || 0,
      orders: pt.order_count || 0,
    }));
  }, [orderData]);

  const topProducts = useMemo(() => {
    if (productStats && productStats.length > 0) {
      return [...productStats].sort((a, b) => b.sales_velocity_30d - a.sales_velocity_30d).slice(0, 5);
    }
    return [];
  }, [productStats]);

  const summary = orderData?.summary;

  const getOrderStatusBadge = (status: string) => {
    const s = status.toLowerCase();
    if (s === 'delivered' || s === 'completed') {
      return <Badge variant="secondary" className="bg-primary/10 text-primary border-primary/20">Delivered</Badge>;
    }
    if (s === 'shipped') {
      return <Badge variant="secondary" className="bg-sky-500/10 text-sky-600 dark:text-sky-400 border-sky-500/20">Shipped</Badge>;
    }
    if (s === 'processing' || s === 'confirmed') {
      return <Badge variant="secondary" className="bg-amber-500/10 text-amber-600 dark:text-amber-400 border-amber-500/20">Processing</Badge>;
    }
    if (s === 'cancelled' || s === 'rejected') {
      return <Badge variant="destructive" className="rounded-md">Cancelled</Badge>;
    }
    return <Badge variant="outline" className="text-muted-foreground">{status}</Badge>;
  };

  return (
    <div className="space-y-10">
      {/* 1. Order Status Pipeline Counters */}
      <div className="grid gap-4 grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 pb-6 border-b border-border/60">
        <div className="p-3.5 rounded-xl bg-muted/40 border border-border/40 flex items-center gap-3">
          <div className="w-8 h-8 rounded-lg bg-amber-500/10 text-amber-600 dark:text-amber-400 flex items-center justify-center shrink-0">
            <Clock className="w-4 h-4" />
          </div>
          <div className="min-w-0">
            <p className="text-[11px] text-muted-foreground font-sans uppercase tracking-wider font-semibold">Pending</p>
            <p className="text-lg font-bold font-display text-foreground">{summary?.pending_count ?? 0}</p>
          </div>
        </div>

        <div className="p-3.5 rounded-xl bg-muted/40 border border-border/40 flex items-center gap-3">
          <div className="w-8 h-8 rounded-lg bg-indigo-500/10 text-indigo-600 dark:text-indigo-400 flex items-center justify-center shrink-0">
            <ShoppingBag className="w-4 h-4" />
          </div>
          <div className="min-w-0">
            <p className="text-[11px] text-muted-foreground font-sans uppercase tracking-wider font-semibold">Processing</p>
            <p className="text-lg font-bold font-display text-foreground">{summary?.processing_count ?? 0}</p>
          </div>
        </div>

        <div className="p-3.5 rounded-xl bg-muted/40 border border-border/40 flex items-center gap-3">
          <div className="w-8 h-8 rounded-lg bg-sky-500/10 text-sky-600 dark:text-sky-400 flex items-center justify-center shrink-0">
            <Truck className="w-4 h-4" />
          </div>
          <div className="min-w-0">
            <p className="text-[11px] text-muted-foreground font-sans uppercase tracking-wider font-semibold">In Transit</p>
            <p className="text-lg font-bold font-display text-foreground">{summary?.shipped_count ?? 0}</p>
          </div>
        </div>

        <div className="p-3.5 rounded-xl bg-muted/40 border border-border/40 flex items-center gap-3">
          <div className="w-8 h-8 rounded-lg bg-primary/10 text-primary flex items-center justify-center shrink-0">
            <CheckCircle2 className="w-4 h-4" />
          </div>
          <div className="min-w-0">
            <p className="text-[11px] text-muted-foreground font-sans uppercase tracking-wider font-semibold">Delivered</p>
            <p className="text-lg font-bold font-display text-foreground">{summary?.delivered_count ?? 0}</p>
          </div>
        </div>

        <div className="p-3.5 rounded-xl bg-muted/40 border border-border/40 flex items-center gap-3">
          <div className="w-8 h-8 rounded-lg bg-rose-500/10 text-rose-600 dark:text-rose-400 flex items-center justify-center shrink-0">
            <XCircle className="w-4 h-4" />
          </div>
          <div className="min-w-0">
            <p className="text-[11px] text-muted-foreground font-sans uppercase tracking-wider font-semibold">Cancelled</p>
            <p className="text-lg font-bold font-display text-foreground">{summary?.cancelled_count ?? 0}</p>
          </div>
        </div>
      </div>

      {/* 2. Revenue Trend Time Series & Top Bouquets */}
      <div className="grid gap-10 lg:grid-cols-7 pb-8 border-b border-border/60">
        {/* Left: Revenue Area Chart */}
        <div className="lg:col-span-4 space-y-4">
          <div className="pb-2 border-b border-border/60 flex items-center justify-between">
            <div>
              <h3 className="font-bold font-display tracking-tight text-lg text-foreground">
                Revenue & Sales Trajectory
              </h3>
              <p className="text-muted-foreground text-xs font-sans">
                Daily transaction volume and gross merchandise value.
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
                No revenue records found for this period.
              </div>
            ) : (
              <ResponsiveContainer width="100%" height="100%">
                <AreaChart data={chartData} margin={{ top: 10, right: 10, left: -20, bottom: 0 }}>
                  <defs>
                    <linearGradient id="colorGmv" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="5%" stopColor="hsl(var(--primary))" stopOpacity={0.4} />
                      <stop offset="95%" stopColor="hsl(var(--primary))" stopOpacity={0.0} />
                    </linearGradient>
                  </defs>
                  <CartesianGrid strokeDasharray="3 3" stroke="currentColor" className="text-border/40" />
                  <XAxis dataKey="date" tick={{ fontSize: 10, fill: 'var(--muted-foreground)' }} axisLine={false} tickLine={false} />
                  <YAxis tick={{ fontSize: 10, fill: 'var(--muted-foreground)' }} axisLine={false} tickLine={false} />
                  <Tooltip
                    content={({ active, payload }) => {
                      if (!active || !payload?.length) return null;
                      const data = payload[0].payload;
                      return (
                        <div className="bg-popover text-popover-foreground border border-border rounded-xl p-3 shadow-lg text-xs font-sans">
                          <p className="font-bold font-display mb-1">{data.date}</p>
                          <p className="text-muted-foreground">
                            GMV: <span className="font-semibold text-primary">{formatCurrency(data.gmv)}</span>
                          </p>
                          <p className="text-muted-foreground">
                            Orders: <span className="font-semibold text-foreground">{data.orders}</span>
                          </p>
                        </div>
                      );
                    }}
                  />
                  <Area
                    type="monotone"
                    dataKey="gmv"
                    stroke="hsl(var(--primary))"
                    strokeWidth={2}
                    fillOpacity={1}
                    fill="url(#colorGmv)"
                  />
                </AreaChart>
              </ResponsiveContainer>
            )}
          </div>
        </div>

        {/* Right: Top Selling Bouquets */}
        <div className="lg:col-span-3 space-y-4">
          <div className="pb-2 border-b border-border/60">
            <h3 className="font-bold font-display tracking-tight text-lg text-foreground">
              Top Selling Bouquets
            </h3>
            <p className="text-muted-foreground text-xs font-sans">
              Highest sales velocity over the past 30 days.
            </p>
          </div>

          <div className="space-y-2.5">
            {loading ? (
              Array.from({ length: 4 }).map((_, i) => (
                <div key={i} className="p-3 rounded-xl bg-muted/20 space-y-2">
                  <Skeleton className="h-4 w-40 bg-muted" />
                  <Skeleton className="h-3 w-28 bg-muted" />
                </div>
              ))
            ) : topProducts.length === 0 ? (
              <p className="text-xs text-muted-foreground py-6 text-center">No product performance records found.</p>
            ) : (
              topProducts.map((p, idx) => (
                <div
                  key={p.id || idx}
                  className="flex items-center justify-between p-3 rounded-xl bg-muted/30 hover:bg-muted/50 border border-border/30 transition-colors"
                >
                  <div className="min-w-0 flex-1 pr-3">
                    <p className="text-xs font-bold text-foreground truncate">{p.name}</p>
                    <p className="text-[11px] text-muted-foreground">
                      {formatCurrency(p.price)} • {p.sales_velocity_30d} units sold
                    </p>
                  </div>
                  <div className="text-right shrink-0">
                    <span className="text-xs font-bold font-display text-primary">
                      {p.revenue_contribution_percentage}%
                    </span>
                    <p className="text-[10px] text-muted-foreground">rev. share</p>
                  </div>
                </div>
              ))
            )}
          </div>
        </div>
      </div>

      {/* 3. Live Orders Stream Table */}
      <div className="space-y-4">
        <div className="flex flex-row items-center justify-between pb-2 border-b border-border/60">
          <div>
            <h3 className="font-bold font-display tracking-tight text-lg text-foreground">
              Recent Customer Orders
            </h3>
            <p className="text-muted-foreground text-xs font-sans">
              Live order stream arriving across all flower boutiques.
            </p>
          </div>
        </div>

        <div>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Order #</TableHead>
                <TableHead>Recipient / Customer</TableHead>
                <TableHead>Shop Branch</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Total</TableHead>
                <TableHead className="text-right">Action</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {loading ? (
                Array.from({ length: 5 }).map((_, i) => (
                  <TableRow key={i}>
                    <TableCell><Skeleton className="h-4 w-20 bg-muted" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-32 bg-muted" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-24 bg-muted" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-16 bg-muted" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-20 bg-muted" /></TableCell>
                    <TableCell className="text-right"><Skeleton className="h-7 w-16 bg-muted ml-auto" /></TableCell>
                  </TableRow>
                ))
              ) : recentOrders.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={6} className="text-center text-xs text-muted-foreground py-8">
                    No orders recorded yet.
                  </TableCell>
                </TableRow>
              ) : (
                recentOrders.map((order) => {
                  const shopName = order.items?.[0]?.shop_name || 'Main Boutique';
                  const receiver = order.address?.receiver_name || 'Valued Customer';
                  return (
                    <TableRow key={order.id} className="hover:bg-muted/30">
                      <TableCell className="font-mono text-xs font-bold text-foreground">
                        {order.number}
                      </TableCell>
                      <TableCell className="text-xs text-foreground font-medium">
                        {receiver}
                      </TableCell>
                      <TableCell className="text-xs text-muted-foreground">
                        <span className="inline-flex items-center gap-1">
                          <Warehouse className="w-3 h-3 text-muted-foreground/70" />
                          {shopName}
                        </span>
                      </TableCell>
                      <TableCell>{getOrderStatusBadge(order.status)}</TableCell>
                      <TableCell className="font-display font-semibold text-xs text-foreground">
                        {formatCurrency(order.total)}
                      </TableCell>
                      <TableCell className="text-right">
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => onInspectOrder(order)}
                          className="h-8 text-xs font-semibold text-primary hover:text-primary hover:bg-primary/10 rounded-lg gap-1"
                        >
                          <Eye className="w-3.5 h-3.5" />
                          Inspect
                        </Button>
                      </TableCell>
                    </TableRow>
                  );
                })
              )}
            </TableBody>
          </Table>
        </div>
      </div>
    </div>
  );
};

export default DashboardEcommerceView;
