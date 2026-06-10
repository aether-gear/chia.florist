import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardHeader, CardContent, CardTitle, CardDescription } from '@/components/ui/card';
import { Loader2, Store } from 'lucide-react';
import { fetchApi } from '@/lib/api';

const schema = z.object({
  name: z.string().min(2, 'Name must be at least 2 characters'),
  description: z.string().min(10, 'Description must be at least 10 characters'),
  logo_url: z.string().url('Must be a valid URL'),
  banner_url: z.string().url('Must be a valid URL'),
});

type FormData = z.infer<typeof schema>;

export default function CreateMerchantPage() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  const { register, handleSubmit, formState: { errors }, reset } = useForm<FormData>({
    resolver: zodResolver(schema),
  });

  const onSubmit = async (data: FormData) => {
    setLoading(true);
    setError(null);
    setSuccess(null);

    try {
      const response = await fetchApi('/merchants', {
        method: 'POST',
        body: JSON.stringify(data),
      });
      setSuccess(response.message || 'Merchant created successfully!');
      reset();
    } catch (err: any) {
      setError(err.message || 'Failed to create merchant');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-2xl mx-auto py-8">
      <div className="mb-8 flex items-center space-x-3">
        <div className="h-10 w-10 bg-indigo-100 rounded-full flex items-center justify-center">
          <Store className="h-5 w-5 text-indigo-600" />
        </div>
        <div>
          <h2 className="text-2xl font-bold tracking-tight">Create Merchant</h2>
          <p className="text-muted-foreground">Register a new merchant entity in the system.</p>
        </div>
      </div>

      <Card className="border-0 shadow-sm">
        <CardHeader>
          <CardTitle>Merchant Details</CardTitle>
          <CardDescription>Enter the basic information for the new merchant.</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
            <div className="space-y-2">
              <label className="text-sm font-medium">Merchant Name</label>
              <Input {...register('name')} placeholder="e.g. Chia Florist" />
              {errors.name && <p className="text-sm text-red-500">{errors.name.message}</p>}
            </div>

            <div className="space-y-2">
              <label className="text-sm font-medium">Description</label>
              <textarea 
                {...register('description')}
                className="flex min-h-[100px] w-full rounded-md border border-slate-200 bg-white px-3 py-2 text-sm ring-offset-white placeholder:text-slate-500 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-indigo-500"
                placeholder="Briefly describe the business..."
              />
              {errors.description && <p className="text-sm text-red-500">{errors.description.message}</p>}
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <label className="text-sm font-medium">Logo URL</label>
                <Input {...register('logo_url')} placeholder="https://..." />
                {errors.logo_url && <p className="text-sm text-red-500">{errors.logo_url.message}</p>}
              </div>

              <div className="space-y-2">
                <label className="text-sm font-medium">Banner URL</label>
                <Input {...register('banner_url')} placeholder="https://..." />
                {errors.banner_url && <p className="text-sm text-red-500">{errors.banner_url.message}</p>}
              </div>
            </div>

            {error && <div className="p-3 text-sm text-red-500 bg-red-50 rounded-md border border-red-100">{error}</div>}
            {success && <div className="p-3 text-sm text-green-700 bg-green-50 rounded-md border border-green-100">{success}</div>}

            <Button type="submit" disabled={loading} className="w-full bg-indigo-600 hover:bg-indigo-700">
              {loading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Register Merchant
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
