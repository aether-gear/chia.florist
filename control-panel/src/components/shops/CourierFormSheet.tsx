import React, { useState, useEffect } from 'react';
import {
  Truck,
  MapPin,
  Loader2,
  Save,
  ShieldCheck,
  AlertCircle,
  CheckCircle2,
  XCircle,
  Info,
} from 'lucide-react';
import { Button } from '../ui/button';
import { Input } from '../ui/input';
import { Label } from '../ui/label';
import { Switch } from '../ui/switch';
import { Badge } from '../ui/badge';
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetDescription,
  SheetClose,
} from '../ui/sheet';
import { type ShopCourier, type ShopAddress } from '../../models/Shop';

interface CourierFormSheetProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  shopId: string;
  shopName?: string;
  courier: ShopCourier | null;
  shopAddresses?: ShopAddress[];
  isAdmin?: boolean;
  onSave: (
    shopId: string,
    code: string,
    data: { active: boolean; name?: string; location_address?: string }
  ) => Promise<boolean>;
  onVerify?: (
    shopId: string,
    code: string,
    action: 'verify' | 'reject',
    rejectionReason?: string
  ) => Promise<boolean>;
  onSuccess: () => void;
}

export default function CourierFormSheet({
  open,
  onOpenChange,
  shopId,
  shopName,
  courier,
  shopAddresses = [],
  isAdmin = false,
  onSave,
  onVerify,
  onSuccess,
}: CourierFormSheetProps) {
  const [isActive, setIsActive] = useState(false);
  const [courierName, setCourierName] = useState('');
  const [locationAddress, setLocationAddress] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isVerifying, setIsVerifying] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Admin rejection modal/inline state
  const [showRejectInput, setShowRejectInput] = useState(false);
  const [rejectionReason, setRejectionReason] = useState('');

  // Sync state with incoming courier prop
  useEffect(() => {
    if (!open || !courier) return;
    setError(null);
    setShowRejectInput(false);
    setRejectionReason('');
    setIsActive(courier.active ?? false);
    setCourierName(courier.name || '');
    setLocationAddress(courier.location_address || '');
  }, [open, courier]);

  const handleCopyShopAddress = (fullAddr: string) => {
    if (!fullAddr) return;
    setLocationAddress(fullAddr);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!courier) return;

    setError(null);

    // Form validation when activated
    if (isActive) {
      if (!courierName.trim()) {
        setError('Please provide a branch / drop point name (e.g. JNE Pulo Gebang).');
        return;
      }
      if (!locationAddress.trim()) {
        setError('Please provide a pickup location address.');
        return;
      }
    }

    setIsSubmitting(true);
    try {
      await onSave(shopId, courier.code, {
        active: isActive,
        name: courierName.trim() || undefined,
        location_address: locationAddress.trim() || undefined,
      });
      onSuccess();
      onOpenChange(false);
    } catch (err: any) {
      setError(err.message || 'Failed to save courier configuration');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleAdminVerify = async (action: 'verify' | 'reject') => {
    if (!courier || !onVerify) return;
    setIsVerifying(true);
    setError(null);
    try {
      await onVerify(
        shopId,
        courier.code,
        action,
        action === 'reject' ? rejectionReason.trim() : undefined
      );
      onSuccess();
      onOpenChange(false);
    } catch (err: any) {
      setError(err.message || `Failed to ${action} courier`);
    } finally {
      setIsVerifying(false);
    }
  };

  const getStatusBadge = () => {
    if (!courier) return null;
    const status = courier.verification_status || 'unconfigured';
    switch (status) {
      case 'verified':
        return (
          <Badge className="bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border-0 rounded-lg">
            <CheckCircle2 className="h-3 w-3 mr-1" />
            Verified & Ready
          </Badge>
        );
      case 'pending':
        return (
          <Badge className="bg-amber-500/10 text-amber-600 dark:text-amber-400 border-0 rounded-lg animate-pulse">
            <AlertCircle className="h-3 w-3 mr-1" />
            Pending Verification
          </Badge>
        );
      case 'rejected':
        return (
          <Badge className="bg-rose-500/10 text-rose-600 dark:text-rose-400 border-0 rounded-lg">
            <XCircle className="h-3 w-3 mr-1" />
            Rejected
          </Badge>
        );
      default:
        return (
          <Badge className="bg-muted text-muted-foreground border-0 rounded-lg">
            Unconfigured
          </Badge>
        );
    }
  };

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent className="w-full sm:max-w-none md:w-[45vw] md:min-w-[45vw] overflow-y-auto border-l border-border/60 bg-background shadow-2xl p-6 sm:p-8">
        <SheetHeader className="mb-6 pb-4 border-b border-border/60">
          <div className="flex items-center justify-between gap-3">
            <div className="flex items-center gap-2.5">
              <div className="p-2 rounded-xl bg-primary/10 text-primary">
                <Truck className="h-5 w-5" />
              </div>
              <div>
                <SheetTitle className="font-display text-xl text-foreground flex items-center gap-2">
                  {courier?.branch_name || courier?.code?.toUpperCase()}
                  <Badge variant="outline" className="font-mono text-xs uppercase ml-1">
                    {courier?.code}
                  </Badge>
                </SheetTitle>
                <SheetDescription className="text-xs text-muted-foreground mt-0.5">
                  Configure delivery courier under {shopName || 'Shop'}
                </SheetDescription>
              </div>
            </div>
            <div>{getStatusBadge()}</div>
          </div>
        </SheetHeader>

        {error && (
          <div className="p-3 text-sm text-destructive bg-destructive/10 rounded-xl border border-destructive/20 mb-5 flex items-start gap-2">
            <AlertCircle className="h-4 w-4 mt-0.5 shrink-0" />
            <span>{error}</span>
          </div>
        )}

        {courier?.verification_status === 'rejected' && courier.rejection_reason && (
          <div className="p-3.5 mb-5 rounded-xl border border-rose-500/20 bg-rose-500/5 text-rose-700 dark:text-rose-400 text-xs space-y-1">
            <div className="font-semibold flex items-center gap-1.5">
              <XCircle className="h-3.5 w-3.5 shrink-0" />
              Previous Rejection Feedback
            </div>
            <p className="pl-5 text-muted-foreground">{courier.rejection_reason}</p>
          </div>
        )}

        {/* Admin Quick Verification Box for Pending Couriers */}
        {isAdmin && courier?.verification_status === 'pending' && onVerify && (
          <div className="p-4 mb-6 rounded-xl border border-amber-500/30 bg-amber-500/5 space-y-3">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2 text-sm font-semibold text-amber-700 dark:text-amber-400">
                <ShieldCheck className="h-4 w-4" />
                <span>Admin Review Action</span>
              </div>
              <Badge variant="outline" className="text-amber-700 border-amber-500/30 text-[10px]">
                Pending Staff Submission
              </Badge>
            </div>
            <p className="text-xs text-muted-foreground">
              This courier record was submitted for verification. As an Administrator, you can approve it to
              activate shipping or reject it with feedback.
            </p>

            {showRejectInput ? (
              <div className="space-y-2 pt-2">
                <Label htmlFor="rejectionReason" className="text-xs">
                  Rejection Reason (visible to staff)
                </Label>
                <Input
                  id="rejectionReason"
                  placeholder="e.g. Pickup address requires specific floor/gate details"
                  value={rejectionReason}
                  onChange={(e) => setRejectionReason(e.target.value)}
                  className="text-xs"
                />
                <div className="flex justify-end gap-2 pt-1">
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    className="h-8 text-xs rounded-lg"
                    onClick={() => setShowRejectInput(false)}
                    disabled={isVerifying}
                  >
                    Cancel
                  </Button>
                  <Button
                    type="button"
                    variant="destructive"
                    size="sm"
                    className="h-8 text-xs rounded-lg"
                    disabled={isVerifying}
                    onClick={() => handleAdminVerify('reject')}
                  >
                    {isVerifying ? <Loader2 className="h-3.5 w-3.5 animate-spin mr-1" /> : null}
                    Confirm Reject
                  </Button>
                </div>
              </div>
            ) : (
              <div className="flex items-center justify-end gap-2 pt-1">
                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  className="h-8 text-xs text-rose-600 border-rose-200 hover:bg-rose-50 dark:hover:bg-rose-950/30 rounded-lg"
                  disabled={isVerifying}
                  onClick={() => setShowRejectInput(true)}
                >
                  <XCircle className="h-3.5 w-3.5 mr-1" />
                  Reject
                </Button>
                <Button
                  type="button"
                  size="sm"
                  className="h-8 text-xs bg-emerald-600 hover:bg-emerald-700 text-white rounded-lg"
                  disabled={isVerifying}
                  onClick={() => handleAdminVerify('verify')}
                >
                  {isVerifying ? <Loader2 className="h-3.5 w-3.5 animate-spin mr-1" /> : <CheckCircle2 className="h-3.5 w-3.5 mr-1" />}
                  Verify & Activate
                </Button>
              </div>
            )}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-6">
          {/* Main Activation Toggle */}
          <div className="p-4 rounded-xl border border-border/80 bg-muted/20 flex items-center justify-between gap-4">
            <div className="space-y-0.5">
              <Label htmlFor="courierToggle" className="text-sm font-semibold cursor-pointer text-foreground">
                Activate Courier for this Shop
              </Label>
              <p className="text-xs text-muted-foreground">
                {isActive
                  ? 'Courier is enabled. Pickup name and location address are required below.'
                  : 'Courier is currently inactive (default status). Toggle on to configure record.'}
              </p>
            </div>
            <Switch
              id="courierToggle"
              checked={isActive}
              onCheckedChange={setIsActive}
            />
          </div>

          {/* Conditional Record Form - Shown only when toggle is ACTIVE */}
          {isActive ? (
            <div className="space-y-5 animate-in fade-in slide-in-from-top-2 duration-200">
              <div className="space-y-2">
                <Label htmlFor="courierName" className="text-sm font-medium">
                  Branch / Drop Point Name <span className="text-destructive">*</span>
                </Label>
                <Input
                  id="courierName"
                  placeholder="e.g. JNE Agen Pulo Gebang, SiCepat Hub Cakung"
                  value={courierName}
                  onChange={(e) => setCourierName(e.target.value)}
                  className="rounded-xl"
                  required
                />
                <p className="text-[11px] text-muted-foreground">
                  The specific branch, agent, or drop point assigned to this shop (e.g. JNE Pulo Gebang).
                </p>
              </div>

              <div className="space-y-2">
                <div className="flex items-center justify-between">
                  <Label htmlFor="locationAddress" className="text-sm font-medium">
                    Location Address <span className="text-destructive">*</span>
                  </Label>
                  {shopAddresses.length > 0 && (
                    <div className="flex items-center gap-1">
                      <span className="text-[11px] text-muted-foreground hidden sm:inline">Pre-fill:</span>
                      <select
                        aria-label="Pre-fill location address from shop address"
                        className="text-[11px] bg-background border border-border rounded-lg px-2 py-1 text-muted-foreground hover:text-foreground cursor-pointer focus:outline-none focus:ring-1 focus:ring-primary/40"
                        onChange={(e) => {
                          if (e.target.value) {
                            handleCopyShopAddress(e.target.value);
                            e.target.value = '';
                          }
                        }}
                        defaultValue=""
                      >
                        <option value="" disabled>
                          Select Shop Address...
                        </option>
                        {shopAddresses.map((addr) => (
                          <option key={addr.id} value={addr.full_address}>
                            {addr.label}: {addr.full_address.slice(0, 30)}...
                          </option>
                        ))}
                      </select>
                    </div>
                  )}
                </div>

                <textarea
                  id="locationAddress"
                  className="flex min-h-[90px] w-full rounded-xl border border-border bg-background px-3 py-2 text-sm text-foreground focus:outline-none focus:ring-1 focus:ring-primary/40 leading-relaxed placeholder:text-muted-foreground/60"
                  placeholder="Full street address, building/unit, district, postal code..."
                  value={locationAddress}
                  onChange={(e) => setLocationAddress(e.target.value)}
                  required
                />
                <p className="text-[11px] text-muted-foreground flex items-center gap-1">
                  <MapPin className="h-3 w-3 shrink-0 text-slate-400" />
                  <span>Physical pickup location for package dispatch (coordinates / long-lat not yet required).</span>
                </p>
              </div>

              {/* Notice regarding verification */}
              {!isAdmin ? (
                <div className="p-3 bg-muted/40 rounded-xl border border-border text-xs text-muted-foreground flex items-start gap-2.5">
                  <ShieldCheck className="h-4 w-4 text-primary shrink-0 mt-0.5" />
                  <div className="space-y-0.5">
                    <span className="font-semibold text-foreground">Verification Required:</span>
                    <p>
                      When you submit, this courier record will wait in <strong>Pending Verification</strong> status until an Administrator verifies it before order dispatching becomes active.
                    </p>
                  </div>
                </div>
              ) : (
                <div className="p-3 bg-emerald-500/5 rounded-xl border border-emerald-500/20 text-xs text-emerald-800 dark:text-emerald-300 flex items-start gap-2.5">
                  <ShieldCheck className="h-4 w-4 text-emerald-600 shrink-0 mt-0.5" />
                  <div className="space-y-0.5">
                    <span className="font-semibold">Administrator Privileges:</span>
                    <p>
                      As an Administrator, saving active configuration will automatically verify and activate this courier for immediate store operations.
                    </p>
                  </div>
                </div>
              )}
            </div>
          ) : (
            <div className="py-6 px-4 text-center border border-dashed border-border/80 rounded-xl text-xs text-muted-foreground space-y-1">
              <Info className="h-5 w-5 mx-auto text-muted-foreground/60 mb-1" />
              <p className="font-medium text-foreground">Courier is currently inactive</p>
              <p>Turn on the toggle above to configure the branch name and pickup location address.</p>
            </div>
          )}

          {/* Footer actions */}
          <div className="pt-4 border-t border-border/60 flex items-center justify-end gap-2">
            <SheetClose asChild>
              <Button type="button" variant="outline" className="rounded-xl">
                Cancel
              </Button>
            </SheetClose>
            <Button
              type="submit"
              disabled={isSubmitting || isVermittingDisabled(isActive, courierName, locationAddress)}
              className="bg-primary hover:bg-primary/90 text-primary-foreground rounded-xl"
            >
              {isSubmitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              <Save className="h-4 w-4 mr-1.5" />
              Save Configuration
            </Button>
          </div>
        </form>
      </SheetContent>
    </Sheet>
  );
}

function isVermittingDisabled(isActive: boolean, name: string, addr: string): boolean {
  if (!isActive) return false;
  return !name.trim() || !addr.trim();
}
