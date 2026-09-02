export const API_BASE_URL = '/api/core';

export function clearAuthStorage(): void {
  localStorage.removeItem('isAuthenticated');
  localStorage.removeItem('userEmail');
  localStorage.removeItem('authUser');
  localStorage.removeItem('authProfile');
  sessionStorage.removeItem('isAuthenticated');
  sessionStorage.removeItem('userEmail');
  sessionStorage.removeItem('authUser');
  sessionStorage.removeItem('authProfile');
}

const PUBLIC_AUTH_ENDPOINTS = [
  '/auth/staff/signin',
  '/auth/signin',
  '/auth/staff/forgot-password',
  '/auth/forgot-password/verify',
  '/auth/forgot-password/reset',
];

export async function fetchApi(endpoint: string, options: RequestInit = {}) {
  const url = `${API_BASE_URL}${endpoint}`;
  
  const headers: Record<string, string> = {
    'X-Account-Type': 'staff',
  };
  if (!(options.body instanceof FormData)) {
    headers['Content-Type'] = 'application/json';
  }
  Object.assign(headers, options.headers);

  const executeFetch = () => fetch(url, {
    ...options,
    headers,
    credentials: 'include', // Important to send/receive cookies
  });

  let response = await executeFetch();

  // If unauthorized, retry once after a short delay to allow parallel token refresh to update the cookies
  if (response.status === 401) {
    await new Promise((resolve) => setTimeout(resolve, 200));
    response = await executeFetch();
  }

  if (!response.ok) {
    let errorMsg = 'An error occurred';
    try {
      const errorData = await response.json();
      errorMsg = errorData.message || errorMsg;
    } catch {
      // ignore parsing error
    }

    // If 401 on an authenticated endpoint, invalidate the session and redirect to login
    const isPublicAuthEndpoint = PUBLIC_AUTH_ENDPOINTS.some((pub) => endpoint.startsWith(pub));
    if (response.status === 401 && !isPublicAuthEndpoint) {
      clearAuthStorage();
      if (typeof window !== 'undefined') {
        window.dispatchEvent(new CustomEvent('auth:session-expired'));
        if (!window.location.pathname.startsWith('/login')) {
          window.location.href = '/login';
        }
      }
    }

    throw new Error(errorMsg);
  }

  // Handle empty responses
  if (response.status === 204 || response.headers.get('content-length') === '0') {
    return null;
  }

  return response.json();
}

