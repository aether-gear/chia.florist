import { useState, useEffect, useCallback } from 'react';
import type { StaffProfile } from '../models/StaffProfile';
import { fetchApi } from '@/lib/api';
import { useAuth } from '@/context/AuthContext';

let cachedProfile: StaffProfile | null = null;
let cachedProfileAccountId: string | null = null;
let inFlightProfilePromise: Promise<StaffProfile | null> | null = null;

function getStoredProfile(): StaffProfile | null {
  try {
    const raw = localStorage.getItem('authProfile') || sessionStorage.getItem('authProfile');
    if (raw) return JSON.parse(raw);
  } catch {
    // ignore parse error
  }
  return null;
}

export function clearStaffProfileCache(): void {
  cachedProfile = null;
  cachedProfileAccountId = null;
  inFlightProfilePromise = null;
  localStorage.removeItem('authProfile');
  sessionStorage.removeItem('authProfile');
}

export async function prefetchStaffProfile(targetStorage?: Storage): Promise<StaffProfile | null> {
  try {
    const data = await fetchApi('/profile');
    if (data && data.profile) {
      cachedProfile = data.profile;
      cachedProfileAccountId = data.profile.account_id || data.profile.user_id || null;
      const storage =
        targetStorage ||
        (localStorage.getItem('isAuthenticated') === 'true' ? localStorage : sessionStorage);
      storage.setItem('authProfile', JSON.stringify(data.profile));
      return data.profile;
    }
  } catch (err) {
    console.warn('Failed to prefetch profile:', err);
  }
  return null;
}

export function useStaffProfileViewModel() {
  const { isAuthenticated, user } = useAuth();
  const currentAccountId = user?.account_id || null;

  // If cached profile doesn't match current user, invalidate cache
  if (cachedProfileAccountId && currentAccountId && cachedProfileAccountId !== currentAccountId) {
    clearStaffProfileCache();
  }

  const [profile, setProfile] = useState<StaffProfile | null>(() => {
    if (!isAuthenticated) return null;
    return cachedProfile || getStoredProfile();
  });
  const [loading, setLoading] = useState<boolean>(() => {
    if (!isAuthenticated) return false;
    const initial = cachedProfile || getStoredProfile();
    return !initial;
  });
  const [error, setError] = useState<string | null>(null);

  const loadProfile = useCallback(async (force = false) => {
    if (!isAuthenticated) {
      setProfile(null);
      setLoading(false);
      return null;
    }

    const currentCached = cachedProfile || getStoredProfile();
    if (!force && currentCached && (!currentAccountId || cachedProfileAccountId === currentAccountId)) {
      setProfile(currentCached);
      setLoading(false);
      return currentCached;
    }

    if (inFlightProfilePromise) {
      const p = await inFlightProfilePromise;
      if (p) setProfile(p);
      setLoading(false);
      return p;
    }

    if (!currentCached) {
      setLoading(true);
    }
    setError(null);

    inFlightProfilePromise = (async () => {
      try {
        const data = await fetchApi('/profile');
        if (data && data.profile) {
          cachedProfile = data.profile;
          cachedProfileAccountId = currentAccountId || data.profile.account_id || data.profile.user_id || null;
          const isRemembered = localStorage.getItem('isAuthenticated') === 'true';
          const storage = isRemembered ? localStorage : sessionStorage;
          storage.setItem('authProfile', JSON.stringify(data.profile));
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
  }, [isAuthenticated, currentAccountId]);

  // Sync profile when authentication state or active user account changes
  useEffect(() => {
    if (!isAuthenticated) {
      setProfile(null);
      setLoading(false);
      clearStaffProfileCache();
      return;
    }

    loadProfile();
  }, [isAuthenticated, currentAccountId, loadProfile]);

  // Listen for session expiration events
  useEffect(() => {
    const handleSessionExpired = () => {
      clearStaffProfileCache();
      setProfile(null);
      setLoading(false);
    };

    window.addEventListener('auth:session-expired', handleSessionExpired);
    return () => {
      window.removeEventListener('auth:session-expired', handleSessionExpired);
    };
  }, []);

  const saveProfile = async (updateData: { name?: string; phone?: string; avatar_url?: string }) => {
    try {
      setLoading(true);
      const res = await fetchApi('/profile', {
        method: 'PUT',
        body: JSON.stringify(updateData)
      });
      if (res && res.profile) {
        cachedProfile = res.profile;
        cachedProfileAccountId = currentAccountId || res.profile.account_id || res.profile.user_id || null;
        const isRemembered = localStorage.getItem('isAuthenticated') === 'true';
        const storage = isRemembered ? localStorage : sessionStorage;
        storage.setItem('authProfile', JSON.stringify(res.profile));
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
