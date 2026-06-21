import { useState, useEffect, useCallback } from 'react';
import { fetchApi } from '../lib/api';
import type { ShopAddressesResponse, ShopCouriersResponse, ShopProductsResponse } from '../models/Shop';

const mockShopId = "c3d4e5f6-a7b8-9012-cdef-123456789012";

const mockAddresses: ShopAddressesResponse = {
  shop_id: mockShopId,
  addresses: [
    {
      id: "d4e5f6a7-b8c9-0123-defa-234567890123",
      label: "Main Branch",
      phone: "+6281234567890",
      is_active: true,
      province_id: "32",
      city_id: "3204",
      district_id: "320401",
      village_id: "3204010001",
      full_address: "Jl. Bunga Indah No. 10, Bekasi",
      postal_code: "17520",
      created_at: "2026-01-15T09:00:00Z",
      updated_at: "2026-05-01T12:00:00Z"
    }
  ]
};

const mockCouriers: ShopCouriersResponse = {
  shop_id: mockShopId,
  couriers: [
    { code: "jne", active: true },
    { code: "sicepat", active: false },
    { code: "gojek", active: true },
    { code: "grab", active: true }
  ]
};

const mockProducts: ShopProductsResponse = {
  shop_id: mockShopId,
  products: [
    {
      id: "9886edf6-087b-48e7-b00a-d79dd092e8d4",
      sku: "EVT-ANV-001",
      name: "Anniversary",
      slug: "anniversary",
      description: "A beautiful anniversary bouquet.",
      status: "active",
      price: 85000,
      weight: 1.5,
      inventory: {
        total_stock: 100,
        reserved_stock: 19,
        available: 81
      },
      created_at: "2026-01-10T08:00:00Z",
      updated_at: "2026-06-01T10:30:00Z"
    },
    {
      id: "2ceea56c-352f-4a48-a262-f60e9ee85b1c",
      sku: "EVT-GOP-007",
      name: "Grand Opening",
      slug: "grand-opening",
      description: "A towering two-tier floral stand featuring red anthuriums.",
      status: "active",
      price: 150000,
      weight: 4.5,
      inventory: {
        total_stock: 50,
        reserved_stock: 5,
        available: 45
      },
      created_at: "2026-02-15T09:15:00Z",
      updated_at: null
    }
  ]
};

export function useShopViewModel() {
  const [addresses, setAddresses] = useState<ShopAddressesResponse | null>(null);
  const [couriers, setCouriers] = useState<ShopCouriersResponse | null>(null);
  const [products, setProducts] = useState<ShopProductsResponse | null>(null);
  const [shopId, setShopId] = useState<string | null>(null);
  
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  const fetchShopDetails = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      
      // Fetch shops to get a shopId
      const shopsData = await fetchApi('/shops');
      const id = shopsData?.shops?.[0]?.id;
      
      if (!id) {
        throw new Error('No shop found');
      }
      setShopId(id);

      // Fetch each resource individually and gracefully
      try {
        const addrRes = await fetchApi(`/shops/${id}/addresses`);
        setAddresses(addrRes);
      } catch (err) {
        console.warn('Failed to load shop addresses, using mock', err);
        setAddresses(mockAddresses);
      }

      try {
        const courRes = await fetchApi(`/shops/${id}/couriers`);
        setCouriers(courRes);
      } catch (err) {
        console.warn('Failed to load shop couriers, using mock', err);
        setCouriers(mockCouriers);
      }

      try {
        const prodRes = await fetchApi(`/shops/${id}/products`);
        setProducts(prodRes);
      } catch (err) {
        console.warn('Failed to load shop products, using mock', err);
        setProducts(mockProducts);
      }
      
    } catch (err: any) {
      console.warn('Backend /shops failed or not implemented, falling back to mock data', err);
      setShopId(mockShopId);
      setAddresses(mockAddresses);
      setCouriers(mockCouriers);
      setProducts(mockProducts);
      setError(null);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchShopDetails();
  }, [fetchShopDetails]);

  const createAddress = async (data: any) => {
    const activeShopId = shopId || mockShopId;
    const isActiveValue = data.is_active ? "true" : "false";
    console.log('[createAddress] is_active:', data.is_active, '=> sending:', isActiveValue);
    try {
      setLoading(true);
      await fetchApi(`/shops/${activeShopId}/addresses`, {
        method: 'POST',
        body: JSON.stringify({
          ...data,
          shop_id: activeShopId,
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

  return {
    shopId,
    addresses: addresses?.addresses || [],
    couriers: couriers?.couriers || [],
    products: products?.products || [],
    loading,
    error,
    createAddress,
    refresh: fetchShopDetails
  };
}
