import { Wallet, Loader2, Plus } from 'lucide-react';
import { Button } from '../../../components/ui/button';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '../../../components/ui/table';
// Removed Card component imports since sections are now borderless and backgroundless
import { usePaymentsViewModel } from '../../../viewmodels/usePaymentsViewModel';
import { Link } from 'react-router-dom';

export default function PaymentAccountsPage() {
  const { accounts, methods, loading, error } = usePaymentsViewModel();

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
            <h2 className="text-3xl font-bold tracking-tight">Payment Accounts</h2>
            <p className="text-muted-foreground">
              Manage your merchant payment accounts and receiving channels.
            </p>
          </div>
          <div className="flex items-center space-x-2">
            <Button asChild>
              <Link to="/admin/payments/accounts/create">
                <Plus className="mr-2 h-4 w-4" /> Add Account
              </Link>
            </Button>
          </div>
        </div>

        <div className="space-y-6">
          <div className="pb-4 border-b border-border/60 mb-6">
            <h3 className="font-bold text-lg">Configured Accounts</h3>
            <p className="text-muted-foreground text-sm">
              These accounts are used to receive settlements from customers.
            </p>
          </div>
          <div>
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
          </div>
        </div>
      </div>
    </div>
  );
}
