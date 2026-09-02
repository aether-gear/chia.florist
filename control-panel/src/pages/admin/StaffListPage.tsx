import { useState, useMemo, useRef, useEffect } from 'react';
import {
  Users,
  Building,
  RefreshCw,
  Copy,
  Check,
  Calendar,
  Clock,
  MoreHorizontal,
  Edit,
  Trash2,
  UserPlus,
  Link as LinkIcon,
  AlertTriangle,
  Loader2,
} from 'lucide-react';
import { Button } from '../../components/ui/button';
import SearchInput from '../../components/SearchInput';
import { useStaffViewModel } from '../../viewmodels/useStaffViewModel';
import Pagination from '../../components/Pagination';
import { Skeleton } from '../../components/ui/skeleton';
import { DataCard, DataCardGridHeader, DataCardList } from '../../components/DataCard';
import EmptyState from '../../components/EmptyState';
import { Avatar, AvatarImage, AvatarFallback } from '../../components/ui/avatar';
import {
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
} from '../../components/ui/dropdown-menu';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '../../components/ui/dialog';
import StaffFormSheet from '../../components/staff/StaffFormSheet';
import StaffDetailsView from '../../components/staff/StaffDetailsView';
import { formatRelativeLastLogin } from '../../lib/timeUtils';
import type { Staff } from '@/models/Staff';

export default function StaffListPage() {
  const {
    staff,
    total,
    accountsMap,
    loading,
    error,
    page,
    limit,
    setPage,
    refresh,
    createStaff,
    createStaffWithAccount,
    updateStaff,
    deleteStaff,
    addStaffAccount,
    removeStaffAccount,
    fetchStaffShopPermissions,
    saveStaffShopPermission,
    deleteStaffShopPermission,
  } = useStaffViewModel();


  const [searchQuery, setSearchQuery] = useState('');
  const [copiedId, setCopiedId] = useState<string | null>(null);

  // Master-Detail Split State
  const [activeStaffId, setActiveStaffId] = useState<string | null>(null);
  const isDetailOpen = Boolean(activeStaffId);
  const detailSectionRef = useRef<HTMLDivElement | null>(null);

  // Staff Form Sheet states
  const [isSheetOpen, setIsSheetOpen] = useState(false);
  const [sheetMode, setSheetMode] = useState<'create' | 'edit' | 'bind'>('create');
  const [staffToForm, setStaffToForm] = useState<Staff | null>(null);

  // Delete Dialog state
  const [staffToDelete, setStaffToDelete] = useState<Staff | null>(null);
  const [deleting, setDeleting] = useState(false);

  useEffect(() => {
    if (activeStaffId && detailSectionRef.current) {
      detailSectionRef.current.scrollIntoView({ behavior: 'smooth', block: 'start' });
    }
  }, [activeStaffId]);

  const filteredStaff = useMemo(() => {
    if (!staff) return [];
    return staff.filter((item) => {
      const query = searchQuery.toLowerCase();
      const accounts = accountsMap[item.id] || [];
      const hasMatchingAccount = accounts.some(
        (a) =>
          a.email.toLowerCase().includes(query) ||
          a.name.toLowerCase().includes(query) ||
          a.username.toLowerCase().includes(query)
      );

      return (
        (item.name && item.name.toLowerCase().includes(query)) ||
        (item.username && item.username.toLowerCase().includes(query)) ||
        (item.phone && item.phone.toLowerCase().includes(query)) ||
        (item.id && item.id.toLowerCase().includes(query)) ||
        hasMatchingAccount
      );
    });
  }, [staff, accountsMap, searchQuery]);

  const activeStaff = useMemo(() => {
    if (!activeStaffId) return null;
    return staff.find((s) => s.id === activeStaffId) || null;
  }, [staff, activeStaffId]);

  const handleCopyId = (id: string, e: React.MouseEvent) => {
    e.stopPropagation();
    navigator.clipboard.writeText(id);
    setCopiedId(id);
    setTimeout(() => setCopiedId(null), 1500);
  };

  const openCreateStaff = () => {
    setStaffToForm(null);
    setSheetMode('create');
    setIsSheetOpen(true);
  };

  const openEditStaff = (item: Staff) => {
    setStaffToForm(item);
    setSheetMode('edit');
    setIsSheetOpen(true);
  };

  const openBindForStaff = (item: Staff) => {
    setStaffToForm(item);
    setSheetMode('bind');
    setIsSheetOpen(true);
  };

  const handleDeleteStaff = async () => {
    if (!staffToDelete) return;
    setDeleting(true);
    try {
      await deleteStaff(staffToDelete.id);
      if (activeStaffId === staffToDelete.id) {
        setActiveStaffId(null);
      }
      setStaffToDelete(null);
    } catch (err: any) {
      alert(err.message || 'Failed to delete staff unit');
    } finally {
      setDeleting(false);
    }
  };

  return (
    <div className="flex-col md:flex">
      <div className="flex-1 space-y-8 p-6 sm:p-8 lg:p-12 animate-in fade-in slide-in-from-left-4 duration-300">
        {/* Header */}
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <h2 className="text-3xl font-bold font-display tracking-tight text-foreground">
              Staff Management
            </h2>
            <p className="text-muted-foreground text-sm">
              Create staff units, manage bound login accounts, and configure permissions
            </p>
          </div>
          <div className="flex items-center gap-2">
            <Button
              size="sm"
              onClick={openCreateStaff}
              className="bg-primary hover:bg-primary/90 text-primary-foreground rounded-xl text-xs font-semibold shadow-sm"
            >
              <Building className="mr-1.5 h-3.5 w-3.5" /> + New Staff
            </Button>
          </div>
        </div>

        {/* Filter and Search Bar */}
        <div className="space-y-6">
          <div className="pb-4 border-b border-border/60">
            <h3 className="text-xl font-bold font-display tracking-tight text-foreground">
              All Staff Units
            </h3>
            <p className="text-muted-foreground text-sm">
              Showing {staff.length} of {total} registered staff entities.
            </p>
          </div>

          <div className="flex flex-col sm:flex-row items-center justify-between gap-4 w-full">
            <SearchInput
              value={searchQuery}
              onChange={setSearchQuery}
              placeholder="Search by staff name, @username, email, or UUID..."
              className="relative flex-1 max-w-md w-full"
            />
            <div className="flex items-center gap-2 justify-end w-full sm:w-auto">
              <Button
                variant="outline"
                onClick={() => refresh()}
                disabled={loading}
                className="flex items-center gap-1.5 border-border text-foreground hover:text-primary hover:bg-primary/5 rounded-xl transition-colors text-xs"
              >
                <RefreshCw className={`h-3.5 w-3.5 ${loading ? 'animate-spin' : ''}`} />
                Refresh
              </Button>
            </div>
          </div>

          {/* Master-Detail Split Grid Layout */}
          <div className="grid grid-cols-12 gap-6 items-start">
            {/* Left Master List Column */}
            <div className={isDetailOpen ? 'col-span-12 lg:col-span-5 space-y-4' : 'col-span-12 space-y-4'}>
              <div className="flex items-center justify-between px-1 text-xs text-muted-foreground">
                <span>Found {filteredStaff.length} matching staff</span>
                {!isDetailOpen && <span className="hidden sm:inline">Click any card to inspect full details</span>}
              </div>

              {!isDetailOpen && (
                <DataCardGridHeader>
                  <span className="col-span-4">Staff Entity</span>
                  <span className="col-span-2">Binding Status</span>
                  <span className="col-span-3">Last Login</span>
                  <span className="col-span-2">Registered Date</span>
                  <span className="col-span-1 text-right">Actions</span>
                </DataCardGridHeader>
              )}

              <DataCardList>
                {loading ? (
                  Array.from({ length: 4 }).map((_, i) => (
                    <DataCard key={`skeleton-${i}`}>
                      <div className="col-span-12 flex items-center gap-3">
                        <Skeleton className="h-10 w-10 rounded-xl bg-muted animate-pulse" />
                        <div className="space-y-1.5 flex-1">
                          <Skeleton className="h-4 w-32 bg-muted animate-pulse" />
                          <Skeleton className="h-3 w-20 bg-muted animate-pulse" />
                        </div>
                      </div>
                    </DataCard>
                  ))
                ) : error ? (
                  <EmptyState
                    title="Failed to load staff"
                    description={error}
                    className="py-12 border-0 bg-transparent text-destructive"
                  />
                ) : filteredStaff.length === 0 ? (
                  <EmptyState
                    icon={<Users className="h-8 w-8 text-muted-foreground/60 mb-2 mx-auto" />}
                    title="No staff members found"
                    description={
                      searchQuery
                        ? 'No staff units match your search criteria.'
                        : 'No staff units exist yet. Click "+ New Staff" to create one.'
                    }
                    className="py-12 border border-dashed border-border/80 rounded-2xl bg-zinc-50/20"
                  />
                ) : (
                  filteredStaff.map((item) => {
                    const fallbackInitials = item.name
                      ? item.name
                          .split(' ')
                          .map((n) => n[0])
                          .join('')
                          .substring(0, 2)
                          .toUpperCase()
                      : 'ST';

                    const accounts = accountsMap[item.id] || [];
                    const isBound = accounts.length > 0;
                    const isSelected = activeStaffId === item.id;

                    const latestLogin =
                      item.last_login_at ||
                      (accounts.length > 0
                        ? accounts.reduce<string | null>((latest, acc) => {
                            if (!acc.last_login_at) return latest;
                            if (!latest) return acc.last_login_at;
                            return new Date(acc.last_login_at) > new Date(latest)
                              ? acc.last_login_at
                              : latest;
                          }, null)
                        : null);

                    return (
                      <DataCard
                        key={item.id}
                        selected={isSelected}
                        onClick={() => setActiveStaffId(isSelected ? null : item.id)}
                        className="cursor-pointer transition-all duration-200"
                      >
                        {isDetailOpen ? (
                          /* Shrunken Compact Row (Avatar + Name + Binding Status + Last Login + Registered Date) */
                          <div className="col-span-12 flex items-center justify-between gap-3 min-w-0">
                            {/* Left: Avatar & Identity */}
                            <div className="flex items-center gap-3 min-w-0 flex-1">
                              <Avatar className="h-9 w-9 shrink-0 rounded-xl ring-1 ring-primary/20">
                                {item.logo_url || item.avatar_url ? (
                                  <AvatarImage
                                    src={item.logo_url || item.avatar_url || ''}
                                    alt={item.name}
                                    className="object-cover rounded-xl"
                                  />
                                ) : null}
                                <AvatarFallback className="font-bold text-xs uppercase bg-primary text-primary-foreground rounded-xl">
                                  {fallbackInitials}
                                </AvatarFallback>
                              </Avatar>
                              <div className="min-w-0">
                                <h4 className="font-bold font-display text-sm text-foreground truncate">
                                  {item.name}
                                </h4>
                                {item.username && (
                                  <p className="text-[11px] text-muted-foreground truncate">
                                    @{item.username}
                                  </p>
                                )}
                              </div>
                            </div>

                            {/* Right: Shrink Fields (Binding Status & Last Login / Registered Date) */}
                            <div className="flex flex-col items-end gap-0.5 shrink-0 text-right">
                              {isBound ? (
                                <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-semibold bg-emerald-500/10 text-emerald-700 dark:text-emerald-400 border border-emerald-500/20">
                                  <LinkIcon className="h-2.5 w-2.5" />
                                  Bound ({accounts.length})
                                </span>
                              ) : (
                                <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-semibold bg-amber-500/10 text-amber-700 dark:text-amber-400 border border-amber-500/20">
                                  <AlertTriangle className="h-2.5 w-2.5" />
                                  Unbound
                                </span>
                              )}

                              {latestLogin ? (
                                (() => {
                                  const text = formatRelativeLastLogin(latestLogin);
                                  const isActiveNow = text === 'Active now';
                                  return (
                                    <span className={`text-[10px] font-medium flex items-center gap-1 ${isActiveNow ? 'text-primary' : 'text-foreground/80'}`}>
                                      <Clock className={`h-2.5 w-2.5 ${isActiveNow ? 'text-primary' : 'text-primary/70'}`} />
                                      {isActiveNow && <span className="h-1.5 w-1.5 rounded-full bg-primary animate-pulse" />}
                                      <span>{text}</span>
                                    </span>
                                  );
                                })()
                              ) : (
                                <span className="text-[10px] text-muted-foreground/60 flex items-center gap-1">
                                  <Clock className="h-2.5 w-2.5 text-muted-foreground/40" />
                                  Never logged in
                                </span>
                              )}
                            </div>
                          </div>
                        ) : (
                          /* Full Expanded 12-Column Row */
                          <>
                            {/* Col 4: Avatar & Staff Identity */}
                            <div className="col-span-1 md:col-span-4 flex items-center gap-3 min-w-0">
                              <Avatar className="h-10 w-10 shrink-0 rounded-xl ring-1 ring-primary/20">
                                {item.logo_url || item.avatar_url ? (
                                  <AvatarImage
                                    src={item.logo_url || item.avatar_url || ''}
                                    alt={item.name}
                                    className="object-cover rounded-xl"
                                  />
                                ) : null}
                                <AvatarFallback className="font-bold text-xs uppercase bg-primary text-primary-foreground rounded-xl">
                                  {fallbackInitials}
                                </AvatarFallback>
                              </Avatar>
                              <div className="min-w-0">
                                <h4 className="font-bold font-display text-sm text-foreground truncate">
                                  {item.name}
                                </h4>
                                <div className="flex items-center gap-1.5 mt-0.5">
                                  {item.username && (
                                    <span className="text-xs text-muted-foreground font-medium truncate">
                                      @{item.username}
                                    </span>
                                  )}
                                  <button
                                    type="button"
                                    onClick={(e) => handleCopyId(item.id, e)}
                                    className="text-[10px] font-mono text-muted-foreground/70 hover:text-foreground flex items-center gap-1 px-1.5 py-0.5 rounded bg-muted/60 hover:bg-muted transition-colors"
                                    title="Click to copy Staff UUID"
                                  >
                                    <span>{item.id.slice(0, 8)}...</span>
                                    {copiedId === item.id ? (
                                      <Check className="h-2.5 w-2.5 text-primary" />
                                    ) : (
                                      <Copy className="h-2.5 w-2.5" />
                                    )}
                                  </button>
                                </div>
                              </div>
                            </div>

                            {/* Col 2: Binding Status */}
                            <div className="col-span-1 md:col-span-2 flex items-center gap-2 text-xs">
                              {isBound ? (
                                <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-semibold bg-emerald-500/10 text-emerald-700 dark:text-emerald-400 border border-emerald-500/20">
                                  <LinkIcon className="h-3 w-3" />
                                  Bound
                                </span>
                              ) : (
                                <span className="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-semibold bg-amber-500/10 text-amber-700 dark:text-amber-400 border border-amber-500/20">
                                  <AlertTriangle className="h-3 w-3" />
                                  Unbound
                                </span>
                              )}
                            </div>

                            {/* Col 3: Last Login At (Beside Binding Status) */}
                            <div className="col-span-1 md:col-span-3 text-xs flex items-center gap-1.5">
                              {latestLogin ? (
                                (() => {
                                  const text = formatRelativeLastLogin(latestLogin);
                                  const isActiveNow = text === 'Active now';
                                  return (
                                    <span className={`flex items-center gap-1.5 font-medium ${isActiveNow ? 'text-primary font-semibold' : 'text-foreground/90'}`}>
                                      <Clock className={`h-3.5 w-3.5 ${isActiveNow ? 'text-primary' : 'text-primary/70'} shrink-0`} />
                                      {isActiveNow && <span className="h-2 w-2 rounded-full bg-primary animate-pulse" />}
                                      <span>{text}</span>
                                    </span>
                                  );
                                })()
                              ) : (
                                <span className="flex items-center gap-1.5 text-muted-foreground/60">
                                  <Clock className="h-3.5 w-3.5 text-muted-foreground/40 shrink-0" />
                                  <span>Never logged in</span>
                                </span>
                              )}
                            </div>

                            {/* Col 2: Registered Date */}
                            <div className="col-span-1 md:col-span-2 text-xs text-muted-foreground flex items-center gap-1.5">
                              <Calendar className="h-3.5 w-3.5 text-muted-foreground/60 shrink-0" />
                              <span>
                                {item.created_at
                                  ? new Date(item.created_at).toLocaleDateString('en-GB', {
                                      day: 'numeric',
                                      month: 'short',
                                      year: 'numeric',
                                    })
                                  : 'N/A'}
                              </span>
                            </div>

                            {/* Col 1: Triple Dots Actions Menu */}
                            <div
                              className="col-span-1 md:col-span-1 text-right"
                              onClick={(e) => e.stopPropagation()}
                            >
                              <DropdownMenu>
                                <DropdownMenuTrigger asChild>
                                  <Button
                                    variant="ghost"
                                    size="icon"
                                    className="h-8 w-8 rounded-lg hover:bg-muted text-muted-foreground hover:text-foreground ml-auto"
                                  >
                                    <MoreHorizontal className="h-4 w-4" />
                                    <span className="sr-only">Open actions</span>
                                  </Button>
                                </DropdownMenuTrigger>
                                <DropdownMenuContent align="end" className="w-48 p-1 rounded-xl shadow-lg border-border/70">
                                  <DropdownMenuItem
                                    onClick={() => openEditStaff(item)}
                                    className="cursor-pointer flex items-center gap-2 px-3 py-2 text-xs font-medium rounded-lg"
                                  >
                                    <Edit className="h-3.5 w-3.5 text-muted-foreground" />
                                    <span>Edit Staff & Account</span>
                                  </DropdownMenuItem>

                                  {!isBound && (
                                    <DropdownMenuItem
                                      onClick={() => openBindForStaff(item)}
                                      className="cursor-pointer flex items-center gap-2 px-3 py-2 text-xs font-medium rounded-lg text-primary"
                                    >
                                      <UserPlus className="h-3.5 w-3.5" />
                                      <span>Bind Account</span>
                                    </DropdownMenuItem>
                                  )}

                                  <DropdownMenuSeparator className="my-1 bg-border/40" />

                                  <DropdownMenuItem
                                    onClick={() => setStaffToDelete(item)}
                                    className="focus:bg-destructive/10 focus:text-destructive text-destructive cursor-pointer flex items-center gap-2 px-3 py-2 text-xs font-medium rounded-lg"
                                  >
                                    <Trash2 className="h-3.5 w-3.5" />
                                    <span>Delete Staff Unit</span>
                                  </DropdownMenuItem>
                                </DropdownMenuContent>
                              </DropdownMenu>
                            </div>
                          </>
                        )}
                      </DataCard>
                    );
                  })
                )}
              </DataCardList>

              <Pagination
                currentPage={page}
                totalPages={Math.ceil(total / limit) || 1}
                totalItems={total}
                limit={limit}
                onPageChange={setPage}
                itemNamePlural="staff units"
              />
            </div>

            {/* Right Column: LinkedIn-Style Staff Details View */}
            {isDetailOpen && activeStaff && (
              <div ref={detailSectionRef} className="col-span-12 lg:col-span-7 sticky top-20 animate-in fade-in slide-in-from-right-4 duration-300">
                <StaffDetailsView
                  staff={activeStaff}
                  accounts={accountsMap[activeStaff.id] || []}
                  onClose={() => setActiveStaffId(null)}
                  onEditStaff={() => openEditStaff(activeStaff)}
                  onBindAccount={() => openBindForStaff(activeStaff)}
                  onUnbindAccount={async (accountId) => {
                    await removeStaffAccount(activeStaff.id, accountId);
                  }}
                  fetchShopPermissions={fetchStaffShopPermissions}
                  saveShopPermission={saveStaffShopPermission}
                  deleteShopPermission={deleteStaffShopPermission}
                />

              </div>
            )}
          </div>
        </div>
      </div>

      {/* Right Overlay Panel Drawer (Create / Edit / Bind) */}
      <StaffFormSheet
        open={isSheetOpen}
        onOpenChange={setIsSheetOpen}
        staffToEdit={staffToForm}
        accounts={staffToForm ? accountsMap[staffToForm.id] || [] : []}
        mode={sheetMode}
        createStaff={createStaff}
        createStaffWithAccount={createStaffWithAccount}
        updateStaff={updateStaff}
        addStaffAccount={addStaffAccount}
        removeStaffAccount={removeStaffAccount}
        onSuccess={() => refresh()}
      />

      {/* Delete Confirmation Dialog */}
      <Dialog open={Boolean(staffToDelete)} onOpenChange={(open) => !open && setStaffToDelete(null)}>
        <DialogContent className="max-w-md rounded-2xl">
          <DialogHeader>
            <DialogTitle className="text-base font-bold font-display text-foreground">
              Delete Staff Unit?
            </DialogTitle>
            <DialogDescription className="text-xs text-muted-foreground mt-1">
              Are you sure you want to delete <strong className="text-foreground">{staffToDelete?.name}</strong>? This will soft-delete the staff unit and cascade unbind all associated user accounts.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter className="gap-2 sm:gap-0 mt-4">
            <Button
              type="button"
              variant="outline"
              size="sm"
              onClick={() => setStaffToDelete(null)}
              disabled={deleting}
              className="text-xs rounded-xl"
            >
              Cancel
            </Button>
            <Button
              type="button"
              variant="destructive"
              size="sm"
              onClick={handleDeleteStaff}
              disabled={deleting}
              className="text-xs rounded-xl font-semibold"
            >
              {deleting ? <Loader2 className="h-3.5 w-3.5 animate-spin mr-1.5" /> : null}
              Confirm Delete
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
