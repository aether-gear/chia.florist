import { useAuth, type AuthRole, type AuthPermission, type AuthMeResponse } from '../context/AuthContext';

export type { AuthRole, AuthPermission, AuthMeResponse };

export function useAuthMeViewModel() {
  const { user, isLoading, error, isAdmin, checkSession } = useAuth();

  return {
    data: user,
    loading: isLoading,
    error,
    isAdmin,
    refresh: checkSession,
  };
}

