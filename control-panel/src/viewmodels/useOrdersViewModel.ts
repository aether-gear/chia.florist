import { useState, useEffect, useCallback, useRef } from 'react';
import { fetchApi } from '../lib/api';
import type { OrdersResponse } from '../models/Order';

export function useOrdersViewModel(fixedShopId?: string) {
  const [data, setData] = useState<OrdersResponse | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [isSwitchingCategory, setIsSwitchingCategory] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  const [page, setPage] = useState<number>(1);
  const [limit, setLimit] = useState<number>(20);
  const [sort, setSort] = useState<string>('latest:desc');
  const [searchNumber, setSearchNumber] = useState<string>('');
  const [statusFilter, setStatusFilterState] = useState<string>('');
  const [shopFilter, setShopFilterState] = useState<string>(fixedShopId || '');
  const [fromDate, setFromDateState] = useState<string>('');
  const [toDate, setToDateState] = useState<string>('');

  const activeRequestId = useRef<number>(0);
  const debounceTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const effectiveShop = fixedShopId !== undefined ? fixedShopId : shopFilter;

  const fetchOrders = useCallback(async (
    targetPage: number,
    targetLimit: number,
    targetSort: string,
    targetSearch: string,
    targetStatus: string,
    targetShop: string,
    targetFromDate: string,
    targetToDate: string
  ) => {
    const requestId = ++activeRequestId.current;
    try {
      setLoading(true);
      setError(null);

      const queryParams = new URLSearchParams();
      queryParams.append('page', targetPage.toString());
      queryParams.append('limit', targetLimit.toString());
      queryParams.append('sort', targetSort);

      if (targetSearch) {
        queryParams.append('number', targetSearch);
      }
      if (targetStatus && targetStatus !== 'all') {
        queryParams.append('status', targetStatus);
      }
      if (targetShop && targetShop !== 'all') {
        queryParams.append('shop', targetShop);
      }
      if (targetFromDate) {
        queryParams.append('from_date', targetFromDate);
      }
      if (targetToDate) {
        queryParams.append('to_date', targetToDate);
      }

      const response = await fetchApi(`/orders?${queryParams.toString()}`);
      if (requestId === activeRequestId.current) {
        setData(response);
      }
    } catch (err: any) {
      if (requestId === activeRequestId.current) {
        console.error('Failed to fetch orders', err);
        setData(null);
        setError(err.message || 'Failed to fetch orders');
      }
    } finally {
      if (requestId === activeRequestId.current) {
        setLoading(false);
        setIsSwitchingCategory(false);
      }
    }
  }, []);

  // Update shop state if fixedShopId prop changes
  useEffect(() => {
    if (fixedShopId !== undefined) {
      setShopFilterState(fixedShopId);
      setPage(1);
    }
  }, [fixedShopId]);

  // Throttled category switching with a locked spinner window
  const setStatusFilter = useCallback((newStatus: string) => {
    setStatusFilterState(newStatus);
    setIsSwitchingCategory(true);
    setLoading(true);
    setPage(1);

    if (debounceTimerRef.current) {
      clearTimeout(debounceTimerRef.current);
    }

    const minThrottleDelay = new Promise(resolve => setTimeout(resolve, 350));

    debounceTimerRef.current = setTimeout(async () => {
      const fetchPromise = fetchOrders(1, limit, sort, searchNumber, newStatus, effectiveShop, fromDate, toDate);
      await Promise.all([fetchPromise, minThrottleDelay]);
    }, 50);
  }, [limit, sort, searchNumber, effectiveShop, fromDate, toDate, fetchOrders]);

  const setShopFilter = useCallback((newShop: string) => {
    if (fixedShopId !== undefined) return;
    setShopFilterState(newShop);
    setLoading(true);
    setPage(1);
    fetchOrders(1, limit, sort, searchNumber, statusFilter, newShop, fromDate, toDate);
  }, [fixedShopId, limit, sort, searchNumber, statusFilter, fromDate, toDate, fetchOrders]);

  const setFromDate = useCallback((date: string) => {
    setFromDateState(date);
    setLoading(true);
    setPage(1);
    fetchOrders(1, limit, sort, searchNumber, statusFilter, effectiveShop, date, toDate);
  }, [limit, sort, searchNumber, statusFilter, effectiveShop, toDate, fetchOrders]);

  const setToDate = useCallback((date: string) => {
    setToDateState(date);
    setLoading(true);
    setPage(1);
    fetchOrders(1, limit, sort, searchNumber, statusFilter, effectiveShop, fromDate, date);
  }, [limit, sort, searchNumber, statusFilter, effectiveShop, fromDate, fetchOrders]);

  useEffect(() => {
    fetchOrders(page, limit, sort, searchNumber, statusFilter, effectiveShop, fromDate, toDate);
    return () => {
      if (debounceTimerRef.current) {
        clearTimeout(debounceTimerRef.current);
      }
    };
  }, [page, limit, sort, searchNumber, effectiveShop, fromDate, toDate, fetchOrders]);

  const refresh = useCallback(() => {
    return fetchOrders(page, limit, sort, searchNumber, statusFilter, effectiveShop, fromDate, toDate);
  }, [fetchOrders, page, limit, sort, searchNumber, statusFilter, effectiveShop, fromDate, toDate]);

  return {
    data,
    loading: loading || isSwitchingCategory,
    isSwitchingCategory,
    error,
    page,
    limit,
    sort,
    searchNumber,
    statusFilter,
    shopFilter: effectiveShop,
    fromDate,
    toDate,
    isShopLocked: fixedShopId !== undefined,
    setPage,
    setLimit,
    setSort,
    setSearchNumber,
    setStatusFilter,
    setShopFilter,
    setFromDate,
    setToDate,
    refresh,
  };
}
