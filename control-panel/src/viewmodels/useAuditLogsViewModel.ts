import { useState, useEffect, useCallback } from 'react';
import { fetchApi } from '../lib/api';
import type { AuditLogsResponse } from '../models/AuditLog';

export function useAuditLogsViewModel() {
  const [data, setData] = useState<AuditLogsResponse | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  const [page, setPage] = useState<number>(1);
  const [limit, setLimit] = useState<number>(10);
  const [sort, setSort] = useState<string>('date:desc');
  
  // Filters
  const [actionFilter, setActionFilter] = useState<string>('');
  const [userIdFilter, setUserIdFilter] = useState<string>('');
  const [startDate, setStartDate] = useState<string>('');
  const [endDate, setEndDate] = useState<string>('');

  const fetchAuditLogs = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);

      const queryParams = new URLSearchParams();
      queryParams.append('page', page.toString());
      queryParams.append('limit', limit.toString());
      queryParams.append('sort', sort);

      if (actionFilter) {
        queryParams.append('action', actionFilter);
      }
      if (userIdFilter) {
        queryParams.append('user_id', userIdFilter);
      }
      if (startDate) {
        queryParams.append('start_date', startDate);
      }
      if (endDate) {
        queryParams.append('end_date', endDate);
      }

      // Backend route is /api/stats
      const response = await fetchApi(`/api/stats?${queryParams.toString()}`);
      setData(response);
    } catch (err: any) {
      console.error('Failed to fetch audit logs', err);
      setData(null);
      setError(err.message || 'Failed to fetch audit logs');
    } finally {
      setLoading(false);
    }
  }, [page, limit, sort, actionFilter, userIdFilter, startDate, endDate]);

  useEffect(() => {
    fetchAuditLogs();
  }, [fetchAuditLogs]);

  return {
    data,
    loading,
    error,
    page,
    limit,
    sort,
    actionFilter,
    userIdFilter,
    startDate,
    endDate,
    setPage,
    setLimit,
    setSort,
    setActionFilter,
    setUserIdFilter,
    setStartDate,
    setEndDate,
    refresh: fetchAuditLogs
  };
}
