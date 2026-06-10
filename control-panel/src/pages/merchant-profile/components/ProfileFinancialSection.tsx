import { useFormContext } from 'react-hook-form';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';

export function ProfileFinancialSection() {
  const { register } = useFormContext();

  return (
    <Card className="border-0 shadow-sm">
      <CardHeader>
        <CardTitle>Financial Information</CardTitle>
        <CardDescription>Bank accounts and tax details for payouts.</CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="space-y-2">
            <label className="text-sm font-medium">Bank Name</label>
            <Input {...register('financial.bankName')} placeholder="e.g. BCA, Mandiri" />
          </div>
          <div className="space-y-2">
            <label className="text-sm font-medium">Bank Account Name</label>
            <Input {...register('financial.bankAccountName')} placeholder="Name on the account" />
          </div>
        </div>

        <div className="space-y-2 md:w-1/2">
          <label className="text-sm font-medium">Bank Account Number</label>
          <Input {...register('financial.bankAccountNumber')} placeholder="Account number" />
        </div>

        <div className="space-y-2">
          <label className="text-sm font-medium">E-Wallet Information (Optional)</label>
          <Input {...register('financial.eWalletInformation')} placeholder="e.g. OVO: 08123456789" />
        </div>

        <div className="pt-4 border-t border-slate-100">
          <div className="space-y-2 md:w-1/2">
            <label className="text-sm font-medium">Tax Number / NPWP (Optional)</label>
            <Input {...register('financial.taxNumber')} placeholder="00.000.000.0-000.000" />
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
