import { useState, useEffect } from 'react';
import { MapPin, Truck, Package, Loader2, Plus, Info, RefreshCw, Eye } from 'lucide-react';
import { Button } from '../../components/ui/button';
import { Badge } from '../../components/ui/badge';
import { Input } from '../../components/ui/input';
import { Label } from '../../components/ui/label';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '../../components/ui/table';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '../../components/ui/card';
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

export default function ShopManagementPage() {
  const {
    shops,
    selectedShopId,
    selectedShopInfo,
    addresses,
    couriers,
    products,
    loading,
    detailsLoading,
    error,
    detailsError,
    createAddress,
    saveShop,
    createShop,
    selectShop,
    refresh
  } = useShopViewModel();
  
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
  
  const [provinceId, setProvinceId] = useState('');
  const [cityId, setCityId] = useState('');
  const [districtId, setDistrictId] = useState('');
  const [villageId, setVillageId] = useState('');
  
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
      selectShop(selectedShopInfo);
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
          const res = await fetchApi('/provinces');
          if (res && res.provinces) {
            setProvinces(res.provinces);
          } else {
            setFetchError('No provinces list found in response');
          }
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
      const res = await fetchApi(`/provinces/${provId}/cities`);
      setCities(res.cities || []);
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
      const res = await fetchApi(`/cities/${cId}/districts`);
      setDistricts(res.districts || []);
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
      const res = await fetchApi(`/districts/${distId}/villages`);
      setVillages(res.villages || []);
    } catch (err: any) {
      console.error(err);
      setFetchError(err.message || 'Failed to load villages');
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
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

  const handleOpenDetails = (shop: any) => {
    selectShop(shop);
    setIsDetailsOpen(true);
  };

  if (loading && shops.length === 0) {
    return (
      <div className="flex h-[50vh] items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex h-[50vh] items-center justify-center">
        <p className="text-destructive">{error}</p>
      </div>
    );
  }

  return (
    <div className="flex-col md:flex">
      <div className="flex-1 space-y-4 p-8 pt-6">
        <div className="flex items-center justify-between space-y-2">
          <div>
            <h2 className="text-3xl font-bold tracking-tight">Shop Management</h2>
            <p className="text-muted-foreground">
              View and manage all your store branches and their details.
            </p>
          </div>
        </div>
        <Card className="shadow-md border-0 bg-white/70 backdrop-blur-md">
          <CardHeader className="flex flex-row items-center justify-between pb-4">
            <div>
              <CardTitle className="text-xl font-semibold">Store Locations</CardTitle>
              <CardDescription>
                You have {shops.length} shop locations registered.
              </CardDescription>
            </div>
            <div className="flex items-center gap-2">
              <Button
                variant="outline"
                size="sm"
                onClick={() => refresh()}
                className="flex items-center gap-1.5 border-slate-200 text-slate-600 hover:text-indigo-600 transition-colors"
              >
                <RefreshCw className="h-4 w-4" />
                Refresh
              </Button>
              <Button
                size="sm"
                onClick={() => setIsAddShopOpen(true)}
                className="bg-indigo-600 hover:bg-indigo-700 text-white flex items-center gap-1.5"
              >
                <Plus className="h-4 w-4" />
                Add Shop
              </Button>
            </div>
          </CardHeader>
          <CardContent>
            <div className="rounded-md border border-slate-100 overflow-hidden">
              <Table>
                <TableHeader className="bg-slate-50/70">
                  <TableRow>
                    <TableHead>Shop Name</TableHead>
                    <TableHead>Description</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead className="text-right">Action</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {shops.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={4} className="text-center h-24 text-slate-500">
                        No shops found.
                      </TableCell>
                    </TableRow>
                  ) : (
                    shops.map((shop) => (
                      <TableRow
                        key={shop.id}
                        className="hover:bg-slate-50/50 cursor-pointer transition-colors"
                        onClick={() => handleOpenDetails(shop)}
                      >
                        <TableCell className="font-semibold text-slate-800">
                          {shop.name}
                        </TableCell>
                        <TableCell className="text-slate-600 max-w-sm truncate">
                          {shop.description || '-'}
                        </TableCell>
                        <TableCell>
                          <Badge
                            variant={shop.is_active ? 'default' : 'secondary'}
                            className={
                              shop.is_active
                                ? 'bg-emerald-100 text-emerald-800 hover:bg-emerald-100/80 border-0'
                                : 'bg-slate-100 text-slate-800 hover:bg-slate-100/80 border-0'
                            }
                          >
                            {shop.is_active ? 'Active' : 'Inactive'}
                          </Badge>
                        </TableCell>
                        <TableCell className="text-right">
                          <Button
                            variant="ghost"
                            size="sm"
                            className="text-indigo-600 hover:text-indigo-700 hover:bg-indigo-50"
                            onClick={(e) => {
                              e.stopPropagation();
                              handleOpenDetails(shop);
                            }}
                          >
                            <Eye className="mr-1.5 h-4 w-4" />
                            View Details
                          </Button>
                        </TableCell>
                      </TableRow>
                    ))
                  )}
                </TableBody>
              </Table>
            </div>
          </CardContent>
        </Card>
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
                <Card className="border-0 shadow-none bg-slate-50/50">
                  <CardHeader>
                    <CardTitle className="text-lg">Edit Shop Profile</CardTitle>
                    <CardDescription>Update general settings for this shop branch.</CardDescription>
                  </CardHeader>
                  <CardContent>
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
                        <input
                          id="shopIsActive"
                          type="checkbox"
                          className="h-4 w-4 rounded border-slate-300 bg-white text-indigo-600 focus:ring-indigo-500"
                          checked={shopIsActive}
                          onChange={(e) => setShopIsActive(e.target.checked)}
                        />
                        <Label htmlFor="shopIsActive" className="text-sm font-medium leading-none cursor-pointer">
                          Shop is Active
                        </Label>
                      </div>
                      <div className="pt-2">
                        <Button type="submit" disabled={isSavingShop} className="bg-indigo-600 hover:bg-indigo-700">
                          {isSavingShop && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                          Save Settings
                        </Button>
                      </div>
                    </form>
                  </CardContent>
                </Card>
              </TabsContent>

              {/* Products Tab */}
              <TabsContent value="products" className="space-y-4">
                <Card className="border-slate-100 shadow-none">
                  <CardHeader className="flex flex-row items-center justify-between pb-4">
                    <div>
                      <CardTitle className="text-lg">Product Inventory</CardTitle>
                      <CardDescription>Listed products and current stock levels at this location.</CardDescription>
                    </div>
                    
                    <Sheet open={isInventoryOpen} onOpenChange={setIsInventoryOpen}>
                      <SheetTrigger asChild>
                        <Button size="sm" className="bg-indigo-600 hover:bg-indigo-700">
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
                            <Button type="submit" disabled={isInventorySubmitting} className="bg-indigo-600 hover:bg-indigo-700">
                              {isInventorySubmitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                              Add Stock
                            </Button>
                          </div>
                        </form>
                      </SheetContent>
                    </Sheet>
                  </CardHeader>
                  <CardContent>
                    <div className="rounded-md border border-slate-100 overflow-hidden">
                      <Table>
                        <TableHeader className="bg-slate-50/50">
                          <TableRow>
                            <TableHead>Product Name</TableHead>
                            <TableHead>SKU</TableHead>
                            <TableHead>Price</TableHead>
                            <TableHead className="text-right">Available Stock</TableHead>
                          </TableRow>
                        </TableHeader>
                        <TableBody>
                          {products.length === 0 ? (
                            <TableRow>
                              <TableCell colSpan={4} className="text-center h-24 text-slate-500">
                                No products found for this shop.
                              </TableCell>
                            </TableRow>
                          ) : (
                            products.map((product) => (
                              <TableRow key={product.id}>
                                <TableCell className="font-medium text-slate-800">
                                  {product.name}
                                  <div className="text-xs text-slate-400 font-normal">{product.slug}</div>
                                </TableCell>
                                <TableCell className="text-slate-600">{product.sku}</TableCell>
                                <TableCell className="text-slate-700">
                                  {new Intl.NumberFormat('id-ID', {
                                    style: 'currency',
                                    currency: 'IDR',
                                    minimumFractionDigits: 0,
                                  }).format(product.price)}
                                </TableCell>
                                <TableCell className="text-right">
                                  <span className="font-bold text-slate-800">
                                    {product.inventory.available}
                                  </span>
                                  <span className="text-xs text-slate-400 ml-2">
                                    (Total: {product.inventory.total_stock})
                                  </span>
                                </TableCell>
                              </TableRow>
                            ))
                          )}
                        </TableBody>
                      </Table>
                    </div>
                  </CardContent>
                </Card>
              </TabsContent>

              {/* Addresses Tab */}
              <TabsContent value="addresses" className="space-y-4">
                <Card className="border-slate-100 shadow-none">
                  <CardHeader className="flex flex-row items-center justify-between pb-4">
                    <div>
                      <CardTitle className="text-lg">Shop Addresses</CardTitle>
                      <CardDescription>Physical locations where this branch operates.</CardDescription>
                    </div>
                    
                    <Sheet open={isAddressOpen} onOpenChange={setIsAddressOpen}>
                      <SheetTrigger asChild>
                        <Button size="sm" className="bg-indigo-600 hover:bg-indigo-700">
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
                            <input
                              id="isActive"
                              type="checkbox"
                              className="h-4 w-4 rounded border-slate-300 bg-white text-indigo-600 focus:ring-indigo-500"
                              checked={isActive}
                              onChange={(e) => setIsActive(e.target.checked)}
                            />
                            <Label htmlFor="isActive" className="text-sm font-medium leading-none cursor-pointer">
                              Set as active address
                            </Label>
                          </div>

                          <div className="pt-4 flex justify-end gap-2">
                            <SheetClose asChild>
                              <Button type="button" variant="outline">Cancel</Button>
                            </SheetClose>
                            <Button type="submit" disabled={isSubmitting} className="bg-indigo-600 hover:bg-indigo-700">
                              {isSubmitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                              Save Address
                            </Button>
                          </div>
                        </form>
                      </SheetContent>
                    </Sheet>
                  </CardHeader>
                  <CardContent>
                    <div className="rounded-md border border-slate-100 overflow-hidden">
                      <Table>
                        <TableHeader className="bg-slate-50/50">
                          <TableRow>
                            <TableHead>Label</TableHead>
                            <TableHead>Full Address</TableHead>
                            <TableHead>Phone</TableHead>
                            <TableHead>Status</TableHead>
                          </TableRow>
                        </TableHeader>
                        <TableBody>
                          {addresses.length === 0 ? (
                            <TableRow>
                              <TableCell colSpan={4} className="text-center h-24 text-slate-500">No addresses found.</TableCell>
                            </TableRow>
                          ) : (
                            addresses.map((addr) => (
                              <TableRow key={addr.id}>
                                <TableCell className="font-semibold text-slate-800">{addr.label}</TableCell>
                                <TableCell className="max-w-xs truncate text-slate-600">{addr.full_address}</TableCell>
                                <TableCell className="text-slate-600">{addr.phone || '-'}</TableCell>
                                <TableCell>
                                  <Badge
                                    variant={addr.is_active ? "default" : "secondary"}
                                    className={
                                      addr.is_active
                                        ? "bg-indigo-50 text-indigo-700 hover:bg-indigo-50 border-0"
                                        : "bg-slate-100 text-slate-600 hover:bg-slate-100 border-0"
                                    }
                                  >
                                    {addr.is_active ? 'Active' : 'Inactive'}
                                  </Badge>
                                </TableCell>
                              </TableRow>
                            ))
                          )}
                        </TableBody>
                      </Table>
                    </div>
                  </CardContent>
                </Card>
              </TabsContent>

              {/* Couriers Tab */}
              <TabsContent value="couriers" className="space-y-4">
                <Card className="border-slate-100 shadow-none">
                  <CardHeader>
                    <CardTitle className="text-lg">Configured Couriers</CardTitle>
                    <CardDescription>Shipping providers configured for this branch.</CardDescription>
                  </CardHeader>
                  <CardContent>
                    <div className="rounded-md border border-slate-100 overflow-hidden">
                      <Table>
                        <TableHeader className="bg-slate-50/50">
                          <TableRow>
                            <TableHead>Courier Code</TableHead>
                            <TableHead className="text-right">Status</TableHead>
                          </TableRow>
                        </TableHeader>
                        <TableBody>
                          {couriers.length === 0 ? (
                            <TableRow>
                              <TableCell colSpan={2} className="text-center h-24 text-slate-500">No couriers found.</TableCell>
                            </TableRow>
                          ) : (
                            couriers.map((courier) => (
                              <TableRow key={courier.code}>
                                <TableCell className="font-semibold uppercase text-slate-800">{courier.code}</TableCell>
                                <TableCell className="text-right">
                                  <Badge
                                    variant={courier.active ? "default" : "secondary"}
                                    className={
                                      courier.active
                                        ? "bg-emerald-100 text-emerald-800 border-0"
                                        : "bg-slate-100 text-slate-600 border-0"
                                    }
                                  >
                                    {courier.active ? 'Active' : 'Disabled'}
                                  </Badge>
                                </TableCell>
                              </TableRow>
                            ))
                          )}
                        </TableBody>
                      </Table>
                    </div>
                  </CardContent>
                </Card>
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
              <input
                id="newShopIsActive"
                type="checkbox"
                className="h-4 w-4 rounded border-slate-300 bg-white text-indigo-600 focus:ring-indigo-500"
                checked={newShopIsActive}
                onChange={(e) => setNewShopIsActive(e.target.checked)}
              />
              <Label htmlFor="newShopIsActive" className="text-sm font-medium leading-none cursor-pointer">
                Shop is Active (Staff Admin Only)
              </Label>
            </div>

            <div className="pt-4 flex justify-end gap-2">
              <SheetClose asChild>
                <Button type="button" variant="outline">Cancel</Button>
              </SheetClose>
              <Button type="submit" disabled={isCreatingShop} className="bg-indigo-600 hover:bg-indigo-700">
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
