import { useState, useEffect, useCallback } from 'react';
import { fetchApi } from '../lib/api';

export function useShopViewModel() {
  const [shops, setShops] = useState<any[]>([]);
  const [selectedShopId, setSelectedShopId] = useState<string | null>(null);
  const [selectedShopInfo, setSelectedShopInfo] = useState<any | null>(null);
  const [addresses, setAddresses] = useState<any[]>([]);
  const [couriers, setCouriers] = useState<any[]>([]);
  const [products, setProducts] = useState<any[]>([]);
  
  const [loading, setLoading] = useState<boolean>(true);
  const [detailsLoading, setDetailsLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [detailsError, setDetailsError] = useState<string | null>(null);

  const fetchShops = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      
      const shopsData = await fetchApi('/shops');
      setShops(shopsData?.shops || []);
    } catch (err: any) {
      setError(err.message || 'Failed to fetch shops');
      setShops([]);
    } finally {
      setLoading(false);
    }
  }, []);

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

  const selectShop = useCallback((shop: any) => {
    setSelectedShopInfo(shop);
    if (shop?.id) {
      fetchShopDetails(shop.id);
    } else {
      setSelectedShopId(null);
      setAddresses([]);
      setCouriers([]);
      setProducts([]);
    }
  }, [fetchShopDetails]);

  useEffect(() => {
    fetchShops();
  }, [fetchShops]);

  const createAddress = async (data: any) => {
    if (!selectedShopId) return false;
    const isActiveValue = data.is_active ? "true" : "false";
    try {
      setDetailsLoading(true);
      await fetchApi(`/shops/${selectedShopId}/addresses`, {
        method: 'POST',
        body: JSON.stringify({
          ...data,
          shop_id: selectedShopId,
          is_active: isActiveValue
        })
      });
      await fetchShopDetails(selectedShopId);
      return true;
    } catch (err) {
      console.error(err);
      return false;
    } finally {
      setDetailsLoading(false);
    }
  };

  const saveShop = async (data: { name: string; description?: string; is_active: string }) => {
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
      const shopsData = await fetchApi('/shops');
      const updatedShops = shopsData?.shops || [];
      setShops(updatedShops);
      
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

  const createShop = async (data: { name: string; description?: string; is_active: string }) => {
    try {
      setLoading(true);
      await fetchApi('/shops', {
        method: 'POST',
        body: JSON.stringify(data)
      });
      // Refresh the list of shops
      const shopsData = await fetchApi('/shops');
      setShops(shopsData?.shops || []);
      return true;
    } catch (err) {
      console.error(err);
      return false;
    } finally {
      setLoading(false);
    }
  };

  return {
    shops,
    selectedShopId,
    selectedShopInfo,
    addresses,
    couriers,
    products,
    loading,
    detailsLoading,
    error,
    detailsError,
    createAddress,
    saveShop,
    createShop,
    selectShop,
    fetchShopDetails,
    refresh: fetchShops
  };
}

