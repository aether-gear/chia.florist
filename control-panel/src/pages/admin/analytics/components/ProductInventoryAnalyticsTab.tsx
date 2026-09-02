import { Link } from 'react-router-dom';
import { Package, AlertTriangle, TrendingUp, DollarSign, Sparkles } from 'lucide-react';
import type { InventoryMetricsResponse, ProductMetricsResponse } from '../../../../models/Analytics';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '../../../../components/ui/table';
import { Skeleton } from '../../../../components/ui/skeleton';

interface Props {
  inventoryData: InventoryMetricsResponse | null;
  productData: ProductMetricsResponse | null;
  loading: boolean;
}

export default function ProductInventoryAnalyticsTab({
  inventoryData,
  productData,
  loading,
}: Props) {
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

  const topByRevenue = Array.isArray(productData?.top_by_revenue) ? productData.top_by_revenue : [];
  const topByVolume = Array.isArray(productData?.top_by_volume) ? productData.top_by_volume : [];

  return (
    <div className="space-y-8 animate-in fade-in slide-in-from-left-4 duration-300">
      {/* Inventory & Financial KPI Grid - SecurityPage Borderless Style */}
      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4 border-b border-border/60 pb-8">
        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <span className="text-sm font-bold font-display text-foreground">Total Inventory Stock</span>
            <Package className="h-4 w-4 text-primary" />
          </div>
          <div className="text-3xl font-bold tracking-tight text-foreground">{inventoryData?.total_stock || 0}</div>
          <p className="text-xs text-muted-foreground">
            Available: <span className="font-semibold text-primary">{inventoryData?.total_available || 0}</span> (Reserved: {inventoryData?.total_reserved || 0})
          </p>
        </div>

        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <span className="text-sm font-bold font-display text-foreground">Low Stock Warnings</span>
            <AlertTriangle className="h-4 w-4 text-amber-500" />
          </div>
          <div className="text-3xl font-bold tracking-tight text-amber-600 dark:text-amber-400">
            {inventoryData?.low_stock_count || 0}
          </div>
          <p className="text-xs text-amber-600 font-medium">Stockouts: {inventoryData?.stockout_count || 0} products</p>
        </div>

        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <span className="text-sm font-bold font-display text-foreground">Avg Gross Margin</span>
            <DollarSign className="h-4 w-4 text-emerald-600" />
          </div>
          <div className="text-3xl font-bold tracking-tight text-emerald-600 dark:text-emerald-400">
            {typeof topByRevenue[0]?.gross_margin_pct === 'number' ? `${topByRevenue[0].gross_margin_pct.toFixed(1)}%` : 'N/A'}
          </div>
          <p className="text-xs text-muted-foreground">High-profit margin line items</p>
        </div>

        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <span className="text-sm font-bold font-display text-foreground">Invoice Void Rate</span>
            <TrendingUp className="h-4 w-4 text-muted-foreground" />
          </div>
          <div className="text-3xl font-bold tracking-tight text-foreground">
            {((productData?.invoice_void_rate || 0) * 100).toFixed(1)}%
          </div>
          <p className="text-xs text-muted-foreground">Canceled / voided invoices</p>
        </div>
      </div>

      {/* AI Stockout & Demand Callout Banner */}
      <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between p-4 rounded-2xl bg-primary/5 border border-primary/20 gap-4">
        <div className="flex items-center gap-3">
          <div className="p-2.5 rounded-xl bg-primary/10 text-primary">
            <Sparkles className="h-5 w-5" />
          </div>
          <div>
            <h4 className="font-bold font-display text-sm text-foreground">AI Stockout Risk Scanner & Demand Forecasting Available</h4>
            <p className="text-xs text-muted-foreground">
              Evaluate real-time inventory depletion risks and projected 7-day SKU volume forecasts.
            </p>
          </div>
        </div>
        <Link
          to="/admin/analytics?tab=ai-intelligence"
          className="text-xs font-semibold text-primary hover:underline flex items-center gap-1 shrink-0"
        >
          Open AI Intelligence →
        </Link>
      </div>

      {/* Product Tables: By Revenue vs By Volume */}
      <div className="grid gap-8 md:grid-cols-2">
        {/* Top by Revenue */}
        <div className="space-y-4">
          <h3 className="font-bold font-display text-base text-foreground flex items-center gap-2">
            <DollarSign className="h-4 w-4 text-primary" />
            Top Ranked Products by Revenue
          </h3>
          {topByRevenue.length === 0 ? (
            <p className="text-xs text-muted-foreground">No product metrics found.</p>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Product</TableHead>
                  <TableHead className="text-right">Revenue</TableHead>
                  <TableHead className="text-right">Velocity (7d/30d)</TableHead>
                  <TableHead className="text-right">Margin</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {topByRevenue.map((p) => (
                  <TableRow key={p.product_id}>
                    <TableCell className="font-medium text-xs text-foreground">{p.product_name}</TableCell>
                    <TableCell className="text-right text-xs font-semibold text-primary">{formatIDR(p.revenue)}</TableCell>
                    <TableCell className="text-right text-xs text-muted-foreground">
                      {p.sales_velocity_7d} / {p.sales_velocity_30d}
                    </TableCell>
                    <TableCell className="text-right text-xs text-emerald-600 font-medium">
                      {p.gross_margin_pct}%
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </div>

        {/* Top by Volume */}
        <div className="space-y-4">
          <h3 className="font-bold font-display text-base text-foreground flex items-center gap-2">
            <TrendingUp className="h-4 w-4 text-primary" />
            Top Ranked Products by Volume Sold
          </h3>
          {topByVolume.length === 0 ? (
            <p className="text-xs text-muted-foreground">No volume metrics found.</p>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Product</TableHead>
                  <TableHead className="text-right">Units Sold</TableHead>
                  <TableHead className="text-right">Velocity (7d/30d)</TableHead>
                  <TableHead className="text-right">Revenue</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {topByVolume.map((p) => (
                  <TableRow key={p.product_id}>
                    <TableCell className="font-medium text-xs text-foreground">{p.product_name}</TableCell>
                    <TableCell className="text-right text-xs font-bold text-foreground">{p.units_sold}</TableCell>
                    <TableCell className="text-right text-xs text-muted-foreground">
                      {p.sales_velocity_7d} / {p.sales_velocity_30d}
                    </TableCell>
                    <TableCell className="text-right text-xs font-semibold text-primary">{formatIDR(p.revenue)}</TableCell>
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
