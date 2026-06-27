import { useState, useEffect } from 'react';
import { Package, Search, Plus, Loader2 } from 'lucide-react';
import { Button } from '../../components/ui/button';
import { Input } from '../../components/ui/input';
import { Badge } from '../../components/ui/badge';
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
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetClose,
} from '../../components/ui/sheet';
import { useProductsViewModel } from '../../viewmodels/useProductsViewModel';
import { useNavigate } from 'react-router-dom';
import { fetchApi } from '../../lib/api';

export default function ProductsPage() {
  const { data, loading, error, refresh } = useProductsViewModel();
  const navigate = useNavigate();

  // Sheet form states
  const [selectedProduct, setSelectedProduct] = useState<any | null>(null);
  const [isAddInventoryOpen, setIsAddInventoryOpen] = useState(false);
  const [shops, setShops] = useState<any[]>([]);
  const [selectedShopId, setSelectedShopId] = useState('');
  const [stock, setStock] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [formError, setFormError] = useState<string | null>(null);

  // Load shops on Sheet open
  useEffect(() => {
    if (isAddInventoryOpen) {
      const loadShops = async () => {
        try {
          const res = await fetchApi('/shops');
          setShops(res.shops || []);
        } catch (err: any) {
          console.error('Failed to load shops', err);
        }
      };
      loadShops();
    }
  }, [isAddInventoryOpen]);

  const handleAddInventory = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedProduct || !selectedShopId || !stock) return;

    setIsSubmitting(true);
    setFormError(null);

    try {
      await fetchApi(`/shops/${selectedShopId}/products/${selectedProduct.id}/inventories`, {
        method: 'POST',
        body: JSON.stringify({ stock: Number(stock) }),
      });
      setIsAddInventoryOpen(false);
      setSelectedShopId('');
      setStock('');
      refresh();
    } catch (err: any) {
      setFormError(err.message || 'Failed to add inventory');
    } finally {
      setIsSubmitting(false);
    }
  };

  if (loading) {
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
            <h2 className="text-3xl font-bold tracking-tight">Products</h2>
            <p className="text-muted-foreground">
              Manage your product catalog and inventory
            </p>
          </div>
          <div className="flex items-center space-x-2">
            <Button onClick={() => navigate('/products/create')}>
              <Plus className="mr-2 h-4 w-4" /> Add Product
            </Button>
          </div>
        </div>

        <Card>
          <CardHeader>
            <CardTitle>All Products</CardTitle>
            <CardDescription>
              You have {data?.total || 0} total products in your catalog.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="mb-4 flex items-center gap-4">
              <div className="relative flex-1 max-w-sm">
                <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
                <Input
                  type="search"
                  placeholder="Search products..."
                  className="pl-8"
                />
              </div>
            </div>

            <div className="rounded-md border">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead className="w-[80px]">Image</TableHead>
                    <TableHead>Product</TableHead>
                    <TableHead>SKU</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead>Price</TableHead>
                    <TableHead className="text-right">Stock</TableHead>
                    <TableHead className="w-[150px]"></TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {!data?.products || data.products.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={7} className="h-24 text-center">
                        No products found.
                      </TableCell>
                    </TableRow>
                  ) : (
                    data.products.map((product) => (
                      <TableRow key={product.id}>
                        <TableCell>
                          <div className="h-10 w-10 overflow-hidden rounded-md border">
                            {product.banner?.thumbnail ? (
                              <img
                                src={product.banner.thumbnail}
                                alt={product.name}
                                className="h-full w-full object-cover"
                              />
                            ) : (
                              <div className="flex h-full w-full items-center justify-center bg-muted">
                                <Package className="h-4 w-4 text-muted-foreground" />
                              </div>
                            )}
                          </div>
                        </TableCell>
                        <TableCell className="font-medium">
                          {product.name}
                          <div className="text-xs text-muted-foreground">
                            {product.slug}
                          </div>
                        </TableCell>
                        <TableCell>{product.sku}</TableCell>
                        <TableCell>
                          <Badge
                            variant={
                              product.status === 'active'
                                ? 'default'
                                : product.status === 'archived'
                                ? 'secondary'
                                : 'destructive'
                            }
                          >
                            {product.status}
                          </Badge>
                        </TableCell>
                        <TableCell>
                          {new Intl.NumberFormat('id-ID', {
                            style: 'currency',
                            currency: 'IDR',
                            minimumFractionDigits: 0,
                          }).format(product.price)}
                        </TableCell>
                        <TableCell className="text-right">
                          <span className={product.stock < 10 ? 'text-destructive font-medium' : ''}>
                            {product.stock}
                          </span>
                        </TableCell>
                        <TableCell className="text-right">
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => {
                              setSelectedProduct(product);
                              setIsAddInventoryOpen(true);
                            }}
                          >
                            <Plus className="mr-1 h-3.5 w-3.5" /> Add Inventory
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

      {/* Add Inventory Sheet */}
      <Sheet open={isAddInventoryOpen} onOpenChange={setIsAddInventoryOpen}>
        <SheetContent className="sm:max-w-md overflow-y-auto">
          <SheetHeader className="mb-4">
            <SheetTitle>Add Inventory</SheetTitle>
            <SheetDescription>
              Assign new inventory stock to a shop location for <strong>{selectedProduct?.name}</strong>.
            </SheetDescription>
          </SheetHeader>

          {formError && (
            <div className="p-3 text-sm text-red-500 bg-red-50 rounded-md border border-red-100 mb-4">
              {formError}
            </div>
          )}

          <form onSubmit={handleAddInventory} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="shop">Select Shop</Label>
              <select
                id="shop"
                className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
                value={selectedShopId}
                onChange={(e) => setSelectedShopId(e.target.value)}
                required
              >
                <option value="">Select Shop Branch</option>
                {shops.map((s) => (
                  <option key={s.id} value={s.id}>
                    {s.name}
                  </option>
                ))}
              </select>
            </div>

            <div className="space-y-2">
              <Label htmlFor="stock">Initial Stock</Label>
              <Input
                id="stock"
                type="number"
                min="0"
                placeholder="e.g. 50"
                value={stock}
                onChange={(e) => setStock(e.target.value)}
                required
              />
            </div>

            <div className="pt-4 flex justify-end gap-2">
              <SheetClose asChild>
                <Button type="button" variant="outline">
                  Cancel
                </Button>
              </SheetClose>
              <Button type="submit" disabled={isSubmitting}>
                {isSubmitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                Add Stock
              </Button>
            </div>
          </form>
        </SheetContent>
      </Sheet>
    </div>
  );
}

