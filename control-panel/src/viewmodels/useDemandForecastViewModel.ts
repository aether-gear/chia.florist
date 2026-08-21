import { useState, useCallback } from 'react';
import { fetchApi } from '../lib/api';
import type { DemandForecast } from '../models/Analytics';

export function useDemandForecastViewModel() {
  const [forecast, setForecast] = useState<DemandForecast | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  const fetchForecast = useCallback(async (productId: string, shopId?: string) => {
    if (!productId) {
      setForecast(null);
      return;
    }

    try {
      setLoading(true);
      setError(null);

      const params = new URLSearchParams();
      params.append('product_id', productId);
      if (shopId) {
        params.append('shop_id', shopId);
      }

      const response: DemandForecast = await fetchApi(
        `/analytics/forecasts/demand?${params.toString()}`
      );
      setForecast(response);
    } catch (err: any) {
      setError(err.message || 'Failed to generate demand forecast');
      setForecast(null);
    } finally {
      setLoading(false);
    }
  }, []);

  const reset = useCallback(() => {
    setForecast(null);
    setError(null);
    setLoading(false);
  }, []);

  return {
    forecast,
    loading,
    error,
    fetchForecast,
    reset,
  };
}
