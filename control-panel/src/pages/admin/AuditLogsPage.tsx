import { useState } from 'react';
import { History, ArrowUpDown, Eye, RefreshCw } from 'lucide-react';
import { Button } from '../../components/ui/button';
import { Input } from '../../components/ui/input';
import { Label } from '../../components/ui/label';
import { Badge } from '../../components/ui/badge';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '../../components/ui/table';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '../../components/ui/card';
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetDescription } from '../../components/ui/sheet';
import { useAuditLogsViewModel } from '../../viewmodels/useAuditLogsViewModel';
import type { AuditLog } from '../../models/AuditLog';
import EmptyState from '../../components/EmptyState';
import SearchInput from '../../components/SearchInput';
import StatusBadge from '../../components/StatusBadge';
import Pagination from '../../components/Pagination';
import { Skeleton } from '../../components/ui/skeleton';

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

        <Card className="border-0 shadow-none bg-zinc-50/40 dark:bg-slate-900/40">
          <CardHeader>
            <CardTitle className="font-bold font-display tracking-tight text-lg text-foreground">All Audit Logs</CardTitle>
            <CardDescription className="text-muted-foreground text-sm">
              View system audit trails and security activity logs.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="mb-4 flex flex-col md:flex-row md:items-end justify-between gap-4">
              {/* Left Side: Filter and Search */}
              <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-4 flex-1 w-full md:w-auto">
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
                    placeholder="UUID of the actor"
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

              {/* Right Side: Refresh */}
              <div className="flex items-center gap-2 justify-end w-full md:w-auto pb-0.5">
                <Button
                  variant="outline"
                  onClick={() => refresh()}
                  className="flex items-center gap-1.5 border-border text-foreground hover:text-primary hover:bg-primary/5 rounded-xl transition-colors"
                >
                  <RefreshCw className="h-4 w-4" />
                  Refresh
                </Button>
              </div>
            </div>

            {error ? (
              <EmptyState
                title="Failed to load audit logs"
                description={error}
                actionLabel="Retry"
                onAction={() => refresh()}
                className="flex h-48 flex-col items-center justify-center text-center p-4 gap-2 border-0 bg-transparent text-destructive"
              />
            ) : (
              <div className="rounded-2xl border border-border overflow-hidden">
                <Table>
                  <TableHeader className="bg-muted/50">
                    <TableRow>
                      <TableHead 
                        className="cursor-pointer hover:bg-muted/70 transition-colors text-foreground"
                        onClick={() => handleSort('date')}
                      >
                        <div className="flex items-center gap-1.5 font-semibold">
                          Date & Time
                          <ArrowUpDown className="h-3.5 w-3.5" />
                        </div>
                      </TableHead>
                      <TableHead 
                        className="cursor-pointer hover:bg-muted/70 transition-colors text-foreground"
                        onClick={() => handleSort('action')}
                      >
                        <div className="flex items-center gap-1.5 font-semibold">
                          Action
                          <ArrowUpDown className="h-3.5 w-3.5" />
                        </div>
                      </TableHead>
                      <TableHead className="font-semibold text-foreground">Category / Resource</TableHead>
                      <TableHead className="font-semibold text-foreground">Actor ID</TableHead>
                      <TableHead className="font-semibold text-foreground text-center">Outcome</TableHead>
                      <TableHead className="font-semibold text-foreground">Client IP</TableHead>
                      <TableHead className="text-right font-semibold text-foreground">Actions</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {loading ? (
                      Array.from({ length: 5 }).map((_, i) => (
                        <TableRow key={`skeleton-${i}`}>
                          <TableCell><Skeleton className="h-5 w-32 animate-pulse bg-muted" /></TableCell>
                          <TableCell><Skeleton className="h-5 w-20 animate-pulse bg-muted" /></TableCell>
                          <TableCell><Skeleton className="h-5 w-40 animate-pulse bg-muted" /></TableCell>
                          <TableCell><Skeleton className="h-5 w-28 animate-pulse bg-muted" /></TableCell>
                          <TableCell className="text-center"><Skeleton className="h-5 w-16 mx-auto animate-pulse bg-muted" /></TableCell>
                          <TableCell><Skeleton className="h-5 w-24 animate-pulse bg-muted" /></TableCell>
                          <TableCell className="text-right"><Skeleton className="h-5 w-8 ml-auto animate-pulse bg-muted" /></TableCell>
                        </TableRow>
                      ))
                    ) : logs.length === 0 ? (
                      <TableRow>
                        <TableCell colSpan={7} className="p-0">
                          <EmptyState
                            icon={<History className="h-8 w-8 mb-2 mx-auto text-muted-foreground" />}
                            title="No audit logs found"
                            description="No audit logs found matching your criteria."
                            className="flex h-48 flex-col items-center justify-center text-center p-4 gap-1.5 border-0 bg-transparent"
                          />
                        </TableCell>
                      </TableRow>
                    ) : (
                      logs.map((log) => (
                        <TableRow 
                          key={log.id} 
                          className="hover:bg-muted/55 cursor-pointer transition-colors"
                          onClick={() => setSelectedLog(log)}
                        >
                          <TableCell className="font-medium text-foreground text-sm">
                            {formatDate(log.created_at)}
                          </TableCell>
                          <TableCell>
                            <Badge className="bg-muted text-muted-foreground border-0 font-semibold uppercase text-xs rounded-lg">
                              {log.action}
                            </Badge>
                          </TableCell>
                          <TableCell className="text-muted-foreground text-xs font-mono">
                            {log.category} / {log.resource}
                          </TableCell>
                          <TableCell className="text-muted-foreground font-mono text-xs max-w-[120px] truncate" title={log.actor_id}>
                            {log.actor_id}
                          </TableCell>
                          <TableCell className="text-center">
                            <StatusBadge status={log.outcome} />
                          </TableCell>
                          <TableCell className="text-muted-foreground text-xs font-mono">
                            {log.client_ip}
                          </TableCell>
                          <TableCell className="text-right" onClick={(e) => e.stopPropagation()}>
                            <Button
                              variant="ghost"
                              size="icon"
                              className="text-muted-foreground hover:text-primary hover:bg-primary/5"
                              onClick={() => setSelectedLog(log)}
                            >
                              <Eye className="h-4 w-4" />
                            </Button>
                          </TableCell>
                        </TableRow>
                      ))
                    )}
                  </TableBody>
                </Table>
              </div>
            )}

            <Pagination
              currentPage={page}
              totalPages={totalPages}
              totalItems={total}
              limit={limit}
              onPageChange={setPage}
              itemNamePlural="logs"
              className="flex items-center justify-between border-t border-border p-4"
            />
          </CardContent>
        </Card>
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

