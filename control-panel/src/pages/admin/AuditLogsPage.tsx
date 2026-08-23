import { useState } from 'react';
import { History, ArrowUpDown, Eye, RefreshCw, Zap, ShieldAlert } from 'lucide-react';
import { Button } from '../../components/ui/button';
import { Input } from '../../components/ui/input';
import { Label } from '../../components/ui/label';
import { Badge } from '../../components/ui/badge';
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetDescription } from '../../components/ui/sheet';
import { useAuditLogsViewModel } from '../../viewmodels/useAuditLogsViewModel';
import type { AuditLog } from '../../models/AuditLog';
import EmptyState from '../../components/EmptyState';
import SearchInput from '../../components/SearchInput';
import StatusBadge from '../../components/StatusBadge';
import Pagination from '../../components/Pagination';
import { Skeleton } from '../../components/ui/skeleton';
import { DataCard, DataCardGridHeader, DataCardList } from '../../components/DataCard';

export default function AuditLogsPage() {
  const {
    data,
    loading,
    error,
    page,
    sort,
    actionFilter,
    userIdFilter,
    startDate,
    endDate,
    setPage,
    setSort,
    setActionFilter,
    setUserIdFilter,
    setStartDate,
    setEndDate,
    refresh
  } = useAuditLogsViewModel();

  const [selectedLog, setSelectedLog] = useState<AuditLog | null>(null);

  const handleSort = (field: string) => {
    let newDirection = 'desc';
    const [currentField, currentDirection] = sort.split(':');
    if (currentField === field && currentDirection === 'desc') {
      newDirection = 'asc';
    }
    setSort(`${field}:${newDirection}`);
    setPage(1);
  };

  const formatDate = (dateStr: string) => {
    try {
      return new Date(dateStr).toLocaleString('id-ID', {
        dateStyle: 'medium',
        timeStyle: 'medium'
      });
    } catch {
      return dateStr;
    }
  };

  const logs = data?.audit_logs || [];
  const total = data?.total || 0;
  const limit = data?.limit || 10;
  const totalPages = Math.ceil(total / limit);

  return (
    <div className="flex-col md:flex">
      <div className="flex-1 space-y-8 p-6 sm:p-8 lg:p-12 animate-in fade-in duration-300">
        <div className="flex items-center justify-between space-y-2">
          <div>
            <h2 className="text-3xl font-bold font-display tracking-tight text-foreground">Audit Logs</h2>
            <p className="text-muted-foreground text-sm">
              Monitor user actions, resource changes, and security outcomes.
            </p>
          </div>
        </div>

        <div className="space-y-6">
          <div className="pb-4 border-b border-border/60">
            <h3 className="text-xl font-bold font-display tracking-tight text-foreground">All Audit Logs</h3>
            <p className="text-muted-foreground text-sm">View system audit trails and security activity logs.</p>
          </div>

          <div className="space-y-4">
            <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-4 items-end">
              <div className="space-y-1.5">
                <Label htmlFor="action-filter" className="text-xs font-bold uppercase tracking-wider text-muted-foreground">Action Name</Label>
                <SearchInput
                  id="action-filter"
                  placeholder="e.g. signin, save_shop"
                  className="relative w-full text-foreground"
                  value={actionFilter}
                  onChange={(val) => {
                    setActionFilter(val);
                    setPage(1);
                  }}
                />
              </div>

              <div className="space-y-1.5">
                <Label htmlFor="user-filter" className="text-xs font-bold uppercase tracking-wider text-muted-foreground">User / Actor ID</Label>
                <SearchInput
                  id="user-filter"
                  placeholder="UUID of actor"
                  className="relative w-full text-foreground"
                  value={userIdFilter}
                  onChange={(val) => {
                    setUserIdFilter(val);
                    setPage(1);
                  }}
                />
              </div>

              <div className="space-y-1.5">
                <Label htmlFor="start-date" className="text-xs font-bold uppercase tracking-wider text-muted-foreground">Start Date</Label>
                <Input
                  id="start-date"
                  type="date"
                  className="text-sm rounded-xl border border-border bg-background text-foreground"
                  value={startDate}
                  onChange={(e) => {
                    setStartDate(e.target.value);
                    setPage(1);
                  }}
                />
              </div>

              <div className="space-y-1.5">
                <Label htmlFor="end-date" className="text-xs font-bold uppercase tracking-wider text-muted-foreground">End Date</Label>
                <Input
                  id="end-date"
                  type="date"
                  className="text-sm rounded-xl border border-border bg-background text-foreground"
                  value={endDate}
                  onChange={(e) => {
                    setEndDate(e.target.value);
                    setPage(1);
                  }}
                />
              </div>
            </div>

            <div className="flex items-center justify-between px-1 text-xs text-muted-foreground">
              <span>Found {total} audit logs</span>
              <Button
                variant="outline"
                size="sm"
                onClick={() => refresh()}
                disabled={loading}
                className="flex items-center gap-1.5 border-border text-foreground hover:text-primary hover:bg-primary/5 rounded-xl transition-colors"
              >
                <RefreshCw className={`h-3.5 w-3.5 ${loading ? 'animate-spin' : ''}`} />
                Refresh
              </Button>
            </div>
          </div>

          {/* Content */}
          <div className="flex flex-col gap-2">
            <DataCardGridHeader>
              <span className="col-span-3">
                <button className="flex items-center gap-1 hover:text-primary" onClick={() => handleSort('created_at')}>
                  Timestamp <ArrowUpDown className="h-3 w-3" />
                </button>
              </span>
              <span className="col-span-3">
                <button className="flex items-center gap-1 hover:text-primary" onClick={() => handleSort('action')}>
                  Action <ArrowUpDown className="h-3 w-3" />
                </button>
              </span>
              <span className="col-span-3">Actor / User ID</span>
              <span className="col-span-2">Outcome</span>
              <span className="col-span-1 text-right">Details</span>
            </DataCardGridHeader>

            <DataCardList>
              {loading ? (
                Array.from({ length: 4 }).map((_, i) => (
                  <DataCard key={`skeleton-${i}`}>
                    <div className="col-span-3"><Skeleton className="h-4 w-28 bg-muted animate-pulse" /></div>
                    <div className="col-span-3"><Skeleton className="h-5 w-24 bg-muted animate-pulse" /></div>
                    <div className="col-span-3"><Skeleton className="h-4 w-32 bg-muted animate-pulse" /></div>
                    <div className="col-span-2"><Skeleton className="h-5 w-16 bg-muted animate-pulse rounded-full" /></div>
                    <div className="col-span-1 text-right"><Skeleton className="h-8 w-8 ml-auto bg-muted animate-pulse rounded-lg" /></div>
                  </DataCard>
                ))
              ) : error ? (
                <EmptyState title="Failed to load audit logs" description={error} className="py-12 border-0 bg-transparent text-destructive" />
              ) : logs.length === 0 ? (
                <EmptyState icon={<History className="h-8 w-8 text-slate-400 mb-2 mx-auto" />} title="No audit logs found" description="No activity logs match the current search filters." className="py-12 border border-dashed border-border/80 rounded-2xl bg-zinc-50/10" />
              ) : (
                logs.map((log) => (
                  <DataCard key={log.id} onClick={() => setSelectedLog(log)}>
                    <div className="col-span-1 md:col-span-3 text-xs text-muted-foreground">
                      <span className="md:hidden font-sans text-muted-foreground mr-1">Time:</span>
                      {formatDate(log.created_at)}
                    </div>

                    <div className="col-span-1 md:col-span-3">
                      {log.action === 'critical_stockout_risk_alert' ? (
                        <Badge className="font-mono text-xs font-semibold bg-destructive/15 text-destructive border-destructive/30 flex items-center gap-1 w-fit">
                          <Zap className="h-3 w-3" /> critical_stockout_risk_alert
                        </Badge>
                      ) : log.action === 'payment_anomaly_detected' ? (
                        <Badge className="font-mono text-xs font-semibold bg-amber-500/15 text-amber-600 dark:text-amber-400 border-amber-500/30 flex items-center gap-1 w-fit">
                          <ShieldAlert className="h-3 w-3" /> payment_anomaly_detected
                        </Badge>
                      ) : (
                        <Badge variant="outline" className="font-mono text-xs font-semibold bg-muted/40 text-foreground border-border/80">
                          {log.action}
                        </Badge>
                      )}
                    </div>

                    <div className="col-span-1 md:col-span-3 font-mono text-xs text-muted-foreground truncate">
                      <span className="md:hidden font-sans text-muted-foreground mr-1">Actor:</span>
                      {log.actor_id ? log.actor_id : <span className="italic opacity-60">System</span>}
                    </div>

                    <div className="col-span-1 md:col-span-2">
                      <StatusBadge status={log.outcome} className="scale-90 origin-left" />
                    </div>

                    <div className="col-span-1 md:col-span-1 text-right" onClick={(e) => e.stopPropagation()}>
                      <Button
                        variant="ghost"
                        size="icon"
                        className="h-8 w-8 text-muted-foreground hover:text-primary hover:bg-primary/5 rounded-lg"
                        onClick={() => setSelectedLog(log)}
                      >
                        <Eye className="h-4 w-4" />
                      </Button>
                    </div>
                  </DataCard>
                ))
              )}
            </DataCardList>

            <Pagination
              currentPage={page}
              totalPages={totalPages}
              totalItems={total}
              limit={limit}
              onPageChange={setPage}
              itemNamePlural="logs"
            />
          </div>
        </div>
      </div>

      {/* Details Sheet Overlay */}
      <Sheet open={!!selectedLog} onOpenChange={(open) => !open && setSelectedLog(null)}>
        <SheetContent className="sm:max-w-xl overflow-y-auto">
          {selectedLog && (
            <>
              <SheetHeader className="mb-6">
                <SheetTitle className="text-xl flex items-center gap-2 font-display">
                  <History className="h-5 w-5 text-primary" />
                  Audit Log Details
                </SheetTitle>
                <SheetDescription>
                  Full audit entry parameters and request payload.
                </SheetDescription>
              </SheetHeader>

              <div className="space-y-6">
                {/* Status and Action summary */}
                <div className="grid grid-cols-2 gap-4 rounded-2xl bg-muted/40 p-4 border border-border">
                  <div>
                    <span className="block text-xs font-bold text-muted-foreground uppercase tracking-wider">ACTION</span>
                    <span className="font-bold font-display text-foreground uppercase tracking-wide text-sm mt-1 block">{selectedLog.action}</span>
                  </div>
                  <div>
                    <span className="block text-xs font-bold text-muted-foreground uppercase tracking-wider">OUTCOME</span>
                    <StatusBadge status={selectedLog.outcome} className="mt-1" />
                  </div>
                </div>

                {/* Properties table */}
                <div className="space-y-3">
                  <h4 className="text-sm font-bold font-display text-foreground border-b border-border pb-2">Properties</h4>
                  <div className="grid grid-cols-3 gap-2 text-sm">
                    <span className="font-medium text-muted-foreground">Log ID</span>
                    <span className="col-span-2 font-mono text-foreground text-xs break-all">{selectedLog.id}</span>

                    <span className="font-medium text-muted-foreground">Timestamp</span>
                    <span className="col-span-2 text-foreground">{formatDate(selectedLog.created_at)}</span>

                    <span className="font-medium text-muted-foreground">Category</span>
                    <span className="col-span-2 text-foreground">{selectedLog.category}</span>

                    <span className="font-medium text-muted-foreground">Resource</span>
                    <span className="col-span-2 text-foreground">{selectedLog.resource} (ID: <span className="font-mono text-xs bg-muted px-1.5 py-0.5 rounded-lg">{selectedLog.resource_id}</span>)</span>

                    <span className="font-medium text-muted-foreground">Actor ID</span>
                    <span className="col-span-2 font-mono text-foreground text-xs break-all">{selectedLog.actor_id}</span>

                    <span className="font-medium text-muted-foreground">Request ID</span>
                    <span className="col-span-2 font-mono text-foreground text-xs">{selectedLog.request_id}</span>

                    <span className="font-medium text-muted-foreground">Client IP</span>
                    <span className="col-span-2 font-mono text-foreground text-xs">{selectedLog.client_ip}</span>

                    <span className="font-medium text-muted-foreground">User Agent</span>
                    <span className="col-span-2 text-muted-foreground text-xs break-all">{selectedLog.metadata?.user_agent || 'N/A'}</span>
                  </div>
                </div>

                {/* Metadata JSON tree */}
                <div className="space-y-2">
                  <h4 className="text-sm font-bold font-display text-foreground">Metadata Payload</h4>
                  <div className="rounded-xl bg-slate-950 p-4 text-xs font-mono text-emerald-400 overflow-x-auto max-h-[300px]">
                    <pre>{JSON.stringify(selectedLog.metadata, null, 2)}</pre>
                  </div>
                </div>
              </div>
            </>
          )}
        </SheetContent>
      </Sheet>
    </div>
  );
}
