import { useState, useEffect, useCallback } from 'react';
import { fetchApi } from '../lib/api';
import type { ProductStatsResponse } from '../models/Product';

export function useProductStatsViewModel() {
  const [data, setData] = useState<ProductStatsResponse | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  const fetchStats = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      // Fetch up to 100 products for full chart comparison
      const result = await fetchApi('/products/stats?page=1&limit=100');
      setData(result);
    } catch (err: any) {
      setError(err.message || 'Failed to fetch product stats');
      setData(null);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchStats();
  }, [fetchStats]);

  return {
    data,
    loading,
    error,
    refresh: fetchStats,
  };
}
