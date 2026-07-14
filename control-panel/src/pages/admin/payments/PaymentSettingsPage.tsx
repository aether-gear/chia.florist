import { useState, useMemo } from 'react';
import { CreditCard, Wallet, Plus, Pencil, RefreshCw } from 'lucide-react';
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
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '../../../components/ui/card';
import { Badge } from '../../../components/ui/badge';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '../../../components/ui/tabs';
import { Skeleton } from '../../../components/ui/skeleton';
import { usePaymentsViewModel } from '../../../viewmodels/usePaymentsViewModel';
import { Link } from 'react-router-dom';
import type { PaymentMethod } from '../../../models/Payment';
import PaymentMethodFormSheet from '../../../components/payments/PaymentMethodFormSheet';
import PaymentMethodDetailOverlay from '../../../components/payments/PaymentMethodDetailOverlay';

export default function PaymentSettingsPage() {
  const { methods, accounts, loading, error, savePaymentMethod, refetch } = usePaymentsViewModel();
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [editingMethod, setEditingMethod] = useState<PaymentMethod | null>(null);
  const [isDetailOpen, setIsDetailOpen] = useState(false);
  const [detailMethod, setDetailMethod] = useState<PaymentMethod | null>(null);

  const [accountSearch, setAccountSearch] = useState('');
  const [methodSearch, setMethodSearch] = useState('');

  const filteredAccounts = useMemo(() => {
    if (!accounts) return [];
    return accounts.filter(account =>
      account.account_name.toLowerCase().includes(accountSearch.toLowerCase()) ||
      (account.account_number && account.account_number.toLowerCase().includes(accountSearch.toLowerCase())) ||
      (account.phone_number && account.phone_number.toLowerCase().includes(accountSearch.toLowerCase()))
    );
  }, [accounts, accountSearch]);

  const filteredMethods = useMemo(() => {
    if (!methods) return [];
    return methods.filter(method =>
      method.name.toLowerCase().includes(methodSearch.toLowerCase()) ||
      method.code.toLowerCase().includes(methodSearch.toLowerCase()) ||
      (method.description && method.description.toLowerCase().includes(methodSearch.toLowerCase()))
    );
  }, [methods, methodSearch]);



  const getMethodName = (methodId: string) => {
    const method = methods.find(m => m.id === methodId);
    return method ? method.name : 'Unknown Method';
  };

  return (
    <div className="flex-col md:flex">
      <div className="flex-1 space-y-4 p-8 pt-6">
        <div className="flex items-center justify-between space-y-2">
          <div>
            <h2 className="text-3xl font-bold tracking-tight">Payment Settings</h2>
            <p className="text-muted-foreground">
              Manage supported payment methods and configured settlement accounts.
            </p>
          </div>
        </div>

        <Tabs defaultValue="accounts" className="space-y-4">
          <TabsList>
            <TabsTrigger value="accounts">Payment Accounts</TabsTrigger>
            <TabsTrigger value="methods">Payment Methods</TabsTrigger>
          </TabsList>

          <TabsContent value="accounts" className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle className="font-bold font-display tracking-tight text-lg text-foreground">Configured Accounts</CardTitle>
                <CardDescription className="text-muted-foreground text-sm">
                  These accounts are used to receive settlements from customers.
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="mb-4 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
                  {/* Left Side: Filter and Search */}
                  <div className="flex flex-col sm:flex-row items-center gap-4 w-full sm:w-auto">
                    <SearchInput
                      value={accountSearch}
                      onChange={setAccountSearch}
                      placeholder="Search accounts..."
                    />
                  </div>

                  {/* Right Side: Refresh & Add */}
                  <div className="flex items-center gap-2 justify-end w-full sm:w-auto">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => refetch()}
                      className="flex items-center gap-1.5 border-border text-foreground hover:text-primary hover:bg-primary/5 rounded-xl transition-colors"
                    >
                      <RefreshCw className="h-4 w-4" />
                      Refresh
                    </Button>
                    <Button asChild className="bg-primary hover:bg-primary/90 text-primary-foreground rounded-xl">
                      <Link to="/admin/payments/accounts/create">
                        <Plus className="mr-2 h-4 w-4" /> Add Account
                      </Link>
                    </Button>
                  </div>
                </div>

                <div className="rounded-md border">
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead className="w-[50px]"></TableHead>
                        <TableHead>Account Name</TableHead>
                        <TableHead>Payment Method</TableHead>
                        <TableHead>Account No / Phone</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {loading ? (
                        Array.from({ length: 3 }).map((_, i) => (
                          <TableRow key={`accounts-skeleton-${i}`}>
                            <TableCell><Skeleton className="h-5 w-5 rounded bg-muted animate-pulse" /></TableCell>
                            <TableCell><Skeleton className="h-5 w-40 bg-muted animate-pulse" /></TableCell>
                            <TableCell><Skeleton className="h-5 w-28 bg-muted animate-pulse" /></TableCell>
                            <TableCell><Skeleton className="h-5 w-32 bg-muted animate-pulse" /></TableCell>
                          </TableRow>
                        ))
                      ) : error ? (
                        <TableRow>
                          <TableCell colSpan={4} className="text-center h-24 text-destructive">
                            {error}
                          </TableCell>
                        </TableRow>
                      ) : filteredAccounts.length === 0 ? (
                        <TableRow>
                          <TableCell colSpan={4} className="text-center h-24 text-muted-foreground">
                            {accountSearch ? `No accounts match "${accountSearch}"` : "No payment accounts configured."}
                          </TableCell>
                        </TableRow>
                      ) : (
                        filteredAccounts.map((account) => (
                          <TableRow key={account.id}>
                            <TableCell>
                              <Wallet className="h-5 w-5 text-muted-foreground" />
                            </TableCell>
                            <TableCell className="font-medium">{account.account_name}</TableCell>
                            <TableCell>{getMethodName(account.method_id)}</TableCell>
                            <TableCell className="font-mono text-sm">
                              {account.account_number || account.phone_number || '-'}
                            </TableCell>
                          </TableRow>
                        ))
                      )}
                    </TableBody>
                  </Table>
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="methods" className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle className="font-bold font-display tracking-tight text-lg text-foreground">Available Methods</CardTitle>
                <CardDescription className="text-muted-foreground text-sm">
                  These are the payment channels available for processing customer payments.
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="mb-4 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
                  {/* Left Side: Filter and Search */}
                  <div className="flex flex-col sm:flex-row items-center gap-4 w-full sm:w-auto">
                    <SearchInput
                      value={methodSearch}
                      onChange={setMethodSearch}
                      placeholder="Search methods..."
                    />
                  </div>

                  {/* Right Side: Refresh & Add */}
                  <div className="flex items-center gap-2 justify-end w-full sm:w-auto">
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => refetch()}
                      className="flex items-center gap-1.5 border-border text-foreground hover:text-primary hover:bg-primary/5 rounded-xl transition-colors"
                    >
                      <RefreshCw className="h-4 w-4" />
                      Refresh
                    </Button>
                    <Button 
                      onClick={() => { setEditingMethod(null); setIsFormOpen(true); }}
                      className="bg-primary hover:bg-primary/90 text-primary-foreground rounded-xl"
                    >
                      <Plus className="mr-2 h-4 w-4" /> Add Method
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
                        <TableHead className="text-right">Status</TableHead>
                        <TableHead className="w-[100px] text-right">Actions</TableHead>
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
                            <TableCell className="text-right"><Skeleton className="h-8 w-8 rounded-xl ml-auto bg-muted animate-pulse" /></TableCell>
                          </TableRow>
                        ))
                      ) : error ? (
                        <TableRow>
                          <TableCell colSpan={8} className="text-center h-24 text-destructive">
                            {error}
                          </TableCell>
                        </TableRow>
                      ) : filteredMethods.length === 0 ? (
                        <TableRow>
                          <TableCell colSpan={8} className="text-center h-24 text-muted-foreground">
                            {methodSearch ? `No methods match "${methodSearch}"` : "No payment methods configured."}
                          </TableCell>
                        </TableRow>
                      ) : (
                        filteredMethods.map((method) => (
                          <TableRow key={method.id} className="cursor-pointer" onClick={() => {
                            setDetailMethod(method);
                            setIsDetailOpen(true);
                          }}>
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
                              <Badge variant={method.is_active ? "default" : "secondary"}>
                                {method.is_active ? 'Active' : 'Inactive'}
                              </Badge>
                            </TableCell>
                            <TableCell className="text-right space-x-1" onClick={(e) => e.stopPropagation()}>
                              <Button
                                variant="ghost"
                                size="icon"
                                className="h-8 w-8"
                                onClick={() => {
                                  setEditingMethod(method);
                                  setIsFormOpen(true);
                                }}
                              >
                                <Pencil className="h-4 w-4" />
                                <span className="sr-only">Edit</span>
                              </Button>
                            </TableCell>
                          </TableRow>
                        ))
                      )}
                    </TableBody>
                  </Table>
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>

        <PaymentMethodFormSheet
          isOpen={isFormOpen}
          onOpenChange={setIsFormOpen}
          method={editingMethod}
          onSave={savePaymentMethod}
        />

        <PaymentMethodDetailOverlay
          isOpen={isDetailOpen}
          onOpenChange={setIsDetailOpen}
          method={detailMethod}
        />
      </div>
    </div>
  );
}
