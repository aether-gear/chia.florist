import { useState, useEffect } from 'react';
import type { CustomersResponse } from '../models/Customer';
import { fetchApi } from '../lib/api';

export function useCustomersViewModel() {
  const [data, setData] = useState<CustomersResponse | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchCustomers = async () => {
      try {
        setLoading(true);
        setError(null);
        // The endpoint is /customers according to Staff Admin API
        const result = await fetchApi('/customers');
        setData(result);
      } catch (err: any) {
        console.error('Failed to fetch customers', err);
        setError(err.message || 'Failed to load customers');
      } finally {
        setLoading(false);
      }
    };

    fetchCustomers();
  }, []);

  return {
    data,
    loading,
    error
  };
}
