import React from 'react';
import { DollarSign, Truck, Sparkles, ShieldCheck, ShieldAlert } from 'lucide-react';
import type {
  OrderMetricsResponse,
  ShipmentMetricsResponse,
  StockoutRiskItem,
} from '../../../models/Analytics';
import type { DashboardWafSummary } from '../../../viewmodels/useDashboardViewModel';
import { Skeleton } from '../../../components/ui/skeleton';

interface DashboardKpiCardsProps {
  orderData: OrderMetricsResponse | null;
  shipmentData: ShipmentMetricsResponse | null;
  stockoutRisks: StockoutRiskItem[];
  wafSummary: DashboardWafSummary;
  loading: boolean;
}

export const DashboardKpiCards: React.FC<DashboardKpiCardsProps> = ({
  orderData,
  shipmentData,
  stockoutRisks,
  wafSummary,
  loading,
}) => {
  const formatCurrency = (val: number) => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      maximumFractionDigits: 0,
    }).format(val);
  };

  const gmv = orderData?.summary?.total_gmv || orderData?.summary?.total_revenue || 0;
  const orderCount = orderData?.summary?.total_orders || 0;
  const aov = orderData?.summary?.aov || (orderCount > 0 ? gmv / orderCount : 0);

  const deliveryRate = shipmentData?.summary?.delivery_rate ?? 0;
  const totalShipments = shipmentData?.summary?.total ?? 0;
  const avgHours = shipmentData?.summary?.avg_fulfillment_sec
    ? (shipmentData.summary.avg_fulfillment_sec / 3600).toFixed(1)
    : '0';

  const criticalStockouts = stockoutRisks.filter((r) => r.risk_level === 'CRITICAL').length;
  const warningStockouts = stockoutRisks.filter((r) => r.risk_level === 'WARNING').length;

  if (loading) {
    return (
      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4 border-b border-border/60 pb-8">
        {Array.from({ length: 4 }).map((_, i) => (
          <div key={i} className="space-y-2">
            <Skeleton className="h-4 w-28 bg-muted" />
            <Skeleton className="h-8 w-36 bg-muted" />
            <Skeleton className="h-3 w-24 bg-muted" />
          </div>
        ))}
      </div>
    );
  }

  return (
    <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4 border-b border-border/60 pb-8">
      {/* 1. E-Commerce Revenue */}
      <div className="space-y-2">
        <div className="flex items-center justify-between">
          <span className="text-sm font-bold font-display text-foreground">Gross Revenue (GMV)</span>
          <div className="w-6 h-6 rounded-md bg-primary/10 flex items-center justify-center text-primary">
            <DollarSign className="h-3.5 w-3.5" />
          </div>
        </div>
        <div className="text-2xl lg:text-3xl font-bold font-display tracking-tight text-foreground">
          {formatCurrency(gmv)}
        </div>
        <p className="text-xs text-muted-foreground font-sans">
          <span className="font-semibold text-foreground">{orderCount}</span> orders • AOV:{' '}
          <span className="font-semibold text-foreground">{formatCurrency(aov)}</span>
        </p>
      </div>

      {/* 2. Logistics & Fulfillment SLA */}
      <div className="space-y-2">
        <div className="flex items-center justify-between">
          <span className="text-sm font-bold font-display text-foreground">Fulfillment SLA</span>
          <div className="w-6 h-6 rounded-md bg-sky-500/10 flex items-center justify-center text-sky-600 dark:text-sky-400">
            <Truck className="h-3.5 w-3.5" />
          </div>
        </div>
        <div className="text-2xl lg:text-3xl font-bold font-display tracking-tight text-foreground">
          {(deliveryRate * 100).toFixed(1)}%
        </div>
        <p className="text-xs text-muted-foreground font-sans">
          <span className="font-semibold text-foreground">{totalShipments}</span> parcels • avg{' '}
          <span className="font-semibold text-foreground">{avgHours}h</span> turnaround
        </p>
      </div>

      {/* 3. AI Stockout Alerts */}
      <div className="space-y-2">
        <div className="flex items-center justify-between">
          <span className="text-sm font-bold font-display text-foreground">AI Stockout Risks</span>
          <div className="w-6 h-6 rounded-md bg-amber-500/10 flex items-center justify-center text-amber-600 dark:text-amber-400">
            <Sparkles className="h-3.5 w-3.5" />
          </div>
        </div>
        <div className="text-2xl lg:text-3xl font-bold font-display tracking-tight text-foreground flex items-baseline gap-2">
          <span className={criticalStockouts > 0 ? 'text-rose-600 dark:text-rose-400' : 'text-foreground'}>
            {criticalStockouts} Critical
          </span>
          {warningStockouts > 0 && (
            <span className="text-xs font-sans text-amber-600 dark:text-amber-400 font-medium">
              +{warningStockouts} Warning
            </span>
          )}
        </div>
        <p className="text-xs text-muted-foreground font-sans">
          {criticalStockouts === 0 && warningStockouts === 0
            ? 'Inventory health optimal'
            : 'Immediate reorder recommended'}
        </p>
      </div>

      {/* 4. Cyber Threat & WAF */}
      <div className="space-y-2">
        <div className="flex items-center justify-between">
          <span className="text-sm font-bold font-display text-foreground">WAF Security Defense</span>
          <div
            className={`w-6 h-6 rounded-md flex items-center justify-center ${
              wafSummary.blocked > 0
                ? 'bg-rose-500/10 text-rose-600 dark:text-rose-400'
                : 'bg-primary/10 text-primary'
            }`}
          >
            {wafSummary.blocked > 0 ? (
              <ShieldAlert className="h-3.5 w-3.5" />
            ) : (
              <ShieldCheck className="h-3.5 w-3.5" />
            )}
          </div>
        </div>
        <div
          className={`text-2xl lg:text-3xl font-bold font-display tracking-tight ${
            wafSummary.blocked > 0 ? 'text-rose-600 dark:text-rose-400' : 'text-primary'
          }`}
        >
          {wafSummary.blocked} Threats Blocked
        </div>
        <p className="text-xs text-muted-foreground font-sans">
          Threat Level: <span className="font-semibold text-foreground">{wafSummary.threatLevel}</span> •{' '}
          <span className="font-semibold text-foreground">{wafSummary.activeRules}</span> active policies
        </p>
      </div>
    </div>
  );
};

export default DashboardKpiCards;
