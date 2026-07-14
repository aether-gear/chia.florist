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
    <div className="flex h-[400px] items-center justify-center">
      <Loader2 className="h-8 w-8 animate-spin text-primary" />
    </div>
  }

  if (error) {
    return (
      <div className="flex h-[400px] flex-col items-center justify-center text-rose-500">
        <AlertCircle className="h-8 w-8 mb-2" />
        <p>{error}</p>
      </div>
    );
  }

  return (
    <div className="space-y-10 animate-in fade-in duration-300">
      <div>
        <h1 className="text-3xl font-bold font-display tracking-tight text-foreground">Profile Settings</h1>
        <p className="text-muted-foreground text-sm">Manage your staff account information.</p>
      </div>

      <div className="grid gap-6 md:grid-cols-2">
        <Card className="border-0 shadow-none bg-zinc-50/40 dark:bg-slate-900/40">
          <CardHeader>
            <CardTitle className="font-bold font-display tracking-tight text-lg text-foreground">Account Details</CardTitle>
            <CardDescription className="text-muted-foreground text-sm">
              Basic information about your staff account. Username cannot be changed.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
              
              <div className="flex items-center space-x-4 mb-6">
                <div className="h-20 w-20 rounded-full bg-primary/10 overflow-hidden border border-border">
                  {profile?.AvatarURL ? (
                    <img src={profile.AvatarURL} alt="Profile" className="h-full w-full object-cover" />
                  ) : (
                    <div className="h-full w-full flex items-center justify-center text-primary font-bold text-xl">
                      {profile?.Name?.charAt(0) || 'S'}
                    </div>
                  )}
                </div>
                <div>
                  <h3 className="font-bold font-display text-foreground">{profile?.Name}</h3>
                  <p className="text-sm text-muted-foreground">@{profile?.Username}</p>
                </div>
              </div>

              <div className="space-y-2">
                <Label htmlFor="username">Username (Read-Only)</Label>
                <Input id="username" value={profile?.Username || ''} disabled className="bg-muted rounded-xl border-border" />
              </div>

              <div className="space-y-2">
                <Label htmlFor="name">Full Name</Label>
                <Input id="name" {...register('name')} placeholder="e.g. Jane Doe" className="rounded-xl border-border bg-background" />
                {errors.name && <p className="text-sm text-rose-500">{errors.name.message}</p>}
              </div>

              <div className="space-y-2">
                <Label htmlFor="phone">Phone Number</Label>
                <Input id="phone" {...register('phone')} placeholder="e.g. 08123456789" className="rounded-xl border-border bg-background" />
                {errors.phone && <p className="text-sm text-rose-500">{errors.phone.message}</p>}
              </div>

              <div className="space-y-2">
                <Label htmlFor="avatar_url">Avatar URL</Label>
                <Input id="avatar_url" {...register('avatar_url')} placeholder="https://example.com/avatar.jpg" className="rounded-xl border-border bg-background" />
                {errors.avatar_url && <p className="text-sm text-rose-500">{errors.avatar_url.message}</p>}
              </div>

              <div className="flex justify-end">
                <Button type="submit" disabled={!isDirty || isSubmitting} className="bg-primary hover:bg-primary/90 text-primary-foreground rounded-xl">
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

        <Card className="border-0 shadow-none bg-zinc-50/40 dark:bg-slate-900/40">
          <CardHeader>
            <CardTitle className="font-bold font-display tracking-tight text-lg text-foreground">Session Information</CardTitle>
            <CardDescription className="text-muted-foreground text-sm">
              Details regarding your current session and account status.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div>
              <p className="text-xs font-bold text-muted-foreground uppercase tracking-wider">Staff ID</p>
              <p className="text-sm text-foreground font-semibold mt-1">{profile?.staff_id}</p>
            </div>
            <div>
              <p className="text-xs font-bold text-muted-foreground uppercase tracking-wider">User ID</p>
              <p className="text-sm text-foreground font-semibold mt-1">{profile?.user_id}</p>
            </div>
            <div>
              <p className="text-xs font-bold text-muted-foreground uppercase tracking-wider">Account Created</p>
              <p className="text-sm text-foreground font-semibold mt-1">
                {profile?.CreatedAt ? new Date(profile.CreatedAt).toLocaleString() : 'N/A'}
              </p>
            </div>
            <div>
              <p className="text-xs font-bold text-muted-foreground uppercase tracking-wider">Last Login</p>
              <p className="text-sm text-foreground font-semibold mt-1">
                {profile?.LastLoginAt ? new Date(profile.LastLoginAt).toLocaleString() : 'Current Session'}
              </p>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
