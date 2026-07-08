import { useState } from 'react';
import { History, Search, Loader2, ChevronLeft, ChevronRight, ArrowUpDown, Eye, RefreshCw } from 'lucide-react';
import { Button } from '../../components/ui/button';
import { Input } from '../../components/ui/input';
import { Badge } from '../../components/ui/badge';
import { Label } from '../../components/ui/label';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '../../components/ui/table';
import { Card, CardContent } from '../../components/ui/card';
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetDescription } from '../../components/ui/sheet';
import { useAuditLogsViewModel } from '../../viewmodels/useAuditLogsViewModel';
import type { AuditLog } from '../../models/AuditLog';

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
      <div className="flex-1 space-y-4 p-8 pt-6">
        <div className="flex items-center justify-between space-y-2">
          <div>
            <h2 className="text-3xl font-bold tracking-tight">Audit Logs</h2>
            <p className="text-muted-foreground">
              Monitor user actions, resource changes, and security outcomes.
            </p>
          </div>
          <div className="flex items-center space-x-2">
            <Button
              variant="outline"
              onClick={() => refresh()}
              className="flex items-center gap-1.5 border-slate-200 text-slate-600 hover:text-indigo-600 transition-colors"
            >
              <RefreshCw className="h-4 w-4" />
              Refresh
            </Button>
          </div>
        </div>

        <Card className="shadow-sm border-slate-100 dark:border-slate-800 bg-white dark:bg-slate-900/50">
          <CardContent className="p-4 space-y-4">
            <div className="grid grid-cols-1 gap-4 md:grid-cols-4">
              <div className="space-y-1">
                <Label htmlFor="action-filter" className="text-xs font-semibold text-slate-500">Action Name</Label>
                <div className="relative">
                  <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-slate-400" />
                  <Input
                    id="action-filter"
                    placeholder="e.g. signin, save_shop"
                    className="pl-8 text-sm"
                    value={actionFilter}
                    onChange={(e) => {
                      setActionFilter(e.target.value);
                      setPage(1);
                    }}
                  />
                </div>
              </div>

              <div className="space-y-1">
                <Label htmlFor="user-filter" className="text-xs font-semibold text-slate-500">User / Actor ID</Label>
                <div className="relative">
                  <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-slate-400" />
                  <Input
                    id="user-filter"
                    placeholder="UUID of the actor"
                    className="pl-8 text-sm"
                    value={userIdFilter}
                    onChange={(e) => {
                      setUserIdFilter(e.target.value);
                      setPage(1);
                    }}
                  />
                </div>
              </div>

              <div className="space-y-1">
                <Label htmlFor="start-date" className="text-xs font-semibold text-slate-500">Start Date</Label>
                <Input
                  id="start-date"
                  type="date"
                  className="text-sm"
                  value={startDate}
                  onChange={(e) => {
                    setStartDate(e.target.value);
                    setPage(1);
                  }}
                />
              </div>

              <div className="space-y-1">
                <Label htmlFor="end-date" className="text-xs font-semibold text-slate-500">End Date</Label>
                <Input
                  id="end-date"
                  type="date"
                  className="text-sm"
                  value={endDate}
                  onChange={(e) => {
                    setEndDate(e.target.value);
                    setPage(1);
                  }}
                />
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Main Content Table */}
        <Card className="shadow-md border-0 bg-white/70 dark:bg-slate-900/40 backdrop-blur-md border border-slate-200/50 dark:border-slate-800/50">
          <CardContent className="p-0">
            {error ? (
              <div className="flex h-48 flex-col items-center justify-center p-4">
                <p className="text-destructive font-medium mb-1">Failed to load audit logs</p>
                <p className="text-xs text-slate-500">{error}</p>
                <Button size="sm" onClick={() => refresh()} className="mt-4 bg-indigo-600">Retry</Button>
              </div>
            ) : loading ? (
              <div className="flex h-64 items-center justify-center">
                <Loader2 className="h-8 w-8 animate-spin text-indigo-500" />
              </div>
            ) : (
              <div className="rounded-md border border-slate-100 dark:border-slate-800 overflow-hidden">
                <Table>
                  <TableHeader className="bg-slate-50/70 dark:bg-slate-900/50">
                    <TableRow>
                      <TableHead 
                        className="cursor-pointer hover:bg-slate-100/50 transition-colors"
                        onClick={() => handleSort('date')}
                      >
                        <div className="flex items-center gap-1.5 font-semibold text-slate-700 dark:text-slate-300">
                          Date & Time
                          <ArrowUpDown className="h-3.5 w-3.5" />
                        </div>
                      </TableHead>
                      <TableHead 
                        className="cursor-pointer hover:bg-slate-100/50 transition-colors"
                        onClick={() => handleSort('action')}
                      >
                        <div className="flex items-center gap-1.5 font-semibold text-slate-700 dark:text-slate-300">
                          Action
                          <ArrowUpDown className="h-3.5 w-3.5" />
                        </div>
                      </TableHead>
                      <TableHead className="font-semibold text-slate-700 dark:text-slate-300">Category / Resource</TableHead>
                      <TableHead className="font-semibold text-slate-700 dark:text-slate-300">Actor ID</TableHead>
                      <TableHead className="font-semibold text-slate-700 dark:text-slate-300 text-center">Outcome</TableHead>
                      <TableHead className="font-semibold text-slate-700 dark:text-slate-300">Client IP</TableHead>
                      <TableHead className="text-right font-semibold text-slate-700 dark:text-slate-300">Actions</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {logs.length === 0 ? (
                      <TableRow>
                        <TableCell colSpan={7} className="text-center h-48 text-slate-400">
                          No audit logs found matching your criteria.
                        </TableCell>
                      </TableRow>
                    ) : (
                      logs.map((log) => (
                        <TableRow 
                          key={log.id} 
                          className="hover:bg-slate-50/50 dark:hover:bg-slate-850/30 cursor-pointer transition-colors"
                          onClick={() => setSelectedLog(log)}
                        >
                          <TableCell className="font-medium text-slate-800 dark:text-slate-200 text-sm">
                            {formatDate(log.created_at)}
                          </TableCell>
                          <TableCell>
                            <Badge className="bg-slate-100 dark:bg-slate-800 text-slate-800 dark:text-slate-200 border-0 font-semibold uppercase text-xs">
                              {log.action}
                            </Badge>
                          </TableCell>
                          <TableCell className="text-slate-600 dark:text-slate-350 text-xs font-mono">
                            {log.category} / {log.resource}
                          </TableCell>
                          <TableCell className="text-slate-500 dark:text-slate-400 font-mono text-xs max-w-[120px] truncate" title={log.actor_id}>
                            {log.actor_id}
                          </TableCell>
                          <TableCell className="text-center">
                            <Badge
                              className={
                                log.outcome === 'success'
                                  ? 'bg-emerald-100 dark:bg-emerald-950/30 text-emerald-800 dark:text-emerald-450 border-0'
                                  : 'bg-rose-100 dark:bg-rose-950/30 text-rose-800 dark:text-rose-450 border-0'
                              }
                            >
                              {log.outcome}
                            </Badge>
                          </TableCell>
                          <TableCell className="text-slate-600 dark:text-slate-350 text-xs font-mono">
                            {log.client_ip}
                          </TableCell>
                          <TableCell className="text-right" onClick={(e) => e.stopPropagation()}>
                            <Button
                              variant="ghost"
                              size="icon"
                              className="text-slate-400 hover:text-indigo-650 dark:hover:text-indigo-400"
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

            {/* Pagination Panel */}
            {!loading && !error && totalPages > 1 && (
              <div className="flex items-center justify-between border-t border-slate-100 p-4">
                <p className="text-sm text-slate-500">
                  Showing page <span className="font-semibold text-slate-800">{page}</span> of{' '}
                  <span className="font-semibold text-slate-800">{totalPages}</span> ({total} logs)
                </p>
                <div className="flex items-center space-x-2">
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setPage(p => Math.max(1, p - 1))}
                    disabled={page === 1}
                    className="border-slate-200"
                  >
                    <ChevronLeft className="h-4 w-4" />
                    Previous
                  </Button>
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setPage(p => Math.min(totalPages, p + 1))}
                    disabled={page === totalPages}
                    className="border-slate-200"
                  >
                    Next
                    <ChevronRight className="h-4 w-4" />
                  </Button>
                </div>
              </div>
            )}
          </CardContent>
        </Card>
      </div>

      {/* Details Sheet Overlay */}
      <Sheet open={!!selectedLog} onOpenChange={(open) => !open && setSelectedLog(null)}>
        <SheetContent className="sm:max-w-xl overflow-y-auto">
          {selectedLog && (
            <>
              <SheetHeader className="mb-6">
                <SheetTitle className="text-xl flex items-center gap-2 text-slate-900 dark:text-slate-100">
                  <History className="h-5 w-5 text-indigo-600 dark:text-indigo-400" />
                  Audit Log Details
                </SheetTitle>
                <SheetDescription className="text-slate-500 dark:text-slate-400">
                  Full audit entry parameters and request payload.
                </SheetDescription>
              </SheetHeader>

              <div className="space-y-6">
                {/* Status and Action summary */}
                <div className="grid grid-cols-2 gap-4 rounded-lg bg-slate-50 dark:bg-slate-900 p-4 border border-slate-100 dark:border-slate-800">
                  <div>
                    <span className="block text-xs font-semibold text-slate-400">ACTION</span>
                    <span className="font-bold text-slate-800 dark:text-slate-200 uppercase tracking-wide text-sm">{selectedLog.action}</span>
                  </div>
                  <div>
                    <span className="block text-xs font-semibold text-slate-400">OUTCOME</span>
                    <Badge
                      className={
                        selectedLog.outcome === 'success'
                          ? 'bg-emerald-100 text-emerald-800 dark:bg-emerald-950/50 dark:text-emerald-400 border-0 font-bold uppercase mt-1'
                          : 'bg-rose-100 text-rose-800 dark:bg-rose-950/50 dark:text-rose-400 border-0 font-bold uppercase mt-1'
                      }
                    >
                      {selectedLog.outcome}
                    </Badge>
                  </div>
                </div>

                {/* Properties table */}
                <div className="space-y-3">
                  <h4 className="text-sm font-semibold text-slate-800 dark:text-slate-200 border-b border-slate-100 dark:border-slate-800 pb-2">Properties</h4>
                  <div className="grid grid-cols-3 gap-2 text-sm">
                    <span className="font-medium text-slate-400">Log ID</span>
                    <span className="col-span-2 font-mono text-slate-700 dark:text-slate-300 text-xs break-all">{selectedLog.id}</span>

                    <span className="font-medium text-slate-400">Timestamp</span>
                    <span className="col-span-2 text-slate-700 dark:text-slate-300">{formatDate(selectedLog.created_at)}</span>

                    <span className="font-medium text-slate-400">Category</span>
                    <span className="col-span-2 text-slate-700 dark:text-slate-300">{selectedLog.category}</span>

                    <span className="font-medium text-slate-400">Resource</span>
                    <span className="col-span-2 text-slate-700 dark:text-slate-300">{selectedLog.resource} (ID: <span className="font-mono text-xs bg-slate-100 dark:bg-slate-800 dark:text-slate-200 px-1 rounded">{selectedLog.resource_id}</span>)</span>

                    <span className="font-medium text-slate-400">Actor ID</span>
                    <span className="col-span-2 font-mono text-slate-700 dark:text-slate-300 text-xs break-all">{selectedLog.actor_id}</span>

                    <span className="font-medium text-slate-400">Request ID</span>
                    <span className="col-span-2 font-mono text-slate-700 dark:text-slate-300 text-xs">{selectedLog.request_id}</span>

                    <span className="font-medium text-slate-400">Client IP</span>
                    <span className="col-span-2 font-mono text-slate-700 dark:text-slate-300 text-xs">{selectedLog.client_ip}</span>

                    <span className="font-medium text-slate-400">User Agent</span>
                    <span className="col-span-2 text-slate-600 dark:text-slate-400 text-xs break-all">{selectedLog.metadata?.user_agent || 'N/A'}</span>
                  </div>
                </div>

                {/* Metadata JSON tree */}
                <div className="space-y-2">
                  <h4 className="text-sm font-semibold text-slate-800 dark:text-slate-200">Metadata Payload</h4>
                  <div className="rounded-md bg-slate-900 p-4 text-xs font-mono text-indigo-250 dark:text-indigo-200 overflow-x-auto max-h-[300px]">
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
