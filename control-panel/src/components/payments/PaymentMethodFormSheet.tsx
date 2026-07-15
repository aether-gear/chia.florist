import { useState, useEffect } from 'react';
import { useForm, useWatch } from 'react-hook-form';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import { Loader2, Info } from 'lucide-react';
import { Button } from '../ui/button';
import { Input } from '../ui/input';
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '../ui/form';
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from '../ui/sheet';
import {
  FEE_TYPES,
  METHOD_TYPES,
  METHOD_CODES,
  type PaymentMethod,
} from '../../models/Payment';

const formSchema = z.object({
  name: z.string().min(1, 'Name is required'),
  code: z.string().min(1, 'Code is required'),
  provider: z.string().min(1, 'Provider is required'),
  type: z.string().min(1, 'Type is required'),
  is_active: z.boolean(),
  description: z.string().min(1, 'Description is required'),
  fee_type: z.string().min(1, 'Fee type is required'),
  fee_amount: z.string().optional().refine((val) => {
    if (!val) return true;
    const num = parseInt(val, 10);
    return !isNaN(num) && num >= 0;
  }, 'Fee amount must be a non-negative integer'),
  fee_percentage: z.string().optional().refine((val) => {
    if (!val) return true;
    const num = parseFloat(val);
    return !isNaN(num) && num >= 0 && num <= 100;
  }, 'Percentage must be between 0 and 100 (e.g. 1.5 for 1.5%)'),
  instruction: z.string().optional(),
});

interface PaymentMethodFormSheetProps {
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
  method: PaymentMethod | null;
  onSave: (
    methodData: Omit<Partial<PaymentMethod>, 'fee_percentage'> & { fee_amount?: string; fee_percentage?: string },
    instructionContent?: string
  ) => Promise<boolean>;
}

export default function PaymentMethodFormSheet({
  isOpen,
  onOpenChange,
  method,
  onSave,
}: PaymentMethodFormSheetProps) {
  const isEdit = !!method;
  const [isSubmitting, setIsSubmitting] = useState(false);

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: '',
      code: 'gopay',
      provider: 'manual',
      type: 'bank_transfer',
      is_active: true,
      description: '',
      fee_type: 'flat',
      fee_amount: '0',
      fee_percentage: '0.0',
      instruction: '',
    },
  });

  const watchFeeType = useWatch({
    control: form.control,
    name: 'fee_type',
  });

  useEffect(() => {
    if (isOpen) {
      if (method) {
        form.reset({
          name: method.name,
          code: method.code,
          provider: method.provider || 'manual',
          type: method.type,
          is_active: method.is_active,
          description: method.description,
          fee_type: method.fee_type,
          fee_amount: method.fee_fixed.toString(),
          fee_percentage: method.fee_percentage ? (method.fee_percentage * 100).toString() : '0.0',
          instruction: method.instruction?.content || '',
        });
      } else {
        form.reset({
          name: '',
          code: 'gopay',
          provider: 'manual',
          type: 'bank_transfer',
          is_active: true,
          description: '',
          fee_type: 'flat',
          fee_amount: '0',
          fee_percentage: '0.0',
          instruction: '',
        });
      }
    }
  }, [isOpen, method, form]);

  const onSubmit = async (values: z.infer<typeof formSchema>) => {
    setIsSubmitting(true);
    const instructionContent = values.instruction?.trim() || undefined;

    const success = await onSave(
      {
        id: method?.id,
        name: values.name,
        code: values.code as any,
        provider: values.provider,
        type: values.type as any,
        is_active: values.is_active,
        description: values.description,
        fee_type: values.fee_type as any,
        fee_amount: values.fee_amount || '0',
        fee_percentage: values.fee_percentage ? (parseFloat(values.fee_percentage) / 100).toString() : '0',
      },
      instructionContent
    );

    setIsSubmitting(false);
    if (success) {
      onOpenChange(false);
    }
  };

  return (
    <Sheet open={isOpen} onOpenChange={onOpenChange}>
      <SheetContent side="right" className="w-full sm:max-w-xl overflow-y-auto border-l border-border/80 shadow-2xl">
        <SheetHeader className="border-b border-border/60 pb-4 mb-6">
          <SheetTitle className="text-2xl font-bold tracking-tight">
            {isEdit ? 'Edit Payment Method' : 'Create Payment Method'}
          </SheetTitle>
          <SheetDescription className="text-muted-foreground text-sm">
            {isEdit ? 'Update details and instruction for this payment method.' : 'Add a new payment channel and instruction to the platform.'}
          </SheetDescription>
        </SheetHeader>

        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6 pb-8">

            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <FormField
                control={form.control}
                name="name"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel className="text-xs uppercase tracking-wider text-muted-foreground font-semibold">Name</FormLabel>
                    <FormControl>
                      <Input placeholder="e.g. GoPay, Bank Mandiri" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="code"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel className="text-xs uppercase tracking-wider text-muted-foreground font-semibold">Code</FormLabel>
                    <FormControl>
                      <select
                        className="flex h-10 w-full items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                        {...field}
                      >
                        {METHOD_CODES.map((c) => (
                          <option key={c.value} value={c.value}>
                            {c.label}
                          </option>
                        ))}
                      </select>
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <FormField
                control={form.control}
                name="type"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel className="text-xs uppercase tracking-wider text-muted-foreground font-semibold">Type</FormLabel>
                    <FormControl>
                      <select
                        className="flex h-10 w-full items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                        {...field}
                      >
                        {METHOD_TYPES.map((t) => (
                          <option key={t.value} value={t.value}>
                            {t.label}
                          </option>
                        ))}
                      </select>
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="provider"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel className="text-xs uppercase tracking-wider text-muted-foreground font-semibold">Provider</FormLabel>
                    <FormControl>
                      <select
                        className="flex h-10 w-full items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                        {...field}
                      >
                        <option value="manual">Manual</option>
                        <option value="midtrans">Midtrans</option>
                      </select>
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            <div className="grid grid-cols-1 gap-4 sm:grid-cols-3 items-end">
              <div className="sm:col-span-2">
                <FormField
                  control={form.control}
                  name="description"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="text-xs uppercase tracking-wider text-muted-foreground font-semibold">Description</FormLabel>
                      <FormControl>
                        <Input placeholder="e.g. Pay securely using GoPay e-wallet" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>

              <FormField
                control={form.control}
                name="is_active"
                render={({ field }) => (
                  <FormItem className="flex flex-row items-center space-x-3 space-y-0 rounded-xl border p-3.5 shadow-sm h-10 bg-background border-input">
                    <FormControl>
                      <input
                        type="checkbox"
                        className="h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary"
                        checked={field.value}
                        onChange={field.onChange}
                      />
                    </FormControl>
                    <div className="space-y-0.5 leading-none">
                      <FormLabel className="font-semibold text-xs text-muted-foreground uppercase tracking-wider">Active</FormLabel>
                    </div>
                  </FormItem>
                )}
              />
            </div>

            <div className="border-t border-border/60 my-6 pt-6 space-y-6">
              <h3 className="text-md font-semibold tracking-tight text-slate-800 dark:text-slate-200">Fee Configuration</h3>

              <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
                <FormField
                  control={form.control}
                  name="fee_type"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel className="text-xs uppercase tracking-wider text-muted-foreground font-semibold">Fee Type</FormLabel>
                      <FormControl>
                        <select
                          className="flex h-10 w-full items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                          {...field}
                        >
                          {FEE_TYPES.map((f) => (
                            <option key={f.value} value={f.value}>
                              {f.label}
                            </option>
                          ))}
                        </select>
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                {(watchFeeType === 'flat' || watchFeeType === 'mixed') && (
                  <FormField
                    control={form.control}
                    name="fee_amount"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel className="text-xs uppercase tracking-wider text-muted-foreground font-semibold">Fixed Fee (IDR)</FormLabel>
                        <FormControl>
                          <Input type="number" placeholder="0" {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                )}

                {(watchFeeType === 'percentage' || watchFeeType === 'mixed') && (
                  <FormField
                    control={form.control}
                    name="fee_percentage"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel className="text-xs uppercase tracking-wider text-muted-foreground font-semibold">Percentage Fee (%)</FormLabel>
                        <FormControl>
                          <div className="relative flex items-center">
                            <Input type="number" step="0.01" placeholder="1.5" className="pr-8" {...field} />
                            <span className="absolute right-3 text-sm font-semibold text-muted-foreground">%</span>
                          </div>
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                )}
              </div>

              {(watchFeeType === 'percentage' || watchFeeType === 'mixed') && (
                <div className="flex items-start gap-2 text-xs text-amber-700 bg-amber-50 dark:bg-amber-950/20 dark:text-amber-400 p-3 rounded-lg border border-amber-200/50">
                  <Info className="h-4 w-4 shrink-0 mt-0.5" />
                  <div>
                    <strong>Note:</strong> Enter percentage as direct number (e.g. 1.5% &rarr; enter <strong>1.5</strong>).
                  </div>
                </div>
              )}
            </div>

            <div className="border-t border-border/60 my-6 pt-6 space-y-4">
              <h3 className="text-md font-semibold tracking-tight text-slate-800 dark:text-slate-200">Payment Instructions</h3>

              <FormField
                control={form.control}
                name="instruction"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel className="text-xs uppercase tracking-wider text-muted-foreground font-semibold">Instruction Content (Markdown)</FormLabel>
                    <FormControl>
                      <textarea
                        className="flex min-h-[120px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                        placeholder="Please transfer exactly **Rp {{amount}}** to Account **{{va_number}}**."
                        {...field}
                      />
                    </FormControl>
                    <FormDescription className="text-[11px] leading-normal text-muted-foreground">
                      Markdown formatting is supported. Use variables like <code>{'{{amount}}'}</code> or <code>{'{{va_number}}'}</code> to dynamically insert settlement values.
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            <div className="flex space-x-3 border-t border-border/60 pt-6">
              <Button type="button" variant="outline" className="w-1/3" onClick={() => onOpenChange(false)}>
                Cancel
              </Button>
              <Button type="submit" disabled={isSubmitting} className="w-2/3">
                {isSubmitting ? (
                  <>
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    Saving...
                  </>
                ) : (
                  'Save Payment Method'
                )}
              </Button>
            </div>
          </form>
        </Form>
      </SheetContent>
    </Sheet>
  );
}
