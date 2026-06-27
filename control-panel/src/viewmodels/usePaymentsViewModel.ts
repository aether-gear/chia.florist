import { useState, useEffect } from 'react';
import { fetchApi } from '../lib/api';
import type { PaymentMethodsResponse, PaymentAccountsResponse, PaymentAccount } from '../models/Payment';

export function usePaymentsViewModel() {
  const [methodsData, setMethods] = useState<PaymentMethodsResponse | null>(null);
  const [accountsData, setAccounts] = useState<PaymentAccountsResponse | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchPayments = async () => {
      try {
        setLoading(true);
        setError(null);
        const [methodsRes, accountsRes] = await Promise.all([
          fetchApi('/payments/methods'),
          fetchApi('/payments/accounts')
        ]);
        
        setMethods(methodsRes);
        setAccounts(accountsRes);
      } catch (err: any) {
        console.error('Backend /payments failed', err);
        setError(err.message || 'Failed to fetch payments data');
      } finally {
        setLoading(false);
      }
    };

    fetchPayments();
  }, []);

  const createPaymentAccount = async (data: Partial<PaymentAccount>) => {
    try {
      setLoading(true);
      const newAccount = await fetchApi('/payments/accounts', {
        method: 'POST',
        body: JSON.stringify({
          ...data,
          is_active: "true"
        })
      });
      
      if (accountsData && newAccount) {
        setAccounts({
          accounts: [...accountsData.accounts, newAccount]
        });
      }
      return true;
    } catch (err: any) {
      console.error('Failed to create payment account', err);
      setError(err.message || 'Failed to create payment account');
      return false;
    } finally {
      setLoading(false);
    }
  };

  return {
    methods: methodsData?.methods || [],
    accounts: accountsData?.accounts || [],
    loading,
    error,
    createPaymentAccount
  };
}
