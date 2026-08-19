import { useState, useEffect, useCallback } from 'react';
import { Link } from 'react-router-dom';
import {
  PackageOpen,
  ArrowUpDown,
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
  BarChart3,
  Plus,
  Trash2,
  Layers,
  Boxes,
  Clock,
  Store,
  Loader2,
  AlertTriangle,
} from 'lucide-react';
import { Button } from '../../components/ui/button';
import { Card, CardContent } from '../../components/ui/card';
import { useOrdersViewModel } from '../../viewmodels/useOrdersViewModel';
import { useOrderActionsViewModel, type ShipmentDispatchPayload, type OrderTrackingResponse } from '../../viewmodels/useOrderActionsViewModel';
import { useAuthMeViewModel } from '../../viewmodels/useAuthMeViewModel';
import { fetchApi } from '../../lib/api';
import EmptyState from '../../components/EmptyState';
import SearchInput from '../../components/SearchInput';
import StatusBadge from '../../components/StatusBadge';
import Pagination from '../../components/Pagination';
import { Skeleton } from '../../components/ui/skeleton';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '../../components/ui/select';
import { Input } from '../../components/ui/input';
import type { Order } from '../../models/Order';

interface ShipmentGroupForm {
  id: string;
  shop_id?: string;
  shop_name?: string;
  fulfillment_method: string;
  courier: string;
  service: string;
  tracking_number: string;
  item_ids: string[];
}

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
    setPage,
    setSort,
    setSearchNumber,
    setStatusFilter,
    setShopFilter,
    refresh
  } = useOrdersViewModel();

  const {
    submitting,
    fetchOrderTracking,
    updateOrderStatus,
    updateShipmentStatus,
    updateShipmentDetails
  } = useOrderActionsViewModel();

  const [selectedOrderId, setSelectedOrderId] = useState<string | null>(null);
  const [shopsList, setShopsList] = useState<Array<{ id: string; name: string; slug?: string }>>([]);

  useEffect(() => {
    fetchApi('/shops?limit=100')
      .then(res => {
        if (res?.shops && res.shops.length > 0) {
          setShopsList(res.shops);
          if (!isAdmin && !shopFilter) {
            const firstShop = res.shops[0];
            setShopFilter(firstShop.slug || firstShop.id);
          }
        }
      })
      .catch(err => {
        console.error('Failed to load shops for order filter', err);
      });
  }, [isAdmin, shopFilter, setShopFilter]);

  // Multi-shipment dispatch state (for processing -> shipped)
  const [shipmentGroups, setShipmentGroups] = useState<ShipmentGroupForm[]>([]);

  // Per-shipment waybill editing state (for shipped status)
  const [editingShipmentId, setEditingShipmentId] = useState<string | null>(null);
  const [editCourier, setEditCourier] = useState<string>('');
  const [editService, setEditService] = useState<string>('');
  const [editTracking, setEditTracking] = useState<string>('');

  // Per-shipment transit updating state (for shipped status)
  const [activeShipmentId, setActiveShipmentId] = useState<string | null>(null);
  const [transitStatus, setTransitStatus] = useState<string>('packed');
  const [transitDescription, setTransitDescription] = useState<string>('');
  const [transitLocation, setTransitLocation] = useState<string>('');

  // Live external courier tracking state & cooldown
  const [liveTracking, setLiveTracking] = useState<OrderTrackingResponse | null>(null);
  const [loadingLiveTracking, setLoadingLiveTracking] = useState<boolean>(false);
  const [syncCooldown, setSyncCooldown] = useState<number>(0);

  const selectedOrder = data?.orders.find(o => o.id === selectedOrderId);

  // Cooldown countdown timer for sync button to prevent rate limit hits
  useEffect(() => {
    if (syncCooldown <= 0) return;
    const timer = setInterval(() => {
      setSyncCooldown(prev => (prev > 0 ? prev - 1 : 0));
    }, 1000);
    return () => clearInterval(timer);
  }, [syncCooldown]);

  // Helper to refresh live tracking from backend
  const loadLiveTracking = useCallback(async (orderId: string) => {
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
  }, [fetchOrderTracking, syncCooldown]);

  // Sync dispatch forms based on shops when selected order changes
  useEffect(() => {
    if (selectedOrder) {
      if (selectedOrder.status === 'processing') {
        // Group items by shop_id to initialize default shipments per shop
        const shopMap = new Map<string, { shop_name: string; items: typeof selectedOrder.items }>();
        selectedOrder.items.forEach(item => {
          const shopId = item.shop_id || 'default-shop';
          const existing = shopMap.get(shopId) || { shop_name: item.shop_name || 'Shop', items: [] };
          existing.items.push(item);
          shopMap.set(shopId, existing);
        });

        const initialGroups: ShipmentGroupForm[] = [];
        let groupIndex = 1;
        shopMap.forEach((shopData, shopId) => {
          const firstItem = shopData.items[0];
          initialGroups.push({
            id: `shipment-shop-${shopId}-${Date.now()}-${groupIndex}`,
            shop_id: shopId,
            shop_name: shopData.shop_name,
            fulfillment_method: 'courier',
            courier: firstItem?.courier_code || 'JNE',
            service: firstItem?.courier_service || 'REG',
            tracking_number: '',
            item_ids: shopData.items.map(i => i.id),
          });
          groupIndex++;
        });

        setShipmentGroups(initialGroups);
      } else {
        setShipmentGroups([]);
      }

      const activeShipments = selectedOrder.shipments && selectedOrder.shipments.length > 0
        ? selectedOrder.shipments
        : (selectedOrder.shipment ? [selectedOrder.shipment] : []);

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
  }, [selectedOrderId, selectedOrder?.id, selectedOrder?.status]);

  const handleSort = () => {
    const currentDirection = sort.split(':')[1];
    const newDirection = currentDirection === 'desc' ? 'asc' : 'desc';
    setSort(`latest:${newDirection}`);
    setPage(1);
  };

  const handleStartProcessing = async (orderId: string) => {
    try {
      await updateOrderStatus(orderId, 'processing');
      refresh();
    } catch (err) {
      // Handled in ViewModel
    }
  };

  const handleAddShipmentGroup = (targetShopId?: string) => {
    if (!selectedOrder) return;
    const shopItems = targetShopId
      ? selectedOrder.items.filter(i => i.shop_id === targetShopId)
      : selectedOrder.items;
    const firstItem = shopItems[0] || selectedOrder.items[0];

    setShipmentGroups(prev => [
      ...prev,
      {
        id: `shipment-${Date.now()}-${prev.length + 1}`,
        shop_id: targetShopId || firstItem?.shop_id,
        shop_name: firstItem?.shop_name || 'Shop',
        fulfillment_method: 'courier',
        courier: firstItem?.courier_code || 'JNE',
        service: firstItem?.courier_service || 'REG',
        tracking_number: '',
        item_ids: [],
      },
    ]);
  };

  const handleRemoveShipmentGroup = (groupId: string) => {
    if (shipmentGroups.length <= 1) return;
    const removedGroup = shipmentGroups.find(g => g.id === groupId);
    const remainingGroups = shipmentGroups.filter(g => g.id !== groupId);

    // Reassign unallocated items from removed group back to the matching shop shipment or first group
    if (removedGroup && removedGroup.item_ids.length > 0 && remainingGroups.length > 0) {
      const matchShopGroup = remainingGroups.find(g => g.shop_id === removedGroup.shop_id) || remainingGroups[0];
      matchShopGroup.item_ids = [
        ...matchShopGroup.item_ids,
        ...removedGroup.item_ids,
      ];
    }
    setShipmentGroups([...remainingGroups]);
  };

  const handleToggleItemInGroup = (groupId: string, itemId: string) => {
    setShipmentGroups(prev =>
      prev.map(g => {
        if (g.id === groupId) {
          const exists = g.item_ids.includes(itemId);
          return {
            ...g,
            item_ids: exists
              ? g.item_ids.filter(id => id !== itemId)
              : [...g.item_ids, itemId],
          };
        } else {
          return {
            ...g,
            item_ids: g.item_ids.filter(id => id !== itemId),
          };
        }
      })
    );
  };

  const handleUpdateGroupField = (
    groupId: string,
    field: keyof ShipmentGroupForm,
    value: any
  ) => {
    setShipmentGroups(prev =>
      prev.map(g => (g.id === groupId ? { ...g, [field]: value } : g))
    );
  };

  const handleDispatchOrder = async (orderId: string) => {
    if (!selectedOrder) return;

    if (shipmentGroups.length === 0) {
      alert('At least one shipment is required to dispatch the order.');
      return;
    }

    const allAssignedItemIds = shipmentGroups.flatMap(g => g.item_ids);
    const unassignedItems = selectedOrder.items.filter(
      item => !allAssignedItemIds.includes(item.id)
    );
    if (unassignedItems.length > 0) {
      alert(
        `All order items must be assigned to a shipment. ${unassignedItems.length} item(s) unassigned.`
      );
      return;
    }

    for (let i = 0; i < shipmentGroups.length; i++) {
      const group = shipmentGroups[i];
      if (group.item_ids.length === 0) {
        alert(`Shipment #${i + 1} (${group.shop_name}) has no items assigned. Please assign items or remove it.`);
        return;
      }
      if (group.fulfillment_method === 'courier' && !group.tracking_number.trim()) {
        alert(`Shipment #${i + 1} (${group.shop_name}) is missing a waybill tracking number.`);
        return;
      }
    }

    const payload: ShipmentDispatchPayload[] = shipmentGroups.map(g => ({
      fulfillment_method: g.fulfillment_method,
      courier: g.courier,
      service: g.service,
      tracking_number: g.tracking_number.trim() || undefined,
      item_ids: g.item_ids,
    }));

    try {
      await updateOrderStatus(orderId, 'shipped', undefined, undefined, payload);
      refresh();
    } catch (err) {
      // Handled in ViewModel
    }
  };

  const handleStartEditingWaybill = (shipment: any) => {
    setEditingShipmentId(shipment.id);
    setEditCourier(shipment.courier || '');
    setEditService(shipment.service || '');
    setEditTracking(shipment.tracking_number || '');
  };

  const handleUpdateWaybill = async (shipmentId: string) => {
    try {
      await updateShipmentDetails(shipmentId, {
        tracking_number: editTracking,
        courier: editCourier,
        service: editService,
      });
      setEditingShipmentId(null);
      refresh();
    } catch (err) {
      // Handled in ViewModel
    }
  };

  const handleUpdateShipmentStatus = async (shipmentId: string) => {
    try {
      await updateShipmentStatus(shipmentId, transitStatus);
      refresh();
    } catch (err) {
      // Handled in ViewModel
    }
  };

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      minimumFractionDigits: 0,
    }).format(amount);
  };

  /**
   * Formats order expiration time for staff:
   * - Uses words when remaining time is between 1 second and 23.59 hours (e.g. "4h 15m left").
   * - Uses formatted date when time is > 23.59 hours.
   * - Shows "Expired" when past due.
   */
  const getOrderExpiry = (order: Order): { label: string; isUrgent?: boolean; isExpired?: boolean } | null => {
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

    // eslint-disable-next-line react-hooks/purity
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

    // 24 hours or greater: display formatted date/time
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

  const statusTabs = [
    { value: 'all', label: 'All Orders' },
    { value: 'pending', label: 'Awaiting Payment' },
    { value: 'confirmed', label: 'To Process' },
    { value: 'processing', label: 'In Packaging' },
    { value: 'shipped', label: 'Shipped' },
    { value: 'delivered', label: 'Delivered' },
    { value: 'cancelled', label: 'Cancelled' },
  ];

  const currentShipments = selectedOrder
    ? selectedOrder.shipments && selectedOrder.shipments.length > 0
      ? selectedOrder.shipments
      : selectedOrder.shipment
      ? [selectedOrder.shipment]
      : []
    : [];

  return (
    <div className="flex-col md:flex">
      <div className="flex-1 space-y-12 p-6 sm:p-8 lg:p-12 animate-in fade-in duration-300">

        {/* Header */}
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <h2 className="text-3xl font-bold font-display tracking-tight text-foreground">Orders</h2>
            <p className="text-muted-foreground text-sm">
              Verify payments, fulfill shipments, and manage orders
            </p>
          </div>
          {isAdmin && (
            <Button asChild variant="outline" size="sm" className="rounded-xl text-xs gap-1.5 border-primary/30 text-primary hover:bg-primary/5 self-start sm:self-auto">
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
            <p className="text-muted-foreground text-sm">Manage fulfillment, check payment updates, and inspect order status logs.</p>
          </div>

          {/* Assigned Shops Segmented Switcher */}
          {shopsList.length > 0 && (
            <div className="flex flex-col gap-2.5 p-4 rounded-2xl border border-border/60 bg-muted/20">
              <div className="flex items-center justify-between">
                <span className="text-xs font-semibold uppercase tracking-wider text-muted-foreground flex items-center gap-1.5">
                  <Store className="w-3.5 h-3.5 text-primary" />
                  {isAdmin ? 'Shops Overview' : 'Your Assigned Shops'}
                </span>
                <span className="text-xs text-muted-foreground font-medium">
                  {shopsList.length} {shopsList.length === 1 ? 'Shop Assigned' : 'Shops Assigned'}
                </span>
              </div>
              <div className="flex items-center gap-2 overflow-x-auto pt-0.5 pb-0.5 scrollbar-none">
                {isAdmin && (
                  <button
                    type="button"
                    onClick={() => setShopFilter('all')}
                    className={`px-3.5 py-1.5 rounded-xl text-xs font-medium transition-all shrink-0 border ${
                      (shopFilter || 'all') === 'all'
                        ? 'bg-primary text-primary-foreground border-primary shadow-sm'
                        : 'bg-background hover:bg-muted text-muted-foreground border-border/80'
                    }`}
                  >
                    All Shops
                  </button>
                )}
                {shopsList.map((shop) => {
                  const shopKey = shop.slug || shop.id;
                  const isSelected = shopFilter === shopKey || (!shopFilter && shopsList.length === 1 && shopKey === (shopsList[0].slug || shopsList[0].id));
                  return (
                    <button
                      key={shop.id}
                      type="button"
                      onClick={() => setShopFilter(shopKey)}
                      className={`px-3.5 py-1.5 rounded-xl text-xs font-medium transition-all shrink-0 border flex items-center gap-2 ${
                        isSelected
                          ? 'bg-primary text-primary-foreground border-primary shadow-sm'
                          : 'bg-background hover:bg-muted text-muted-foreground border-border/80'
                      }`}
                    >
                      <Store className="w-3.5 h-3.5 text-current shrink-0" />
                      <span>{shop.name}</span>
                      {shop.slug && (
                        <span className={`text-[10px] px-1.5 py-0.5 rounded-md font-mono ${
                          isSelected ? 'bg-primary-foreground/20 text-primary-foreground' : 'bg-muted text-muted-foreground'
                        }`}>
                          {shop.slug}
                        </span>
                      )}
                    </button>
                  );
                })}
              </div>
            </div>
          )}

          {/* Full-width Search and Status Filter Tabs */}
          <div className={`flex flex-col gap-4 items-start w-full ${
            selectedOrderId ? 'hidden lg:flex' : 'flex'
          }`}>
            <div className="flex items-center justify-between w-full gap-4">
              <div className="flex-1 md:max-w-xs lg:max-w-sm">
                <SearchInput
                  value={searchNumber}
                  onChange={(val) => {
                    setSearchNumber(val);
                    setPage(1);
                  }}
                  placeholder="Search by Order Number..."
                />
              </div>

              {/* Right Side: Refresh */}
              <div className="flex items-center gap-2 justify-end w-full sm:w-auto">

                <Button
                  variant="outline"
                  onClick={() => refresh()}
                  disabled={loading}
                  size="sm"
                  className="rounded-xl h-9 text-xs"
                >
                  <RefreshCw className={`h-3.5 w-3.5 mr-1.5 ${loading ? 'animate-spin' : ''}`} />
                  Refresh
                </Button>
              </div>
            </div>

            {/* Status Filter Tabs */}
            <div className="flex gap-2 overflow-x-auto pb-1 w-full border-b border-border/40 scrollbar-none">
              {statusTabs.map((tab) => {
                const isActive = (statusFilter || 'all') === tab.value;
                return (
                  <button
                    key={tab.value}
                    disabled={isSwitchingCategory || loading}
                    onClick={() => {
                      setStatusFilter(tab.value === 'all' ? '' : tab.value);
                      setSelectedOrderId(null);
                    }}
                    className={`px-3.5 py-1.5 rounded-lg text-xs font-semibold whitespace-nowrap transition-all ${
                      isActive
                        ? 'bg-primary text-primary-foreground shadow-sm'
                        : 'text-muted-foreground hover:text-foreground hover:bg-muted/40'
                    } ${
                      isSwitchingCategory || loading
                        ? 'opacity-60 cursor-not-allowed pointer-events-none'
                        : 'cursor-pointer'
                    }`}
                  >
                    {tab.label}
                  </button>
                );
              })}
            </div>
          </div>

          {/* Master-Detail Split Workspace Layout */}
          <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">

            {/* LEFT PANE: Master Order Cards List (Collapsible / Shrinkable) */}
            <div className={`lg:col-span-5 flex flex-col space-y-4 ${selectedOrderId ? 'hidden lg:flex' : 'flex'}`}>

              {/* Sorting & Order Count Bar */}
              <div className="flex items-center justify-between px-1 text-xs text-muted-foreground">
                <div className="flex items-center gap-2">
                  <span>{data?.total ? `Found ${data.total} orders` : 'No orders found'}</span>
                  {isSwitchingCategory && (
                    <Loader2 className="h-3 w-3 animate-spin text-primary" />
                  )}
                </div>
                <button
                  disabled={isSwitchingCategory || loading}
                  onClick={handleSort}
                  className={`flex items-center gap-1 hover:text-foreground transition-colors font-medium ${
                    isSwitchingCategory || loading ? 'opacity-50 cursor-not-allowed pointer-events-none' : 'cursor-pointer'
                  }`}
                >
                  Date <ArrowUpDown className="h-3 w-3" />
                </button>
              </div>

              {/* List Container with Cards */}
              <div className="space-y-3 pr-1 min-h-[250px]">
                {isSwitchingCategory ? (
                  <div className="py-24 flex flex-col items-center justify-center gap-3 text-muted-foreground bg-muted/10 rounded-2xl border border-dashed border-border/80 min-h-[320px] select-none pointer-events-none animate-in fade-in duration-150">
                    <Loader2 className="h-8 w-8 animate-spin text-primary" />
                    <div className="text-center space-y-0.5">
                      <p className="text-xs font-bold text-foreground">Loading category orders...</p>
                      <p className="text-[11px] text-muted-foreground">Please wait a moment</p>
                    </div>
                  </div>
                ) : loading && !data ? (
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
                ) : !data?.orders || data.orders.length === 0 ? (
                  <EmptyState
                    icon={<PackageOpen className="h-8 w-8 mb-2 mx-auto text-slate-400" />}
                    title="No orders found"
                    description="No orders found matching current criteria."
                    className="py-12 border border-dashed border-border/80 rounded-2xl bg-zinc-50/10"
                  />
                ) : (
                  data.orders.map((order) => {
                    const isSelected = order.id === selectedOrderId;
                    const itemPreview = order.items.map(item => `${item.product_name} (x${item.quantity})`).join(', ');
                    const countShipments = order.shipments && order.shipments.length > 0
                      ? order.shipments.length
                      : order.shipment
                      ? 1
                      : 0;

                    const expiryInfo = getOrderExpiry(order);

                    return (
                      <Card
                        key={order.id}
                        onClick={() => setSelectedOrderId(order.id)}
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
                          <p className="text-xs text-muted-foreground line-clamp-1">
                            {itemPreview}
                          </p>

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
                          </div>

                          {/* Card Footer: Truncated ID & Total Price */}
                          <div className="flex justify-between items-center border-t border-border/40 pt-2 mt-1">
                            <span className="text-[10px] text-muted-foreground font-mono">
                              ID: {order.id.slice(0, 8)}...
                            </span>
                            <span className="text-sm font-bold text-primary font-mono">
                              {formatCurrency(order.total)}
                            </span>
                          </div>
                        </CardContent>
                      </Card>
                    );
                  })
                )}
              </div>

              {/* Pagination */}
              {data && data.total > limit && (
                <div className="pt-2">
                  <Pagination
                    currentPage={page}
                    totalPages={Math.ceil(data.total / limit)}
                    totalItems={data.total}
                    limit={limit}
                    onPageChange={setPage}
                    itemNamePlural="orders"
                  />
                </div>
              )}
            </div>

            {/* RIGHT PANE: Order Action & Detail Inspector (Sticky & Fully Scrollable) */}
            <div className={`lg:col-span-7 border border-border/80 rounded-2xl bg-card flex flex-col lg:sticky lg:top-24 lg:self-start lg:max-h-[calc(100vh-6rem)] overflow-y-auto overscroll-contain pr-0.5 ${
              !selectedOrderId ? 'hidden lg:flex items-center justify-center p-12 text-center text-muted-foreground min-h-[400px]' : 'flex'
            }`}>

              {!selectedOrderId ? (
                /* No Selection Placeholder */
                <div className="space-y-4 max-w-sm flex flex-col items-center py-12">
                  <PackageOpen className="h-12 w-12 text-muted-foreground/50 stroke-[1.5]" />
                  <h3 className="text-lg font-bold font-display tracking-tight text-foreground">No Order Selected</h3>
                  <p className="text-xs text-muted-foreground leading-relaxed">
                    Select an order from the list on the left to inspect items, allocate split shipments, manage waybills, and update transit status.
                  </p>
                </div>
              ) : selectedOrder ? (
                /* Active Order Inspector */
                <div className="w-full flex-1 space-y-0">

                  {/* 1. Header Toolbar */}
                  <div className="p-6 border-b border-border/60 bg-muted/10 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
                    <div className="flex items-center gap-3">
                      <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => setSelectedOrderId(null)}
                        className="rounded-xl h-8 w-8 text-muted-foreground hover:text-foreground lg:hidden"
                      >
                        <ArrowLeft className="h-4 w-4" />
                      </Button>
                      <div>
                        <div className="flex items-center gap-2">
                          <h4 className="text-lg font-bold font-mono text-foreground">
                            {selectedOrder.number}
                          </h4>
                          <StatusBadge status={selectedOrder.status} />
                        </div>
                        <p className="text-xs text-muted-foreground mt-0.5">
                          Ordered on {new Date(selectedOrder.created_at).toLocaleString()}
                        </p>
                      </div>
                    </div>

                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => setSelectedOrderId(null)}
                      className="rounded-xl text-xs h-8 hidden lg:flex cursor-pointer"
                    >
                      Close Inspector
                    </Button>
                  </div>

                  {/* 2. State Transition Action Bar */}
                  <div className="p-6 border-b border-border/60 space-y-4">
                    <div className="flex items-center justify-between">
                      <h5 className="text-xs font-bold uppercase tracking-wider text-muted-foreground">Fulfillment Actions</h5>
                      <span className="text-[11px] font-medium text-muted-foreground">Status: <strong className="text-foreground capitalize">{selectedOrder.status}</strong></span>
                    </div>

                    {/* CASE A: PENDING PAYMENT */}
                    {selectedOrder.status === 'pending' && (
                      <div className="p-4 rounded-xl border border-amber-500/20 bg-amber-500/5 text-amber-900 dark:text-amber-200 text-xs flex items-center justify-between">
                        <span>Awaiting payment settlement from customer. Status will auto-advance once paid.</span>
                      </div>
                    )}

                    {/* CASE B: CONFIRMED - READY FOR PACKAGING */}
                    {selectedOrder.status === 'confirmed' && (
                      <div className="space-y-3">
                        <div className="p-4 rounded-xl border border-blue-500/20 bg-blue-500/5 text-xs text-blue-950 dark:text-blue-200 flex items-center justify-between">
                          <span>Payment verified. Ready to initiate order assembly and warehouse processing.</span>
                          <Button
                            size="sm"
                            disabled={submitting}
                            onClick={() => handleStartProcessing(selectedOrder.id)}
                            className="bg-primary text-primary-foreground hover:bg-primary/90 text-xs rounded-xl h-8 font-semibold cursor-pointer"
                          >
                            <Sprout className="h-3.5 w-3.5 mr-1" /> Start Processing
                          </Button>
                        </div>
                      </div>
                    )}

                    {/* CASE C: PROCESSING - MULTI / SPLIT SHIPMENT DISPATCH (BY SHOP & SPLITTABLE) */}
                    {selectedOrder.status === 'processing' && (
                      <div className="border border-border/80 bg-background p-5 rounded-xl space-y-5">
                        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2">
                          <div className="space-y-0.5">
                            <h5 className="text-sm font-semibold text-foreground flex items-center gap-1.5">
                              <Truck className="h-4 w-4 text-primary" /> Dispatch Shipping & Assign Courier
                            </h5>
                            <p className="text-xs text-muted-foreground">
                              Shipments are grouped by shop warehouse by default. You can split products into additional shipments if needed.
                            </p>
                          </div>
                          <Button
                            type="button"
                            variant="outline"
                            size="sm"
                            onClick={() => handleAddShipmentGroup()}
                            className="rounded-xl text-xs h-8 gap-1 self-start sm:self-auto border-primary/40 text-primary hover:bg-primary/5 cursor-pointer"
                          >
                            <Plus className="h-3.5 w-3.5" /> Add Split Shipment
                          </Button>
                        </div>

                        {/* Shipment Groups List */}
                        <div className="space-y-4">
                          {shipmentGroups.map((group, groupIdx) => {
                            // Filter items that match this shop or all order items
                            const shopItems = group.shop_id
                              ? selectedOrder.items.filter(i => i.shop_id === group.shop_id)
                              : selectedOrder.items;

                            return (
                              <div
                                key={group.id}
                                className="border border-border/80 rounded-xl p-4 bg-muted/10 space-y-4 relative"
                              >
                                <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2 border-b border-border/40 pb-2.5">
                                  <div className="flex items-center gap-2 flex-wrap">
                                    <span className="h-5 w-5 rounded-full bg-primary/10 text-primary text-[11px] font-bold flex items-center justify-center">
                                      {groupIdx + 1}
                                    </span>
                                    <span className="text-xs font-bold text-foreground">
                                      Shipment #{groupIdx + 1}
                                    </span>
                                    {group.shop_name && (
                                      <span className="inline-flex items-center gap-1 text-[11px] bg-secondary text-secondary-foreground px-2 py-0.5 rounded-md font-medium">
                                        <Store className="h-3 w-3" /> {group.shop_name}
                                      </span>
                                    )}
                                    <span className="text-[10px] text-muted-foreground font-mono">
                                      ({group.item_ids.length} item{group.item_ids.length !== 1 ? 's' : ''})
                                    </span>
                                  </div>

                                  <div className="flex items-center gap-2">
                                    {shipmentGroups.length > 1 && (
                                      <Button
                                        type="button"
                                        variant="ghost"
                                        size="sm"
                                        onClick={() => handleRemoveShipmentGroup(group.id)}
                                        className="h-6 px-2 text-destructive hover:bg-destructive/10 text-xs rounded-lg cursor-pointer"
                                      >
                                        <Trash2 className="h-3.5 w-3.5 mr-1" /> Remove
                                      </Button>
                                    )}
                                  </div>
                                </div>

                                {/* Items Selector in this group */}
                                <div className="space-y-2">
                                  <div className="flex justify-between items-center">
                                    <label className="text-[11px] font-semibold text-muted-foreground uppercase flex items-center gap-1">
                                      <Boxes className="h-3.5 w-3.5" /> Assigned Products ({group.shop_name || 'All Shops'})
                                    </label>
                                    <span className="text-[10px] text-muted-foreground">
                                      {group.item_ids.length} of {shopItems.length} selected
                                    </span>
                                  </div>

                                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
                                    {shopItems.map(item => {
                                      const isChecked = group.item_ids.includes(item.id);
                                      const otherGroupIdx = shipmentGroups.findIndex(
                                        g => g.id !== group.id && g.item_ids.includes(item.id)
                                      );

                                      return (
                                        <label
                                          key={item.id}
                                          className={`flex items-start gap-2.5 p-2.5 rounded-lg border text-xs cursor-pointer transition-all ${
                                            isChecked
                                              ? 'border-primary/60 bg-primary/5 text-foreground font-semibold ring-1 ring-primary/20'
                                              : otherGroupIdx >= 0
                                              ? 'border-border/40 opacity-55 hover:opacity-80 bg-muted/20 text-muted-foreground'
                                              : 'border-border/60 hover:bg-muted/30 text-muted-foreground'
                                          }`}
                                        >
                                          <input
                                            type="checkbox"
                                            checked={isChecked}
                                            onChange={() => handleToggleItemInGroup(group.id, item.id)}
                                            className="mt-0.5 rounded text-primary focus:ring-primary h-3.5 w-3.5 cursor-pointer"
                                          />
                                          <div className="flex-1 min-w-0">
                                            <div className="truncate text-xs">{item.product_name}</div>
                                            <div className="text-[10px] text-muted-foreground font-normal flex items-center justify-between">
                                              <span>Qty: {item.quantity}</span>
                                              {otherGroupIdx >= 0 && (
                                                <span className="text-amber-600 dark:text-amber-400 font-medium">
                                                  (In Pkg #{otherGroupIdx + 1})
                                                </span>
                                              )}
                                            </div>
                                          </div>
                                        </label>
                                      );
                                    })}
                                  </div>
                                </div>

                                {/* Logistics Fields for this shipment */}
                                <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 pt-2">
                                  <div className="space-y-1.5">
                                    <label className="text-[11px] font-semibold text-muted-foreground">Fulfillment Method</label>
                                    <Select
                                      value={group.fulfillment_method}
                                      onValueChange={(val) => handleUpdateGroupField(group.id, 'fulfillment_method', val)}
                                    >
                                      <SelectTrigger className="rounded-xl h-8 text-xs bg-background">
                                        <SelectValue />
                                      </SelectTrigger>
                                      <SelectContent>
                                        <SelectItem value="courier">Courier Logistics</SelectItem>
                                        <SelectItem value="self_delivery">Self Delivery (Shop Rider)</SelectItem>
                                      </SelectContent>
                                    </Select>
                                  </div>

                                  {group.fulfillment_method === 'courier' && (
                                    <>
                                      <div className="space-y-1.5">
                                        <label className="text-[11px] font-semibold text-muted-foreground">Courier Brand</label>
                                        <Input
                                          value={group.courier}
                                          onChange={(e) => handleUpdateGroupField(group.id, 'courier', e.target.value)}
                                          placeholder="JNE, POS, SiCepat..."
                                          className="rounded-xl h-8 text-xs bg-background"
                                        />
                                      </div>
                                      <div className="space-y-1.5">
                                        <label className="text-[11px] font-semibold text-muted-foreground">Service Level</label>
                                        <Input
                                          value={group.service}
                                          onChange={(e) => handleUpdateGroupField(group.id, 'service', e.target.value)}
                                          placeholder="REG, YES, etc."
                                          className="rounded-xl h-8 text-xs bg-background"
                                        />
                                      </div>
                                      <div className="space-y-1.5">
                                        <label className="text-[11px] font-semibold text-muted-foreground">Waybill Tracking Number</label>
                                        <Input
                                          value={group.tracking_number}
                                          onChange={(e) => handleUpdateGroupField(group.id, 'tracking_number', e.target.value)}
                                          placeholder="Insert shipping receipt ID"
                                          className="rounded-xl h-8 text-xs bg-background"
                                        />
                                      </div>
                                    </>
                                  )}
                                </div>
                              </div>
                            );
                          })}
                        </div>

                        <Button
                          disabled={submitting}
                          onClick={() => handleDispatchOrder(selectedOrder.id)}
                          className="rounded-xl bg-primary text-primary-foreground hover:bg-primary/90 text-xs font-semibold h-9 w-full shadow-sm cursor-pointer"
                        >
                          <Truck className="h-4 w-4 mr-1.5" /> Dispatch {shipmentGroups.length > 1 ? `${shipmentGroups.length} Shipments` : 'Shipment'} & Mark Shipped
                        </Button>
                      </div>
                    )}

                    {/* CASE D: SHIPPED - IN TRANSIT LOGISTICS */}
                    {selectedOrder.status === 'shipped' && currentShipments.length > 0 && (
                      <div className="border border-border/80 bg-background p-5 rounded-xl space-y-5">
                        <div className="space-y-1">
                          <h5 className="text-sm font-semibold text-foreground flex items-center gap-1.5">
                            <Layers className="h-4 w-4 text-primary" /> Active Logistics & Shipment Status
                          </h5>
                          <p className="text-xs text-muted-foreground">
                            Manage waybill information and post transit timeline updates per package.
                          </p>
                        </div>

                        {/* Shipments List Accordion / Cards */}
                        <div className="space-y-4">
                          {currentShipments.map((shipment, sIdx) => {
                            const isEditing = editingShipmentId === shipment.id;
                            const isSelectedForTimeline = activeShipmentId === shipment.id;

                            const assignedItems = selectedOrder.items.filter(
                              item => item.shipment_id === shipment.id || shipment.item_ids?.includes(item.id)
                            );

                            return (
                              <div
                                key={shipment.id}
                                className={`border rounded-xl p-4 transition-all ${
                                  isSelectedForTimeline
                                    ? 'border-primary/60 bg-muted/20'
                                    : 'border-border/60 bg-muted/5'
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
                                    </div>
                                    <div className="text-[11px] text-muted-foreground mt-0.5">
                                      Waybill:{' '}
                                      <strong className="font-mono text-foreground">
                                        {shipment.tracking_number || 'Self Delivery'}
                                      </strong>{' '}
                                      · {assignedItems.length} item{assignedItems.length !== 1 ? 's' : ''}
                                    </div>
                                  </div>

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
                                </div>

                                {isEditing ? (
                                  /* EDIT WAYBILL SUB-FORM */
                                  <div className="space-y-3 p-3 mt-3 border border-border/40 rounded-xl bg-background">
                                    <h6 className="text-xs font-bold text-foreground">Edit Shipping Waybill</h6>
                                    <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
                                      <div className="space-y-1">
                                        <label className="text-[10px] font-semibold text-muted-foreground uppercase">Courier</label>
                                        <Input
                                          value={editCourier}
                                          onChange={(e) => setEditCourier(e.target.value)}
                                          className="h-8 text-xs rounded-lg"
                                        />
                                      </div>
                                      <div className="space-y-1">
                                        <label className="text-[10px] font-semibold text-muted-foreground uppercase">Service</label>
                                        <Input
                                          value={editService}
                                          onChange={(e) => setEditService(e.target.value)}
                                          className="h-8 text-xs rounded-lg"
                                        />
                                      </div>
                                      <div className="space-y-1">
                                        <label className="text-[10px] font-semibold text-muted-foreground uppercase">Tracking No</label>
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
                                      onClick={() => handleUpdateWaybill(shipment.id)}
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
                                        <label className="text-[11px] font-semibold text-muted-foreground">New Transit Status</label>
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
                                        <label className="text-[11px] font-semibold text-muted-foreground">Location (Optional)</label>
                                        <Input
                                          value={transitLocation}
                                          onChange={(e) => setTransitLocation(e.target.value)}
                                          placeholder="Warehouse, Hub City, etc."
                                          className="rounded-xl h-8 text-xs bg-background"
                                        />
                                      </div>
                                    </div>
                                    <div className="space-y-1.5">
                                      <label className="text-[11px] font-semibold text-muted-foreground">Description/Note (Optional)</label>
                                      <Input
                                        value={transitDescription}
                                        onChange={(e) => setTransitDescription(e.target.value)}
                                        placeholder="E.g. Package hand over to courier JNE"
                                        className="rounded-xl h-8 text-xs bg-background"
                                      />
                                    </div>
                                    <Button
                                      disabled={submitting}
                                      onClick={() => handleUpdateShipmentStatus(shipment.id)}
                                      className="rounded-xl bg-primary text-primary-foreground hover:bg-primary/90 text-xs font-semibold h-8 w-full cursor-pointer"
                                    >
                                      Update Package #{sIdx + 1} Status
                                    </Button>
                                  </div>
                                )}
                              </div>
                            );
                          })}
                        </div>
                      </div>
                    )}

                    {/* CASE E: DELIVERED (COMPLETED) */}
                    {selectedOrder.status === 'delivered' && (
                      <div className="p-4 rounded-xl border border-emerald-500/20 bg-emerald-500/5 text-emerald-950 dark:text-emerald-200 text-xs flex items-center gap-2">
                        <CheckCircle2 className="h-4 w-4 text-emerald-600" />
                        <span>Order fulfillment complete. All packages delivered to customer.</span>
                      </div>
                    )}

                    {/* CASE F: CANCELLED */}
                    {selectedOrder.status === 'cancelled' && (
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
                            {selectedOrder.items.map((item) => {
                              const matchedShipmentIndex = currentShipments.findIndex(
                                s => s.id === item.shipment_id || s.item_ids?.includes(item.id)
                              );

                              return (
                                <tr key={item.id} className="text-foreground hover:bg-muted/10">
                                  <td className="p-3">
                                    <div className="font-semibold text-foreground">{item.product_name}</div>
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
                            <span className="font-medium text-foreground">{formatCurrency(selectedOrder.subtotal)}</span>
                          </div>
                          <div className="flex justify-between text-muted-foreground">
                            <span>Shipping Fee</span>
                            <span className="font-medium text-foreground">{formatCurrency(selectedOrder.shipping_fee)}</span>
                          </div>
                          <div className="border-t border-border/60 pt-2.5 mt-1 flex justify-between text-sm font-bold text-foreground">
                            <span>Grand Total</span>
                            <span className="text-primary">{formatCurrency(selectedOrder.total)}</span>
                          </div>

                          {selectedOrder.payment && (
                            <div className="border-t border-border/40 pt-2.5 mt-2 space-y-1.5 font-medium text-[10px] text-muted-foreground">
                              <div>Payment Method: <span className="text-foreground font-semibold uppercase">{selectedOrder.payment.provider}</span></div>
                              <div>Transaction Status: <span className="text-foreground font-semibold uppercase">{selectedOrder.payment.status}</span></div>
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
                              <MapPin className="h-3.5 w-3.5 text-muted-foreground" /> {selectedOrder.address ? 'Delivery Address' : 'Address Metadata'}
                            </div>
                            {selectedOrder.address ? (
                              <div className="text-foreground space-y-1 leading-relaxed text-[11px] font-medium">
                                <p><span className="font-semibold text-muted-foreground">Receiver:</span> {selectedOrder.address.receiver_name}</p>
                                <p><span className="font-semibold text-muted-foreground">Phone:</span> {selectedOrder.address.phone}</p>
                                <p><span className="font-semibold text-muted-foreground">Address:</span> {selectedOrder.address.full_address}</p>
                                {selectedOrder.address.postal_code && (
                                  <p><span className="font-semibold text-muted-foreground">Postal Code:</span> {selectedOrder.address.postal_code}</p>
                                )}
                              </div>
                            ) : (
                              <div className="text-muted-foreground space-y-1 leading-relaxed text-[11px]">
                                <p><span className="font-semibold text-foreground">User UUID:</span> <span className="font-mono">{selectedOrder.user_id}</span></p>
                                <p><span className="font-semibold text-foreground">Address ID:</span> <span className="font-mono">{selectedOrder.address_id}</span></p>
                              </div>
                            )}
                          </div>

                          {currentShipments.length > 0 && (
                            <div className="border-t border-border/40 pt-2.5 mt-2 space-y-2 text-[11px] text-muted-foreground">
                              <div className="font-bold text-foreground">Packages Dispatched:</div>
                              {currentShipments.map((s, idx) => (
                                <div key={s.id} className="p-2 rounded bg-background border border-border/40 space-y-0.5">
                                  <div className="font-semibold text-foreground flex justify-between">
                                    <span>Pkg #{idx + 1}: {s.courier.toUpperCase()} {s.service.toUpperCase()}</span>
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
                              onClick={() => selectedOrder && loadLiveTracking(selectedOrder.id)}
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
                          <div className="flex items-start gap-2 p-3 rounded-xl bg-amber-500/10 border border-amber-500/20 text-amber-600 dark:text-amber-400 text-xs font-medium leading-relaxed">
                            <AlertTriangle className="h-4 w-4 shrink-0 text-amber-500 mt-0.5" />
                            <span>{liveTracking.warning}</span>
                          </div>
                        )}

                        {(() => {
                          const activeShipment = currentShipments.find(s => s.id === activeShipmentId) || currentShipments[0];
                          
                          // Prefer live merged timeline from GET /orders/{id}/tracking (which includes live Komerce checkpoints)
                          const liveEvents = liveTracking?.timeline || [];
                          const internalEvents = activeShipment?.events || [];

                          const eventsToDisplay = liveEvents.length > 0
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
                                      <span className={`absolute -left-[31px] top-1 h-4 w-4 rounded-full border-2 bg-background flex items-center justify-center transition-colors ${
                                        isLatest ? 'border-primary ring-4 ring-primary/10' : 'border-muted-foreground/45'
                                      }`}>
                                        <span className={`h-1.5 w-1.5 rounded-full ${isLatest ? 'bg-primary' : 'bg-muted-foreground/45'}`} />
                                      </span>

                                      <div className="space-y-1">
                                        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-1">
                                          <div className="flex items-center gap-2">
                                            <span className={`text-xs font-bold uppercase tracking-wider ${isLatest ? 'text-primary' : 'text-muted-foreground'}`}>
                                              {event.status}
                                            </span>
                                            {event.location && (
                                              <span className="text-[10px] bg-muted px-1.5 py-0.5 rounded text-muted-foreground font-medium">
                                                {event.location}
                                              </span>
                                            )}
                                          </div>
                                          <span className="text-[10px] text-muted-foreground font-mono">
                                            {new Date(event.timestamp).toLocaleString(undefined, { dateStyle: 'short', timeStyle: 'short' })}
                                          </span>
                                        </div>
                                        <p className="text-xs text-muted-foreground leading-relaxed">
                                          {event.description}
                                        </p>
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
                              <span>No shipment updates logged yet for {activeShipment ? `${activeShipment.courier.toUpperCase()} (${activeShipment.tracking_number || 'Self Delivery'})` : 'this package'}.</span>
                              <span className="text-[10px] text-muted-foreground/60">Use the logistics panel above to post the first shipment status update.</span>
                            </div>
                          );
                        })()}
                      </div>
                    )}

                  </div>

                </div>
              ) : null}

            </div>

          </div>

        </div>

      </div>
    </div>
  );
}
