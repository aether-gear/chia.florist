import { PackageOpen } from 'lucide-react';

export default function PlaceholderPage({ title }: { title: string }) {
  return (
    <div className="flex flex-col items-center justify-center h-[70vh] space-y-4 animate-in fade-in slide-in-from-left-4 duration-300">
      <div className="bg-indigo-100 p-6 rounded-full">
        <PackageOpen className="w-12 h-12 text-indigo-600" />
      </div>
      <h2 className="text-2xl font-semibold text-slate-800">{title} Management</h2>
      <p className="text-slate-500 max-w-md text-center">
        This section is currently under development. {title} management features will be coming soon in the next update.
      </p>
    </div>
  );
}
