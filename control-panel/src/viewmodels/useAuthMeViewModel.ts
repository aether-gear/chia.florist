import { useState, useEffect } from 'react';
import { fetchApi } from '../lib/api';

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

export function useAuthMeViewModel() {
  const [data, setData] = useState<AuthMeResponse | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  // Determine admin status from stored email
  const storedEmail = localStorage.getItem('userEmail') || sessionStorage.getItem('userEmail') || '';
  const isAdmin = ADMIN_EMAILS.includes(storedEmail.toLowerCase().trim());

  useEffect(() => {
    const fetchMe = async () => {
      try {
        setLoading(true);
        setError(null);
        const result = await fetchApi('/auth/me');
        setData(result);
      } catch (err: any) {
        console.error('Backend /auth/me failed', err);
        setData(null);
        setError(err.message || 'Not authenticated');
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
