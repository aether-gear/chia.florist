import { useState, useEffect } from 'react';
import type { MerchantsResponse } from '../models/Merchant';
import { fetchApi } from '../lib/api';

export function useMerchantsViewModel() {
  const [data, setData] = useState<MerchantsResponse | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchMerchants = async () => {
      try {
        setLoading(true);
        setError(null);
        // The endpoint is /staff according to Staff Admin API
        const result = await fetchApi('/staff');
        
        // Map staff array to merchants key so frontend UI functions correctly
        setData({
          page: result.page || 1,
          limit: result.limit || 10,
          total: result.total || 0,
          merchants: result.staff || []
        });
      } catch (err: any) {
        console.error('Backend /staff failed', err);
        setError(err.message || 'Failed to load staff list');
      } finally {
        setLoading(false);
      }
    };

    fetchMerchants();
  }, []);

  return {
    data,
    loading,
    error
  };
}
