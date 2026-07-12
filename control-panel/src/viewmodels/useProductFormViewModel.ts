import { useState, useCallback } from 'react';
import { fetchApi } from '../lib/api';
import type { Product } from '../models/Product';

export interface ProductFormValues {
  id?: string;
  sku: string;
  name: string;
  description: string;
  is_available: boolean;
  status: 'active' | 'inactive' | 'archived';
  price: number;
  weight: number | null;
}

export interface InventorySyncItem {
  shopId: string;
  shopName: string;
  stock: number;
  isNew: boolean;
  isModified: boolean;
  isDeleted: boolean;
  originalStock: number;
}

export function useProductFormViewModel() {
  const [product, setProduct] = useState<Product | null>(null);
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);

  const loadProduct = useCallback(async (slug: string) => {
    setLoading(true);
    setError(null);
    try {
      const res = await fetchApi(`/products/${slug}`);
      setProduct(res);
    } catch (err: any) {
      setError(err.message || 'Failed to load product details');
      setProduct(null);
    } finally {
      setLoading(false);
    }
  }, []);

  const clearProduct = useCallback(() => {
    setProduct(null);
    setError(null);
    setSuccess(false);
  }, []);

  const syncInventories = async (productId: string, list: InventorySyncItem[]) => {
    const promises = list.map(async (item) => {
      try {
        if (item.isNew && !item.isDeleted) {
          // Add inventory
          return await fetchApi(`/shops/${item.shopId}/products/${productId}/inventories`, {
            method: 'POST',
            body: JSON.stringify({ stock: Number(item.stock) }),
          });
        } else if (!item.isNew && item.isDeleted) {
          // Delete inventory
          return await fetchApi(`/shops/${item.shopId}/products/${productId}/inventories`, {
            method: 'DELETE',
          });
        } else if (!item.isNew && item.isModified && !item.isDeleted) {
          // Update inventory
          return await fetchApi(`/shops/${item.shopId}/products/${productId}/inventories`, {
            method: 'PUT',
            body: JSON.stringify({ stock: Number(item.stock) }),
          });
        }
      } catch (err) {
        console.error(`Failed to sync inventory for shop ${item.shopName}:`, err);
      }
    });
    await Promise.all(promises);
  };

  const saveProduct = async (
    data: ProductFormValues,
    inventoriesList: InventorySyncItem[]
  ) => {
    setSaving(true);
    setError(null);
    setSuccess(false);
    try {
      const isEdit = !!data.id;
      
      // 1. Save product
      await fetchApi('/products', {
        method: 'POST',
        body: JSON.stringify({
          id: data.id || undefined,
          sku: data.sku,
          name: data.name,
          description: data.description || null,
          is_available: data.is_available,
          status: data.status,
          price: Number(data.price),
          weight: data.weight ? Number(data.weight) : null,
        }),
      });

      // 2. Resolve product ID (if creating)
      let productId = data.id;
      if (!isEdit) {
        const searchRes = await fetchApi(`/products?name=${encodeURIComponent(data.name)}&limit=1`);
        const createdProduct = searchRes?.products?.[0];
        if (createdProduct && createdProduct.name === data.name) {
          productId = createdProduct.id;
        }
      }

      // 3. Sync inventories
      if (productId) {
        await syncInventories(productId, inventoriesList);
      }

      setSuccess(true);
      return true;
    } catch (err: any) {
      setError(err.message || 'Failed to save product');
      return false;
    } finally {
      setSaving(false);
    }
  };

  const uploadImage = async (file: File) => {
    if (!product) return false;
    setUploading(true);
    setError(null);
    try {
      const formData = new FormData();
      formData.append('image', file);

      await fetchApi(`/products/id/${product.id}/images`, {
        method: 'POST',
        body: formData,
      });

      // Reload product details to show the new image
      const res = await fetchApi(`/products/${product.slug}`);
      setProduct(res);
      return true;
    } catch (err: any) {
      setError(err.message || 'Failed to upload product image');
      return false;
    } finally {
      setUploading(false);
    }
  };

  return {
    product,
    loading,
    saving,
    uploading,
    error,
    success,
    loadProduct,
    clearProduct,
    saveProduct,
    uploadImage,
  };
}
