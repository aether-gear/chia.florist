import { useState, useEffect } from 'react';
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
  AlertCircle,
  Sprout
} from 'lucide-react';
import { Button } from '../../components/ui/button';
import { Card, CardContent } from '../../components/ui/card';
import { useOrdersViewModel } from '../../viewmodels/useOrdersViewModel';
import { useOrderActionsViewModel } from '../../viewmodels/useOrderActionsViewModel';
import EmptyState from '../../components/EmptyState';
import SearchInput from '../../components/SearchInput';
import StatusBadge from '../../components/StatusBadge';
import Pagination from '../../components/Pagination';
import { Skeleton } from '../../components/ui/skeleton';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '../../components/ui/select';
import { Input } from '../../components/ui/input';

export default function OrdersPage() {
  const {
    data,
    loading,
    error,
    page,
    limit,
    sort,
    searchNumber,
    statusFilter,
    setPage,
    setSort,
    setSearchNumber,
    setStatusFilter,
    refresh
  } = useOrdersViewModel();

  const {
    submitting,
    confirmPayment,
    rejectPayment,
    updateOrderStatus,
    updateShipmentStatus,
    updateShipmentDetails
  } = useOrderActionsViewModel();

  const [selectedOrderId, setSelectedOrderId] = useState<string | null>(null);
  const [isEditingShipment, setIsEditingShipment] = useState<boolean>(false);

  // Dispatch Form States
  const [dispatchMethod, setDispatchMethod] = useState<string>('courier');
  const [courierCode, setCourierCode] = useState<string>('');
  const [courierService, setCourierService] = useState<string>('');
  const [trackingNumber, setTrackingNumber] = useState<string>('');

  // Shipment Transit Update States
  const [transitStatus, setTransitStatus] = useState<string>('packed');
  const [transitDescription, setTransitDescription] = useState<string>('');
  const [transitLocation, setTransitLocation] = useState<string>('');

  // Shipment Waybill Edit States
  const [editCourier, setEditCourier] = useState<string>('');
  const [editService, setEditService] = useState<string>('');
  const [editTracking, setEditTracking] = useState<string>('');

  const selectedOrder = data?.orders.find(o => o.id === selectedOrderId);

  // Sync dispatch forms when selected order changes
  useEffect(() => {
    if (selectedOrder) {
      const firstItem = selectedOrder.items[0];
      setCourierCode(selectedOrder.shipment?.courier || firstItem?.courier_code || '');
      setCourierService(selectedOrder.shipment?.service || firstItem?.courier_service || '');
      setTrackingNumber(selectedOrder.shipment?.tracking_number || '');
      setDispatchMethod(selectedOrder.shipment?.fulfillment_method || 'courier');

      // Sync waybill edit states
      setEditCourier(selectedOrder.shipment?.courier || firstItem?.courier_code || '');
      setEditService(selectedOrder.shipment?.service || firstItem?.courier_service || '');
      setEditTracking(selectedOrder.shipment?.tracking_number || '');

      // Default next shipment status based on current
      setTransitStatus(selectedOrder.shipment?.status || 'packed');
      setTransitDescription('');
      setTransitLocation('');
      setIsEditingShipment(false);
    }
  }, [selectedOrderId, selectedOrder]);

  const handleSort = () => {
    const currentDirection = sort.split(':')[1];
    const newDirection = currentDirection === 'desc' ? 'asc' : 'desc';
    setSort(`latest:${newDirection}`);
    setPage(1);
  };

  const handleConfirmPayment = async (paymentId: string) => {
    if (confirm('Are you sure you want to confirm this payment manually? This will mark the order as Confirmed.')) {
      try {
        await confirmPayment(paymentId);
        refresh();
      } catch (err) {
        // Error toast is handled in useOrderActionsViewModel
      }
    }
  };

  const handleRejectPayment = async (paymentId: string) => {
    if (confirm('Are you sure you want to reject this payment? This will mark the order as Cancelled.')) {
      try {
        await rejectPayment(paymentId);
        refresh();
      } catch (err) {
        // Error toast is handled in useOrderActionsViewModel
      }
    }
  };

  const handleStartProcessing = async (orderId: string) => {
    try {
      await updateOrderStatus(orderId, 'processing');
      refresh();
    } catch (err) {
      // Handled in ViewModel
    }
  };

  const handleDispatchOrder = async (orderId: string) => {
    if (dispatchMethod === 'courier' && !trackingNumber.trim()) {
      alert('Tracking number is required for courier fulfillment.');
      return;
    }
    try {
      await updateOrderStatus(orderId, 'shipped', trackingNumber, dispatchMethod);
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

  const handleUpdateWaybill = async (shipmentId: string) => {
    try {
      await updateShipmentDetails(shipmentId, {
        tracking_number: editTracking,
        courier: editCourier,
        service: editService
      });
      setIsEditingShipment(false);
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

  const statusTabs = [
    { value: 'all', label: 'All Orders' },
    { value: 'pending', label: 'Awaiting Payment' },
    { value: 'confirmed', label: 'To Process' },
    { value: 'processing', label: 'In Packaging' },
    { value: 'shipped', label: 'Shipped' },
    { value: 'delivered', label: 'Delivered' },
    { value: 'cancelled', label: 'Cancelled' },
  ];

  return (
    <div className="flex-col md:flex">
      <div className="flex-1 space-y-6 p-6 sm:p-8 lg:p-12 animate-in fade-in duration-300">

        {/* Header */}
        <div className="flex items-center justify-between pb-4 border-b border-border/60">
          <div>
            <h2 className="text-3xl font-bold font-display tracking-tight text-foreground">Order Workspace</h2>
            <p className="text-muted-foreground text-sm">
              Verify payments, fulfill shipments, and manage orders
            </p>
          </div>
          <Button
            variant="outline"
            size="icon"
            className="rounded-xl border-border text-foreground hover:text-primary hover:bg-primary/5 h-10 w-10"
            onClick={() => refresh()}
            disabled={loading}
          >
            <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
          </Button>
        </div>

        {/* Full-width Search and Status Filter Tabs */}
        <div className={`flex flex-col gap-4 justify-between bg-zinc-50/20 dark:bg-slate-900/10 p-4 rounded-2xl border border-border/60 ${
          selectedOrderId ? 'hidden lg:flex' : 'flex'
        }`}>
          <div className="w-full md:max-w-xs lg:max-w-sm">
            <SearchInput
              value={searchNumber}
              onChange={(val) => {
                setSearchNumber(val);
                setPage(1);
              }}
              placeholder="Search by Order Number..."
            />
          </div>
          <div className="flex gap-1 overflow-x-auto pb-1 w-full md:w-auto scrollbar-none">
            {statusTabs.map((tab) => {
              const isActive = (statusFilter || 'all') === tab.value;
              return (
                <button
                  key={tab.value}
                  onClick={() => {
                    setStatusFilter(tab.value === 'all' ? '' : tab.value);
                    setPage(1);
                    setSelectedOrderId(null); // Clear selected order
                  }}
                  className={`px-3 py-1.5 text-xs font-semibold whitespace-nowrap rounded-lg transition-all ${
                    isActive
                      ? 'bg-primary/10 text-primary'
                      : 'text-muted-foreground hover:text-foreground hover:bg-muted/50'
                  }`}
                >
                  {tab.label}
                </button>
              );
            })}
          </div>
        </div>

        {/* Workspace Layout */}
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">

          {/* LEFT PANE: Master Order List */}
          <div className={`lg:col-span-5 flex flex-col space-y-4 ${selectedOrderId ? 'hidden lg:flex' : 'flex'}`}>

            {/* Sorting bar */}
            <div className="flex items-center justify-between px-1 text-xs text-muted-foreground">
              <span>{data?.total ? `Found ${data.total} orders` : 'No orders found'}</span>
              <button
                onClick={handleSort}
                className="flex items-center gap-1 hover:text-foreground transition-colors font-medium"
              >
                  Date <ArrowUpDown className="h-3 w-3" />
                </button>
              </div>

            {/* List Container */}
            <div className="space-y-3 pr-1">
              {loading ? (
                Array.from({ length: 4 }).map((_, i) => (
                  <Card key={`skeleton-${i}`} className="border border-border/50 shadow-none">
                    <CardContent className="p-4 space-y-3">
                      <div className="flex justify-between">
                        <Skeleton className="h-4 w-28 bg-muted animate-pulse" />
                        <Skeleton className="h-4 w-16 bg-muted animate-pulse" />
                      </div>
                      <Skeleton className="h-3 w-40 bg-muted animate-pulse" />
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

                  return (
                    <Card
                      key={order.id}
                      onClick={() => setSelectedOrderId(order.id)}
                      className={`cursor-pointer transition-all border shadow-none hover:border-primary/50 select-none ${
                        isSelected
                          ? 'border-primary/60 bg-primary/5 ring-1 ring-primary/45'
                          : 'border-border/60 bg-background'
                      }`}
                    >
                      <CardContent className="p-4 space-y-3">
                        <div className="flex justify-between items-start">
                          <span className="font-semibold font-display text-sm tracking-tight text-foreground">
                            {order.number}
                          </span>
                          <span className="text-[10px] text-muted-foreground">
                            {new Date(order.created_at).toLocaleDateString(undefined, { dateStyle: 'short' })}
                          </span>
                        </div>

                        <p className="text-xs text-muted-foreground line-clamp-1">
                          {itemPreview}
                        </p>

                        <div className="flex flex-wrap gap-1.5 pt-1">
                          <StatusBadge status={order.status} className="scale-90 origin-left" />
                          {order.payment && (
                            <StatusBadge status={`Pay: ${order.payment.status}`} className="scale-90 origin-left bg-zinc-100 dark:bg-zinc-800 text-foreground" />
                          )}
                          {order.shipment && (
                            <StatusBadge status={`Ship: ${order.shipment.status}`} className="scale-90 origin-left bg-zinc-100 dark:bg-zinc-800 text-foreground" />
                          )}
                        </div>

                        <div className="flex justify-between items-center border-t border-border/40 pt-2 mt-2">
                          <span className="text-[10px] text-muted-foreground font-mono">
                            ID: {order.id.slice(0, 8)}...
                          </span>
                          <span className="text-sm font-bold text-primary">
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
            <div className="pt-2">
              <Pagination
                currentPage={page}
                totalPages={Math.ceil((data?.total || 0) / limit)}
                totalItems={data?.total || 0}
                limit={limit}
                onPageChange={setPage}
                itemNamePlural="orders"
              />
            </div>
          </div>

          {/* RIGHT PANE: Detail Workspace */}
          <div className={`lg:col-span-7 border border-border/80 rounded-2xl bg-zinc-50/15 dark:bg-slate-900/10 flex flex-col lg:sticky lg:top-24 lg:self-start lg:max-h-[calc(100vh-10rem)] overflow-hidden ${
            !selectedOrderId ? 'hidden lg:flex items-center justify-center p-12 text-center text-muted-foreground' : 'flex'
          }`}>

            {!selectedOrderId ? (
              <div className="space-y-4 max-w-sm flex flex-col items-center">
                <PackageOpen className="h-12 w-12 text-muted-foreground/60 stroke-[1.5]" />
                <h3 className="text-lg font-bold font-display tracking-tight text-foreground">No Order Selected</h3>
                <p className="text-xs text-muted-foreground">
                  Select an order from the left pane to inspect its items, audit transactions, and process fulfillment state transitions.
                </p>
              </div>
            ) : !selectedOrder ? (
              <div className="p-12 text-center text-destructive">
                <AlertCircle className="h-12 w-12 mx-auto mb-4" />
                <p className="text-sm font-semibold">Error finding order details.</p>
                <Button variant="ghost" onClick={() => setSelectedOrderId(null)} className="mt-4">Back to List</Button>
              </div>
            ) : (
              <div className="flex flex-col flex-1 min-h-0 divide-y divide-border/60 overflow-hidden">

                {/* 1. Detail Header */}
                <div className="p-6 flex items-center justify-between bg-zinc-50/40 dark:bg-slate-900/40">
                  <div className="space-y-1">
                    <div className="flex items-center gap-3">
                      <button
                        onClick={() => setSelectedOrderId(null)}
                        className="lg:hidden p-1 rounded-lg hover:bg-muted text-muted-foreground transition-colors"
                      >
                        <ArrowLeft className="h-5 w-5" />
                      </button>
                      <h3 className="text-xl font-bold font-display tracking-tight text-foreground">
                        {selectedOrder.number}
                      </h3>
                    </div>
                    <p className="text-xs text-muted-foreground">
                      Placed on {new Date(selectedOrder.created_at).toLocaleString(undefined, { dateStyle: 'medium', timeStyle: 'short' })}
                    </p>
                  </div>

                  <div className="flex flex-col sm:flex-row gap-2 items-end sm:items-center">
                    <StatusBadge status={selectedOrder.status} />
                  </div>
                </div>

                {/* Scrollable Container (Banner + Details) */}
                <div className="flex-1 overflow-y-auto min-h-0 divide-y divide-border/60">

                  {/* 2. Interactive Workflow Action Banner */}
                  <div className="p-6 bg-primary/[0.02] dark:bg-primary/[0.01]">
                  <h4 className="text-xs font-bold uppercase tracking-wider text-muted-foreground mb-3 flex items-center gap-1.5 flex-wrap">
                    <Sprout className="h-3.5 w-3.5 text-primary" /> Logistics & <Sprout className="h-3.5 w-3.5 text-primary" /> Payment
                  </h4>

                  {/* SUBMITTING STATE LOADER */}
                  {submitting && (
                    <div className="py-4 flex items-center justify-center gap-2 text-xs text-muted-foreground animate-pulse">
                      <RefreshCw className="h-4 w-4 animate-spin text-primary" />
                      Updating order records, please wait...
                    </div>
                  )}

                  {/* ACTION SWITCH CASE */}
                  {!submitting && (
                    <div className="animate-in fade-in duration-200">

                      {/* CASE A: PENDING MANUAL PAYMENT */}
                      {selectedOrder.status === 'pending' && selectedOrder.payment && selectedOrder.payment.provider === 'manual' && selectedOrder.payment.status === 'pending' && (
                        <div className="border border-amber-200 dark:border-amber-950/40 bg-amber-50/30 dark:bg-amber-950/10 p-4 rounded-xl space-y-4">
                          <div className="flex items-start gap-3">
                            <CreditCard className="h-5 w-5 text-amber-600 dark:text-amber-400 mt-0.5" />
                            <div className="space-y-1">
                              <h5 className="text-sm font-semibold text-foreground">Awaiting Bank Transfer Verification</h5>
                              <p className="text-xs text-muted-foreground">
                                Customer reported paying manual bank transfer. Please audit your merchant account before confirming.
                              </p>
                            </div>
                          </div>

                          <div className="text-xs space-y-1.5 p-3 rounded-lg bg-background border border-border/40 font-medium">
                            <div className="flex justify-between"><span className="text-muted-foreground">Provider:</span> <span className="uppercase text-foreground">{selectedOrder.payment.provider}</span></div>
                            <div className="flex justify-between"><span className="text-muted-foreground">Amount:</span> <span className="text-primary font-bold">{formatCurrency(selectedOrder.payment.amount)}</span></div>
                          </div>

                          <div className="flex gap-3 pt-1">
                            <Button
                              onClick={() => handleConfirmPayment(selectedOrder.payment!.id)}
                              className="rounded-xl flex-1 bg-primary text-primary-foreground hover:bg-primary/90 text-xs font-semibold h-9"
                            >
                              <CheckCircle2 className="h-4 w-4 mr-1.5" /> Confirm Payment
                            </Button>
                            <Button
                              variant="outline"
                              onClick={() => handleRejectPayment(selectedOrder.payment!.id)}
                              className="rounded-xl flex-1 border-destructive text-destructive hover:bg-destructive/5 text-xs font-semibold h-9"
                            >
                              <XCircle className="h-4 w-4 mr-1.5" /> Reject Payment
                            </Button>
                          </div>
                        </div>
                      )}

                      {/* CASE B: CONFIRMED - READY TO PROCESS */}
                      {selectedOrder.status === 'confirmed' && (
                        <div className="border border-border/80 bg-background p-4 rounded-xl flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
                          <div className="space-y-1 max-w-md">
                            <h5 className="text-sm font-semibold text-foreground">Ready for Fulfilling</h5>
                            <p className="text-xs text-muted-foreground">
                              Payment cleared. Move order to packaging to begin preparing the flowers.
                            </p>
                          </div>
                          <Button
                            onClick={() => handleStartProcessing(selectedOrder.id)}
                            className="rounded-xl bg-primary text-primary-foreground hover:bg-primary/90 text-xs font-semibold h-9 w-full sm:w-auto"
                          >
                            <Truck className="h-4 w-4 mr-1.5" /> Start Processing
                          </Button>
                        </div>
                      )}

                      {/* CASE C: PROCESSING - DISPATCH SHIPMENT */}
                      {selectedOrder.status === 'processing' && (
                        <div className="border border-border/80 bg-background p-5 rounded-xl space-y-4">
                          <div className="space-y-1">
                            <h5 className="text-sm font-semibold text-foreground">Dispatch Shipping & Book Courier</h5>
                            <p className="text-xs text-muted-foreground">
                              Fulfill order items and update courier details. If using a manual courier, insert the waybill tracking number below.
                            </p>
                          </div>

                          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                            <div className="space-y-1.5">
                              <label className="text-xs font-semibold text-muted-foreground">Fulfillment Method</label>
                              <Select value={dispatchMethod} onValueChange={setDispatchMethod}>
                                <SelectTrigger className="rounded-xl h-9 text-xs">
                                  <SelectValue />
                                </SelectTrigger>
                                <SelectContent>
                                  <SelectItem value="courier">Courier Logistics</SelectItem>
                                  <SelectItem value="self_delivery">Self Delivery (Shop Rider)</SelectItem>
                                </SelectContent>
                              </Select>
                            </div>

                            {dispatchMethod === 'courier' && (
                              <>
                                <div className="space-y-1.5">
                                  <label className="text-xs font-semibold text-muted-foreground">Courier Brand</label>
                                  <Input
                                    value={courierCode}
                                    onChange={(e) => setCourierCode(e.target.value)}
                                    placeholder="JNE, POS, etc."
                                    className="rounded-xl h-9 text-xs"
                                  />
                                </div>
                                <div className="space-y-1.5">
                                  <label className="text-xs font-semibold text-muted-foreground">Service Level</label>
                                  <Input
                                    value={courierService}
                                    onChange={(e) => setCourierService(e.target.value)}
                                    placeholder="REG, YES, etc."
                                    className="rounded-xl h-9 text-xs"
                                  />
                                </div>
                                <div className="space-y-1.5">
                                  <label className="text-xs font-semibold text-muted-foreground">Waybill Tracking Number</label>
                                  <Input
                                    value={trackingNumber}
                                    onChange={(e) => setTrackingNumber(e.target.value)}
                                    placeholder="Insert shipping receipt ID"
                                    className="rounded-xl h-9 text-xs"
                                  />
                                </div>
                              </>
                            )}
                          </div>

                          <Button
                            onClick={() => handleDispatchOrder(selectedOrder.id)}
                            className="rounded-xl bg-primary text-primary-foreground hover:bg-primary/90 text-xs font-semibold h-9 w-full"
                          >
                            <Truck className="h-4 w-4 mr-1.5" /> Dispatch Shipment & Ship
                          </Button>
                        </div>
                      )}

                      {/* CASE D: SHIPPED - IN TRANSIT LOGISTICS */}
                      {selectedOrder.status === 'shipped' && selectedOrder.shipment && (
                        <div className="border border-border/80 bg-background p-5 rounded-xl space-y-4">

                          <div className="flex justify-between items-start">
                            <div className="space-y-1">
                              <h5 className="text-sm font-semibold text-foreground">Logistics Management</h5>
                              <p className="text-xs text-muted-foreground">
                                Update the active transit status or adjust waybill data.
                              </p>
                            </div>
                            <Button
                              variant="ghost"
                              size="sm"
                              onClick={() => setIsEditingShipment(prev => !prev)}
                              className="text-xs text-primary hover:bg-primary/5 rounded-lg flex items-center gap-1 h-7"
                            >
                              <Edit3 className="h-3.5 w-3.5" /> {isEditingShipment ? 'Cancel Edit' : 'Edit Waybill'}
                            </Button>
                          </div>

                          {isEditingShipment ? (
                            /* EDIT WAYBILL SUB-FORM */
                            <div className="space-y-3 p-4 border border-border/40 rounded-xl bg-zinc-50/20 dark:bg-slate-900/10">
                              <h6 className="text-xs font-bold text-foreground">Edit Shipping Waybill</h6>
                              <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
                                <div className="space-y-1">
                                  <label className="text-[10px] font-semibold text-muted-foreground uppercase">Courier</label>
                                  <Input value={editCourier} onChange={(e) => setEditCourier(e.target.value)} className="h-8 text-xs rounded-lg" />
                                </div>
                                <div className="space-y-1">
                                  <label className="text-[10px] font-semibold text-muted-foreground uppercase">Service</label>
                                  <Input value={editService} onChange={(e) => setEditService(e.target.value)} className="h-8 text-xs rounded-lg" />
                                </div>
                                <div className="space-y-1">
                                  <label className="text-[10px] font-semibold text-muted-foreground uppercase">Tracking No</label>
                                  <Input value={editTracking} onChange={(e) => setEditTracking(e.target.value)} className="h-8 text-xs rounded-lg" />
                                </div>
                              </div>
                              <Button
                                size="sm"
                                onClick={() => handleUpdateWaybill(selectedOrder.shipment!.id)}
                                className="h-8 text-xs rounded-lg w-full bg-primary text-primary-foreground hover:bg-primary/90 mt-2"
                              >
                                Save Waybill Changes
                              </Button>
                            </div>
                          ) : (
                            /* TRANSIT STATE SELECTOR */
                            <div className="space-y-3">
                              <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
                                <div className="space-y-1.5">
                                  <label className="text-xs font-semibold text-muted-foreground">New Transit Status</label>
                                  <Select value={transitStatus} onValueChange={setTransitStatus}>
                                    <SelectTrigger className="rounded-xl h-9 text-xs">
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
                                  <label className="text-xs font-semibold text-muted-foreground">Location (Optional)</label>
                                  <Input
                                    value={transitLocation}
                                    onChange={(e) => setTransitLocation(e.target.value)}
                                    placeholder="Warehouse, Hub City, etc."
                                    className="rounded-xl h-9 text-xs"
                                  />
                                </div>
                              </div>
                              <div className="space-y-1.5">
                                <label className="text-xs font-semibold text-muted-foreground">Description/Note (Optional)</label>
                                <Input
                                  value={transitDescription}
                                  onChange={(e) => setTransitDescription(e.target.value)}
                                  placeholder="E.g. Package hand over to courier JNE"
                                  className="rounded-xl h-9 text-xs"
                                />
                              </div>
                              <Button
                                onClick={() => handleUpdateShipmentStatus(selectedOrder.shipment!.id)}
                                className="rounded-xl bg-primary text-primary-foreground hover:bg-primary/90 text-xs font-semibold h-9 w-full"
                              >
                                Update Shipment Status
                              </Button>
                            </div>
                          )}

                        </div>
                      )}

                      {/* CASE E: DELIVERED (COMPLETED) */}
                      {selectedOrder.status === 'delivered' && (
                        <div className="border border-emerald-100 bg-emerald-50/10 p-4 rounded-xl flex items-center gap-3">
                          <CheckCircle2 className="h-5 w-5 text-emerald-600 dark:text-emerald-400 shrink-0" />
                          <div className="space-y-0.5">
                            <h5 className="text-sm font-semibold text-emerald-800 dark:text-emerald-400">Order Completed</h5>
                            <p className="text-xs text-muted-foreground">
                              This order has been successfully delivered and all fulfillment workflows are resolved.
                            </p>
                          </div>
                        </div>
                      )}

                      {/* CASE F: CANCELLED */}
                      {selectedOrder.status === 'cancelled' && (
                        <div className="border border-destructive/20 bg-destructive/5 p-4 rounded-xl flex items-center gap-3">
                          <XCircle className="h-5 w-5 text-destructive shrink-0" />
                          <div className="space-y-0.5">
                            <h5 className="text-sm font-semibold text-destructive">Order Cancelled</h5>
                            <p className="text-xs text-muted-foreground">
                              This order has been marked as cancelled. No further updates are permitted.
                            </p>
                          </div>
                        </div>
                      )}

                      {/* CASE G: AUTOMATED PAYMENT PENDING (MIDTRANS ETC) */}
                      {selectedOrder.status === 'pending' && selectedOrder.payment && selectedOrder.payment.provider !== 'manual' && (
                        <div className="border border-blue-200 dark:border-blue-950/40 bg-blue-50/10 p-4 rounded-xl flex items-start gap-3">
                          <CreditCard className="h-5 w-5 text-blue-600 dark:text-blue-400 mt-0.5" />
                          <div className="space-y-1">
                            <h5 className="text-sm font-semibold text-foreground">Awaiting Automated Settlement</h5>
                            <p className="text-xs text-muted-foreground">
                              This order is awaiting payment via {selectedOrder.payment.provider.toUpperCase()}. Statuses are synced automatically from payment webhooks.
                            </p>
                          </div>
                        </div>
                      )}

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
                            <th className="p-3 text-center">Quantity</th>
                            <th className="p-3 text-right">Price</th>
                            <th className="p-3 text-right">Subtotal</th>
                          </tr>
                        </thead>
                        <tbody className="divide-y divide-border/40">
                          {selectedOrder.items.map((item) => (
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
                              <td className="p-3 text-center font-medium">{item.quantity}</td>
                              <td className="p-3 text-right font-mono">{formatCurrency(item.unit_price)}</td>
                              <td className="p-3 text-right font-mono font-semibold">{formatCurrency(item.subtotal)}</td>
                            </tr>
                          ))}
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
                      <div className="p-4 rounded-xl border border-border/60 bg-zinc-50/20 dark:bg-slate-900/10 space-y-2.5 text-xs">
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
                        <User className="h-3.5 w-3.5" /> Shipping Address
                      </h5>
                      <div className="p-4 rounded-xl border border-border/60 bg-zinc-50/20 dark:bg-slate-900/10 space-y-3 text-xs">
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

                        {selectedOrder.shipment && (
                          <div className="border-t border-border/40 pt-2.5 mt-2 space-y-1 text-[11px] text-muted-foreground">
                            <div>Fulfillment Method: <span className="font-semibold text-foreground uppercase">{selectedOrder.shipment.fulfillment_method}</span></div>
                            {selectedOrder.shipment.tracking_number && (
                              <div>Waybill Number: <span className="font-semibold text-primary font-mono">{selectedOrder.shipment.tracking_number}</span></div>
                            )}
                            <div>Courier Assigned: <span className="font-semibold text-foreground uppercase">{selectedOrder.shipment.courier} {selectedOrder.shipment.service}</span></div>
                          </div>
                        )}
                      </div>
                    </div>

                  </div>

                  {/* Logistics Shipment History Timeline */}
                  {selectedOrder.shipment && (
                    <div className="space-y-4 pt-2 border-t border-border/60">
                      <h5 className="text-sm font-bold text-foreground font-display flex items-center gap-1.5">
                        <Activity className="h-4 w-4 text-muted-foreground" /> Shipment Tracking History
                      </h5>

                      {selectedOrder.shipment.events && selectedOrder.shipment.events.length > 0 ? (
                        /* EVENTS DOT TIMELINE */
                        <div className="relative pl-6 border-l-2 border-primary/20 space-y-6 ml-2 py-1">
                          {selectedOrder.shipment.events.map((event, idx) => {
                            const isLatest = idx === 0;
                            return (
                              <div key={event.id} className="relative">
                                {/* Dotted marker circle */}
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
                      ) : (
                        /* EMPTY HISTORICAL EVENTS */
                        <div className="p-4 border border-dashed border-border/80 rounded-xl bg-zinc-50/10 text-center text-xs text-muted-foreground flex flex-col items-center gap-1">
                          <Activity className="h-6 w-6 text-muted-foreground/50 stroke-[1.5]" />
                          <span>No shipment updates logged yet.</span>
                          <span className="text-[10px] text-muted-foreground/60">Use the control panel above to post the first shipment status update.</span>
                        </div>
                      )}
                    </div>
                  )}

                </div>

                </div>

              </div>
            )}

          </div>

        </div>

      </div>
    </div>
  );
}
