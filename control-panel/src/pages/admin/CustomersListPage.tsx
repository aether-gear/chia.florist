import { useState, useMemo } from 'react';
import { RefreshCw } from 'lucide-react';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '../../components/ui/table';
// Removed Card component imports since sections are now borderless and backgroundless
import { useCustomersViewModel } from '../../viewmodels/useCustomersViewModel';
import { Avatar, AvatarFallback } from '../../components/ui/avatar';
import Pagination from '../../components/Pagination';
import { Skeleton } from '../../components/ui/skeleton';
import SearchInput from '../../components/SearchInput';
import { Button } from '../../components/ui/button';

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
      <div className="flex-1 space-y-8 p-6 sm:p-8 lg:p-12 animate-in fade-in duration-300">
        <div className="flex items-center justify-between space-y-2">
          <div>
            <h2 className="text-3xl font-bold font-display tracking-tight text-foreground">Customers</h2>
            <p className="text-muted-foreground text-sm">
              Manage registered customer accounts
            </p>
          </div>
        </div>

        <div className="space-y-6">
          <div className="pb-4 border-b border-border/60 mb-6">
            <h3 className="font-bold font-display tracking-tight text-lg text-foreground">All Customers</h3>
            <p className="text-muted-foreground text-sm">
              Showing {data?.customers.length || 0} of {data?.total || 0} customers.
            </p>
          </div>
          <div>
            <div className="mb-4 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
              {/* Left Side: Filter and Search */}
              <div className="flex flex-col sm:flex-row items-center gap-4 w-full sm:w-auto">
                <SearchInput
                  value={searchQuery}
                  onChange={setSearchQuery}
                  placeholder="Search customers..."
                />
              </div>

              {/* Right Side: Refresh */}
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

            <div className="rounded-2xl border border-border overflow-hidden">
              <Table>
                <TableHeader className="bg-muted/50">
                  <TableRow>
                    <TableHead className="w-[250px]">Customer</TableHead>
                    <TableHead>Phone</TableHead>
                    <TableHead>Last Login</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {loading ? (
                    Array.from({ length: 5 }).map((_, i) => (
                      <TableRow key={`skeleton-${i}`}>
                        <TableCell>
                          <div className="flex items-center gap-3">
                            <Skeleton className="h-10 w-10 rounded-full animate-pulse bg-muted" />
                            <div className="space-y-1.5 flex-1">
                              <Skeleton className="h-4 w-28 animate-pulse bg-muted" />
                              <Skeleton className="h-3 w-20 animate-pulse bg-muted" />
                            </div>
                          </div>
                        </TableCell>
                        <TableCell><Skeleton className="h-5 w-24 animate-pulse bg-muted" /></TableCell>
                        <TableCell><Skeleton className="h-5 w-32 animate-pulse bg-muted" /></TableCell>
                      </TableRow>
                    ))
                  ) : error ? (
                    <TableRow>
                      <TableCell colSpan={3} className="h-24 text-center text-destructive">
                        {error}
                      </TableCell>
                    </TableRow>
                  ) : filteredCustomers.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={3} className="h-24 text-center text-muted-foreground">
                        {searchQuery ? `No customers match "${searchQuery}"` : "No customers found."}
                      </TableCell>
                    </TableRow>
                  ) : (
                    filteredCustomers.map((customer) => (
                      <TableRow key={customer.id}>
                        <TableCell className="font-medium">
                          <div className="flex items-center gap-3">
                            <Avatar>
                              <AvatarFallback className="bg-primary/10 text-primary font-bold">{getInitials(customer.name)}</AvatarFallback>
                            </Avatar>
                            <div>
                              <div className="font-semibold text-foreground">{customer.name}</div>
                              <div className="text-xs text-muted-foreground">@{customer.username}</div>
                            </div>
                          </div>
                        </TableCell>
                        <TableCell>
                          {customer.phone || '-'}
                        </TableCell>
                        <TableCell>
                          {customer.last_login_at ? new Date(customer.last_login_at).toLocaleString('en-GB', {
                            day: 'numeric',
                            month: 'short',
                            year: 'numeric',
                            hour: '2-digit',
                            minute: '2-digit'
                          }) : (
                            <span className="text-muted-foreground text-sm italic">Never logged in</span>
                          )}
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
              itemNamePlural="customers"
            />
          </div>
        </div>
      </div>
    </div>
  );
}
