import { useState, useMemo } from 'react';
import { Package, Search, Plus, Loader2, RefreshCw, MoreHorizontal, Edit, Trash2, AlertTriangle } from 'lucide-react';
import { Button } from '../../components/ui/button';
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
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
} from '../../components/ui/dropdown-menu';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from '../../components/ui/dialog';
import { useProductsViewModel } from '../../viewmodels/useProductsViewModel';
import { fetchApi } from '../../lib/api';
import ProductFormSheet from '../../components/products/ProductFormSheet';
import LoadingState from '../../components/LoadingState';
import EmptyState from '../../components/EmptyState';
import SearchInput from '../../components/SearchInput';
import StatusBadge from '../../components/StatusBadge';

export default function ProductsPage() {
  const { data, loading, error, refresh } = useProductsViewModel();

  const [searchQuery, setSearchQuery] = useState('');

  // Product Form Sheet states
  const [isProductSheetOpen, setIsProductSheetOpen] = useState(false);
  const [activeProductSlug, setActiveProductSlug] = useState<string | undefined>(undefined);

  const [productToDelete, setProductToDelete] = useState<any | null>(null);
  const [isDeleting, setIsDeleting] = useState(false);

  const filteredProducts = useMemo(() => {
    if (!data?.products) return [];
    return data.products.filter(product =>
      product.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      (product.sku && product.sku.toLowerCase().includes(searchQuery.toLowerCase())) ||
      (product.slug && product.slug.toLowerCase().includes(searchQuery.toLowerCase()))
    );
  }, [data, searchQuery]);

  const handleDeleteProduct = async () => {
    if (!productToDelete) return;
    setIsDeleting(true);
    try {
      await fetchApi(`/products/id/${productToDelete.id}`, {
        method: 'DELETE',
      });
      setProductToDelete(null);
      refresh();
    } catch (err: any) {
      alert(err.message || 'Failed to delete product');
    } finally {
      setIsDeleting(false);
    }
  };

  if (loading) {
    return <LoadingState message="Loading products..." />;
  }

  if (error) {
    return (
      <EmptyState
        title="Failed to load products"
        description={error}
        actionLabel="Retry"
        onAction={() => refresh()}
      />
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
            <Button
              variant="outline"
              onClick={() => refresh()}
              className="flex items-center gap-1.5 border-slate-200 text-slate-600 hover:text-indigo-600 transition-colors"
            >
              <RefreshCw className="h-4 w-4" />
              Refresh
            </Button>
            <Button onClick={() => { setActiveProductSlug(undefined); setIsProductSheetOpen(true); }}>
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
              <SearchInput
                value={searchQuery}
                onChange={setSearchQuery}
                placeholder="Search products..."
              />
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
                  {filteredProducts.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={7} className="p-0">
                        <EmptyState
                          icon={<Package className="h-8 w-8 mb-2 mx-auto text-slate-400" />}
                          title="No products found"
                          description={searchQuery ? `No products match "${searchQuery}"` : "Try adding a new product to your catalog."}
                          className="flex h-32 flex-col items-center justify-center text-center p-4 gap-1.5 border-0 bg-transparent"
                        />
                      </TableCell>
                    </TableRow>
                  ) : (
                    filteredProducts.map((product) => (
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
                          <StatusBadge status={product.status} />
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
                          <DropdownMenu>
                            <DropdownMenuTrigger asChild>
                              <Button variant="ghost" size="icon" className="h-8 w-8">
                                <MoreHorizontal className="h-4 w-4" />
                              </Button>
                            </DropdownMenuTrigger>
                            <DropdownMenuContent align="end" className="bg-white border border-slate-200 shadow-md rounded-md p-1 min-w-[150px]">
                              <DropdownMenuItem
                                onClick={() => {
                                  setActiveProductSlug(product.slug);
                                  setIsProductSheetOpen(true);
                                }}
                                className="cursor-pointer flex items-center gap-2 px-2 py-1.5 text-sm hover:bg-slate-50 rounded"
                              >
                                <Edit className="h-4 w-4 text-slate-500" />
                                <span>Edit Product</span>
                              </DropdownMenuItem>
                              <DropdownMenuSeparator className="bg-slate-100 my-1" />
                              <DropdownMenuItem
                                onClick={() => setProductToDelete(product)}
                                className="cursor-pointer flex items-center gap-2 px-2 py-1.5 text-sm hover:bg-red-50 text-red-600 rounded"
                              >
                                <Trash2 className="h-4 w-4 text-red-500" />
                                <span>Delete Product</span>
                              </DropdownMenuItem>
                            </DropdownMenuContent>
                          </DropdownMenu>
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

      {/* Product Form Sheet (Add / Edit) */}
      <ProductFormSheet
        open={isProductSheetOpen}
        onOpenChange={setIsProductSheetOpen}
        productSlug={activeProductSlug}
        onSuccess={refresh}
      />

      {/* Delete Confirmation Dialog */}
      <Dialog open={!!productToDelete} onOpenChange={(open) => !open && setProductToDelete(null)}>
        <DialogContent className="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle className="flex items-center gap-2 text-red-600">
              <AlertTriangle className="h-5 w-5" />
              Confirm Deletion
            </DialogTitle>
            <DialogDescription className="py-2">
              Are you sure you want to permanently delete <strong>{productToDelete?.name}</strong>? All association with inventory and shops will be removed. This cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter className="flex gap-2">
            <Button type="button" variant="outline" onClick={() => setProductToDelete(null)} disabled={isDeleting}>
              Cancel
            </Button>
            <Button
              type="button"
              variant="destructive"
              onClick={handleDeleteProduct}
              disabled={isDeleting}
            >
              {isDeleting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Yes, Delete Product
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}

