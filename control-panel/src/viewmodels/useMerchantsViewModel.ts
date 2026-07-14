import { useState, useEffect, useCallback } from 'react';
import type { MerchantsResponse } from '../models/Merchant';
import { fetchApi } from '../lib/api';

export function useMerchantsViewModel() {
  const [data, setData] = useState<MerchantsResponse | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [page, setPage] = useState<number>(1);
  const [limit, setLimit] = useState<number>(40);

  const fetchMerchants = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      // The endpoint is /staff according to Staff Admin API
      const result = await fetchApi(`/staff?page=${page}&limit=${limit}`);
      
      // Map staff array to merchants key so frontend UI functions correctly
      setData({
        page: result.page || page,
        limit: result.limit || limit,
        total: result.total || 0,
        merchants: result.staff || []
      });
    } catch (err: any) {
      console.error('Backend /staff failed', err);
      setError(err.message || 'Failed to load staff list');
    } finally {
      setLoading(false);
    }
  }, [page, limit]);

  useEffect(() => {
    fetchMerchants();
  }, [fetchMerchants]);

  return {
    data,
    loading,
    error,
    page,
    limit,
    setPage,
    setLimit,
    refresh: fetchMerchants
  };
}
