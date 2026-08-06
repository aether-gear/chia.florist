import { useState, useEffect, useMemo } from 'react';
import { MapPin, Truck, Package, Loader2, Plus, Info, RefreshCw, Eye } from 'lucide-react';
import { Button } from '../../components/ui/button';
import { Badge } from '../../components/ui/badge';
import { Input } from '../../components/ui/input';
import { Label } from '../../components/ui/label';
import { Checkbox } from '../../components/ui/checkbox';

// Removed Card component imports since sections are now borderless and backgroundless
import { Skeleton } from '../../components/ui/skeleton';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '../../components/ui/tabs';
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
  SheetClose,
} from '../../components/ui/sheet';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '../../components/ui/dialog';
import { useShopViewModel } from '../../viewmodels/useShopViewModel';
import { fetchApi } from '../../lib/api';
import Pagination from '../../components/Pagination';
import SearchInput from '../../components/SearchInput';
import { DataCard, DataCardGridHeader, DataCardList } from '../../components/DataCard';

export default function ShopManagementPage() {
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
    createAddress,
    saveShop,
    createShop,
    selectShop,
    refresh,
  } = useShopViewModel();

  const [searchQuery, setSearchQuery] = useState('');

  const filteredShops = useMemo(() => {
    if (!shops) return [];
    return shops.filter(shop =>
      shop.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      (shop.description && shop.description.toLowerCase().includes(searchQuery.toLowerCase()))
    );
  }, [shops, searchQuery]);

  // Overlay details Dialog state
  const [isDetailsOpen, setIsDetailsOpen] = useState(false);

  // Sheet control (Add Address)
  const [isAddressOpen, setIsAddressOpen] = useState(false);

  // Form states (Add Address)
  const [label, setLabel] = useState('');
  const [phone, setPhone] = useState('');
  const [fullAddress, setFullAddress] = useState('');
  const [postalCode, setPostalCode] = useState('');
  const [isActive, setIsActive] = useState(false);

  // Shop Info Form states
  const [shopName, setShopName] = useState('');
  const [shopDesc, setShopDesc] = useState('');
  const [shopIsActive, setShopIsActive] = useState(false);
  const [isSavingShop, setIsSavingShop] = useState(false);

  // Add Shop states
  const [isAddShopOpen, setIsAddShopOpen] = useState(false);
  const [newShopName, setNewShopName] = useState('');
  const [newShopDesc, setNewShopDesc] = useState('');
  const [newShopIsActive, setNewShopIsActive] = useState(false);
  const [isCreatingShop, setIsCreatingShop] = useState(false);
  const [addShopError, setAddShopError] = useState<string | null>(null);

  // Sync shop info to form
  useEffect(() => {
    if (selectedShopInfo) {
      setShopName(selectedShopInfo.name || '');
      setShopDesc(selectedShopInfo.description || '');
      setShopIsActive(selectedShopInfo.is_active || false);
    }
  }, [selectedShopInfo]);

  const handleAddShopSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newShopName) return;
    setIsCreatingShop(true);
    setAddShopError(null);
    try {
      const success = await createShop({
        name: newShopName,
        description: newShopDesc || undefined,
        is_active: newShopIsActive ? "true" : "false"
      });
      if (success) {
        setIsAddShopOpen(false);
        setNewShopName('');
        setNewShopDesc('');
        setNewShopIsActive(false);
      } else {
        setAddShopError("Failed to add shop. Please check your inputs or try again.");
      }
    } catch (err: any) {
      setAddShopError(err.message || "Failed to add shop");
    } finally {
      setIsCreatingShop(false);
    }
  };

  const handleSaveShop = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSavingShop(true);
    await saveShop({
      name: shopName,
      description: shopDesc,
      is_active: shopIsActive ? "true" : "false"
    });
    setIsSavingShop(false);
  };

  // Location dropdown states
  const [provinces, setProvinces] = useState<any[]>([]);
  const [cities, setCities] = useState<any[]>([]);
  const [districts, setDistricts] = useState<any[]>([]);
  const [villages, setVillages] = useState<any[]>([]);

  const [provinceId, setProvinceId] = useState<string>('');
  const [cityId, setCityId] = useState<string>('');
  const [districtId, setDistrictId] = useState<string>('');
  const [villageId, setVillageId] = useState<string>('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [fetchError, setFetchError] = useState<string | null>(null);

  // Add Inventory Sheet states
  const [isInventoryOpen, setIsInventoryOpen] = useState(false);
  const [productsList, setProductsList] = useState<any[]>([]);
  const [selectedProductId, setSelectedProductId] = useState('');
  const [inventoryStock, setInventoryStock] = useState('');
  const [isInventorySubmitting, setIsInventorySubmitting] = useState(false);
  const [inventoryError, setInventoryError] = useState<string | null>(null);

  // Load all products when Add Inventory Sheet is open
  useEffect(() => {
    if (isInventoryOpen) {
      const loadProducts = async () => {
        try {
          const res = await fetchApi('/products?page=1&limit=100');
          setProductsList(res.products || []);
        } catch (err: any) {
          console.error('Failed to load products', err);
        }
      };
      loadProducts();
    }
  }, [isInventoryOpen]);

  const handleInventorySubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedShopId || !selectedProductId || !inventoryStock) return;

    setIsInventorySubmitting(true);
    setInventoryError(null);

    try {
      await fetchApi(`/shops/${selectedShopId}/products/${selectedProductId}/inventories`, {
        method: 'POST',
        body: JSON.stringify({ stock: Number(inventoryStock) }),
      });
      setIsInventoryOpen(false);
      setSelectedProductId('');
      setInventoryStock('');
      // Re-fetch details of the active shop
      if (selectedShopInfo) {
        selectShop(selectedShopInfo);
      }
    } catch (err: any) {
      setInventoryError(err.message || 'Failed to add inventory');
    } finally {
      setIsInventorySubmitting(false);
    }
  };

  // Load provinces on Add Address Sheet open
  useEffect(() => {
    if (isAddressOpen) {
      const loadProvinces = async () => {
        try {
          setFetchError(null);
          const res = await fetchApi('/locations/provinces');
          setProvinces(res.provinces || res || []);
        } catch (err: any) {
          console.error('Failed to load provinces', err);
          setFetchError(err.message || 'Failed to load provinces');
        }
      };
      loadProvinces();
    }
  }, [isAddressOpen]);

  const handleProvinceChange = async (provId: string) => {
    setProvinceId(provId);
    setCityId('');
    setDistrictId('');
    setVillageId('');
    setCities([]);
    setDistricts([]);
    setVillages([]);
    if (!provId) return;
    try {
      setFetchError(null);
      const res = await fetchApi(`/locations/provinces/${provId}/cities`);
      setCities(res.cities || res || []);
    } catch (err: any) {
      console.error(err);
      setFetchError(err.message || 'Failed to load cities');
    }
  };

  const handleCityChange = async (cId: string) => {
    setCityId(cId);
    setDistrictId('');
    setVillageId('');
    setDistricts([]);
    setVillages([]);
    if (!cId) return;
    try {
      setFetchError(null);
      const res = await fetchApi(`/locations/cities/${cId}/districts`);
      setDistricts(res.districts || res || []);
    } catch (err: any) {
      console.error(err);
      setFetchError(err.message || 'Failed to load districts');
    }
  };

  const handleDistrictChange = async (distId: string) => {
    setDistrictId(distId);
    setVillageId('');
    setVillages([]);
    if (!distId) return;
    try {
      setFetchError(null);
      const res = await fetchApi(`/locations/districts/${distId}/villages`);
      setVillages(res.villages || res || []);
    } catch (err: any) {
      console.error(err);
      setFetchError(err.message || 'Failed to load villages');
    }
  };

  const handleCreateAddress = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    const success = await createAddress({
      label,
      phone: phone || null,
      province_id: provinceId,
      city_id: cityId,
      district_id: districtId,
      village_id: villageId,
      full_address: fullAddress,
      postal_code: postalCode,
      is_active: isActive,
    });
    setIsSubmitting(false);
    if (success) {
      setIsAddressOpen(false);
      // Reset form
      setLabel('');
      setPhone('');
      setFullAddress('');
      setPostalCode('');
      setProvinceId('');
      setCityId('');
      setDistrictId('');
      setVillageId('');
      setIsActive(false);
    }
  };

  const handleSubmit = handleCreateAddress;

  const handleOpenDetails = (shop: any) => {
    selectShop(shop);
    setIsDetailsOpen(true);
  };

  return (
    <div className="flex-col md:flex">
      <div className="flex-1 space-y-12 p-6 sm:p-8 lg:p-12 animate-in fade-in duration-300">
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

          <div className="flex items-center justify-between px-1 text-xs text-muted-foreground">
            <span>Found {filteredShops.length} shops</span>
          </div>

          {/* Content */}
          <div className="flex flex-col gap-2">
            <DataCardGridHeader>
              <span className="col-span-4">Shop Name</span>
              <span className="col-span-4">Description</span>
              <span className="col-span-2">Status</span>
              <span className="col-span-2 text-right">Action</span>
            </DataCardGridHeader>

            <DataCardList>
              {loading ? (
                Array.from({ length: 4 }).map((_, i) => (
                  <DataCard key={`skeleton-${i}`}>
                    <div className="col-span-4"><Skeleton className="h-5 w-40 bg-muted animate-pulse" /></div>
                    <div className="col-span-4"><Skeleton className="h-4 w-60 bg-muted animate-pulse" /></div>
                    <div className="col-span-2"><Skeleton className="h-5 w-16 bg-muted animate-pulse rounded-full" /></div>
                    <div className="col-span-2 text-right"><Skeleton className="h-8 w-24 ml-auto bg-muted animate-pulse rounded-lg" /></div>
                  </DataCard>
                ))
              ) : error ? (
                <div className="py-12 border-0 bg-transparent text-destructive text-center">Failed to load shops: {error}</div>
              ) : filteredShops.length === 0 ? (
                <div className="py-12 border border-dashed border-border/80 rounded-2xl bg-zinc-50/10 text-center text-muted-foreground">
                  <MapPin className="h-8 w-8 text-slate-400 mb-2 mx-auto" />
                  <p>No shops found</p>
                  <p className="text-sm">No store locations match your search.</p>
                </div>
              ) : (
                filteredShops.map((shop) => (
                  <DataCard key={shop.id} onClick={() => handleOpenDetails(shop)}>
                    <div className="col-span-1 md:col-span-4 min-w-0">
                      <h4 className="font-semibold font-display text-sm text-foreground truncate">{shop.name}</h4>
                      <p className="text-[10px] text-muted-foreground font-mono truncate">ID: {shop.id}</p>
                    </div>

                    <div className="col-span-1 md:col-span-4 text-xs text-muted-foreground truncate">
                      {shop.description || 'No description provided.'}
                    </div>

                    <div className="col-span-1 md:col-span-2">
                      <Badge
                        variant={shop.is_active ? 'default' : 'secondary'}
                        className={
                          shop.is_active
                            ? 'bg-primary/10 text-primary hover:bg-primary/10 border-0 rounded-lg scale-90 origin-left'
                            : 'bg-muted text-muted-foreground hover:bg-muted border-0 rounded-lg scale-90 origin-left'
                        }
                      >
                        {shop.is_active ? 'Active' : 'Inactive'}
                      </Badge>
                    </div>

                    <div className="col-span-1 md:col-span-2 text-right" onClick={(e) => e.stopPropagation()}>
                      <Button
                        variant="ghost"
                        size="sm"
                        className="text-primary hover:text-primary/90 hover:bg-primary/5 rounded-lg text-xs"
                        onClick={() => handleOpenDetails(shop)}
                      >
                        <Eye className="mr-1.5 h-3.5 w-3.5" />
                        View Details
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

      {/* Shop Details Central Overlay Dialog */}
      <Dialog open={isDetailsOpen} onOpenChange={setIsDetailsOpen}>
        <DialogContent className="max-w-4xl max-h-[85vh] overflow-y-auto">
          <DialogHeader className="border-b pb-4">
            <DialogTitle className="text-2xl font-bold text-slate-800">
              {selectedShopInfo?.name || 'Shop Details'}
            </DialogTitle>
            <DialogDescription className="text-slate-500">
              {selectedShopInfo?.description || 'Manage parameters, addresses, couriers, and product inventory.'}
            </DialogDescription>
          </DialogHeader>

          {detailsLoading && (
            <div className="flex justify-center py-6">
              <Loader2 className="h-6 w-6 animate-spin text-indigo-600" />
            </div>
          )}

          {detailsError && (
            <div className="p-3 text-sm text-red-500 bg-red-50 rounded-md border border-red-100 my-2">
              {detailsError}
            </div>
          )}

          {!detailsLoading && selectedShopInfo && (
            <Tabs defaultValue="products" className="space-y-4 pt-2">
              <TabsList className="grid grid-cols-4 bg-slate-50 p-1 rounded-lg border border-slate-100">
                <TabsTrigger value="info" className="flex items-center gap-1.5 data-[state=active]:bg-white data-[state=active]:shadow-sm">
                  <Info className="h-4 w-4" /> General
                </TabsTrigger>
                <TabsTrigger value="products" className="flex items-center gap-1.5 data-[state=active]:bg-white data-[state=active]:shadow-sm">
                  <Package className="h-4 w-4" /> Products
                </TabsTrigger>
                <TabsTrigger value="addresses" className="flex items-center gap-1.5 data-[state=active]:bg-white data-[state=active]:shadow-sm">
                  <MapPin className="h-4 w-4" /> Addresses
                </TabsTrigger>
                <TabsTrigger value="couriers" className="flex items-center gap-1.5 data-[state=active]:bg-white data-[state=active]:shadow-sm">
                  <Truck className="h-4 w-4" /> Couriers
                </TabsTrigger>
              </TabsList>

              {/* General Tab */}
              <TabsContent value="info" className="space-y-4">
                <div className="space-y-6">
                  <div className="pb-4 border-b border-border/60">
                    <h3 className="text-lg font-bold text-foreground">Edit Shop Profile</h3>
                    <p className="text-muted-foreground text-sm">Update general settings for this shop branch.</p>
                  </div>
                  <div>
                    <form onSubmit={handleSaveShop} className="space-y-4 max-w-md">
                      <div className="space-y-2">
                        <Label htmlFor="shopName">Shop Name</Label>
                        <Input
                          id="shopName"
                          value={shopName}
                          onChange={(e) => setShopName(e.target.value)}
                          required
                          placeholder="e.g. Chia Branch 1"
                        />
                      </div>
                      <div className="space-y-2">
                        <Label htmlFor="shopDesc">Description</Label>
                        <textarea
                          id="shopDesc"
                          className="flex min-h-[80px] w-full rounded-md border border-slate-200 bg-white px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
                          value={shopDesc}
                          onChange={(e) => setShopDesc(e.target.value)}
                          placeholder="Brief description of the shop..."
                        />
                      </div>
                      <div className="flex items-center space-x-2 pt-2">
                        <Checkbox
                          id="shopIsActive"
                          checked={shopIsActive}
                          onCheckedChange={(checked) => setShopIsActive(checked === true)}
                        />
                        <Label htmlFor="shopIsActive" className="text-sm font-medium leading-none cursor-pointer">
                          Shop is Active
                        </Label>
                      </div>
                      <div className="pt-2">
                        <Button type="submit" disabled={isSavingShop}>
                          {isSavingShop && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                          Save Settings
                        </Button>
                      </div>
                    </form>
                  </div>
                </div>
              </TabsContent>

              {/* Products Tab */}
              <TabsContent value="products" className="space-y-4">
                <div className="space-y-6">
                  <div className="flex flex-row items-center justify-between pb-4 border-b border-border/60">
                    <div>
                      <h3 className="text-lg font-bold text-foreground">Product Inventory</h3>
                      <p className="text-muted-foreground text-sm">Listed products and current stock levels at this location.</p>
                    </div>

                    <Sheet open={isInventoryOpen} onOpenChange={setIsInventoryOpen}>
                      <SheetTrigger asChild>
                        <Button size="sm">
                          <Plus className="mr-1.5 h-4 w-4" /> Add Inventory
                        </Button>
                      </SheetTrigger>
                      <SheetContent className="sm:max-w-md overflow-y-auto">
                        <SheetHeader className="mb-4">
                          <SheetTitle>Add Product Inventory</SheetTitle>
                          <SheetDescription>Assign stock for a catalog product at this store.</SheetDescription>
                        </SheetHeader>

                        {inventoryError && (
                          <div className="p-3 text-sm text-red-500 bg-red-50 rounded-md border border-red-100 mb-4">
                            {inventoryError}
                          </div>
                        )}

                        <form onSubmit={handleInventorySubmit} className="space-y-4">
                          <div className="space-y-2">
                            <Label htmlFor="productSelect">Select Product</Label>
                            <select
                              id="productSelect"
                              className="flex h-10 w-full rounded-md border border-slate-200 bg-white px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
                              value={selectedProductId}
                              onChange={(e) => setSelectedProductId(e.target.value)}
                              required
                            >
                              <option value="">Select Product</option>
                              {productsList.map((p) => (
                                <option key={p.id} value={p.id}>
                                  {p.name} ({p.sku})
                                </option>
                              ))}
                            </select>
                          </div>

                          <div className="space-y-2">
                            <Label htmlFor="inventoryStockInput">Initial Stock</Label>
                            <Input
                              id="inventoryStockInput"
                              type="number"
                              min="0"
                              placeholder="e.g. 50"
                              value={inventoryStock}
                              onChange={(e) => setInventoryStock(e.target.value)}
                              required
                            />
                          </div>

                          <div className="pt-4 flex justify-end gap-2">
                            <SheetClose asChild>
                              <Button type="button" variant="outline">Cancel</Button>
                            </SheetClose>
                            <Button type="submit" disabled={isInventorySubmitting}>
                              {isInventorySubmitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                              Add Stock
                            </Button>
                          </div>
                        </form>
                      </SheetContent>
                    </Sheet>
                  </div>

                  {/* Content */}
                  <div className="flex flex-col gap-2">
                    <DataCardGridHeader>
                      <span className="col-span-4">Product Name</span>
                      <span className="col-span-3">SKU</span>
                      <span className="col-span-3">Price</span>
                      <span className="col-span-2 text-right">Available Stock</span>
                    </DataCardGridHeader>

                    <DataCardList>
                      {products.length === 0 ? (
                        <div className="py-8 border border-dashed border-border/80 rounded-2xl bg-zinc-50/10 text-center text-muted-foreground text-sm">
                          No products found for this shop.
                        </div>
                      ) : (
                        products.map((product) => (
                          <DataCard key={product.id}>
                            <div className="col-span-1 md:col-span-4 min-w-0">
                              <h4 className="font-semibold text-sm text-foreground truncate">{product.name}</h4>
                              <p className="text-xs text-muted-foreground font-mono truncate">{product.slug}</p>
                            </div>
                            <div className="col-span-1 md:col-span-3 text-xs font-mono text-muted-foreground">
                              <span className="md:hidden font-sans mr-1">SKU:</span>
                              {product.sku || '-'}
                            </div>
                            <div className="col-span-1 md:col-span-3 text-xs font-bold text-primary">
                              {new Intl.NumberFormat('id-ID', {
                                style: 'currency',
                                currency: 'IDR',
                                minimumFractionDigits: 0,
                              }).format(product.price)}
                            </div>
                            <div className="col-span-1 md:col-span-2 text-right text-xs">
                              <span className="font-bold text-foreground">{product.inventory.available}</span>
                              <span className="text-[11px] text-muted-foreground ml-1.5">(Total: {product.inventory.total_stock})</span>
                            </div>
                          </DataCard>
                        ))
                      )}
                    </DataCardList>
                  </div>
                </div>
              </TabsContent>

              {/* Addresses Tab */}
              <TabsContent value="addresses" className="space-y-4">
                <div className="space-y-6">
                  <div className="flex flex-row items-center justify-between pb-4 border-b border-border/60">
                    <div>
                      <h3 className="text-lg font-bold text-foreground">Shop Addresses</h3>
                      <p className="text-muted-foreground text-sm">Physical locations where this branch operates.</p>
                    </div>

                    <Sheet open={isAddressOpen} onOpenChange={setIsAddressOpen}>
                      <SheetTrigger asChild>
                        <Button size="sm">
                          <Plus className="mr-1.5 h-4 w-4" /> Add Address
                        </Button>
                      </SheetTrigger>
                      <SheetContent className="sm:max-w-md overflow-y-auto">
                        <SheetHeader className="mb-4">
                          <SheetTitle>Add Shop Address</SheetTitle>
                          <SheetDescription>Create a physical pickup location for this shop.</SheetDescription>
                        </SheetHeader>

                        {fetchError && (
                          <div className="p-3 text-sm text-red-500 bg-red-50 rounded-md border border-red-100 mb-4">
                            {fetchError}
                          </div>
                        )}

                        <form onSubmit={handleSubmit} className="space-y-4">
                          <div className="space-y-2">
                            <Label htmlFor="label">Address Label / Branch Name</Label>
                            <Input
                              id="label"
                              placeholder="e.g. Main Branch, Jakarta Warehouse"
                              value={label}
                              onChange={(e) => setLabel(e.target.value)}
                              required
                            />
                          </div>

                          <div className="space-y-2">
                            <Label htmlFor="phone">Phone Number</Label>
                            <Input
                              id="phone"
                              placeholder="e.g. +62812345678"
                              value={phone}
                              onChange={(e) => setPhone(e.target.value)}
                            />
                          </div>

                          <div className="space-y-2">
                            <Label htmlFor="province">Province</Label>
                            <select
                              id="province"
                              className="flex h-10 w-full rounded-md border border-slate-200 bg-white px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
                              value={provinceId}
                              onChange={(e) => handleProvinceChange(e.target.value)}
                              required
                            >
                              <option value="">Select Province</option>
                              {provinces.map((p) => (
                                <option key={p.id} value={p.id}>{p.name}</option>
                              ))}
                            </select>
                          </div>

                          <div className="space-y-2">
                            <Label htmlFor="city">City / Regency</Label>
                            <select
                              id="city"
                              className="flex h-10 w-full rounded-md border border-slate-200 bg-white px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 disabled:opacity-50"
                              value={cityId}
                              onChange={(e) => handleCityChange(e.target.value)}
                              disabled={!provinceId}
                              required
                            >
                              <option value="">Select City</option>
                              {cities.map((c) => (
                                <option key={c.id} value={c.id}>{c.name}</option>
                              ))}
                            </select>
                          </div>

                          <div className="space-y-2">
                            <Label htmlFor="district">District</Label>
                            <select
                              id="district"
                              className="flex h-10 w-full rounded-md border border-slate-200 bg-white px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 disabled:opacity-50"
                              value={districtId}
                              onChange={(e) => handleDistrictChange(e.target.value)}
                              disabled={!cityId}
                              required
                            >
                              <option value="">Select District</option>
                              {districts.map((d) => (
                                <option key={d.id} value={d.id}>{d.name}</option>
                              ))}
                            </select>
                          </div>

                          <div className="space-y-2">
                            <Label htmlFor="village">Village</Label>
                            <select
                              id="village"
                              className="flex h-10 w-full rounded-md border border-slate-200 bg-white px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 disabled:opacity-50"
                              value={villageId}
                              onChange={(e) => setVillageId(e.target.value)}
                              disabled={!districtId}
                              required
                            >
                              <option value="">Select Village</option>
                              {villages.map((v) => (
                                <option key={v.id} value={v.id}>{v.name}</option>
                              ))}
                            </select>
                          </div>

                          <div className="space-y-2">
                            <Label htmlFor="fullAddress">Full Address</Label>
                            <textarea
                              id="fullAddress"
                              className="flex min-h-[80px] w-full rounded-md border border-slate-200 bg-white px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
                              placeholder="Street name, building number, unit, etc."
                              value={fullAddress}
                              onChange={(e) => setFullAddress(e.target.value)}
                              required
                            />
                          </div>

                          <div className="space-y-2">
                            <Label htmlFor="postalCode">Postal Code</Label>
                            <Input
                              id="postalCode"
                              placeholder="e.g. 17520"
                              value={postalCode}
                              onChange={(e) => setPostalCode(e.target.value)}
                              required
                            />
                          </div>

                          <div className="flex items-center space-x-2 pt-2">
                            <Checkbox
                              id="isActive"
                              checked={isActive}
                              onCheckedChange={(checked) => setIsActive(checked === true)}
                            />
                            <Label htmlFor="isActive" className="text-sm font-medium leading-none cursor-pointer">
                              Set as active address
                            </Label>
                          </div>

                          <div className="pt-4 flex justify-end gap-2">
                            <SheetClose asChild>
                              <Button type="button" variant="outline">Cancel</Button>
                            </SheetClose>
                            <Button type="submit" disabled={isSubmitting}>
                              {isSubmitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                              Save Address
                            </Button>
                          </div>
                        </form>
                      </SheetContent>
                    </Sheet>
                  </div>

                  {/* Content */}
                  <div className="flex flex-col gap-2">
                    <DataCardGridHeader>
                      <span className="col-span-3">Label</span>
                      <span className="col-span-5">Full Address</span>
                      <span className="col-span-2">Phone</span>
                      <span className="col-span-2 text-right">Status</span>
                    </DataCardGridHeader>

                    <DataCardList>
                      {addresses.length === 0 ? (
                        <div className="py-8 border border-dashed border-border/80 rounded-2xl bg-zinc-50/10 text-center text-muted-foreground text-sm">
                          No addresses found.
                        </div>
                      ) : (
                        addresses.map((addr) => (
                          <DataCard key={addr.id}>
                            <div className="col-span-1 md:col-span-3 font-semibold text-foreground text-sm truncate">{addr.label}</div>
                            <div className="col-span-1 md:col-span-5 text-xs text-muted-foreground truncate">{addr.full_address}</div>
                            <div className="col-span-1 md:col-span-2 text-xs text-muted-foreground">{addr.phone || '-'}</div>
                            <div className="col-span-1 md:col-span-2 text-right">
                              <Badge
                                variant={addr.is_active ? "default" : "secondary"}
                                className={
                                  addr.is_active
                                    ? "bg-primary/10 text-primary border-0 scale-90 origin-right"
                                    : "bg-muted text-muted-foreground border-0 scale-90 origin-right"
                                }
                              >
                                {addr.is_active ? 'Active' : 'Inactive'}
                              </Badge>
                            </div>
                          </DataCard>
                        ))
                      )}
                    </DataCardList>
                  </div>
                </div>
              </TabsContent>

              {/* Couriers Tab */}
              <TabsContent value="couriers" className="space-y-4">
                <div className="space-y-6">
                  <div className="pb-4 border-b border-border/60">
                    <h3 className="text-lg font-bold text-foreground">Configured Couriers</h3>
                    <p className="text-muted-foreground text-sm">Shipping providers configured for this branch.</p>
                  </div>

                  {/* Content */}
                  <div className="flex flex-col gap-2">
                    <DataCardGridHeader>
                      <span className="col-span-8">Courier Code</span>
                      <span className="col-span-4 text-right">Status</span>
                    </DataCardGridHeader>

                    <DataCardList>
                      {couriers.length === 0 ? (
                        <div className="py-8 border border-dashed border-border/80 rounded-2xl bg-zinc-50/10 text-center text-muted-foreground text-sm">
                          No couriers found.
                        </div>
                      ) : (
                        couriers.map((courier) => (
                          <DataCard key={courier.code}>
                            <div className="col-span-1 md:col-span-8 font-semibold uppercase text-foreground text-sm">{courier.code}</div>
                            <div className="col-span-1 md:col-span-4 text-right">
                              <Badge
                                variant={courier.active ? "default" : "secondary"}
                                className={
                                  courier.active
                                    ? "bg-primary/10 text-primary border-0 scale-90 origin-right"
                                    : "bg-muted text-muted-foreground border-0 scale-90 origin-right"
                                }
                              >
                                {courier.active ? 'Active' : 'Disabled'}
                              </Badge>
                            </div>
                          </DataCard>
                        ))
                      )}
                    </DataCardList>
                  </div>
                </div>
              </TabsContent>
            </Tabs>
          )}
        </DialogContent>
      </Dialog>

      {/* Add Shop Sheet */}
      <Sheet open={isAddShopOpen} onOpenChange={setIsAddShopOpen}>
        <SheetContent className="sm:max-w-md overflow-y-auto">
          <SheetHeader className="mb-4">
            <SheetTitle>Add New Shop</SheetTitle>
            <SheetDescription>
              Create a new store branch location.
            </SheetDescription>
          </SheetHeader>

          {addShopError && (
            <div className="p-3 text-sm text-red-500 bg-red-50 rounded-md border border-red-100 mb-4">
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
                className="flex min-h-[80px] w-full rounded-md border border-slate-200 bg-white px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500"
                placeholder="Brief description of the shop..."
                value={newShopDesc}
                onChange={(e) => setNewShopDesc(e.target.value)}
              />
            </div>

            <div className="flex items-center space-x-2 pt-2">
              <Checkbox
                id="newShopIsActive"
                checked={newShopIsActive}
                onCheckedChange={(checked) => setNewShopIsActive(checked === true)}
              />
              <Label htmlFor="newShopIsActive" className="text-sm font-medium leading-none cursor-pointer">
                Shop is Active (Staff Admin Only)
              </Label>
            </div>

            <div className="pt-4 flex justify-end gap-2">
              <SheetClose asChild>
                <Button type="button" variant="outline">Cancel</Button>
              </SheetClose>
              <Button type="submit" disabled={isCreatingShop}>
                {isCreatingShop && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                Add Shop
              </Button>
            </div>
          </form>
        </SheetContent>
      </Sheet>
    </div>
  );
}
