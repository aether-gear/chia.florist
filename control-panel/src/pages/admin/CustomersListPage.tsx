import { Search, Loader2 } from 'lucide-react';
import { Input } from '../../components/ui/input';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '../../components/ui/table';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '../../components/ui/card';
import { useCustomersViewModel } from '../../viewmodels/useCustomersViewModel';
import { Avatar, AvatarFallback } from '../../components/ui/avatar';

export default function CustomersListPage() {
  const { data, loading, error } = useCustomersViewModel();

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

        <Card className="border-0 shadow-none bg-zinc-50/40 dark:bg-slate-900/40">
          <CardHeader>
            <CardTitle className="font-bold font-display tracking-tight text-lg text-foreground">All Customers</CardTitle>
            <CardDescription className="text-muted-foreground text-sm">
              Showing {data?.customers.length || 0} of {data?.total || 0} customers.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="mb-4 flex items-center gap-4">
              <div className="relative flex-1 max-w-sm">
                <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
                <Input
                  type="search"
                  placeholder="Search customers by name, username, or email..."
                  className="pl-8 rounded-xl border border-border bg-background text-foreground"
                />
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
                  {data?.customers.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={3} className="h-24 text-center">
                        No customers found.
                      </TableCell>
                    </TableRow>
                  ) : (
                    data?.customers.map((customer) => (
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
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
