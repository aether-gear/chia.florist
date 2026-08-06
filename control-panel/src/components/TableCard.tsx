import * as React from 'react';
import { cn } from '../lib/utils';

export interface TableCardProps extends React.HTMLAttributes<HTMLDivElement> {
  children: React.ReactNode;
  className?: string;
}

export function TableCard({ children, className, ...props }: TableCardProps) {
  return (
    <div
      className={cn(
        "rounded-2xl border border-border/60 bg-card overflow-hidden shadow-none transition-all",
        className
      )}
      {...props}
    >
      {children}
    </div>
  );
}

export interface TableCardHeaderProps extends React.HTMLAttributes<HTMLDivElement> {
  children: React.ReactNode;
  className?: string;
}

export function TableCardHeader({ children, className, ...props }: TableCardHeaderProps) {
  return (
    <div
      className={cn(
        "p-4 border-b border-border/60 flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-4 bg-muted/20",
        className
      )}
      {...props}
    >
      {children}
    </div>
  );
}

export interface TableCardTitleProps extends React.HTMLAttributes<HTMLDivElement> {
  title?: React.ReactNode;
  subtitle?: React.ReactNode;
  children?: React.ReactNode;
  className?: string;
}

export function TableCardTitle({ title, subtitle, children, className, ...props }: TableCardTitleProps) {
  return (
    <div className={cn("space-y-0.5", className)} {...props}>
      {title && (
        <h3 className="font-bold font-display tracking-tight text-base sm:text-lg text-foreground">
          {title}
        </h3>
      )}
      {subtitle && (
        <p className="text-muted-foreground text-xs sm:text-sm">
          {subtitle}
        </p>
      )}
      {children}
    </div>
  );
}

export interface TableCardContentProps extends React.HTMLAttributes<HTMLDivElement> {
  children: React.ReactNode;
  className?: string;
}

export function TableCardContent({ children, className, ...props }: TableCardContentProps) {
  return (
    <div className={cn("relative w-full overflow-x-auto", className)} {...props}>
      {children}
    </div>
  );
}

export interface TableCardFooterProps extends React.HTMLAttributes<HTMLDivElement> {
  children: React.ReactNode;
  className?: string;
}

export function TableCardFooter({ children, className, ...props }: TableCardFooterProps) {
  return (
    <div
      className={cn(
        "p-4 border-t border-border/60 flex flex-col sm:flex-row items-center justify-between gap-4 bg-muted/10",
        className
      )}
      {...props}
    >
      {children}
    </div>
  );
}
