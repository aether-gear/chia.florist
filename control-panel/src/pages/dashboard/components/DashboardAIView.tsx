import React from 'react';
import { Sparkles, Percent } from 'lucide-react';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '../../../components/ui/table';
import { Badge } from '../../../components/ui/badge';
import { Skeleton } from '../../../components/ui/skeleton';
import type { StockoutRiskItem, DemandForecast } from '../../../models/Analytics';
import type { ProductStat } from '../../../models/Product';

interface DashboardAIViewProps {
  stockoutRisks: StockoutRiskItem[];
  productStats: ProductStat[];
  selectedForecastProductId: string;
  onSelectForecastProduct: (id: string) => void;
  forecastData: DemandForecast | null;
  forecastLoading: boolean;
  loading: boolean;
}

export const DashboardAIView: React.FC<DashboardAIViewProps> = ({
  stockoutRisks,
  productStats,
  selectedForecastProductId,
  onSelectForecastProduct,
  forecastData,
  forecastLoading,
  loading,
}) => {
  const getRiskBadge = (level: string) => {
    switch (level) {
      case 'CRITICAL':
        return <Badge variant="destructive" className="rounded-md uppercase">Critical</Badge>;
      case 'WARNING':
        return <Badge variant="secondary" className="bg-amber-500/15 text-amber-700 dark:text-amber-300 border-amber-500/20 rounded-md uppercase">Warning</Badge>;
      default:
        return <Badge variant="secondary" className="bg-primary/10 text-primary border-primary/20 rounded-md uppercase">Nominal</Badge>;
    }
  };

  const getConfidenceBadge = (tier: string) => {
    switch (tier?.toLowerCase()) {
      case 'high':
        return <Badge className="bg-primary/15 text-primary border-primary/20">High Confidence</Badge>;
      case 'medium':
        return <Badge className="bg-amber-500/15 text-amber-700 dark:text-amber-300 border-amber-500/20">Medium Confidence</Badge>;
      default:
        return <Badge variant="outline" className="text-muted-foreground">Baseline Model</Badge>;
    }
  };

  return (
    <div className="space-y-10">
      {/* 1. Interactive AI Demand Forecast Analyzer */}
      <div className="p-6 rounded-2xl bg-primary/5 border border-primary/15 space-y-6">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b border-primary/15">
          <div className="flex items-center gap-2.5">
            <div className="w-8 h-8 rounded-lg bg-primary/15 text-primary flex items-center justify-center">
              <Sparkles className="w-4 h-4" />
            </div>
            <div>
              <h3 className="font-bold font-display tracking-tight text-lg text-foreground">
                Predictive AI Demand Forecaster
              </h3>
              <p className="text-muted-foreground text-xs font-sans">
                Predicts next 7-day sales velocity using time-series moving averages and inventory burn rate.
              </p>
            </div>
          </div>

          <div className="flex items-center gap-2">
            <span className="text-xs font-medium text-muted-foreground">Bouquet:</span>
            <select
              value={selectedForecastProductId}
              onChange={(e) => onSelectForecastProduct(e.target.value)}
              className="h-9 rounded-xl border border-border bg-background px-3 py-1 text-xs font-medium shadow-sm transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring text-foreground max-w-[220px]"
            >
              {productStats.map((p) => (
                <option key={p.id} value={p.id}>
                  {p.name}
                </option>
              ))}
            </select>
          </div>
        </div>

        {/* Forecast Metric Cards */}
        {forecastLoading ? (
          <div className="grid gap-4 sm:grid-cols-3">
            {Array.from({ length: 3 }).map((_, i) => (
              <div key={i} className="p-4 rounded-xl bg-background/60 border border-border/40 space-y-2">
                <Skeleton className="h-4 w-28 bg-muted" />
                <Skeleton className="h-7 w-20 bg-muted" />
              </div>
            ))}
          </div>
        ) : forecastData ? (
          <div className="grid gap-4 sm:grid-cols-3">
            <div className="p-4 rounded-xl bg-background/80 border border-primary/20 space-y-1">
              <span className="text-[11px] font-sans font-semibold uppercase tracking-wider text-muted-foreground">
                Predicted 7d Demand
              </span>
              <div className="text-2xl font-bold font-display text-primary">
                {forecastData.predicted_units_sold_7d}{' '}
                <span className="text-xs font-normal text-muted-foreground font-sans">units</span>
              </div>
              <div className="pt-1">{getConfidenceBadge(forecastData.confidence_tier)}</div>
            </div>

            <div className="p-4 rounded-xl bg-background/80 border border-border/40 space-y-1">
              <span className="text-[11px] font-sans font-semibold uppercase tracking-wider text-muted-foreground">
                Historical 7d Velocity
              </span>
              <div className="text-2xl font-bold font-display text-foreground">
                {forecastData.historical_velocity_7d}{' '}
                <span className="text-xs font-normal text-muted-foreground font-sans">units sold</span>
              </div>
              <p className="text-[11px] text-muted-foreground pt-1">Actual historical baseline</p>
            </div>

            <div className="p-4 rounded-xl bg-background/80 border border-border/40 space-y-1">
              <span className="text-[11px] font-sans font-semibold uppercase tracking-wider text-muted-foreground">
                Current Stock In Hand
              </span>
              <div className="text-2xl font-bold font-display text-foreground">
                {forecastData.current_stock}{' '}
                <span className="text-xs font-normal text-muted-foreground font-sans">units</span>
              </div>
              <p className="text-[11px] text-muted-foreground pt-1">
                {forecastData.current_stock < forecastData.predicted_units_sold_7d ? (
                  <span className="text-rose-600 font-semibold">⚠️ Projected stockout within 7d</span>
                ) : (
                  <span className="text-primary font-semibold">✅ Sufficient stock coverage</span>
                )}
              </p>
            </div>
          </div>
        ) : (
          <p className="text-xs text-muted-foreground text-center py-4">Select a bouquet above to generate AI demand forecast.</p>
        )}
      </div>

      {/* 2. Stockout Risk Assessment Matrix */}
      <div className="space-y-4">
        <div className="flex flex-row items-center justify-between pb-2 border-b border-border/60">
          <div>
            <h3 className="font-bold font-display tracking-tight text-lg text-foreground">
              Inventory Stockout Risk Matrix
            </h3>
            <p className="text-muted-foreground text-xs font-sans">
              Algorithmic run-out prediction based on supplier lead times and 7-day sales burn rates.
            </p>
          </div>
        </div>

        <div>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Product / Bouquet</TableHead>
                <TableHead>Shop Branch</TableHead>
                <TableHead>Available Stock</TableHead>
                <TableHead>7d Burn Rate</TableHead>
                <TableHead>Lead Time</TableHead>
                <TableHead>Est. Run-Out</TableHead>
                <TableHead>Risk Level</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {loading ? (
                Array.from({ length: 4 }).map((_, i) => (
                  <TableRow key={i}>
                    <TableCell><Skeleton className="h-4 w-36 bg-muted" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-28 bg-muted" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-16 bg-muted" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-16 bg-muted" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-16 bg-muted" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-20 bg-muted" /></TableCell>
                    <TableCell><Skeleton className="h-6 w-16 bg-muted" /></TableCell>
                  </TableRow>
                ))
              ) : stockoutRisks.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={7} className="text-center text-xs text-muted-foreground py-8">
                    No critical stockout risks detected across boutique branches.
                  </TableCell>
                </TableRow>
              ) : (
                stockoutRisks.map((risk, idx) => (
                  <TableRow key={`${risk.product_id}-${risk.shop_id}-${idx}`} className="hover:bg-muted/30">
                    <TableCell className="font-medium text-xs text-foreground">
                      {risk.product_name}
                    </TableCell>
                    <TableCell className="text-xs text-muted-foreground">
                      {risk.shop_name}
                    </TableCell>
                    <TableCell className="font-mono text-xs text-foreground font-semibold">
                      {risk.available_stock}
                    </TableCell>
                    <TableCell className="text-xs text-muted-foreground">
                      {risk.stock_burn_rate_7d} / day
                    </TableCell>
                    <TableCell className="text-xs text-muted-foreground">
                      {risk.supplier_lead_time_days} days
                    </TableCell>
                    <TableCell className="text-xs font-semibold text-foreground">
                      {risk.estimated_days_to_stockout} days
                    </TableCell>
                    <TableCell>{getRiskBadge(risk.risk_level)}</TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </div>
      </div>

      {/* 3. Product Conversion & Margin Intelligence */}
      <div className="space-y-4">
        <div className="flex flex-row items-center justify-between pb-2 border-b border-border/60">
          <div>
            <h3 className="font-bold font-display tracking-tight text-lg text-foreground">
              Conversion & Margin Performance Intelligence
            </h3>
            <p className="text-muted-foreground text-xs font-sans">
              Real-time conversion signals (detail views to purchase) and gross profit margins.
            </p>
          </div>
        </div>

        <div>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Bouquet</TableHead>
                <TableHead>Catalog Views</TableHead>
                <TableHead>30d Sales</TableHead>
                <TableHead>Conversion Rate</TableHead>
                <TableHead>Gross Margin</TableHead>
                <TableHead>Revenue Contribution</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {loading ? (
                Array.from({ length: 4 }).map((_, i) => (
                  <TableRow key={i}>
                    <TableCell><Skeleton className="h-4 w-36 bg-muted" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-16 bg-muted" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-16 bg-muted" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-20 bg-muted" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-20 bg-muted" /></TableCell>
                    <TableCell><Skeleton className="h-4 w-20 bg-muted" /></TableCell>
                  </TableRow>
                ))
              ) : productStats.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={6} className="text-center text-xs text-muted-foreground py-8">
                    No conversion statistics recorded.
                  </TableCell>
                </TableRow>
              ) : (
                productStats.slice(0, 8).map((p) => {
                  const convPct = (p.conversion_rate * 100).toFixed(1);
                  const marginPct = p.gross_margin_pct != null ? `${p.gross_margin_pct.toFixed(1)}%` : '—';
                  return (
                    <TableRow key={p.id} className="hover:bg-muted/30">
                      <TableCell className="font-medium text-xs text-foreground">
                        {p.name}
                      </TableCell>
                      <TableCell className="font-mono text-xs text-muted-foreground">
                        {p.view_count.toLocaleString()}
                      </TableCell>
                      <TableCell className="font-mono text-xs text-foreground font-semibold">
                        {p.sales_velocity_30d}
                      </TableCell>
                      <TableCell>
                        <span className="inline-flex items-center gap-1 font-semibold text-xs text-primary">
                          <Percent className="w-3 h-3" />
                          {convPct}%
                        </span>
                      </TableCell>
                      <TableCell className="font-mono text-xs text-foreground">
                        {marginPct}
                      </TableCell>
                      <TableCell className="font-display text-xs font-semibold text-primary">
                        {p.revenue_contribution_percentage}%
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

export default DashboardAIView;
