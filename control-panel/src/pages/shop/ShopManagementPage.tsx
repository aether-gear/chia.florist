import { useState, useMemo, useEffect } from 'react';
import { useParams, useSearchParams, useNavigate } from 'react-router-dom';
import {
  MapPin,
  Truck,
  Package,
  Loader2,
  Plus,
  Info,
  RefreshCw,
  MoreHorizontal,
  Edit,
  Trash2,
  AlertTriangle,
  ArrowRight,
  ShieldCheck,
  ShoppingBag,
  CheckCircle2,
  XCircle,
  AlertCircle,
} from 'lucide-react';
import { Button } from '../../components/ui/button';
import { Badge } from '../../components/ui/badge';
import { Input } from '../../components/ui/input';
import { Label } from '../../components/ui/label';
import { Checkbox } from '../../components/ui/checkbox';
import { Skeleton } from '../../components/ui/skeleton';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '../../components/ui/tabs';
import ShopOrdersTab from '../../components/shops/ShopOrdersTab';
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetClose,
} from '../../components/ui/sheet';
import {
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
} from '../../components/ui/dropdown-menu';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '../../components/ui/dialog';
import { useShopViewModel } from '../../viewmodels/useShopViewModel';
import { useAuthMeViewModel } from '../../viewmodels/useAuthMeViewModel';
import Pagination from '../../components/Pagination';
import SearchInput from '../../components/SearchInput';
import { DataCard, DataCardList, DataCardGridHeader } from '../../components/DataCard';
import InventoryFormSheet from '../../components/shops/InventoryFormSheet';
import AddressFormSheet from '../../components/shops/AddressFormSheet';
import CourierFormSheet from '../../components/shops/CourierFormSheet';
import Breadcrumb from '../../components/Breadcrumb';

export default function ShopManagementPage() {
  const { shopId: urlShopId } = useParams<{ shopId?: string }>();
  const [searchParams, setSearchParams] = useSearchParams();
  const navigate = useNavigate();

  const activeTab = searchParams.get('tab') || 'info';

  const { isAdmin } = useAuthMeViewModel();
  const {
    shops,
    total,
    page,
    limit,
    selectedShopId,
    selectedShopInfo,
    addresses,
    couriers,
    products,
    loading,
    detailsLoading,
    error,
    detailsError,
    setPage,
    deleteAddress,
    updateCourier,
    verifyCourier,
    saveShop,
    createShop,
    deleteShop,
    removeInventory,
    selectShop,
    loadShopById,
    refresh,
  } = useShopViewModel(urlShopId);

  // Courier states
  const [courierSearchQuery, setCourierSearchQuery] = useState('');
  const [courierStatusFilter, setCourierStatusFilter] = useState<'active' | 'all' | 'pending' | 'inactive'>('active');
  const [editingCourier, setEditingCourier] = useState<any | null>(null);
  const [isCourierSheetOpen, setIsCourierSheetOpen] = useState(false);

  // Admin Verification Dialog states
  const [verifyingCourier, setVerifyingCourier] = useState<any | null>(null);
  const [verifyAction, setVerifyAction] = useState<'verify' | 'reject'>('verify');
  const [verifyRejectionReason, setVerifyRejectionReason] = useState('');
  const [isProcessingVerification, setIsProcessingVerification] = useState(false);
  const [verificationError, setVerificationError] = useState<string | null>(null);

  const handleOpenEditCourier = (courier: any) => {
    setEditingCourier(courier);
    setIsCourierSheetOpen(true);
  };

  const handleOpenVerifyModal = (courier: any, action: 'verify' | 'reject') => {
    setVerifyingCourier(courier);
    setVerifyAction(action);
    setVerifyRejectionReason('');
    setVerificationError(null);
  };

  const handleConfirmVerification = async () => {
    if (!selectedShopId || !verifyingCourier) return;
    setIsProcessingVerification(true);
    setVerificationError(null);
    try {
      await verifyCourier(
        selectedShopId,
        verifyingCourier.code,
        verifyAction,
        verifyAction === 'reject' ? verifyRejectionReason.trim() : undefined
      );
      setVerifyingCourier(null);
      if (selectedShopInfo) {
        selectShop(selectedShopInfo);
      }
    } catch (err: any) {
      setVerificationError(err.message || `Failed to ${verifyAction} courier`);
    } finally {
      setIsProcessingVerification(false);
    }
  };

  const filteredCouriers = useMemo(() => {
    if (!couriers) return [];
    let result = couriers;

    // Filter by display filter (by default: 'active')
    if (courierStatusFilter === 'active') {
      result = result.filter((c) => c.active);
    } else if (courierStatusFilter === 'pending') {
      result = result.filter((c) => c.verification_status === 'pending');
    } else if (courierStatusFilter === 'inactive') {
      result = result.filter((c) => !c.active);
    }
    // If 'all', keep all couriers

    if (courierSearchQuery.trim()) {
      const q = courierSearchQuery.toLowerCase().trim();
      result = result.filter(
        (c) =>
          (c.code && c.code.toLowerCase().includes(q)) ||
          (c.branch_name && c.branch_name.toLowerCase().includes(q)) ||
          (c.name && c.name.toLowerCase().includes(q)) ||
          (c.location_address && c.location_address.toLowerCase().includes(q))
      );
    }
    return result;
  }, [couriers, courierStatusFilter, courierSearchQuery]);

  const getCourierVerificationBadge = (status?: string) => {
    switch (status) {
      case 'verified':
        return (
          <Badge
            variant="default"
            className="bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border-0 rounded-lg scale-90"
          >
            <CheckCircle2 className="h-3 w-3 mr-1" />
            Verified
          </Badge>
        );
      case 'pending':
        return (
          <Badge
            variant="secondary"
            className="bg-amber-500/10 text-amber-600 dark:text-amber-400 border-0 rounded-lg scale-90 animate-pulse"
          >
            <AlertCircle className="h-3 w-3 mr-1" />
            Pending Review
          </Badge>
        );
      case 'rejected':
        return (
          <Badge
            variant="destructive"
            className="bg-rose-500/10 text-rose-600 dark:text-rose-400 border-0 rounded-lg scale-90"
          >
            <XCircle className="h-3 w-3 mr-1" />
            Rejected
          </Badge>
        );
      default:
        return (
          <Badge
            variant="outline"
            className="bg-muted text-muted-foreground border-0 rounded-lg scale-90"
          >
            Unconfigured
          </Badge>
        );
    }
  };

  // Sync state with URL parameter on mount or URL change
  useEffect(() => {
    if (urlShopId) {
      loadShopById(urlShopId);
    } else {
      selectShop(null);
    }
  }, [urlShopId, loadShopById, selectShop]);

  const isDetailView = Boolean(urlShopId || selectedShopId);
  const currentShopId = urlShopId || selectedShopId || '';

  const [searchQuery, setSearchQuery] = useState('');

  const filteredShops = useMemo(() => {
    if (!shops) return [];
    return shops.filter(
      (shop) =>
        (shop.name && shop.name.toLowerCase().includes(searchQuery.toLowerCase())) ||
        (shop.description && shop.description.toLowerCase().includes(searchQuery.toLowerCase()))
    );
  }, [shops, searchQuery]);

  // Address states
  const [isAddressOpen, setIsAddressOpen] = useState(false);
  const [editingAddress, setEditingAddress] = useState<any | null>(null);
  const [addressToDelete, setAddressToDelete] = useState<any | null>(null);
  const [isDeletingAddress, setIsDeletingAddress] = useState(false);
  const [deleteAddressError, setDeleteAddressError] = useState<string | null>(null);

  const handleDeleteAddressConfirm = async () => {
    if (!selectedShopId || !addressToDelete) return;

    setIsDeletingAddress(true);
    setDeleteAddressError(null);

    try {
      await deleteAddress(selectedShopId, addressToDelete.id);
      setAddressToDelete(null);
    } catch (err: any) {
      setDeleteAddressError(err.message || 'Failed to delete address');
    } finally {
      setIsDeletingAddress(false);
    }
  };

  // Inventory states
  const [isInventoryOpen, setIsInventoryOpen] = useState(false);
  const [editingInventoryProduct, setEditingInventoryProduct] = useState<any | null>(null);

  // Delete Inventory states
  const [inventoryToDelete, setInventoryToDelete] = useState<any | null>(null);
  const [isDeletingInventory, setIsDeletingInventory] = useState(false);
  const [deleteInventoryError, setDeleteInventoryError] = useState<string | null>(null);

  const handleDeleteInventoryConfirm = async () => {
    if (!selectedShopId || !inventoryToDelete) return;

    setIsDeletingInventory(true);
    setDeleteInventoryError(null);

    try {
      await removeInventory(selectedShopId, inventoryToDelete.id);
      setInventoryToDelete(null);
    } catch (err: any) {
      setDeleteInventoryError(err.message || 'Failed to delete inventory');
    } finally {
      setIsDeletingInventory(false);
    }
  };

  // Add Shop states
  const [isAddShopOpen, setIsAddShopOpen] = useState(false);
  const [newShopName, setNewShopName] = useState('');
  const [newShopDesc, setNewShopDesc] = useState('');
  const [newShopIsActive, setNewShopIsActive] = useState(false);
  const [newShopApprovalStatus, setNewShopApprovalStatus] = useState<string>('pending');
  const [isCreatingShop, setIsCreatingShop] = useState(false);
  const [addShopError, setAddShopError] = useState<string | null>(null);

  const handleAddShopSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newShopName) return;
    setIsCreatingShop(true);
    setAddShopError(null);
    try {
      const payload: any = {
        name: newShopName,
        description: newShopDesc || undefined,
      };
      if (isAdmin) {
        payload.is_active = newShopIsActive ? 'true' : 'false';
        payload.approval_status = newShopApprovalStatus;
      }
      const success = await createShop(payload);
      if (success) {
        setIsAddShopOpen(false);
        setNewShopName('');
        setNewShopDesc('');
        setNewShopIsActive(false);
        setNewShopApprovalStatus('pending');
      } else {
        setAddShopError('Failed to add shop. Please check your inputs or try again.');
      }
    } catch (err: any) {
      setAddShopError(err.message || 'Failed to add shop');
    } finally {
      setIsCreatingShop(false);
    }
  };

  // Edit Shop states
  const [isEditShopOpen, setIsEditShopOpen] = useState(false);
  const [editShopName, setEditShopName] = useState('');
  const [editShopDesc, setEditShopDesc] = useState('');
  const [editShopIsActive, setEditShopIsActive] = useState(false);
  const [editShopApprovalStatus, setEditShopApprovalStatus] = useState<string>('pending');
  const [isSavingShop, setIsSavingShop] = useState(false);
  const [editShopError, setEditShopError] = useState<string | null>(null);

  const handleOpenEditShop = () => {
    if (!selectedShopInfo) return;
    setEditShopName(selectedShopInfo.name || '');
    setEditShopDesc(selectedShopInfo.description || '');
    setEditShopIsActive(selectedShopInfo.is_active || false);
    setEditShopApprovalStatus(selectedShopInfo.approval_status || 'pending');
    setEditShopError(null);
    setIsEditShopOpen(true);
  };

  const handleEditShopSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editShopName) return;
    setIsSavingShop(true);
    setEditShopError(null);
    try {
      const payload: any = {
        name: editShopName,
        description: editShopDesc || undefined,
      };
      if (isAdmin) {
        payload.is_active = editShopIsActive ? 'true' : 'false';
        payload.approval_status = editShopApprovalStatus;
      }
      const success = await saveShop(payload);
      if (success) {
        setIsEditShopOpen(false);
      } else {
        setEditShopError('Failed to update shop. Please check your inputs.');
      }
    } catch (err: any) {
      setEditShopError(err.message || 'Failed to update shop');
    } finally {
      setIsSavingShop(false);
    }
  };

  // Delete Shop states
  const [isDeleteShopOpen, setIsDeleteShopOpen] = useState(false);
  const [isDeletingShop, setIsDeletingShop] = useState(false);
  const [deleteShopError, setDeleteShopError] = useState<string | null>(null);

  const handleDeleteShopConfirm = async () => {
    if (!selectedShopId) return;
    setIsDeletingShop(true);
    setDeleteShopError(null);
    try {
      await deleteShop(selectedShopId);
      setIsDeleteShopOpen(false);
      navigate('/shop');
    } catch (err: any) {
      setDeleteShopError(err.message || 'Failed to delete shop');
    } finally {
      setIsDeletingShop(false);
    }
  };

  const handleOpenDetails = (shop: any) => {
    navigate(`/shop/${shop.id}?tab=${searchParams.get('tab') || 'info'}`);
  };

  const handleBackToList = () => {
    navigate('/shop');
  };

  const handleTabChange = (newTab: string) => {
    setSearchParams({ tab: newTab }, { replace: true });
  };

  const getApprovalBadge = (status?: string) => {
    switch (status) {
      case 'approved':
        return (
          <Badge
            variant="default"
            className="bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border-0 rounded-lg scale-90"
          >
            Approved
          </Badge>
        );
      case 'rejected':
        return (
          <Badge
            variant="destructive"
            className="bg-rose-500/10 text-rose-600 dark:text-rose-400 border-0 rounded-lg scale-90"
          >
            Rejected
          </Badge>
        );
      case 'pending':
      default:
        return (
          <Badge
            variant="secondary"
            className="bg-amber-500/10 text-amber-600 dark:text-amber-400 border-0 rounded-lg scale-90"
          >
            Pending Approval
          </Badge>
        );
    }
  };

  return (
    <div className="flex-col md:flex">
      <div className="flex-1 space-y-8 p-6 sm:p-8 lg:p-12">
        {!isDetailView ? (
          /* View 1: Store Locations List View */
          <div className="space-y-6 animate-in fade-in slide-in-from-left-4 duration-300">
            <div className="flex items-center justify-between space-y-2">
              <div>
                <h2 className="text-3xl font-bold font-display tracking-tight text-foreground">Shop Management</h2>
                <p className="text-muted-foreground text-sm">
                  View and manage all your store branches and their details.
                </p>
              </div>
            </div>

            <div className="space-y-6">
              <div className="pb-4 border-b border-border/60">
                <h3 className="text-xl font-bold font-display tracking-tight text-foreground">Store Locations</h3>
                <p className="text-muted-foreground text-sm">You have {total} shop locations registered.</p>
              </div>

              <div className="flex flex-col sm:flex-row items-center justify-between gap-4 w-full">
                <SearchInput
                  value={searchQuery}
                  onChange={setSearchQuery}
                  placeholder="Search shops by name..."
                  className="relative flex-1 max-w-sm w-full"
                />
                <div className="flex items-center gap-2 justify-end w-full sm:w-auto">
                  <Button
                    variant="outline"
                    onClick={() => refresh()}
                    disabled={loading}
                    className="flex items-center gap-1.5 border-border text-foreground hover:text-primary hover:bg-primary/5 rounded-xl transition-colors"
                  >
                    <RefreshCw className={`h-4 w-4 ${loading ? 'animate-spin' : ''}`} />
                    Refresh
                  </Button>
                  <Button
                    onClick={() => setIsAddShopOpen(true)}
                    className="flex items-center gap-1.5 bg-primary hover:bg-primary/90 text-primary-foreground rounded-xl"
                  >
                    <Plus className="h-4 w-4" />
                    Add Shop
                  </Button>
                </div>
              </div>

              <div className="space-y-4">
                <div className="flex items-center justify-between px-1 text-xs text-muted-foreground">
                  <span>Found {filteredShops.length} shops</span>
                </div>

                <DataCardList>
                  {loading ? (
                    Array.from({ length: 4 }).map((_, i) => (
                      <DataCard key={`skeleton-${i}`}>
                        <div className="col-span-12">
                          <Skeleton className="h-5 w-40 bg-muted animate-pulse" />
                        </div>
                      </DataCard>
                    ))
                  ) : error ? (
                    <div className="py-12 border-0 bg-transparent text-destructive text-center">
                      Failed to load shops: {error}
                    </div>
                  ) : filteredShops.length === 0 ? (
                    <div className="py-12 border border-dashed border-border/80 rounded-2xl bg-zinc-50/10 text-center text-muted-foreground">
                      <MapPin className="h-8 w-8 text-slate-400 mb-2 mx-auto" />
                      <p>No shops found</p>
                      <p className="text-sm">No store locations match your search.</p>
                    </div>
                  ) : (
                    filteredShops.map((shop) => (
                      <DataCard
                        key={shop.id}
                        onClick={() => handleOpenDetails(shop)}
                        className="cursor-pointer hover:border-primary/40 transition-all duration-200"
                      >
                        <div className="col-span-1 md:col-span-4 min-w-0">
                          <h4 className="font-semibold font-display text-sm text-foreground truncate group-hover:text-primary">
                            {shop.name}
                          </h4>
                          <p className="text-[10px] text-muted-foreground font-mono truncate">ID: {shop.id}</p>
                        </div>

                        <div className="col-span-1 md:col-span-4 text-xs text-muted-foreground truncate">
                          {shop.description || 'No description provided.'}
                        </div>

                        <div className="col-span-1 md:col-span-4 flex items-center justify-end gap-2 flex-wrap">
                          {getApprovalBadge(shop.approval_status)}
                          <Badge
                            variant={shop.is_active ? 'default' : 'secondary'}
                            className={
                              shop.is_active
                                ? 'bg-primary/10 text-primary border-0 rounded-lg scale-90'
                                : 'bg-muted text-muted-foreground border-0 rounded-lg scale-90'
                            }
                          >
                            {shop.is_active ? 'Active' : 'Inactive'}
                          </Badge>

                          <Button
                            variant="ghost"
                            size="sm"
                            className="h-8 px-2.5 text-xs text-muted-foreground hover:text-primary hover:bg-primary/10 rounded-lg transition-colors"
                          >
                            Manage <ArrowRight className="ml-1 h-3.5 w-3.5" />
                          </Button>
                        </div>
                      </DataCard>
                    ))
                  )}
                </DataCardList>

                <Pagination
                  currentPage={page}
                  totalPages={Math.ceil(total / limit)}
                  totalItems={total}
                  limit={limit}
                  onPageChange={setPage}
                  itemNamePlural="shops"
                />
              </div>
            </div>
          </div>
        ) : (
          /* View 2: Selected Shop Details View */
          <div className="space-y-6 animate-in fade-in slide-in-from-right-4 duration-300">
            {/* Global Breadcrumb Navigation */}
            <Breadcrumb
              items={[
                { label: 'Shop Management', onClick: handleBackToList },
                { label: 'Store Locations', onClick: handleBackToList },
                { label: selectedShopInfo?.name || 'Shop Details' },
              ]}
            />

            {/* Header with Back button and Shop Info */}
            <div className="flex flex-col sm:flex-row sm:items-center justify-between pb-4 border-b border-border/60 gap-4">
              <div className="space-y-1">
                <div className="flex items-center gap-3 flex-wrap">
                  <h2 className="text-2xl sm:text-3xl font-bold font-display tracking-tight text-foreground">
                    {selectedShopInfo?.name || 'Shop Details'}
                  </h2>
                  {selectedShopInfo && (
                    <div className="flex items-center gap-2">
                      {getApprovalBadge(selectedShopInfo.approval_status)}
                      <Badge
                        variant={selectedShopInfo.is_active ? 'default' : 'secondary'}
                        className={
                          selectedShopInfo.is_active
                            ? 'bg-primary/10 text-primary border-0 rounded-lg'
                            : 'bg-muted text-muted-foreground border-0 rounded-lg'
                        }
                      >
                        {selectedShopInfo.is_active ? 'Active' : 'Inactive'}
                      </Badge>
                    </div>
                  )}
                </div>
                <p className="text-muted-foreground text-sm pt-1">
                  {selectedShopInfo?.description || 'Manage parameters, addresses, couriers, and product inventory.'}
                </p>
              </div>

              <div className="flex items-center gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={handleOpenEditShop}
                  className="rounded-xl border-border text-foreground hover:bg-muted"
                >
                  <Edit className="h-4 w-4 mr-1.5" />
                  Edit Shop
                </Button>

                {isAdmin && (
                  <Button
                    variant="destructive"
                    size="sm"
                    onClick={() => {
                      setDeleteShopError(null);
                      setIsDeleteShopOpen(true);
                    }}
                    className="rounded-xl"
                  >
                    <Trash2 className="h-4 w-4 mr-1.5" />
                    Delete Shop
                  </Button>
                )}
              </div>
            </div>

            {detailsLoading && (
              <div className="flex justify-center py-12">
                <Loader2 className="h-6 w-6 animate-spin text-primary" />
              </div>
            )}

            {detailsError && (
              <div className="p-3 text-sm text-destructive bg-destructive/10 rounded-xl border border-destructive/20 font-sans">
                {detailsError}
              </div>
            )}

            {!detailsLoading && selectedShopInfo && (
              <Tabs value={activeTab} onValueChange={handleTabChange} className="space-y-6 pt-2">
                <TabsList className="grid grid-cols-2 sm:grid-cols-5 max-w-2xl bg-muted/50 p-1 rounded-xl border border-border/60">
                  <TabsTrigger
                    value="info"
                    className="flex items-center gap-1.5 data-[state=active]:bg-background data-[state=active]:shadow-sm rounded-lg text-xs font-medium"
                  >
                    <Info className="h-3.5 w-3.5" /> General
                  </TabsTrigger>
                  <TabsTrigger
                    value="orders"
                    className="flex items-center gap-1.5 data-[state=active]:bg-background data-[state=active]:shadow-sm rounded-lg text-xs font-medium"
                  >
                    <ShoppingBag className="h-3.5 w-3.5" /> Orders
                  </TabsTrigger>
                  <TabsTrigger
                    value="products"
                    className="flex items-center gap-1.5 data-[state=active]:bg-background data-[state=active]:shadow-sm rounded-lg text-xs font-medium"
                  >
                    <Package className="h-3.5 w-3.5" /> Products
                  </TabsTrigger>
                  <TabsTrigger
                    value="addresses"
                    className="flex items-center gap-1.5 data-[state=active]:bg-background data-[state=active]:shadow-sm rounded-lg text-xs font-medium"
                  >
                    <MapPin className="h-3.5 w-3.5" /> Addresses
                  </TabsTrigger>
                  <TabsTrigger
                    value="couriers"
                    className="flex items-center gap-1.5 data-[state=active]:bg-background data-[state=active]:shadow-sm rounded-lg text-xs font-medium"
                  >
                    <Truck className="h-3.5 w-3.5" /> Couriers
                  </TabsTrigger>
                </TabsList>

                {/* General Tab */}
                <TabsContent value="info" className="space-y-4 pt-2">
                  <div className="space-y-6">
                    <div className="pb-4 border-b border-border/60">
                      <h3 className="text-lg font-bold font-display text-foreground">Shop Overview</h3>
                      <p className="text-muted-foreground text-sm">General details for this store location.</p>
                    </div>

                    <div className="grid gap-4 max-w-lg">
                      <div className="p-4 rounded-xl border border-border/60 bg-muted/20 space-y-4">
                        <div className="grid grid-cols-3 gap-2 text-sm">
                          <span className="text-muted-foreground font-medium">Shop Name</span>
                          <span className="col-span-2 font-semibold text-foreground">
                            {selectedShopInfo?.name || '-'}
                          </span>
                        </div>

                        <div className="grid grid-cols-3 gap-2 text-sm">
                          <span className="text-muted-foreground font-medium">Shop ID</span>
                          <span className="col-span-2 font-mono text-xs text-muted-foreground">
                            {selectedShopInfo?.id || '-'}
                          </span>
                        </div>

                        <div className="grid grid-cols-3 gap-2 text-sm">
                          <span className="text-muted-foreground font-medium">Approval Status</span>
                          <div className="col-span-2">
                            {getApprovalBadge(selectedShopInfo?.approval_status)}
                          </div>
                        </div>

                        <div className="grid grid-cols-3 gap-2 text-sm">
                          <span className="text-muted-foreground font-medium">Operating Status</span>
                          <div className="col-span-2">
                            <Badge
                              variant={selectedShopInfo?.is_active ? 'default' : 'secondary'}
                              className={
                                selectedShopInfo?.is_active
                                  ? 'bg-primary/10 text-primary border-0 rounded-lg'
                                  : 'bg-muted text-muted-foreground border-0 rounded-lg'
                              }
                            >
                              {selectedShopInfo?.is_active ? 'Active' : 'Inactive'}
                            </Badge>
                          </div>
                        </div>

                        <div className="grid grid-cols-3 gap-2 text-sm">
                          <span className="text-muted-foreground font-medium">Description</span>
                          <span className="col-span-2 text-muted-foreground">
                            {selectedShopInfo?.description || 'No description provided.'}
                          </span>
                        </div>
                      </div>
                    </div>
                  </div>
                </TabsContent>

                {/* Products Tab */}
                <TabsContent value="products" className="space-y-4 pt-2">
                  <div className="space-y-6">
                    <div className="flex flex-row items-center justify-between pb-4 border-b border-border/60">
                      <div>
                        <h3 className="text-lg font-bold font-display text-foreground">Product Inventory</h3>
                        <p className="text-muted-foreground text-sm">
                          Listed products and current stock levels at this location.
                        </p>
                      </div>

                      <Button
                        size="sm"
                        onClick={() => setIsInventoryOpen(true)}
                        className="bg-primary hover:bg-primary/90 text-primary-foreground rounded-xl"
                      >
                        <Plus className="mr-1.5 h-4 w-4" /> Add Inventory
                      </Button>
                    </div>

                    <div className="space-y-2">
                      <DataCardList>
                        {products.length === 0 ? (
                          <div className="py-8 border border-dashed border-border/80 rounded-2xl bg-zinc-50/10 text-center text-muted-foreground text-sm">
                            No products found for this shop.
                          </div>
                        ) : (
                          products.map((product) => (
                            <DataCard key={product.id} className="bg-muted/20">
                              <div className="col-span-1 md:col-span-4 min-w-0 space-y-0.5">
                                <h4 className="font-semibold text-sm text-foreground line-clamp-2 leading-tight">
                                  {product.name}
                                </h4>
                                <p className="text-xs text-muted-foreground font-mono truncate">{product.slug}</p>
                              </div>

                              <div className="col-span-1 md:col-span-3 min-w-0 text-xs font-mono text-muted-foreground truncate">
                                <span className="md:hidden font-sans mr-1 text-foreground font-medium">SKU:</span>
                                {product.sku || '-'}
                              </div>

                              <div className="col-span-1 md:col-span-3 min-w-0 space-y-0.5 text-xs">
                                <div className="font-bold text-primary truncate">
                                  {new Intl.NumberFormat('id-ID', {
                                    style: 'currency',
                                    currency: 'IDR',
                                    minimumFractionDigits: 0,
                                  }).format(product.price)}
                                </div>
                                <div className="text-[11px] text-muted-foreground truncate">
                                  <span className="font-semibold text-foreground">
                                    {product.inventory.available}
                                  </span>{' '}
                                  available{' '}
                                  <span className="opacity-80">(Total: {product.inventory.total_stock})</span>
                                </div>
                              </div>

                              <div className="col-span-1 md:col-span-2 min-w-0 flex items-center justify-end">
                                <DropdownMenu>
                                  <DropdownMenuTrigger asChild>
                                    <Button variant="ghost" size="icon" className="h-8 w-8 rounded-xl">
                                      <MoreHorizontal className="h-4 w-4" />
                                      <span className="sr-only">Open menu</span>
                                    </Button>
                                  </DropdownMenuTrigger>
                                  <DropdownMenuContent align="end" className="w-44 rounded-xl">
                                    <DropdownMenuItem
                                      onClick={() => {
                                        setEditingInventoryProduct(product);
                                      }}
                                    >
                                      <Edit className="h-4 w-4 mr-2 text-muted-foreground" />
                                      Edit Inventory
                                    </DropdownMenuItem>
                                    <DropdownMenuSeparator className="my-1" />
                                    <DropdownMenuItem
                                      onClick={() => {
                                        setInventoryToDelete(product);
                                        setDeleteInventoryError(null);
                                      }}
                                      className="text-destructive focus:text-destructive"
                                    >
                                      <Trash2 className="h-4 w-4 mr-2 text-destructive" />
                                      Delete Inventory
                                    </DropdownMenuItem>
                                  </DropdownMenuContent>
                                </DropdownMenu>
                              </div>
                            </DataCard>
                          ))
                        )}
                      </DataCardList>
                    </div>
                  </div>
                </TabsContent>

                {/* Addresses Tab */}
                <TabsContent value="addresses" className="space-y-4 pt-2">
                  <div className="space-y-6">
                    <div className="flex flex-row items-center justify-between pb-4 border-b border-border/60">
                      <div>
                        <h3 className="text-lg font-bold font-display text-foreground">Shop Addresses</h3>
                        <p className="text-muted-foreground text-sm">
                          Physical locations where this branch operates.
                        </p>
                      </div>

                      <Button
                        size="sm"
                        onClick={() => setIsAddressOpen(true)}
                        className="bg-primary hover:bg-primary/90 text-primary-foreground rounded-xl"
                      >
                        <Plus className="mr-1.5 h-4 w-4" /> Add Address
                      </Button>
                    </div>

                    <div className="space-y-2">
                      <DataCardList>
                        {addresses.length === 0 ? (
                          <div className="py-8 border border-dashed border-border/80 rounded-2xl bg-zinc-50/10 text-center text-muted-foreground text-sm">
                            No addresses found.
                          </div>
                        ) : (
                          addresses.map((addr) => (
                            <DataCard key={addr.id} className="bg-muted/20">
                              <div className="col-span-1 md:col-span-3 min-w-0 font-semibold text-foreground text-sm truncate">
                                {addr.label}
                              </div>
                              <div className="col-span-1 md:col-span-4 min-w-0 text-xs text-muted-foreground truncate">
                                {addr.full_address}
                              </div>
                              <div className="col-span-1 md:col-span-2 min-w-0 text-xs text-muted-foreground truncate">
                                {addr.phone || '-'}
                              </div>
                              <div className="col-span-1 md:col-span-2 min-w-0 text-right">
                                <Badge
                                  variant={addr.is_active ? 'default' : 'secondary'}
                                  className={
                                    addr.is_active
                                      ? 'bg-primary/10 text-primary border-0 scale-90 origin-right'
                                      : 'bg-muted text-muted-foreground border-0 scale-90 origin-right'
                                  }
                                >
                                  {addr.is_active ? 'Active' : 'Inactive'}
                                </Badge>
                              </div>
                              <div className="col-span-1 md:col-span-1 min-w-0 flex items-center justify-end">
                                <DropdownMenu>
                                  <DropdownMenuTrigger asChild>
                                    <Button variant="ghost" size="icon" className="h-8 w-8 rounded-xl">
                                      <MoreHorizontal className="h-4 w-4" />
                                      <span className="sr-only">Open menu</span>
                                    </Button>
                                  </DropdownMenuTrigger>
                                  <DropdownMenuContent align="end" className="w-44 rounded-xl">
                                    <DropdownMenuItem
                                      onClick={() => {
                                        setEditingAddress(addr);
                                      }}
                                    >
                                      <Edit className="h-4 w-4 mr-2 text-muted-foreground" />
                                      Edit Address
                                    </DropdownMenuItem>
                                    <DropdownMenuSeparator className="my-1" />
                                    <DropdownMenuItem
                                      onClick={() => {
                                        setAddressToDelete(addr);
                                        setDeleteAddressError(null);
                                      }}
                                      className="text-destructive focus:text-destructive"
                                    >
                                      <Trash2 className="h-4 w-4 mr-2 text-destructive" />
                                      Delete Address
                                    </DropdownMenuItem>
                                  </DropdownMenuContent>
                                </DropdownMenu>
                              </div>
                            </DataCard>
                          ))
                        )}
                      </DataCardList>
                    </div>
                  </div>
                </TabsContent>

                {/* Couriers Tab */}
                <TabsContent value="couriers" className="space-y-4 pt-2">
                  <div className="space-y-6">
                    <div className="flex flex-col sm:flex-row sm:items-center justify-between pb-4 border-b border-border/60 gap-4">
                      <div>
                        <h3 className="text-lg font-bold font-display text-foreground">Shop Courier Management</h3>
                        <p className="text-muted-foreground text-sm">
                          Shipping providers configured for this branch (bound to Komerce logistics).
                        </p>
                      </div>

                      <div className="flex items-center gap-2 flex-wrap">
                        <button
                          type="button"
                          onClick={() => setCourierStatusFilter('all')}
                          className="focus:outline-none"
                        >
                          <Badge
                            variant={courierStatusFilter === 'all' ? 'default' : 'outline'}
                            className="px-2.5 py-1 text-xs rounded-lg cursor-pointer hover:bg-muted/80 transition-colors"
                          >
                            Total: {couriers.length}
                          </Badge>
                        </button>
                        <button
                          type="button"
                          onClick={() => setCourierStatusFilter('active')}
                          className="focus:outline-none"
                        >
                          <Badge
                            variant="default"
                            className={`px-2.5 py-1 text-xs rounded-lg cursor-pointer transition-colors ${
                              courierStatusFilter === 'active'
                                ? 'bg-emerald-600 text-white ring-2 ring-emerald-400/30'
                                : 'bg-emerald-500/10 text-emerald-600 hover:bg-emerald-500/20 border-0'
                            }`}
                          >
                            Active: {couriers.filter((c) => c.active).length}
                          </Badge>
                        </button>
                        <button
                          type="button"
                          onClick={() => setCourierStatusFilter('pending')}
                          className="focus:outline-none"
                        >
                          <Badge
                            variant="secondary"
                            className={`px-2.5 py-1 text-xs rounded-lg cursor-pointer transition-colors ${
                              courierStatusFilter === 'pending'
                                ? 'bg-amber-600 text-white ring-2 ring-amber-400/30'
                                : 'bg-amber-500/10 text-amber-600 hover:bg-amber-500/20 border-0'
                            }`}
                          >
                            Pending: {couriers.filter((c) => c.verification_status === 'pending').length}
                          </Badge>
                        </button>
                      </div>
                    </div>

                    <div className="flex flex-col sm:flex-row items-center justify-between gap-4 w-full">
                      <SearchInput
                        value={courierSearchQuery}
                        onChange={setCourierSearchQuery}
                        placeholder="Filter couriers by brand, name, or address..."
                        className="relative flex-1 max-w-sm w-full"
                      />

                      {/* Display Filter controls (Default: Active) */}
                      <div className="flex items-center gap-1 bg-muted/60 p-1 rounded-xl border border-border/60 self-stretch sm:self-auto overflow-x-auto">
                        <Button
                          type="button"
                          variant="ghost"
                          size="sm"
                          onClick={() => setCourierStatusFilter('active')}
                          className={`h-7 px-3 text-xs font-medium rounded-lg transition-all ${
                            courierStatusFilter === 'active'
                              ? 'bg-background text-foreground shadow-sm font-semibold'
                              : 'text-muted-foreground hover:text-foreground'
                          }`}
                        >
                          Active ({couriers.filter((c) => c.active).length})
                        </Button>
                        <Button
                          type="button"
                          variant="ghost"
                          size="sm"
                          onClick={() => setCourierStatusFilter('all')}
                          className={`h-7 px-3 text-xs font-medium rounded-lg transition-all ${
                            courierStatusFilter === 'all'
                              ? 'bg-background text-foreground shadow-sm font-semibold'
                              : 'text-muted-foreground hover:text-foreground'
                          }`}
                        >
                          All ({couriers.length})
                        </Button>
                        <Button
                          type="button"
                          variant="ghost"
                          size="sm"
                          onClick={() => setCourierStatusFilter('pending')}
                          className={`h-7 px-3 text-xs font-medium rounded-lg transition-all ${
                            courierStatusFilter === 'pending'
                              ? 'bg-background text-foreground shadow-sm font-semibold'
                              : 'text-muted-foreground hover:text-foreground'
                          }`}
                        >
                          Pending ({couriers.filter((c) => c.verification_status === 'pending').length})
                        </Button>
                        <Button
                          type="button"
                          variant="ghost"
                          size="sm"
                          onClick={() => setCourierStatusFilter('inactive')}
                          className={`h-7 px-3 text-xs font-medium rounded-lg transition-all ${
                            courierStatusFilter === 'inactive'
                              ? 'bg-background text-foreground shadow-sm font-semibold'
                              : 'text-muted-foreground hover:text-foreground'
                          }`}
                        >
                          Inactive ({couriers.filter((c) => !c.active).length})
                        </Button>
                      </div>
                    </div>

                    <div className="space-y-2">
                      <DataCardGridHeader className="hidden md:grid">
                        <div className="col-span-3">Provider / Brand</div>
                        <div className="col-span-3">Branch / Drop Point</div>
                        <div className="col-span-3">Location Address</div>
                        <div className="col-span-2 text-right">Status</div>
                        <div className="col-span-1 text-right">Action</div>
                      </DataCardGridHeader>

                      <DataCardList>
                        {filteredCouriers.length === 0 ? (
                          <div className="py-10 border border-dashed border-border/80 rounded-2xl bg-zinc-50/10 text-center text-muted-foreground text-sm space-y-3">
                            <Truck className="h-8 w-8 text-slate-400 mb-1 mx-auto" />
                            {courierStatusFilter === 'active' && !courierSearchQuery ? (
                              <div className="space-y-1">
                                <p className="font-semibold text-foreground">No active couriers configured</p>
                                <p className="text-xs text-muted-foreground max-w-md mx-auto">
                                  Currently showing only active couriers by default. Switch to "All" to view and activate available couriers.
                                </p>
                                <div className="pt-2">
                                  <Button
                                    variant="outline"
                                    size="sm"
                                    onClick={() => setCourierStatusFilter('all')}
                                    className="rounded-xl text-xs"
                                  >
                                    View All Couriers ({couriers.length})
                                  </Button>
                                </div>
                              </div>
                            ) : (
                              <div className="space-y-1">
                                <p className="font-semibold text-foreground">No couriers found</p>
                                <p className="text-xs text-muted-foreground">No couriers match your search or filter selection.</p>
                                <div className="pt-2">
                                  <Button
                                    variant="ghost"
                                    size="sm"
                                    onClick={() => {
                                      setCourierSearchQuery('');
                                      setCourierStatusFilter('all');
                                    }}
                                    className="rounded-xl text-xs"
                                  >
                                    Reset Filters
                                  </Button>
                                </div>
                              </div>
                            )}
                          </div>
                        ) : (
                          filteredCouriers.map((courier) => (
                            <DataCard key={courier.code} className="bg-muted/20 hover:border-border transition-colors">
                              {/* Col 1: Provider / Brand */}
                              <div className="col-span-1 md:col-span-3 min-w-0 flex items-center gap-2.5">
                                <div className="p-1.5 rounded-lg bg-primary/10 text-primary shrink-0">
                                  <Truck className="h-4 w-4" />
                                </div>
                                <div className="min-w-0">
                                  <div className="font-semibold text-foreground text-sm truncate">
                                    {courier.branch_name || courier.code.toUpperCase()}
                                  </div>
                                  <div className="flex items-center gap-1.5 mt-0.5">
                                    <Badge variant="outline" className="font-mono text-[10px] uppercase px-1.5 py-0">
                                      {courier.code}
                                    </Badge>
                                  </div>
                                </div>
                              </div>

                              {/* Col 2: Branch / Drop Point */}
                              <div className="col-span-1 md:col-span-3 min-w-0 text-xs">
                                <span className="md:hidden font-medium text-foreground mr-1">Branch / Drop Point:</span>
                                {courier.name ? (
                                  <span className="font-medium text-foreground truncate block">{courier.name}</span>
                                ) : (
                                  <span className="text-muted-foreground italic">Not configured</span>
                                )}
                              </div>

                              {/* Col 3: Location Address */}
                              <div className="col-span-1 md:col-span-3 min-w-0 text-xs">
                                <span className="md:hidden font-medium text-foreground mr-1">Address:</span>
                                {courier.location_address ? (
                                  <div className="flex items-center gap-1.5 text-muted-foreground truncate" title={courier.location_address}>
                                    <MapPin className="h-3.5 w-3.5 shrink-0 text-slate-400" />
                                    <span className="truncate">{courier.location_address}</span>
                                  </div>
                                ) : (
                                  <span className="text-muted-foreground italic">No address registered</span>
                                )}
                              </div>

                              {/* Col 4: Status Badges */}
                              <div className="col-span-1 md:col-span-2 min-w-0 flex items-center gap-1.5 flex-wrap justify-start md:justify-end">
                                {getCourierVerificationBadge(courier.verification_status)}
                                <Badge
                                  variant={courier.active ? 'default' : 'secondary'}
                                  className={
                                    courier.active
                                      ? 'bg-primary/10 text-primary border-0 scale-90 origin-right'
                                      : 'bg-muted text-muted-foreground border-0 scale-90 origin-right'
                                  }
                                >
                                  {courier.active ? 'Active' : 'Disabled'}
                                </Badge>
                              </div>

                              {/* Col 5: Action Triple Dots */}
                              <div className="col-span-1 md:col-span-1 min-w-0 flex items-center justify-end">
                                <DropdownMenu>
                                  <DropdownMenuTrigger asChild>
                                    <Button variant="ghost" size="icon" className="h-8 w-8 rounded-xl hover:bg-muted">
                                      <MoreHorizontal className="h-4 w-4" />
                                      <span className="sr-only">Open menu</span>
                                    </Button>
                                  </DropdownMenuTrigger>
                                  <DropdownMenuContent align="end" className="w-48 rounded-xl shadow-lg">
                                    <DropdownMenuItem onClick={() => handleOpenEditCourier(courier)}>
                                      <Edit className="h-4 w-4 mr-2 text-muted-foreground" />
                                      Edit Configuration
                                    </DropdownMenuItem>

                                    {isAdmin && courier.verification_status === 'pending' && (
                                      <>
                                        <DropdownMenuSeparator className="my-1" />
                                        <DropdownMenuItem
                                          onClick={() => handleOpenVerifyModal(courier, 'verify')}
                                          className="text-emerald-600 focus:text-emerald-600 font-medium"
                                        >
                                          <CheckCircle2 className="h-4 w-4 mr-2 text-emerald-600" />
                                          Verify & Activate
                                        </DropdownMenuItem>
                                        <DropdownMenuItem
                                          onClick={() => handleOpenVerifyModal(courier, 'reject')}
                                          className="text-destructive focus:text-destructive"
                                        >
                                          <XCircle className="h-4 w-4 mr-2 text-destructive" />
                                          Reject Courier
                                        </DropdownMenuItem>
                                      </>
                                    )}
                                  </DropdownMenuContent>
                                </DropdownMenu>
                              </div>
                            </DataCard>
                          ))
                        )}
                      </DataCardList>
                    </div>
                  </div>
                </TabsContent>

                {/* Orders Tab */}
                <TabsContent value="orders" className="space-y-4 pt-2">
                  <ShopOrdersTab shopId={currentShopId} shopName={selectedShopInfo?.name} />
                </TabsContent>
              </Tabs>
            )}
          </div>
        )}

        {/* Add Shop Sheet */}
        <Sheet open={isAddShopOpen} onOpenChange={setIsAddShopOpen}>
          <SheetContent className="w-full sm:max-w-none md:w-[45vw] md:min-w-[45vw] overflow-y-auto border-l border-border/60 bg-background shadow-2xl">
            <SheetHeader className="mb-4">
              <SheetTitle className="font-display">Add New Shop</SheetTitle>
              <SheetDescription>Create a new store branch location.</SheetDescription>
            </SheetHeader>

            {addShopError && (
              <div className="p-3 text-sm text-destructive bg-destructive/10 rounded-xl border border-destructive/20 mb-4 font-sans">
                {addShopError}
              </div>
            )}

            <form onSubmit={handleAddShopSubmit} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="newShopName">Shop Name</Label>
                <Input
                  id="newShopName"
                  placeholder="e.g. Depok Warehouse, Surabaya Branch"
                  value={newShopName}
                  onChange={(e) => setNewShopName(e.target.value)}
                  required
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="newShopDesc">Description</Label>
                <textarea
                  id="newShopDesc"
                  className="flex min-h-[80px] w-full rounded-xl border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-1 focus:ring-primary/40 text-foreground"
                  placeholder="Brief description of the shop..."
                  value={newShopDesc}
                  onChange={(e) => setNewShopDesc(e.target.value)}
                />
              </div>

              {isAdmin ? (
                <div className="space-y-4 pt-2 border-t border-border/60">
                  <div className="space-y-2">
                    <Label htmlFor="newShopApprovalStatus">Approval Status</Label>
                    <select
                      id="newShopApprovalStatus"
                      className="flex h-10 w-full rounded-xl border border-border bg-background px-3 py-2 text-sm text-foreground focus:outline-none focus:ring-1 focus:ring-primary/40"
                      value={newShopApprovalStatus}
                      onChange={(e) => setNewShopApprovalStatus(e.target.value)}
                    >
                      <option value="pending">Pending Approval</option>
                      <option value="approved">Approved</option>
                      <option value="rejected">Rejected</option>
                    </select>
                  </div>

                  <div className="flex items-center space-x-2 pt-1">
                    <Checkbox
                      id="newShopIsActive"
                      checked={newShopIsActive}
                      onCheckedChange={(checked) => setNewShopIsActive(checked === true)}
                    />
                    <Label htmlFor="newShopIsActive" className="text-sm font-medium leading-none cursor-pointer">
                      Shop is Active
                    </Label>
                  </div>
                </div>
              ) : (
                <div className="p-3 bg-muted/40 rounded-xl border border-border text-xs text-muted-foreground flex items-center gap-2">
                  <ShieldCheck className="h-4 w-4 text-primary shrink-0" />
                  New shops are created as Pending Approval and Inactive until reviewed by an Admin.
                </div>
              )}

              <div className="pt-4 flex justify-end gap-2">
                <SheetClose asChild>
                  <Button type="button" variant="outline" className="rounded-xl">
                    Cancel
                  </Button>
                </SheetClose>
                <Button
                  type="submit"
                  disabled={isCreatingShop}
                  className="bg-primary hover:bg-primary/90 text-primary-foreground rounded-xl"
                >
                  {isCreatingShop && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                  Add Shop
                </Button>
              </div>
            </form>
          </SheetContent>
        </Sheet>

        {/* Edit Shop Sheet */}
        <Sheet open={isEditShopOpen} onOpenChange={setIsEditShopOpen}>
          <SheetContent className="w-full sm:max-w-none md:w-[45vw] md:min-w-[45vw] overflow-y-auto border-l border-border/60 bg-background shadow-2xl">
            <SheetHeader className="mb-4">
              <SheetTitle className="font-display">Edit Shop</SheetTitle>
              <SheetDescription>Update store branch parameters and status.</SheetDescription>
            </SheetHeader>

            {editShopError && (
              <div className="p-3 text-sm text-destructive bg-destructive/10 rounded-xl border border-destructive/20 mb-4 font-sans">
                {editShopError}
              </div>
            )}

            <form onSubmit={handleEditShopSubmit} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="editShopName">Shop Name</Label>
                <Input
                  id="editShopName"
                  placeholder="e.g. Depok Warehouse"
                  value={editShopName}
                  onChange={(e) => setEditShopName(e.target.value)}
                  required
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="editShopDesc">Description</Label>
                <textarea
                  id="editShopDesc"
                  className="flex min-h-[80px] w-full rounded-xl border border-border bg-background px-3 py-2 text-sm focus:outline-none focus:ring-1 focus:ring-primary/40 text-foreground"
                  placeholder="Brief description of the shop..."
                  value={editShopDesc}
                  onChange={(e) => setEditShopDesc(e.target.value)}
                />
              </div>

              {isAdmin ? (
                <div className="space-y-4 pt-2 border-t border-border/60">
                  <div className="space-y-2">
                    <Label htmlFor="editShopApprovalStatus">Approval Status</Label>
                    <select
                      id="editShopApprovalStatus"
                      className="flex h-10 w-full rounded-xl border border-border bg-background px-3 py-2 text-sm text-foreground focus:outline-none focus:ring-1 focus:ring-primary/40"
                      value={editShopApprovalStatus}
                      onChange={(e) => setEditShopApprovalStatus(e.target.value)}
                    >
                      <option value="pending">Pending Approval</option>
                      <option value="approved">Approved</option>
                      <option value="rejected">Rejected</option>
                    </select>
                  </div>

                  <div className="flex items-center space-x-2 pt-1">
                    <Checkbox
                      id="editShopIsActive"
                      checked={editShopIsActive}
                      onCheckedChange={(checked) => setEditShopIsActive(checked === true)}
                    />
                    <Label htmlFor="editShopIsActive" className="text-sm font-medium leading-none cursor-pointer">
                      Shop is Active
                    </Label>
                  </div>
                </div>
              ) : (
                <div className="p-3 bg-muted/40 rounded-xl border border-border text-xs text-muted-foreground flex items-center gap-2">
                  <ShieldCheck className="h-4 w-4 text-muted-foreground shrink-0" />
                  Approval and operating statuses are managed exclusively by Administrator accounts.
                </div>
              )}

              <div className="pt-4 flex justify-end gap-2">
                <Button
                  type="button"
                  variant="outline"
                  className="rounded-xl"
                  onClick={() => setIsEditShopOpen(false)}
                >
                  Cancel
                </Button>
                <Button
                  type="submit"
                  disabled={isSavingShop}
                  className="bg-primary hover:bg-primary/90 text-primary-foreground rounded-xl"
                >
                  {isSavingShop && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                  Save Changes
                </Button>
              </div>
            </form>
          </SheetContent>
        </Sheet>

        {/* Delete Shop Confirmation Dialog */}
        <Dialog open={isDeleteShopOpen} onOpenChange={setIsDeleteShopOpen}>
          <DialogContent className="max-w-md rounded-2xl">
            <DialogHeader>
              <DialogTitle className="flex items-center gap-2 text-destructive font-display">
                <AlertTriangle className="h-5 w-5" /> Delete Shop Location
              </DialogTitle>
              <DialogDescription className="pt-2 text-sm text-muted-foreground">
                Are you sure you want to permanently delete the shop branch <strong>{selectedShopInfo?.name}</strong>?
                This action will mark the shop as deleted and disallow all subsequent transactions.
              </DialogDescription>
            </DialogHeader>

            {deleteShopError && (
              <div className="p-3 text-sm text-destructive bg-destructive/10 rounded-xl border border-destructive/20 my-2 font-sans">
                {deleteShopError}
              </div>
            )}

            <DialogFooter className="gap-2 sm:gap-0">
              <Button
                type="button"
                variant="outline"
                className="rounded-xl"
                onClick={() => setIsDeleteShopOpen(false)}
                disabled={isDeletingShop}
              >
                Cancel
              </Button>
              <Button
                type="button"
                variant="destructive"
                className="rounded-xl"
                onClick={handleDeleteShopConfirm}
                disabled={isDeletingShop}
              >
                {isDeletingShop && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                Delete Shop
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>

        {/* Global Inventory Sheet - Add Mode */}
        <InventoryFormSheet
          open={isInventoryOpen}
          onOpenChange={setIsInventoryOpen}
          shopId={selectedShopId || ''}
          shopName={selectedShopInfo?.name}
          existingProducts={products}
          onSuccess={() => {
            setIsInventoryOpen(false);
            if (selectedShopId && selectedShopInfo) {
              selectShop(selectedShopInfo);
            }
          }}
        />

        {/* Global Inventory Sheet - Edit Mode */}
        <InventoryFormSheet
          open={!!editingInventoryProduct}
          onOpenChange={(open) => !open && setEditingInventoryProduct(null)}
          shopId={selectedShopId || ''}
          shopName={selectedShopInfo?.name}
          product={editingInventoryProduct}
          onSuccess={() => {
            setEditingInventoryProduct(null);
            if (selectedShopId && selectedShopInfo) {
              selectShop(selectedShopInfo);
            }
          }}
        />

        {/* Delete Inventory Confirmation Dialog */}
        <Dialog open={!!inventoryToDelete} onOpenChange={(open) => !open && setInventoryToDelete(null)}>
          <DialogContent className="max-w-md rounded-2xl">
            <DialogHeader>
              <DialogTitle className="flex items-center gap-2 text-destructive font-display">
                <AlertTriangle className="h-5 w-5" /> Delete Inventory
              </DialogTitle>
              <DialogDescription className="pt-2 text-sm text-muted-foreground">
                Are you sure you want to remove the inventory record for <strong>{inventoryToDelete?.name}</strong> from{' '}
                <strong>{selectedShopInfo?.name}</strong>? This action cannot be undone.
              </DialogDescription>
            </DialogHeader>

            {deleteInventoryError && (
              <div className="p-3 text-sm text-destructive bg-destructive/10 rounded-xl border border-destructive/20 my-2 font-sans">
                {deleteInventoryError}
              </div>
            )}

            <DialogFooter className="gap-2 sm:gap-0">
              <Button
                type="button"
                variant="outline"
                className="rounded-xl"
                onClick={() => setInventoryToDelete(null)}
                disabled={isDeletingInventory}
              >
                Cancel
              </Button>
              <Button
                type="button"
                variant="destructive"
                className="rounded-xl"
                onClick={handleDeleteInventoryConfirm}
                disabled={isDeletingInventory}
              >
                {isDeletingInventory && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                Delete
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>

        {/* Global Address Sheet - Add Mode */}
        <AddressFormSheet
          open={isAddressOpen}
          onOpenChange={setIsAddressOpen}
          shopId={selectedShopId || ''}
          shopName={selectedShopInfo?.name}
          onSuccess={() => {
            setIsAddressOpen(false);
            if (selectedShopId && selectedShopInfo) {
              selectShop(selectedShopInfo);
            }
          }}
        />

        {/* Global Address Sheet - Edit Mode */}
        <AddressFormSheet
          open={!!editingAddress}
          onOpenChange={(open) => !open && setEditingAddress(null)}
          shopId={selectedShopId || ''}
          shopName={selectedShopInfo?.name}
          address={editingAddress}
          onSuccess={() => {
            setEditingAddress(null);
            if (selectedShopId && selectedShopInfo) {
              selectShop(selectedShopInfo);
            }
          }}
        />

        {/* Delete Address Confirmation Dialog */}
        <Dialog open={!!addressToDelete} onOpenChange={(open) => !open && setAddressToDelete(null)}>
          <DialogContent className="max-w-md rounded-2xl">
            <DialogHeader>
              <DialogTitle className="flex items-center gap-2 text-destructive font-display">
                <AlertTriangle className="h-5 w-5" /> Delete Shop Address
              </DialogTitle>
              <DialogDescription className="pt-2 text-sm text-muted-foreground">
                Are you sure you want to remove the address <strong>{addressToDelete?.label}</strong> from{' '}
                <strong>{selectedShopInfo?.name}</strong>? Active addresses cannot be deleted directly without activating another address first.
              </DialogDescription>
            </DialogHeader>

            {deleteAddressError && (
              <div className="p-3 text-sm text-destructive bg-destructive/10 rounded-xl border border-destructive/20 my-2 font-sans">
                {deleteAddressError}
              </div>
            )}

            <DialogFooter className="gap-2 sm:gap-0">
              <Button
                type="button"
                variant="outline"
                className="rounded-xl"
                onClick={() => setAddressToDelete(null)}
                disabled={isDeletingAddress}
              >
                Cancel
              </Button>
              <Button
                type="button"
                variant="destructive"
                className="rounded-xl"
                onClick={handleDeleteAddressConfirm}
                disabled={isDeletingAddress}
              >
                {isDeletingAddress && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                Delete Address
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>

        {/* Global Courier Sheet */}
        <CourierFormSheet
          open={isCourierSheetOpen}
          onOpenChange={setIsCourierSheetOpen}
          shopId={selectedShopId || ''}
          shopName={selectedShopInfo?.name}
          courier={editingCourier}
          shopAddresses={addresses}
          isAdmin={isAdmin}
          onSave={updateCourier}
          onVerify={verifyCourier}
          onSuccess={() => {
            if (selectedShopId && selectedShopInfo) {
              selectShop(selectedShopInfo);
            }
          }}
        />

        {/* Admin Courier Verification Dialog */}
        <Dialog open={!!verifyingCourier} onOpenChange={(open) => !open && setVerifyingCourier(null)}>
          <DialogContent className="max-w-md rounded-2xl">
            <DialogHeader>
              <DialogTitle
                className={`flex items-center gap-2 font-display ${
                  verifyAction === 'verify' ? 'text-emerald-600' : 'text-destructive'
                }`}
              >
                {verifyAction === 'verify' ? (
                  <>
                    <CheckCircle2 className="h-5 w-5" /> Verify & Activate Courier
                  </>
                ) : (
                  <>
                    <XCircle className="h-5 w-5" /> Reject Courier
                  </>
                )}
              </DialogTitle>
              <DialogDescription className="pt-2 text-sm text-muted-foreground">
                {verifyAction === 'verify' ? (
                  <>
                    Are you sure you want to verify and activate{' '}
                    <strong>{verifyingCourier?.branch_name || verifyingCourier?.code?.toUpperCase()}</strong> for{' '}
                    <strong>{selectedShopInfo?.name}</strong>? This courier will immediately become available for checkout
                    orders.
                  </>
                ) : (
                  <>
                    Are you sure you want to reject{' '}
                    <strong>{verifyingCourier?.branch_name || verifyingCourier?.code?.toUpperCase()}</strong>? The courier
                    will remain inactive and staff can view your rejection reason.
                  </>
                )}
              </DialogDescription>
            </DialogHeader>

            {verificationError && (
              <div className="p-3 text-sm text-destructive bg-destructive/10 rounded-xl border border-destructive/20 my-2 font-sans">
                {verificationError}
              </div>
            )}

            {verifyAction === 'reject' && (
              <div className="space-y-2 py-2">
                <Label htmlFor="quickRejectionReason" className="text-xs">
                  Rejection Reason (visible to store staff)
                </Label>
                <Input
                  id="quickRejectionReason"
                  placeholder="e.g. Pickup address requires specific gate/room number"
                  value={verifyRejectionReason}
                  onChange={(e) => setVerifyRejectionReason(e.target.value)}
                  className="rounded-xl text-sm"
                />
              </div>
            )}

            <DialogFooter className="gap-2 sm:gap-0">
              <Button
                type="button"
                variant="outline"
                className="rounded-xl"
                onClick={() => setVerifyingCourier(null)}
                disabled={isProcessingVerification}
              >
                Cancel
              </Button>
              <Button
                type="button"
                variant={verifyAction === 'verify' ? 'default' : 'destructive'}
                className={`rounded-xl ${
                  verifyAction === 'verify' ? 'bg-emerald-600 hover:bg-emerald-700 text-white' : ''
                }`}
                onClick={handleConfirmVerification}
                disabled={isProcessingVerification}
              >
                {isProcessingVerification && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                {verifyAction === 'verify' ? 'Confirm & Activate' : 'Confirm Rejection'}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
    </div>
  );
}
