import { useState, useEffect } from 'react';
import type { StaffProfile } from '../models/StaffProfile';
import { fetchApi } from '@/lib/api';

export function useMerchantProfileViewModel() {
  const [profile, setProfile] = useState<StaffProfile | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const loadProfile = async () => {
      try {
        setLoading(true);
        setError(null);
        // The endpoint is /profile according to Staff API
        const data = await fetchApi('/profile');
        setProfile(data.profile);
      } catch (err: any) {
        console.error('Failed to fetch profile', err);
        setError(err.message || 'Failed to load profile');
      } finally {
        setLoading(false);
      }
    };

    loadProfile();
  }, []);

  const saveProfile = async (updateData: { name?: string; phone?: string; avatar_url?: string }) => {
    try {
      setLoading(true);
      const res = await fetchApi('/profile', {
        method: 'PUT',
        body: JSON.stringify(updateData)
      });
      // The API returns the updated profile
      if (res && res.profile) {
        setProfile(res.profile);
      }
      return true;
    } catch (err: any) {
      console.error('Failed to save profile', err);
      setError(err.message || 'Failed to update profile');
      return false;
    } finally {
      setLoading(false);
    }
  };

  return {
    profile,
    loading,
    error,
    saveProfile
  };
}
