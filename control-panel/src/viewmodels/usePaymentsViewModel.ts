import { useState, useEffect } from 'react';
import { fetchApi } from '../lib/api';
import type { PaymentMethodsResponse, PaymentAccountsResponse, PaymentAccount, PaymentMethod } from '../models/Payment';

export function usePaymentsViewModel() {
  const [methodsData, setMethods] = useState<PaymentMethodsResponse | null>(null);
  const [accountsData, setAccounts] = useState<PaymentAccountsResponse | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

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

  useEffect(() => {
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

  const savePaymentMethodWithInstruction = async (
    methodData: Omit<Partial<PaymentMethod>, 'fee_percentage'> & { fee_amount?: string; fee_percentage?: string },
    instructionContent?: string
  ) => {
    try {
      setLoading(true);
      
      // Generate ID on client side if not present (Create mode) so we can link the instruction to the same ID
      const methodId = methodData.id || (typeof crypto !== 'undefined' && crypto.randomUUID ? crypto.randomUUID() : 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
        const r = (Math.random() * 16) | 0;
        const v = c === 'x' ? r : (r & 0x3) | 0x8;
        return v.toString(16);
      }));

      // 1. Save payment method
      await fetchApi('/payments/methods', {
        method: 'POST',
        body: JSON.stringify({
          id: methodId,
          name: methodData.name,
          code: methodData.code,
          provider: methodData.provider,
          type: methodData.type,
          is_active: methodData.is_active ? "true" : "false",
          description: methodData.description,
          fee_type: methodData.fee_type,
          fee_amount: methodData.fee_amount,
          fee_percentage: methodData.fee_percentage,
        })
      });

      // 2. Save payment instruction if provided
      if (instructionContent !== undefined) {
        await fetchApi(`/payments/methods/${methodId}/instruction`, {
          method: 'POST',
          body: JSON.stringify({
            content: instructionContent
          })
        });
      }

      // Refresh the methods and accounts
      await fetchPayments();
      return true;
    } catch (err: any) {
      console.error('Failed to save payment method with instruction', err);
      setError(err.message || 'Failed to save payment method with instruction');
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
    createPaymentAccount,
    savePaymentMethod: savePaymentMethodWithInstruction,
    refetch: fetchPayments
  };
}
