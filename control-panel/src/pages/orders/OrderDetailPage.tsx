import { useState, useEffect, useCallback } from 'react';
import { useParams, useSearchParams, useNavigate } from 'react-router-dom';
import {
  ArrowLeft,
  PackageOpen,
  CheckCircle2,
  XCircle,
  Truck,
  CreditCard,
  User,
  RefreshCw,
  Edit3,
  Activity,
  Sprout,
  Layers,
  Loader2,
  AlertTriangle,
  Clock,
  ArrowRight,
  Warehouse,
  Terminal,
} from 'lucide-react';
import { Button } from '../../components/ui/button';
import { Input } from '../../components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '../../components/ui/select';
import { Skeleton } from '../../components/ui/skeleton';
import StatusBadge from '../../components/StatusBadge';
import Breadcrumb from '../../components/Breadcrumb';
import EmptyState from '../../components/EmptyState';
import { useAuthMeViewModel } from '../../viewmodels/useAuthMeViewModel';
import { useOrderDetailViewModel } from '../../viewmodels/useOrderDetailViewModel';
import { useOrderActionsViewModel, type ShipmentDispatchPayload, type OrderTrackingResponse } from '../../viewmodels/useOrderActionsViewModel';
import ShipmentDispatchSection from '../../components/orders/ShipmentDispatchSection';
import OrderItemDetailOverlay from '../../components/orders/OrderItemDetailOverlay';
import { getOrderExpiry } from '../../components/orders/OrdersTable';
import type { OrderItem } from '../../models/Order';

const formatCurrency = (amount: number) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
  }).format(amount);
};

export default function OrderDetailPage() {
  const { orderId } = useParams<{ orderId: string }>();
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const { isAdmin } = useAuthMeViewModel();

  const fromContext = searchParams.get('from');
  const contextShopId = searchParams.get('shopId');
  const contextShopName = searchParams.get('shopName');

  const { order, loading, error, refresh } = useOrderDetailViewModel(orderId);
  const {
    submitting,
    fetchOrderTracking,
    updateOrderStatus,
    dispatchShopShipment,
    updateShipmentStatus,
    updateShipmentDetails,
  } = useOrderActionsViewModel();

  // Item inspection drawer state
  const [selectedItemForDetail, setSelectedItemForDetail] = useState<OrderItem | null>(null);
  const [isItemDetailOpen, setIsItemDetailOpen] = useState<boolean>(false);

  // Per-shipment waybill editing state
  const [editingShipmentId, setEditingShipmentId] = useState<string | null>(null);
  const [editCourier, setEditCourier] = useState<string>('');
  const [editService, setEditService] = useState<string>('');
  const [editTracking, setEditTracking] = useState<string>('');

  // Per-shipment transit updating state
  const [activeShipmentId, setActiveShipmentId] = useState<string | null>(null);
  const [transitStatus, setTransitStatus] = useState<string>('packed');
  const [transitDescription, setTransitDescription] = useState<string>('');
  const [transitLocation, setTransitLocation] = useState<string>('');

  // Live external courier tracking state & cooldown
  const [liveTracking, setLiveTracking] = useState<OrderTrackingResponse | null>(null);
  const [loadingLiveTracking, setLoadingLiveTracking] = useState<boolean>(false);
  const [syncCooldown, setSyncCooldown] = useState<number>(0);

  // Sync state when order changes
  useEffect(() => {
    if (order) {
      const activeShipments =
        order.shipments && order.shipments.length > 0
          ? order.shipments
          : order.shipment
          ? [order.shipment]
          : [];

      if (activeShipments.length > 0) {
        setActiveShipmentId(activeShipments[0].id);
        setTransitStatus(activeShipments[0].status || 'packed');
      } else {
        setActiveShipmentId(null);
      }

      setLiveTracking(null);
      setEditingShipmentId(null);
      setTransitDescription('');
      setTransitLocation('');
    } else {
      setLiveTracking(null);
    }
  }, [order?.id, order?.status]);

  // Cooldown countdown timer for sync button
  useEffect(() => {
    if (syncCooldown <= 0) return;
    const timer = setInterval(() => {
      setSyncCooldown((prev) => (prev > 0 ? prev - 1 : 0));
    }, 1000);
    return () => clearInterval(timer);
  }, [syncCooldown]);

  const loadLiveTracking = useCallback(
    async (targetOrderId: string) => {
      if (syncCooldown > 0) return;
      setLoadingLiveTracking(true);
      try {
        const res = await fetchOrderTracking(targetOrderId);
        setLiveTracking(res);
        setSyncCooldown(15);
      } catch (err) {
        console.error('Failed to load live tracking', err);
        setLiveTracking(null);
      } finally {
        setLoadingLiveTracking(false);
      }
    },
    [fetchOrderTracking, syncCooldown]
  );

  const handleOpenItemDetail = (item: OrderItem) => {
    setSelectedItemForDetail(item);
    setIsItemDetailOpen(true);
  };

  const handleStartEditingWaybill = (shipment: any) => {
    setEditingShipmentId(shipment.id);
    setEditCourier(shipment.courier || '');
    setEditService(shipment.service || '');
    setEditTracking(shipment.tracking_number || '');
  };

  const handleSaveWaybill = async (shipmentId: string) => {
    await updateShipmentDetails(shipmentId, {
      tracking_number: editTracking,
      courier: editCourier,
      service: editService,
    });
    setEditingShipmentId(null);
    refresh();
  };

  const handleStartProcessing = async () => {
    if (!order) return;
    await updateOrderStatus(order.id, 'processing');
    refresh();
  };

  const handleDispatchOrder = async (targetOrderId: string, shipments: ShipmentDispatchPayload[]) => {
    await updateOrderStatus(targetOrderId, 'shipped', undefined, undefined, shipments);
    refresh();
  };

  const handleDispatchShopShipment = async (
    targetOrderId: string,
    payload: {
      shop_id: string;
      fulfillment_method: string;
      courier: string;
      service: string;
      tracking_number?: string;
      item_ids: string[];
    }
  ) => {
    await dispatchShopShipment(targetOrderId, payload);
    refresh();
  };

  const handleUpdateShipmentStatus = async (shipmentId: string, status: string) => {
    await updateShipmentStatus(shipmentId, status);
    refresh();
  };

  const handleBack = () => {
    if (fromContext === 'shop' && contextShopId) {
      navigate(`/shop/${contextShopId}?tab=orders`);
    } else {
      navigate('/orders');
    }
  };

  // Breadcrumb generation
  const breadcrumbItems =
    fromContext === 'shop' && contextShopId
      ? [
          { label: 'Shop Management', href: '/shop' },
          {
            label: contextShopName || 'Shop Details',
            href: `/shop/${contextShopId}?tab=orders`,
          },
          { label: 'Orders', href: `/shop/${contextShopId}?tab=orders` },
          { label: order?.number ? `Order ${order.number}` : 'Order Detail' },
        ]
      : [
          { label: 'Orders', href: '/orders' },
          { label: order?.number ? `Order ${order.number}` : 'Order Detail' },
        ];

  if (loading) {
    return (
      <div className="flex-1 space-y-8 p-6 sm:p-8 xl:p-12 animate-in fade-in duration-200">
        <div className="space-y-4">
          <Skeleton className="h-5 w-48 bg-muted rounded" />
          <div className="flex justify-between items-center">
            <Skeleton className="h-10 w-64 bg-muted rounded-xl" />
            <Skeleton className="h-10 w-28 bg-muted rounded-xl" />
          </div>
        </div>
        <div className="grid grid-cols-1 xl:grid-cols-12 gap-8">
          <div className="xl:col-span-8 space-y-6">
            <Skeleton className="h-48 w-full bg-muted rounded-2xl" />
            <Skeleton className="h-64 w-full bg-muted rounded-2xl" />
          </div>
          <div className="xl:col-span-4 space-y-6">
            <Skeleton className="h-48 w-full bg-muted rounded-2xl" />
            <Skeleton className="h-48 w-full bg-muted rounded-2xl" />
          </div>
        </div>
      </div>
    );
  }

  if (error || !order) {
    return (
      <div className="flex-1 space-y-8 p-6 sm:p-8 xl:p-12">
        <Breadcrumb items={breadcrumbItems} />
        <div className="py-16 text-center">
          <EmptyState
            icon={<PackageOpen className="h-12 w-12 mb-3 mx-auto text-muted-foreground/60 stroke-[1.5]" />}
            title="Order Not Found"
            description={error || 'The requested order does not exist or you do not have permission to view it.'}
            actionLabel="Back to Orders"
            onAction={handleBack}
          />
        </div>
      </div>
    );
  }

  const currentShipments =
    order.shipments && order.shipments.length > 0 ? order.shipments : order.shipment ? [order.shipment] : [];
  const expiryInfo = getOrderExpiry(order);

  const matchedShipmentForItem = (item: OrderItem) => {
    const idx = currentShipments.findIndex(
      (s) => s.id === item.shipment_id || s.item_ids?.includes(item.id)
    );
    if (idx >= 0) {
      return { shipment: currentShipments[idx], index: idx };
    }
    return null;
  };

  const selectedItemMatchedShipment = selectedItemForDetail
    ? matchedShipmentForItem(selectedItemForDetail)
    : null;

  return (
    <div className="flex-col md:flex">
      <div className="flex-1 space-y-8 p-6 sm:p-8 xl:p-12 animate-in fade-in duration-300">
        {/* 1. Breadcrumb Navigation */}
        <Breadcrumb items={breadcrumbItems} className="text-sm font-medium" />

        {/* 2. Page Header */}
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-6 border-b border-border/60">
          <div className="flex items-start sm:items-center gap-4">
            <Button
              variant="outline"
              size="icon"
              onClick={handleBack}
              className="rounded-xl h-10 w-10 text-muted-foreground hover:text-foreground shrink-0 border-border hover:bg-primary/5 cursor-pointer"
              title="Go back"
            >
              <ArrowLeft className="h-5 w-5" />
            </Button>
            <div className="space-y-1.5">
              <div className="flex items-center gap-3 flex-wrap">
                <h2 className="text-3xl sm:text-4xl font-bold font-mono tracking-tight text-foreground">
                  {order.number}
                </h2>
                <StatusBadge status={order.status} className="scale-110 origin-left" />
                {expiryInfo && (
                  <span
                    className={`inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-bold tracking-tight ${
                      expiryInfo.isExpired
                        ? 'bg-destructive/15 text-destructive border border-destructive/30'
                        : expiryInfo.isUrgent
                        ? 'bg-amber-500/15 text-amber-700 dark:text-amber-400 border border-amber-500/30'
                        : 'bg-muted text-muted-foreground'
                    }`}
                    title="Order Expiration Target"
                  >
                    <Clock className="h-3.5 w-3.5" />
                    {expiryInfo.label}
                  </span>
                )}
              </div>
              <p className="text-muted-foreground text-sm font-sans">
                Placed on {new Date(order.created_at).toLocaleString()}
              </p>
            </div>
          </div>

          <div className="flex items-center gap-2 self-start sm:self-auto">
            <Button
              variant="outline"
              size="sm"
              onClick={() => refresh()}
              disabled={submitting}
              className="rounded-xl text-sm h-10 px-4 gap-2 border-border text-foreground hover:text-primary hover:bg-primary/5 cursor-pointer"
            >
              <RefreshCw className={`h-4 w-4 ${submitting ? 'animate-spin' : ''}`} />
              Refresh
            </Button>
          </div>
        </div>

        {/* 3. State Transition Action Bar */}
        <div className="space-y-4">
          {/* CASE A: PENDING PAYMENT */}
          {order.status === 'pending' && (
            <div className="p-5 rounded-2xl border border-amber-500/20 bg-amber-500/5 text-amber-900 dark:text-amber-200 text-sm flex items-center justify-between">
              <span className="font-medium">Awaiting payment settlement from customer. Status will auto-advance once paid.</span>
            </div>
          )}

          {/* CASE B: CONFIRMED - READY FOR PACKAGING */}
          {order.status === 'confirmed' && (
            <div className="p-5 rounded-2xl border border-blue-500/20 bg-blue-500/5 text-sm text-blue-950 dark:text-blue-200 flex flex-col sm:flex-row sm:items-center justify-between gap-4">
              <span className="font-medium">Payment verified. Ready to initiate order assembly and warehouse processing.</span>
              <Button
                size="sm"
                disabled={submitting}
                onClick={handleStartProcessing}
                className="bg-primary text-primary-foreground hover:bg-primary/90 text-sm rounded-xl h-10 px-4 font-semibold cursor-pointer shrink-0"
              >
                <Sprout className="h-4 w-4 mr-1.5" /> Start Processing
              </Button>
            </div>
          )}

          {/* CASE C: PROCESSING - MULTI / SPLIT SHIPMENT DISPATCH */}
          {order.status === 'processing' && (
            <ShipmentDispatchSection
              order={order}
              submitting={submitting}
              shopId={contextShopId || undefined}
              onDispatchOrder={handleDispatchOrder}
              onDispatchShopShipment={handleDispatchShopShipment}
            />
          )}

          {/* CASE E: DELIVERED */}
          {order.status === 'delivered' && (
            <div className="p-5 rounded-2xl border border-emerald-500/20 bg-emerald-500/5 text-emerald-950 dark:text-emerald-200 text-sm flex items-center gap-2.5">
              <CheckCircle2 className="h-5 w-5 text-emerald-600 shrink-0" />
              <span className="font-medium">Order fulfillment complete. All packages delivered to customer.</span>
            </div>
          )}

          {/* CASE F: CANCELLED */}
          {order.status === 'cancelled' && (
            <div className="p-5 rounded-2xl border border-destructive/20 bg-destructive/5 text-destructive text-sm flex items-center gap-2.5">
              <XCircle className="h-5 w-5 text-destructive shrink-0" />
              <span className="font-medium">This order was cancelled. Reserved inventory was released back to stock.</span>
            </div>
          )}
        </div>

        {/* 4. Detailed Two-Column Grid (xl:grid-cols-12 for spacious wide screens) */}
        <div className="grid grid-cols-1 xl:grid-cols-12 gap-8 items-start">
          {/* Left / Main Column (Items + Logistics + Tracking) */}
          <div className="xl:col-span-8 space-y-8">
            {/* Clean Order Items Cards List */}
            <div className="border border-border/80 rounded-2xl bg-card p-6 sm:p-8 space-y-6 shadow-sm">
              <div className="flex items-center justify-between pb-4 border-b border-border/60">
                <div>
                  <h4 className="text-xl font-bold text-foreground font-display flex items-center gap-2.5">
                    <PackageOpen className="h-5 w-5 text-primary" /> Order Items ({order.items.length})
                  </h4>
                  <p className="text-sm text-muted-foreground mt-0.5">
                    Click any item to inspect its complete custom specifications, inscription texts, and logistics allocation.
                  </p>
                </div>
              </div>

              {/* Individual Item Cards */}
              <div className="space-y-3.5">
                {order.items.map((item, idx) => {
                  const match = matchedShipmentForItem(item);
                  const isMyItem = contextShopId ? item.shop_id === contextShopId : false;
                  const isCustom = item.is_custom || item.product_variant_type === 'custom' || !item.product_id;

                  return (
                    <div
                      key={item.id}
                      onClick={() => handleOpenItemDetail(item)}
                      className={`group cursor-pointer p-4 sm:p-5 rounded-2xl border transition-all duration-200 hover:border-primary/60 hover:shadow-sm ${
                        isMyItem ? 'bg-primary/5 border-primary/30' : 'bg-muted/10 border-border/60 hover:bg-muted/20'
                      }`}
                    >
                      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
                        {/* Item Title & Specs */}
                        <div className="space-y-2 flex-1 min-w-0">
                          <div className="flex items-center gap-2 flex-wrap">
                            <span className="text-xs font-mono font-bold text-muted-foreground">
                              #{idx + 1}
                            </span>
                            <h5 className="font-semibold text-base text-foreground group-hover:text-primary transition-colors truncate">
                              {item.product_name}
                            </h5>
                            {isCustom && (
                              <span className="px-2.5 py-0.5 rounded-full text-xs font-semibold bg-purple-500/10 text-purple-700 dark:text-purple-300 border border-purple-500/20">
                                Custom Board
                              </span>
                            )}
                            <span className="px-2.5 py-0.5 rounded-full text-xs font-semibold bg-primary/10 text-primary border border-primary/20 flex items-center gap-1">
                              <Warehouse className="h-3 w-3" />
                              {item.shop_name}
                            </span>
                          </div>

                          <div className="flex items-center gap-3 text-xs text-muted-foreground flex-wrap">
                            {item.item_options?.size && (
                              <span className="font-medium text-foreground">
                                Size: <strong className="capitalize">{item.item_options.size}</strong>
                              </span>
                            )}
                            {item.item_options?.jambul && (
                              <span>
                                Jambul: <strong className="capitalize">{item.item_options.jambul}</strong>
                              </span>
                            )}
                            {item.courier_code && (
                              <span className="uppercase font-medium">
                                Courier: {item.courier_code} {item.courier_service}
                              </span>
                            )}
                          </div>
                        </div>

                        {/* Package Tag & Pricing & Hover Arrow */}
                        <div className="flex items-center justify-between sm:justify-end gap-4 shrink-0 border-t sm:border-t-0 pt-3 sm:pt-0 border-border/40">
                          <div className="text-left sm:text-right space-y-0.5">
                            <div className="flex items-center sm:justify-end gap-2">
                              {match ? (
                                <span className="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-bold bg-primary/10 text-primary">
                                  <Truck className="h-3.5 w-3.5" /> Pkg #{match.index + 1}
                                </span>
                              ) : (
                                <span className="text-xs text-muted-foreground italic">Unassigned Pkg</span>
                              )}
                              <span className="text-xs text-muted-foreground font-semibold">
                                × {item.quantity}
                              </span>
                            </div>
                            <div className="font-mono font-bold text-sm text-foreground">
                              {formatCurrency(item.subtotal)}
                            </div>
                          </div>

                          <div className="w-8 h-8 rounded-full bg-primary/10 text-primary flex items-center justify-center opacity-0 group-hover:opacity-100 transition-all duration-200 group-hover:translate-x-0.5 shrink-0">
                            <ArrowRight className="h-4 w-4" />
                          </div>
                        </div>
                      </div>
                    </div>
                  );
                })}
              </div>
            </div>

            {/* Active Logistics & Waybill Management */}
            {currentShipments.length > 0 && (
              <div className="border border-border/80 rounded-2xl bg-card p-6 sm:p-8 space-y-6 shadow-sm">
                <div className="space-y-1 pb-4 border-b border-border/60">
                  <h4 className="text-xl font-bold text-foreground font-display flex items-center gap-2.5">
                    <Layers className="h-5 w-5 text-primary" /> Active Logistics & Packages ({currentShipments.length})
                  </h4>
                  <p className="text-sm text-muted-foreground">
                    Manage waybill tracking numbers and post transit timeline updates per package.
                  </p>
                </div>

                <div className="space-y-4">
                  {currentShipments.map((shipment, sIdx) => {
                    const isEditing = editingShipmentId === shipment.id;
                    const isSelectedForTimeline = activeShipmentId === shipment.id;

                    const assignedItems = order.items.filter(
                      (item) => item.shipment_id === shipment.id || shipment.item_ids?.includes(item.id)
                    );

                    const isMyShopShipment = !contextShopId || assignedItems.some((i) => i.shop_id === contextShopId);

                    return (
                      <div
                        key={shipment.id}
                        className={`border rounded-xl p-5 transition-all space-y-4 ${
                          isSelectedForTimeline ? 'border-primary/60 bg-muted/20 shadow-xs' : 'border-border/60 bg-muted/5'
                        }`}
                      >
                        <div className="flex flex-col sm:flex-row justify-between sm:items-center gap-3 pb-3 border-b border-border/40">
                          <div>
                            <div className="flex items-center gap-2.5 flex-wrap">
                              <span className="text-sm font-bold text-foreground">
                                Package #{sIdx + 1}: {shipment.courier.toUpperCase()} {shipment.service.toUpperCase()}
                              </span>
                              <span className="px-2.5 py-0.5 rounded-full text-xs font-bold bg-primary/10 text-primary uppercase">
                                {shipment.status}
                              </span>
                              {!isMyShopShipment && (
                                <span className="px-2.5 py-0.5 rounded-md text-xs font-semibold bg-muted text-muted-foreground">
                                  Other Shop Branch
                                </span>
                              )}
                            </div>
                            <div className="text-xs text-muted-foreground mt-1">
                              Waybill Tracking:{' '}
                              <strong className="font-mono text-foreground text-sm">
                                {shipment.tracking_number || 'Self Delivery'}
                              </strong>{' '}
                              · {assignedItems.length} item{assignedItems.length !== 1 ? 's' : ''}
                            </div>
                          </div>

                          {isMyShopShipment && (
                            <div className="flex items-center gap-2">
                              <Button
                                variant="outline"
                                size="sm"
                                onClick={() => {
                                  if (isEditing) {
                                    setEditingShipmentId(null);
                                  } else {
                                    handleStartEditingWaybill(shipment);
                                  }
                                }}
                                className="text-xs h-8 px-3 text-primary hover:bg-primary/5 rounded-lg flex items-center gap-1.5 cursor-pointer border-primary/30"
                              >
                                <Edit3 className="h-3.5 w-3.5" /> {isEditing ? 'Cancel Editing' : 'Edit Waybill'}
                              </Button>
                            </div>
                          )}
                        </div>

                        {isMyShopShipment ? (
                          isEditing ? (
                            /* EDIT WAYBILL SUB-FORM */
                            <div className="space-y-4 p-4 border border-border/60 rounded-xl bg-background">
                              <h6 className="text-sm font-bold text-foreground">Edit Shipping Waybill</h6>
                              <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                                <div className="space-y-1.5">
                                  <label className="text-xs font-semibold text-muted-foreground uppercase">
                                    Courier Brand
                                  </label>
                                  <Input
                                    value={editCourier}
                                    onChange={(e) => setEditCourier(e.target.value)}
                                    className="h-9 text-sm rounded-lg"
                                  />
                                </div>
                                <div className="space-y-1.5">
                                  <label className="text-xs font-semibold text-muted-foreground uppercase">
                                    Service Tier
                                  </label>
                                  <Input
                                    value={editService}
                                    onChange={(e) => setEditService(e.target.value)}
                                    className="h-9 text-sm rounded-lg"
                                  />
                                </div>
                                <div className="space-y-1.5">
                                  <label className="text-xs font-semibold text-muted-foreground uppercase">
                                    Tracking Number
                                  </label>
                                  <Input
                                    value={editTracking}
                                    onChange={(e) => setEditTracking(e.target.value)}
                                    className="h-9 text-sm rounded-lg font-mono"
                                  />
                                </div>
                              </div>
                              <Button
                                size="sm"
                                disabled={submitting}
                                onClick={() => handleSaveWaybill(shipment.id)}
                                className="h-9 text-sm rounded-lg w-full bg-primary text-primary-foreground hover:bg-primary/90 mt-2 font-semibold cursor-pointer"
                              >
                                Save Waybill Changes
                              </Button>
                            </div>
                          ) : (
                            /* TRANSIT STATE SELECTOR */
                            <div className="space-y-4 pt-1">
                              <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                                <div className="space-y-1.5">
                                  <label className="text-xs font-semibold text-muted-foreground">
                                    Update Transit Status
                                  </label>
                                  <Select
                                    value={transitStatus}
                                    onValueChange={(val) => {
                                      setActiveShipmentId(shipment.id);
                                      setTransitStatus(val);
                                    }}
                                  >
                                    <SelectTrigger className="rounded-xl h-9 text-sm bg-background">
                                      <SelectValue />
                                    </SelectTrigger>
                                    <SelectContent>
                                      <SelectItem value="packed">Packed</SelectItem>
                                      <SelectItem value="labelled">Labelled</SelectItem>
                                      <SelectItem value="picked_up">Picked Up</SelectItem>
                                      <SelectItem value="in_transit">In Transit</SelectItem>
                                      <SelectItem value="out_for_delivery">Out For Delivery</SelectItem>
                                      <SelectItem value="delivered">Delivered (Fulfillment Complete)</SelectItem>
                                      <SelectItem value="failed">Failed Delivery</SelectItem>
                                      <SelectItem value="returned">Returned to Shop</SelectItem>
                                      <SelectItem value="cancelled">Cancelled</SelectItem>
                                    </SelectContent>
                                  </Select>
                                </div>
                                <div className="space-y-1.5 sm:col-span-2">
                                  <label className="text-xs font-semibold text-muted-foreground">
                                    Location (Optional)
                                  </label>
                                  <Input
                                    value={transitLocation}
                                    onChange={(e) => setTransitLocation(e.target.value)}
                                    placeholder="Warehouse, Hub City, Distribution Center..."
                                    className="rounded-xl h-9 text-sm bg-background"
                                  />
                                </div>
                              </div>
                              <div className="space-y-1.5">
                                <label className="text-xs font-semibold text-muted-foreground">
                                  Description / Checkpoint Note (Optional)
                                </label>
                                <Input
                                  value={transitDescription}
                                  onChange={(e) => setTransitDescription(e.target.value)}
                                  placeholder="E.g. Package handed over to courier driver"
                                  className="rounded-xl h-9 text-sm bg-background"
                                />
                              </div>
                              <Button
                                disabled={submitting}
                                onClick={() => handleUpdateShipmentStatus(shipment.id, transitStatus)}
                                className="rounded-xl bg-primary text-primary-foreground hover:bg-primary/90 text-sm font-semibold h-9 w-full cursor-pointer"
                              >
                                Update Package #{sIdx + 1} Status
                              </Button>
                            </div>
                          )
                        ) : (
                          <div className="mt-2 text-xs text-muted-foreground italic">
                            This package is fulfilled and managed by another store branch.
                          </div>
                        )}
                      </div>
                    );
                  })}
                </div>
              </div>
            )}

            {/* Shipment Tracking History Timeline */}
            {currentShipments.length > 0 && (
              <div className="border border-border/80 rounded-2xl bg-card p-6 sm:p-8 space-y-6 shadow-sm">
                <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b border-border/60">
                  <div>
                    <div className="flex items-center gap-2.5 flex-wrap">
                      <h4 className="text-xl font-bold text-foreground font-display flex items-center gap-2.5">
                        <Activity className="h-5 w-5 text-primary" /> Shipment Tracking History
                      </h4>
                      {liveTracking?.timeline && liveTracking.timeline.length > 0 && (
                        <span className="inline-flex items-center gap-1.5 px-3 py-0.5 rounded-full text-xs font-semibold bg-emerald-500/10 text-emerald-600 border border-emerald-500/20">
                          Live Courier Feed
                        </span>
                      )}
                    </div>
                    <p className="text-sm text-muted-foreground mt-0.5">
                      Chronological transit milestones and courier tracking checkpoints.
                    </p>
                  </div>
                  <div className="flex items-center gap-2.5 self-start sm:self-auto flex-wrap">
                    {loadingLiveTracking && <Loader2 className="h-4 w-4 animate-spin text-primary" />}
                    <Button
                      variant="outline"
                      size="sm"
                      disabled={loadingLiveTracking || syncCooldown > 0}
                      onClick={() => order && loadLiveTracking(order.id)}
                      className="h-8 px-3 text-xs font-semibold rounded-xl gap-1.5 border-primary/30 text-primary hover:bg-primary/5 cursor-pointer disabled:opacity-60"
                    >
                      <RefreshCw className={`h-3.5 w-3.5 ${loadingLiveTracking ? 'animate-spin' : ''}`} />
                      {syncCooldown > 0 ? `Synced (${syncCooldown}s)` : 'Sync Courier Feed'}
                    </Button>
                    {currentShipments.length > 1 && (
                      <div className="flex gap-1.5">
                        {currentShipments.map((s, idx) => (
                          <button
                            key={s.id}
                            onClick={() => setActiveShipmentId(s.id)}
                            className={`px-3 py-1.5 rounded-xl text-xs font-bold transition-all cursor-pointer ${
                              activeShipmentId === s.id
                                ? 'bg-primary text-primary-foreground shadow-xs'
                                : 'bg-muted text-muted-foreground hover:text-foreground'
                            }`}
                          >
                            Pkg #{idx + 1}
                          </button>
                        ))}
                      </div>
                    )}
                  </div>
                </div>

                {liveTracking?.warning && (
                  <div className="flex items-start gap-3 p-4 rounded-xl bg-amber-500/10 border border-amber-500/20 text-amber-800 dark:text-amber-300 text-sm font-medium leading-relaxed">
                    <AlertTriangle className="h-5 w-5 shrink-0 text-amber-500 mt-0.5" />
                    <div className="space-y-0.5">
                      <div className="font-semibold text-foreground">Manual Status Updates Required</div>
                      <p>
                        {liveTracking.warning.toLowerCase().includes('unsupported') ||
                        liveTracking.warning.toLowerCase().includes('not supported') ||
                        liveTracking.warning.toLowerCase().includes('invalid') ||
                        liveTracking.warning.toLowerCase().includes('unavailable')
                          ? 'Live tracking is not supported for this courier. Please update the package transit status manually using the shipment timeline controls.'
                          : liveTracking.warning}
                      </p>
                    </div>
                  </div>
                )}

                {(() => {
                  const activeShipment =
                    currentShipments.find((s) => s.id === activeShipmentId) || currentShipments[0];

                  const liveEvents = liveTracking?.timeline || [];
                  const internalEvents = activeShipment?.events || [];

                  const eventsToDisplay =
                    liveEvents.length > 0
                      ? liveEvents.map((e, idx) => ({
                          id: `live-event-${idx}`,
                          status: e.status,
                          description: e.description,
                          location: e.location,
                          timestamp: e.timestamp,
                        }))
                      : internalEvents;

                  if (eventsToDisplay.length > 0) {
                    return (
                      <div className="relative pl-7 border-l-2 border-primary/20 space-y-7 ml-3 py-2">
                        {eventsToDisplay.map((event, idx) => {
                          const isLatest = idx === eventsToDisplay.length - 1 || idx === 0;
                          return (
                            <div key={event.id} className="relative">
                              <span
                                className={`absolute -left-[35px] top-1 h-4.5 w-4.5 rounded-full border-2 bg-background flex items-center justify-center transition-colors ${
                                  isLatest ? 'border-primary ring-4 ring-primary/10' : 'border-muted-foreground/45'
                                }`}
                              >
                                <span
                                  className={`h-2 w-2 rounded-full ${
                                    isLatest ? 'bg-primary' : 'bg-muted-foreground/45'
                                  }`}
                                />
                              </span>

                              <div className="space-y-1.5">
                                <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-1">
                                  <div className="flex items-center gap-2.5">
                                    <span
                                      className={`text-sm font-bold uppercase tracking-wider ${
                                        isLatest ? 'text-primary' : 'text-muted-foreground'
                                      }`}
                                    >
                                      {event.status}
                                    </span>
                                    {event.location && (
                                      <span className="text-xs bg-muted px-2 py-0.5 rounded-md text-muted-foreground font-semibold">
                                        {event.location}
                                      </span>
                                    )}
                                  </div>
                                  <span className="text-xs text-muted-foreground font-mono">
                                    {new Date(event.timestamp).toLocaleString(undefined, {
                                      dateStyle: 'medium',
                                      timeStyle: 'short',
                                    })}
                                  </span>
                                </div>
                                <p className="text-sm text-muted-foreground leading-relaxed">{event.description}</p>
                              </div>
                            </div>
                          );
                        })}
                      </div>
                    );
                  }

                  return (
                    <div className="p-8 border border-dashed border-border/80 rounded-2xl bg-muted/10 text-center text-sm text-muted-foreground flex flex-col items-center gap-2">
                      <Activity className="h-8 w-8 text-muted-foreground/50 stroke-[1.5]" />
                      <span className="font-semibold text-foreground">
                        No shipment updates logged yet for{' '}
                        {activeShipment
                          ? `${activeShipment.courier.toUpperCase()} (${
                              activeShipment.tracking_number || 'Self Delivery'
                            })`
                          : 'this package'}
                      </span>
                      <span className="text-xs text-muted-foreground max-w-sm">
                        Use the logistics panel above to post the first shipment status update.
                      </span>
                    </div>
                  );
                })()}
              </div>
            )}
          </div>

          {/* Right / Sidebar Column (Billing + Address + Metadata) */}
          <div className="xl:col-span-4 space-y-6">
            {/* Payment & Billing Summary */}
            <div className="border border-border/80 rounded-2xl bg-card p-6 sm:p-7 space-y-5 shadow-sm">
              <h5 className="text-sm font-bold uppercase tracking-wider text-muted-foreground flex items-center gap-2">
                <CreditCard className="h-4 w-4 text-primary" /> Billing Details
              </h5>
              <div className="space-y-3.5 text-sm">
                <div className="flex justify-between text-muted-foreground">
                  <span>Items Subtotal</span>
                  <span className="font-semibold text-foreground font-mono text-sm">{formatCurrency(order.subtotal)}</span>
                </div>
                <div className="flex justify-between text-muted-foreground">
                  <span>Shipping Fee</span>
                  <span className="font-semibold text-foreground font-mono text-sm">{formatCurrency(order.shipping_fee)}</span>
                </div>
                <div className="border-t border-border/60 pt-3.5 flex justify-between items-center text-base font-bold text-foreground">
                  <span>Grand Total</span>
                  <span className="text-primary font-mono text-xl">{formatCurrency(order.total)}</span>
                </div>

                {order.payment && (
                  <div className="border-t border-border/40 pt-4 space-y-2.5 text-xs bg-muted/20 -mx-6 sm:-mx-7 -mb-6 sm:-mb-7 p-5 rounded-b-2xl">
                    <div className="flex justify-between">
                      <span className="text-muted-foreground font-medium">Provider:</span>
                      <span className="text-foreground font-bold uppercase text-xs">{order.payment.provider}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-muted-foreground font-medium">Payment Status:</span>
                      <span className="text-foreground font-bold uppercase text-xs">{order.payment.status}</span>
                    </div>
                    {order.payment.expires_at && (
                      <div className="flex justify-between">
                        <span className="text-muted-foreground font-medium">Payment Due:</span>
                        <span className="text-foreground font-mono text-xs font-semibold">
                          {new Date(order.payment.expires_at).toLocaleString(undefined, { dateStyle: 'short', timeStyle: 'short' })}
                        </span>
                      </div>
                    )}
                  </div>
                )}
              </div>
            </div>

            {/* Shipping & Recipient Address */}
            <div className="border border-border/80 rounded-2xl bg-card p-6 sm:p-7 space-y-5 shadow-sm">
              <h5 className="text-sm font-bold uppercase tracking-wider text-muted-foreground flex items-center gap-2">
                <User className="h-4 w-4 text-primary" /> Shipping & Recipient
              </h5>
              <div className="space-y-4 text-sm">
                {order.address ? (
                  <div className="text-foreground space-y-3 leading-relaxed text-sm">
                    <div>
                      <span className="text-xs uppercase font-bold text-muted-foreground block mb-0.5">Receiver Name</span>
                      <span className="font-semibold text-base text-foreground">{order.address.receiver_name}</span>
                    </div>
                    <div>
                      <span className="text-xs uppercase font-bold text-muted-foreground block mb-0.5">Phone Number</span>
                      <span className="font-mono text-sm text-foreground">{order.address.phone}</span>
                    </div>
                    <div>
                      <span className="text-xs uppercase font-bold text-muted-foreground block mb-0.5">Full Destination Address</span>
                      <p className="text-muted-foreground text-sm leading-relaxed">{order.address.full_address}</p>
                    </div>
                    {order.address.postal_code && (
                      <div>
                        <span className="text-xs uppercase font-bold text-muted-foreground block mb-0.5">Postal Code</span>
                        <span className="font-mono text-sm text-foreground">{order.address.postal_code}</span>
                      </div>
                    )}
                  </div>
                ) : (
                  <div className="text-muted-foreground space-y-3 leading-relaxed text-sm">
                    <div>
                      <span className="text-xs uppercase font-bold text-muted-foreground block mb-0.5">Customer ID</span>
                      <span className="font-mono text-xs text-foreground">{order.customer_id || order.user_id}</span>
                    </div>
                    <div>
                      <span className="text-xs uppercase font-bold text-muted-foreground block mb-0.5">Address ID</span>
                      <span className="font-mono text-xs text-foreground">{order.address_id}</span>
                    </div>
                  </div>
                )}

                {isAdmin && (
                  <div className="pt-4 border-t border-border/60 space-y-1.5">
                    <div className="flex items-center justify-between text-xs text-muted-foreground">
                      <span className="font-mono font-bold uppercase tracking-wider text-[11px] flex items-center gap-1.5 text-muted-foreground">
                        <Terminal className="h-3.5 w-3.5 text-primary" /> Internal System Order UUID
                      </span>
                      <span className="text-[10px] px-1.5 py-0.5 rounded bg-amber-500/10 text-amber-600 dark:text-amber-400 font-semibold border border-amber-500/20">
                        Admin Only
                      </span>
                    </div>
                    <div className="flex items-center gap-2 bg-muted/40 border border-border/60 rounded-xl p-2.5">
                      <code className="text-[11px] font-mono text-foreground break-all select-all flex-1">
                        {order.id}
                      </code>
                    </div>
                  </div>
                )}
              </div>
            </div>

            {/* Order Lifecycle Metadata */}
            <div className="border border-border/80 rounded-2xl bg-card p-6 sm:p-7 space-y-4 text-sm shadow-sm">
              <h5 className="text-sm font-bold uppercase tracking-wider text-muted-foreground">Order Metadata</h5>
              <div className="space-y-2.5 text-xs text-muted-foreground">
                <div className="flex justify-between">
                  <span>Created:</span>
                  <span className="text-foreground font-mono font-medium">
                    {new Date(order.created_at).toLocaleString(undefined, { dateStyle: 'medium', timeStyle: 'short' })}
                  </span>
                </div>
                {order.confirmed_at && (
                  <div className="flex justify-between">
                    <span>Confirmed:</span>
                    <span className="text-foreground font-mono font-medium">
                      {new Date(order.confirmed_at).toLocaleString(undefined, { dateStyle: 'medium', timeStyle: 'short' })}
                    </span>
                  </div>
                )}
                {order.handling_expires_at && (
                  <div className="flex justify-between">
                    <span>Handling Deadline:</span>
                    <span className="text-foreground font-mono font-semibold text-amber-600 dark:text-amber-400">
                      {new Date(order.handling_expires_at).toLocaleString(undefined, { dateStyle: 'medium', timeStyle: 'short' })}
                    </span>
                  </div>
                )}
                {order.updated_at && (
                  <div className="flex justify-between">
                    <span>Last Updated:</span>
                    <span className="text-foreground font-mono font-medium">
                      {new Date(order.updated_at).toLocaleString(undefined, { dateStyle: 'medium', timeStyle: 'short' })}
                    </span>
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* 5. Order Item Details Overlay Panel */}
      <OrderItemDetailOverlay
        item={selectedItemForDetail}
        isOpen={isItemDetailOpen}
        onOpenChange={setIsItemDetailOpen}
        shipment={selectedItemMatchedShipment?.shipment}
        packageIndex={selectedItemMatchedShipment?.index}
      />
    </div>
  );
}
