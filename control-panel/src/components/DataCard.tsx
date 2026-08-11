import * as React from 'react';
import { cn } from '../lib/utils';

export interface DataCardProps extends React.HTMLAttributes<HTMLDivElement> {
  children: React.ReactNode;
  className?: string;
  onClick?: (e: React.MouseEvent<HTMLDivElement>) => void;
  selected?: boolean;
}

export function DataCard({ children, className, onClick, selected, ...props }: DataCardProps) {
  return (
    <div
      onClick={onClick}
      className={cn(
        "grid grid-cols-1 md:grid-cols-12 gap-4 items-center px-4 py-3 rounded-xl border bg-background text-sm shadow-none transition-all select-none overflow-hidden",
        onClick ? "cursor-pointer hover:border-primary/50" : "",
        selected
          ? "border-primary/60 bg-primary/5 ring-1 ring-primary/45"
          : "border-border/60",
        className
      )}
      {...props}
    >
      {children}
    </div>
  );
}

export interface DataCardGridHeaderProps extends React.HTMLAttributes<HTMLDivElement> {
  children: React.ReactNode;
  className?: string;
}

export function DataCardGridHeader({ children, className, ...props }: DataCardGridHeaderProps) {
  return (
    <div
      className={cn(
        "grid grid-cols-12 gap-4 py-1.5 text-xs font-semibold uppercase tracking-wider text-muted-foreground select-none",
        className
      )}
      {...props}
    >
      {children}
    </div>
  );
}

export interface DataCardListProps extends React.HTMLAttributes<HTMLDivElement> {
  children: React.ReactNode;
  className?: string;
}

export function DataCardList({ children, className, ...props }: DataCardListProps) {
  return (
    <div className={cn("space-y-2.5 pr-1", className)} {...props}>
      {children}
    </div>
  );
}
