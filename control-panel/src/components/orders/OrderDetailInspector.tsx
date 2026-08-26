import React, { useState, useEffect, useCallback } from 'react';
import {
  PackageOpen,
  ArrowLeft,
  CheckCircle2,
  XCircle,
  Truck,
  CreditCard,
  MapPin,
  User,
  RefreshCw,
  Edit3,
  Activity,
  Sprout,
  Layers,
  Loader2,
  AlertTriangle,
} from 'lucide-react';
import { Button } from '../ui/button';
import { Input } from '../ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '../ui/select';
import StatusBadge from '../StatusBadge';
import type { Order } from '../../models/Order';
import type { ShipmentDispatchPayload, OrderTrackingResponse } from '../../viewmodels/useOrderActionsViewModel';
import ShipmentDispatchSection from './ShipmentDispatchSection';

export interface OrderDetailInspectorProps {
  order: Order | null;
  submitting: boolean;
  shopId?: string;
  onClose?: () => void;
  onStartProcessing: (orderId: string) => Promise<void>;
  onDispatchOrder?: (orderId: string, shipments: ShipmentDispatchPayload[]) => Promise<void>;
  onDispatchShopShipment?: (
    orderId: string,
    payload: {
      shop_id: string;
      fulfillment_method: string;
      courier: string;
      service: string;
      tracking_number?: string;
      item_ids: string[];
    }
  ) => Promise<void>;
  onUpdateShipmentStatus: (shipmentId: string, status: string) => Promise<void>;
  onUpdateWaybill: (
    shipmentId: string,
    details: { tracking_number?: string; courier?: string; service?: string }
  ) => Promise<void>;
  fetchOrderTracking: (orderId: string) => Promise<OrderTrackingResponse | null>;
}

const formatCurrency = (amount: number) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
  }).format(amount);
};

export const OrderDetailInspector: React.FC<OrderDetailInspectorProps> = ({
  order,
  submitting,
  shopId,
  onClose,
  onStartProcessing,
  onDispatchOrder,
  onDispatchShopShipment,
  onUpdateShipmentStatus,
  onUpdateWaybill,
  fetchOrderTracking,
}) => {
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
    async (orderId: string) => {
      if (syncCooldown > 0) return;
      setLoadingLiveTracking(true);
      try {
        const res = await fetchOrderTracking(orderId);
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

  if (!order) {
    return (
      <div className="space-y-4 max-w-sm flex flex-col items-center py-12 text-center text-muted-foreground mx-auto">
        <PackageOpen className="h-12 w-12 text-muted-foreground/50 stroke-[1.5]" />
        <h3 className="text-lg font-bold font-display tracking-tight text-foreground">No Order Selected</h3>
        <p className="text-xs text-muted-foreground leading-relaxed">
          Select an order from the list on the left to inspect items, allocate split shipments, manage waybills, and
          update transit status.
        </p>
      </div>
    );
  }

  const currentShipments =
    order.shipments && order.shipments.length > 0 ? order.shipments : order.shipment ? [order.shipment] : [];

  const handleStartEditingWaybill = (shipment: any) => {
    setEditingShipmentId(shipment.id);
    setEditCourier(shipment.courier || '');
    setEditService(shipment.service || '');
    setEditTracking(shipment.tracking_number || '');
  };

  const handleSaveWaybill = async (shipmentId: string) => {
    await onUpdateWaybill(shipmentId, {
      tracking_number: editTracking,
      courier: editCourier,
      service: editService,
    });
    setEditingShipmentId(null);
  };

  return (
    <div className="w-full flex-1 space-y-0">
      {/* 1. Header Toolbar */}
      <div className="p-6 border-b border-border/60 bg-muted/10 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div className="flex items-center gap-3">
          {onClose && (
            <Button
              variant="ghost"
              size="icon"
              onClick={onClose}
              className="rounded-xl h-8 w-8 text-muted-foreground hover:text-foreground lg:hidden cursor-pointer"
            >
              <ArrowLeft className="h-4 w-4" />
            </Button>
          )}
          <div>
            <div className="flex items-center gap-2">
              <h4 className="text-lg font-bold font-mono text-foreground">{order.number}</h4>
              <StatusBadge status={order.status} />
            </div>
            <p className="text-xs text-muted-foreground mt-0.5">
              Ordered on {new Date(order.created_at).toLocaleString()}
            </p>
          </div>
        </div>

        {onClose && (
          <Button
            variant="outline"
            size="sm"
            onClick={onClose}
            className="rounded-xl text-xs h-8 hidden lg:flex cursor-pointer"
          >
            Close Inspector
          </Button>
        )}
      </div>

      {/* 2. State Transition Action Bar */}
      <div className="p-6 border-b border-border/60 space-y-4">
        <div className="flex items-center justify-between">
          <h5 className="text-xs font-bold uppercase tracking-wider text-muted-foreground">Fulfillment Actions</h5>
          <span className="text-[11px] font-medium text-muted-foreground">
            Status: <strong className="text-foreground capitalize">{order.status}</strong>
          </span>
        </div>

        {/* CASE A: PENDING PAYMENT */}
        {order.status === 'pending' && (
          <div className="p-4 rounded-xl border border-amber-500/20 bg-amber-500/5 text-amber-900 dark:text-amber-200 text-xs flex items-center justify-between">
            <span>Awaiting payment settlement from customer. Status will auto-advance once paid.</span>
          </div>
        )}

        {/* CASE B: CONFIRMED - READY FOR PACKAGING */}
        {order.status === 'confirmed' && (
          <div className="p-4 rounded-xl border border-blue-500/20 bg-blue-500/5 text-xs text-blue-950 dark:text-blue-200 flex items-center justify-between">
            <span>Payment verified. Ready to initiate order assembly and warehouse processing.</span>
            <Button
              size="sm"
              disabled={submitting}
              onClick={() => onStartProcessing(order.id)}
              className="bg-primary text-primary-foreground hover:bg-primary/90 text-xs rounded-xl h-8 font-semibold cursor-pointer"
            >
              <Sprout className="h-3.5 w-3.5 mr-1" /> Start Processing
            </Button>
          </div>
        )}

        {/* CASE C: PROCESSING - MULTI / SPLIT SHIPMENT DISPATCH */}
        {order.status === 'processing' && (
          <ShipmentDispatchSection
            order={order}
            submitting={submitting}
            shopId={shopId}
            onDispatchOrder={onDispatchOrder}
            onDispatchShopShipment={onDispatchShopShipment}
          />
        )}

        {/* ACTIVE LOGISTICS & SHIPMENT STATUS */}
        {(order.status === 'shipped' || (order.status === 'processing' && currentShipments.length > 0)) &&
          currentShipments.length > 0 && (
            <div className="border border-border/80 bg-background p-5 rounded-xl space-y-5">
              <div className="space-y-1">
                <h5 className="text-sm font-semibold text-foreground flex items-center gap-1.5">
                  <Layers className="h-4 w-4 text-primary" /> Active Logistics & Shipment Status
                </h5>
                <p className="text-xs text-muted-foreground">
                  Manage waybill information and post transit timeline updates per package.
                </p>
              </div>

              {/* Shipments List */}
              <div className="space-y-4">
                {currentShipments.map((shipment, sIdx) => {
                  const isEditing = editingShipmentId === shipment.id;
                  const isSelectedForTimeline = activeShipmentId === shipment.id;

                  const assignedItems = order.items.filter(
                    (item) => item.shipment_id === shipment.id || shipment.item_ids?.includes(item.id)
                  );

                  const isMyShopShipment = !shopId || assignedItems.some((i) => i.shop_id === shopId);

                  return (
                    <div
                      key={shipment.id}
                      className={`border rounded-xl p-4 transition-all ${
                        isSelectedForTimeline ? 'border-primary/60 bg-muted/20' : 'border-border/60 bg-muted/5'
                      }`}
                    >
                      <div className="flex flex-col sm:flex-row justify-between sm:items-center gap-2 pb-3 border-b border-border/40">
                        <div>
                          <div className="flex items-center gap-2">
                            <span className="text-xs font-bold text-foreground">
                              Package #{sIdx + 1}: {shipment.courier.toUpperCase()} {shipment.service.toUpperCase()}
                            </span>
                            <span className="px-2 py-0.5 rounded-full text-[10px] font-bold bg-primary/10 text-primary uppercase">
                              {shipment.status}
                            </span>
                            {!isMyShopShipment && (
                              <span className="px-2 py-0.5 rounded text-[10px] font-semibold bg-muted text-muted-foreground">
                                Other Shop
                              </span>
                            )}
                          </div>
                          <div className="text-[11px] text-muted-foreground mt-0.5">
                            Waybill:{' '}
                            <strong className="font-mono text-foreground">
                              {shipment.tracking_number || 'Self Delivery'}
                            </strong>{' '}
                            · {assignedItems.length} item{assignedItems.length !== 1 ? 's' : ''}
                          </div>
                        </div>

                        {isMyShopShipment && (
                          <div className="flex items-center gap-2">
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => {
                                if (isEditing) {
                                  setEditingShipmentId(null);
                                } else {
                                  handleStartEditingWaybill(shipment);
                                }
                              }}
                              className="text-xs text-primary hover:bg-primary/5 rounded-lg flex items-center gap-1 h-7 cursor-pointer"
                            >
                              <Edit3 className="h-3 w-3" /> {isEditing ? 'Cancel' : 'Edit Waybill'}
                            </Button>
                          </div>
                        )}
                      </div>

                      {isMyShopShipment ? (
                        isEditing ? (
                          /* EDIT WAYBILL SUB-FORM */
                          <div className="space-y-3 p-3 mt-3 border border-border/40 rounded-xl bg-background">
                            <h6 className="text-xs font-bold text-foreground">Edit Shipping Waybill</h6>
                            <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
                              <div className="space-y-1">
                                <label className="text-[10px] font-semibold text-muted-foreground uppercase">
                                  Courier
                                </label>
                                <Input
                                  value={editCourier}
                                  onChange={(e) => setEditCourier(e.target.value)}
                                  className="h-8 text-xs rounded-lg"
                                />
                              </div>
                              <div className="space-y-1">
                                <label className="text-[10px] font-semibold text-muted-foreground uppercase">
                                  Service
                                </label>
                                <Input
                                  value={editService}
                                  onChange={(e) => setEditService(e.target.value)}
                                  className="h-8 text-xs rounded-lg"
                                />
                              </div>
                              <div className="space-y-1">
                                <label className="text-[10px] font-semibold text-muted-foreground uppercase">
                                  Tracking No
                                </label>
                                <Input
                                  value={editTracking}
                                  onChange={(e) => setEditTracking(e.target.value)}
                                  className="h-8 text-xs rounded-lg"
                                />
                              </div>
                            </div>
                            <Button
                              size="sm"
                              disabled={submitting}
                              onClick={() => handleSaveWaybill(shipment.id)}
                              className="h-8 text-xs rounded-lg w-full bg-primary text-primary-foreground hover:bg-primary/90 mt-2 font-semibold cursor-pointer"
                            >
                              Save Waybill Changes
                            </Button>
                          </div>
                        ) : (
                          /* TRANSIT STATE SELECTOR */
                          <div className="space-y-3 mt-3">
                            <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
                              <div className="space-y-1.5">
                                <label className="text-[11px] font-semibold text-muted-foreground">
                                  New Transit Status
                                </label>
                                <Select
                                  value={transitStatus}
                                  onValueChange={(val) => {
                                    setActiveShipmentId(shipment.id);
                                    setTransitStatus(val);
                                  }}
                                >
                                  <SelectTrigger className="rounded-xl h-8 text-xs bg-background">
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
                                <label className="text-[11px] font-semibold text-muted-foreground">
                                  Location (Optional)
                                </label>
                                <Input
                                  value={transitLocation}
                                  onChange={(e) => setTransitLocation(e.target.value)}
                                  placeholder="Warehouse, Hub City, etc."
                                  className="rounded-xl h-8 text-xs bg-background"
                                />
                              </div>
                            </div>
                            <div className="space-y-1.5">
                              <label className="text-[11px] font-semibold text-muted-foreground">
                                Description/Note (Optional)
                              </label>
                              <Input
                                value={transitDescription}
                                onChange={(e) => setTransitDescription(e.target.value)}
                                placeholder="E.g. Package hand over to courier JNE"
                                className="rounded-xl h-8 text-xs bg-background"
                              />
                            </div>
                            <Button
                              disabled={submitting}
                              onClick={() => onUpdateShipmentStatus(shipment.id, transitStatus)}
                              className="rounded-xl bg-primary text-primary-foreground hover:bg-primary/90 text-xs font-semibold h-8 w-full cursor-pointer"
                            >
                              Update Package #{sIdx + 1} Status
                            </Button>
                          </div>
                        )
                      ) : (
                        <div className="mt-2 text-[11px] text-muted-foreground italic">
                          This package is fulfilled and managed by another store branch.
                        </div>
                      )}
                    </div>
                  );
                })}
              </div>
            </div>
          )}

        {/* CASE E: DELIVERED */}
        {order.status === 'delivered' && (
          <div className="p-4 rounded-xl border border-emerald-500/20 bg-emerald-500/5 text-emerald-950 dark:text-emerald-200 text-xs flex items-center gap-2">
            <CheckCircle2 className="h-4 w-4 text-emerald-600" />
            <span>Order fulfillment complete. All packages delivered to customer.</span>
          </div>
        )}

        {/* CASE F: CANCELLED */}
        {order.status === 'cancelled' && (
          <div className="p-4 rounded-xl border border-destructive/20 bg-destructive/5 text-destructive text-xs flex items-center gap-2">
            <XCircle className="h-4 w-4 text-destructive" />
            <span>This order was cancelled. Reserved inventory was released back to stock.</span>
          </div>
        )}
      </div>

      {/* 3. Detailed Panels Grid */}
      <div className="p-6 space-y-8">
        {/* Order Items Section */}
        <div className="space-y-3">
          <h5 className="text-sm font-bold text-foreground font-display flex items-center gap-2">
            <PackageOpen className="h-4 w-4 text-muted-foreground" /> Order Items
          </h5>
          <div className="border border-border/60 rounded-xl overflow-hidden bg-background">
            <table className="w-full text-left border-collapse text-xs">
              <thead>
                <tr className="bg-muted/40 text-muted-foreground border-b border-border/40 font-semibold">
                  <th className="p-3">Product Name</th>
                  <th className="p-3 text-center">Package Assigned</th>
                  <th className="p-3 text-center">Quantity</th>
                  <th className="p-3 text-right">Price</th>
                  <th className="p-3 text-right">Subtotal</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border/40">
                {order.items.map((item) => {
                  const matchedShipmentIndex = currentShipments.findIndex(
                    (s) => s.id === item.shipment_id || s.item_ids?.includes(item.id)
                  );
                  const isMyItem = shopId ? item.shop_id === shopId : false;

                  return (
                    <tr
                      key={item.id}
                      className={`text-foreground hover:bg-muted/10 transition-colors ${
                        isMyItem ? 'bg-primary/5 font-medium' : ''
                      }`}
                    >
                      <td className="p-3">
                        <div className="font-semibold text-foreground flex items-center gap-1.5 flex-wrap">
                          <span>{item.product_name}</span>
                          {item.item_options?.size && (
                            <span className="px-1.5 py-0.5 rounded text-[9px] font-bold bg-muted text-muted-foreground border border-border/50">
                              Size: {item.item_options.size === 'small' ? '1.5 × 2.0m' : item.item_options.size === 'medium' ? '1.8 × 2.5m' : item.item_options.size === 'large' ? '2.0 × 3.0m' : item.item_options.size}
                            </span>
                          )}
                          {item.item_options?.jambul && (
                            <span className="px-1.5 py-0.5 rounded text-[9px] font-bold bg-emerald-500/10 text-emerald-700 dark:text-emerald-400 border border-emerald-500/20">
                              Jambul: {item.item_options.jambul === 'none' ? 'None' : item.item_options.jambul === 'top' ? 'Top' : item.item_options.jambul === 'bottom' ? 'Bottom' : item.item_options.jambul === 'both' ? 'Both' : item.item_options.jambul}
                            </span>
                          )}
                          {isMyItem && (
                            <span className="px-1.5 py-0.5 rounded text-[9px] font-bold bg-primary/15 text-primary">
                              Your Shop
                            </span>
                          )}
                        </div>
                        <div className="text-[10px] text-muted-foreground">Shop: {item.shop_name}</div>
                        {item.courier_code && (
                          <div className="text-[10px] text-muted-foreground uppercase">
                            Expected: {item.courier_code} {item.courier_service}
                          </div>
                        )}
                      </td>
                      <td className="p-3 text-center">
                        {matchedShipmentIndex >= 0 ? (
                          <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-bold bg-primary/10 text-primary">
                            <Truck className="h-3 w-3" /> Pkg #{matchedShipmentIndex + 1}
                          </span>
                        ) : (
                          <span className="text-muted-foreground/60 text-[10px]">Unassigned</span>
                        )}
                      </td>
                      <td className="p-3 text-center font-medium">{item.quantity}</td>
                      <td className="p-3 text-right font-mono">{formatCurrency(item.unit_price)}</td>
                      <td className="p-3 text-right font-mono font-semibold">{formatCurrency(item.subtotal)}</td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        </div>

        {/* Payment & Customer Details Split */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {/* Payment Summary */}
          <div className="space-y-3">
            <h5 className="text-xs font-bold uppercase tracking-wider text-muted-foreground flex items-center gap-1.5">
              <CreditCard className="h-3.5 w-3.5" /> Billing details
            </h5>
            <div className="p-4 rounded-xl border border-border/60 bg-muted/10 space-y-2.5 text-xs">
              <div className="flex justify-between text-muted-foreground">
                <span>Items Subtotal</span>
                <span className="font-medium text-foreground">{formatCurrency(order.subtotal)}</span>
              </div>
              <div className="flex justify-between text-muted-foreground">
                <span>Shipping Fee</span>
                <span className="font-medium text-foreground">{formatCurrency(order.shipping_fee)}</span>
              </div>
              <div className="border-t border-border/60 pt-2.5 mt-1 flex justify-between text-sm font-bold text-foreground">
                <span>Grand Total</span>
                <span className="text-primary">{formatCurrency(order.total)}</span>
              </div>

              {order.payment && (
                <div className="border-t border-border/40 pt-2.5 mt-2 space-y-1.5 font-medium text-[10px] text-muted-foreground">
                  <div>
                    Payment Method:{' '}
                    <span className="text-foreground font-semibold uppercase">{order.payment.provider}</span>
                  </div>
                  <div>
                    Transaction Status:{' '}
                    <span className="text-foreground font-semibold uppercase">{order.payment.status}</span>
                  </div>
                </div>
              )}
            </div>
          </div>

          {/* Customer & Shipping Address */}
          <div className="space-y-3">
            <h5 className="text-xs font-bold uppercase tracking-wider text-muted-foreground flex items-center gap-1.5">
              <User className="h-3.5 w-3.5" /> Shipping Address & Packages
            </h5>
            <div className="p-4 rounded-xl border border-border/60 bg-muted/10 space-y-3 text-xs">
              <div className="space-y-1">
                <div className="font-semibold text-foreground flex items-center gap-1">
                  <MapPin className="h-3.5 w-3.5 text-muted-foreground" />{' '}
                  {order.address ? 'Delivery Address' : 'Address Metadata'}
                </div>
                {order.address ? (
                  <div className="text-foreground space-y-1 leading-relaxed text-[11px] font-medium">
                    <p>
                      <span className="font-semibold text-muted-foreground">Receiver:</span>{' '}
                      {order.address.receiver_name}
                    </p>
                    <p>
                      <span className="font-semibold text-muted-foreground">Phone:</span> {order.address.phone}
                    </p>
                    <p>
                      <span className="font-semibold text-muted-foreground">Address:</span>{' '}
                      {order.address.full_address}
                    </p>
                    {order.address.postal_code && (
                      <p>
                        <span className="font-semibold text-muted-foreground">Postal Code:</span>{' '}
                        {order.address.postal_code}
                      </p>
                    )}
                  </div>
                ) : (
                  <div className="text-muted-foreground space-y-1 leading-relaxed text-[11px]">
                    <p>
                      <span className="font-semibold text-foreground">Customer UUID:</span>{' '}
                      <span className="font-mono">{order.customer_id || order.user_id}</span>
                    </p>
                    <p>
                      <span className="font-semibold text-foreground">Address ID:</span>{' '}
                      <span className="font-mono">{order.address_id}</span>
                    </p>
                  </div>
                )}
              </div>

              {currentShipments.length > 0 && (
                <div className="border-t border-border/40 pt-2.5 mt-2 space-y-2 text-[11px] text-muted-foreground">
                  <div className="font-bold text-foreground">Packages Dispatched:</div>
                  {currentShipments.map((s, idx) => (
                    <div key={s.id} className="p-2 rounded bg-background border border-border/40 space-y-0.5">
                      <div className="font-semibold text-foreground flex justify-between">
                        <span>
                          Pkg #{idx + 1}: {s.courier.toUpperCase()} {s.service.toUpperCase()}
                        </span>
                        <span className="uppercase text-[10px] text-primary font-bold">{s.status}</span>
                      </div>
                      <div className="font-mono text-primary text-[10px]">
                        {s.tracking_number || 'Self Delivery'}
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>
        </div>

        {/* Logistics Shipment History Timeline */}
        {currentShipments.length > 0 && (
          <div className="space-y-4 pt-2 border-t border-border/60">
            <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2">
              <div className="flex items-center gap-2">
                <h5 className="text-sm font-bold text-foreground font-display flex items-center gap-1.5">
                  <Activity className="h-4 w-4 text-muted-foreground" /> Shipment Tracking History
                </h5>
                {liveTracking?.timeline && liveTracking.timeline.length > 0 && (
                  <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-bold bg-emerald-500/10 text-emerald-600 border border-emerald-500/20">
                    Live Courier Feed
                  </span>
                )}
              </div>
              <div className="flex items-center gap-2 self-start sm:self-auto">
                {loadingLiveTracking && <Loader2 className="h-3.5 w-3.5 animate-spin text-primary" />}
                <Button
                  variant="outline"
                  size="sm"
                  disabled={loadingLiveTracking || syncCooldown > 0}
                  onClick={() => order && loadLiveTracking(order.id)}
                  className="h-7 px-2.5 text-[11px] rounded-lg gap-1 border-primary/30 text-primary hover:bg-primary/5 cursor-pointer disabled:opacity-60"
                >
                  <RefreshCw className={`h-3 w-3 ${loadingLiveTracking ? 'animate-spin' : ''}`} />
                  {syncCooldown > 0 ? `Synced (${syncCooldown}s)` : 'Sync Courier Status'}
                </Button>
                {currentShipments.length > 1 && (
                  <div className="flex gap-1.5">
                    {currentShipments.map((s, idx) => (
                      <button
                        key={s.id}
                        onClick={() => setActiveShipmentId(s.id)}
                        className={`px-2.5 py-1 rounded-lg text-xs font-semibold transition-all cursor-pointer ${
                          activeShipmentId === s.id
                            ? 'bg-primary text-primary-foreground'
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
              <div className="flex items-start gap-2.5 p-3.5 rounded-xl bg-amber-500/10 border border-amber-500/20 text-amber-700 dark:text-amber-300 text-xs font-medium leading-relaxed">
                <AlertTriangle className="h-4 w-4 shrink-0 text-amber-500 mt-0.5" />
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
                  <div className="relative pl-6 border-l-2 border-primary/20 space-y-6 ml-2 py-1">
                    {eventsToDisplay.map((event, idx) => {
                      const isLatest = idx === eventsToDisplay.length - 1 || idx === 0;
                      return (
                        <div key={event.id} className="relative">
                          <span
                            className={`absolute -left-[31px] top-1 h-4 w-4 rounded-full border-2 bg-background flex items-center justify-center transition-colors ${
                              isLatest ? 'border-primary ring-4 ring-primary/10' : 'border-muted-foreground/45'
                            }`}
                          >
                            <span
                              className={`h-1.5 w-1.5 rounded-full ${
                                isLatest ? 'bg-primary' : 'bg-muted-foreground/45'
                              }`}
                            />
                          </span>

                          <div className="space-y-1">
                            <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-1">
                              <div className="flex items-center gap-2">
                                <span
                                  className={`text-xs font-bold uppercase tracking-wider ${
                                    isLatest ? 'text-primary' : 'text-muted-foreground'
                                  }`}
                                >
                                  {event.status}
                                </span>
                                {event.location && (
                                  <span className="text-[10px] bg-muted px-1.5 py-0.5 rounded text-muted-foreground font-medium">
                                    {event.location}
                                  </span>
                                )}
                              </div>
                              <span className="text-[10px] text-muted-foreground font-mono">
                                {new Date(event.timestamp).toLocaleString(undefined, {
                                  dateStyle: 'short',
                                  timeStyle: 'short',
                                })}
                              </span>
                            </div>
                            <p className="text-xs text-muted-foreground leading-relaxed">{event.description}</p>
                          </div>
                        </div>
                      );
                    })}
                  </div>
                );
              }

              return (
                <div className="p-4 border border-dashed border-border/80 rounded-xl bg-muted/10 text-center text-xs text-muted-foreground flex flex-col items-center gap-1">
                  <Activity className="h-6 w-6 text-muted-foreground/50 stroke-[1.5]" />
                  <span>
                    No shipment updates logged yet for{' '}
                    {activeShipment
                      ? `${activeShipment.courier.toUpperCase()} (${
                          activeShipment.tracking_number || 'Self Delivery'
                        })`
                      : 'this package'}
                    .
                  </span>
                  <span className="text-[10px] text-muted-foreground/60">
                    Use the logistics panel above to post the first shipment status update.
                  </span>
                </div>
              );
            })()}
          </div>
        )}
      </div>
    </div>
  );
};

export default OrderDetailInspector;
