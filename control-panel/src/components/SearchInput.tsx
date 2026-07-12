import { Search } from 'lucide-react';
import { Input } from './ui/input';

interface SearchInputProps {
  value: string;
  onChange: (val: string) => void;
  placeholder?: string;
  className?: string;
  id?: string;
}

export default function SearchInput({
  value,
  onChange,
  placeholder = 'Search...',
  className = 'relative flex-1 max-w-sm w-full',
  id
}: SearchInputProps) {
  return (
    <div className={className}>
      <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
      <Input
        id={id}
        type="search"
        placeholder={placeholder}
        className="pl-8"
        value={value}
        onChange={(e) => onChange(e.target.value)}
      />
    </div>
  );
}
