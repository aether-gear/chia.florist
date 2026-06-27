import { useState, useEffect, useCallback } from 'react';
import { fetchApi } from '../lib/api';
import type { ShopAddressesResponse, ShopCouriersResponse, ShopProductsResponse } from '../models/Shop';

export function useShopViewModel() {
  const [addresses, setAddresses] = useState<ShopAddressesResponse | null>(null);
  const [couriers, setCouriers] = useState<ShopCouriersResponse | null>(null);
  const [products, setProducts] = useState<ShopProductsResponse | null>(null);
  const [shopId, setShopId] = useState<string | null>(null);
  const [shopInfo, setShopInfo] = useState<any | null>(null);
  
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  const fetchShopDetails = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      
      // Fetch shops to get a shopId
      const shopsData = await fetchApi('/shops');
      const shop = shopsData?.shops?.[0];
      const id = shop?.id;
      
      if (!id) {
        throw new Error('No shop found');
      }
      setShopId(id);
      setShopInfo(shop);

      // Fetch each resource individually and propagate errors
      const addrRes = await fetchApi(`/shops/${id}/addresses`);
      setAddresses(addrRes);

      const courRes = await fetchApi(`/shops/${id}/couriers`);
      setCouriers(courRes);

      const prodRes = await fetchApi(`/shops/${id}/products`);
      setProducts(prodRes);
      
    } catch (err: any) {
      setError(err.message || 'Failed to fetch shop details');
      setShopId(null);
      setShopInfo(null);
      setAddresses(null);
      setCouriers(null);
      setProducts(null);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchShopDetails();
  }, [fetchShopDetails]);

  const createAddress = async (data: any) => {
    if (!shopId) return false;
    const isActiveValue = data.is_active ? "true" : "false";
    try {
      setLoading(true);
      await fetchApi(`/shops/${shopId}/addresses`, {
        method: 'POST',
        body: JSON.stringify({
          ...data,
          shop_id: shopId,
          is_active: isActiveValue
        })
      });
      await fetchShopDetails();
      return true;
    } catch (err) {
      console.error(err);
      return false;
    } finally {
      setLoading(false);
    }
  };

  const saveShop = async (data: { name: string; description?: string; is_active: string }) => {
    try {
      setLoading(true);
      await fetchApi('/shops', {
        method: 'POST',
        body: JSON.stringify({
          id: shopId,
          ...data
        })
      });
      await fetchShopDetails();
      return true;
    } catch (err) {
      console.error(err);
      return false;
    } finally {
      setLoading(false);
    }
  };

  return {
    shopId,
    shopInfo,
    addresses: addresses?.addresses || [],
    couriers: couriers?.couriers || [],
    products: products?.products || [],
    loading,
    error,
    createAddress,
    saveShop,
    refresh: fetchShopDetails
  };
}

