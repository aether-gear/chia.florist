import { useState, useEffect, useCallback } from 'react';
import { fetchApi } from '../lib/api';
import type { Order } from '../models/Order';

export function useOrderDetailViewModel(orderId?: string) {
  const [order, setOrder] = useState<Order | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  const fetchOrder = useCallback(async () => {
    if (!orderId) {
      setOrder(null);
      setLoading(false);
      return;
    }
    try {
      setLoading(true);
      setError(null);
      const response = await fetchApi(`/orders/${orderId}`);
      setOrder(response);
    } catch (err: any) {
      console.error('Failed to fetch order detail', err);
      setError(err.message || 'Failed to load order details');
      setOrder(null);
    } finally {
      setLoading(false);
    }
  }, [orderId]);

  useEffect(() => {
    fetchOrder();
  }, [fetchOrder]);

  return {
    order,
    loading,
    error,
    refresh: fetchOrder,
    setOrder,
  };
}
