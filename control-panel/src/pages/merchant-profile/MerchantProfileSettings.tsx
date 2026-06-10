import { useEffect } from 'react';
import { useForm, FormProvider } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { Button } from '@/components/ui/button';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Save, AlertCircle, Loader2 } from 'lucide-react';

import { useMerchantProfileViewModel } from '@/viewmodels/useMerchantProfileViewModel';
import type { MerchantProfile } from '@/models/MerchantProfile';

import { ProfileIdentitySection } from './components/ProfileIdentitySection';
import { ProfileContactSection } from './components/ProfileContactSection';
import { ProfileSettingsSection } from './components/ProfileSettingsSection';
import { ProfileOperationalSection } from './components/ProfileOperationalSection';
import { ProfileFinancialSection } from './components/ProfileFinancialSection';

// Optional: basic validation schema with zod
const profileSchema = z.object({
  identity: z.object({
    accountId: z.string().optional(),
    merchantId: z.string().optional(),
    name: z.string().min(2, 'Name is required (min 2 chars)'),
    slug: z.string().min(2, 'Slug is required'),
    profilePhoto: z.string().url('Must be a valid URL').or(z.literal('')),
    coverBanner: z.string().url('Must be a valid URL').or(z.literal('')),
    description: z.string().optional(),
  }),
  contact: z.object({
    email: z.string().email('Invalid email address'),
    phone: z.string().min(5, 'Phone is required'),
    whatsapp: z.string().optional(),
    customerServiceContact: z.string().optional(),
    address: z.string().min(5, 'Address is required'),
    country: z.string().min(2, 'Country is required'),
    province: z.string().min(2, 'Province is required'),
    city: z.string().min(2, 'City is required'),
    district: z.string().optional(),
    postalCode: z.string().optional(),
    fullAddress: z.string().optional(),
    latitude: z.number().nullable().optional(),
    longitude: z.number().nullable().optional(),
  }),
  settings: z.object({
    preferredLanguage: z.string(),
    preferredCurrency: z.string(),
    timezone: z.string(),
  }),
  operational: z.object({
    openingHours: z.string(),
    closingHours: z.string(),
    businessDays: z.array(z.string()),
    deliveryRadius: z.number(),
    pickupAvailable: z.boolean(),
  }),
  financial: z.object({
    bankAccountName: z.string().min(2, 'Account name is required'),
    bankName: z.string().min(2, 'Bank name is required'),
    bankAccountNumber: z.string().min(5, 'Account number is required'),
    eWalletInformation: z.string().optional(),
    taxNumber: z.string().optional(),
  }),
});

export default function MerchantProfileSettings() {
  const { profile, loading, error, saveProfile } = useMerchantProfileViewModel();

  const methods = useForm<MerchantProfile>({
    resolver: zodResolver(profileSchema),
    defaultValues: {
      operational: {
        businessDays: [],
      }
    }
  });

  const { handleSubmit, reset, formState: { isSubmitting, isDirty } } = methods;

  useEffect(() => {
    if (profile) {
      reset(profile);
    }
  }, [profile, reset]);

  const onSubmit = async (data: MerchantProfile) => {
    await saveProfile(data);
  };

  if (loading && !profile) {
    return (
      <div className="flex h-[400px] items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-indigo-600" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex h-[400px] flex-col items-center justify-center text-red-500">
        <AlertCircle className="h-8 w-8 mb-2" />
        <p>{error}</p>
      </div>
    );
  }

  return (
    <div className="space-y-6 pb-10">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold tracking-tight">Profile Settings</h2>
          <p className="text-muted-foreground mt-1">
            Manage your store information, operational hours, and financial details.
          </p>
        </div>
        <Button 
          onClick={handleSubmit(onSubmit)} 
          disabled={isSubmitting || !isDirty}
          className="bg-indigo-600 hover:bg-indigo-700"
        >
          {isSubmitting ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : <Save className="mr-2 h-4 w-4" />}
          Save Changes
        </Button>
      </div>

      <FormProvider {...methods}>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-8">
          <Tabs defaultValue="identity" className="w-full">
            <TabsList className="mb-4">
              <TabsTrigger value="identity">Identity</TabsTrigger>
              <TabsTrigger value="contact">Contact</TabsTrigger>
              <TabsTrigger value="settings">Settings</TabsTrigger>
              <TabsTrigger value="operational">Operational</TabsTrigger>
              <TabsTrigger value="financial">Financial</TabsTrigger>
            </TabsList>
            
            <TabsContent value="identity">
              <ProfileIdentitySection />
            </TabsContent>
            
            <TabsContent value="contact">
              <ProfileContactSection />
            </TabsContent>
            
            <TabsContent value="settings">
              <ProfileSettingsSection />
            </TabsContent>
            
            <TabsContent value="operational">
              <ProfileOperationalSection />
            </TabsContent>
            
            <TabsContent value="financial">
              <ProfileFinancialSection />
            </TabsContent>
          </Tabs>
        </form>
      </FormProvider>
    </div>
  );
}
