import { useState, useEffect, useCallback } from 'react';
import type { StaffProfile } from '../models/StaffProfile';
import { fetchApi } from '@/lib/api';

let cachedProfile: StaffProfile | null = null;
let inFlightProfilePromise: Promise<StaffProfile | null> | null = null;

export function useStaffProfileViewModel() {
  const [profile, setProfile] = useState<StaffProfile | null>(cachedProfile);
  const [loading, setLoading] = useState<boolean>(!cachedProfile);
  const [error, setError] = useState<string | null>(null);

  const loadProfile = useCallback(async (force = false) => {
    if (!force && cachedProfile) {
      setProfile(cachedProfile);
      setLoading(false);
      return cachedProfile;
    }

    if (inFlightProfilePromise) {
      const p = await inFlightProfilePromise;
      if (p) setProfile(p);
      setLoading(false);
      return p;
    }

    setLoading(true);
    setError(null);

    inFlightProfilePromise = (async () => {
      try {
        // The endpoint is /profile according to Staff API
        const data = await fetchApi('/profile');
        if (data && data.profile) {
          cachedProfile = data.profile;
          return data.profile;
        }
        return null;
      } catch (err: any) {
        console.error('Failed to fetch staff profile', err);
        setError(err.message || 'Failed to load profile');
        return null;
      } finally {
        inFlightProfilePromise = null;
        setLoading(false);
      }
    })();

    const result = await inFlightProfilePromise;
    if (result) {
      setProfile(result);
    }
    return result;
  }, []);

  useEffect(() => {
    loadProfile();
  }, [loadProfile]);

  const saveProfile = async (updateData: { name?: string; phone?: string; avatar_url?: string }) => {
    try {
      setLoading(true);
      const res = await fetchApi('/profile', {
        method: 'PUT',
        body: JSON.stringify(updateData)
      });
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
    refresh: loadProfile,
    saveProfile
  };
}
