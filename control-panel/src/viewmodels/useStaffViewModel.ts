import { useState, useEffect, useCallback } from 'react';
import type { StaffListResponse, AddStaffAccountPayload, CreateStaffPayload } from '../models/Staff';
import { fetchApi } from '../lib/api';

export function useStaffViewModel() {
  const [data, setData] = useState<StaffListResponse | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [page, setPage] = useState<number>(1);
  const [limit, setLimit] = useState<number>(20);

  const fetchStaff = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const result = await fetchApi(`/staff?page=${page}&limit=${limit}`);
      
      setData({
        page: result.page || page,
        limit: result.limit || limit,
        total: result.total || 0,
        staff: result.staff || []
      });
    } catch (err: any) {
      console.error('Backend /staff fetch failed', err);
      setError(err.message || 'Failed to load staff list');
    } finally {
      setLoading(false);
    }
  }, [page, limit]);

  useEffect(() => {
    fetchStaff();
  }, [fetchStaff]);

  const addStaffAccount = async (staffId: string, payload: AddStaffAccountPayload) => {
    const res = await fetchApi(`/staff/${staffId}/accounts`, {
      method: 'POST',
      body: JSON.stringify(payload)
    });
    await fetchStaff();
    return res;
  };

  const createStaff = async (payload: CreateStaffPayload) => {
    const res = await fetchApi('/staff', {
      method: 'POST',
      body: JSON.stringify(payload)
    });
    await fetchStaff();
    return res;
  };

  return {
    data,
    staff: data?.staff || [],
    total: data?.total || 0,
    loading,
    error,
    page,
    limit,
    setPage,
    setLimit,
    refresh: fetchStaff,
    addStaffAccount,
    createStaff
  };
}
