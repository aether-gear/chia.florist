import { useState, useEffect } from 'react';
import { MapPin, Truck, Package, Loader2, Plus, Info } from 'lucide-react';
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
import { useShopViewModel } from '../../viewmodels/useShopViewModel';
import { fetchApi } from '../../lib/api';

export default function ShopManagementPage() {
  const { shopId, shopInfo, addresses, couriers, products, loading, error, createAddress, saveShop, refresh } = useShopViewModel();
  
  // Sheet control
  const [isOpen, setIsOpen] = useState(false);
  
  // Form states
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

  // Sync shop info to form
  useEffect(() => {
    if (shopInfo) {
      setShopName(shopInfo.name || '');
      setShopDesc(shopInfo.description || '');
      setShopIsActive(shopInfo.is_active || false);
    }
  }, [shopInfo]);

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
    const activeShopId = shopId || "c3d4e5f6-a7b8-9012-cdef-123456789012";
    if (!activeShopId || !selectedProductId || !inventoryStock) return;

    setIsInventorySubmitting(true);
    setInventoryError(null);

    try {
      await fetchApi(`/shops/${activeShopId}/products/${selectedProductId}/inventories`, {
        method: 'POST',
        body: JSON.stringify({ stock: Number(inventoryStock) }),
      });
      setIsInventoryOpen(false);
      setSelectedProductId('');
      setInventoryStock('');
      refresh();
    } catch (err: any) {
      setInventoryError(err.message || 'Failed to add inventory');
    } finally {
      setIsInventorySubmitting(false);
    }
  };


  // Load provinces on Sheet open
  useEffect(() => {
    if (isOpen) {
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
  }, [isOpen]);

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
      setIsOpen(false);
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

  if (loading && addresses.length === 0) {
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
            <h2 className="text-3xl font-bold tracking-tight">Shop Details</h2>
            <p className="text-muted-foreground">
              Manage your shop's addresses, couriers, and listed products.
            </p>
          </div>
        </div>

        <Tabs defaultValue="addresses" className="space-y-4">
          <TabsList>
            <TabsTrigger value="info" className="flex items-center gap-2">
              <Info className="h-4 w-4" /> General Info
            </TabsTrigger>
            <TabsTrigger value="addresses" className="flex items-center gap-2">
              <MapPin className="h-4 w-4" /> Addresses
            </TabsTrigger>
            <TabsTrigger value="couriers" className="flex items-center gap-2">
              <Truck className="h-4 w-4" /> Couriers
            </TabsTrigger>
            <TabsTrigger value="products" className="flex items-center gap-2">
              <Package className="h-4 w-4" /> Shop Products
            </TabsTrigger>
          </TabsList>
          
          <TabsContent value="info" className="space-y-4">
             <Card>
              <CardHeader>
                <CardTitle>General Information</CardTitle>
                <CardDescription>
                  Basic details and configuration for your shop.
                </CardDescription>
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
                      className="flex min-h-[80px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
                      value={shopDesc}
                      onChange={(e) => setShopDesc(e.target.value)}
                      placeholder="Brief description of the shop..."
                    />
                  </div>
                  <div className="flex items-center space-x-2 pt-2">
                    <input
                      id="shopIsActive"
                      type="checkbox"
                      className="h-4 w-4 rounded border-input bg-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
                      checked={shopIsActive}
                      onChange={(e) => setShopIsActive(e.target.checked)}
                    />
                    <Label htmlFor="shopIsActive" className="text-sm font-medium leading-none cursor-pointer">
                      Shop is Active
                    </Label>
                  </div>
                  <div className="pt-4">
                    <Button type="submit" disabled={isSavingShop}>
                      {isSavingShop && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                      Save Shop Settings
                    </Button>
                  </div>
                </form>
              </CardContent>
            </Card>
          </TabsContent>

          <TabsContent value="addresses" className="space-y-4">
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-4">
                <div>
                  <CardTitle>Shop Addresses</CardTitle>
                  <CardDescription>
                    Manage physical locations where this shop operates.
                  </CardDescription>
                </div>
                
                <Sheet open={isOpen} onOpenChange={setIsOpen}>
                  <SheetTrigger asChild>
                    <Button size="sm">
                      <Plus className="mr-2 h-4 w-4" /> Add Address
                    </Button>
                  </SheetTrigger>
                  <SheetContent className="sm:max-w-md overflow-y-auto">
                    <SheetHeader className="mb-4">
                      <SheetTitle>Add Shop Address</SheetTitle>
                      <SheetDescription>
                        Create a physical branch or pick-up location for your shop.
                      </SheetDescription>
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
                          className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
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
                          className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:opacity-50"
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
                          className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:opacity-50"
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
                          className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:opacity-50"
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
                          className="flex min-h-[80px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
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
                          className="h-4 w-4 rounded border-input bg-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
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
                        <Button type="submit" disabled={isSubmitting}>
                          {isSubmitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                          Save Address
                        </Button>
                      </div>
                    </form>
                  </SheetContent>
                </Sheet>
              </CardHeader>
              <CardContent>
                <div className="rounded-md border">
                  <Table>
                    <TableHeader>
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
                          <TableCell colSpan={4} className="text-center h-24">No addresses found.</TableCell>
                        </TableRow>
                      ) : (
                        addresses.map((addr) => (
                          <TableRow key={addr.id}>
                            <TableCell className="font-medium">{addr.label}</TableCell>
                            <TableCell className="max-w-xs truncate">{addr.full_address}</TableCell>
                            <TableCell>{addr.phone || '-'}</TableCell>
                            <TableCell>
                              <Badge variant={addr.is_active ? "default" : "secondary"}>
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

          <TabsContent value="couriers" className="space-y-4">
             <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-4">
                <div>
                  <CardTitle>Configured Couriers</CardTitle>
                  <CardDescription>
                    Enable or disable shipping providers for your shop.
                  </CardDescription>
                </div>
              </CardHeader>
              <CardContent>
                <div className="rounded-md border">
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>Courier Code</TableHead>
                        <TableHead className="text-right">Status</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {couriers.length === 0 ? (
                        <TableRow>
                          <TableCell colSpan={2} className="text-center h-24">No couriers found.</TableCell>
                        </TableRow>
                      ) : (
                        couriers.map((courier) => (
                          <TableRow key={courier.code}>
                            <TableCell className="font-medium uppercase">{courier.code}</TableCell>
                            <TableCell className="text-right">
                              <Badge variant={courier.active ? "default" : "secondary"}>
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

          <TabsContent value="products" className="space-y-4">
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-4">
                <div>
                  <CardTitle>Shop Products Inventory</CardTitle>
                  <CardDescription>
                    Products associated with this specific shop location and their stock levels.
                  </CardDescription>
                </div>
                
                <Sheet open={isInventoryOpen} onOpenChange={setIsInventoryOpen}>
                  <SheetTrigger asChild>
                    <Button size="sm">
                      <Plus className="mr-2 h-4 w-4" /> Add Inventory
                    </Button>
                  </SheetTrigger>
                  <SheetContent className="sm:max-w-md overflow-y-auto">
                    <SheetHeader className="mb-4">
                      <SheetTitle>Add Inventory</SheetTitle>
                      <SheetDescription>
                        Assign stock for a product at this shop location.
                      </SheetDescription>
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
                          className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
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
              </CardHeader>
              <CardContent>
                <div className="rounded-md border">
                  <Table>
                    <TableHeader>
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
                          <TableCell colSpan={4} className="text-center h-24">No products found for this shop.</TableCell>
                        </TableRow>
                      ) : (
                        products.map((product) => (
                          <TableRow key={product.id}>
                            <TableCell className="font-medium">
                              {product.name}
                              <div className="text-xs text-muted-foreground">{product.slug}</div>
                            </TableCell>
                            <TableCell>{product.sku}</TableCell>
                            <TableCell>
                              {new Intl.NumberFormat('id-ID', {
                                style: 'currency',
                                currency: 'IDR',
                                minimumFractionDigits: 0,
                              }).format(product.price)}
                            </TableCell>
                            <TableCell className="text-right">
                              <span className="font-bold">
                                {product.inventory.available}
                              </span>
                              <span className="text-xs text-muted-foreground ml-2">
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
        </Tabs>
      </div>
    </div>
  );
}
