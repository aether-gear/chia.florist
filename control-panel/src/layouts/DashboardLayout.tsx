import React from 'react';
import { Outlet, Link, useLocation, useNavigate } from 'react-router-dom';
import { ShieldAlert, LayoutDashboard, ShoppingBag, Package, FileText, Truck, LogOut, Menu, Store, Users, Wallet, Crown, History, User, ChevronDown, Check, Globe } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Sheet, SheetContent, SheetTrigger } from '@/components/ui/sheet';
import { useAuthMeViewModel } from '../viewmodels/useAuthMeViewModel';
import { useMerchantProfileViewModel } from '../viewmodels/useMerchantProfileViewModel';
import { fetchApi } from '../lib/api';
import {
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuSub,
  DropdownMenuSubTrigger,
  DropdownMenuSubContent,
} from '../components/ui/dropdown-menu';
import { Avatar, AvatarImage, AvatarFallback } from '../components/ui/avatar';

type NavigationItem = {
  name: string;
  href: string;
  icon: React.ComponentType<{ className?: string }>;
  adminOnly?: boolean;
};

type NavigationGroup = {
  title: string | null;
  items: NavigationItem[];
};

const navigationGroups: NavigationGroup[] = [
  {
    title: null,
    items: [
      { name: 'Dashboard', href: '/', icon: LayoutDashboard },
    ]
  },
  {
    title: 'OPERATIONS',
    items: [
      { name: 'Orders', href: '/orders', icon: FileText },
      { name: 'Products', href: '/products', icon: Package },
      { name: 'Shop', href: '/shop', icon: ShoppingBag },
      { name: 'Shipments', href: '/shipments', icon: Truck, adminOnly: true },
    ]
  },
  {
    title: 'ADMIN',
    items: [
      { name: 'Customers', href: '/admin/customers', icon: Users, adminOnly: true },
      { name: 'Merchants', href: '/admin/merchants', icon: Store, adminOnly: true },
      { name: 'Audit Logs', href: '/admin/audit-logs', icon: History, adminOnly: true },
      { name: 'Security', href: '/security', icon: ShieldAlert },
    ]
  },
  {
    title: 'SETTINGS',
    items: [
      { name: 'Payment Settings', href: '/admin/payments', icon: Wallet },
    ]
  }
];

const allNavigation = navigationGroups.flatMap(g => g.items);

export default function DashboardLayout() {
  const location = useLocation();
  const navigate = useNavigate();

  const { data: authData, isAdmin } = useAuthMeViewModel();
  const { profile: staffProfile } = useMerchantProfileViewModel();

  const visibleNavigation = allNavigation.filter(n => isAdmin || !n.adminOnly);

  const handleLogout = async (e?: React.MouseEvent) => {
    if (e) e.preventDefault();
    try {
      await fetchApi('/auth/logout', { method: 'POST' });
    } catch (err) {
      console.error(err);
    } finally {
      localStorage.removeItem('isAuthenticated');
      localStorage.removeItem('userEmail');
      navigate('/login');
    }
  };

  const userEmail = localStorage.getItem('userEmail') || '';

  const renderProfileDropdown = () => {
    const fallbackInitials = staffProfile?.Name
      ? staffProfile.Name.split(' ').map(n => n[0]).join('').substring(0, 2).toUpperCase()
      : (isAdmin ? 'AD' : 'ME');

    return (
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <button className="flex items-center gap-2 hover:bg-slate-100 p-1.5 px-2.5 rounded-lg transition-colors cursor-pointer outline-none focus-visible:ring-2 focus-visible:ring-indigo-500 text-left">
            <Avatar className="h-8 w-8 ring-2 ring-indigo-500/10">
              {staffProfile?.AvatarURL && (
                <AvatarImage src={staffProfile.AvatarURL} alt={staffProfile.Name} className="object-cover" />
              )}
              <AvatarFallback className={`text-white font-bold text-xs uppercase ${isAdmin ? 'bg-amber-600' : 'bg-slate-700'}`}>
                {isAdmin ? '★' : fallbackInitials}
              </AvatarFallback>
            </Avatar>
            <span className="hidden sm:inline text-sm font-semibold text-slate-700 max-w-[120px] truncate">
              {staffProfile?.Name || (isAdmin ? 'Administrator' : 'User')}
            </span>
            <ChevronDown className="h-4 w-4 text-slate-400" />
          </button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" className="bg-white border border-slate-200 shadow-lg rounded-lg w-64 p-1">
          <div className="px-3 py-2.5 bg-slate-50 rounded-t-md border-b border-slate-100 mb-1">
            <div className="text-[10px] font-semibold text-slate-400 uppercase tracking-wider mb-0.5">Signed in as</div>
            <div className="text-sm font-bold text-slate-900 truncate" title={userEmail}>
              {userEmail}
            </div>
            {staffProfile?.Username && (
              <div className="text-xs text-slate-500 font-medium truncate mt-0.5">
                @{staffProfile.Username}
              </div>
            )}
            <div className="mt-1.5">
              <span className={`inline-flex items-center rounded-full px-2 py-0.5 text-[10px] font-medium ${
                isAdmin ? 'bg-amber-100 text-amber-800 border border-amber-200' : 'bg-indigo-100 text-indigo-800 border border-indigo-200'
              }`}>
                {isAdmin ? 'Administrator' : (authData?.roles[0]?.name || 'Merchant')}
              </span>
            </div>
          </div>

          <DropdownMenuItem asChild className="focus:bg-slate-50 cursor-pointer flex w-full items-center gap-2 px-3 py-2 text-sm text-slate-700 rounded-md">
            <Link to="/merchant/settings">
              <User className="h-4 w-4 text-slate-500" />
              <span>Account</span>
            </Link>
          </DropdownMenuItem>

          <DropdownMenuSeparator className="bg-slate-100" />

          <DropdownMenuItem
            onClick={() => handleLogout()}
            className="focus:bg-red-50 focus:text-red-600 text-red-500 cursor-pointer flex w-full items-center gap-2 px-3 py-2 text-sm font-medium rounded-md"
          >
            <LogOut className="h-4 w-4" />
            <span>Log out</span>
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    );
  };

  const renderSidebarContent = () => (
    <div className="flex h-full flex-col bg-slate-900">
      <div className="flex h-16 shrink-0 items-center px-6">
        <ShieldAlert className="h-8 w-8 text-indigo-500" />
        <span className="ml-3 text-lg font-bold text-white tracking-wide">
          Control Panel
        </span>
      </div>

      {/* Admin badge */}
      {isAdmin && (
        <div className="mx-3 mb-1 flex items-center gap-1.5 rounded-md bg-amber-500/15 px-3 py-1.5 text-xs font-semibold text-amber-400">
          <Crown className="h-3.5 w-3.5" />
          Administrator
        </div>
      )}

      <div className="flex-1 overflow-y-auto py-2">
        <nav className="space-y-1 px-3">
          {navigationGroups.map((group, idx) => {
            const visibleItems = group.items.filter(item => isAdmin || !item.adminOnly);
            if (visibleItems.length === 0) return null;

            return (
              <div key={group.title || `group-${idx}`} className="mb-6">
                {group.title && (
                  <div className="px-3 pt-4 pb-2">
                    <p className="text-xs font-semibold uppercase tracking-wider text-slate-500">
                      {group.title}
                    </p>
                  </div>
                )}
                <div className="space-y-1">
                  {visibleItems.map((item) => {
                    const isActive = location.pathname === item.href;
                    return (
                      <Link
                        key={item.name}
                        to={item.href}
                        className={`flex items-center px-3 py-2 text-sm font-medium rounded-md transition-colors ${
                          isActive ? 'bg-indigo-600 text-white' : 'text-slate-300 hover:bg-slate-800 hover:text-white'
                        }`}
                      >
                        <item.icon
                          className={`flex-shrink-0 -ml-1 mr-3 h-5 w-5 ${
                            isActive ? 'text-white' : 'text-slate-400 group-hover:text-white'
                          }`}
                          aria-hidden="true"
                        />
                        <span className="truncate">{item.name}</span>
                      </Link>
                    );
                  })}
                </div>
              </div>
            );
          })}
        </nav>
      </div>
    </div>
  );

  return (
    <div className="min-h-screen bg-slate-50 flex">
      {/* Mobile sidebar */}
      <Sheet>
        <div className="lg:hidden flex items-center justify-between p-4 border-b bg-white w-full fixed top-0 z-10 h-16 text-slate-800">
          <div className="flex items-center">
            <SheetTrigger asChild>
              <Button variant="ghost" size="icon" className="-ml-2">
                <Menu className="h-6 w-6" />
                <span className="sr-only">Open sidebar</span>
              </Button>
            </SheetTrigger>
            <div className="ml-4 font-bold">WAF Control Panel</div>
          </div>
          <div className="flex items-center">
            {renderProfileDropdown()}
          </div>
        </div>
        <SheetContent side="left" className="p-0 w-64 border-r-0">
          {renderSidebarContent()}
        </SheetContent>
      </Sheet>

      {/* Desktop sidebar */}
      <div className="hidden lg:flex lg:w-64 lg:flex-col lg:fixed lg:inset-y-0 z-20">
        {renderSidebarContent()}
      </div>

      {/* Main Content */}
      <div className="flex-1 lg:pl-64 flex flex-col pt-16 lg:pt-0 min-w-0">
        <header className="hidden lg:flex h-16 flex-shrink-0 items-center justify-between border-b bg-white px-8 text-slate-800">
          <h1 className="text-xl font-semibold">
            {visibleNavigation.find(n => n.href === location.pathname)?.name || 'Dashboard'}
          </h1>
          <div className="flex items-center space-x-4">
            {renderProfileDropdown()}
          </div>
        </header>
        <main className="flex-1 overflow-y-auto">
          <div className="mx-auto max-w-7xl p-4 sm:p-6 lg:p-8">
            <Outlet />
          </div>
        </main>
      </div>
    </div>
  );
}
