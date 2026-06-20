import { useState, useEffect } from 'react';
import type { ProductsResponse } from '../models/Product';

const mockProductsResponse: ProductsResponse = {
  limit: 10,
  page: 1,
  total: 2,
  products: [
    {
      id: "9886edf6-087b-48e7-b00a-d79dd092e8d4",
      sku: "EVT-ANV-001",
      name: "Anniversary",
      slug: "anniversary",
      status: "active",
      is_available: true,
      price: 85000,
      stock: 1291,
      banner: {
        thumbnail: "https://images.unsplash.com/photo-1563241598-f21cd494cecb?w=100&h=100&fit=crop",
        preview: null,
        detail: null
      },
      availability: [
        { slug: "Chia Medan Satria", name: "chia-medan-satria", stock: 981 },
        { slug: "Chia Cikarang", name: "chia-cikarang", stock: 310 }
      ]
    },
    {
      id: "2ceea56c-352f-4a48-a262-f60e9ee85b1c",
      sku: "EVT-GOP-007",
      name: "Grand Opening",
      slug: "grand-opening",
      status: "active",
      is_available: true,
      price: 150000,
      stock: 837,
      banner: {
        thumbnail: "https://images.unsplash.com/photo-1582735689369-4fe89db7114c?w=100&h=100&fit=crop",
        preview: "https://images.unsplash.com/photo-1582735689369-4fe89db7114c?w=600&h=600&fit=crop",
        detail: null
      },
      availability: [
        { slug: "Chia Cipinang", name: "chia-cipinang", stock: 109 },
        { slug: "Chia Medan Satria", name: "chia-medan-satria", stock: 728 }
      ]
    }
  ]
};

export function useProductsViewModel() {
  const [data, setData] = useState<ProductsResponse | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchProducts = async () => {
      try {
        setLoading(true);
        setError(null);
        const response = await fetch('/api/core/products?page=1&limit=10');
        if (response.ok) {
          const result = await response.json();
          setData(result);
          return;
        }
        throw new Error('Failed to fetch products');
      } catch (err: any) {
        console.warn('Backend /products failed or not implemented, falling back to mock data');
        setData(mockProductsResponse);
        setError(null);
      } finally {
        setLoading(false);
      }
    };

    fetchProducts();
  }, []);

  return {
    data,
    loading,
    error
  };
}
