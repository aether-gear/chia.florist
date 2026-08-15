import { useState, useMemo } from 'react';
import { Users, UserPlus, RefreshCw, Copy, Check, Phone, Calendar } from 'lucide-react';
import { Button } from '../../components/ui/button';
import SearchInput from '../../components/SearchInput';
import { useStaffViewModel } from '../../viewmodels/useStaffViewModel';
import Pagination from '../../components/Pagination';
import { Skeleton } from '../../components/ui/skeleton';
import { DataCard, DataCardGridHeader, DataCardList } from '../../components/DataCard';
import EmptyState from '../../components/EmptyState';
import { Avatar, AvatarImage, AvatarFallback } from '../../components/ui/avatar';
import StaffFormSheet from '../../components/staff/StaffFormSheet';

export default function StaffListPage() {
  const { staff, total, loading, error, page, limit, setPage, refresh } = useStaffViewModel();

  const [searchQuery, setSearchQuery] = useState('');
  const [copiedId, setCopiedId] = useState<string | null>(null);

  // Staff Form Sheet state
  const [isSheetOpen, setIsSheetOpen] = useState(false);
  const [selectedStaffId, setSelectedStaffId] = useState<string | undefined>(undefined);
  const [selectedStaffName, setSelectedStaffName] = useState<string | undefined>(undefined);
  const [sheetMode, setSheetMode] = useState<'assign' | 'create'>('assign');

  const filteredStaff = useMemo(() => {
    if (!staff) return [];
    return staff.filter((item) =>
      (item.name && item.name.toLowerCase().includes(searchQuery.toLowerCase())) ||
      (item.username && item.username.toLowerCase().includes(searchQuery.toLowerCase())) ||
      (item.phone && item.phone.toLowerCase().includes(searchQuery.toLowerCase())) ||
      (item.id && item.id.toLowerCase().includes(searchQuery.toLowerCase()))
    );
  }, [staff, searchQuery]);

  const handleCopyId = (id: string, e: React.MouseEvent) => {
    e.stopPropagation();
    navigator.clipboard.writeText(id);
    setCopiedId(id);
    setTimeout(() => setCopiedId(null), 1500);
  };

  const openAssignForStaff = (staffId: string, staffName: string) => {
    setSelectedStaffId(staffId);
    setSelectedStaffName(staffName);
    setSheetMode('assign');
    setIsSheetOpen(true);
  };

  const openCreateStaff = () => {
    setSelectedStaffId(undefined);
    setSelectedStaffName(undefined);
    setSheetMode('create');
    setIsSheetOpen(true);
  };

  const openAssignGeneral = () => {
    setSelectedStaffId(staff.length > 0 ? staff[0].id : undefined);
    setSelectedStaffName(staff.length > 0 ? staff[0].name : undefined);
    setSheetMode('assign');
    setIsSheetOpen(true);
  };

  return (
    <div className="flex-col md:flex">
      <div className="flex-1 space-y-8 p-6 sm:p-8 lg:p-12 animate-in fade-in duration-300">
        {/* Header */}
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <h2 className="text-3xl font-bold font-display tracking-tight text-foreground">Staff Management</h2>
            <p className="text-muted-foreground text-sm">
              Manage staff entities, view accounts, and assign new team members
            </p>
          </div>
          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={openCreateStaff}
              className="rounded-xl text-xs font-semibold border-border text-foreground hover:bg-muted"
            >
              + New Staff Unit
            </Button>
            <Button
              size="sm"
              onClick={openAssignGeneral}
              className="bg-primary hover:bg-primary/90 text-primary-foreground rounded-xl text-xs font-semibold shadow-sm"
            >
              <UserPlus className="mr-1.5 h-3.5 w-3.5" /> Assign Staff Account
            </Button>
          </div>
        </div>

        {/* Filter and Search Bar */}
        <div className="space-y-6">
          <div className="pb-4 border-b border-border/60">
            <h3 className="text-xl font-bold font-display tracking-tight text-foreground">All Staff Members</h3>
            <p className="text-muted-foreground text-sm">
              Showing {staff.length} of {total} staff accounts.
            </p>
          </div>

          <div className="flex flex-col sm:flex-row items-center justify-between gap-4 w-full">
            <SearchInput
              value={searchQuery}
              onChange={setSearchQuery}
              placeholder="Search by name, @username, or UUID..."
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

          <div className="flex items-center justify-between px-1 text-xs text-muted-foreground">
            <span>Found {filteredStaff.length} matching staff</span>
          </div>

          {/* List Content */}
          <div className="flex flex-col gap-2">
            <DataCardGridHeader>
              <span className="col-span-4">Staff Member</span>
              <span className="col-span-3">Contact & Phone</span>
              <span className="col-span-3">Registered At</span>
              <span className="col-span-2 text-right">Action</span>
            </DataCardGridHeader>

            <DataCardList>
              {loading ? (
                Array.from({ length: 4 }).map((_, i) => (
                  <DataCard key={`skeleton-${i}`}>
                    <div className="col-span-4 flex items-center gap-3">
                      <Skeleton className="h-10 w-10 rounded-full bg-muted animate-pulse" />
                      <div className="space-y-1.5">
                        <Skeleton className="h-4 w-32 bg-muted animate-pulse" />
                        <Skeleton className="h-3 w-20 bg-muted animate-pulse" />
                      </div>
                    </div>
                    <div className="col-span-3"><Skeleton className="h-4 w-28 bg-muted animate-pulse" /></div>
                    <div className="col-span-3"><Skeleton className="h-4 w-24 bg-muted animate-pulse" /></div>
                    <div className="col-span-2 text-right"><Skeleton className="h-8 w-24 ml-auto bg-muted animate-pulse rounded-xl" /></div>
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
                      ? 'No staff members match your search criteria.'
                      : 'No staff accounts have been assigned yet. Click Assign Staff Account to get started.'
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

                  return (
                    <DataCard key={item.id}>
                      {/* Avatar & Name */}
                      <div className="col-span-1 md:col-span-4 flex items-center gap-3 min-w-0">
                        <Avatar className="h-10 w-10 shrink-0 ring-2 ring-primary/10">
                          {item.avatar_url && (
                            <AvatarImage src={item.avatar_url} alt={item.name} className="object-cover" />
                          )}
                          <AvatarFallback className="font-bold text-xs uppercase bg-primary text-primary-foreground">
                            {fallbackInitials}
                          </AvatarFallback>
                        </Avatar>
                        <div className="min-w-0">
                          <h4 className="font-bold font-display text-sm text-foreground truncate">
                            {item.name}
                          </h4>
                          <div className="flex items-center gap-1.5 mt-0.5">
                            {item.username && (
                              <span className="text-xs text-muted-foreground font-medium">
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

                      {/* Contact / Phone */}
                      <div className="col-span-1 md:col-span-3 text-xs text-muted-foreground flex items-center gap-1.5">
                        <Phone className="h-3.5 w-3.5 text-muted-foreground/60 shrink-0" />
                        <span className="truncate">{item.phone || 'No phone number'}</span>
                      </div>

                      {/* Registered Date */}
                      <div className="col-span-1 md:col-span-3 text-xs text-muted-foreground flex items-center gap-1.5">
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

                      {/* Actions */}
                      <div className="col-span-1 md:col-span-2 text-right" onClick={(e) => e.stopPropagation()}>
                        <Button
                          size="sm"
                          variant="outline"
                          onClick={() => openAssignForStaff(item.id, item.name)}
                          className="rounded-xl text-xs border-border hover:border-primary hover:text-primary hover:bg-primary/5 transition-all"
                        >
                          <UserPlus className="mr-1.5 h-3.5 w-3.5" /> Assign Account
                        </Button>
                      </div>
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
              itemNamePlural="staff"
            />
          </div>
        </div>
      </div>

      {/* Right Panel Drawer */}
      <StaffFormSheet
        open={isSheetOpen}
        onOpenChange={setIsSheetOpen}
        selectedStaffId={selectedStaffId}
        selectedStaffName={selectedStaffName}
        staffList={staff}
        defaultMode={sheetMode}
        onSuccess={() => refresh()}
      />
    </div>
  );
}
