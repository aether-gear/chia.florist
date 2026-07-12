import { useState, useEffect, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { fetchApi } from '../lib/api';
import type { Product } from '../models/Product';

export interface EditProductForm {
  sku: string;
  name: string;
  description: string;
  is_available: boolean;
  status: 'active' | 'inactive' | 'archived';
  price: number;
  weight: number | null;
}

export function useEditProductViewModel(slug: string | undefined) {
  const [product, setProduct] = useState<Product | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [saving, setSaving] = useState(false);
  const [saveError, setSaveError] = useState<string | null>(null);

  const [uploading, setUploading] = useState(false);
  const [uploadError, setUploadError] = useState<string | null>(null);

  const [deleting, setDeleting] = useState(false);
  const [deleteError, setDeleteError] = useState<string | null>(null);

  const navigate = useNavigate();

  const fetchProductDetail = useCallback(async () => {
    if (!slug) return;
    try {
      setLoading(true);
      setError(null);
      const res = await fetchApi(`/products/${slug}`);
      setProduct(res);
    } catch (err: any) {
      setError(err.message || 'Failed to fetch product details');
      setProduct(null);
    } finally {
      setLoading(false);
    }
  }, [slug]);

  useEffect(() => {
    fetchProductDetail();
  }, [fetchProductDetail]);

  const updateProduct = async (data: EditProductForm) => {
    if (!product) return false;
    setSaving(true);
    setSaveError(null);
    try {
      await fetchApi('/products', {
        method: 'POST',
        body: JSON.stringify({
          id: product.id,
          sku: data.sku,
          name: data.name,
          description: data.description || null,
          is_available: data.is_available,
          status: data.status,
          price: Number(data.price),
          weight: data.weight ? Number(data.weight) : null,
        }),
      });
      await fetchProductDetail();
      return true;
    } catch (err: any) {
      setSaveError(err.message || 'An error occurred while saving the product');
      return false;
    } finally {
      setSaving(false);
    }
  };

  const uploadImage = async (file: File) => {
    if (!product) return false;
    setUploading(true);
    setUploadError(null);
    try {
      const formData = new FormData();
      formData.append('image', file);

      await fetchApi(`/products/id/${product.id}/images`, {
        method: 'POST',
        body: formData,
      });

      await fetchProductDetail();
      return true;
    } catch (err: any) {
      setUploadError(err.message || 'Failed to upload product image');
      return false;
    } finally {
      setUploading(false);
    }
  };

  const deleteProduct = async () => {
    if (!product) return false;
    setDeleting(true);
    setDeleteError(null);
    try {
      await fetchApi(`/products/id/${product.id}`, {
        method: 'DELETE',
      });
      navigate('/products');
      return true;
    } catch (err: any) {
      setDeleteError(err.message || 'Failed to delete product');
      return false;
    } finally {
      setDeleting(false);
    }
  };

  return {
    product,
    loading,
    error,
    saving,
    saveError,
    uploading,
    uploadError,
    deleting,
    deleteError,
    updateProduct,
    uploadImage,
    deleteProduct,
    refresh: fetchProductDetail,
  };
}
