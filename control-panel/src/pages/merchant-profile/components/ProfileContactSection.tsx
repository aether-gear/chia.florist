import { useFormContext } from 'react-hook-form';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';

import type { MerchantProfile } from '@/models/MerchantProfile';

export function ProfileContactSection() {
  const { register, formState: { errors } } = useFormContext<MerchantProfile>();

  return (
    <Card className="border-0 shadow-sm">
      <CardHeader>
        <CardTitle>Contact Information</CardTitle>
        <CardDescription>How customers and administrators can reach you.</CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="space-y-2">
            <label className="text-sm font-medium">Email Address</label>
            <Input type="email" {...register('contact.email')} placeholder="store@example.com" />
            {errors.contact?.email && (
              <p className="text-sm text-red-500">{errors.contact.email.message as string}</p>
            )}
          </div>
          <div className="space-y-2">
            <label className="text-sm font-medium">Phone Number</label>
            <Input {...register('contact.phone')} placeholder="+62..." />
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="space-y-2">
            <label className="text-sm font-medium">WhatsApp Number</label>
            <Input {...register('contact.whatsapp')} placeholder="+62..." />
          </div>
          <div className="space-y-2">
            <label className="text-sm font-medium">Customer Service Contact</label>
            <Input {...register('contact.customerServiceContact')} placeholder="CS Email or Phone" />
          </div>
        </div>

        <div className="pt-4 pb-2">
          <h4 className="text-sm font-semibold border-b pb-2">Location Details</h4>
        </div>

        <div className="space-y-2">
          <label className="text-sm font-medium">Street Address</label>
          <Input {...register('contact.address')} placeholder="123 Main St" />
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div className="space-y-2">
            <label className="text-sm font-medium">Country</label>
            <Input {...register('contact.country')} placeholder="Indonesia" />
          </div>
          <div className="space-y-2">
            <label className="text-sm font-medium">Province / State</label>
            <Input {...register('contact.province')} placeholder="Jakarta" />
          </div>
          <div className="space-y-2">
            <label className="text-sm font-medium">City</label>
            <Input {...register('contact.city')} placeholder="South Jakarta" />
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="space-y-2">
            <label className="text-sm font-medium">District</label>
            <Input {...register('contact.district')} placeholder="Kebayoran Baru" />
          </div>
          <div className="space-y-2">
            <label className="text-sm font-medium">Postal Code</label>
            <Input {...register('contact.postalCode')} placeholder="12345" />
          </div>
        </div>

        <div className="space-y-2">
          <label className="text-sm font-medium">Full Formatted Address</label>
          <textarea 
            {...register('contact.fullAddress')} 
            className="flex min-h-[60px] w-full rounded-md border border-slate-200 bg-white px-3 py-2 text-sm ring-offset-white placeholder:text-slate-500 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-indigo-500 disabled:cursor-not-allowed disabled:opacity-50"
            placeholder="Complete address for display"
          />
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="space-y-2">
            <label className="text-sm font-medium">Latitude</label>
            <Input type="number" step="any" {...register('contact.latitude', { valueAsNumber: true })} placeholder="-6.200000" />
          </div>
          <div className="space-y-2">
            <label className="text-sm font-medium">Longitude</label>
            <Input type="number" step="any" {...register('contact.longitude', { valueAsNumber: true })} placeholder="106.816666" />
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
