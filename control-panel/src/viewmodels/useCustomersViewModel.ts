import { useState, useEffect } from 'react';
import type { CustomersResponse } from '../models/Customer';

const mockCustomersResponse: CustomersResponse = {
  page: 1,
  limit: 10,
  total: 4,
  users: [
    {
      id: "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      name: "Jane Doe",
      username: "janedoe",
      phone: "+6281234567890",
      last_login_at: "2026-06-17T10:00:00Z"
    },
    {
      id: "f8e7d6c5-b4a3-0987-6543-21fedcba0987",
      name: "John Smith",
      username: "johnsmith99",
      phone: "+6289876543210",
      last_login_at: "2026-06-18T14:30:00Z"
    },
    {
      id: "12345678-90ab-cdef-1234-567890abcdef",
      name: "Budi Santoso",
      username: "budis",
      phone: "+628111222333",
      last_login_at: null
    },
    {
      id: "87654321-fedc-ba09-8765-43210fedcba9",
      name: "Siti Rahma",
      username: "sitir",
      phone: "+628555666777",
      last_login_at: "2026-06-19T08:15:00Z"
    }
  ]
};

export function useCustomersViewModel() {
  const [data, setData] = useState<CustomersResponse | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchCustomers = async () => {
      try {
        setLoading(true);
        setError(null);
        const response = await fetch('/api/core/customers');
        if (response.ok) {
          const result = await response.json();
          setData(result);
          return;
        }
        throw new Error('Failed to fetch customers');
      } catch (err: any) {
        console.warn('Backend /customers failed or not implemented, falling back to mock data');
        setData(mockCustomersResponse);
        setError(null);
      } finally {
        setLoading(false);
      }
    };

    fetchCustomers();
  }, []);

  return {
    data,
    loading,
    error
  };
}
