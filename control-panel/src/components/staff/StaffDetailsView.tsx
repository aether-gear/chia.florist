import { useState } from 'react';
import {
  Edit,
  X,
  Building,
  Calendar,
  Link as LinkIcon,
  AlertTriangle,
  Copy,
  Check,
  UserPlus,
  ShieldCheck,
} from 'lucide-react';
import { Button } from '../ui/button';
import { Avatar, AvatarImage, AvatarFallback } from '../ui/avatar';
import StaffAccountsList from './StaffAccountsList';
import StaffShopPermissionsSection from './StaffShopPermissionsSection';
import type { Staff, StaffAccountMember, StaffShopPermission, SaveStaffShopPermissionPayload } from '@/models/Staff';

interface StaffDetailsViewProps {
  staff: Staff;
  accounts: StaffAccountMember[];
  onClose: () => void;
  onEditStaff: () => void;
  onBindAccount: () => void;
  onUnbindAccount: (accountId: string) => Promise<void>;
  fetchShopPermissions?: (staffId: string) => Promise<StaffShopPermission[]>;
  saveShopPermission?: (staffId: string, payload: SaveStaffShopPermissionPayload) => Promise<void>;
  deleteShopPermission?: (staffId: string, shopId: string) => Promise<void>;
  loading?: boolean;
}

export default function StaffDetailsView({
  staff,
  accounts,
  onClose,
  onEditStaff,
  onBindAccount,
  onUnbindAccount,
  fetchShopPermissions,
  saveShopPermission,
  deleteShopPermission,
  loading = false,
}: StaffDetailsViewProps) {

  const [copiedId, setCopiedId] = useState(false);

  const fallbackInitials = staff.name
    ? staff.name
        .split(' ')
        .map((n) => n[0])
        .join('')
        .substring(0, 2)
        .toUpperCase()
    : 'ST';

  const isBound = accounts.length > 0;

  const handleCopyId = () => {
    navigator.clipboard.writeText(staff.id);
    setCopiedId(true);
    setTimeout(() => setCopiedId(false), 1500);
  };

  return (
    <div className="border border-border/60 rounded-3xl bg-background overflow-hidden shadow-sm animate-in fade-in slide-in-from-right-2 duration-200">
      {/* Top Banner Section (Clean top section with 15% reduced height) */}
      <div className="relative h-28 sm:h-36 w-full bg-gradient-to-r from-primary/20 via-emerald-900/30 to-primary/15 overflow-hidden border-b border-border/40">
        {staff.banner_url ? (
          <img
            src={staff.banner_url}
            alt={staff.name}
            className="w-full h-full object-cover"
          />
        ) : (
          <div className="w-full h-full flex items-center justify-center opacity-25">
            <svg
              className="w-16 h-16 text-primary"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="1.5"
            >
              <path d="M12 22c0-5.523-4.477-10-10-10 5.523 0 10-4.477 10-10 0 5.523 4.477 10 10 10-5.523 0-10 4.477-10 10z" />
            </svg>
          </div>
        )}

        {/* Close Button Top Right */}
        <Button
          variant="ghost"
          size="icon"
          onClick={onClose}
          className="absolute top-2.5 right-2.5 h-7 w-7 rounded-full bg-background/80 hover:bg-background text-foreground backdrop-blur-md shadow-xs border border-border/40"
          title="Close detail view"
        >
          <X className="h-3.5 w-3.5" />
        </Button>
      </div>

      {/* Main Body */}
      <div className="px-6 py-5 space-y-6">
        {/* Header: Avatar, Name, and Quick Actions */}
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b border-border/50">
          <div className="flex items-center gap-3.5 min-w-0">
            {/* Clean Avatar */}
            <Avatar className="h-14 w-14 sm:h-16 sm:w-16 rounded-2xl ring-2 ring-primary/20 shadow-xs shrink-0 bg-muted">
              {staff.logo_url || staff.avatar_url ? (
                <AvatarImage
                  src={staff.logo_url || staff.avatar_url || ''}
                  alt={staff.name}
                  className="object-cover rounded-2xl"
                />
              ) : null}
              <AvatarFallback className="font-bold text-base sm:text-lg uppercase bg-primary text-primary-foreground rounded-2xl">
                {fallbackInitials}
              </AvatarFallback>
            </Avatar>

            {/* Title & Subtitle */}
            <div className="space-y-0.5 min-w-0">
              <h3 className="font-bold font-display text-lg sm:text-xl text-foreground truncate">
                {staff.name}
              </h3>
              <div className="flex flex-wrap items-center gap-2 text-xs text-muted-foreground">
                {staff.username && (
                  <span className="font-medium text-foreground/80">@{staff.username}</span>
                )}
                <span>•</span>
                <button
                  type="button"
                  onClick={handleCopyId}
                  className="font-mono text-[11px] text-muted-foreground hover:text-foreground flex items-center gap-1 px-1.5 py-0.5 rounded bg-muted/60 hover:bg-muted transition-colors"
                  title="Click to copy Staff UUID"
                >
                  <span>{staff.id.slice(0, 8)}...</span>
                  {copiedId ? (
                    <Check className="h-2.5 w-2.5 text-primary" />
                  ) : (
                    <Copy className="h-2.5 w-2.5" />
                  )}
                </button>
              </div>
            </div>
          </div>

          {/* Action Buttons */}
          <div className="flex items-center gap-2 shrink-0">
            {!isBound && (
              <Button
                type="button"
                size="sm"
                onClick={onBindAccount}
                className="text-xs rounded-xl font-semibold bg-primary hover:bg-primary/90 text-primary-foreground"
              >
                <UserPlus className="h-3.5 w-3.5 mr-1.5" />
                Bind Account
              </Button>
            )}
            <Button
              type="button"
              variant="outline"
              size="sm"
              onClick={onEditStaff}
              className="text-xs rounded-xl font-semibold border-border hover:bg-muted"
            >
              <Edit className="h-3.5 w-3.5 mr-1.5 text-muted-foreground" />
              Edit Staff
            </Button>
          </div>
        </div>

        {/* Quick Metadata Bar */}
        <div className="grid grid-cols-2 sm:grid-cols-3 gap-3 p-3.5 rounded-2xl bg-muted/20 border border-border/40 text-xs">
          <div>
            <div className="text-[10px] font-bold uppercase tracking-wider text-muted-foreground">
              Binding Status
            </div>
            <div className="mt-1">
              {isBound ? (
                <span className="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-xs font-semibold bg-emerald-500/10 text-emerald-700 dark:text-emerald-400 border border-emerald-500/20">
                  <LinkIcon className="h-2.5 w-2.5" />
                  Bound
                </span>
              ) : (
                <span className="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-xs font-semibold bg-amber-500/10 text-amber-700 dark:text-amber-400 border border-amber-500/20">
                  <AlertTriangle className="h-2.5 w-2.5" />
                  Unbound
                </span>
              )}
            </div>
          </div>

          <div>
            <div className="text-[10px] font-bold uppercase tracking-wider text-muted-foreground">
              Registered Date
            </div>
            <div className="font-semibold text-foreground mt-1 flex items-center gap-1">
              <Calendar className="h-3 w-3 text-muted-foreground" />
              {staff.created_at
                ? new Date(staff.created_at).toLocaleDateString('en-GB', {
                    day: 'numeric',
                    month: 'short',
                    year: 'numeric',
                  })
                : 'N/A'}
            </div>
          </div>

          <div className="col-span-2 sm:col-span-1">
            <div className="text-[10px] font-bold uppercase tracking-wider text-muted-foreground">
              Entity Type
            </div>
            <div className="font-semibold text-foreground mt-1 flex items-center gap-1">
              <Building className="h-3 w-3 text-muted-foreground" />
              Staff Unit
            </div>
          </div>
        </div>

        {/* Description / Bio */}
        {staff.description ? (
          <div className="space-y-1.5">
            <h4 className="text-xs font-bold font-display uppercase tracking-wider text-muted-foreground">
              About Staff Unit
            </h4>
            <p className="text-xs text-foreground/90 leading-relaxed bg-muted/10 p-3.5 rounded-2xl border border-border/30">
              {staff.description}
            </p>
          </div>
        ) : null}

        {/* Bound Accounts Section */}
        <div className="space-y-3 pt-2">
          <div className="flex items-center justify-between pb-2 border-b border-border/40">
            <div className="flex items-center gap-1.5">
              <ShieldCheck className="h-4 w-4 text-primary" />
              <h4 className="text-xs font-bold font-display uppercase tracking-wider text-foreground">
                {isBound ? 'Bound User Account' : 'User Account'}
              </h4>
            </div>
          </div>

          <StaffAccountsList
            staffId={staff.id}
            accounts={accounts}
            onUnbind={onUnbindAccount}
            onBindNew={onBindAccount}
            loading={loading}
          />
        </div>

        {/* Shop Access & Permissions Section */}
        {fetchShopPermissions && saveShopPermission && deleteShopPermission && (
          <StaffShopPermissionsSection
            staffId={staff.id}
            staffName={staff.name}
            fetchPermissions={fetchShopPermissions}
            savePermission={saveShopPermission}
            deletePermission={deleteShopPermission}
          />
        )}
      </div>
    </div>
  );
}

