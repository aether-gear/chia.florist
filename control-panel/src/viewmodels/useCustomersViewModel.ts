import { useState, useEffect, useCallback } from 'react';
import type { CustomersResponse } from '../models/Customer';
import { fetchApi } from '../lib/api';

export function useCustomersViewModel() {
  const [data, setData] = useState<CustomersResponse | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [page, setPage] = useState<number>(1);
  const [limit, setLimit] = useState<number>(40);

  const fetchCustomers = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      // The endpoint is /customers according to Staff Admin API
      const result = await fetchApi(`/customers?page=${page}&limit=${limit}`);
      setData(result);
    } catch (err: any) {
      console.error('Failed to fetch customers', err);
      setError(err.message || 'Failed to load customers');
    } finally {
      setLoading(false);
    }
  }, [page, limit]);

  useEffect(() => {
    fetchCustomers();
  }, [fetchCustomers]);

  return {
    data,
    loading,
    error,
    page,
    limit,
    setPage,
    setLimit,
    refresh: fetchCustomers
  };
}
