import { useState, useEffect, useCallback } from 'react';
import { fetchApi } from '../lib/api';
import type { ProductsResponse } from '../models/Product';

export type ProductStatusFilter = 'all' | 'active' | 'inactive' | 'archived';

export function useProductsViewModel() {
  const [data, setData] = useState<ProductsResponse | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [page, setPage] = useState<number>(1);
  const [limit, setLimit] = useState<number>(15);
  const [statusFilter, setStatusFilter] = useState<ProductStatusFilter>('all');

  const fetchProducts = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const url = `/products?page=${page}&limit=${limit}&status=${statusFilter}`;
      const result = await fetchApi(url);
      setData(result);
    } catch (err: any) {
      setError(err.message || 'Failed to fetch products');
      setData(null);
    } finally {
      setLoading(false);
    }
  }, [page, limit, statusFilter]);

  useEffect(() => {
    fetchProducts();
  }, [fetchProducts]);

  const handleSetStatusFilter = (status: ProductStatusFilter) => {
    setStatusFilter(status);
    setPage(1);
  };

  return {
    data,
    loading,
    error,
    page,
    limit,
    statusFilter,
    setStatusFilter: handleSetStatusFilter,
    setPage,
    setLimit,
    refresh: fetchProducts,
  };
}
