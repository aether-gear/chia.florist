import { Skeleton } from './ui/skeleton';

interface LoadingStateProps {
  message?: string;
  className?: string;
  rows?: number;
}

export default function LoadingState({
  message = 'Loading...',
  className = 'w-full p-6 space-y-4',
  rows = 3
}: LoadingStateProps) {
  return (
    <div className={className} aria-label={message}>
      <div className="flex items-center justify-between pb-2 border-b border-muted">
        <div className="flex items-center gap-3">
          <Skeleton className="h-5 w-32" />
          <Skeleton className="h-5 w-20" />
        </div>
        <Skeleton className="h-8 w-24" />
      </div>
      <div className="space-y-3">
        {Array.from({ length: rows }).map((_, i) => (
          <div key={i} className="flex items-center justify-between py-2 border-b border-muted/50 last:border-0">
            <div className="flex items-center gap-3 w-2/3">
              <Skeleton className="h-4 w-12" />
              <Skeleton className="h-4 w-full max-w-[200px]" />
              <Skeleton className="h-4 w-full max-w-[120px] hidden md:block" />
            </div>
            <div className="flex items-center gap-2">
              <Skeleton className="h-4 w-16" />
              <Skeleton className="h-7 w-7 rounded-full" />
            </div>
          </div>
        ))}
      </div>
      <div className="flex items-center justify-center text-xs text-muted-foreground animate-pulse gap-2 pt-2">
        <div className="h-3 w-3 rounded-full border-2 border-primary border-t-transparent animate-spin" />
        <span>{message}</span>
      </div>
    </div>
  );
}
