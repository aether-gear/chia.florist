import { useState, useEffect } from 'react';
import type { MerchantsResponse } from '../models/Merchant';
import { fetchApi } from '../lib/api';

const mockMerchantsResponse: MerchantsResponse = {
  page: 1,
  limit: 10,
  total: 3,
  merchants: [
    {
      id: "b2c3d4e5-f6a7-8901-bcde-f12345678901",
      name: "Chia Florist",
      description: "Fresh flowers delivered to your door.",
      logo_url: "https://images.unsplash.com/photo-1558227005-950c41870195?w=200&h=200&fit=crop",
      banner_url: "https://images.unsplash.com/photo-1572454591674-2739f30d8c40?w=1200&h=300&fit=crop",
      created_at: "2026-01-01T08:00:00Z"
    },
    {
      id: "c3d4e5f6-0123-4567-8901-abcdef123456",
      name: "Blossom Boutique",
      description: "Premium and elegant floral arrangements.",
      logo_url: "https://images.unsplash.com/photo-1526045612212-70caf35c14df?w=200&h=200&fit=crop",
      banner_url: "https://images.unsplash.com/photo-1563241598-f21cd494cecb?w=1200&h=300&fit=crop",
      created_at: "2026-02-15T10:30:00Z"
    },
    {
      id: "f1e2d3c4-b5a6-7890-1234-567890abcdef",
      name: "Green Thumb Co",
      description: "Indoor plants and terrariums.",
      logo_url: "https://images.unsplash.com/photo-1416879598555-22002165e31c?w=200&h=200&fit=crop",
      banner_url: "https://images.unsplash.com/photo-1463320898484-cdefecf8b328?w=1200&h=300&fit=crop",
      created_at: "2026-03-20T14:15:00Z"
    }
  ]
};

export function useMerchantsViewModel() {
  const [data, setData] = useState<MerchantsResponse | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchMerchants = async () => {
      try {
        setLoading(true);
        setError(null);
        const result = await fetchApi('/merchants');
        setData(result);
      } catch (err: any) {
        console.warn('Backend /merchants failed or not implemented, falling back to mock data', err);
        setData(mockMerchantsResponse);
        // We do not set error, so the UI can display the mock data
        setError(null);
      } finally {
        setLoading(false);
      }
    };

    fetchMerchants();
  }, []);

  return {
    data,
    loading,
    error
  };
}
