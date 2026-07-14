import { Badge } from './ui/badge';

interface StatusBadgeProps {
  status: string;
  className?: string;
}

const getBadgeVariant = (status: string) => {
  const normalized = status.toLowerCase();
  switch (normalized) {
    case 'active':
    case 'success':
    case 'delivered':
    case 'confirmed':
    case 'completed':
    case 'whitelisted':
      return 'success' as const;
    case 'pending':
    case 'processing':
    case 'shipped':
    case 'ignored':
    case 'archived':
      return 'secondary' as const;
    case 'cancelled':
    case 'failed':
    case 'banned':
    case 'inactive':
      return 'danger' as const;
    default:
      return 'outline' as const;
  }
};

export default function StatusBadge({ status, className = '' }: StatusBadgeProps) {
  const variant = getBadgeVariant(status);
  return (
    <Badge
      variant={variant}
      className={`font-semibold uppercase tracking-wide text-xs px-2.5 py-0.5 rounded-full transition-colors ${className}`}
    >
      {status}
    </Badge>
  );
}
