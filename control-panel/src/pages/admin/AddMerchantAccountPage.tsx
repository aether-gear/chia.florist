import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { useParams } from 'react-router-dom';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardHeader, CardContent, CardTitle, CardDescription } from '@/components/ui/card';
import { Loader2, Users } from 'lucide-react';
import { fetchApi } from '@/lib/api';

const schema = z.object({
  merchantId: z.string().min(1, 'Merchant ID is required'),
  email: z.string().email('Must be a valid email'),
  name: z.string().min(2, 'Name must be at least 2 characters'),
  username: z.string().min(3, 'Username must be at least 3 characters'),
  password: z.string().min(6, 'Password must be at least 6 characters'),
  phone: z.string().min(5, 'Phone number is required'),
});

type FormData = z.infer<typeof schema>;

export default function AddMerchantAccountPage() {
  const { merchantId: paramMerchantId } = useParams();

  
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  const { register, handleSubmit, formState: { errors }, reset } = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: {
      merchantId: paramMerchantId || '',
    }
  });

  const onSubmit = async (data: FormData) => {
    setLoading(true);
    setError(null);
    setSuccess(null);

    try {
      const response = await fetchApi(`/staff/${data.merchantId}/accounts`, {
        method: 'POST',
        body: JSON.stringify({
          email: data.email,
          name: data.name,
          username: data.username,
          password: data.password,
          phone: data.phone,
        }),
      });
      setSuccess(response?.message || 'Account successfully added!');
      reset({ merchantId: data.merchantId }); // Keep merchantId, reset the rest
    } catch (err: any) {
      setError(err.message || 'Failed to add merchant account');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-2xl mx-auto py-8">
      <div className="mb-8 flex items-center space-x-3">
        <div className="h-10 w-10 bg-indigo-100 rounded-full flex items-center justify-center">
          <Users className="h-5 w-5 text-indigo-600" />
        </div>
        <div>
          <h2 className="text-2xl font-bold tracking-tight">Add Merchant (Staff) Account</h2>
          <p className="text-muted-foreground">Register a new staff/admin account for a specific merchant staff entity.</p>
        </div>
      </div>

      <Card className="border-0 shadow-sm">
        <CardHeader>
          <CardTitle>Account Details</CardTitle>
          <CardDescription>Enter the credentials and contact info for the new user.</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
            <div className="space-y-2">
              <label className="text-sm font-medium">Merchant ID</label>
              <Input {...register('merchantId')} placeholder="Enter merchant ID" />
              {errors.merchantId && <p className="text-sm text-red-500">{errors.merchantId.message}</p>}
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <label className="text-sm font-medium">Full Name</label>
                <Input {...register('name')} placeholder="John Doe" />
                {errors.name && <p className="text-sm text-red-500">{errors.name.message}</p>}
              </div>

              <div className="space-y-2">
                <label className="text-sm font-medium">Username</label>
                <Input {...register('username')} placeholder="johndoe" />
                {errors.username && <p className="text-sm text-red-500">{errors.username.message}</p>}
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <label className="text-sm font-medium">Email Address</label>
                <Input type="email" {...register('email')} placeholder="john@example.com" />
                {errors.email && <p className="text-sm text-red-500">{errors.email.message}</p>}
              </div>

              <div className="space-y-2">
                <label className="text-sm font-medium">Phone Number</label>
                <Input {...register('phone')} placeholder="+62..." />
                {errors.phone && <p className="text-sm text-red-500">{errors.phone.message}</p>}
              </div>
            </div>

            <div className="space-y-2">
              <label className="text-sm font-medium">Password</label>
              <Input type="password" {...register('password')} placeholder="••••••••" />
              {errors.password && <p className="text-sm text-red-500">{errors.password.message}</p>}
            </div>

            {error && <div className="p-3 text-sm text-red-500 bg-red-50 rounded-md border border-red-100">{error}</div>}
            {success && <div className="p-3 text-sm text-green-700 bg-green-50 rounded-md border border-green-100">{success}</div>}

            <Button type="submit" disabled={loading} className="w-full bg-indigo-600 hover:bg-indigo-700">
              {loading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Create Account
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
