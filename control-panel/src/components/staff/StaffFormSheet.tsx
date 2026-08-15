import { useState, useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { Loader2, UserPlus, Building, X, CheckCircle2, AlertCircle, Eye, EyeOff } from 'lucide-react';
import { Button } from '../ui/button';
import { Input } from '../ui/input';
import { Label } from '../ui/label';
import { Textarea } from '../ui/textarea';
import {
  Sheet,
  SheetContent,
} from '../ui/sheet';
import { fetchApi } from '@/lib/api';
import type { Staff } from '@/models/Staff';

const assignAccountSchema = z.object({
  staffId: z.string().min(1, 'Staff ID is required'),
  email: z.string().email('Please enter a valid email address'),
  name: z.string().min(2, 'Name must be at least 2 characters'),
  username: z.string().min(3, 'Username must be at least 3 characters'),
  password: z.string().min(6, 'Password must be at least 6 characters'),
  phone: z.string().min(5, 'Phone number is required (min 5 digits)'),
});

const createStaffSchema = z.object({
  name: z.string().min(2, 'Staff entity name must be at least 2 characters'),
  description: z.string().optional(),
  logo_url: z.string().url('Must be a valid URL').optional().or(z.literal('')),
  banner_url: z.string().url('Must be a valid URL').optional().or(z.literal('')),
});

type AssignAccountFormData = z.infer<typeof assignAccountSchema>;
type CreateStaffFormData = z.infer<typeof createStaffSchema>;

interface StaffFormSheetProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  selectedStaffId?: string;
  selectedStaffName?: string;
  staffList?: Staff[];
  onSuccess: () => void;
  defaultMode?: 'assign' | 'create';
}

export default function StaffFormSheet({
  open,
  onOpenChange,
  selectedStaffId,
  selectedStaffName,
  staffList = [],
  onSuccess,
  defaultMode = 'assign',
}: StaffFormSheetProps) {
  const [activeTab, setActiveTab] = useState<'assign' | 'create'>(defaultMode);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [successMsg, setSuccessMsg] = useState<string | null>(null);
  const [showPassword, setShowPassword] = useState(false);

  const {
    register: registerAssign,
    handleSubmit: handleSubmitAssign,
    formState: { errors: errorsAssign },
    reset: resetAssign,
  } = useForm<AssignAccountFormData>({
    resolver: zodResolver(assignAccountSchema),
    defaultValues: {
      staffId: selectedStaffId || '',
      email: '',
      name: '',
      username: '',
      password: '',
      phone: '',
    },
  });

  const {
    register: registerCreate,
    handleSubmit: handleSubmitCreate,
    formState: { errors: errorsCreate },
    reset: resetCreate,
  } = useForm<CreateStaffFormData>({
    resolver: zodResolver(createStaffSchema),
    defaultValues: {
      name: '',
      description: '',
      logo_url: '',
      banner_url: '',
    },
  });

  useEffect(() => {
    if (open) {
      setActiveTab(defaultMode);
      setError(null);
      setSuccessMsg(null);
      setShowPassword(false);
      resetAssign({
        staffId: selectedStaffId || (staffList.length > 0 ? staffList[0].id : ''),
        email: '',
        name: '',
        username: '',
        password: '',
        phone: '',
      });
      resetCreate({
        name: '',
        description: '',
        logo_url: '',
        banner_url: '',
      });
    }
  }, [open, selectedStaffId, staffList, defaultMode, resetAssign, resetCreate]);

  const onAssignSubmit = async (data: AssignAccountFormData) => {
    setLoading(true);
    setError(null);
    setSuccessMsg(null);

    try {
      const res = await fetchApi(`/staff/${data.staffId}/accounts`, {
        method: 'POST',
        body: JSON.stringify({
          email: data.email.trim(),
          name: data.name.trim(),
          username: data.username.trim(),
          password: data.password,
          phone: data.phone.trim(),
        }),
      });

      setSuccessMsg(res?.message || 'Staff account assigned successfully!');
      setTimeout(() => {
        onSuccess();
        onOpenChange(false);
      }, 900);
    } catch (err: any) {
      console.error('Failed to assign staff account', err);
      setError(err.message || 'Failed to assign staff account');
    } finally {
      setLoading(false);
    }
  };

  const onCreateSubmit = async (data: CreateStaffFormData) => {
    setLoading(true);
    setError(null);
    setSuccessMsg(null);

    try {
      const payload: Record<string, any> = {
        name: data.name.trim(),
      };
      if (data.description) payload.description = data.description.trim();
      if (data.logo_url) payload.logo_url = data.logo_url.trim();
      if (data.banner_url) payload.banner_url = data.banner_url.trim();

      const res = await fetchApi('/staff', {
        method: 'POST',
        body: JSON.stringify(payload),
      });

      setSuccessMsg(res?.message || 'Staff entity created successfully!');
      setTimeout(() => {
        onSuccess();
        onOpenChange(false);
      }, 900);
    } catch (err: any) {
      console.error('Failed to create staff entity', err);
      setError(err.message || 'Failed to create staff entity');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent className="w-full sm:max-w-none md:w-[48vw] md:min-w-[460px] p-0 flex flex-col h-full border-l border-border/60 bg-background shadow-2xl">
        {/* Header */}
        <div className="flex items-center justify-between px-6 py-5 border-b border-border/60">
          <div>
            <h2 className="text-xl font-bold font-display tracking-tight text-foreground">
              {activeTab === 'assign' ? 'Assign Staff Account' : 'Create Staff Entity'}
            </h2>
            <p className="text-xs text-muted-foreground mt-0.5">
              {activeTab === 'assign'
                ? 'Register and link an email account to a staff entity'
                : 'Create a new staff organizational unit'}
            </p>
          </div>
          <Button
            variant="ghost"
            size="icon"
            onClick={() => onOpenChange(false)}
            className="rounded-lg h-8 w-8 hover:bg-muted text-muted-foreground hover:text-foreground"
          >
            <X className="h-4 w-4" />
          </Button>
        </div>

        {/* Tab Switcher */}
        <div className="px-6 pt-4 pb-2 border-b border-border/40 bg-muted/20">
          <div className="flex rounded-xl bg-muted/60 p-1 gap-1">
            <button
              type="button"
              onClick={() => {
                setActiveTab('assign');
                setError(null);
                setSuccessMsg(null);
              }}
              className={`flex-1 flex items-center justify-center gap-2 py-2 text-xs font-semibold rounded-lg transition-all ${
                activeTab === 'assign'
                  ? 'bg-background text-foreground shadow-sm'
                  : 'text-muted-foreground hover:text-foreground'
              }`}
            >
              <UserPlus className="h-3.5 w-3.5 text-primary" />
              Assign Account
            </button>
            <button
              type="button"
              onClick={() => {
                setActiveTab('create');
                setError(null);
                setSuccessMsg(null);
              }}
              className={`flex-1 flex items-center justify-center gap-2 py-2 text-xs font-semibold rounded-lg transition-all ${
                activeTab === 'create'
                  ? 'bg-background text-foreground shadow-sm'
                  : 'text-muted-foreground hover:text-foreground'
              }`}
            >
              <Building className="h-3.5 w-3.5 text-primary" />
              New Staff Unit
            </button>
          </div>
        </div>

        {/* Scrollable Content Body */}
        <div className="flex-1 overflow-y-auto px-6 py-6 space-y-6">
          {error && (
            <div className="flex items-start gap-2.5 p-3.5 text-xs text-rose-600 bg-rose-500/10 rounded-xl border border-rose-500/20">
              <AlertCircle className="h-4 w-4 shrink-0 mt-0.5" />
              <span>{error}</span>
            </div>
          )}

          {successMsg && (
            <div className="flex items-center gap-2 p-3.5 text-xs text-primary bg-primary/10 rounded-xl border border-primary/20">
              <CheckCircle2 className="h-4 w-4 shrink-0 text-primary" />
              <span>{successMsg}</span>
            </div>
          )}

          {activeTab === 'assign' ? (
            <form id="assign-staff-form" onSubmit={handleSubmitAssign(onAssignSubmit)} className="space-y-4">
              {/* Target Staff Entity */}
              <div className="space-y-2">
                <Label htmlFor="staffId" className="text-xs font-semibold text-foreground">
                  Target Staff Entity <span className="text-rose-500">*</span>
                </Label>
                {staffList.length > 0 ? (
                  <div className="space-y-1.5">
                    <select
                      id="staffId"
                      {...registerAssign('staffId')}
                      className="w-full h-10 px-3 rounded-xl border border-border bg-background text-foreground text-sm focus:outline-none focus:ring-2 focus:ring-ring"
                    >
                      <option value="" disabled>Select a staff entity...</option>
                      {staffList.map((s) => (
                        <option key={s.id} value={s.id}>
                          {s.name} ({s.id.slice(0, 8)}...)
                        </option>
                      ))}
                    </select>
                    {selectedStaffName && (
                      <p className="text-[11px] text-muted-foreground">
                        Selected: <span className="font-semibold text-foreground">{selectedStaffName}</span>
                      </p>
                    )}
                  </div>
                ) : (
                  <Input
                    id="staffId"
                    {...registerAssign('staffId')}
                    placeholder="Enter Staff UUID"
                    className="rounded-xl border-border bg-background text-sm font-mono"
                  />
                )}
                {errorsAssign.staffId && (
                  <p className="text-xs text-rose-500">{errorsAssign.staffId.message}</p>
                )}
              </div>

              {/* Email Address */}
              <div className="space-y-2">
                <Label htmlFor="email" className="text-xs font-semibold text-foreground">
                  Email Address <span className="text-rose-500">*</span>
                </Label>
                <Input
                  id="email"
                  type="email"
                  {...registerAssign('email')}
                  placeholder="staff@chia.florist"
                  className="rounded-xl border-border bg-background text-sm"
                />
                {errorsAssign.email && (
                  <p className="text-xs text-rose-500">{errorsAssign.email.message}</p>
                )}
              </div>

              {/* Full Name & Username */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="name" className="text-xs font-semibold text-foreground">
                    Full Name <span className="text-rose-500">*</span>
                  </Label>
                  <Input
                    id="name"
                    {...registerAssign('name')}
                    placeholder="Jane Doe"
                    className="rounded-xl border-border bg-background text-sm"
                  />
                  {errorsAssign.name && (
                    <p className="text-xs text-rose-500">{errorsAssign.name.message}</p>
                  )}
                </div>

                <div className="space-y-2">
                  <Label htmlFor="username" className="text-xs font-semibold text-foreground">
                    Username <span className="text-rose-500">*</span>
                  </Label>
                  <Input
                    id="username"
                    {...registerAssign('username')}
                    placeholder="janedoe"
                    className="rounded-xl border-border bg-background text-sm"
                  />
                  {errorsAssign.username && (
                    <p className="text-xs text-rose-500">{errorsAssign.username.message}</p>
                  )}
                </div>
              </div>

              {/* Phone & Password */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="phone" className="text-xs font-semibold text-foreground">
                    Phone Number <span className="text-rose-500">*</span>
                  </Label>
                  <Input
                    id="phone"
                    {...registerAssign('phone')}
                    placeholder="+628123456789"
                    className="rounded-xl border-border bg-background text-sm"
                  />
                  {errorsAssign.phone && (
                    <p className="text-xs text-rose-500">{errorsAssign.phone.message}</p>
                  )}
                </div>

                <div className="space-y-2">
                  <Label htmlFor="password" className="text-xs font-semibold text-foreground">
                    Temporary Password <span className="text-rose-500">*</span>
                  </Label>
                  <div className="relative">
                    <Input
                      id="password"
                      type={showPassword ? 'text' : 'password'}
                      {...registerAssign('password')}
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
                  {errorsAssign.password && (
                    <p className="text-xs text-rose-500">{errorsAssign.password.message}</p>
                  )}
                </div>
              </div>
            </form>
          ) : (
            <form id="create-staff-form" onSubmit={handleSubmitCreate(onCreateSubmit)} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="staffName" className="text-xs font-semibold text-foreground">
                  Staff Entity Name <span className="text-rose-500">*</span>
                </Label>
                <Input
                  id="staffName"
                  {...registerCreate('name')}
                  placeholder="e.g. Floral Operations & Logistics"
                  className="rounded-xl border-border bg-background text-sm"
                />
                {errorsCreate.name && (
                  <p className="text-xs text-rose-500">{errorsCreate.name.message}</p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="description" className="text-xs font-semibold text-foreground">
                  Description <span className="text-muted-foreground font-normal">(Optional)</span>
                </Label>
                <Textarea
                  id="description"
                  {...registerCreate('description')}
                  placeholder="Primary floral workshop and fulfillment branch..."
                  rows={3}
                  className="rounded-xl border-border bg-background text-sm resize-none"
                />
                {errorsCreate.description && (
                  <p className="text-xs text-rose-500">{errorsCreate.description.message}</p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="logo_url" className="text-xs font-semibold text-foreground">
                  Logo URL <span className="text-muted-foreground font-normal">(Optional)</span>
                </Label>
                <Input
                  id="logo_url"
                  {...registerCreate('logo_url')}
                  placeholder="https://example.com/logo.png"
                  className="rounded-xl border-border bg-background text-sm font-mono"
                />
                {errorsCreate.logo_url && (
                  <p className="text-xs text-rose-500">{errorsCreate.logo_url.message}</p>
                )}
              </div>

              <div className="space-y-2">
                <Label htmlFor="banner_url" className="text-xs font-semibold text-foreground">
                  Banner URL <span className="text-muted-foreground font-normal">(Optional)</span>
                </Label>
                <Input
                  id="banner_url"
                  {...registerCreate('banner_url')}
                  placeholder="https://example.com/banner.png"
                  className="rounded-xl border-border bg-background text-sm font-mono"
                />
                {errorsCreate.banner_url && (
                  <p className="text-xs text-rose-500">{errorsCreate.banner_url.message}</p>
                )}
              </div>
            </form>
          )}
        </div>

        {/* Footer */}
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
            form={activeTab === 'assign' ? 'assign-staff-form' : 'create-staff-form'}
            size="sm"
            disabled={loading}
            className="text-xs rounded-xl font-semibold bg-primary hover:bg-primary/90 text-primary-foreground min-w-[130px]"
          >
            {loading ? (
              <>
                <Loader2 className="mr-2 h-3.5 w-3.5 animate-spin" />
                Saving...
              </>
            ) : activeTab === 'assign' ? (
              <>
                <UserPlus className="mr-1.5 h-3.5 w-3.5" />
                Assign Account
              </>
            ) : (
              <>
                <Building className="mr-1.5 h-3.5 w-3.5" />
                Create Entity
              </>
            )}
          </Button>
        </div>
      </SheetContent>
    </Sheet>
  );
}
