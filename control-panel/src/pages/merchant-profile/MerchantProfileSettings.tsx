import { useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Save, AlertCircle, Loader2 } from 'lucide-react';

import { useMerchantProfileViewModel } from '@/viewmodels/useMerchantProfileViewModel';

const profileSchema = z.object({
  name: z.string().min(2, 'Name is required (min 2 chars)'),
  phone: z.string().min(5, 'Phone is required').optional().or(z.literal('')),
  avatar_url: z.string().url('Must be a valid URL').optional().or(z.literal('')),
});

type ProfileFormValues = z.infer<typeof profileSchema>;

export default function MerchantProfileSettings() {
  const { profile, loading, error, saveProfile } = useMerchantProfileViewModel();

  const { register, handleSubmit, reset, formState: { isSubmitting, errors, isDirty } } = useForm<ProfileFormValues>({
    resolver: zodResolver(profileSchema),
    defaultValues: {
      name: '',
      phone: '',
      avatar_url: ''
    }
  });

  useEffect(() => {
    if (profile) {
      reset({
        name: profile.Name || '',
        phone: profile.Phone || '',
        avatar_url: profile.AvatarURL || ''
      });
    }
  }, [profile, reset]);

  const onSubmit = async (data: ProfileFormValues) => {
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
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold tracking-tight text-slate-900 dark:text-slate-100">Profile Settings</h1>
        <p className="text-muted-foreground">Manage your staff account information.</p>
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Account Details</CardTitle>
            <CardDescription>
              Basic information about your staff account. Username cannot be changed.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
              
              <div className="flex items-center space-x-4 mb-6">
                <div className="h-20 w-20 rounded-full bg-indigo-100 overflow-hidden border border-gray-200">
                  {profile?.AvatarURL ? (
                    <img src={profile.AvatarURL} alt="Profile" className="h-full w-full object-cover" />
                  ) : (
                    <div className="h-full w-full flex items-center justify-center text-indigo-400 font-bold text-xl">
                      {profile?.Name?.charAt(0) || 'S'}
                    </div>
                  )}
                </div>
                <div>
                  <h3 className="font-medium text-slate-900 dark:text-slate-100">{profile?.Name}</h3>
                  <p className="text-sm text-gray-500">@{profile?.Username}</p>
                </div>
              </div>

              <div className="space-y-2">
                <Label htmlFor="username">Username (Read-Only)</Label>
                <Input id="username" value={profile?.Username || ''} disabled className="bg-gray-50 dark:bg-slate-900/50 dark:border-slate-800" />
              </div>

              <div className="space-y-2">
                <Label htmlFor="name">Full Name</Label>
                <Input id="name" {...register('name')} placeholder="e.g. Jane Doe" />
                {errors.name && <p className="text-sm text-red-500">{errors.name.message}</p>}
              </div>

              <div className="space-y-2">
                <Label htmlFor="phone">Phone Number</Label>
                <Input id="phone" {...register('phone')} placeholder="e.g. 08123456789" />
                {errors.phone && <p className="text-sm text-red-500">{errors.phone.message}</p>}
              </div>

              <div className="space-y-2">
                <Label htmlFor="avatar_url">Avatar URL</Label>
                <Input id="avatar_url" {...register('avatar_url')} placeholder="https://example.com/avatar.jpg" />
                {errors.avatar_url && <p className="text-sm text-red-500">{errors.avatar_url.message}</p>}
              </div>

              <div className="flex justify-end">
                <Button type="submit" disabled={!isDirty || isSubmitting}>
                  {isSubmitting ? (
                    <>
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" /> Saving...
                    </>
                  ) : (
                    <>
                      <Save className="mr-2 h-4 w-4" /> Save Profile
                    </>
                  )}
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Session Information</CardTitle>
            <CardDescription>
              Details regarding your current session and account status.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div>
              <p className="text-sm font-medium text-gray-500">Staff ID</p>
              <p className="text-sm text-slate-900 dark:text-slate-200">{profile?.staff_id}</p>
            </div>
            <div>
              <p className="text-sm font-medium text-gray-500">User ID</p>
              <p className="text-sm text-slate-900 dark:text-slate-200">{profile?.user_id}</p>
            </div>
            <div>
              <p className="text-sm font-medium text-gray-500">Account Created</p>
              <p className="text-sm text-slate-900 dark:text-slate-200">
                {profile?.CreatedAt ? new Date(profile.CreatedAt).toLocaleString() : 'N/A'}
              </p>
            </div>
            <div>
              <p className="text-sm font-medium text-gray-500">Last Login</p>
              <p className="text-sm text-slate-900 dark:text-slate-200">
                {profile?.LastLoginAt ? new Date(profile.LastLoginAt).toLocaleString() : 'Current Session'}
              </p>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
