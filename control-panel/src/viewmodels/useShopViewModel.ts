import { useState, useEffect, useCallback, useRef } from 'react';
import { fetchApi } from '../lib/api';

export function useShopViewModel(initialShopId?: string) {
  const [shops, setShops] = useState<any[]>([]);
  const [total, setTotal] = useState<number>(0);
  const [page, setPage] = useState<number>(1);
  const [limit, setLimit] = useState<number>(15);
  const [selectedShopId, setSelectedShopId] = useState<string | null>(initialShopId || null);
  const [selectedShopInfo, setSelectedShopInfo] = useState<any | null>(null);
  const [addresses, setAddresses] = useState<any[]>([]);
  const [couriers, setCouriers] = useState<any[]>([]);
  const [products, setProducts] = useState<any[]>([]);
  
  const [loading, setLoading] = useState<boolean>(!initialShopId);
  const [detailsLoading, setDetailsLoading] = useState<boolean>(Boolean(initialShopId));
  const [error, setError] = useState<string | null>(null);
  const [detailsError, setDetailsError] = useState<string | null>(null);

  const fetchShops = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      
      const shopsData = await fetchApi(`/shops?page=${page}&limit=${limit}`);
      setShops(shopsData?.shops || []);
      setTotal(shopsData?.total || 0);
    } catch (err: any) {
      setError(err.message || 'Failed to fetch shops');
      setShops([]);
      setTotal(0);
    } finally {
      setLoading(false);
    }
  }, [page, limit]);

  const activeShopIdRef = useRef<string | null>(null);

  const fetchShopDetails = useCallback(async (id: string) => {
    try {
      setDetailsLoading(true);
      setDetailsError(null);
      setSelectedShopId(id);

      const [addrRes, courRes, prodRes] = await Promise.all([
        fetchApi(`/shops/${id}/addresses`).catch(e => { console.error(e); return null; }),
        fetchApi(`/shops/${id}/couriers`).catch(e => { console.error(e); return null; }),
        fetchApi(`/shops/${id}/products`).catch(e => { console.error(e); return null; })
      ]);

      setAddresses(addrRes?.addresses || []);
      setCouriers(courRes?.couriers || []);
      setProducts(prodRes?.products || []);
    } catch (err: any) {
      setDetailsError(err.message || 'Failed to fetch shop details');
    } finally {
      setDetailsLoading(false);
    }
  }, []);

  const loadShopById = useCallback(async (id: string) => {
    if (!id) {
      activeShopIdRef.current = null;
      setSelectedShopId(null);
      setSelectedShopInfo(null);
      setAddresses([]);
      setCouriers([]);
      setProducts([]);
      return;
    }

    activeShopIdRef.current = id;
    setSelectedShopId(id);
    setDetailsLoading(true);
    setDetailsError(null);

    try {
      let shopErrorMsg: string | null = null;
      const [shopRes, addrRes, courRes, prodRes] = await Promise.all([
        fetchApi(`/shops/${id}`).catch(e => {
          console.error('Failed to fetch shop info', e);
          shopErrorMsg = e?.message || 'Failed to load shop details.';
          return null;
        }),
        fetchApi(`/shops/${id}/addresses`).catch(e => { console.error('Failed to fetch addresses', e); return null; }),
        fetchApi(`/shops/${id}/couriers`).catch(e => { console.error('Failed to fetch couriers', e); return null; }),
        fetchApi(`/shops/${id}/products`).catch(e => { console.error('Failed to fetch products', e); return null; })
      ]);

      // Only apply state if this is still the active shop request
      if (activeShopIdRef.current === id) {
        const shopData = shopRes?.shop || (shopRes && shopRes.id ? shopRes : null);
        if (shopData) {
          setSelectedShopInfo(shopData);
          if (shopData.id) {
            setSelectedShopId(shopData.id);
          }
        } else {
          setDetailsError(shopErrorMsg || 'Failed to load shop details.');
        }

        setAddresses(addrRes?.addresses || []);
        setCouriers(courRes?.couriers || []);
        setProducts(prodRes?.products || []);
      }
    } catch (err: any) {
      if (activeShopIdRef.current === id) {
        setDetailsError(err.message || 'Failed to fetch shop details');
      }
    } finally {
      if (activeShopIdRef.current === id) {
        setDetailsLoading(false);
      }
    }
  }, []);

  const selectShop = useCallback((shop: any) => {
    if (shop?.id) {
      activeShopIdRef.current = shop.id;
      setSelectedShopId(shop.id);
      setSelectedShopInfo(shop);
      fetchShopDetails(shop.id);
    } else {
      activeShopIdRef.current = null;
      setSelectedShopId(null);
      setSelectedShopInfo(null);
      setAddresses([]);
      setCouriers([]);
      setProducts([]);
    }
  }, [fetchShopDetails]);

  useEffect(() => {
    if (!initialShopId) {
      fetchShops();
    }
  }, [initialShopId, fetchShops]);

  const createAddress = async (shopId: string, data: any) => {
    const targetId = shopId || selectedShopId;
    if (!targetId) throw new Error('No shop ID specified');
    const isActiveValue = data.is_active ? "true" : "false";
    try {
      setDetailsLoading(true);
      await fetchApi(`/shops/${targetId}/addresses`, {
        method: 'POST',
        body: JSON.stringify({
          ...data,
          is_active: isActiveValue
        })
      });
      await fetchShopDetails(targetId);
      return true;
    } catch (err: any) {
      console.error(err);
      throw err;
    } finally {
      setDetailsLoading(false);
    }
  };

  const saveShop = async (data: { name: string; description?: string; is_active?: string; approval_status?: string }) => {
    if (!selectedShopId) return false;
    try {
      setDetailsLoading(true);
      await fetchApi('/shops', {
        method: 'POST',
        body: JSON.stringify({
          id: selectedShopId,
          ...data
        })
      });
      // Refresh the list of shops
      const shopsData = await fetchApi(`/shops?page=${page}&limit=${limit}`);
      const updatedShops = shopsData?.shops || [];
      setShops(updatedShops);
      setTotal(shopsData?.total || 0);
      
      // Update selected shop info
      const updated = updatedShops.find((s: any) => s.id === selectedShopId);
      if (updated) {
        setSelectedShopInfo(updated);
      }
      return true;
    } catch (err) {
      console.error(err);
      return false;
    } finally {
      setDetailsLoading(false);
    }
  };

  const createShop = async (data: { name: string; description?: string; is_active?: string; approval_status?: string }) => {
    try {
      setLoading(true);
      await fetchApi('/shops', {
        method: 'POST',
        body: JSON.stringify(data)
      });
      // Refresh the list of shops
      const shopsData = await fetchApi(`/shops?page=${page}&limit=${limit}`);
      setShops(shopsData?.shops || []);
      setTotal(shopsData?.total || 0);
      return true;
    } catch (err) {
      console.error(err);
      return false;
    } finally {
      setLoading(false);
    }
  };

  const deleteShop = async (shopId: string) => {
    try {
      setLoading(true);
      await fetchApi(`/shops/${shopId}`, {
        method: 'DELETE',
      });
      if (selectedShopId === shopId) {
        setSelectedShopId(null);
        setSelectedShopInfo(null);
      }
      await fetchShops();
      return true;
    } catch (err: any) {
      console.error(err);
      throw err;
    } finally {
      setLoading(false);
    }
  };


  const updateInventory = async (shopId: string, productId: string, stock: number) => {
    try {
      setDetailsLoading(true);
      await fetchApi(`/shops/${shopId}/products/${productId}/inventories`, {
        method: 'PUT',
        body: JSON.stringify({ stock }),
      });
      await fetchShopDetails(shopId);
      return true;
    } catch (err: any) {
      console.error(err);
      throw err;
    } finally {
      setDetailsLoading(false);
    }
  };

  const removeInventory = async (shopId: string, productId: string) => {
    try {
      setDetailsLoading(true);
      await fetchApi(`/shops/${shopId}/products/${productId}/inventories`, {
        method: 'DELETE',
      });
      await fetchShopDetails(shopId);
      return true;
    } catch (err: any) {
      console.error(err);
      throw err;
    } finally {
      setDetailsLoading(false);
    }
  };

  const updateAddress = async (shopId: string, addressId: string, data: any) => {
    const isActiveValue = data.is_active ? "true" : "false";
    try {
      setDetailsLoading(true);
      await fetchApi(`/shops/${shopId}/addresses/${addressId}`, {
        method: 'PUT',
        body: JSON.stringify({
          ...data,
          shop_id: shopId,
          is_active: isActiveValue,
        }),
      });
      await fetchShopDetails(shopId);
      return true;
    } catch (err: any) {
      console.error(err);
      throw err;
    } finally {
      setDetailsLoading(false);
    }
  };

  const deleteAddress = async (shopId: string, addressId: string) => {
    try {
      setDetailsLoading(true);
      await fetchApi(`/shops/${shopId}/addresses/${addressId}`, {
        method: 'DELETE',
      });
      await fetchShopDetails(shopId);
      return true;
    } catch (err: any) {
      console.error(err);
      throw err;
    } finally {
      setDetailsLoading(false);
    }
  };

  const updateCourier = async (
    shopId: string,
    code: string,
    data: { active: boolean; name?: string; location_address?: string }
  ) => {
    try {
      setDetailsLoading(true);
      await fetchApi(`/shops/${shopId}/couriers/${code}`, {
        method: 'PUT',
        body: JSON.stringify(data),
      });
      await fetchShopDetails(shopId);
      return true;
    } catch (err: any) {
      console.error(err);
      throw err;
    } finally {
      setDetailsLoading(false);
    }
  };

  const verifyCourier = async (
    shopId: string,
    code: string,
    action: 'verify' | 'reject',
    rejectionReason?: string
  ) => {
    try {
      setDetailsLoading(true);
      await fetchApi(`/shops/${shopId}/couriers/${code}/verify`, {
        method: 'POST',
        body: JSON.stringify({
          action,
          rejection_reason: rejectionReason || undefined,
        }),
      });
      await fetchShopDetails(shopId);
      return true;
    } catch (err: any) {
      console.error(err);
      throw err;
    } finally {
      setDetailsLoading(false);
    }
  };

  return {
    shops,
    total,
    page,
    limit,
    selectedShopId,
    selectedShopInfo,
    addresses,
    couriers,
    products,
    loading,
    detailsLoading,
    error,
    detailsError,
    setPage,
    setLimit,
    createAddress,
    updateAddress,
    deleteAddress,
    updateCourier,
    verifyCourier,
    saveShop,
    createShop,
    deleteShop,
    updateInventory,
    removeInventory,
    selectShop,
    fetchShopDetails,
    loadShopById,
    refresh: fetchShops
  };
}

