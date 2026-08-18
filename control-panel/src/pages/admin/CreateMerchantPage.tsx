import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
// Removed Card component imports since sections are now borderless and backgroundless
import { Loader2, Store } from 'lucide-react';
import { fetchApi } from '@/lib/api';

const schema = z.object({
  name: z.string().min(2, 'Name must be at least 2 characters'),
  username: z.string().min(3, 'Username must be at least 3 characters'),
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
      const response = await fetchApi('/staff', {
        method: 'POST',
        body: JSON.stringify(data),
      });
      setSuccess(response?.message || 'Staff created successfully!');
      reset();
    } catch (err: any) {
      setError(err.message || 'Failed to create staff');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-2xl mx-auto py-10 animate-in fade-in duration-300">
      <div className="mb-8 flex items-center space-x-4">
        <div className="h-12 w-12 bg-primary/10 rounded-2xl flex items-center justify-center">
          <Store className="h-6 w-6 text-primary" />
        </div>
        <div>
          <h2 className="text-2xl font-bold font-display tracking-tight text-foreground">Create Merchant (Staff)</h2>
          <p className="text-muted-foreground text-sm">Register a new merchant staff entity in the system.</p>
        </div>
      </div>

      <div className="space-y-6">
        <div className="pb-4 border-b border-border/60 mb-6">
          <h3 className="font-bold font-display tracking-tight text-lg text-foreground">Merchant Details</h3>
          <p className="text-muted-foreground text-sm">Enter the basic information for the new merchant.</p>
        </div>
        <div>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <label className="text-sm font-semibold text-foreground">Merchant Name</label>
                <Input {...register('name')} placeholder="e.g. Chia Florist" className="rounded-xl border border-border bg-background" />
                {errors.name && <p className="text-sm text-rose-500">{errors.name.message}</p>}
              </div>

              <div className="space-y-2">
                <label className="text-sm font-semibold text-foreground">Username</label>
                <Input {...register('username')} placeholder="e.g. chia-florist" className="rounded-xl border border-border bg-background" />
                {errors.username && <p className="text-sm text-rose-500">{errors.username.message}</p>}
              </div>
            </div>

            <div className="space-y-2">
              <label className="text-sm font-semibold text-foreground">Description</label>
              <Textarea
                {...register('description')}
                placeholder="Briefly describe the business..."
                className="min-h-[100px] rounded-xl border border-border bg-background"
              />
              {errors.description && <p className="text-sm text-rose-500">{errors.description.message}</p>}
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <label className="text-sm font-semibold text-foreground">Logo URL</label>
                <Input {...register('logo_url')} placeholder="https://..." className="rounded-xl border border-border bg-background" />
                {errors.logo_url && <p className="text-sm text-rose-500">{errors.logo_url.message}</p>}
              </div>

              <div className="space-y-2">
                <label className="text-sm font-semibold text-foreground">Banner URL</label>
                <Input {...register('banner_url')} placeholder="https://..." className="rounded-xl border border-border bg-background" />
                {errors.banner_url && <p className="text-sm text-rose-500">{errors.banner_url.message}</p>}
              </div>
            </div>

            {error && <div className="p-3.5 text-sm text-rose-600 bg-rose-50 rounded-xl border border-rose-100">{error}</div>}
            {success && <div className="p-3.5 text-sm text-primary bg-primary/10 rounded-xl border border-primary/20">{success}</div>}

            <Button type="submit" disabled={loading} className="w-full bg-primary hover:bg-primary/90 text-primary-foreground rounded-xl font-semibold">
              {loading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Register Merchant
            </Button>
          </form>
        </div>
      </div>
    </div>
  );
}
