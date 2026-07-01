import { Outlet, Link, useLocation, useNavigate } from 'react-router-dom';
import { ShieldAlert, LayoutDashboard, ShoppingBag, Package, FileText, Activity, Truck, LogOut, Menu, Settings, Store, Users, Wallet, Crown } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Sheet, SheetContent, SheetTrigger } from '@/components/ui/sheet';
import { useAuthMeViewModel } from '../viewmodels/useAuthMeViewModel';

type NavigationItem = {
  name: string;
  href: string;
  icon: any;
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
      { name: 'Transactions', href: '/transactions', icon: Activity },
      { name: 'Shipments', href: '/shipments', icon: Truck, adminOnly: true },
    ]
  },
  {
    title: 'ADMIN',
    items: [
      { name: 'Customers', href: '/admin/customers', icon: Users, adminOnly: true },
      { name: 'Merchants', href: '/admin/merchants', icon: Store, adminOnly: true },
      { name: 'Security', href: '/security', icon: ShieldAlert },
    ]
  },
  {
    title: 'SETTINGS',
    items: [
      { name: 'Payment Settings', href: '/admin/payments', icon: Wallet },
      { name: 'Profile', href: '/merchant/settings', icon: Settings },
    ]
  }
];

const allNavigation = navigationGroups.flatMap(g => g.items);

export default function DashboardLayout() {
  const location = useLocation();
  const navigate = useNavigate();
  const { data: authData, isAdmin } = useAuthMeViewModel();

  const visibleNavigation = allNavigation.filter(n => isAdmin || !n.adminOnly);

  const handleLogout = async (e: React.MouseEvent) => {
    e.preventDefault();
    try {
      await fetch('/api/core/auth/logout', { method: 'POST' });
    } catch (err) {
      console.error(err);
    } finally {
      localStorage.removeItem('isAuthenticated');
      localStorage.removeItem('userEmail');
      navigate('/login');
    }
  };

  const userEmail = localStorage.getItem('userEmail') || '';

  const SidebarContent = () => (
    <div className="flex h-full flex-col bg-slate-900">
      <div className="flex h-16 shrink-0 items-center px-6">
        <ShieldAlert className="h-8 w-8 text-indigo-500" />
        <span className="ml-3 text-lg font-bold text-white tracking-wide">
          WAF Control
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
      <div className="p-4 border-t border-slate-800">
        <div className="flex items-center">
          <div className="flex-shrink-0">
            <div className={`h-8 w-8 rounded-full flex items-center justify-center text-white font-bold text-xs uppercase ${isAdmin ? 'bg-amber-600' : 'bg-slate-700'}`}>
              {isAdmin ? '★' : (authData?.account_type ? authData.account_type.substring(0, 2) : 'U')}
            </div>
          </div>
          <div className="ml-3 min-w-0">
            <p className="text-sm font-medium text-white truncate">
              {userEmail || (authData ? authData.account_type : 'Loading...')}
            </p>
            <p className="text-xs font-medium text-slate-400">
              {isAdmin ? 'Administrator' : (authData?.roles[0]?.name || 'Merchant')}
            </p>
          </div>
        </div>
      </div>
    </div>
  );

  return (
    <div className="min-h-screen bg-slate-50 flex">
      {/* Mobile sidebar */}
      <Sheet>
        <div className="lg:hidden flex items-center p-4 border-b bg-white w-full fixed top-0 z-10 h-16">
          <SheetTrigger asChild>
            <Button variant="ghost" size="icon" className="-ml-2">
              <Menu className="h-6 w-6" />
              <span className="sr-only">Open sidebar</span>
            </Button>
          </SheetTrigger>
          <div className="ml-4 font-bold">WAF Control Panel</div>
        </div>
        <SheetContent side="left" className="p-0 w-64 border-r-0">
          <SidebarContent />
        </SheetContent>
      </Sheet>

      {/* Desktop sidebar */}
      <div className="hidden lg:flex lg:w-64 lg:flex-col lg:fixed lg:inset-y-0 z-20">
        <SidebarContent />
      </div>

      {/* Main Content */}
      <div className="flex-1 lg:pl-64 flex flex-col pt-16 lg:pt-0 min-w-0">
        <header className="hidden lg:flex h-16 flex-shrink-0 items-center justify-between border-b bg-white px-8">
          <h1 className="text-xl font-semibold text-slate-800">
            {visibleNavigation.find(n => n.href === location.pathname)?.name || 'Dashboard'}
          </h1>
          <div className="flex items-center space-x-4">
            <Button variant="outline" size="sm" onClick={handleLogout}>
              <LogOut className="mr-2 h-4 w-4" />
              Sign out
            </Button>
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
