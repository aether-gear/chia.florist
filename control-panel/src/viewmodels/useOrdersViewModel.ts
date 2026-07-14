import { useState, useEffect, useCallback } from 'react';
import { fetchApi } from '../lib/api';
import type { OrdersResponse } from '../models/Order';

export function useOrdersViewModel() {
  const [data, setData] = useState<OrdersResponse | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  const [page, setPage] = useState<number>(1);
  const [limit, setLimit] = useState<number>(20);
  const [sort, setSort] = useState<string>('latest:desc');
  const [searchNumber, setSearchNumber] = useState<string>('');
  const [statusFilter, setStatusFilter] = useState<string>('');

  const fetchOrders = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);

      const queryParams = new URLSearchParams();
      queryParams.append('page', page.toString());
      queryParams.append('limit', limit.toString());
      queryParams.append('sort', sort);

      if (searchNumber) {
        queryParams.append('number', searchNumber);
      }
      if (statusFilter && statusFilter !== 'all') {
        queryParams.append('status', statusFilter);
      }

      const response = await fetchApi(`/orders?${queryParams.toString()}`);
      setData(response);
    } catch (err: any) {
      console.error('Failed to fetch orders', err);
      setData(null);
      setError(err.message || 'Failed to fetch orders');
    } finally {
      setLoading(false);
    }
  }, [page, limit, sort, searchNumber, statusFilter]);

  useEffect(() => {
    fetchOrders();
  }, [fetchOrders]);

  return {
    data,
    loading,
    error,
    page,
    limit,
    sort,
    searchNumber,
    statusFilter,
    setPage,
    setLimit,
    setSort,
    setSearchNumber,
    setStatusFilter,
    refresh: fetchOrders
  };
}
