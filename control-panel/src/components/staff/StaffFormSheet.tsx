import { useState, useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import {
  Loader2,
  Building,
  UserPlus,
  X,
  CheckCircle2,
  AlertCircle,
  Eye,
  EyeOff,
  Save,
  Link as LinkIcon,
  ShieldCheck,
} from 'lucide-react';
import { Button } from '../ui/button';
import { Input } from '../ui/input';
import { Label } from '../ui/label';
import { Textarea } from '../ui/textarea';
import { Switch } from '../ui/switch';
import {
  Sheet,
  SheetContent,
} from '../ui/sheet';
import StaffAccountsList from './StaffAccountsList';
import type { Staff, StaffAccountMember } from '@/models/Staff';

const staffEntitySchema = z.object({
  name: z.string().min(2, 'Staff entity name must be at least 2 characters'),
  username: z.string().min(3, 'Username must be at least 3 characters'),
  description: z.string().optional(),
  logo_url: z.string().url('Must be a valid URL').optional().or(z.literal('')),
  banner_url: z.string().url('Must be a valid URL').optional().or(z.literal('')),
});

const accountSchema = z.object({
  email: z.string().email('Please enter a valid email address'),
  password: z.string().min(6, 'Password must be at least 6 characters'),
});

type StaffEntityFormData = z.infer<typeof staffEntitySchema>;
type AccountFormData = z.infer<typeof accountSchema>;

interface StaffFormSheetProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  staffToEdit?: Staff | null;
  accounts?: StaffAccountMember[];
  onSuccess: () => void;
  mode?: 'create' | 'edit' | 'bind';
  // ViewModel action handlers
  createStaff?: (payload: any) => Promise<any>;
  createStaffWithAccount?: (staffPayload: any, accountPayload: any) => Promise<any>;
  updateStaff?: (staffId: string, payload: any) => Promise<any>;
  addStaffAccount?: (staffId: string, payload: any) => Promise<any>;
  removeStaffAccount?: (staffId: string, accountId: string) => Promise<any>;
}

export default function StaffFormSheet({
  open,
  onOpenChange,
  staffToEdit,
  accounts = [],
  onSuccess,
  mode = 'create',
  createStaff,
  createStaffWithAccount,
  updateStaff,
  addStaffAccount,
  removeStaffAccount,
}: StaffFormSheetProps) {
  const isEdit = mode === 'edit' && Boolean(staffToEdit);
  const isBindOnly = mode === 'bind' && Boolean(staffToEdit);

  // States
  const [bindDirectly, setBindDirectly] = useState(true);
  const [showBindInlineForm, setShowBindInlineForm] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [successMsg, setSuccessMsg] = useState<string | null>(null);
  const [showPassword, setShowPassword] = useState(false);

  // Form for Staff Entity
  const {
    register: registerStaff,
    handleSubmit: handleSubmitStaff,
    formState: { errors: errorsStaff, isDirty: isStaffDirty },
    reset: resetStaff,
    setValue: setStaffValue,
  } = useForm<StaffEntityFormData>({
    resolver: zodResolver(staffEntitySchema),
    defaultValues: {
      name: '',
      username: '',
      description: '',
      logo_url: '',
      banner_url: '',
    },
  });

  // Form for Bound Account
  const {
    register: registerAccount,
    handleSubmit: handleSubmitAccount,
    formState: { errors: errorsAccount },
    reset: resetAccount,
  } = useForm<AccountFormData>({
    resolver: zodResolver(accountSchema),
    defaultValues: {
      email: '',
      password: '',
    },
  });

  useEffect(() => {
    if (open) {
      setError(null);
      setSuccessMsg(null);
      setShowPassword(false);
      setShowBindInlineForm(isBindOnly || (isEdit && accounts.length === 0));

      if (isEdit && staffToEdit) {
        resetStaff({
          name: staffToEdit.name || '',
          username: staffToEdit.username || '',
          description: staffToEdit.description || '',
          logo_url: staffToEdit.logo_url || staffToEdit.avatar_url || '',
          banner_url: staffToEdit.banner_url || '',
        });
      } else {
        resetStaff({
          name: '',
          username: '',
          description: '',
          logo_url: '',
          banner_url: '',
        });
        setBindDirectly(true);
      }

      resetAccount({
        email: '',
        password: '',
      });
    }
  }, [open, isEdit, isBindOnly, staffToEdit, accounts, resetStaff, resetAccount]);

  // Handle Staff Entity name change -> auto-suggest username if empty
  const handleStaffNameBlur = (e: React.FocusEvent<HTMLInputElement>) => {
    if (!isEdit && e.target.value) {
      const slug = e.target.value
        .toLowerCase()
        .replace(/[^a-z0-9]/g, '-')
        .replace(/-+/g, '-')
        .replace(/^-|-$/g, '');
      setStaffValue('username', slug, { shouldValidate: true });
    }
  };

  // Submit Handler for Creation (Single or Combined)
  const handleCreateSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    setSuccessMsg(null);

    try {
      await handleSubmitStaff(async (staffData) => {
        const staffPayload = {
          name: staffData.name.trim(),
          username: staffData.username.trim(),
          description: staffData.description?.trim() || undefined,
          logo_url: staffData.logo_url?.trim() || undefined,
          banner_url: staffData.banner_url?.trim() || undefined,
        };

        if (bindDirectly && createStaffWithAccount) {
          // Validate and submit account form alongside staff
          await handleSubmitAccount(async (accountData) => {
            const accountPayload = {
              email: accountData.email.trim(),
              password: accountData.password,
            };

            await createStaffWithAccount(staffPayload, accountPayload);
            setSuccessMsg('Staff entity and bound account successfully created!');
            setTimeout(() => {
              onSuccess();
              onOpenChange(false);
            }, 800);
          })();
        } else if (createStaff) {
          await createStaff(staffPayload);
          setSuccessMsg('Staff entity created successfully!');
          setTimeout(() => {
            onSuccess();
            onOpenChange(false);
          }, 800);
        }
      })();
    } catch (err: any) {
      console.error('Failed to create staff', err);
      setError(err.message || 'Failed to create staff');
    } finally {
      setLoading(false);
    }
  };

  // Submit Handler for Edit Staff Entity
  const handleEditStaffSubmit = async (data: StaffEntityFormData) => {
    if (!staffToEdit || !updateStaff) return;
    setLoading(true);
    setError(null);
    setSuccessMsg(null);

    try {
      await updateStaff(staffToEdit.id, {
        name: data.name.trim(),
        description: data.description?.trim() || undefined,
        logo_url: data.logo_url?.trim() || undefined,
        banner_url: data.banner_url?.trim() || undefined,
      });

      setSuccessMsg('Staff entity details updated successfully!');
      setTimeout(() => {
        onSuccess();
        onOpenChange(false);
      }, 700);
    } catch (err: any) {
      console.error('Failed to update staff', err);
      setError(err.message || 'Failed to update staff');
    } finally {
      setLoading(false);
    }
  };

  // Submit Handler for Adding Account to Existing Staff in Edit Mode
  const handleAddAccountToExistingStaff = async (data: AccountFormData) => {
    if (!staffToEdit || !addStaffAccount) return;
    setLoading(true);
    setError(null);
    setSuccessMsg(null);

    try {
      await addStaffAccount(staffToEdit.id, {
        email: data.email.trim(),
        password: data.password,
      });

      setSuccessMsg('Staff account successfully bound!');
      setShowBindInlineForm(false);
      resetAccount({
        email: '',
        password: '',
      });
      onSuccess();
    } catch (err: any) {
      console.error('Failed to bind account', err);
      setError(err.message || 'Failed to bind account');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent className="w-full sm:max-w-none md:w-[48vw] md:min-w-[480px] p-0 flex flex-col h-full border-l border-border/60 bg-background shadow-2xl">
        {/* Header */}
        <div className="flex items-center justify-between px-6 py-5 border-b border-border/60 pr-12">
          <div>
            <h2 className="text-xl font-bold font-display tracking-tight text-foreground flex items-center gap-2">
              <Building className="h-5 w-5 text-primary" />
              {isEdit ? 'Edit Staff Entity' : isBindOnly ? 'Bind Account to Staff' : 'Create New Staff'}
            </h2>
            <p className="text-xs text-muted-foreground mt-0.5">
              {isEdit
                ? `Update metadata and manage bound accounts for ${staffToEdit?.name}`
                : isBindOnly
                ? `Attach user credentials to ${staffToEdit?.name}`
                : 'Introduce a new staff unit with direct or deferred account binding'}
            </p>
          </div>
        </div>

        {/* Scrollable Content Body */}
        <div className="flex-1 overflow-y-auto px-6 py-6 space-y-6">
          {error && (
            <div className="flex items-start gap-2.5 p-3.5 text-xs text-rose-600 bg-rose-500/10 rounded-xl border border-rose-500/20 animate-in fade-in">
              <AlertCircle className="h-4 w-4 shrink-0 mt-0.5" />
              <span>{error}</span>
            </div>
          )}

          {successMsg && (
            <div className="flex items-center gap-2 p-3.5 text-xs text-primary bg-primary/10 rounded-xl border border-primary/20 animate-in fade-in">
              <CheckCircle2 className="h-4 w-4 shrink-0 text-primary" />
              <span>{successMsg}</span>
            </div>
          )}

          {/* CREATE MODE FORM */}
          {!isEdit && !isBindOnly && (
            <form id="create-staff-form" onSubmit={handleCreateSubmit} className="space-y-6">
              {/* Section 1: Staff Entity */}
              <div className="space-y-4">
                <div className="pb-2 border-b border-border/40">
                  <h3 className="text-xs font-bold font-display uppercase tracking-wider text-muted-foreground">
                    Staff Unit Details
                  </h3>
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor="create-name" className="text-xs font-semibold text-foreground">
                      Staff Entity Name <span className="text-rose-500">*</span>
                    </Label>
                    <Input
                      id="create-name"
                      {...registerStaff('name', { onBlur: handleStaffNameBlur })}
                      placeholder="e.g. Floral Workshop Central"
                      className="rounded-xl border-border bg-background text-sm"
                    />
                    {errorsStaff.name && (
                      <p className="text-xs text-rose-500">{errorsStaff.name.message}</p>
                    )}
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="create-username" className="text-xs font-semibold text-foreground">
                      Username <span className="text-rose-500">*</span>
                    </Label>
                    <Input
                      id="create-username"
                      {...registerStaff('username')}
                      placeholder="e.g. floral-workshop"
                      className="rounded-xl border-border bg-background text-sm"
                    />
                    {errorsStaff.username && (
                      <p className="text-xs text-rose-500">{errorsStaff.username.message}</p>
                    )}
                  </div>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="create-description" className="text-xs font-semibold text-foreground">
                    Description <span className="text-muted-foreground font-normal">(Optional)</span>
                  </Label>
                  <Textarea
                    id="create-description"
                    {...registerStaff('description')}
                    placeholder="Primary fulfillment and inventory handling branch..."
                    rows={2}
                    className="rounded-xl border-border bg-background text-sm resize-none"
                  />
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor="create-logo" className="text-xs font-semibold text-foreground">
                      Logo URL <span className="text-muted-foreground font-normal">(Optional)</span>
                    </Label>
                    <Input
                      id="create-logo"
                      {...registerStaff('logo_url')}
                      placeholder="https://.../logo.png"
                      className="rounded-xl border-border bg-background text-sm font-mono text-xs"
                    />
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="create-banner" className="text-xs font-semibold text-foreground">
                      Banner URL <span className="text-muted-foreground font-normal">(Optional)</span>
                    </Label>
                    <Input
                      id="create-banner"
                      {...registerStaff('banner_url')}
                      placeholder="https://.../banner.png"
                      className="rounded-xl border-border bg-background text-sm font-mono text-xs"
                    />
                  </div>
                </div>
              </div>

              {/* Section 2: Immediate Account Binding Toggle */}
              <div className="space-y-4 pt-2">
                <div className="flex items-center justify-between p-4 rounded-2xl bg-muted/20 border border-border/50">
                  <div className="space-y-0.5 pr-4">
                    <Label htmlFor="bind-toggle" className="text-xs font-bold text-foreground cursor-pointer flex items-center gap-1.5">
                      <LinkIcon className="h-3.5 w-3.5 text-primary" />
                      Bind Account Directly
                    </Label>
                    <p className="text-[11px] text-muted-foreground">
                      Create login credentials for this staff unit immediately instead of later.
                    </p>
                  </div>
                  <Switch
                    id="bind-toggle"
                    checked={bindDirectly}
                    onCheckedChange={setBindDirectly}
                  />
                </div>

                {/* Conditional Account Fields */}
                {bindDirectly && (
                  <div className="p-4 rounded-2xl bg-muted/10 border border-border/40 space-y-4 animate-in fade-in duration-200">
                    <div className="pb-2 border-b border-border/30 flex items-center gap-1.5">
                      <ShieldCheck className="h-3.5 w-3.5 text-primary" />
                      <span className="text-xs font-bold font-display text-foreground">
                        Account Credentials
                      </span>
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="account-email" className="text-xs font-semibold text-foreground">
                        Email Address <span className="text-rose-500">*</span>
                      </Label>
                      <Input
                        id="account-email"
                        type="email"
                        {...registerAccount('email')}
                        placeholder="staff@chia.florist"
                        className="rounded-xl border-border bg-background text-sm"
                      />
                      {errorsAccount.email && (
                        <p className="text-xs text-rose-500">{errorsAccount.email.message}</p>
                      )}
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="account-password" className="text-xs font-semibold text-foreground">
                        Password <span className="text-rose-500">*</span>
                      </Label>
                      <div className="relative">
                        <Input
                          id="account-password"
                          type={showPassword ? 'text' : 'password'}
                          {...registerAccount('password')}
                          placeholder="••••••••"
                          className="rounded-xl border-border bg-background text-sm pr-9"
                        />
                        <button
                          type="button"
                          onClick={() => setShowPassword(!showPassword)}
                          className="absolute right-2.5 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground p-0.5"
                        >
                          {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                        </button>
                      </div>
                      {errorsAccount.password && (
                        <p className="text-xs text-rose-500">{errorsAccount.password.message}</p>
                      )}
                    </div>
                  </div>
                )}
              </div>
            </form>
          )}

          {/* EDIT MODE FORM */}
          {isEdit && staffToEdit && (
            <div className="space-y-8">
              {/* Staff Entity Metadata Form */}
              <form id="edit-staff-form" onSubmit={handleSubmitStaff(handleEditStaffSubmit)} className="space-y-4">
                <div className="pb-2 border-b border-border/40 flex items-center justify-between">
                  <h3 className="text-xs font-bold font-display uppercase tracking-wider text-muted-foreground">
                    Staff Unit Information
                  </h3>
                  <Button
                    type="submit"
                    size="sm"
                    disabled={loading || !isStaffDirty}
                    className="h-7 text-xs rounded-lg font-semibold bg-primary hover:bg-primary/90 text-primary-foreground"
                  >
                    <Save className="h-3 w-3 mr-1" /> Save Staff Details
                  </Button>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="edit-name" className="text-xs font-semibold text-foreground">
                    Staff Entity Name <span className="text-rose-500">*</span>
                  </Label>
                  <Input
                    id="edit-name"
                    {...registerStaff('name')}
                    className="rounded-xl border-border bg-background text-sm"
                  />
                  {errorsStaff.name && (
                    <p className="text-xs text-rose-500">{errorsStaff.name.message}</p>
                  )}
                </div>

                <div className="space-y-2">
                  <Label htmlFor="edit-description" className="text-xs font-semibold text-foreground">
                    Description
                  </Label>
                  <Textarea
                    id="edit-description"
                    {...registerStaff('description')}
                    rows={2}
                    className="rounded-xl border-border bg-background text-sm resize-none"
                  />
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor="edit-logo" className="text-xs font-semibold text-foreground">
                      Logo URL
                    </Label>
                    <Input
                      id="edit-logo"
                      {...registerStaff('logo_url')}
                      className="rounded-xl border-border bg-background text-sm font-mono text-xs"
                    />
                  </div>

                  <div className="space-y-2">
                    <Label htmlFor="edit-banner" className="text-xs font-semibold text-foreground">
                      Banner URL
                    </Label>
                    <Input
                      id="edit-banner"
                      {...registerStaff('banner_url')}
                      className="rounded-xl border-border bg-background text-sm font-mono text-xs"
                    />
                  </div>
                </div>
              </form>

              {/* Bound Accounts Section */}
              <div className="space-y-4 pt-2">
                <div className="pb-2 border-b border-border/40 flex items-center justify-between">
                  <div className="flex items-center gap-1.5">
                    <LinkIcon className="h-3.5 w-3.5 text-primary" />
                    <h3 className="text-xs font-bold font-display uppercase tracking-wider text-muted-foreground">
                      {accounts.length > 0 ? 'Bound Account' : 'Account Binding'}
                    </h3>
                  </div>
                  {!showBindInlineForm && accounts.length === 0 && (
                    <Button
                      type="button"
                      variant="outline"
                      size="sm"
                      onClick={() => setShowBindInlineForm(true)}
                      className="h-7 text-xs rounded-lg font-semibold border-border hover:bg-muted"
                    >
                      <UserPlus className="h-3 w-3 mr-1" /> Bind Account
                    </Button>
                  )}
                </div>

                {accounts.length > 0 && (
                  <p className="text-[11px] text-muted-foreground">
                    1 account per staff unit. To attach a different account, unbind the current credentials first.
                  </p>
                )}

                {/* Bound Accounts List */}
                <StaffAccountsList
                  staffId={staffToEdit.id}
                  accounts={accounts}
                  onUnbind={
                    removeStaffAccount
                      ? async (accountId) => {
                          await removeStaffAccount(staffToEdit.id, accountId);
                          onSuccess();
                        }
                      : undefined
                  }
                  onBindNew={() => setShowBindInlineForm(true)}
                />

                {/* Inline Bind Account Form */}
                {showBindInlineForm && (
                  <form
                    onSubmit={handleSubmitAccount(handleAddAccountToExistingStaff)}
                    className="p-4 rounded-2xl bg-muted/20 border border-primary/20 space-y-4 animate-in fade-in"
                  >
                    <div className="flex items-center justify-between pb-2 border-b border-border/30">
                      <div className="flex items-center gap-1.5">
                        <UserPlus className="h-3.5 w-3.5 text-primary" />
                        <span className="text-xs font-bold font-display text-foreground">
                          Bind New Account to Staff
                        </span>
                      </div>
                      <Button
                        type="button"
                        variant="ghost"
                        size="icon"
                        onClick={() => setShowBindInlineForm(false)}
                        className="h-6 w-6 rounded-md hover:bg-muted"
                      >
                        <X className="h-3 w-3" />
                      </Button>
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="inline-email" className="text-xs font-semibold text-foreground">
                        Email Address <span className="text-rose-500">*</span>
                      </Label>
                      <Input
                        id="inline-email"
                        type="email"
                        {...registerAccount('email')}
                        placeholder="newuser@chia.florist"
                        className="rounded-xl border-border bg-background text-sm"
                      />
                      {errorsAccount.email && (
                        <p className="text-xs text-rose-500">{errorsAccount.email.message}</p>
                      )}
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="inline-password" className="text-xs font-semibold text-foreground">
                        Password <span className="text-rose-500">*</span>
                      </Label>
                      <div className="relative">
                        <Input
                          id="inline-password"
                          type={showPassword ? 'text' : 'password'}
                          {...registerAccount('password')}
                          placeholder="••••••••"
                          className="rounded-xl border-border bg-background text-sm pr-9"
                        />
                        <button
                          type="button"
                          onClick={() => setShowPassword(!showPassword)}
                          className="absolute right-2.5 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground p-0.5"
                        >
                          {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                        </button>
                      </div>
                      {errorsAccount.password && (
                        <p className="text-xs text-rose-500">{errorsAccount.password.message}</p>
                      )}
                    </div>

                    <div className="flex justify-end gap-2 pt-2">
                      <Button
                        type="button"
                        variant="ghost"
                        size="sm"
                        onClick={() => setShowBindInlineForm(false)}
                        className="text-xs rounded-xl"
                      >
                        Cancel
                      </Button>
                      <Button
                        type="submit"
                        size="sm"
                        disabled={loading}
                        className="text-xs rounded-xl font-semibold bg-primary hover:bg-primary/90 text-primary-foreground"
                      >
                        {loading ? <Loader2 className="h-3 w-3 animate-spin mr-1.5" /> : <UserPlus className="h-3 w-3 mr-1.5" />}
                        Bind Account
                      </Button>
                    </div>
                  </form>
                )}
              </div>
            </div>
          )}
        </div>

        {/* Footer (For Create Mode) */}
        {!isEdit && !isBindOnly && (
          <div className="flex items-center justify-end gap-2 px-6 py-4 border-t border-border/60 bg-muted/20">
            <Button
              type="button"
              variant="outline"
              size="sm"
              onClick={() => onOpenChange(false)}
              disabled={loading}
              className="text-xs rounded-xl"
            >
              Cancel
            </Button>
            <Button
              type="submit"
              form="create-staff-form"
              size="sm"
              disabled={loading}
              className="text-xs rounded-xl font-semibold bg-primary hover:bg-primary/90 text-primary-foreground min-w-[130px]"
            >
              {loading ? (
                <>
                  <Loader2 className="mr-2 h-3.5 w-3.5 animate-spin" />
                  Creating...
                </>
              ) : (
                <>
                  <Building className="mr-1.5 h-3.5 w-3.5" />
                  Create Staff
                </>
              )}
            </Button>
          </div>
        )}
      </SheetContent>
    </Sheet>
  );
}
