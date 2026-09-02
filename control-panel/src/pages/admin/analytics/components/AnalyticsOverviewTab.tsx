import { DollarSign, ShoppingBag, CreditCard, Truck, AlertTriangle } from 'lucide-react';
import { ResponsiveContainer, AreaChart, Area, XAxis, YAxis, Tooltip } from 'recharts';
import type {
  OrderMetricsResponse,
  PaymentMetricsResponse,
  ShipmentMetricsResponse,
  InventoryMetricsResponse,
} from '../../../../models/Analytics';
import { Skeleton } from '../../../../components/ui/skeleton';

interface Props {
  orderData: OrderMetricsResponse | null;
  paymentData: PaymentMetricsResponse | null;
  shipmentData: ShipmentMetricsResponse | null;
  inventoryData: InventoryMetricsResponse | null;
  loading: boolean;
}

export default function AnalyticsOverviewTab({
  orderData,
  paymentData,
  shipmentData,
  inventoryData,
  loading,
}: Props) {
  const formatIDR = (val?: number | null) => {
    if (val === undefined || val === null || isNaN(val)) return 'Rp 0';
    return `Rp ${val.toLocaleString('id-ID')}`;
  };

  const timeSeries = Array.isArray(orderData?.time_series) ? orderData.time_series : [];
  const chartPoints = timeSeries.map(p => ({
    date: p && p.date ? new Date(p.date).toLocaleDateString(undefined, { month: 'short', day: 'numeric' }) : '',
    gmv: p && typeof p.gmv === 'number' ? p.gmv : 0,
    orders: p && typeof p.order_count === 'number' ? p.order_count : 0,
  }));

  if (loading) {
    return (
      <div className="space-y-6">
        <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-4 border-b border-border/60 pb-8">
          {Array.from({ length: 4 }).map((_, i) => (
            <div key={i} className="space-y-2">
              <Skeleton className="h-4 w-28 bg-muted animate-pulse" />
              <Skeleton className="h-8 w-36 bg-muted animate-pulse" />
              <Skeleton className="h-3 w-40 bg-muted animate-pulse" />
            </div>
          ))}
        </div>
        <Skeleton className="h-72 rounded-2xl bg-muted animate-pulse" />
      </div>
    );
  }

  const orderSummary = orderData?.summary;
  const paymentSummary = paymentData?.summary;
  const shipmentSummary = shipmentData?.summary;

  return (
    <div className="space-y-8 animate-in fade-in slide-in-from-left-4 duration-300">
      {/* Overview KPI Cards - SecurityPage Borderless Style */}
      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4 border-b border-border/60 pb-8">
        {/* Gross Merchandise Value */}
        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <span className="text-sm font-bold font-display text-foreground">Gross Merchandise Value</span>
            <DollarSign className="h-4 w-4 text-primary" />
          </div>
          <div className="text-3xl font-bold tracking-tight text-foreground">
            {formatIDR(orderSummary?.total_gmv)}
          </div>
          <p className="text-xs text-muted-foreground">
            Revenue: <span className="font-semibold text-primary">{formatIDR(orderSummary?.total_revenue)}</span>
          </p>
        </div>

        {/* Total Orders */}
        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <span className="text-sm font-bold font-display text-foreground">Total Orders</span>
            <ShoppingBag className="h-4 w-4 text-primary" />
          </div>
          <div className="text-3xl font-bold tracking-tight text-foreground">
            {orderSummary?.total_orders || 0}
          </div>
          <p className="text-xs text-muted-foreground">
            AOV: <span className="font-semibold text-foreground">{formatIDR(orderSummary?.aov)}</span>
          </p>
        </div>

        {/* Payment Success */}
        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <span className="text-sm font-bold font-display text-foreground">Payment Success</span>
            <CreditCard className="h-4 w-4 text-emerald-600" />
          </div>
          <div className="text-3xl font-bold tracking-tight text-emerald-600 dark:text-emerald-400">
            {((paymentSummary?.payment_success_rate || 0) * 100).toFixed(1)}%
          </div>
          <p className="text-xs text-muted-foreground">
            Total Paid: <span className="font-semibold text-foreground">{formatIDR(paymentSummary?.total_paid)}</span>
          </p>
        </div>

        {/* Delivery Rate */}
        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <span className="text-sm font-bold font-display text-foreground">Delivery Rate</span>
            <Truck className="h-4 w-4 text-blue-600" />
          </div>
          <div className="text-3xl font-bold tracking-tight text-foreground">
            {((shipmentSummary?.delivery_rate || 0) * 100).toFixed(1)}%
          </div>
          <p className="text-xs text-amber-600 font-medium flex items-center gap-1">
            <AlertTriangle className="h-3 w-3 inline shrink-0" />
            Low Stock Alert: {inventoryData?.low_stock_count || 0} items
          </p>
        </div>
      </div>

      {/* GMV & Orders Trend Chart */}
      <div className="space-y-4">
        <div className="flex items-center justify-between border-b border-border/40 pb-3">
          <div>
            <h3 className="font-bold font-display tracking-tight text-lg text-foreground">Sales & GMV Performance Trend</h3>
            <p className="text-xs text-muted-foreground">Historical order volume and revenue generation</p>
          </div>
        </div>
        <div className="h-[280px] w-full">
          {chartPoints.length === 0 ? (
            <div className="h-full flex items-center justify-center text-xs text-muted-foreground">
              No time-series data available for selected period.
            </div>
          ) : (
            <ResponsiveContainer width="100%" height="100%">
              <AreaChart data={chartPoints} margin={{ top: 10, right: 10, left: 0, bottom: 0 }}>
                <defs>
                  <linearGradient id="gmvGrad" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="hsl(142.4, 71.8%, 29.2%)" stopOpacity={0.4} />
                    <stop offset="95%" stopColor="hsl(142.4, 71.8%, 29.2%)" stopOpacity={0} />
                  </linearGradient>
                </defs>
                <XAxis dataKey="date" tickLine={false} axisLine={false} tick={{ fontSize: 11, fill: 'var(--muted-foreground)' }} />
                <YAxis
                  axisLine={false}
                  tickLine={false}
                  tickFormatter={(val) => `Rp ${(val / 1000000).toFixed(1)}M`}
                  tick={{ fontSize: 11, fill: 'var(--muted-foreground)' }}
                />
                <Tooltip
                  content={({ active, payload }) => {
                    if (active && payload && payload.length) {
                      const data = payload[0].payload;
                      return (
                        <div className="bg-popover text-popover-foreground border border-border rounded-xl p-3 shadow-md text-xs font-sans">
                          <p className="font-semibold font-display mb-1">{data.date}</p>
                          <p className="text-primary font-semibold">GMV: {formatIDR(data.gmv)}</p>
                          <p className="text-muted-foreground">Orders: {data.orders}</p>
                        </div>
                      );
                    }
                    return null;
                  }}
                />
                <Area type="monotone" dataKey="gmv" stroke="hsl(142.4, 71.8%, 29.2%)" strokeWidth={2.5} fillOpacity={1} fill="url(#gmvGrad)" />
              </AreaChart>
            </ResponsiveContainer>
          )}
        </div>
      </div>
    </div>
  );
}
