import { CreditCard, Loader2 } from 'lucide-react';
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
import { usePaymentsViewModel } from '../../../viewmodels/usePaymentsViewModel';

export default function PaymentMethodsPage() {
  const { methods, loading, error } = usePaymentsViewModel();

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
      <div className="flex-1 space-y-4 p-8 pt-6">
        <div className="flex items-center justify-between space-y-2">
          <div>
            <h2 className="text-3xl font-bold tracking-tight">Payment Methods</h2>
            <p className="text-muted-foreground">
              View all supported payment methods on the platform.
            </p>
          </div>
        </div>

        <Card>
          <CardHeader>
            <CardTitle>Available Methods</CardTitle>
            <CardDescription>
              These are the payment channels available for processing customer payments.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="rounded-md border">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead className="w-[50px]"></TableHead>
                    <TableHead>Name</TableHead>
                    <TableHead>Type</TableHead>
                    <TableHead>Description</TableHead>
                    <TableHead className="text-right">Status</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {methods.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={5} className="text-center h-24">
                        No payment methods configured.
                      </TableCell>
                    </TableRow>
                  ) : (
                    methods.map((method) => (
                      <TableRow key={method.id}>
                        <TableCell>
                          <CreditCard className="h-5 w-5 text-muted-foreground" />
                        </TableCell>
                        <TableCell className="font-medium">{method.name}</TableCell>
                        <TableCell className="uppercase text-xs">{method.type.replace('_', ' ')}</TableCell>
                        <TableCell className="max-w-xs truncate text-muted-foreground">
                          {method.description}
                        </TableCell>
                        <TableCell className="text-right">
                          <Badge variant={method.is_active ? "default" : "secondary"}>
                            {method.is_active ? 'Active' : 'Inactive'}
                          </Badge>
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
