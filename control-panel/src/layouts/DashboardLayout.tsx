import React from 'react';
import { Outlet, Link, useLocation, useNavigate } from 'react-router-dom';
import { LayoutDashboard, ShoppingBag, Package, FileText, Truck, LogOut, Menu, Store, Users, Wallet, Crown, User, ChevronDown, Shield, ClipboardClock, BarChart3 } from 'lucide-react';
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
      { name: 'Analytics', href: '/admin/analytics', icon: BarChart3, adminOnly: true },
      { name: 'Customers', href: '/admin/customers', icon: Users, adminOnly: true },
      { name: 'Merchants', href: '/admin/merchants', icon: Store, adminOnly: true },
      { name: 'Audit Logs', href: '/admin/audit-logs', icon: ClipboardClock, adminOnly: true },
      { name: 'Security', href: '/security', icon: Shield },
    ]
  },
  {
    title: 'SETTINGS',
    items: [
      { name: 'Payment Settings', href: '/admin/payments', icon: Wallet },
    ]
  }
];

export default function DashboardLayout() {
  const location = useLocation();
  const navigate = useNavigate();

  const { data: authData, isAdmin } = useAuthMeViewModel();
  const { profile: staffProfile } = useMerchantProfileViewModel();

  const handleLogout = async (e?: React.MouseEvent) => {
    if (e) e.preventDefault();
    try {
      await fetchApi('/auth/logout', { method: 'POST' });
    } catch (err) {
      console.error(err);
    } finally {
      localStorage.removeItem('isAuthenticated');
      localStorage.removeItem('userEmail');
      sessionStorage.removeItem('isAuthenticated');
      sessionStorage.removeItem('userEmail');
      navigate('/login');
    }
  };

  const userEmail = localStorage.getItem('userEmail') || sessionStorage.getItem('userEmail') || '';

  const renderProfileDropdown = () => {
    const fallbackInitials = staffProfile?.Name
      ? staffProfile.Name.split(' ').map(n => n[0]).join('').substring(0, 2).toUpperCase()
      : (isAdmin ? 'AD' : 'ME');

    return (
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <button className="flex items-center gap-2 hover:bg-muted p-1.5 px-2.5 rounded-lg transition-colors cursor-pointer outline-none focus-visible:ring-2 focus-visible:ring-ring text-left">
            <Avatar className="h-8 w-8 ring-2 ring-primary/10">
              {staffProfile?.AvatarURL && (
                <AvatarImage src={staffProfile.AvatarURL} alt={staffProfile.Name} className="object-cover" />
              )}
              <AvatarFallback className={`font-bold text-xs uppercase bg-primary text-primary-foreground`}>
                {isAdmin ? '★' : fallbackInitials}
              </AvatarFallback>
            </Avatar>
            <span className="hidden sm:inline text-sm font-semibold text-foreground max-w-[120px] truncate">
              {staffProfile?.Name || (isAdmin ? 'Administrator' : 'User')}
            </span>
            <ChevronDown className="h-4 w-4 text-muted-foreground/60" />
          </button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" className="w-64 p-1">
          <div className="px-3 py-2.5 bg-muted/50 rounded-t-md border-b border-border/85 mb-1">
            <div className="text-[10px] font-semibold text-muted-foreground/80 uppercase tracking-wider mb-0.5">Signed in as</div>
            <div className="text-sm font-bold text-foreground truncate" title={userEmail}>
              {userEmail}
            </div>
            {staffProfile?.Username && (
              <div className="text-xs text-muted-foreground font-medium truncate mt-0.5">
                @{staffProfile.Username}
              </div>
            )}
            <div className="mt-1.5">
              <span className={`inline-flex items-center rounded-full px-2 py-0.5 text-[10px] font-medium bg-primary/10 text-primary border border-primary/20`}>
                {isAdmin ? 'Administrator' : (authData?.roles[0]?.name || 'Merchant')}
              </span>
            </div>
          </div>

          <DropdownMenuItem asChild className="cursor-pointer flex w-full items-center gap-2 px-3 py-2 text-sm rounded-md">
            <Link to="/merchant/settings">
              <User className="h-4 w-4 text-muted-foreground" />
              <span>Account</span>
            </Link>
          </DropdownMenuItem>

          <DropdownMenuSeparator />

          <DropdownMenuItem
            onClick={() => handleLogout()}
            className="focus:bg-destructive/10 focus:text-destructive text-destructive cursor-pointer flex w-full items-center gap-2 px-3 py-2 text-sm font-medium rounded-md"
          >
            <LogOut className="h-4 w-4" />
            <span>Log out</span>
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    );
  };

  const renderSidebarContent = () => (
    <div className="flex h-full flex-col bg-zinc-100">
      <div className="flex h-16 shrink-0 items-center px-6 gap-2.5 border-b border-border/30">
        <div className="w-8 h-8 rounded-lg bg-primary/10 flex items-center justify-center border border-primary/20">
          <svg className="w-4.5 h-4.5 text-primary" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
            <path d="M12 22c0-5.523-4.477-10-10-10 5.523 0 10-4.477 10-10 0 5.523 4.477 10 10 10-5.523 0-10 4.477-10 10z" />
          </svg>
        </div>
        <span className="font-display font-bold text-base tracking-tight text-foreground">
          chia.florist
        </span>
      </div>

      {/* Admin badge */}
      {isAdmin && (
        <div className="mx-4 mt-4 mb-1 flex items-center gap-1.5 rounded-lg bg-amber-500/10 px-2.5 py-1 text-[11px] font-semibold text-amber-700 dark:text-amber-300 border border-amber-500/20">
          <Crown className="h-3 w-3" />
          <span>Administrator</span>
        </div>
      )}

      <div className="flex-1 overflow-y-auto py-4">
        <nav className="space-y-1 px-3">
          {navigationGroups.map((group, idx) => {
            const visibleItems = group.items.filter(item => isAdmin || !item.adminOnly);
            if (visibleItems.length === 0) return null;

            return (
              <div key={group.title || `group-${idx}`} className="mb-6">
                {group.title && (
                  <div className="px-3 pt-3 pb-1">
                    <p className="text-[10px] font-bold font-sans uppercase tracking-widest text-muted-foreground/60">
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
                        className={`flex items-center px-3 py-2 text-sm font-medium rounded-md transition-all duration-200 ease-out active:scale-[0.98] ${
                          isActive ? 'bg-primary text-primary-foreground font-semibold shadow-sm shadow-primary/5' : 'text-muted-foreground hover:bg-muted/80 hover:text-foreground'
                        }`}
                      >
                        <item.icon
                          className={`flex-shrink-0 -ml-1 mr-3 h-5 w-5 transition-colors ${
                            isActive ? 'text-primary-foreground' : 'text-muted-foreground group-hover:text-foreground'
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
    <div className="min-h-screen bg-background flex text-foreground font-sans antialiased">
      {/* Mobile sidebar */}
      <Sheet>
        <div className="lg:hidden flex items-center justify-between p-4 border-b border-border/40 bg-background/95 backdrop-blur-md w-full fixed top-0 z-10 h-16 text-foreground">
          <div className="flex items-center">
            <SheetTrigger asChild>
              <Button variant="ghost" size="icon" className="-ml-2">
                <Menu className="h-6 w-6" />
                <span className="sr-only">Open sidebar</span>
              </Button>
            </SheetTrigger>
            <div className="ml-2 flex items-center gap-1.5">
              <div className="w-6 h-6 rounded-md bg-primary/10 flex items-center justify-center border border-primary/20">
                <svg className="w-3.5 h-3.5 text-primary" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
                  <path d="M12 22c0-5.523-4.477-10-10-10 5.523 0 10-4.477 10-10 0 5.523 4.477 10 10 10-5.523 0-10 4.477-10 10z" />
                </svg>
              </div>
              <span className="font-display font-bold text-sm tracking-tight text-foreground">
                chia.florist
              </span>
            </div>
          </div>
          <div className="flex items-center">
            {renderProfileDropdown()}
          </div>
        </div>
        <SheetContent side="left" className="p-0 w-64 border-r border-border/40">
          {renderSidebarContent()}
        </SheetContent>
      </Sheet>

      {/* Desktop sidebar */}
      <div className="hidden lg:flex lg:w-64 lg:flex-col lg:fixed lg:inset-y-0 z-20 border-r border-border/50">
        {renderSidebarContent()}
      </div>

      {/* Main Content */}
      <div className="flex-1 lg:pl-64 flex flex-col pt-16 lg:pt-0 min-w-0">
        <header className="hidden lg:flex h-16 flex-shrink-0 items-center justify-between border-b border-border/40 bg-background/95 backdrop-blur-md px-8 text-foreground sticky top-0 z-10">
          <h1 className="text-xl font-bold font-display tracking-tight text-foreground">
            {(() => {
              const matchedGroup = navigationGroups.find(g =>
                g.items.some(item => location.pathname === item.href || location.pathname.startsWith(item.href + '/'))
              );
              return matchedGroup?.title
                ? matchedGroup.title.charAt(0) + matchedGroup.title.slice(1).toLowerCase()
                : 'Dashboard';
            })()}
          </h1>
          <div className="flex items-center space-x-4">
            {renderProfileDropdown()}
          </div>
        </header>
        <main className="flex-1">
          <div className="mx-auto max-w-7xl py-4 sm:p-8 animate-in fade-in duration-300">
            <Outlet />
          </div>
        </main>
      </div>
    </div>
  );
}
