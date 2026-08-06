import { useState, useMemo } from 'react';
import { Store, RefreshCw } from 'lucide-react';
import { Button } from '../../components/ui/button';
import { Link } from 'react-router-dom';
import SearchInput from '../../components/SearchInput';
import { useMerchantsViewModel } from '../../viewmodels/useMerchantsViewModel';
import Pagination from '../../components/Pagination';
import { Skeleton } from '../../components/ui/skeleton';
import { DataCard, DataCardGridHeader, DataCardList } from '../../components/DataCard';

export default function MerchantsListPage() {
  const { data, loading, error, page, limit, setPage, refresh } = useMerchantsViewModel();

  const [searchQuery, setSearchQuery] = useState('');

  const filteredMerchants = useMemo(() => {
    if (!data?.merchants) return [];
    return data.merchants.filter(merchant =>
      merchant.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      (merchant.description && merchant.description.toLowerCase().includes(searchQuery.toLowerCase())) ||
      merchant.id.toLowerCase().includes(searchQuery.toLowerCase())
    );
  }, [data, searchQuery]);

  return (
    <div className="flex-col md:flex">
      <div className="flex-1 space-y-8 p-6 sm:p-8 lg:p-12 animate-in fade-in duration-300">
        <div className="flex items-center justify-between space-y-2">
          <div>
            <h2 className="text-3xl font-bold font-display tracking-tight text-foreground">Merchants</h2>
            <p className="text-muted-foreground text-sm">
              Manage merchants registered on the platform
            </p>
          </div>
        </div>

        <div className="space-y-6">
          <div className="pb-4 border-b border-border/60">
            <h3 className="text-xl font-bold font-display tracking-tight text-foreground">All Merchants</h3>
            <p className="text-muted-foreground text-sm">Showing {data?.merchants.length || 0} of {data?.total || 0} merchants.</p>
          </div>

          <div className="flex flex-col sm:flex-row items-center justify-between gap-4 w-full">
            <SearchInput
              value={searchQuery}
              onChange={setSearchQuery}
              placeholder="Search merchants..."
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
              <Button asChild className="bg-primary hover:bg-primary/90 text-primary-foreground rounded-xl">
                <Link to="/admin/merchants/create">
                  <Store className="mr-2 h-4 w-4" /> Create Merchant
                </Link>
              </Button>
            </div>
          </div>

          <div className="flex items-center justify-between px-1 text-xs text-muted-foreground">
            <span>Found {filteredMerchants.length} merchants</span>
          </div>

          {/* Content */}
          <div className="flex flex-col gap-2">
            <DataCardGridHeader>
              <span className="col-span-4">Merchant</span>
              <span className="col-span-4">Description</span>
              <span className="col-span-2">Registered Date</span>
              <span className="col-span-2 text-right">Action</span>
            </DataCardGridHeader>

            <DataCardList>
              {loading ? (
                Array.from({ length: 4 }).map((_, i) => (
                  <DataCard key={`skeleton-${i}`}>
                    <div className="col-span-4 flex items-center gap-3">
                      <Skeleton className="h-9 w-9 rounded-md bg-muted animate-pulse" />
                      <Skeleton className="h-4 w-32 bg-muted animate-pulse" />
                    </div>
                    <div className="col-span-4"><Skeleton className="h-4 w-48 bg-muted animate-pulse" /></div>
                    <div className="col-span-2"><Skeleton className="h-4 w-24 bg-muted animate-pulse" /></div>
                    <div className="col-span-2 text-right"><Skeleton className="h-8 w-24 ml-auto bg-muted animate-pulse rounded-xl" /></div>
                  </DataCard>
                ))
              ) : error ? (
                <EmptyState title="Failed to load merchants" description={error} className="py-12 border-0 bg-transparent text-destructive" />
              ) : filteredMerchants.length === 0 ? (
                <EmptyState icon={<Store className="h-8 w-8 text-slate-400 mb-2 mx-auto" />} title="No merchants found" description="No merchants match your search criteria." className="py-12 border border-dashed border-border/80 rounded-2xl bg-zinc-50/10" />
              ) : (
                filteredMerchants.map((merchant) => (
                  <DataCard key={merchant.id}>
                    <div className="col-span-1 md:col-span-4 flex items-center gap-3 min-w-0">
                      <div className="h-9 w-9 overflow-hidden rounded-md border shrink-0 bg-muted flex items-center justify-center">
                        {merchant.logo_url ? (
                          <img src={merchant.logo_url} alt={merchant.name} className="h-full w-full object-cover" />
                        ) : (
                          <Store className="h-4 w-4 text-muted-foreground" />
                        )}
                      </div>
                      <div className="min-w-0">
                        <h4 className="font-semibold font-display text-sm text-foreground truncate">{merchant.name}</h4>
                        <p className="text-xs text-muted-foreground font-mono truncate">{merchant.id}</p>
                      </div>
                    </div>

                    <div className="col-span-1 md:col-span-4 text-xs text-muted-foreground truncate">
                      {merchant.description || 'No description provided.'}
                    </div>

                    <div className="col-span-1 md:col-span-2 text-xs text-muted-foreground">
                      <span className="md:hidden font-sans text-muted-foreground mr-1">Registered:</span>
                      {new Date(merchant.created_at).toLocaleDateString('en-GB', { day: 'numeric', month: 'short', year: 'numeric' })}
                    </div>

                    <div className="col-span-1 md:col-span-2 text-right" onClick={(e) => e.stopPropagation()}>
                      <Button size="sm" variant="outline" className="rounded-xl text-xs" asChild>
                        <Link to={`/admin/merchants/${merchant.id}/accounts/add`}>
                          Add Account
                        </Link>
                      </Button>
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
              itemNamePlural="merchants"
            />
          </div>
        </div>
      </div>
    </div>
  );
}
