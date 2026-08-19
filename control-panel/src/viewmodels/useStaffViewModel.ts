import { useState, useEffect, useCallback } from 'react';
import type {
  StaffListResponse,
  StaffAccountMember,
  StaffAccountsResponse,
  AddStaffAccountPayload,
  CreateStaffPayload,
  UpdateStaffPayload,
  StaffShopPermission,
  SaveStaffShopPermissionPayload,
  StaffShopPermissionsResponse,
} from '../models/Staff';
import { fetchApi } from '../lib/api';


export function useStaffViewModel() {
  const [data, setData] = useState<StaffListResponse | null>(null);
  const [accountsMap, setAccountsMap] = useState<Record<string, StaffAccountMember[]>>({});
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [page, setPage] = useState<number>(1);
  const [limit, setLimit] = useState<number>(20);

  const fetchStaffAccounts = useCallback(async (staffId: string): Promise<StaffAccountMember[]> => {
    try {
      const res: StaffAccountsResponse = await fetchApi(`/staff/${staffId}/accounts`);
      const accounts = res.accounts || [];
      setAccountsMap((prev) => ({ ...prev, [staffId]: accounts }));
      return accounts;
    } catch (err) {
      console.error(`Failed to fetch accounts for staff ${staffId}`, err);
      return [];
    }
  }, []);

  const fetchStaff = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const result: StaffListResponse = await fetchApi(`/staff?page=${page}&limit=${limit}`);
      const staffList = result.staff || [];
      // Concurrently fetch accounts for each staff member to determine binding status
      // and wait until BOTH staff data and all account bindings have resolved
      const accountsMapResult: Record<string, StaffAccountMember[]> = {};
      if (staffList.length > 0) {
        const results = await Promise.all(
          staffList.map((s) =>
            fetchApi(`/staff/${s.id}/accounts`)
              .then((res: StaffAccountsResponse) => ({ id: s.id, accounts: res.accounts || [] }))
              .catch(() => ({ id: s.id, accounts: [] }))
          )
        );

        results.forEach((r) => {
          accountsMapResult[r.id] = r.accounts;
        });
      }

      setAccountsMap((prev) => ({ ...prev, ...accountsMapResult }));
      setData({
        page: result.page || page,
        limit: result.limit || limit,
        total: result.total || 0,
        staff: staffList,
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
      body: JSON.stringify(payload),
    });
    await fetchStaffAccounts(staffId);
    return res;
  };

  const removeStaffAccount = async (staffId: string, accountId: string) => {
    const res = await fetchApi(`/staff/${staffId}/accounts/${accountId}`, {
      method: 'DELETE',
    });
    await fetchStaffAccounts(staffId);
    return res;
  };

  const createStaff = async (payload: CreateStaffPayload) => {
    const res = await fetchApi('/staff', {
      method: 'POST',
      body: JSON.stringify(payload),
    });
    await fetchStaff();
    return res;
  };

  const createStaffWithAccount = async (
    staffPayload: CreateStaffPayload,
    accountPayload: AddStaffAccountPayload
  ) => {
    // 1. Create staff unit
    await fetchApi('/staff', {
      method: 'POST',
      body: JSON.stringify(staffPayload),
    });

    // 2. Fetch latest staff list to retrieve the newly created staff entity ID
    const listRes: StaffListResponse = await fetchApi('/staff?page=1&limit=5&sort=latest:desc');
    const newlyCreatedStaff = listRes.staff?.find(
      (s) => s.name.trim().toLowerCase() === staffPayload.name.trim().toLowerCase()
    ) || listRes.staff?.[0];

    if (!newlyCreatedStaff) {
      throw new Error('Staff entity created, but unable to locate new staff ID for account binding.');
    }

    // 3. Bind user account to the new staff entity
    const accountRes = await fetchApi(`/staff/${newlyCreatedStaff.id}/accounts`, {
      method: 'POST',
      body: JSON.stringify(accountPayload),
    });

    await fetchStaff();
    return accountRes;
  };

  const updateStaff = async (staffId: string, payload: UpdateStaffPayload) => {
    const res = await fetchApi(`/staff/${staffId}`, {
      method: 'PUT',
      body: JSON.stringify(payload),
    });
    await fetchStaff();
    return res;
  };

  const deleteStaff = async (staffId: string) => {
    const res = await fetchApi(`/staff/${staffId}`, {
      method: 'DELETE',
    });
    await fetchStaff();
    return res;
  };

  const fetchStaffShopPermissions = useCallback(async (staffId: string): Promise<StaffShopPermission[]> => {
    try {
      const res: StaffShopPermissionsResponse = await fetchApi(`/staff/${staffId}/shops`);
      return res.permissions || [];
    } catch (err) {
      console.error(`Failed to fetch shop permissions for staff ${staffId}`, err);
      return [];
    }
  }, []);

  const saveStaffShopPermission = async (staffId: string, payload: SaveStaffShopPermissionPayload) => {
    const res = await fetchApi(`/staff/${staffId}/shops`, {
      method: 'POST',
      body: JSON.stringify(payload),
    });
    return res;
  };

  const deleteStaffShopPermission = async (staffId: string, shopId: string) => {
    const res = await fetchApi(`/staff/${staffId}/shops/${shopId}`, {
      method: 'DELETE',
    });
    return res;
  };

  return {
    data,
    staff: data?.staff || [],
    total: data?.total || 0,
    accountsMap,
    loading,
    error,
    page,
    limit,
    setPage,
    setLimit,
    refresh: fetchStaff,
    fetchStaffAccounts,
    fetchStaffShopPermissions,
    saveStaffShopPermission,
    deleteStaffShopPermission,
    addStaffAccount,
    removeStaffAccount,
    createStaff,
    createStaffWithAccount,
    updateStaff,
    deleteStaff,
  };
}

