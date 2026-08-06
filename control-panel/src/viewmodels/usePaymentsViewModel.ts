import { useState, useEffect } from 'react';
import { fetchApi } from '../lib/api';
import type { PaymentMethodsResponse } from '../models/Payment';

export function usePaymentsViewModel() {
  const [methodsData, setMethods] = useState<PaymentMethodsResponse | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [toggling, setToggling] = useState<boolean>(false);

  const fetchPayments = async () => {
    try {
      setLoading(true);
      setError(null);
      const methodsRes = await fetchApi('/payments/methods');
      
      setMethods(methodsRes);
    } catch (err: any) {
      console.error('Backend /payments failed', err);
      setError(err.message || 'Failed to fetch payments data');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPayments();
  }, []);

  const togglePaymentMethodActive = async (methodId: string, isActive: boolean) => {
    if (toggling) return false;
    try {
      setToggling(true);
      await fetchApi(`/payments/methods/${methodId}`, {
        method: 'PATCH',
        body: JSON.stringify({
          is_active: isActive
        })
      });
      await fetchPayments();
      return true;
    } catch (err: any) {
      console.error('Failed to toggle payment method state', err);
      setError(err.message || 'Failed to update payment method status');
      return false;
    } finally {
      setToggling(false);
    }
  };

  return {
    methods: methodsData?.methods || [],
    loading,
    toggling,
    error,
    togglePaymentMethodActive,
    refetch: fetchPayments
  };
}
