import { useState, useEffect, useCallback } from 'react';
import { fetchApi } from '../lib/api';
import type {
  OrderMetricsResponse,
  PaymentMetricsResponse,
  ShipmentMetricsResponse,
  InventoryMetricsResponse,
  ProductMetricsResponse,
} from '../models/Analytics';
import type { Shop } from '../models/Shop';

export type DatePreset = '7d' | '30d' | '90d' | 'custom';

export function useAnalyticsViewModel() {
  const [preset, setPreset] = useState<DatePreset>('30d');
  const [fromDate, setFromDate] = useState<string>('');
  const [toDate, setToDate] = useState<string>('');
  const [granularity, setGranularity] = useState<'daily' | 'weekly' | 'monthly'>('daily');
  const [shopId, setShopId] = useState<string>('');
  const [topN, setTopN] = useState<number>(10);

  // Data states
  const [orderMetrics, setOrderMetrics] = useState<OrderMetricsResponse | null>(null);
  const [paymentMetrics, setPaymentMetrics] = useState<PaymentMetricsResponse | null>(null);
  const [shipmentMetrics, setShipmentMetrics] = useState<ShipmentMetricsResponse | null>(null);
  const [inventoryMetrics, setInventoryMetrics] = useState<InventoryMetricsResponse | null>(null);
  const [productMetrics, setProductMetrics] = useState<ProductMetricsResponse | null>(null);
  const [shops, setShops] = useState<Shop[]>([]);

  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  // Helper to compute date string range based on preset
  const getComputedDates = useCallback(() => {
    if (preset === 'custom') {
      return { from: fromDate, to: toDate };
    }
    const now = new Date();
    const toStr = now.toISOString().split('T')[0];
    const past = new Date(now);
    if (preset === '7d') past.setDate(past.getDate() - 7);
    else if (preset === '30d') past.setDate(past.getDate() - 30);
    else if (preset === '90d') past.setDate(past.getDate() - 90);
    const fromStr = past.toISOString().split('T')[0];
    return { from: fromStr, to: toStr };
  }, [preset, fromDate, toDate]);

  // Fetch list of shops for dropdown filter
  const fetchShopsList = useCallback(async () => {
    try {
      const res = await fetchApi('/shops?limit=100');
      if (res && res.shops) {
        setShops(res.shops);
      }
    } catch {
      // Ignore shop list error gracefully
    }
  }, []);

  // Fetch all analytics data
  const fetchAllAnalytics = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);

      const { from, to } = getComputedDates();

      // Build query string helper
      const buildQueryParams = (extra: Record<string, string | number | undefined>) => {
        const params = new URLSearchParams();
        if (from) params.append('from', from);
        if (to) params.append('to', to);
        Object.entries(extra).forEach(([k, v]) => {
          if (v !== undefined && v !== '') {
            params.append(k, String(v));
          }
        });
        const str = params.toString();
        return str ? `?${str}` : '';
      };

      const [ordersRes, paymentsRes, shipmentsRes, inventoryRes, productsRes] = await Promise.allSettled([
        fetchApi(`/analytics/orders${buildQueryParams({ granularity, shop_id: shopId, top_n: topN })}`),
        fetchApi(`/analytics/payments${buildQueryParams({})}`),
        fetchApi(`/analytics/shipments${buildQueryParams({ top_n: topN })}`),
        fetchApi(`/analytics/inventory${buildQueryParams({ shop_id: shopId })}`),
        fetchApi(`/analytics/products${buildQueryParams({ top_n: topN })}`),
      ]);

      if (ordersRes.status === 'fulfilled') setOrderMetrics(ordersRes.value);
      if (paymentsRes.status === 'fulfilled') setPaymentMetrics(paymentsRes.value);
      if (shipmentsRes.status === 'fulfilled') setShipmentMetrics(shipmentsRes.value);
      if (inventoryRes.status === 'fulfilled') setInventoryMetrics(inventoryRes.value);
      if (productsRes.status === 'fulfilled') setProductMetrics(productsRes.value);

      // Check if all failed
      const allRejected = [ordersRes, paymentsRes, shipmentsRes, inventoryRes, productsRes].every(
        r => r.status === 'rejected'
      );
      if (allRejected) {
        const firstError = (ordersRes as PromiseRejectedResult).reason?.message || 'Failed to load analytics data';
        setError(firstError);
      }
    } catch (err: any) {
      setError(err.message || 'An error occurred while loading analytics');
    } finally {
      setLoading(false);
    }
  }, [getComputedDates, granularity, shopId, topN]);

  useEffect(() => {
    fetchShopsList();
  }, [fetchShopsList]);

  useEffect(() => {
    fetchAllAnalytics();
  }, [fetchAllAnalytics]);

  return {
    preset,
    setPreset,
    fromDate,
    setFromDate,
    toDate,
    setToDate,
    granularity,
    setGranularity,
    shopId,
    setShopId,
    topN,
    setTopN,
    shops,
    orderMetrics,
    paymentMetrics,
    shipmentMetrics,
    inventoryMetrics,
    productMetrics,
    loading,
    error,
    refresh: fetchAllAnalytics,
  };
}
