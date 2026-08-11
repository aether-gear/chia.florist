import React from 'react';
import { Link } from 'react-router-dom';
import { ChevronRight, Home } from 'lucide-react';

export interface BreadcrumbItem {
  label: string;
  onClick?: () => void;
  href?: string;
  icon?: React.ComponentType<{ className?: string }>;
}

export interface BreadcrumbProps {
  items: BreadcrumbItem[];
  showHome?: boolean;
  homeHref?: string;
  className?: string;
  separator?: React.ReactNode;
}

export default function Breadcrumb({
  items,
  showHome = false,
  homeHref = '/',
  className = '',
  separator = <ChevronRight className="h-3.5 w-3.5 text-muted-foreground/60 flex-shrink-0" />,
}: BreadcrumbProps) {
  const allItems: BreadcrumbItem[] = showHome
    ? [{ label: 'Home', href: homeHref, icon: Home }, ...items]
    : items;

  return (
    <nav aria-label="Breadcrumb" className={`flex items-center text-xs md:text-sm font-sans ${className}`}>
      <ol className="flex items-center flex-wrap gap-1.5 text-muted-foreground">
        {allItems.map((item, index) => {
          const isLast = index === allItems.length - 1;
          const Icon = item.icon;

          return (
            <li key={index} className="inline-flex items-center gap-1.5">
              {index > 0 && <span aria-hidden="true" className="select-none">{separator}</span>}

              {isLast ? (
                <span className="font-semibold text-foreground flex items-center gap-1.5" aria-current="page">
                  {Icon && <Icon className="h-3.5 w-3.5 text-foreground" />}
                  {item.label}
                </span>
              ) : item.onClick ? (
                <button
                  type="button"
                  onClick={item.onClick}
                  className="hover:text-primary transition-colors flex items-center gap-1.5 font-medium focus:outline-none focus:ring-1 focus:ring-primary/40 rounded-sm"
                >
                  {Icon && <Icon className="h-3.5 w-3.5" />}
                  {item.label}
                </button>
              ) : item.href ? (
                <Link
                  to={item.href}
                  className="hover:text-primary transition-colors flex items-center gap-1.5 font-medium"
                >
                  {Icon && <Icon className="h-3.5 w-3.5" />}
                  {item.label}
                </Link>
              ) : (
                <span className="flex items-center gap-1.5 font-medium">
                  {Icon && <Icon className="h-3.5 w-3.5" />}
                  {item.label}
                </span>
              )}
            </li>
          );
        })}
      </ol>
    </nav>
  );
}
