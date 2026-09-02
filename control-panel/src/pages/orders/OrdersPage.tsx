import { useState, useEffect, useRef } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { BarChart3 } from 'lucide-react';
import { Button } from '../../components/ui/button';
import { useOrdersViewModel } from '../../viewmodels/useOrdersViewModel';
import { useAuthMeViewModel } from '../../viewmodels/useAuthMeViewModel';
import { fetchApi } from '../../lib/api';
import OrderFilters, { type ShopOption } from '../../components/orders/OrderFilters';
import OrdersTable from '../../components/orders/OrdersTable';

export default function OrdersPage() {
  const navigate = useNavigate();
  const { isAdmin, loading: authLoading } = useAuthMeViewModel();
  const [shopsList, setShopsList] = useState<ShopOption[]>([]);
  const [shopsLoaded, setShopsLoaded] = useState<boolean>(false);
  const shopsFetchedRef = useRef(false);

  // For admins: orders can be loaded immediately across all shops
  // For staff: wait until shops are loaded so the initial fetch directly targets their assigned shop
  const isOrdersEnabled = !authLoading && (isAdmin || shopsLoaded);

  const {
    data,
    loading,
    isSwitchingCategory,
    error,
    page,
    limit,
    sort,
    searchNumber,
    statusFilter,
    shopFilter,
    fromDate,
    toDate,
    setPage,
    setSort,
    setSearchNumber,
    setStatusFilter,
    setShopFilter,
    setFromDate,
    setToDate,
    refresh,
  } = useOrdersViewModel(undefined, { enabled: isOrdersEnabled });

  useEffect(() => {
    if (shopsFetchedRef.current) return;
    shopsFetchedRef.current = true;

    fetchApi('/shops?limit=100')
      .then((res) => {
        if (res?.shops && res.shops.length > 0) {
          setShopsList(res.shops);
          if (!isAdmin) {
            const firstShop = res.shops[0];
            setShopFilter(firstShop.slug || firstShop.id);
          }
        }
      })
      .catch((err) => {
        console.error('Failed to load shops for order filter', err);
      })
      .finally(() => {
        setShopsLoaded(true);
      });
  }, [isAdmin, setShopFilter]);

  const handleSort = () => {
    const currentDirection = sort.split(':')[1];
    const newDirection = currentDirection === 'desc' ? 'asc' : 'desc';
    setSort(`latest:${newDirection}`);
    setPage(1);
  };

  const handleSelectOrder = (orderId: string) => {
    navigate(`/orders/${orderId}`);
  };

  return (
    <div className="flex-col md:flex">
      <div className="flex-1 space-y-8 p-6 sm:p-8 lg:p-12 animate-in fade-in slide-in-from-left-4 duration-300">
        {/* Header */}
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <h2 className="text-3xl font-bold font-display tracking-tight text-foreground">Orders</h2>
            <p className="text-muted-foreground text-sm">
              Verify payments, fulfill shipments, and manage customer orders
            </p>
          </div>
          {isAdmin && (
            <Button
              asChild
              variant="outline"
              size="sm"
              className="rounded-xl text-xs gap-1.5 border-primary/30 text-primary hover:bg-primary/5 self-start sm:self-auto"
            >
              <Link to="/admin/analytics?tab=orders">
                <BarChart3 className="h-3.5 w-3.5" /> View Order Analytics →
              </Link>
            </Button>
          )}
        </div>

        {/* Order Workspace Section */}
        <div className="space-y-6">
          <div className="pb-4 border-b border-border/60">
            <h3 className="text-xl font-bold font-display tracking-tight text-foreground">Order Workspace</h3>
            <p className="text-muted-foreground text-sm">
              Manage fulfillment, check payment updates, and inspect order status logs.
            </p>
          </div>

          {/* Filters Bar */}
          <OrderFilters
            searchNumber={searchNumber}
            statusFilter={statusFilter}
            shopFilter={shopFilter}
            fromDate={fromDate}
            toDate={toDate}
            shopsList={shopsList}
            isAdmin={isAdmin}
            loading={loading}
            isSwitchingCategory={isSwitchingCategory}
            onSearchChange={(val) => {
              setSearchNumber(val);
              setPage(1);
            }}
            onStatusChange={(status) => {
              setStatusFilter(status);
            }}
            onShopChange={(shop) => setShopFilter(shop)}
            onFromDateChange={(date) => setFromDate(date)}
            onToDateChange={(date) => setToDate(date)}
            onRefresh={refresh}
          />

          {/* Full-Width Orders List */}
          <div className="w-full">
            <OrdersTable
              orders={data?.orders || []}
              total={data?.total || 0}
              page={page}
              limit={limit}
              sort={sort}
              loading={loading}
              isSwitchingCategory={isSwitchingCategory}
              error={error}
              onSelectOrder={handleSelectOrder}
              onSortChange={handleSort}
              onPageChange={setPage}
            />
          </div>
        </div>
      </div>
    </div>
  );
}
