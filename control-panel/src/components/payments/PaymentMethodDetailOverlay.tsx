import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription } from '../ui/dialog';
import { Badge } from '../ui/badge';
import { CreditCard, FileText, Landmark, ShieldCheck, Wallet } from 'lucide-react';
import type { PaymentMethod } from '../../models/Payment';

interface PaymentMethodDetailOverlayProps {
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
  method: PaymentMethod | null;
}

export default function PaymentMethodDetailOverlay({
  isOpen,
  onOpenChange,
  method,
}: PaymentMethodDetailOverlayProps) {
  if (!method) return null;

  const getMethodIcon = (type: string) => {
    switch (type) {
      case 'bank_transfer':
        return <Landmark className="h-6 w-6 text-indigo-500" />;
      case 'ewallet':
        return <Wallet className="h-6 w-6 text-emerald-500" />;
      case 'qr_code':
        return <CreditCard className="h-6 w-6 text-amber-500" />;
      default:
        return <CreditCard className="h-6 w-6 text-slate-500" />;
    }
  };

  const formatFee = (method: PaymentMethod) => {
    const feeFixed = method.fee_fixed ? `Rp ${method.fee_fixed.toLocaleString()}` : '';
    const feePercentage = method.fee_percentage ? `${parseFloat((method.fee_percentage * 100).toFixed(4))}%` : '';

    switch (method.fee_type) {
      case 'flat':
        return feeFixed ? `${feeFixed} Flat` : 'Free';
      case 'percentage':
        return feePercentage ? `${feePercentage}` : 'Free';
      case 'mixed':
        return feeFixed && feePercentage ? `${feeFixed} + ${feePercentage}` : 'Free';
      default:
        return 'Free';
    }
  };

  const renderMarkdown = (text: string) => {
    let html = text
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;');

    // Headings
    html = html.replace(/^### (.*?)$/gm, '<h5 class="text-sm font-bold mt-3 mb-1 text-slate-800 dark:text-slate-200">$1</h5>');
    html = html.replace(/^## (.*?)$/gm, '<h4 class="text-base font-bold mt-4 mb-2 text-slate-800 dark:text-slate-200">$1</h4>');
    html = html.replace(/^# (.*?)$/gm, '<h3 class="text-lg font-bold mt-5 mb-2 text-slate-800 dark:text-slate-200">$1</h3>');

    // Bold text
    html = html.replace(/\*\*(.*?)\*\*/g, '<strong class="font-bold text-slate-900 dark:text-white">$1</strong>');

    // Inline Code
    html = html.replace(/`(.*?)`/g, '<code class="bg-slate-100 dark:bg-zinc-800 px-1.5 py-0.5 rounded font-mono text-xs text-indigo-600 dark:text-indigo-400">$1</code>');

    // Bullet lists
    html = html.replace(/^\s*[-*]\s+(.*?)$/gm, '<li class="list-disc ml-5 my-1 text-slate-700 dark:text-slate-300">$1</li>');

    // Numbered lists
    html = html.replace(/^\s*(\d+)\.\s+(.*?)$/gm, '<li class="list-decimal ml-5 my-1 text-slate-700 dark:text-slate-300">$2</li>');

    // Links: [text](url)
    html = html.replace(/\[(.*?)\]\((.*?)\)/g, '<a href="$2" target="_blank" rel="noopener noreferrer" class="text-indigo-600 hover:underline font-medium">$1</a>');

    // Newlines to break lines (but not around list items or headers if we can avoid, simple break is fine)
    html = html.replace(/\n/g, '<br />');

    return (
      <div 
        dangerouslySetInnerHTML={{ __html: html }} 
        className="prose dark:prose-invert max-w-none text-sm leading-relaxed text-slate-700 dark:text-slate-300 font-sans space-y-1" 
      />
    );
  };

  return (
    <Dialog open={isOpen} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-2xl border border-border/80 shadow-2xl backdrop-blur-sm">
        <DialogHeader className="flex flex-row items-start space-x-4 border-b border-border/60 pb-4">
          <div className="rounded-xl bg-slate-100 dark:bg-slate-900 p-3 shadow-inner mt-1">
            {getMethodIcon(method.type)}
          </div>
          <div className="flex-1 space-y-1">
            <div className="flex items-center space-x-2">
              <DialogTitle className="text-2xl font-bold tracking-tight">{method.name}</DialogTitle>
              <Badge variant={method.is_active ? 'default' : 'secondary'} className="h-5">
                {method.is_active ? 'Active' : 'Inactive'}
              </Badge>
            </div>
            <DialogDescription className="font-mono text-xs text-muted-foreground uppercase tracking-wider">
              Code: <span className="text-foreground dark:text-white font-semibold">{method.code}</span> | Type: {method.type.replace('_', ' ')}
            </DialogDescription>
          </div>
        </DialogHeader>

        <div className="space-y-6 pt-4">
          <div className="space-y-2">
            <h4 className="text-sm font-semibold text-slate-500 uppercase tracking-wider">Description</h4>
            <p className="text-sm text-foreground leading-relaxed">{method.description}</p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 bg-slate-50 dark:bg-zinc-900/50 p-4 rounded-xl border border-border/40">
            <div className="space-y-1">
              <span className="text-xs text-muted-foreground font-medium uppercase">Fee Structure</span>
              <p className="text-sm font-semibold capitalize">{method.fee_type || 'None'}</p>
            </div>
            <div className="space-y-1">
              <span className="text-xs text-muted-foreground font-medium uppercase">Customer Fee</span>
              <p className="text-sm font-semibold text-indigo-600 dark:text-indigo-400">{formatFee(method)}</p>
            </div>
          </div>

          <div className="space-y-3">
            <div className="flex items-center space-x-2 border-b border-border/40 pb-2">
              <FileText className="h-4 w-4 text-muted-foreground" />
              <h4 className="text-sm font-semibold text-slate-500 uppercase tracking-wider">Payment Instructions</h4>
            </div>
            
            {method.instruction && method.instruction.content ? (
              <div className="bg-slate-100/50 dark:bg-zinc-900/30 p-5 rounded-xl border border-border/60 shadow-inner relative overflow-hidden">
                <div className="absolute top-0 left-0 w-1 h-full bg-indigo-500" />
                {renderMarkdown(method.instruction.content)}
              </div>
            ) : (
              <div className="flex items-center space-x-2.5 p-4 rounded-xl bg-amber-50/50 dark:bg-amber-950/10 border border-amber-200/40 text-amber-700 dark:text-amber-400">
                <ShieldCheck className="h-5 w-5 shrink-0" />
                <p className="text-sm">No payment instructions configured for this method.</p>
              </div>
            )}
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
