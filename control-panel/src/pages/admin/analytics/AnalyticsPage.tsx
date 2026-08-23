import { useSearchParams } from 'react-router-dom';
import { RefreshCw, BarChart3, ShoppingBag, CreditCard, Truck, Package, Sparkles } from 'lucide-react';
import { useAnalyticsViewModel } from '../../../viewmodels/useAnalyticsViewModel';
import AnalyticsOverviewTab from './components/AnalyticsOverviewTab';
import OrderAnalyticsTab from './components/OrderAnalyticsTab';
import PaymentAnalyticsTab from './components/PaymentAnalyticsTab';
import ShipmentAnalyticsTab from './components/ShipmentAnalyticsTab';
import ProductInventoryAnalyticsTab from './components/ProductInventoryAnalyticsTab';
import AIForecastsAnalyticsTab from './components/AIForecastsAnalyticsTab';
import { Button } from '../../../components/ui/button';

export default function AnalyticsPage() {
  const [searchParams, setSearchParams] = useSearchParams();
  const activeTab = searchParams.get('tab') || 'overview';

  const {
    preset,
    setPreset,
    fromDate,
    setFromDate,
    toDate,
    setToDate,
    shopId,
    setShopId,
    shops,
    orderMetrics,
    paymentMetrics,
    shipmentMetrics,
    inventoryMetrics,
    productMetrics,
    loading,
    error,
    refresh,
  } = useAnalyticsViewModel();

  const handleTabChange = (tab: string) => {
    setSearchParams({ tab });
  };

  const tabs = [
    { id: 'overview', label: 'Executive Overview', icon: BarChart3 },
    { id: 'ai-intelligence', label: 'AI Intelligence & Forecasting', icon: Sparkles },
    { id: 'orders', label: 'Orders & Sales', icon: ShoppingBag },
    { id: 'payments', label: 'Payments & Revenue', icon: CreditCard },
    { id: 'shipments', label: 'Fulfillment & Logistics', icon: Truck },
    { id: 'products', label: 'Products & Inventory', icon: Package },
  ];

  return (
    <div className="flex-col md:flex">
      <div className="flex-1 space-y-12 p-6 sm:p-8 lg:p-12 animate-in fade-in duration-300">

        {/* Page Header */}
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <h2 className="text-3xl font-bold font-display tracking-tight text-foreground">
              Analytics
            </h2>
            <p className="text-muted-foreground text-sm">
              Real-time business performance, order volume, logistics, and inventory metrics
            </p>
          </div>
        </div>

        {/* Analytics Workspace Section */}
        <div className="space-y-6">
          <div className="pb-4 border-b border-border/60">
            <h3 className="text-xl font-bold font-display tracking-tight text-foreground">
              Analytics Workspace
            </h3>
            <p className="text-muted-foreground text-sm">
              Filter performance data by domain tabs, shop branch, and time period.
            </p>
          </div>

          {/* Global Filters & Controls (Below Selection Pills) */}
          <div className="flex flex-wrap items-center justify-between gap-3 pt-1">
            <div className="flex flex-wrap items-center gap-3">
              {/* Shop Selector */}
              <select
                value={shopId}
                onChange={(e) => setShopId(e.target.value)}
                className="h-9 rounded-xl border border-border bg-background px-3 py-1 text-xs font-medium shadow-sm transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring text-foreground"
              >
                <option value="">All Shop Branches</option>
                {(Array.isArray(shops) ? shops : []).map((s) => (
                  <option key={s.id} value={s.id}>
                    {s.name}
                  </option>
                ))}
              </select>

              {/* Date Preset Pills */}
              <div className="flex items-center gap-1 bg-muted p-1 rounded-xl text-xs">
                {(['7d', '30d', '90d', 'custom'] as const).map((p) => (
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

              {/* Custom Date Inputs */}
              {preset === 'custom' && (
                <div className="flex items-center gap-2">
                  <input
                    type="date"
                    value={fromDate}
                    onChange={(e) => setFromDate(e.target.value)}
                    className="h-9 rounded-xl border border-border bg-background px-2.5 py-1 text-xs text-foreground"
                  />
                  <span className="text-xs text-muted-foreground">to</span>
                  <input
                    type="date"
                    value={toDate}
                    onChange={(e) => setToDate(e.target.value)}
                    className="h-9 rounded-xl border border-border bg-background px-2.5 py-1 text-xs text-foreground"
                  />
                </div>
              )}
            </div>

            {/* Refresh Button */}
            <Button
              variant="outline"
              onClick={() => refresh()}
              disabled={loading}
              className="flex items-center gap-1.5 border-border text-foreground hover:text-primary hover:bg-primary/5 rounded-xl transition-colors h-9 text-xs"
            >
              <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
              Refresh
            </Button>
          </div>

          {/* Selection Pills (Borderless) */}
          <div className="flex gap-1 overflow-x-auto pb-1 w-full scrollbar-none">
            {tabs.map((tab) => {
              const Icon = tab.icon;
              const isActive = activeTab === tab.id;
              return (
                <button
                  key={tab.id}
                  onClick={() => handleTabChange(tab.id)}
                  className={`flex items-center gap-2 px-3 py-1.5 text-xs font-semibold whitespace-nowrap rounded-lg transition-all ${
                    isActive
                      ? 'bg-primary/10 text-primary'
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

          {/* Active Tab View Content */}
          <div className="pt-2">
            {activeTab === 'overview' && (
              <AnalyticsOverviewTab
                orderData={orderMetrics}
                paymentData={paymentMetrics}
                shipmentData={shipmentMetrics}
                inventoryData={inventoryMetrics}
                loading={loading}
              />
            )}
            {activeTab === 'ai-intelligence' && (
              <AIForecastsAnalyticsTab shopId={shopId} />
            )}
            {activeTab === 'orders' && (
              <OrderAnalyticsTab data={orderMetrics} loading={loading} />
            )}
            {activeTab === 'payments' && (
              <PaymentAnalyticsTab data={paymentMetrics} loading={loading} />
            )}
            {activeTab === 'shipments' && (
              <ShipmentAnalyticsTab data={shipmentMetrics} loading={loading} />
            )}
            {activeTab === 'products' && (
              <ProductInventoryAnalyticsTab
                inventoryData={inventoryMetrics}
                productData={productMetrics}
                loading={loading}
              />
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
