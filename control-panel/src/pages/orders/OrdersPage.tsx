import { useState, useEffect, useRef } from 'react';
import { Link } from 'react-router-dom';
import { BarChart3 } from 'lucide-react';
import { Button } from '../../components/ui/button';
import { useOrdersViewModel } from '../../viewmodels/useOrdersViewModel';
import { useOrderActionsViewModel, type ShipmentDispatchPayload } from '../../viewmodels/useOrderActionsViewModel';
import { useAuthMeViewModel } from '../../viewmodels/useAuthMeViewModel';
import { fetchApi } from '../../lib/api';
import OrderFilters, { type ShopOption } from '../../components/orders/OrderFilters';
import OrdersTable from '../../components/orders/OrdersTable';
import OrderDetailInspector from '../../components/orders/OrderDetailInspector';

export default function OrdersPage() {
  const { isAdmin } = useAuthMeViewModel();
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
  } = useOrdersViewModel();

  const {
    submitting,
    fetchOrderTracking,
    updateOrderStatus,
    updateShipmentStatus,
    updateShipmentDetails,
  } = useOrderActionsViewModel();

  const [selectedOrderId, setSelectedOrderId] = useState<string | null>(null);
  const [shopsList, setShopsList] = useState<ShopOption[]>([]);
  const shopsFetchedRef = useRef(false);

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
      });
  }, [isAdmin, setShopFilter]);

  const selectedOrder = data?.orders.find((o) => o.id === selectedOrderId) || null;

  const handleSort = () => {
    const currentDirection = sort.split(':')[1];
    const newDirection = currentDirection === 'desc' ? 'asc' : 'desc';
    setSort(`latest:${newDirection}`);
    setPage(1);
  };

  const handleStartProcessing = async (orderId: string) => {
    await updateOrderStatus(orderId, 'processing');
    refresh();
  };

  const handleDispatchOrder = async (orderId: string, shipments: ShipmentDispatchPayload[]) => {
    await updateOrderStatus(orderId, 'shipped', undefined, undefined, shipments);
    refresh();
  };

  const handleUpdateShipmentStatus = async (shipmentId: string, status: string) => {
    await updateShipmentStatus(shipmentId, status);
    refresh();
  };

  const handleUpdateWaybill = async (
    shipmentId: string,
    details: { tracking_number?: string; courier?: string; service?: string }
  ) => {
    await updateShipmentDetails(shipmentId, details);
    refresh();
  };

  return (
    <div className="flex-col md:flex">
      <div className="flex-1 space-y-8 p-6 sm:p-8 lg:p-12 animate-in fade-in duration-300">
        {/* Header */}
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <h2 className="text-3xl font-bold font-display tracking-tight text-foreground">Orders</h2>
            <p className="text-muted-foreground text-sm">
              Verify payments, fulfill shipments, and manage orders
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
              setSelectedOrderId(null);
            }}
            onShopChange={(shop) => setShopFilter(shop)}
            onFromDateChange={(date) => setFromDate(date)}
            onToDateChange={(date) => setToDate(date)}
            onRefresh={refresh}
          />

          {/* Master-Detail Split Workspace Layout */}
          <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">
            {/* LEFT PANE: Master Order Cards List */}
            <div className={`lg:col-span-5 flex flex-col space-y-4 ${selectedOrderId ? 'hidden lg:flex' : 'flex'}`}>
              <OrdersTable
                orders={data?.orders || []}
                total={data?.total || 0}
                page={page}
                limit={limit}
                sort={sort}
                loading={loading}
                isSwitchingCategory={isSwitchingCategory}
                error={error}
                selectedOrderId={selectedOrderId}
                onSelectOrder={(id) => setSelectedOrderId(id)}
                onSortChange={handleSort}
                onPageChange={setPage}
              />
            </div>

            {/* RIGHT PANE: Order Action & Detail Inspector */}
            <div
              className={`lg:col-span-7 border border-border/80 rounded-2xl bg-card flex flex-col lg:sticky lg:top-24 lg:self-start lg:max-h-[calc(100vh-6rem)] overflow-y-auto overscroll-contain pr-0.5 ${
                !selectedOrderId
                  ? 'hidden lg:flex items-center justify-center p-12 text-center text-muted-foreground min-h-[400px]'
                  : 'flex'
              }`}
            >
              <OrderDetailInspector
                order={selectedOrder}
                submitting={submitting}
                onClose={() => setSelectedOrderId(null)}
                onStartProcessing={handleStartProcessing}
                onDispatchOrder={handleDispatchOrder}
                onUpdateShipmentStatus={handleUpdateShipmentStatus}
                onUpdateWaybill={handleUpdateWaybill}
                fetchOrderTracking={fetchOrderTracking}
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
