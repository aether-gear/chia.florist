import { useState, useEffect } from 'react';
import {
  Sparkles,
  AlertTriangle,
  Clock,
  TrendingUp,
  RefreshCw,
  Cpu,
  Package,
  Store,
  ShieldCheck,
  Zap,
} from 'lucide-react';
import {
  ResponsiveContainer,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip,
  Cell,
} from 'recharts';
import { useStockoutRisksViewModel } from '../../../../viewmodels/useStockoutRisksViewModel';
import { useDemandForecastViewModel } from '../../../../viewmodels/useDemandForecastViewModel';
import { useProductsViewModel } from '../../../../viewmodels/useProductsViewModel';
import type { StockoutRiskItem } from '../../../../models/Analytics';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '../../../../components/ui/table';
import { Badge } from '../../../../components/ui/badge';
import { Button } from '../../../../components/ui/button';
import { Skeleton } from '../../../../components/ui/skeleton';
import SearchInput from '../../../../components/SearchInput';
import EmptyState from '../../../../components/EmptyState';

interface Props {
  shopId?: string;
  onSelectProduct?: (productId: string) => void;
}

export default function AIForecastsAnalyticsTab({ shopId: propShopId }: Props) {
  const {
    shopId,
    setShopId,
    riskFilter,
    setRiskFilter,
    searchQuery,
    setSearchQuery,
    filteredRisks,
    criticalCount,
    warningCount,
    avgDaysToStockout,
    loading: risksLoading,
    error: risksError,
    refresh: refreshRisks,
  } = useStockoutRisksViewModel(propShopId);

  // Sync propShopId if provided
  useEffect(() => {
    if (propShopId !== undefined) {
      setShopId(propShopId);
    }
  }, [propShopId, setShopId]);

  // Product catalog for Live Demand Forecasting
  const { data: productsData } = useProductsViewModel();
  const products = productsData?.products || [];

  const [selectedProductId, setSelectedProductId] = useState<string>('');
  const {
    forecast,
    loading: forecastLoading,
    error: forecastError,
    fetchForecast,
  } = useDemandForecastViewModel();

  // Auto-select first product when list is loaded
  useEffect(() => {
    if (products.length > 0 && !selectedProductId) {
      const firstId = products[0].id;
      setSelectedProductId(firstId);
      fetchForecast(firstId, shopId || undefined);
    }
  }, [products, selectedProductId, shopId, fetchForecast]);

  const handleProductSelect = (productId: string) => {
    setSelectedProductId(productId);
    fetchForecast(productId, shopId || undefined);
  };

  const getRiskBadge = (level: string) => {
    switch (level?.toUpperCase()) {
      case 'CRITICAL':
        return (
          <Badge className="bg-destructive/15 text-destructive border-destructive/30 text-xs font-semibold px-2 py-0.5 animate-pulse">
            CRITICAL
          </Badge>
        );
      case 'WARNING':
        return (
          <Badge className="bg-amber-500/15 text-amber-600 dark:text-amber-400 border-amber-500/30 text-xs font-semibold px-2 py-0.5">
            WARNING
          </Badge>
        );
      case 'NORMAL':
      default:
        return (
          <Badge className="bg-emerald-500/15 text-emerald-600 dark:text-emerald-400 border-emerald-500/30 text-xs font-semibold px-2 py-0.5">
            NORMAL
          </Badge>
        );
    }
  };

  const getConfidenceBadge = (tier: string) => {
    switch (tier?.toLowerCase()) {
      case 'high':
        return <Badge className="bg-emerald-500/15 text-emerald-600 border-emerald-500/30 text-[10px] uppercase">High Confidence</Badge>;
      case 'medium':
        return <Badge className="bg-blue-500/15 text-blue-600 border-blue-500/30 text-[10px] uppercase">Medium Confidence</Badge>;
      default:
        return <Badge className="bg-muted text-muted-foreground border-border text-[10px] uppercase">Baseline Fallback</Badge>;
    }
  };

  // Forecast comparison chart data
  const forecastChartData = forecast
    ? [
        {
          name: 'Historical 7d Sold',
          units: forecast.historical_velocity_7d,
          color: 'hsl(var(--chart-1))',
        },
        {
          name: 'AI Forecast (Next 7d)',
          units: Math.round(forecast.predicted_units_sold_7d * 10) / 10,
          color: 'hsl(var(--primary))',
        },
      ]
    : [];

  return (
    <div className="space-y-10 animate-in fade-in duration-300">
      {/* 1. Top AI KPI Executive Summary Grid */}
      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4 border-b border-border/60 pb-8">
        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <span className="text-sm font-bold font-display text-foreground">Critical Stockout Alerts</span>
            <AlertTriangle className="h-4 w-4 text-destructive" />
          </div>
          <div className="text-3xl font-bold tracking-tight text-destructive">
            {risksLoading ? <Skeleton className="h-9 w-16 bg-muted" /> : criticalCount}
          </div>
          <p className="text-xs text-muted-foreground">
            {criticalCount > 0 ? (
              <span className="font-semibold text-destructive">Imminent stockout within 48-72h</span>
            ) : (
              <span className="text-emerald-600 font-medium">No critical stock risks</span>
            )}
          </p>
        </div>

        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <span className="text-sm font-bold font-display text-foreground">Warning Watchlist</span>
            <Clock className="h-4 w-4 text-amber-500" />
          </div>
          <div className="text-3xl font-bold tracking-tight text-amber-600 dark:text-amber-400">
            {risksLoading ? <Skeleton className="h-9 w-16 bg-muted" /> : warningCount}
          </div>
          <p className="text-xs text-muted-foreground">
            Lead-time threshold exceeded
          </p>
        </div>

        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <span className="text-sm font-bold font-display text-foreground">Avg Stock Runway</span>
            <TrendingUp className="h-4 w-4 text-primary" />
          </div>
          <div className="text-3xl font-bold tracking-tight text-foreground">
            {risksLoading ? (
              <Skeleton className="h-9 w-20 bg-muted" />
            ) : avgDaysToStockout !== null ? (
              `${avgDaysToStockout} days`
            ) : (
              '14+ days'
            )}
          </div>
          <p className="text-xs text-muted-foreground">Average days until stock depletion</p>
        </div>

        <div className="space-y-2">
          <div className="flex items-center justify-between">
            <span className="text-sm font-bold font-display text-foreground">AI Inference Engine</span>
            <Cpu className="h-4 w-4 text-primary" />
          </div>
          <div className="text-3xl font-bold tracking-tight text-foreground flex items-center gap-2">
            <span>Online</span>
            <span className="h-2.5 w-2.5 rounded-full bg-emerald-500 animate-ping" />
          </div>
          <p className="text-xs text-muted-foreground">Hybrid ML models active (500ms SLA)</p>
        </div>
      </div>

      {/* 2. Interactive Live SKU Demand Forecaster */}
      <div className="space-y-4">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2 pb-2 border-b border-border/60">
          <div>
            <h3 className="font-bold font-display text-lg text-foreground flex items-center gap-2">
              <Sparkles className="h-5 w-5 text-primary" />
              Live SKU Demand Forecasting & Velocity
            </h3>
            <p className="text-xs text-muted-foreground">
              Predict projected 7-day unit sales based on 30-day lag features, rolling statistics, and margin signals.
            </p>
          </div>

          <div className="flex items-center gap-2">
            <select
              value={selectedProductId}
              onChange={(e) => handleProductSelect(e.target.value)}
              className="h-9 rounded-xl border border-border bg-background px-3 py-1 text-xs font-medium shadow-sm transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring text-foreground max-w-[240px]"
            >
              {products.map((p) => (
                <option key={p.id} value={p.id}>
                  {p.name}
                </option>
              ))}
            </select>

            <Button
              variant="outline"
              size="sm"
              onClick={() => fetchForecast(selectedProductId, shopId || undefined)}
              disabled={forecastLoading || !selectedProductId}
              className="h-9 rounded-xl text-xs gap-1.5 border-border text-foreground hover:text-primary hover:bg-primary/5"
            >
              <RefreshCw className={`h-3.5 w-3.5 ${forecastLoading ? 'animate-spin' : ''}`} />
              Run Inference
            </Button>
          </div>
        </div>

        {forecastLoading ? (
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6 p-6 border border-border/60 rounded-2xl bg-muted/10">
            <Skeleton className="h-28 bg-muted rounded-xl" />
            <Skeleton className="h-28 bg-muted rounded-xl" />
            <Skeleton className="h-28 bg-muted rounded-xl" />
          </div>
        ) : forecastError ? (
          <div className="p-4 rounded-xl bg-destructive/10 border border-destructive/20 text-xs text-destructive">
            {forecastError}
          </div>
        ) : forecast ? (
          <div className="grid grid-cols-1 lg:grid-cols-12 gap-6 p-6 border border-border/60 rounded-2xl bg-background shadow-xs">
            {/* Forecast Metric Summary */}
            <div className="lg:col-span-7 space-y-5">
              <div className="flex items-center justify-between">
                <div>
                  <h4 className="text-xl font-bold font-display text-foreground">{forecast.product_name}</h4>
                  <p className="text-xs text-muted-foreground font-mono mt-0.5">Product ID: {forecast.product_id}</p>
                </div>
                {getConfidenceBadge(forecast.confidence_tier)}
              </div>

              <div className="grid grid-cols-3 gap-4 pt-2 border-t border-border/60">
                <div className="space-y-1">
                  <span className="text-[11px] font-bold text-muted-foreground uppercase tracking-wider">7d Forecast</span>
                  <div className="text-2xl font-bold font-display text-primary">
                    {forecast.predicted_units_sold_7d.toFixed(1)}{' '}
                    <span className="text-xs font-normal text-muted-foreground font-sans">units</span>
                  </div>
                  <p className="text-[10px] text-muted-foreground">Expected weekly volume</p>
                </div>

                <div className="space-y-1">
                  <span className="text-[11px] font-bold text-muted-foreground uppercase tracking-wider">7d Velocity</span>
                  <div className="text-2xl font-bold font-display text-foreground">
                    {forecast.historical_velocity_7d}{' '}
                    <span className="text-xs font-normal text-muted-foreground font-sans">units</span>
                  </div>
                  <p className="text-[10px] text-muted-foreground">Historical trailing sold</p>
                </div>

                <div className="space-y-1">
                  <span className="text-[11px] font-bold text-muted-foreground uppercase tracking-wider">Current Stock</span>
                  <div className={`text-2xl font-bold font-display ${forecast.current_stock < forecast.predicted_units_sold_7d ? 'text-destructive' : 'text-foreground'}`}>
                    {forecast.current_stock}{' '}
                    <span className="text-xs font-normal text-muted-foreground font-sans">available</span>
                  </div>
                  <p className="text-[10px] text-muted-foreground">
                    {forecast.current_stock >= forecast.predicted_units_sold_7d ? (
                      <span className="text-emerald-600 font-medium">✓ Safe runway</span>
                    ) : (
                      <span className="text-destructive font-bold">⚠ Deficit projected</span>
                    )}
                  </p>
                </div>
              </div>

              <div className="p-3.5 rounded-xl bg-muted/40 border border-border/60 text-xs flex items-center justify-between">
                <span className="text-muted-foreground">Inventory Buffer Status:</span>
                <span className="font-semibold text-foreground">
                  {forecast.current_stock >= forecast.predicted_units_sold_7d * 1.5 ? (
                    <span className="text-emerald-600 flex items-center gap-1.5"><ShieldCheck className="h-4 w-4" /> Optimal Stock Coverage</span>
                  ) : forecast.current_stock >= forecast.predicted_units_sold_7d ? (
                    <span className="text-amber-600 flex items-center gap-1.5"><Clock className="h-4 w-4" /> Reorder Suggested Soon</span>
                  ) : (
                    <span className="text-destructive flex items-center gap-1.5 font-bold"><AlertTriangle className="h-4 w-4" /> High Risk of Depletion</span>
                  )}
                </span>
              </div>
            </div>

            {/* Recharts Bar Visualizer */}
            <div className="lg:col-span-5 flex flex-col justify-between pt-4 lg:pt-0 lg:border-l lg:border-border/60 lg:pl-6">
              <span className="text-xs font-bold font-display text-foreground mb-2">Demand vs Velocity Comparison</span>
              <div className="h-[180px] w-full">
                <ResponsiveContainer width="100%" height="100%">
                  <BarChart data={forecastChartData} margin={{ top: 10, right: 10, left: -20, bottom: 0 }}>
                    <XAxis dataKey="name" tick={{ fontSize: 10, fill: 'var(--muted-foreground)' }} axisLine={false} tickLine={false} />
                    <YAxis tick={{ fontSize: 10, fill: 'var(--muted-foreground)' }} axisLine={false} tickLine={false} />
                    <Tooltip content={<ForecastTooltip />} />
                    <Bar dataKey="units" radius={[6, 6, 0, 0]}>
                      {forecastChartData.map((entry, index) => (
                        <Cell key={`cell-${index}`} fill={entry.color} />
                      ))}
                    </Bar>
                  </BarChart>
                </ResponsiveContainer>
              </div>
            </div>
          </div>
        ) : null}
      </div>

      {/* 3. AI Stockout Risk Scanner Table Feed */}
      <div className="space-y-4">
        <div className="flex flex-col md:flex-row md:items-center justify-between gap-4 pb-2 border-b border-border/60">
          <div>
            <h3 className="font-bold font-display text-lg text-foreground flex items-center gap-2">
              <Zap className="h-5 w-5 text-amber-500" />
              Stockout Risk Scanner & Depletion Rates
            </h3>
            <p className="text-xs text-muted-foreground">
              Evaluated by gradient boosted trees ranking inventory depletion velocity against supplier replenishment lead times.
            </p>
          </div>

          <div className="flex flex-wrap items-center gap-3">
            {/* Risk Filter Pills */}
            <div className="flex items-center gap-1 bg-muted p-1 rounded-xl text-xs">
              {(['ALL', 'CRITICAL', 'WARNING', 'NORMAL'] as const).map((lvl) => (
                <button
                  key={lvl}
                  onClick={() => setRiskFilter(lvl)}
                  className={`px-2.5 py-1 rounded-lg transition-all font-medium font-sans uppercase ${
                    riskFilter === lvl
                      ? 'bg-primary text-primary-foreground shadow-sm'
                      : 'text-muted-foreground hover:text-foreground'
                  }`}
                >
                  {lvl} {lvl === 'CRITICAL' && criticalCount > 0 ? `(${criticalCount})` : ''}
                </button>
              ))}
            </div>

            {/* Search Input */}
            <div className="w-full sm:w-60">
              <SearchInput
                value={searchQuery}
                onChange={setSearchQuery}
                placeholder="Search product or shop..."
                className="relative w-full text-foreground text-xs"
              />
            </div>

            <Button
              variant="outline"
              size="sm"
              onClick={() => refreshRisks()}
              disabled={risksLoading}
              className="h-9 rounded-xl text-xs gap-1.5 border-border text-foreground hover:text-primary hover:bg-primary/5"
            >
              <RefreshCw className={`h-3.5 w-3.5 ${risksLoading ? 'animate-spin' : ''}`} />
              Refresh
            </Button>
          </div>
        </div>

        {risksLoading ? (
          <div className="space-y-3">
            {Array.from({ length: 4 }).map((_, i) => (
              <Skeleton key={i} className="h-16 w-full rounded-xl bg-muted" />
            ))}
          </div>
        ) : risksError ? (
          <div className="p-4 rounded-xl bg-destructive/10 border border-destructive/20 text-xs text-destructive">
            {risksError}
          </div>
        ) : filteredRisks.length === 0 ? (
          <EmptyState
            icon={<Package className="h-8 w-8 text-slate-400 mb-2 mx-auto" />}
            title="No stockout risks found"
            description="No inventory line items match your current filter criteria."
            className="py-12 border border-dashed border-border/80 rounded-2xl bg-zinc-50/10"
          />
        ) : (
          <div className="border border-border/60 rounded-2xl overflow-hidden bg-background shadow-xs">
            <Table>
              <TableHeader>
                <TableRow className="hover:bg-transparent">
                  <TableHead className="w-[280px]">Product & Shop</TableHead>
                  <TableHead className="text-right">Available Stock</TableHead>
                  <TableHead className="text-right">Daily Burn Rate</TableHead>
                  <TableHead className="text-right">Days to Stockout</TableHead>
                  <TableHead className="text-right">Probability</TableHead>
                  <TableHead className="text-right">Urgency Ratio</TableHead>
                  <TableHead className="text-center">Risk Level</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {filteredRisks.map((item: StockoutRiskItem) => (
                  <TableRow key={`${item.product_id}-${item.shop_id}`} className="hover:bg-muted/30">
                    <TableCell>
                      <div className="min-w-0">
                        <div className="font-semibold text-xs text-foreground truncate">{item.product_name}</div>
                        <div className="text-[10px] text-muted-foreground flex items-center gap-1 mt-0.5">
                          <Store className="h-3 w-3" />
                          <span>{item.shop_name}</span>
                        </div>
                      </div>
                    </TableCell>

                    <TableCell className="text-right font-mono text-xs">
                      <span className={item.available_stock <= 2 ? 'text-destructive font-bold' : 'text-foreground font-medium'}>
                        {item.available_stock}
                      </span>
                      <span className="text-muted-foreground text-[10px] block">
                        (Total: {item.stock} / Res: {item.reserved_stock})
                      </span>
                    </TableCell>

                    <TableCell className="text-right font-mono text-xs text-foreground">
                      {item.stock_burn_rate_7d.toFixed(2)}{' '}
                      <span className="text-[10px] text-muted-foreground">/day</span>
                    </TableCell>

                    <TableCell className="text-right font-mono text-xs">
                      <span className={item.estimated_days_to_stockout <= 2 ? 'text-destructive font-bold' : item.estimated_days_to_stockout <= 7 ? 'text-amber-600 font-semibold' : 'text-foreground'}>
                        {item.estimated_days_to_stockout.toFixed(1)} days
                      </span>
                    </TableCell>

                    <TableCell className="text-right font-mono text-xs">
                      <span className={item.stockout_probability >= 0.8 ? 'text-destructive font-bold' : item.stockout_probability >= 0.5 ? 'text-amber-600' : 'text-muted-foreground'}>
                        {(item.stockout_probability * 100).toFixed(1)}%
                      </span>
                    </TableCell>

                    <TableCell className="text-right font-mono text-xs">
                      <Badge variant="outline" className="text-[10px] font-mono bg-muted/40">
                        {item.reorder_urgency_ratio.toFixed(2)}x
                      </Badge>
                    </TableCell>

                    <TableCell className="text-center">
                      {getRiskBadge(item.risk_level)}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
        )}
      </div>
    </div>
  );
}

// Custom Tooltip for Forecast Bar Chart
const ForecastTooltip = ({ active, payload }: any) => {
  if (active && payload && payload.length) {
    const data = payload[0].payload;
    return (
      <div className="bg-popover text-popover-foreground border border-border rounded-xl p-3 shadow-md text-xs font-sans">
        <p className="font-semibold font-display mb-1">{data.name}</p>
        <p className="text-muted-foreground">
          Quantity: <span className="font-semibold text-primary">{data.units} units</span>
        </p>
      </div>
    );
  }
  return null;
};
