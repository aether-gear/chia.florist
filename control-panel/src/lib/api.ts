export const API_BASE_URL = '/api/core';

export async function fetchApi(endpoint: string, options: RequestInit = {}) {
  const url = `${API_BASE_URL}${endpoint}`;
  
  const headers: Record<string, string> = {};
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
    throw new Error(errorMsg);
  }

  // Handle empty responses
  if (response.status === 204 || response.headers.get('content-length') === '0') {
    return null;
  }

  return response.json();
}
