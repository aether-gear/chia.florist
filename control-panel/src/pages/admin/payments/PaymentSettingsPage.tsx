import { useState, useMemo } from 'react';
import { CreditCard, RefreshCw } from 'lucide-react';
import SearchInput from '../../../components/SearchInput';
import { Button } from '../../../components/ui/button';
import { Switch } from '../../../components/ui/switch';
import { Skeleton } from '../../../components/ui/skeleton';
import { usePaymentsViewModel } from '../../../viewmodels/usePaymentsViewModel';

import { DataCard, DataCardGridHeader, DataCardList } from '../../../components/DataCard';
import EmptyState from '@/components/EmptyState';

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
          <div className="pb-4 border-b border-border/60">
            <h3 className="text-xl font-bold font-display tracking-tight text-foreground">Available Methods</h3>
            <p className="text-muted-foreground text-sm">These are the payment channels available for processing customer payments.</p>
          </div>

          <div className="flex flex-col sm:flex-row items-center justify-between gap-4 w-full">
            <SearchInput
              value={methodSearch}
              onChange={setMethodSearch}
              placeholder="Search methods..."
              className="relative flex-1 max-w-sm w-full"
            />
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

          <div className="flex items-center justify-between px-1 text-xs text-muted-foreground">
            <span>Found {filteredMethods.length} methods</span>
          </div>

          {/* Content */}
          <div className="flex flex-col gap-2">
            <DataCardGridHeader>
              <span className="col-span-4">Method</span>
              <span className="col-span-3">Provider & Type</span>
              <span className="col-span-3">Description</span>
              <span className="col-span-2 text-right">Active Status</span>
            </DataCardGridHeader>

            <DataCardList>
              {loading ? (
                Array.from({ length: 3 }).map((_, i) => (
                  <DataCard key={`skeleton-${i}`}>
                    <div className="col-span-4 flex items-center gap-3">
                      <Skeleton className="h-9 w-9 rounded-lg bg-muted animate-pulse" />
                      <Skeleton className="h-4 w-32 bg-muted animate-pulse" />
                    </div>
                    <div className="col-span-3"><Skeleton className="h-4 w-24 bg-muted animate-pulse" /></div>
                    <div className="col-span-3"><Skeleton className="h-4 w-40 bg-muted animate-pulse" /></div>
                    <div className="col-span-2 text-right"><Skeleton className="h-5 w-12 ml-auto bg-muted animate-pulse" /></div>
                  </DataCard>
                ))
              ) : error ? (
                <EmptyState title="Failed to load payment methods" description={error} className="py-12 border-0 bg-transparent text-destructive" />
              ) : filteredMethods.length === 0 ? (
                <EmptyState icon={<CreditCard className="h-8 w-8 text-slate-400 mb-2 mx-auto" />} title="No methods found" description="No payment methods configured matching your search." className="py-12 border border-dashed border-border/80 rounded-2xl bg-zinc-50/10" />
              ) : (
                filteredMethods.map((method) => (
                  <DataCard key={method.id}>
                    <div className="col-span-1 md:col-span-4 flex items-center gap-3 min-w-0">
                      <div className="h-9 w-9 rounded-lg bg-primary/10 flex items-center justify-center text-primary shrink-0">
                        <CreditCard className="h-4 w-4" />
                      </div>
                      <div className="min-w-0">
                        <h4 className="font-semibold font-display text-sm text-foreground truncate">{method.name}</h4>
                        <p className="text-xs font-mono text-muted-foreground truncate">{method.code}</p>
                      </div>
                    </div>

                    <div className="col-span-1 md:col-span-3 text-xs">
                      <span className="font-semibold text-primary capitalize">{method.provider}</span>
                      <span className="text-muted-foreground uppercase ml-1.5 font-sans">({method.type.replace('_', ' ')})</span>
                    </div>

                    <div className="col-span-1 md:col-span-3 text-xs text-muted-foreground truncate">
                      {method.description || '-'}
                    </div>

                    <div className="col-span-1 md:col-span-2 flex items-center justify-end gap-2" onClick={(e) => e.stopPropagation()}>
                      <span className="text-xs text-muted-foreground md:hidden mr-1">Active:</span>
                      <Switch
                        checked={method.is_active}
                        disabled={loading || toggling}
                        onCheckedChange={(checked) => togglePaymentMethodActive(method.id, checked)}
                        className="data-[state=checked]:bg-primary"
                      />
                    </div>
                  </DataCard>
                ))
              )}
            </DataCardList>
          </div>
        </div>
      </div>
    </div>
  );
}
