import React from 'react';
import { useNavigate } from 'react-router-dom';
import { useShopOrdersViewModel } from '../../viewmodels/useShopOrdersViewModel';
import OrderFilters from '../orders/OrderFilters';
import OrdersTable from '../orders/OrdersTable';

export interface ShopOrdersTabProps {
  shopId: string;
  shopName?: string;
}

export const ShopOrdersTab: React.FC<ShopOrdersTabProps> = ({ shopId, shopName }) => {
  const navigate = useNavigate();
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

  const handleSort = () => {
    const currentDirection = sort.split(':')[1];
    const newDirection = currentDirection === 'desc' ? 'asc' : 'desc';
    setSort(`latest:${newDirection}`);
    setPage(1);
  };

  const handleSelectOrder = (orderId: string) => {
    const shopNameParam = shopName ? `&shopName=${encodeURIComponent(shopName)}` : '';
    navigate(`/orders/${orderId}?from=shop&shopId=${shopId}${shopNameParam}`);
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
        }}
        onFromDateChange={(date) => setFromDate(date)}
        onToDateChange={(date) => setToDate(date)}
        onRefresh={refresh}
      />

      {/* Full-Width Shop Orders List */}
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
          hideShopColumn={true}
        />
      </div>
    </div>
  );
};

export default ShopOrdersTab;
