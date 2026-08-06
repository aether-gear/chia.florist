import { Search } from 'lucide-react';
import { Input } from './ui/input';
import { cn } from '../lib/utils';

interface SearchInputProps {
  value: string;
  onChange: (val: string) => void;
  placeholder?: string;
  className?: string;
  id?: string;
  variant?: 'default' | 'borderless';
  inputClassName?: string;
}

export default function SearchInput({
  value,
  onChange,
  placeholder = 'Search...',
  className = 'relative flex-1 max-w-sm w-full',
  id,
  variant = 'default',
  inputClassName
}: SearchInputProps) {
  return (
    <div className={className}>
      <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
      <Input
        id={id}
        type="search"
        placeholder={placeholder}
        className={cn(
          "pl-8 transition-colors",
          variant === 'borderless'
            ? "border-0 bg-transparent shadow-none focus-visible:ring-1 focus-visible:ring-primary/40 hover:bg-muted/30"
            : "",
          inputClassName
        )}
        value={value}
        onChange={(e) => onChange(e.target.value)}
      />
    </div>
  );
}

