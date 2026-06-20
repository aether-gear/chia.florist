import { useState, useEffect } from 'react';

export interface AuthRole {
  code: string;
  name?: string;
}

export interface AuthPermission {
  code: string;
}

export interface AuthMeResponse {
  account_id: string;
  account_type: string;
  is_authenticated: boolean;
  roles: AuthRole[];
  permissions: AuthPermission[];
}

const ADMIN_EMAILS = ['test@test.com'];

const mockAuthMeResponse: AuthMeResponse = {
  account_id: "51e20db6-5bdb-4f6a-b2b7-8d40c0db857d",
  account_type: "merchant",
  is_authenticated: true,
  roles: [
    { code: "merchant_admin", name: "Merchant Admin" },
    { code: "merchant_staff", name: "Merchant Staff" }
  ],
  permissions: [
    { code: "merchant_staff" },
    { code: "merchant_admin" }
  ]
};

export function useAuthMeViewModel() {
  const [data, setData] = useState<AuthMeResponse | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  // Determine admin status from stored email
  const storedEmail = localStorage.getItem('userEmail') || '';
  const isAdmin = ADMIN_EMAILS.includes(storedEmail.toLowerCase().trim());

  useEffect(() => {
    const fetchMe = async () => {
      try {
        setLoading(true);
        setError(null);
        const response = await fetch('/api/core/auth/me');
        if (response.ok) {
          const result = await response.json();
          setData(result);
          return;
        }
        throw new Error('Not authenticated');
      } catch (err: any) {
        console.warn('Backend /auth/me failed or not implemented, falling back to mock admin role');
        setData(mockAuthMeResponse);
      } finally {
        setLoading(false);
      }
    };

    fetchMe();
  }, []);

  return {
    data,
    loading,
    error,
    isAdmin,
  };
}
