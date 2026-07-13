import { Loader2 } from 'lucide-react';

interface LoadingStateProps {
  message?: string;
  className?: string;
}

export default function LoadingState({ message = 'Loading...', className = 'flex h-[50vh] flex-col items-center justify-center gap-2' }: LoadingStateProps) {
  return (
    <div className={className}>
      <Loader2 className="h-8 w-8 animate-spin text-indigo-600" />
      <span className="text-sm font-medium text-slate-500">{message}</span>
    </div>
  );
}
