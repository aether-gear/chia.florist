import { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import {
  Save,
  AlertCircle,
  Loader2,
  CheckCircle2,
  Calendar,
  Clock,
  ShieldCheck,
  Crown,
  Copy,
  Check,
  IdCard,
  Mail,
  User as UserIcon,
  Sparkles,
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { useStaffProfileViewModel } from '@/viewmodels/useStaffProfileViewModel';
import { useAuthMeViewModel } from '@/viewmodels/useAuthMeViewModel';
import AvatarUpload from '@/components/AvatarUpload';
import { useToast } from '@/hooks/use-toast';
import { formatRelativeLastLogin } from '@/lib/timeUtils';

const profileSchema = z.object({
  name: z.string().min(2, 'Full Name is required (minimum 2 characters)'),
  phone: z.string().min(5, 'Phone number is required (minimum 5 digits)').optional().or(z.literal('')),
});

type ProfileFormValues = z.infer<typeof profileSchema>;

export default function StaffProfileSettings() {
  const { profile, loading, error, saveProfile } = useStaffProfileViewModel();
  const { data: authData, isAdmin } = useAuthMeViewModel();
  const { toast } = useToast();

  const [copiedKey, setCopiedKey] = useState<string | null>(null);
  const [successBanner, setSuccessBanner] = useState<string | null>(null);

  const storedEmail =
    localStorage.getItem('userEmail') || sessionStorage.getItem('userEmail') || '';

  const {
    register,
    handleSubmit,
    reset,
    formState: { isSubmitting, errors, isDirty },
  } = useForm<ProfileFormValues>({
    resolver: zodResolver(profileSchema),
    defaultValues: {
      name: '',
      phone: '',
    },
  });

  useEffect(() => {
    if (profile) {
      reset({
        name: profile.Name || '',
        phone: profile.Phone || '',
      });
    }
  }, [profile, reset]);

  const onSubmit = async (data: ProfileFormValues) => {
    setSuccessBanner(null);
    const ok = await saveProfile(data);
    if (ok) {
      setSuccessBanner('Profile details updated successfully.');
      toast({
        title: 'Profile Updated',
        description: 'Your profile changes have been saved successfully.',
      });
      setTimeout(() => setSuccessBanner(null), 3500);
    }
  };

  const handleCopy = (text: string, key: string) => {
    navigator.clipboard.writeText(text);
    setCopiedKey(key);
    setTimeout(() => setCopiedKey(null), 1500);
  };

  if (loading && !profile) {
    return (
      <div className="flex h-[400px] items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    );
  }

  if (error && !profile) {
    return (
      <div className="flex h-[400px] flex-col items-center justify-center text-rose-500">
        <AlertCircle className="h-8 w-8 mb-2" />
        <p className="font-semibold text-sm">{error}</p>
      </div>
    );
  }

  const lastActiveText = formatRelativeLastLogin(profile?.LastLoginAt, 'Current session');
  const isActiveNow = lastActiveText === 'Active now' || lastActiveText === 'Current session';

  return (
    <div className="flex-col md:flex">
      <div className="flex-1 max-w-5xl w-full mx-auto space-y-10 p-6 sm:p-8 lg:p-12 animate-in fade-in duration-300">
        {/* Header */}
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <h2 className="text-3xl font-bold font-display tracking-tight text-foreground">
              Staff Profile
            </h2>
            <p className="text-muted-foreground text-sm">
              Manage your display name, contact phone, and avatar
            </p>
          </div>
          <div className="flex items-center gap-2">
            <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold bg-emerald-500/10 text-emerald-700 dark:text-emerald-400 border border-emerald-500/20">
              <span className="h-2 w-2 rounded-full bg-emerald-500 animate-pulse" />
              Account Active
            </span>
            {isAdmin && (
              <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-semibold bg-amber-500/10 text-amber-700 dark:text-amber-400 border border-amber-500/20">
                <Crown className="h-3 w-3" />
                Administrator
              </span>
            )}
          </div>
        </div>

        {successBanner && (
          <div className="flex items-center gap-2 p-4 text-sm text-primary bg-primary/10 rounded-2xl border border-primary/20 animate-in fade-in">
            <CheckCircle2 className="h-4 w-4 shrink-0 text-primary" />
            <span>{successBanner}</span>
          </div>
        )}

        <div className="space-y-10">
          {/* Top Section: Personal Information & Account Details */}
          <div className="space-y-6">
            <div className="pb-4 border-b border-border/60">
              <h3 className="font-bold font-display tracking-tight text-lg text-foreground">
                Personal Information
              </h3>
              <p className="text-muted-foreground text-xs">
                Update your display name, contact number, and profile picture.
              </p>
            </div>

            {/* Profile Picture Upload & Monogram Banner */}
            <div className="flex flex-col sm:flex-row items-start sm:items-center gap-5 p-5 rounded-2xl bg-muted/20 border border-border/40">
              {profile && (
                <AvatarUpload
                  userId={profile.user_id}
                  currentAvatarUrl={profile.AvatarURL || ''}
                  displayName={profile.Name}
                  onUploadComplete={(url) => saveProfile({ avatar_url: url })}
                  onRemoveComplete={() => saveProfile({ avatar_url: '' })}
                />
              )}
              <div className="space-y-1">
                <h4 className="font-bold font-display text-base text-foreground">
                  {profile?.Name || 'Staff Member'}
                </h4>
                <div className="flex items-center gap-2 text-xs text-muted-foreground">
                  <span className="font-medium text-foreground/80">@{profile?.Username}</span>
                  <span>•</span>
                  <span>{storedEmail || 'No email associated'}</span>
                </div>
                <p className="text-[11px] text-muted-foreground/80">
                  Click the camera icon to upload a square JPEG or PNG avatar.
                </p>
              </div>
            </div>

            {/* Editable Form */}
            <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                {/* Username (Read Only) */}
                <div className="space-y-1.5">
                  <Label htmlFor="username" className="text-xs font-semibold text-foreground flex items-center gap-1.5">
                    <UserIcon className="h-3.5 w-3.5 text-muted-foreground" />
                    Username
                  </Label>
                  <Input
                    id="username"
                    value={profile?.Username || ''}
                    disabled
                    className="bg-muted/70 rounded-xl border-border text-muted-foreground font-mono text-sm cursor-not-allowed"
                  />
                  <p className="text-[11px] text-muted-foreground">Unique account sign-in handle.</p>
                </div>

                {/* Email (Read Only) */}
                <div className="space-y-1.5">
                  <Label htmlFor="email" className="text-xs font-semibold text-foreground flex items-center gap-1.5">
                    <Mail className="h-3.5 w-3.5 text-muted-foreground" />
                    Email Address
                  </Label>
                  <Input
                    id="email"
                    value={storedEmail}
                    disabled
                    className="bg-muted/70 rounded-xl border-border text-muted-foreground text-sm cursor-not-allowed"
                  />
                  <p className="text-[11px] text-muted-foreground">Primary contact email managed by administrator.</p>
                </div>
              </div>

              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                {/* Full Name */}
                <div className="space-y-1.5">
                  <Label htmlFor="name" className="text-xs font-semibold text-foreground">
                    Full Name <span className="text-rose-500">*</span>
                  </Label>
                  <Input
                    id="name"
                    {...register('name')}
                    placeholder="e.g. Jane Doe"
                    className="rounded-xl border-border bg-background text-sm"
                  />
                  {errors.name && <p className="text-xs text-rose-500">{errors.name.message}</p>}
                </div>

                {/* Phone Number */}
                <div className="space-y-1.5">
                  <Label htmlFor="phone" className="text-xs font-semibold text-foreground">
                    Phone Number
                  </Label>
                  <Input
                    id="phone"
                    {...register('phone')}
                    placeholder="e.g. +628123456789"
                    className="rounded-xl border-border bg-background text-sm"
                  />
                  {errors.phone && <p className="text-xs text-rose-500">{errors.phone.message}</p>}
                </div>
              </div>

              <div className="flex justify-end pt-2">
                <Button
                  type="submit"
                  disabled={!isDirty || isSubmitting}
                  className="bg-primary hover:bg-primary/90 text-primary-foreground rounded-xl text-xs font-semibold shadow-sm min-w-[130px]"
                >
                  {isSubmitting ? (
                    <>
                      <Loader2 className="mr-1.5 h-3.5 w-3.5 animate-spin" /> Saving...
                    </>
                  ) : (
                    <>
                      <Save className="mr-1.5 h-3.5 w-3.5" /> Save Changes
                    </>
                  )}
                </Button>
              </div>
            </form>
          </div>

          {/* Bottom Section: Account Activity & Access (Non-Technical) */}
          <div className="space-y-6 pt-4">
            <div className="pb-4 border-b border-border/60">
              <h3 className="font-bold font-display tracking-tight text-lg text-foreground flex items-center gap-2">
                <Sparkles className="h-4 w-4 text-primary" />
                Account Activity & Access
              </h3>
              <p className="text-muted-foreground text-xs">
                Overview of your account status, active history, and reference codes.
              </p>
            </div>

            <div className="space-y-6 rounded-2xl bg-muted/10 border border-border/40 p-6">
              {/* Access Roles & Permissions */}
              <div className="pb-5 border-b border-border/30">
                <div className="flex items-center gap-1.5 text-xs font-bold text-muted-foreground uppercase tracking-wider mb-2">
                  <ShieldCheck className="h-3.5 w-3.5 text-primary" />
                  Your Access Level
                </div>
                <div className="flex flex-wrap gap-2">
                  {authData?.roles && authData.roles.length > 0 ? (
                    authData.roles.map((r) => (
                      <span
                        key={r.code}
                        className="inline-flex items-center gap-1 px-3 py-1 rounded-xl text-xs font-semibold bg-primary/10 text-primary border border-primary/20"
                      >
                        <ShieldCheck className="h-3 w-3" />
                        {r.name || r.code}
                      </span>
                    ))
                  ) : (
                    <span className="inline-flex items-center px-3 py-1 rounded-xl text-xs font-medium bg-muted text-muted-foreground">
                      Staff Member
                    </span>
                  )}
                </div>
              </div>

              {/* Activity Timeline Bar */}
              <div className="grid grid-cols-1 sm:grid-cols-3 gap-6 pb-5 border-b border-border/30">
                {/* Member Since */}
                <div>
                  <div className="flex items-center gap-1.5 text-[11px] font-bold text-muted-foreground uppercase tracking-wider mb-1">
                    <Calendar className="h-3 w-3 text-muted-foreground" />
                    Member Since
                  </div>
                  <p className="text-xs font-semibold text-foreground">
                    {profile?.CreatedAt
                      ? new Date(profile.CreatedAt).toLocaleDateString('en-GB', {
                          day: 'numeric',
                          month: 'short',
                          year: 'numeric',
                        })
                      : 'N/A'}
                  </p>
                </div>

                {/* Last Active */}
                <div>
                  <div className="flex items-center gap-1.5 text-[11px] font-bold text-muted-foreground uppercase tracking-wider mb-1">
                    <Clock className="h-3 w-3 text-muted-foreground" />
                    Last Active
                  </div>
                  <p className={`text-xs font-semibold flex items-center gap-1.5 ${isActiveNow ? 'text-primary' : 'text-foreground'}`}>
                    {isActiveNow && <span className="h-1.5 w-1.5 rounded-full bg-primary animate-pulse" />}
                    {lastActiveText}
                  </p>
                </div>

                {/* Profile Last Updated */}
                <div>
                  <div className="flex items-center gap-1.5 text-[11px] font-bold text-muted-foreground uppercase tracking-wider mb-1">
                    <Clock className="h-3 w-3 text-muted-foreground" />
                    Profile Last Updated
                  </div>
                  <p className="text-xs font-medium text-muted-foreground">
                    {profile?.UpdatedAt
                      ? new Date(profile.UpdatedAt).toLocaleDateString('en-GB', {
                          day: 'numeric',
                          month: 'short',
                          year: 'numeric',
                        })
                      : 'No recent updates'}
                  </p>
                </div>
              </div>

              {/* Account Reference Codes */}
              <div className="space-y-4 pt-1">
                <div>
                  <h4 className="text-xs font-bold font-display uppercase tracking-wider text-muted-foreground">
                    Account Reference Codes
                  </h4>
                  <p className="text-[11px] text-muted-foreground/80 mt-0.5">
                    Unique references for your staff assignment and user account. Use these if you ever need support.
                  </p>
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  {/* Staff ID */}
                  <div className="p-3.5 rounded-xl bg-background border border-border/60">
                    <div className="flex items-center justify-between text-[11px] font-bold text-muted-foreground mb-1.5">
                      <span className="flex items-center gap-1">
                        <IdCard className="h-3 w-3 text-primary" />
                        Staff Assignment ID
                      </span>
                      <button
                        type="button"
                        onClick={() => profile?.staff_id && handleCopy(profile.staff_id, 'staffId')}
                        className="text-[10px] font-semibold text-primary hover:underline flex items-center gap-1"
                      >
                        {copiedKey === 'staffId' ? (
                          <>
                            <Check className="h-2.5 w-2.5" /> Copied
                          </>
                        ) : (
                          <>
                            <Copy className="h-2.5 w-2.5" /> Copy
                          </>
                        )}
                      </button>
                    </div>
                    <p className="text-xs font-mono text-muted-foreground truncate select-all">
                      {profile?.staff_id || 'N/A'}
                    </p>
                  </div>

                  {/* User ID */}
                  <div className="p-3.5 rounded-xl bg-background border border-border/60">
                    <div className="flex items-center justify-between text-[11px] font-bold text-muted-foreground mb-1.5">
                      <span className="flex items-center gap-1">
                        <UserIcon className="h-3 w-3 text-primary" />
                        User Account ID
                      </span>
                      <button
                        type="button"
                        onClick={() => profile?.user_id && handleCopy(profile.user_id, 'userId')}
                        className="text-[10px] font-semibold text-primary hover:underline flex items-center gap-1"
                      >
                        {copiedKey === 'userId' ? (
                          <>
                            <Check className="h-2.5 w-2.5" /> Copied
                          </>
                        ) : (
                          <>
                            <Copy className="h-2.5 w-2.5" /> Copy
                          </>
                        )}
                      </button>
                    </div>
                    <p className="text-xs font-mono text-muted-foreground truncate select-all">
                      {profile?.user_id || 'N/A'}
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
