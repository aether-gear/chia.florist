import { useState, useEffect, useCallback, useMemo } from 'react';
import { fetchApi } from '../lib/api';
import type {
  OrderMetricsResponse,
  PaymentMetricsResponse,
  ShipmentMetricsResponse,
  InventoryMetricsResponse,
  ProductMetricsResponse,
  StockoutRiskItem,
  StockoutRisksResponse,
  DemandForecast,
} from '../models/Analytics';
import type { Order, OrdersResponse } from '../models/Order';
import type { ProductStat, ProductStatsResponse } from '../models/Product';
import type { Shop } from '../models/Shop';

export type DashboardViewTab = 'overview' | 'ecommerce' | 'ai' | 'cyber';
export type DashboardDatePreset = '7d' | '30d' | '90d';

export interface DashboardWafSummary {
  total: number;
  blocked: number;
  allowed: number;
  threatLevel: 'Low' | 'Medium' | 'High' | 'Critical';
  activeRules: number;
  bannedIps: number;
  whitelistedIps: number;
}

export interface SecurityEventLog {
  id: string;
  timestamp: string;
  ip: string;
  method: string;
  url: string;
  status: 'Blocked' | 'Allowed';
  ruleId: string;
  reason: string;
  payload: string;
  userAgent: string;
}

export interface DynamicAIInsight {
  id: string;
  title: string;
  description: string;
  category: 'ecommerce' | 'ai' | 'cyber';
  severity: 'info' | 'warning' | 'critical' | 'success';
}

export function useDashboardViewModel() {
  // Navigation and Filter States
  const [activeTab, setActiveTab] = useState<DashboardViewTab>('overview');
  const [preset, setPreset] = useState<DashboardDatePreset>('30d');
  const [shopId, setShopId] = useState<string>('');

  // Domain Data States
  const [orderMetrics, setOrderMetrics] = useState<OrderMetricsResponse | null>(null);
  const [paymentMetrics, setPaymentMetrics] = useState<PaymentMetricsResponse | null>(null);
  const [shipmentMetrics, setShipmentMetrics] = useState<ShipmentMetricsResponse | null>(null);
  const [inventoryMetrics, setInventoryMetrics] = useState<InventoryMetricsResponse | null>(null);
  const [productMetrics, setProductMetrics] = useState<ProductMetricsResponse | null>(null);
  const [productStats, setProductStats] = useState<ProductStat[]>([]);
  const [recentOrders, setRecentOrders] = useState<Order[]>([]);
  const [stockoutRisks, setStockoutRisks] = useState<StockoutRiskItem[]>([]);
  const [shops, setShops] = useState<Shop[]>([]);

  // Cyber / WAF States
  const [securityLogs, setSecurityLogs] = useState<SecurityEventLog[]>([]);
  const [wafRulesCount, setWafRulesCount] = useState<number>(0);
  const [bannedIpsCount, setBannedIpsCount] = useState<number>(0);
  const [whitelistedIpsCount, setWhitelistedIpsCount] = useState<number>(0);

  // Demand Forecast Interactive State
  const [selectedForecastProductId, setSelectedForecastProductId] = useState<string>('');
  const [forecastData, setForecastData] = useState<DemandForecast | null>(null);
  const [forecastLoading, setForecastLoading] = useState<boolean>(false);

  // Inspection Drawer States
  const [selectedOrder, setSelectedOrder] = useState<Order | null>(null);
  const [isOrderDetailOpen, setIsOrderDetailOpen] = useState<boolean>(false);
  const [selectedSecurityLog, setSelectedSecurityLog] = useState<SecurityEventLog | null>(null);
  const [isSecurityDetailOpen, setIsSecurityDetailOpen] = useState<boolean>(false);

  // Global Loading & Error
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  // Compute ISO Dates
  const getDates = useCallback(() => {
    const now = new Date();
    const toStr = now.toISOString().split('T')[0];
    const past = new Date(now);
    if (preset === '7d') past.setDate(past.getDate() - 7);
    else if (preset === '30d') past.setDate(past.getDate() - 30);
    else if (preset === '90d') past.setDate(past.getDate() - 90);
    const fromStr = past.toISOString().split('T')[0];
    return { from: fromStr, to: toStr };
  }, [preset]);

  // Fetch Shops
  const fetchShops = useCallback(async () => {
    try {
      const res = await fetchApi('/shops?limit=100');
      if (res?.shops) setShops(res.shops);
    } catch {
      // Gracefully ignore
    }
  }, []);

  // Main Dashboard Data Loader
  const loadDashboardData = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const { from, to } = getDates();

      const buildQuery = (params: Record<string, string | number | undefined>) => {
        const p = new URLSearchParams();
        if (from) p.append('from', from);
        if (to) p.append('to', to);
        Object.entries(params).forEach(([k, v]) => {
          if (v !== undefined && v !== '') p.append(k, String(v));
        });
        const str = p.toString();
        return str ? `?${str}` : '';
      };

      const [
        ordersRes,
        paymentsRes,
        shipmentsRes,
        inventoryRes,
        productsRes,
        recentOrdersRes,
        productStatsRes,
        stockoutRes,
        statsRes,
        rulesRes,
        ipRes,
      ] = await Promise.allSettled([
        fetchApi(`/analytics/orders${buildQuery({ granularity: 'daily', shop_id: shopId, top_n: 8 })}`),
        fetchApi(`/analytics/payments${buildQuery({})}`),
        fetchApi(`/analytics/shipments${buildQuery({ top_n: 8 })}`),
        fetchApi(`/analytics/inventory${buildQuery({ shop_id: shopId })}`),
        fetchApi(`/analytics/products${buildQuery({ top_n: 8 })}`),
        fetchApi(`/orders?limit=10&sort=latest:desc${shopId ? `&shop=${shopId}` : ''}`),
        fetchApi('/products/stats?page=1&limit=50'),
        fetchApi(`/analytics/inventory/stockout-risks${shopId ? `?shop_id=${shopId}` : ''}`),
        fetchApi('/api/stats?limit=300'),
        fetchApi('/api/rules'),
        fetchApi('/api/ip'),
      ]);

      if (ordersRes.status === 'fulfilled') setOrderMetrics(ordersRes.value);
      if (paymentsRes.status === 'fulfilled') setPaymentMetrics(paymentsRes.value);
      if (shipmentsRes.status === 'fulfilled') setShipmentMetrics(shipmentsRes.value);
      if (inventoryRes.status === 'fulfilled') setInventoryMetrics(inventoryRes.value);
      if (productsRes.status === 'fulfilled') setProductMetrics(productsRes.value);

      if (recentOrdersRes.status === 'fulfilled') {
        const orderData = recentOrdersRes.value as OrdersResponse;
        setRecentOrders(orderData?.orders || []);
      }

      if (productStatsRes.status === 'fulfilled') {
        const pStatsData = productStatsRes.value as ProductStatsResponse;
        const statsList = pStatsData?.stats || [];
        setProductStats(statsList);
        if (statsList.length > 0 && !selectedForecastProductId) {
          setSelectedForecastProductId(statsList[0].id);
        }
      }

      if (stockoutRes.status === 'fulfilled') {
        const stockData = stockoutRes.value as StockoutRisksResponse;
        setStockoutRisks(stockData?.risks || []);
      }

      if (rulesRes.status === 'fulfilled') {
        setWafRulesCount(rulesRes.value?.rules?.length || 0);
      }

      if (ipRes.status === 'fulfilled') {
        const entries = ipRes.value?.entries || [];
        const banned = entries.filter((e: any) => String(e.status).toLowerCase().includes('ban')).length;
        const white = entries.filter((e: any) => String(e.status).toLowerCase().includes('white')).length;
        setBannedIpsCount(banned);
        setWhitelistedIpsCount(white);
      }

      if (statsRes.status === 'fulfilled') {
        const rawLogs = statsRes.value?.audit_logs || [];
        const mapped: SecurityEventLog[] = rawLogs.map((log: any) => {
          const isBlocked = log.outcome === 'blocked' || log.action === 'request_blocked';
          return {
            id: log.id,
            timestamp: log.created_at,
            ip: log.client_ip || '127.0.0.1',
            method: log.metadata?.method || 'GET',
            url: log.metadata?.path || '/',
            status: isBlocked ? 'Blocked' : 'Allowed',
            ruleId: log.metadata?.rule_id || '-',
            reason: log.metadata?.reason || '-',
            payload: log.metadata?.payload || '-',
            userAgent: log.metadata?.user_agent || '-',
          };
        });
        setSecurityLogs(mapped);
      }
    } catch (err: any) {
      setError(err.message || 'Failed to refresh dashboard');
    } finally {
      setLoading(false);
    }
  }, [getDates, shopId, selectedForecastProductId]);

  // Initial mount
  useEffect(() => {
    fetchShops();
  }, [fetchShops]);

  useEffect(() => {
    loadDashboardData();
  }, [loadDashboardData]);

  // Load interactive AI Demand Forecast when selected product changes
  useEffect(() => {
    if (!selectedForecastProductId) return;
    let isCancelled = false;

    const fetchForecast = async () => {
      try {
        setForecastLoading(true);
        const params = new URLSearchParams();
        params.append('product_id', selectedForecastProductId);
        if (shopId) params.append('shop_id', shopId);

        const res: DemandForecast = await fetchApi(`/analytics/forecasts/demand?${params.toString()}`);
        if (!isCancelled) {
          setForecastData(res);
        }
      } catch {
        if (!isCancelled) setForecastData(null);
      } finally {
        if (!isCancelled) setForecastLoading(false);
      }
    };

    fetchForecast();
    return () => {
      isCancelled = true;
    };
  }, [selectedForecastProductId, shopId]);

  // Computed WAF Summary
  const wafSummary: DashboardWafSummary = useMemo(() => {
    const total = securityLogs.length;
    const blocked = securityLogs.filter((l) => l.status === 'Blocked').length;
    const allowed = total - blocked;

    let threatLevel: 'Low' | 'Medium' | 'High' | 'Critical' = 'Low';
    if (blocked > 50) threatLevel = 'Critical';
    else if (blocked > 20) threatLevel = 'High';
    else if (blocked > 5) threatLevel = 'Medium';

    return {
      total,
      blocked,
      allowed,
      threatLevel,
      activeRules: wafRulesCount,
      bannedIps: bannedIpsCount,
      whitelistedIps: whitelistedIpsCount,
    };
  }, [securityLogs, wafRulesCount, bannedIpsCount, whitelistedIpsCount]);

  // Dynamic AI Advisories generated from real metrics
  const aiInsights: DynamicAIInsight[] = useMemo(() => {
    const list: DynamicAIInsight[] = [];

    // 1. Stockout Risk Insight
    const criticalStockouts = stockoutRisks.filter((r) => r.risk_level === 'CRITICAL');
    if (criticalStockouts.length > 0) {
      const topCritical = criticalStockouts[0];
      list.push({
        id: 'stockout-crit',
        title: 'Critical Restock Alert',
        description: `"${topCritical.product_name}" has only ${topCritical.available_stock} units left at ${topCritical.shop_name} (est. ${topCritical.estimated_days_to_stockout} days until stockout). Supplier lead time is ${topCritical.supplier_lead_time_days} days.`,
        category: 'ai',
        severity: 'critical',
      });
    }

    // 2. Cyber Threat Anomaly
    if (wafSummary.blocked > 0) {
      list.push({
        id: 'cyber-threat',
        title: 'WAF Security Telemetry',
        description: `WAF intercepted and mitigated ${wafSummary.blocked} malicious payloads across ${wafSummary.total} requests. Threat assessment is currently rated at ${wafSummary.threatLevel}.`,
        category: 'cyber',
        severity: wafSummary.threatLevel === 'Critical' || wafSummary.threatLevel === 'High' ? 'warning' : 'info',
      });
    }

    // 3. Top Bouquet Sales Velocity
    if (productStats.length > 0) {
      const topSeller = [...productStats].sort((a, b) => b.sales_velocity_7d - a.sales_velocity_7d)[0];
      if (topSeller && topSeller.sales_velocity_7d > 0) {
        list.push({
          id: 'top-product',
          title: 'Top Sales Velocity Leader',
          description: `"${topSeller.name}" generated ${topSeller.sales_velocity_7d} sales over the last 7 days, contributing to ${topSeller.revenue_contribution_percentage}% of total catalog revenue.`,
          category: 'ecommerce',
          severity: 'success',
        });
      }
    }

    // 4. Logistics / Delivery Success Rate
    if (shipmentMetrics?.summary) {
      const { delivery_rate, avg_fulfillment_sec } = shipmentMetrics.summary;
      const hours = (avg_fulfillment_sec / 3600).toFixed(1);
      list.push({
        id: 'logistics-rate',
        title: 'Fulfillment & Courier Health',
        description: `Current fleet delivery success rate is ${(delivery_rate * 100).toFixed(1)}% with an average dispatch-to-delivery turnaround of ${hours} hours.`,
        category: 'ecommerce',
        severity: delivery_rate >= 0.9 ? 'success' : 'warning',
      });
    }

    return list;
  }, [stockoutRisks, wafSummary, productStats, shipmentMetrics]);

  return {
    activeTab,
    setActiveTab,
    preset,
    setPreset,
    shopId,
    setShopId,
    shops,
    orderMetrics,
    paymentMetrics,
    shipmentMetrics,
    inventoryMetrics,
    productMetrics,
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
    refresh: loadDashboardData,
  };
}
