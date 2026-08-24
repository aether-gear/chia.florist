import React from 'react';
import { PackageOpen, ArrowUpDown, Clock, Truck, Loader2 } from 'lucide-react';
import { Card, CardContent } from '../ui/card';
import { Skeleton } from '../ui/skeleton';
import EmptyState from '../EmptyState';
import StatusBadge from '../StatusBadge';
import Pagination from '../Pagination';
import type { Order } from '../../models/Order';

export interface OrdersTableProps {
  orders: Order[];
  total: number;
  page: number;
  limit: number;
  sort: string;
  loading: boolean;
  isSwitchingCategory?: boolean;
  error?: string | null;
  selectedOrderId?: string | null;
  onSelectOrder: (orderId: string) => void;
  onSortChange: () => void;
  onPageChange: (page: number) => void;
  hideShopColumn?: boolean;
}

const formatCurrency = (amount: number) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
  }).format(amount);
};

export const getOrderExpiry = (order: Order): { label: string; isUrgent?: boolean; isExpired?: boolean } | null => {
  let expiryStr: string | undefined | null = null;
  if (order.status === 'pending') {
    expiryStr = order.payment?.expires_at;
  } else if (order.status === 'confirmed' || order.status === 'processing') {
    expiryStr = order.handling_expires_at;
    if (!expiryStr && order.confirmed_at) {
      const confTime = new Date(order.confirmed_at).getTime();
      if (!isNaN(confTime)) {
        expiryStr = new Date(confTime + 72 * 60 * 60 * 1000).toISOString();
      }
    }
  }

  if (!expiryStr) return null;
  const target = new Date(expiryStr);
  if (isNaN(target.getTime())) return null;

  const now = Date.now();
  const diffMs = target.getTime() - now;

  if (diffMs <= 0) {
    return {
      label: 'Expired',
      isExpired: true,
      isUrgent: true,
    };
  }

  const oneDayMs = 24 * 60 * 60 * 1000;
  if (diffMs < oneDayMs) {
    const totalSecs = Math.floor(diffMs / 1000);
    const hours = Math.floor(totalSecs / 3600);
    const minutes = Math.floor((totalSecs % 3600) / 60);

    const timeWords =
      hours > 0
        ? `${hours}h ${minutes > 0 ? `${minutes}m` : ''}`.trim() + ' left'
        : minutes > 0
        ? `${minutes}m left`
        : `${totalSecs}s left`;

    return {
      label: timeWords,
      isUrgent: true,
      isExpired: false,
    };
  }

  const formattedDate = target.toLocaleDateString(undefined, {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });

  return {
    label: `Due ${formattedDate}`,
    isUrgent: false,
    isExpired: false,
  };
};

export const OrdersTable: React.FC<OrdersTableProps> = ({
  orders,
  total,
  page,
  limit,
  sort,
  loading,
  isSwitchingCategory = false,
  error = null,
  selectedOrderId,
  onSelectOrder,
  onSortChange,
  onPageChange,
  hideShopColumn = false,
}) => {
  return (
    <div className="flex flex-col space-y-4 w-full">
      {/* Sorting & Order Count Bar */}
      <div className="flex items-center justify-between px-1 text-xs text-muted-foreground">
        <div className="flex items-center gap-2">
          <span>{total ? `Found ${total} orders` : 'No orders found'}</span>
          {isSwitchingCategory && <Loader2 className="h-3 w-3 animate-spin text-primary" />}
        </div>
        <button
          disabled={isSwitchingCategory || loading}
          onClick={onSortChange}
          className={`flex items-center gap-1 hover:text-foreground transition-colors font-medium cursor-pointer ${
            isSwitchingCategory || loading ? 'opacity-50 cursor-not-allowed pointer-events-none' : ''
          }`}
        >
          Date ({sort.endsWith(':asc') ? 'Asc' : 'Desc'}) <ArrowUpDown className="h-3 w-3" />
        </button>
      </div>

      {/* List Container with Cards */}
      <div className="space-y-3 min-h-[250px]">
        {isSwitchingCategory ? (
          <div className="py-24 flex flex-col items-center justify-center gap-3 text-muted-foreground bg-muted/10 rounded-2xl border border-dashed border-border/80 min-h-[320px] select-none pointer-events-none animate-in fade-in duration-150">
            <Loader2 className="h-8 w-8 animate-spin text-primary" />
            <div className="text-center space-y-0.5">
              <p className="text-xs font-bold text-foreground">Loading category orders...</p>
              <p className="text-[11px] text-muted-foreground">Please wait a moment</p>
            </div>
          </div>
        ) : loading && orders.length === 0 ? (
          Array.from({ length: 4 }).map((_, i) => (
            <Card key={`skeleton-${i}`} className="border border-border/50 shadow-none rounded-xl">
              <CardContent className="p-4 space-y-3">
                <div className="flex justify-between items-center">
                  <Skeleton className="h-4 w-28 bg-muted animate-pulse rounded" />
                  <Skeleton className="h-4 w-16 bg-muted animate-pulse rounded" />
                </div>
                <Skeleton className="h-3 w-40 bg-muted animate-pulse rounded" />
                <div className="flex gap-2">
                  <Skeleton className="h-5 w-16 rounded-full bg-muted animate-pulse" />
                  <Skeleton className="h-5 w-16 rounded-full bg-muted animate-pulse" />
                </div>
              </CardContent>
            </Card>
          ))
        ) : error ? (
          <EmptyState
            title="Failed to load orders"
            description={error}
            className="py-12 border-0 bg-transparent text-destructive"
          />
        ) : orders.length === 0 ? (
          <EmptyState
            icon={<PackageOpen className="h-8 w-8 mb-2 mx-auto text-slate-400" />}
            title="No orders found"
            description="No orders found matching current criteria."
            className="py-12 border border-dashed border-border/80 rounded-2xl bg-zinc-50/10"
          />
        ) : (
          orders.map((order) => {
            const isSelected = order.id === selectedOrderId;
            const itemPreview = order.items.map((item) => `${item.product_name} (x${item.quantity})`).join(', ');
            const countShipments =
              order.shipments && order.shipments.length > 0 ? order.shipments.length : order.shipment ? 1 : 0;

            const expiryInfo = getOrderExpiry(order);

            return (
              <Card
                key={order.id}
                onClick={() => onSelectOrder(order.id)}
                className={`cursor-pointer transition-all border shadow-none hover:border-primary/50 select-none rounded-xl ${
                  isSelected
                    ? 'border-primary/60 bg-primary/5 ring-1 ring-primary/45'
                    : 'border-border/60 bg-card hover:bg-muted/10'
                }`}
              >
                <CardContent className="p-4 space-y-2.5">
                  {/* Card Header: Order Number & Date / Expiration */}
                  <div className="flex justify-between items-start gap-2">
                    <span className="font-semibold font-mono text-sm tracking-tight text-foreground">
                      {order.number}
                    </span>
                    <div className="flex items-center gap-1.5 flex-wrap justify-end">
                      {expiryInfo && (
                        <span
                          className={`inline-flex items-center gap-1 px-1.5 py-0.5 rounded text-[10px] font-bold tracking-tight ${
                            expiryInfo.isExpired
                              ? 'bg-destructive/15 text-destructive border border-destructive/30'
                              : expiryInfo.isUrgent
                              ? 'bg-amber-500/15 text-amber-700 dark:text-amber-400 border border-amber-500/30'
                              : 'bg-muted text-muted-foreground'
                          }`}
                          title="Order Expiration Target"
                        >
                          <Clock className="h-3 w-3" />
                          {expiryInfo.label}
                        </span>
                      )}
                      <span className="text-[10px] text-muted-foreground font-mono">
                        {new Date(order.created_at).toLocaleDateString(undefined, { dateStyle: 'short' })}
                      </span>
                    </div>
                  </div>

                  {/* Items Summary Line */}
                  <p className="text-xs text-muted-foreground line-clamp-1">{itemPreview}</p>

                  {/* Status & Package Badges Row */}
                  <div className="flex flex-wrap gap-1.5 pt-0.5">
                    <StatusBadge status={order.status} className="scale-90 origin-left" />
                    {order.payment && (
                      <StatusBadge
                        status={`Pay: ${order.payment.status}`}
                        className="scale-90 origin-left bg-muted text-muted-foreground font-medium"
                      />
                    )}
                    {countShipments > 0 ? (
                      <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-bold bg-primary/10 text-primary">
                        <Truck className="h-3 w-3" /> {countShipments} {countShipments === 1 ? 'pkg' : 'pkgs'}
                      </span>
                    ) : (
                      order.shipment && (
                        <StatusBadge
                          status={`Ship: ${order.shipment.status}`}
                          className="scale-90 origin-left bg-muted text-muted-foreground font-medium"
                        />
                      )
                    )}
                    {!hideShopColumn && order.items[0]?.shop_name && (
                      <span className="text-[10px] px-2 py-0.5 rounded bg-muted text-muted-foreground">
                        {order.items[0].shop_name}
                        {order.items.some((i) => i.shop_id !== order.items[0].shop_id) ? ' (+more)' : ''}
                      </span>
                    )}
                  </div>

                  {/* Card Footer: Truncated ID & Total Price */}
                  <div className="flex justify-between items-center border-t border-border/40 pt-2 mt-1">
                    <span className="text-[10px] text-muted-foreground font-mono">
                      ID: {order.id.slice(0, 8)}...
                    </span>
                    <span className="text-sm font-bold text-primary font-mono">{formatCurrency(order.total)}</span>
                  </div>
                </CardContent>
              </Card>
            );
          })
        )}
      </div>

      {/* Pagination */}
      {total > limit && (
        <div className="pt-2">
          <Pagination
            currentPage={page}
            totalPages={Math.ceil(total / limit)}
            totalItems={total}
            limit={limit}
            onPageChange={onPageChange}
            itemNamePlural="orders"
          />
        </div>
      )}
    </div>
  );
};

export default OrdersTable;
