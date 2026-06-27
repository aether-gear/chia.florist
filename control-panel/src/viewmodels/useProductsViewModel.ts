import { useState, useEffect, useCallback } from 'react';
import { fetchApi } from '../lib/api';
import type { ProductsResponse } from '../models/Product';

export function useProductsViewModel() {
  const [data, setData] = useState<ProductsResponse | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  const fetchProducts = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const result = await fetchApi('/products?page=1&limit=100');
      setData(result);
    } catch (err: any) {
      setError(err.message || 'Failed to fetch products');
      setData(null);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchProducts();
  }, [fetchProducts]);

  return {
    data,
    loading,
    error,
    refresh: fetchProducts
  };
}

