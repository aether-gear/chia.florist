import { useState, useMemo } from 'react';
import { Store, RefreshCw } from 'lucide-react';
import { Button } from '../../components/ui/button';
import { Link } from 'react-router-dom';
import SearchInput from '../../components/SearchInput';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '../../components/ui/table';
// Removed Card component imports since sections are now borderless and backgroundless
import { useMerchantsViewModel } from '../../viewmodels/useMerchantsViewModel';
import Pagination from '../../components/Pagination';
import { Skeleton } from '../../components/ui/skeleton';

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
          <div className="pb-4 border-b border-border/60 mb-6">
            <h3 className="font-bold font-display tracking-tight text-lg text-foreground">All Merchants</h3>
            <p className="text-muted-foreground text-sm">
              Showing {data?.merchants.length || 0} of {data?.total || 0} merchants.
            </p>
          </div>
          <div>
            <div className="mb-4 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
              {/* Left Side: Filter and Search */}
              <div className="flex flex-col sm:flex-row items-center gap-4 w-full sm:w-auto">
                <SearchInput
                  value={searchQuery}
                  onChange={setSearchQuery}
                  placeholder="Search merchants..."
                />
              </div>

              {/* Right Side: Adding and Refresh */}
              <div className="flex items-center gap-2 justify-end w-full sm:w-auto">
                <Button
                  variant="outline"
                  onClick={() => refresh()}
                  disabled={loading}
                  className="flex items-center gap-1.5 border-border text-foreground hover:text-primary hover:bg-primary/5 rounded-xl transition-colors animate-in fade-in duration-200"
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

            <div className="rounded-2xl border border-border overflow-hidden">
              <Table>
                <TableHeader className="bg-muted/50">
                  <TableRow>
                    <TableHead className="w-[80px]">Logo</TableHead>
                    <TableHead>Merchant Name</TableHead>
                    <TableHead>Description</TableHead>
                    <TableHead>Registered On</TableHead>
                    <TableHead className="text-right">Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {loading ? (
                    Array.from({ length: 5 }).map((_, i) => (
                      <TableRow key={`skeleton-${i}`}>
                        <TableCell><Skeleton className="h-10 w-10 rounded-md bg-muted animate-pulse" /></TableCell>
                        <TableCell>
                          <Skeleton className="h-5 w-36 animate-pulse bg-muted mb-1.5" />
                          <Skeleton className="h-3.5 w-24 animate-pulse bg-muted" />
                        </TableCell>
                        <TableCell><Skeleton className="h-5 w-48 animate-pulse bg-muted" /></TableCell>
                        <TableCell><Skeleton className="h-5 w-28 animate-pulse bg-muted" /></TableCell>
                        <TableCell className="text-right"><Skeleton className="h-8 w-24 ml-auto animate-pulse bg-muted" /></TableCell>
                      </TableRow>
                    ))
                  ) : error ? (
                    <TableRow>
                      <TableCell colSpan={5} className="h-24 text-center text-destructive">
                        {error}
                      </TableCell>
                    </TableRow>
                  ) : filteredMerchants.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={5} className="h-24 text-center text-muted-foreground">
                        {searchQuery ? `No merchants match "${searchQuery}"` : "No merchants found."}
                      </TableCell>
                    </TableRow>
                  ) : (
                    filteredMerchants.map((merchant) => (
                      <TableRow key={merchant.id}>
                        <TableCell>
                          <div className="h-10 w-10 overflow-hidden rounded-md border">
                            {merchant.logo_url ? (
                              <img
                                src={merchant.logo_url}
                                alt={merchant.name}
                                className="h-full w-full object-cover"
                              />
                            ) : (
                              <div className="flex h-full w-full items-center justify-center bg-muted">
                                <Store className="h-4 w-4 text-muted-foreground" />
                              </div>
                            )}
                          </div>
                        </TableCell>
                        <TableCell className="font-medium">
                          {merchant.name}
                          <div className="text-xs text-muted-foreground mt-1">
                            {merchant.id}
                          </div>
                        </TableCell>
                        <TableCell className="max-w-xs truncate">
                          {merchant.description || '-'}
                        </TableCell>
                        <TableCell>
                          {new Date(merchant.created_at).toLocaleDateString('en-GB', {
                            day: 'numeric',
                            month: 'short',
                            year: 'numeric'
                          })}
                        </TableCell>
                        <TableCell className="text-right">
                          <Button size="sm" variant="outline" asChild>
                            <Link to={`/admin/merchants/${merchant.id}/accounts/add`}>
                              Add Account
                            </Link>
                          </Button>
                        </TableCell>
                      </TableRow>
                    ))
                  )}
                </TableBody>
              </Table>
            </div>

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
