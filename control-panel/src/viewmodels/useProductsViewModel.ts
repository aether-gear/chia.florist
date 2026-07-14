import { useState, useEffect, useCallback } from 'react';
import { fetchApi } from '../lib/api';
import type { ProductsResponse } from '../models/Product';

export function useProductsViewModel() {
  const [data, setData] = useState<ProductsResponse | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [page, setPage] = useState<number>(1);
  const [limit, setLimit] = useState<number>(15);

  const fetchProducts = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const result = await fetchApi(`/products?page=${page}&limit=${limit}`);
      setData(result);
    } catch (err: any) {
      setError(err.message || 'Failed to fetch products');
      setData(null);
    } finally {
      setLoading(false);
    }
  }, [page, limit]);

  useEffect(() => {
    fetchProducts();
  }, [fetchProducts]);

  return {
    data,
    loading,
    error,
    page,
    limit,
    setPage,
    setLimit,
    refresh: fetchProducts
  };
}

