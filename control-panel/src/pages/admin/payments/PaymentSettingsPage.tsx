import { useState, useMemo } from 'react';
import { CreditCard, RefreshCw } from 'lucide-react';
import SearchInput from '../../../components/SearchInput';
import { Button } from '../../../components/ui/button';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '../../../components/ui/table';
import { Switch } from '../../../components/ui/switch';
import { Skeleton } from '../../../components/ui/skeleton';
import { usePaymentsViewModel } from '../../../viewmodels/usePaymentsViewModel';

export default function PaymentSettingsPage() {
  const { methods, loading, toggling, error, togglePaymentMethodActive, refetch } = usePaymentsViewModel();
  const [methodSearch, setMethodSearch] = useState('');

  const filteredMethods = useMemo(() => {
    if (!methods) return [];
    return methods.filter(method =>
      method.name.toLowerCase().includes(methodSearch.toLowerCase()) ||
      method.code.toLowerCase().includes(methodSearch.toLowerCase()) ||
      (method.description && method.description.toLowerCase().includes(methodSearch.toLowerCase()))
    );
  }, [methods, methodSearch]);

  return (
    <div className="flex-col md:flex">
      <div className="flex-1 space-y-12 p-6 sm:p-8 lg:p-12 animate-in fade-in duration-300">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <h2 className="text-3xl font-bold font-display tracking-tight text-foreground">Payment Settings</h2>
            <p className="text-muted-foreground text-sm">
              Manage supported payment methods.
            </p>
          </div>
        </div>

        <div className="space-y-6">
          <div className="pb-4 border-b border-border/60 mb-6">
            <h3 className="font-bold font-display tracking-tight text-lg text-foreground">Available Methods</h3>
            <p className="text-muted-foreground text-sm">
              These are the payment channels available for processing customer payments.
            </p>
          </div>
          <div>
            <div className="mb-4 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
              {/* Left Side: Filter and Search */}
              <div className="flex flex-col sm:flex-row items-center gap-4 w-full sm:w-auto">
                <SearchInput
                  value={methodSearch}
                  onChange={setMethodSearch}
                  placeholder="Search methods..."
                />
              </div>

              {/* Right Side: Refresh */}
              <div className="flex items-center gap-2 justify-end w-full sm:w-auto">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => refetch()}
                  disabled={loading || toggling}
                  className="flex items-center gap-1.5 border-border text-foreground hover:text-primary hover:bg-primary/5 rounded-xl transition-colors"
                >
                  <RefreshCw className={`h-4 w-4 ${loading || toggling ? 'animate-spin' : ''}`} />
                  Refresh
                </Button>
              </div>
            </div>

            <div className="rounded-md border">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead className="w-[50px]"></TableHead>
                    <TableHead>Name</TableHead>
                    <TableHead>Code</TableHead>
                    <TableHead>Provider</TableHead>
                    <TableHead>Type</TableHead>
                    <TableHead>Description</TableHead>
                    <TableHead className="text-right">Active</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {loading ? (
                    Array.from({ length: 3 }).map((_, i) => (
                      <TableRow key={`methods-skeleton-${i}`}>
                        <TableCell><Skeleton className="h-5 w-5 rounded bg-muted animate-pulse" /></TableCell>
                        <TableCell><Skeleton className="h-5 w-32 bg-muted animate-pulse" /></TableCell>
                        <TableCell><Skeleton className="h-5 w-16 bg-muted animate-pulse" /></TableCell>
                        <TableCell><Skeleton className="h-5 w-20 bg-muted animate-pulse" /></TableCell>
                        <TableCell><Skeleton className="h-5 w-24 bg-muted animate-pulse" /></TableCell>
                        <TableCell><Skeleton className="h-5 w-40 bg-muted animate-pulse" /></TableCell>
                        <TableCell className="text-right"><Skeleton className="h-5 w-16 ml-auto bg-muted animate-pulse" /></TableCell>
                      </TableRow>
                    ))
                  ) : error ? (
                    <TableRow>
                      <TableCell colSpan={7} className="text-center h-24 text-destructive">
                        {error}
                      </TableCell>
                    </TableRow>
                  ) : filteredMethods.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={7} className="text-center h-24 text-muted-foreground">
                        {methodSearch ? `No methods match "${methodSearch}"` : "No payment methods configured."}
                      </TableCell>
                    </TableRow>
                  ) : (
                    filteredMethods.map((method) => (
                      <TableRow key={method.id}>
                        <TableCell>
                          <CreditCard className="h-5 w-5 text-muted-foreground" />
                        </TableCell>
                        <TableCell className="font-medium">{method.name}</TableCell>
                        <TableCell className="font-mono text-xs">{method.code}</TableCell>
                        <TableCell className="capitalize text-xs font-semibold text-indigo-600 dark:text-indigo-400">{method.provider}</TableCell>
                        <TableCell className="uppercase text-xs">{method.type.replace('_', ' ')}</TableCell>
                        <TableCell className="max-w-xs truncate text-muted-foreground">
                          {method.description}
                        </TableCell>
                        <TableCell className="text-right">
                          <Switch
                            checked={method.is_active}
                            disabled={loading || toggling}
                            onCheckedChange={(checked) => togglePaymentMethodActive(method.id, checked)}
                            className="ml-auto data-[state=checked]:bg-primary"
                          />
                        </TableCell>
                      </TableRow>
                    ))
                  )}
                </TableBody>
              </Table>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
