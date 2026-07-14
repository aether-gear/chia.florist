import { Store, Search, Loader2 } from 'lucide-react';
import { Input } from '../../components/ui/input';
import { Button } from '../../components/ui/button';
import { Link } from 'react-router-dom';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '../../components/ui/table';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '../../components/ui/card';
import { useMerchantsViewModel } from '../../viewmodels/useMerchantsViewModel';
import Pagination from '../../components/Pagination';

export default function MerchantsListPage() {
  const { data, loading, error, page, limit, setPage } = useMerchantsViewModel();

  if (loading) {
    return (
      <div className="flex h-[50vh] items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex h-[50vh] items-center justify-center">
        <p className="text-destructive">{error}</p>
      </div>
    );
  }

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
          <div className="flex items-center space-x-2">
            <Button asChild className="bg-primary hover:bg-primary/90 text-primary-foreground rounded-xl">
              <Link to="/admin/merchants/create">
                <Store className="mr-2 h-4 w-4" /> Create Merchant
              </Link>
            </Button>
          </div>
        </div>

        <Card className="border-0 shadow-none bg-zinc-50/40 dark:bg-slate-900/40">
          <CardHeader>
            <CardTitle className="font-bold font-display tracking-tight text-lg text-foreground">All Merchants</CardTitle>
            <CardDescription className="text-muted-foreground text-sm">
              Showing {data?.merchants.length || 0} of {data?.total || 0} merchants.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="mb-4 flex items-center gap-4">
              <div className="relative flex-1 max-w-sm">
                <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
                <Input
                  type="search"
                  placeholder="Search merchants..."
                  className="pl-8 rounded-xl border border-border bg-background text-foreground"
                />
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
                  {data?.merchants.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={5} className="h-24 text-center">
                        No merchants found.
                      </TableCell>
                    </TableRow>
                  ) : (
                    data?.merchants.map((merchant) => (
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
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
