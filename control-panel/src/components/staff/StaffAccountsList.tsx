import { useState } from 'react';
import { Mail, Phone, Clock, Calendar, ShieldCheck, Trash2, Loader2, UserPlus, AlertTriangle } from 'lucide-react';
import { Button } from '../ui/button';
import { Avatar, AvatarFallback } from '../ui/avatar';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '../ui/dialog';
import type { StaffAccountMember } from '@/models/Staff';
import { formatRelativeLastLogin } from '@/lib/timeUtils';

interface StaffAccountsListProps {
  staffId: string;
  accounts: StaffAccountMember[];
  onUnbind?: (accountId: string) => Promise<void>;
  onBindNew?: () => void;
  readOnly?: boolean;
  loading?: boolean;
}

export default function StaffAccountsList({
  staffId: _staffId,
  accounts,
  onUnbind,
  onBindNew,
  readOnly = false,
  loading = false,
}: StaffAccountsListProps) {
  const [accountToUnbind, setAccountToUnbind] = useState<StaffAccountMember | null>(null);
  const [unbinding, setUnbinding] = useState(false);

  const handleConfirmUnbind = async () => {
    if (!accountToUnbind || !onUnbind) return;
    setUnbinding(true);
    try {
      await onUnbind(accountToUnbind.account_id);
      setAccountToUnbind(null);
    } catch (err: any) {
      alert(err.message || 'Failed to unbind staff account');
    } finally {
      setUnbinding(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center py-8">
        <Loader2 className="h-5 w-5 animate-spin text-primary" />
        <span className="ml-2 text-xs text-muted-foreground">Loading bound accounts...</span>
      </div>
    );
  }

  if (accounts.length === 0) {
    return (
      <div className="rounded-2xl border border-dashed border-border/80 p-6 text-center bg-zinc-50/20 dark:bg-zinc-900/20">
        <div className="mx-auto w-9 h-9 rounded-full bg-amber-500/10 text-amber-600 dark:text-amber-400 flex items-center justify-center mb-2">
          <AlertTriangle className="h-4 w-4" />
        </div>
        <h4 className="text-xs font-bold text-foreground">No Account Bound</h4>
        <p className="text-[11px] text-muted-foreground mt-0.5 max-w-xs mx-auto">
          This staff entity does not have any active user login credentials bound yet.
        </p>
        {!readOnly && onBindNew && (
          <Button
            type="button"
            size="sm"
            onClick={onBindNew}
            className="mt-3 text-xs rounded-xl bg-primary hover:bg-primary/90 text-primary-foreground font-semibold"
          >
            <UserPlus className="mr-1.5 h-3.5 w-3.5" /> Bind Account Now
          </Button>
        )}
      </div>
    );
  }

  return (
    <div className="space-y-3">
      {accounts.map((acc) => {
        const initials = acc.name
          ? acc.name
              .split(' ')
              .map((n) => n[0])
              .join('')
              .substring(0, 2)
              .toUpperCase()
          : 'ST';

        const isAdminRole = acc.role?.code === 'staff_admin';

        return (
          <div
            key={acc.account_id}
            className="flex flex-col sm:flex-row sm:items-center justify-between gap-3 p-3.5 rounded-xl bg-background border border-border/50 hover:border-border transition-all"
          >
            {/* Account Info */}
            <div className="flex items-center gap-3 min-w-0">
              <Avatar className="h-9 w-9 shrink-0 ring-1 ring-primary/20">
                <AvatarFallback className="font-bold text-xs uppercase bg-primary text-primary-foreground">
                  {initials}
                </AvatarFallback>
              </Avatar>
              <div className="min-w-0 space-y-0.5">
                <div className="flex items-center gap-2">
                  <span className="font-bold font-display text-xs text-foreground truncate">
                    {acc.name}
                  </span>
                  {acc.username && (
                    <span className="text-[11px] text-muted-foreground font-medium">
                      @{acc.username}
                    </span>
                  )}
                  <span
                    className={`inline-flex items-center gap-0.5 px-2 py-0.5 rounded-md text-[10px] font-semibold border ${
                      isAdminRole
                        ? 'bg-amber-500/10 text-amber-700 dark:text-amber-300 border-amber-500/20'
                        : 'bg-primary/10 text-primary border-primary/20'
                    }`}
                  >
                    <ShieldCheck className="h-2.5 w-2.5" />
                    {acc.role?.name || (isAdminRole ? 'Staff Admin' : 'Staff Member')}
                  </span>
                </div>

                <div className="flex flex-wrap items-center gap-3 text-[11px] text-muted-foreground">
                  <span className="flex items-center gap-1">
                    <Mail className="h-3 w-3 text-muted-foreground/70" />
                    {acc.email}
                  </span>
                  {acc.phone && (
                    <span className="flex items-center gap-1">
                      <Phone className="h-3 w-3 text-muted-foreground/70" />
                      {acc.phone}
                    </span>
                  )}
                  <span className="flex items-center gap-1">
                    <Clock className="h-3 w-3 text-muted-foreground/70" />
                    {acc.last_login_at
                      ? `Last active: ${formatRelativeLastLogin(acc.last_login_at)}`
                      : 'Never logged in'}
                  </span>
                  <span className="flex items-center gap-1">
                    <Calendar className="h-3 w-3 text-muted-foreground/70" />
                    {acc.created_at
                      ? `Added: ${new Date(acc.created_at).toLocaleDateString('en-GB', {
                          day: 'numeric',
                          month: 'short',
                          year: 'numeric',
                        })}`
                      : ''}
                  </span>
                </div>
              </div>
            </div>

            {/* Unbind Action */}
            {!readOnly && onUnbind && (
              <Button
                type="button"
                variant="ghost"
                size="sm"
                onClick={() => setAccountToUnbind(acc)}
                className="h-8 px-2 text-xs text-rose-500 hover:text-rose-600 hover:bg-rose-50 dark:hover:bg-rose-950/30 rounded-lg shrink-0 self-end sm:self-center"
                title="Unbind this account from staff entity"
              >
                <Trash2 className="h-3.5 w-3.5 mr-1" />
                Unbind
              </Button>
            )}
          </div>
        );
      })}

      {/* Confirmation Dialog */}
      <Dialog open={Boolean(accountToUnbind)} onOpenChange={(open) => !open && setAccountToUnbind(null)}>
        <DialogContent className="max-w-md rounded-2xl">
          <DialogHeader>
            <DialogTitle className="text-base font-bold font-display text-foreground">
              Unbind Staff Account?
            </DialogTitle>
            <DialogDescription className="text-xs text-muted-foreground mt-1">
              Are you sure you want to unbind the account for{' '}
              <strong className="text-foreground">{accountToUnbind?.name}</strong> ({accountToUnbind?.email})? This user will no longer have access to this staff entity.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter className="gap-2 sm:gap-0 mt-4">
            <Button
              type="button"
              variant="outline"
              size="sm"
              onClick={() => setAccountToUnbind(null)}
              disabled={unbinding}
              className="text-xs rounded-xl"
            >
              Cancel
            </Button>
            <Button
              type="button"
              variant="destructive"
              size="sm"
              onClick={handleConfirmUnbind}
              disabled={unbinding}
              className="text-xs rounded-xl font-semibold"
            >
              {unbinding ? <Loader2 className="h-3.5 w-3.5 animate-spin mr-1.5" /> : null}
              Confirm Unbind
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
