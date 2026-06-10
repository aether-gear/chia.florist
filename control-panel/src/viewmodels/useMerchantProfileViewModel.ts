import { useState, useEffect } from 'react';
import type { MerchantProfile } from '../models/MerchantProfile';

const mockMerchantProfile: MerchantProfile = {
  identity: {
    accountId: 'ACC-100293',
    merchantId: 'MER-884920',
    name: 'Chia Florist Jakarta',
    slug: 'chia-florist-jkt',
    profilePhoto: 'https://images.unsplash.com/photo-1558227005-950c41870195?w=200&h=200&fit=crop',
    coverBanner: 'https://images.unsplash.com/photo-1572454591674-2739f30d8c40?w=1200&h=300&fit=crop',
    description: 'Premium florist based in Jakarta providing fresh flowers and arrangements for all your special occasions.',
  },
  contact: {
    email: 'hello@chia.florist',
    phone: '+62 811 2233 4455',
    whatsapp: '+62 811 2233 4455',
    customerServiceContact: 'cs@chia.florist',
    address: 'Jl. Sudirman No. 45, Senayan',
    country: 'Indonesia',
    province: 'DKI Jakarta',
    city: 'South Jakarta',
    district: 'Kebayoran Baru',
    postalCode: '12190',
    fullAddress: 'Jl. Sudirman No. 45, Senayan, Kebayoran Baru, South Jakarta, DKI Jakarta 12190, Indonesia',
    latitude: -6.2274,
    longitude: 106.8055,
  },
  settings: {
    preferredLanguage: 'id',
    preferredCurrency: 'IDR',
    timezone: 'Asia/Jakarta',
  },
  operational: {
    openingHours: '08:00',
    closingHours: '20:00',
    businessDays: ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'],
    deliveryRadius: 15,
    pickupAvailable: true,
  },
  financial: {
    bankAccountName: 'PT Chia Florist',
    bankName: 'BCA',
    bankAccountNumber: '1234567890',
    eWalletInformation: 'OVO: 081122334455',
    taxNumber: '01.234.567.8-901.000',
  }
};

export function useMerchantProfileViewModel() {
  const [profile, setProfile] = useState<MerchantProfile | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    // Simulate API fetch delay
    const loadMerchantProfile = async () => {
      try {
        setLoading(true);
        // TODO: Replace with API call when service-core endpoint becomes available
        // const response = await fetch('/api/v1/merchants/profile');
        // const data = await response.json();
        
        await new Promise(resolve => setTimeout(resolve, 800)); // fake delay
        setProfile(mockMerchantProfile);
      } catch (err) {
        setError('Failed to load merchant profile');
        console.error(err);
      } finally {
        setLoading(false);
      }
    };

    loadMerchantProfile();
  }, []);

  const saveProfile = async (data: MerchantProfile) => {
    try {
      setLoading(true);
      // TODO: Replace with actual API call to update profile
      // await fetch('/api/v1/merchants/profile', { method: 'PUT', body: JSON.stringify(data) });
      
      await new Promise(resolve => setTimeout(resolve, 1000));
      setProfile(data); // update local state
      console.log('Profile saved successfully', data);
      return true;
    } catch (err) {
      console.error('Failed to save profile', err);
      return false;
    } finally {
      setLoading(false);
    }
  };

  return {
    profile,
    loading,
    error,
    saveProfile
  };
}
