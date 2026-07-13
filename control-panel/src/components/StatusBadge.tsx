import { Badge } from './ui/badge';

interface StatusBadgeProps {
  status: string;
  className?: string;
}

const getBadgeConfig = (status: string) => {
  const normalized = status.toLowerCase();
  switch (normalized) {
    case 'active':
    case 'success':
    case 'delivered':
    case 'confirmed':
    case 'completed':
    case 'whitelisted':
      return {
        variant: 'default' as const,
        style: 'bg-emerald-100 text-emerald-800 border-0 hover:bg-emerald-200/80'
      };
    case 'pending':
    case 'processing':
    case 'shipped':
    case 'ignored':
    case 'archived':
      return {
        variant: 'secondary' as const,
        style: 'bg-slate-100 text-slate-800 border-0 hover:bg-slate-200/80'
      };
    case 'cancelled':
    case 'failed':
    case 'banned':
    case 'inactive':
      return {
        variant: 'destructive' as const,
        style: 'bg-rose-100 text-rose-800 border-0 hover:bg-rose-200/80'
      };
    default:
      return {
        variant: 'outline' as const,
        style: 'bg-slate-50 text-slate-600 border-slate-200'
      };
  }
};

export default function StatusBadge({ status, className = '' }: StatusBadgeProps) {
  const config = getBadgeConfig(status);
  return (
    <Badge variant={config.variant} className={`${config.style} ${className} font-semibold uppercase tracking-wide text-xs px-2.5 py-0.5 rounded-full transition-colors`}>
      {status}
    </Badge>
  );
}
