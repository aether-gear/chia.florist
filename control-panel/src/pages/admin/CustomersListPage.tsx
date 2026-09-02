import { useState, useMemo } from 'react';
import { RefreshCw } from 'lucide-react';
import { useCustomersViewModel } from '../../viewmodels/useCustomersViewModel';
import { Avatar, AvatarFallback } from '../../components/ui/avatar';
import Pagination from '../../components/Pagination';
import { Skeleton } from '../../components/ui/skeleton';
import SearchInput from '../../components/SearchInput';
import { Button } from '../../components/ui/button';
import { DataCard, DataCardGridHeader, DataCardList } from '../../components/DataCard';
import EmptyState from '../../components/EmptyState';

export default function CustomersListPage() {
  const { data, loading, error, page, limit, setPage, refresh } = useCustomersViewModel();

  const [searchQuery, setSearchQuery] = useState('');

  const filteredCustomers = useMemo(() => {
    if (!data?.customers) return [];
    return data.customers.filter(customer =>
      customer.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      customer.username.toLowerCase().includes(searchQuery.toLowerCase()) ||
      (customer.phone && customer.phone.toLowerCase().includes(searchQuery.toLowerCase()))
    );
  }, [data, searchQuery]);

  // Generate initials for avatar fallback
  const getInitials = (name: string) => {
    return name
      .split(' ')
      .map(n => n[0])
      .join('')
      .substring(0, 2)
      .toUpperCase();
  };

  return (
    <div className="flex-col md:flex">
      <div className="flex-1 space-y-8 p-6 sm:p-8 lg:p-12 animate-in fade-in slide-in-from-left-4 duration-300">
        <div className="flex items-center justify-between space-y-2">
          <div>
            <h2 className="text-3xl font-bold font-display tracking-tight text-foreground">Customers</h2>
            <p className="text-muted-foreground text-sm">
              Manage registered customer accounts
            </p>
          </div>
        </div>

        <div className="space-y-6">
          <div className="pb-4 border-b border-border/60">
            <h3 className="text-xl font-bold font-display tracking-tight text-foreground">All Customers</h3>
            <p className="text-muted-foreground text-sm">Showing {data?.customers.length || 0} of {data?.total || 0} customers.</p>
          </div>

          <div className="flex flex-col sm:flex-row items-center justify-between gap-4 w-full">
            <SearchInput
              value={searchQuery}
              onChange={setSearchQuery}
              placeholder="Search customers..."
              className="relative flex-1 max-w-sm w-full"
            />
            <div className="flex items-center gap-2 justify-end w-full sm:w-auto">
              <Button
                variant="outline"
                onClick={() => refresh()}
                disabled={loading}
                className="flex items-center gap-1.5 border-border text-foreground hover:text-primary hover:bg-primary/5 rounded-xl transition-colors"
              >
                <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
                Refresh
              </Button>
            </div>
          </div>

          <div className="flex items-center justify-between px-1 text-xs text-muted-foreground">
            <span>Found {filteredCustomers.length} customers</span>
          </div>

          {/* Content */}
          <div className="flex flex-col gap-2">
            <DataCardGridHeader>
              <span className="col-span-5">Customer</span>
              <span className="col-span-3">Phone</span>
              <span className="col-span-4">Last Login</span>
            </DataCardGridHeader>

            <DataCardList>
              {loading ? (
                Array.from({ length: 4 }).map((_, i) => (
                  <DataCard key={`skeleton-${i}`}>
                    <div className="col-span-5 flex items-center gap-3">
                      <Skeleton className="h-9 w-9 rounded-full bg-muted animate-pulse" />
                      <Skeleton className="h-4 w-32 bg-muted animate-pulse" />
                    </div>
                    <div className="col-span-3"><Skeleton className="h-4 w-24 bg-muted animate-pulse" /></div>
                    <div className="col-span-4"><Skeleton className="h-4 w-32 bg-muted animate-pulse" /></div>
                  </DataCard>
                ))
              ) : error ? (
                <EmptyState title="Failed to load customers" description={error} className="py-12 border-0 bg-transparent text-destructive" />
              ) : filteredCustomers.length === 0 ? (
                <EmptyState title="No customers found" description="No customers match your search criteria." className="py-12 border border-dashed border-border/80 rounded-2xl bg-zinc-50/10" />
              ) : (
                filteredCustomers.map((customer) => (
                  <DataCard key={customer.id}>
                    <div className="col-span-1 md:col-span-5 flex items-center gap-3 min-w-0">
                      <Avatar className="h-9 w-9">
                        <AvatarFallback className="bg-primary/10 text-primary font-bold text-xs">{getInitials(customer.name)}</AvatarFallback>
                      </Avatar>
                      <div className="min-w-0">
                        <h4 className="font-semibold font-display text-sm text-foreground truncate">{customer.name}</h4>
                        <p className="text-xs text-muted-foreground truncate">@{customer.username}</p>
                      </div>
                    </div>

                    <div className="col-span-1 md:col-span-3 text-xs text-muted-foreground">
                      <span className="md:hidden font-sans text-muted-foreground mr-1">Phone:</span>
                      {customer.phone || '-'}
                    </div>

                    <div className="col-span-1 md:col-span-4 text-xs text-muted-foreground">
                      <span className="md:hidden font-sans text-muted-foreground mr-1">Last login:</span>
                      {customer.last_login_at ? new Date(customer.last_login_at).toLocaleString('en-GB', { day: 'numeric', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' }) : <span className="italic opacity-70">Never logged in</span>}
                    </div>
                  </DataCard>
                ))
              )}
            </DataCardList>

            <Pagination
              currentPage={page}
              totalPages={Math.ceil((data?.total || 0) / limit)}
              totalItems={data?.total || 0}
              limit={limit}
              onPageChange={setPage}
              itemNamePlural="customers"
            />
          </div>
        </div>
      </div>
    </div>
  );
}
