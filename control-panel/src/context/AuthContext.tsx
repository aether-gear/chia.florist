import React, { createContext, useContext, useState, useEffect, useCallback } from 'react';
import { fetchApi } from '../lib/api';
import { clearStaffProfileCache, prefetchStaffProfile } from '../viewmodels/useStaffProfileViewModel';

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

interface AuthContextType {
  isAuthenticated: boolean;
  userEmail: string;
  user: AuthMeResponse | null;
  isAdmin: boolean;
  isLoading: boolean;
  error: string | null;
  login: (email: string, rememberMe: boolean) => Promise<void>;
  logout: () => Promise<void>;
  invalidateSession: () => void;
  checkSession: () => Promise<boolean>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

const ADMIN_EMAILS = ['test@test.com'];

export function getStoredAuthState(): {
  isAuthenticated: boolean;
  userEmail: string;
  user: AuthMeResponse | null;
} {
  const localAuth = localStorage.getItem('isAuthenticated') === 'true';
  const sessionAuth = sessionStorage.getItem('isAuthenticated') === 'true';
  const email = localStorage.getItem('userEmail') || sessionStorage.getItem('userEmail') || '';
  let storedUser: AuthMeResponse | null = null;
  try {
    const rawUser = localStorage.getItem('authUser') || sessionStorage.getItem('authUser');
    if (rawUser) {
      storedUser = JSON.parse(rawUser);
    }
  } catch {
    // ignore parse error
  }
  return {
    isAuthenticated: localAuth || sessionAuth,
    userEmail: email,
    user: storedUser,
  };
}

export function clearAuthStorage(): void {
  localStorage.removeItem('isAuthenticated');
  localStorage.removeItem('userEmail');
  localStorage.removeItem('authUser');
  localStorage.removeItem('authProfile');
  sessionStorage.removeItem('isAuthenticated');
  sessionStorage.removeItem('userEmail');
  sessionStorage.removeItem('authUser');
  sessionStorage.removeItem('authProfile');
  clearStaffProfileCache();
}

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const initialStorage = getStoredAuthState();
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(initialStorage.isAuthenticated);
  const [userEmail, setUserEmail] = useState<string>(initialStorage.userEmail);
  const [user, setUser] = useState<AuthMeResponse | null>(initialStorage.user);
  const [isLoading, setIsLoading] = useState<boolean>(initialStorage.isAuthenticated && !initialStorage.user);
  const [error, setError] = useState<string | null>(null);

  const hasAdminRole = user?.roles?.some((r) => r.code === 'staff_admin') ?? false;
  const isAdmin = hasAdminRole || (userEmail ? ADMIN_EMAILS.includes(userEmail.toLowerCase().trim()) : false);

  const invalidateSession = useCallback(() => {
    clearAuthStorage();
    setIsAuthenticated(false);
    setUserEmail('');
    setUser(null);
    setError(null);
    setIsLoading(false);
    clearStaffProfileCache();
  }, []);

  const inFlightSessionPromiseRef = React.useRef<Promise<boolean> | null>(null);

  const checkSession = useCallback(async (): Promise<boolean> => {
    if (inFlightSessionPromiseRef.current) {
      return inFlightSessionPromiseRef.current;
    }

    const sessionPromise = (async () => {
      try {
        setError(null);
        const result: AuthMeResponse = await fetchApi('/auth/staff/me');
        if (result && result.is_authenticated) {
          setUser(result);
          setIsAuthenticated(true);
          const email = localStorage.getItem('userEmail') || sessionStorage.getItem('userEmail') || '';
          setUserEmail(email);
          const isRemembered = localStorage.getItem('isAuthenticated') === 'true';
          const storage = isRemembered ? localStorage : sessionStorage;
          storage.setItem('authUser', JSON.stringify(result));
          return true;
        } else {
          invalidateSession();
          return false;
        }
      } catch (err: any) {
        console.warn('Session verification failed:', err?.message || err);
        invalidateSession();
        return false;
      } finally {
        setIsLoading(false);
        inFlightSessionPromiseRef.current = null;
      }
    })();

    inFlightSessionPromiseRef.current = sessionPromise;
    return sessionPromise;
  }, [invalidateSession]);

  const login = useCallback(async (email: string, rememberMe: boolean) => {
    // Clean both storages first to prevent mixed states
    clearAuthStorage();

    const storage = rememberMe ? localStorage : sessionStorage;
    storage.setItem('isAuthenticated', 'true');
    storage.setItem('userEmail', email);

    setIsAuthenticated(true);
    setUserEmail(email);

    // Warm up session and staff profile in parallel before navigating to dashboard
    await Promise.allSettled([
      checkSession(),
      prefetchStaffProfile(storage),
    ]);
  }, [checkSession]);

  const logout = useCallback(async () => {
    try {
      await fetchApi('/auth/staff/logout', { method: 'POST' });
    } catch (err) {
      console.warn('Logout API error:', err);
    } finally {
      invalidateSession();
    }
  }, [invalidateSession]);

  // Initial check on mount
  useEffect(() => {
    const current = getStoredAuthState();
    if (current.isAuthenticated) {
      checkSession();
    } else {
      setIsLoading(false);
    }
  }, [checkSession]);

  // Listen for session expiration events dispatched by API interceptors or other tabs
  useEffect(() => {
    const handleSessionExpired = () => {
      invalidateSession();
    };

    const handleStorageChange = (e: StorageEvent) => {
      if (e.key === 'isAuthenticated' && e.newValue !== 'true') {
        invalidateSession();
      }
    };

    window.addEventListener('auth:session-expired', handleSessionExpired);
    window.addEventListener('storage', handleStorageChange);

    return () => {
      window.removeEventListener('auth:session-expired', handleSessionExpired);
      window.removeEventListener('storage', handleStorageChange);
    };
  }, [invalidateSession]);

  return (
    <AuthContext.Provider
      value={{
        isAuthenticated,
        userEmail,
        user,
        isAdmin,
        isLoading,
        error,
        login,
        logout,
        invalidateSession,
        checkSession,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
