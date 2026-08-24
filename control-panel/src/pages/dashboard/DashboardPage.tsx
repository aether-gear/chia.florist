import {
  LayoutDashboard,
  ShoppingBag,
  Sparkles,
  Shield,
  RefreshCw,
} from 'lucide-react';
import { Button } from '../../components/ui/button';
import { useDashboardViewModel } from '../../viewmodels/useDashboardViewModel';
import { useOrderActionsViewModel } from '../../viewmodels/useOrderActionsViewModel';
import DashboardKpiCards from './components/DashboardKpiCards';
import DashboardAIInsights from './components/DashboardAIInsights';
import DashboardEcommerceView from './components/DashboardEcommerceView';
import DashboardAIView from './components/DashboardAIView';
import DashboardCyberView from './components/DashboardCyberView';
import SecurityDetailSheet from './components/SecurityDetailSheet';
import OrderDetailSheet from '../../components/orders/OrderDetailSheet';
import type { Order } from '../../models/Order';
import type { SecurityEventLog } from '../../viewmodels/useDashboardViewModel';

export default function DashboardPage() {
  const {
    activeTab,
    setActiveTab,
    preset,
    setPreset,
    shopId,
    setShopId,
    shops,
    orderMetrics,
    shipmentMetrics,
    productStats,
    recentOrders,
    stockoutRisks,
    securityLogs,
    wafSummary,
    aiInsights,
    selectedForecastProductId,
    setSelectedForecastProductId,
    forecastData,
    forecastLoading,
    selectedOrder,
    setSelectedOrder,
    isOrderDetailOpen,
    setIsOrderDetailOpen,
    selectedSecurityLog,
    setSelectedSecurityLog,
    isSecurityDetailOpen,
    setIsSecurityDetailOpen,
    loading,
    error,
    refresh,
  } = useDashboardViewModel();

  const orderActions = useOrderActionsViewModel();

  const handleInspectOrder = (order: Order) => {
    setSelectedOrder(order);
    setIsOrderDetailOpen(true);
  };

  const handleInspectSecurityLog = (log: SecurityEventLog) => {
    setSelectedSecurityLog(log);
    setIsSecurityDetailOpen(true);
  };

  const tabs = [
    { id: 'overview', label: 'Executive Hub', icon: LayoutDashboard },
    { id: 'ecommerce', label: 'E-Commerce Operations', icon: ShoppingBag },
    { id: 'ai', label: 'AI Intelligence & Forecasts', icon: Sparkles },
    { id: 'cyber', label: 'Cyber Defense & WAF', icon: Shield },
  ] as const;

  return (
    <div className="space-y-10 animate-in fade-in duration-300">
      {/* 1. Header & Global Filters */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-6 border-b border-border/60">
        <div>
          <h2 className="text-3xl font-bold font-display tracking-tight text-foreground">
            Command Dashboard
          </h2>
          <p className="text-muted-foreground text-sm font-sans mt-0.5">
            Real-time unified intelligence spanning botanical e-commerce, predictive AI, and cyber defense.
          </p>
        </div>

        {/* Global Toolbar */}
        <div className="flex flex-wrap items-center gap-2.5">
          {/* Shop Selector */}
          <select
            value={shopId}
            onChange={(e) => setShopId(e.target.value)}
            className="h-9 rounded-xl border border-border bg-background px-3 py-1 text-xs font-medium shadow-sm transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring text-foreground"
          >
            <option value="">All Boutique Branches</option>
            {shops.map((s) => (
              <option key={s.id} value={s.id}>
                {s.name}
              </option>
            ))}
          </select>

          {/* Date Range Pills */}
          <div className="flex items-center gap-1 bg-muted p-1 rounded-xl text-xs">
            {(['7d', '30d', '90d'] as const).map((p) => (
              <button
                key={p}
                onClick={() => setPreset(p)}
                className={`px-2.5 py-1 rounded-lg transition-all font-medium font-sans uppercase ${
                  preset === p
                    ? 'bg-primary text-primary-foreground shadow-sm'
                    : 'text-muted-foreground hover:text-foreground'
                }`}
              >
                {p}
              </button>
            ))}
          </div>

          {/* Refresh Button */}
          <Button
            variant="outline"
            size="sm"
            onClick={() => refresh()}
            disabled={loading}
            className="h-9 text-xs gap-1.5 border-border rounded-xl text-foreground hover:text-primary hover:bg-primary/5"
          >
            <RefreshCw className={`w-3.5 h-3.5 ${loading ? 'animate-spin' : ''}`} />
            Refresh
          </Button>
        </div>
      </div>

      {/* 2. Multi-Pillar KPI Banner */}
      <DashboardKpiCards
        orderData={orderMetrics}
        shipmentData={shipmentMetrics}
        stockoutRisks={stockoutRisks}
        wafSummary={wafSummary}
        loading={loading}
      />

      {/* 3. Domain Navigation Tabs (Borderless Selection Pills) */}
      <div className="flex gap-1 overflow-x-auto pb-1 w-full border-b border-border/40 scrollbar-none">
        {tabs.map((tab) => {
          const Icon = tab.icon;
          const isActive = activeTab === tab.id;
          return (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={`flex items-center gap-2 px-4 py-2 text-xs font-semibold whitespace-nowrap rounded-lg transition-all ${
                isActive
                  ? 'bg-primary/10 text-primary font-bold'
                  : 'text-muted-foreground hover:text-foreground hover:bg-muted/50'
              }`}
            >
              <Icon className="h-4 w-4" />
              <span>{tab.label}</span>
            </button>
          );
        })}
      </div>

      {/* Error Alert */}
      {error && (
        <div className="p-4 rounded-xl bg-destructive/10 border border-destructive/20 text-xs text-destructive">
          {error}
        </div>
      )}

      {/* 4. Active Tab Content */}
      <div className="pt-2">
        {activeTab === 'overview' && (
          <div className="space-y-10">
            {/* Overview Mid Grid: E-Commerce Velocity + Live AI Advisories */}
            <div className="grid gap-10 lg:grid-cols-7 pb-8 border-b border-border/60">
              <div className="lg:col-span-4 space-y-4">
                <DashboardEcommerceView
                  orderData={orderMetrics}
                  productStats={productStats}
                  recentOrders={recentOrders}
                  loading={loading}
                  onInspectOrder={handleInspectOrder}
                />
              </div>

              <div className="lg:col-span-3 space-y-6">
                <DashboardAIInsights insights={aiInsights} loading={loading} />
              </div>
            </div>

            {/* Overview Bottom: Cyber Telemetry Stream Preview */}
            <div className="pt-2">
              <DashboardCyberView
                wafSummary={wafSummary}
                securityLogs={securityLogs}
                loading={loading}
                onInspectSecurityLog={handleInspectSecurityLog}
              />
            </div>
          </div>
        )}

        {activeTab === 'ecommerce' && (
          <DashboardEcommerceView
            orderData={orderMetrics}
            productStats={productStats}
            recentOrders={recentOrders}
            loading={loading}
            onInspectOrder={handleInspectOrder}
          />
        )}

        {activeTab === 'ai' && (
          <DashboardAIView
            stockoutRisks={stockoutRisks}
            productStats={productStats}
            selectedForecastProductId={selectedForecastProductId}
            onSelectForecastProduct={setSelectedForecastProductId}
            forecastData={forecastData}
            forecastLoading={forecastLoading}
            loading={loading}
          />
        )}

        {activeTab === 'cyber' && (
          <DashboardCyberView
            wafSummary={wafSummary}
            securityLogs={securityLogs}
            loading={loading}
            onInspectSecurityLog={handleInspectSecurityLog}
          />
        )}
      </div>

      {/* 5. Detail Inspection Drawers (Right Overlay Sheet) */}
      <OrderDetailSheet
        order={selectedOrder}
        isOpen={isOrderDetailOpen}
        onOpenChange={setIsOrderDetailOpen}
        submitting={orderActions.submitting}
        onStartProcessing={async (orderId) => {
          await orderActions.updateOrderStatus(orderId, 'processing');
          refresh();
        }}
        onDispatchOrder={async (orderId, shipments) => {
          await orderActions.updateOrderStatus(orderId, 'shipped', undefined, undefined, shipments);
          refresh();
        }}
        onUpdateShipmentStatus={async (shipmentId, status) => {
          await orderActions.updateShipmentStatus(shipmentId, status);
          refresh();
        }}
        onUpdateWaybill={async (shipmentId, details) => {
          await orderActions.updateShipmentDetails(shipmentId, details);
          refresh();
        }}
        fetchOrderTracking={orderActions.fetchOrderTracking}
      />

      <SecurityDetailSheet
        log={selectedSecurityLog}
        isOpen={isSecurityDetailOpen}
        onOpenChange={setIsSecurityDetailOpen}
      />
    </div>
  );
}
