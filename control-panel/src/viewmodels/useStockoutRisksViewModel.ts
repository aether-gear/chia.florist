import { useState, useEffect, useCallback, useMemo } from 'react';
import { fetchApi } from '../lib/api';
import type { StockoutRiskItem, StockoutRisksResponse } from '../models/Analytics';

export type RiskLevelFilter = 'ALL' | 'CRITICAL' | 'WARNING' | 'NORMAL';

export function useStockoutRisksViewModel(initialShopId?: string) {
  const [shopId, setShopId] = useState<string>(initialShopId || '');
  const [riskFilter, setRiskFilter] = useState<RiskLevelFilter>('ALL');
  const [searchQuery, setSearchQuery] = useState<string>('');
  const [risks, setRisks] = useState<StockoutRiskItem[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  const fetchRisks = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);

      const params = new URLSearchParams();
      if (shopId) {
        params.append('shop_id', shopId);
      }

      const queryString = params.toString() ? `?${params.toString()}` : '';
      const response: StockoutRisksResponse = await fetchApi(`/analytics/inventory/stockout-risks${queryString}`);

      if (response && Array.isArray(response.risks)) {
        setRisks(response.risks);
      } else {
        setRisks([]);
      }
    } catch (err: any) {
      setError(err.message || 'Failed to load stockout risks');
    } finally {
      setLoading(false);
    }
  }, [shopId]);

  useEffect(() => {
    fetchRisks();
  }, [fetchRisks]);

  // Derived counts
  const criticalCount = useMemo(
    () => risks.filter((r) => r.risk_level === 'CRITICAL').length,
    [risks]
  );

  const warningCount = useMemo(
    () => risks.filter((r) => r.risk_level === 'WARNING').length,
    [risks]
  );

  const normalCount = useMemo(
    () => risks.filter((r) => r.risk_level === 'NORMAL').length,
    [risks]
  );

  // Average days to stockout across at-risk items
  const avgDaysToStockout = useMemo(() => {
    const atRisk = risks.filter((r) => r.risk_level === 'CRITICAL' || r.risk_level === 'WARNING');
    if (atRisk.length === 0) return null;
    const sum = atRisk.reduce((acc, curr) => acc + curr.estimated_days_to_stockout, 0);
    return Math.max(0, parseFloat((sum / atRisk.length).toFixed(1)));
  }, [risks]);

  // Filtered risks based on search and level filter
  const filteredRisks = useMemo(() => {
    return risks.filter((item) => {
      const matchesFilter =
        riskFilter === 'ALL' || item.risk_level === riskFilter;

      const query = searchQuery.toLowerCase().trim();
      const matchesSearch =
        !query ||
        item.product_name.toLowerCase().includes(query) ||
        item.shop_name.toLowerCase().includes(query) ||
        item.product_id.toLowerCase().includes(query);

      return matchesFilter && matchesSearch;
    });
  }, [risks, riskFilter, searchQuery]);

  return {
    shopId,
    setShopId,
    riskFilter,
    setRiskFilter,
    searchQuery,
    setSearchQuery,
    risks,
    filteredRisks,
    criticalCount,
    warningCount,
    normalCount,
    avgDaysToStockout,
    totalCount: risks.length,
    loading,
    error,
    refresh: fetchRisks,
  };
}
