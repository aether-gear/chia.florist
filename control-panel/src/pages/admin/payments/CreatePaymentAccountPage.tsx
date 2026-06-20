import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import { ArrowLeft, Loader2 } from 'lucide-react';
import { Button } from '../../../components/ui/button';
import { Input } from '../../../components/ui/input';
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '../../../components/ui/form';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '../../../components/ui/card';
import { usePaymentsViewModel } from '../../../viewmodels/usePaymentsViewModel';

const formSchema = z.object({
  method_id: z.string().min(1, 'Payment method is required'),
  account_name: z.string().min(3, 'Account name must be at least 3 characters'),
  account_number: z.string().optional(),
  phone_number: z.string().optional(),
  qr_string: z.string().optional(),
});

export default function CreatePaymentAccountPage() {
  const navigate = useNavigate();
  const { methods, createPaymentAccount, loading: viewModelLoading } = usePaymentsViewModel();
  const [isSubmitting, setIsSubmitting] = useState(false);

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      method_id: '',
      account_name: '',
      account_number: '',
      phone_number: '',
      qr_string: '',
    },
  });

  const onSubmit = async (values: z.infer<typeof formSchema>) => {
    setIsSubmitting(true);
    const success = await createPaymentAccount(values);
    setIsSubmitting(false);
    
    if (success) {
      navigate('/admin/payments/accounts');
    }
  };

  if (viewModelLoading && methods.length === 0) {
    return (
      <div className="flex h-[50vh] items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  return (
    <div className="flex-col md:flex">
      <div className="flex-1 space-y-4 p-8 pt-6 max-w-3xl mx-auto w-full">
        <div className="flex items-center space-x-4">
          <Button variant="outline" size="icon" onClick={() => navigate(-1)}>
            <ArrowLeft className="h-4 w-4" />
          </Button>
          <div>
            <h2 className="text-3xl font-bold tracking-tight">Create Payment Account</h2>
            <p className="text-muted-foreground">
              Add a new bank account or e-wallet to receive payments.
            </p>
          </div>
        </div>

        <Card>
          <CardHeader>
            <CardTitle>Account Details</CardTitle>
            <CardDescription>
              Please enter the payment gateway integration or bank account details.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Form {...form}>
              <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
                
                <FormField
                  control={form.control}
                  name="method_id"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Payment Method</FormLabel>
                      <FormControl>
                        <select
                          className="flex h-10 w-full items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                          {...field}
                        >
                          <option value="" disabled>Select a payment method</option>
                          {methods.map(method => (
                            <option key={method.id} value={method.id}>
                              {method.name} ({method.type.replace('_', ' ')})
                            </option>
                          ))}
                        </select>
                      </FormControl>
                      <FormDescription>
                        The payment channel provider.
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="account_name"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Account Name</FormLabel>
                      <FormControl>
                        <Input placeholder="e.g. PT Chia Florist" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
                  <FormField
                    control={form.control}
                    name="account_number"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Account Number</FormLabel>
                        <FormControl>
                          <Input placeholder="Required for Bank Transfers" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name="phone_number"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Phone Number</FormLabel>
                        <FormControl>
                          <Input placeholder="Required for E-Wallets" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>

                <Button type="submit" disabled={isSubmitting} className="w-full">
                  {isSubmitting ? (
                    <>
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                      Saving...
                    </>
                  ) : (
                    'Create Account'
                  )}
                </Button>
              </form>
            </Form>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
