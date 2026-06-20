import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { fetchApi } from '../lib/api';

export interface CreateProductForm {
  sku: string;
  name: string;
  description: string;
  is_available: boolean;
  status: 'active' | 'inactive' | 'archived';
  price: number;
  weight: number | null;
}

export function useCreateProductViewModel() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  const createProduct = async (data: CreateProductForm) => {
    setLoading(true);
    setError(null);
    try {
      await fetchApi('/products', {
        method: 'POST',
        body: JSON.stringify({
          ...data,
          description: data.description || null,
          weight: data.weight ? Number(data.weight) : null,
          price: Number(data.price),
        }),
      });

      // Successfully created
      navigate('/products');
    } catch (err: any) {
      setError(err.message || 'An error occurred while creating the product');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  return {
    createProduct,
    loading,
    error,
  };
}
