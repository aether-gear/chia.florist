import React, { useState } from 'react';
import { useShopOrdersViewModel } from '../../viewmodels/useShopOrdersViewModel';
import { useOrderActionsViewModel, type ShipmentDispatchPayload } from '../../viewmodels/useOrderActionsViewModel';
import OrderFilters from '../orders/OrderFilters';
import OrdersTable from '../orders/OrdersTable';
import OrderDetailInspector from '../orders/OrderDetailInspector';

export interface ShopOrdersTabProps {
  shopId: string;
  shopName?: string;
}

export const ShopOrdersTab: React.FC<ShopOrdersTabProps> = ({ shopId, shopName }) => {
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
    fromDate,
    toDate,
    setPage,
    setSort,
    setSearchNumber,
    setStatusFilter,
    setFromDate,
    setToDate,
    refresh,
  } = useShopOrdersViewModel(shopId);

  const {
    submitting,
    fetchOrderTracking,
    updateOrderStatus,
    updateShipmentStatus,
    updateShipmentDetails,
  } = useOrderActionsViewModel();

  const [selectedOrderId, setSelectedOrderId] = useState<string | null>(null);

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
    <div className="space-y-6 pt-2">
      <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between pb-4 border-b border-border/60 gap-2">
        <div>
          <h3 className="text-lg font-bold font-display text-foreground">
            {shopName ? `${shopName} Orders` : 'Shop Orders'}
          </h3>
          <p className="text-muted-foreground text-sm">
            Manage fulfillment, shipments, and customer orders assigned to this store branch.
          </p>
        </div>
      </div>

      {/* Filter Bar with Shop Filter Hidden */}
      <OrderFilters
        searchNumber={searchNumber}
        statusFilter={statusFilter}
        fromDate={fromDate}
        toDate={toDate}
        isShopLocked={true}
        hideShopFilter={true}
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
        onFromDateChange={(date) => setFromDate(date)}
        onToDateChange={(date) => setToDate(date)}
        onRefresh={refresh}
      />

      {/* Workspace Master-Detail Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">
        {/* Left: Orders List */}
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
            hideShopColumn={true}
          />
        </div>

        {/* Right: Inspector */}
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
  );
};

export default ShopOrdersTab;
