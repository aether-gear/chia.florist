import React from 'react';
import { Button } from './ui/button';

interface EmptyStateProps {
  icon?: React.ReactNode;
  title: string;
  description?: string;
  actionLabel?: string;
  onAction?: () => void;
  className?: string;
}

export default function EmptyState({
  icon,
  title,
  description,
  actionLabel,
  onAction,
  className = 'flex h-[50vh] flex-col items-center justify-center text-center p-4 gap-3 border border-dashed rounded-lg bg-slate-50/50'
}: EmptyStateProps) {
  return (
    <div className={className}>
      {icon && <div className="text-slate-400">{icon}</div>}
      <div className="space-y-1">
        <h3 className="text-base font-semibold text-slate-800">{title}</h3>
        {description && <p className="text-sm text-slate-500 max-w-sm">{description}</p>}
      </div>
      {actionLabel && onAction && (
        <Button onClick={onAction} className="mt-2" size="sm">
          {actionLabel}
        </Button>
      )}
    </div>
  );
}
