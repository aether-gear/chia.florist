import { useState } from 'react';
import { CreditCard, Wallet, Loader2, Plus, Pencil, Eye } from 'lucide-react';
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
import { usePaymentsViewModel } from '../../../viewmodels/usePaymentsViewModel';
import { Link } from 'react-router-dom';
import type { PaymentMethod } from '../../../models/Payment';
import PaymentMethodFormSheet from '../../../components/payments/PaymentMethodFormSheet';
import PaymentMethodDetailOverlay from '../../../components/payments/PaymentMethodDetailOverlay';

export default function PaymentSettingsPage() {
  const { methods, accounts, loading, error, savePaymentMethod } = usePaymentsViewModel();
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [editingMethod, setEditingMethod] = useState<PaymentMethod | null>(null);
  const [isDetailOpen, setIsDetailOpen] = useState(false);
  const [detailMethod, setDetailMethod] = useState<PaymentMethod | null>(null);

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
              <CardHeader className="flex flex-row items-center justify-between">
                <div>
                  <CardTitle>Configured Accounts</CardTitle>
                  <CardDescription>
                    These accounts are used to receive settlements from customers.
                  </CardDescription>
                </div>
                <Button asChild size="sm">
                  <Link to="/admin/payments/accounts/create">
                    <Plus className="mr-2 h-4 w-4" /> Add Account
                  </Link>
                </Button>
              </CardHeader>
              <CardContent>
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
                      {accounts.length === 0 ? (
                        <TableRow>
                          <TableCell colSpan={4} className="text-center h-24">
                            No payment accounts configured.
                          </TableCell>
                        </TableRow>
                      ) : (
                        accounts.map((account) => (
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
              <CardHeader className="flex flex-row items-center justify-between">
                <div>
                  <CardTitle>Available Methods</CardTitle>
                  <CardDescription>
                    These are the payment channels available for processing customer payments.
                  </CardDescription>
                </div>
                <Button onClick={() => { setEditingMethod(null); setIsFormOpen(true); }} size="sm">
                  <Plus className="mr-2 h-4 w-4" /> Add Method
                </Button>
              </CardHeader>
              <CardContent>
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
                      {methods.length === 0 ? (
                        <TableRow>
                          <TableCell colSpan={8} className="text-center h-24">
                            No payment methods configured.
                          </TableCell>
                        </TableRow>
                      ) : (
                        methods.map((method) => (
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
