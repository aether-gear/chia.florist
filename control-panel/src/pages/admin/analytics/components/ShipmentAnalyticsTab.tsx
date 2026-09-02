import { Truck, CheckCircle2, AlertCircle, Clock } from 'lucide-react';
import type { ShipmentMetricsResponse } from '../../../../models/Analytics';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '../../../../components/ui/table';
import { Skeleton } from '../../../../components/ui/skeleton';

interface Props {
  data: ShipmentMetricsResponse | null;
  loading: boolean;
}

export default function ShipmentAnalyticsTab({ data, loading }: Props) {
  const formatIDR = (val?: number | null) => {
    if (val === undefined || val === null || isNaN(val)) return 'Rp 0';
    return `Rp ${val.toLocaleString('id-ID')}`;
  };

  const formatFulfillmentTime = (sec?: number) => {
    if (!sec || isNaN(sec)) return 'N/A';
    const hours = Math.floor(sec / 3600);
    const mins = Math.round((sec % 3600) / 60);
    if (hours > 0) return `${hours}h ${mins}m`;
    return `${mins} mins`;
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
  const couriers = Array.isArray(data?.couriers) ? data.couriers : [];

  return (
    <div className="space-y-8 animate-in fade-in slide-in-from-left-4 duration-300">
      {/* Shipment KPI Grid - SecurityPage Borderless Style */}
      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4 border-b border-border/60 pb-8">
        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <span className="text-sm font-bold font-display text-foreground">Total Shipments</span>
            <Truck className="h-4 w-4 text-primary" />
          </div>
          <div className="text-3xl font-bold tracking-tight text-foreground">{summary?.total || 0}</div>
          <p className="text-xs text-muted-foreground">Delivered: {summary?.delivered || 0} packages</p>
        </div>

        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <span className="text-sm font-bold font-display text-foreground">Delivery Success</span>
            <CheckCircle2 className="h-4 w-4 text-emerald-600" />
          </div>
          <div className="text-3xl font-bold tracking-tight text-emerald-600 dark:text-emerald-400">
            {((summary?.delivery_rate || 0) * 100).toFixed(1)}%
          </div>
          <p className="text-xs text-muted-foreground">Successful logistics completion</p>
        </div>

        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <span className="text-sm font-bold font-display text-foreground">Avg Fulfillment Speed</span>
            <Clock className="h-4 w-4 text-primary" />
          </div>
          <div className="text-3xl font-bold tracking-tight text-foreground">
            {formatFulfillmentTime(summary?.avg_fulfillment_sec)}
          </div>
          <p className="text-xs text-muted-foreground">Order to handoff delivery</p>
        </div>

        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <span className="text-sm font-bold font-display text-foreground">Issues & Returns</span>
            <AlertCircle className="h-4 w-4 text-rose-500" />
          </div>
          <div className="text-3xl font-bold tracking-tight text-rose-600 dark:text-rose-400">
            {(summary?.failed || 0) + (summary?.returned || 0)}
          </div>
          <p className="text-xs text-rose-500 font-medium">
            {summary?.failed || 0} Failed / {summary?.returned || 0} Returned
          </p>
        </div>
      </div>

      {/* Couriers Performance Table */}
      <div className="space-y-4">
        <h3 className="font-bold font-display text-lg text-foreground flex items-center gap-2">
          <Truck className="h-5 w-5 text-primary" />
          Courier Service Performance Breakdown
        </h3>
        {couriers.length === 0 ? (
          <p className="text-xs text-muted-foreground">No courier performance records found.</p>
        ) : (
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Courier</TableHead>
                <TableHead>Service</TableHead>
                <TableHead className="text-right">Shipments</TableHead>
                <TableHead className="text-right">Delivery Success</TableHead>
                <TableHead className="text-right">Avg Shipping Cost</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {couriers.map((item, idx) => (
                <TableRow key={`${item.courier}-${item.service}-${idx}`}>
                  <TableCell className="font-bold text-xs text-foreground uppercase">{item.courier}</TableCell>
                  <TableCell className="text-xs text-muted-foreground uppercase">{item.service}</TableCell>
                  <TableCell className="text-right text-xs text-muted-foreground">{item.count}</TableCell>
                  <TableCell className="text-right text-xs font-semibold text-emerald-600">
                    {((item.delivery_rate || 0) * 100).toFixed(1)}%
                  </TableCell>
                  <TableCell className="text-right text-xs font-medium text-foreground">{formatIDR(item.avg_cost)}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        )}
      </div>
    </div>
  );
}
