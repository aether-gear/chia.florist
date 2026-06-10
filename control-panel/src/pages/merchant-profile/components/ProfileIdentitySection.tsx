import { useFormContext } from 'react-hook-form';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';

import type { MerchantProfile } from '@/models/MerchantProfile';

export function ProfileIdentitySection() {
  const { register, formState: { errors } } = useFormContext<MerchantProfile>();

  return (
    <Card className="border-0 shadow-sm">
      <CardHeader>
        <CardTitle>Identity Information</CardTitle>
        <CardDescription>Basic information about your merchant account.</CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="space-y-2">
            <label className="text-sm font-medium">Account ID</label>
            <Input {...register('identity.accountId')} disabled className="bg-slate-50" />
          </div>
          <div className="space-y-2">
            <label className="text-sm font-medium">Merchant ID</label>
            <Input {...register('identity.merchantId')} disabled className="bg-slate-50" />
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="space-y-2">
            <label className="text-sm font-medium">Merchant Name / Display Name</label>
            <Input {...register('identity.name')} placeholder="Your shop name" />
            {errors.identity?.name && (
              <p className="text-sm text-red-500">{errors.identity.name.message as string}</p>
            )}
          </div>
          <div className="space-y-2">
            <label className="text-sm font-medium">Merchant Slug</label>
            <Input {...register('identity.slug')} placeholder="e.g. my-shop-name" />
          </div>
        </div>

        <div className="space-y-2">
          <label className="text-sm font-medium">Description / Bio</label>
          <textarea 
            {...register('identity.description')} 
            className="flex min-h-[80px] w-full rounded-md border border-slate-200 bg-white px-3 py-2 text-sm ring-offset-white placeholder:text-slate-500 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-indigo-500 disabled:cursor-not-allowed disabled:opacity-50"
            placeholder="Tell customers about your business..."
          />
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="space-y-2">
            <label className="text-sm font-medium">Profile Photo URL</label>
            <Input {...register('identity.profilePhoto')} placeholder="https://..." />
          </div>
          <div className="space-y-2">
            <label className="text-sm font-medium">Cover Banner URL</label>
            <Input {...register('identity.coverBanner')} placeholder="https://..." />
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
