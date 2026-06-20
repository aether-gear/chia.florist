import { useState, useEffect } from 'react';
import { fetchApi } from '../lib/api';
import type { PaymentMethodsResponse, PaymentAccountsResponse, PaymentAccount } from '../models/Payment';

const mockMethods: PaymentMethodsResponse = {
  methods: [
    {
      id: "0137d751-5188-447a-b630-1bf858f4f866",
      name: "QRIS",
      type: "qr_code",
      is_active: true,
      description: "QRIS payment via Midtrans",
      fee_type: "",
      fee_fixed: 0,
      fee_percentage: 0
    },
    {
      id: "5de3fdf1-7cf2-4354-bf31-a288a6706c41",
      name: "GoPay",
      type: "ewallet",
      is_active: true,
      description: "GoPay via Midtrans",
      fee_type: "",
      fee_fixed: 0,
      fee_percentage: 0
    },
    {
      id: "074b02e4-e047-4f60-bdb0-cfeb5481d002",
      name: "DANA",
      type: "ewallet",
      is_active: true,
      description: "DANA via Midtrans",
      fee_type: "",
      fee_fixed: 0,
      fee_percentage: 0
    },
    {
      id: "24ce2aac-bd73-4c29-9ab9-2f53282b2679",
      name: "Mandiri",
      type: "bank_transfer",
      is_active: true,
      description: "Mandiri Bill Payment via Midtrans",
      fee_type: "",
      fee_fixed: 0,
      fee_percentage: 0
    }
  ]
};

const mockAccounts: PaymentAccountsResponse = {
  accounts: [
    {
      id: "5672b98b-3474-4bbe-94de-ccf784ae90dc",
      method_id: "24ce2aac-bd73-4c29-9ab9-2f53282b2679",
      account_name: "Mandiri Reyhan",
      account_number: "1690002799366",
      phone_number: "0895326204046",
      qr_string: null
    },
    {
      id: "01198989-6b57-4005-b7e6-50e797ccca04",
      method_id: "074b02e4-e047-4f60-bdb0-cfeb5481d002",
      account_name: "Dana Ilham",
      account_number: null,
      phone_number: "081291302897",
      qr_string: null
    },
    {
      id: "d8242cd0-7e80-4bbb-95cd-7b0061230ed6",
      method_id: "5de3fdf1-7cf2-4354-bf31-a288a6706c41",
      account_name: "GoPay Ilham",
      account_number: null,
      phone_number: "081291302897",
      qr_string: null
    }
  ]
};

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
        console.warn('Backend /payments failed, falling back to mock data', err);
        setMethods(mockMethods);
        setAccounts(mockAccounts);
        setError(null);
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
    } catch (err) {
      console.error(err);
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
