import { useState, useMemo } from 'react';
import { Package, Plus, Loader2, RefreshCw, MoreHorizontal, Edit, Trash2, AlertTriangle } from 'lucide-react';
import { Button } from '../../components/ui/button';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '../../components/ui/table';
// Removed Card imports since sections are now borderless and backgroundless
import {
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
} from '../../components/ui/dropdown-menu';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from '../../components/ui/dialog';
import { useProductsViewModel } from '../../viewmodels/useProductsViewModel';
import { useProductStatsViewModel } from '../../viewmodels/useProductStatsViewModel';
import { fetchApi } from '../../lib/api';
import ProductFormSheet from '../../components/products/ProductFormSheet';
import ProductPerformanceCharts from '../../components/products/ProductPerformanceCharts';
import { Skeleton } from '../../components/ui/skeleton';
import EmptyState from '../../components/EmptyState';
import SearchInput from '../../components/SearchInput';
import StatusBadge from '../../components/StatusBadge';
import Pagination from '../../components/Pagination';

export default function ProductsPage() {
  const { data, loading, error, refresh, page, limit, setPage } = useProductsViewModel();
  const { data: statsData, loading: statsLoading, error: statsError, refresh: refreshStats } = useProductStatsViewModel();

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
      refreshStats();
    } catch (err: any) {
      alert(err.message || 'Failed to delete product');
    } finally {
      setIsDeleting(false);
    }
  };



  return (
    <div className="flex-col md:flex">
      <div className="flex-1 space-y-12 p-6 sm:p-8 lg:p-12 animate-in fade-in duration-300">
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <h2 className="text-3xl font-bold font-display tracking-tight text-foreground">Products</h2>
            <p className="text-muted-foreground text-sm">
              Manage your product catalog and inventory
            </p>
          </div>
        </div>

        {/* Performance Analytics Section */}
        <div className="space-y-4">
          <h3 className="text-xl font-bold font-display tracking-tight text-foreground">Performance Analytics</h3>

          {statsLoading ? (
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
              {Array.from({ length: 3 }).map((_, i) => (
                <div key={`stats-skeleton-${i}`} className="h-[320px] flex flex-col justify-between space-y-4">
                  <div className="space-y-2 pb-2 border-b border-border/60">
                    <Skeleton className="h-5 w-32 animate-pulse bg-muted" />
                    <Skeleton className="h-3.5 w-48 animate-pulse bg-muted" />
                  </div>
                  <Skeleton className="h-[180px] w-full rounded-xl animate-pulse bg-muted my-4" />
                </div>
              ))}
            </div>
          ) : statsError ? (
            <div className="text-sm text-destructive bg-destructive/10 p-4 rounded-xl border border-destructive/20 font-sans">
              Failed to load performance metrics: {statsError}
            </div>
          ) : (
            statsData?.stats && <ProductPerformanceCharts stats={statsData.stats} />
          )}
        </div>

        <div className="space-y-6">
          <div className="pb-4 border-b border-border/60">
            <h3 className="font-bold font-display tracking-tight text-lg text-foreground">All Products</h3>
            <p className="text-muted-foreground text-sm">
              You have {data?.total || 0} total products in your catalog.
            </p>
          </div>
          <div>
            <div className="mb-4 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
              {/* Left Side: Filter and Search */}
              <div className="flex flex-col sm:flex-row items-center gap-4 w-full sm:w-auto">
                <SearchInput
                  value={searchQuery}
                  onChange={setSearchQuery}
                  placeholder="Search products..."
                />
              </div>

              {/* Right Side: Adding and Refresh */}
              <div className="flex items-center gap-2 justify-end w-full sm:w-auto">
                <Button
                  variant="outline"
                  onClick={() => { refresh(); refreshStats(); }}
                  disabled={loading || statsLoading}
                  className="flex items-center gap-1.5 border-border text-foreground hover:text-primary hover:bg-primary/5 rounded-xl transition-colors"
                >
                  <RefreshCw className={`h-4 w-4 ${(loading || statsLoading) ? 'animate-spin' : ''}`} />
                  Refresh
                </Button>
                <Button className="bg-primary hover:bg-primary/90 text-primary-foreground rounded-xl" onClick={() => { setActiveProductSlug(undefined); setIsProductSheetOpen(true); }}>
                  <Plus className="mr-2 h-4 w-4" /> Add Product
                </Button>
              </div>
            </div>

            <div className="rounded-2xl border border-border overflow-hidden">
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
                  {loading ? (
                    Array.from({ length: 5 }).map((_, i) => (
                      <TableRow key={`skeleton-${i}`}>
                        <TableCell><Skeleton className="h-10 w-10 rounded-md bg-muted animate-pulse" /></TableCell>
                        <TableCell>
                          <Skeleton className="h-5 w-32 animate-pulse bg-muted mb-1.5" />
                          <Skeleton className="h-3.5 w-24 animate-pulse bg-muted" />
                        </TableCell>
                        <TableCell><Skeleton className="h-5 w-20 animate-pulse bg-muted" /></TableCell>
                        <TableCell><Skeleton className="h-5 w-16 animate-pulse bg-muted" /></TableCell>
                        <TableCell><Skeleton className="h-5 w-24 animate-pulse bg-muted" /></TableCell>
                        <TableCell className="text-right"><Skeleton className="h-5 w-12 ml-auto animate-pulse bg-muted" /></TableCell>
                        <TableCell><Skeleton className="h-8 w-8 rounded-xl ml-auto animate-pulse bg-muted" /></TableCell>
                      </TableRow>
                    ))
                  ) : error ? (
                    <TableRow>
                      <TableCell colSpan={7} className="p-0">
                        <EmptyState
                          title="Failed to load products"
                          description={error}
                          actionLabel="Retry"
                          onAction={() => refresh()}
                          className="flex h-32 flex-col items-center justify-center text-center p-4 gap-2 border-0 bg-transparent text-destructive"
                        />
                      </TableCell>
                    </TableRow>
                  ) : filteredProducts.length === 0 ? (
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
                            <DropdownMenuContent align="end" className="min-w-[150px] p-1">
                              <DropdownMenuItem
                                onClick={() => {
                                  setActiveProductSlug(product.slug);
                                  setIsProductSheetOpen(true);
                                }}
                                className="cursor-pointer flex items-center gap-2 px-2.5 py-1.5 text-sm rounded-lg hover:bg-muted"
                              >
                                <Edit className="h-4 w-4 text-muted-foreground" />
                                <span>Edit Product</span>
                              </DropdownMenuItem>
                              <DropdownMenuSeparator className="my-1" />
                              <DropdownMenuItem
                                onClick={() => setProductToDelete(product)}
                                className="cursor-pointer flex items-center gap-2 px-2.5 py-1.5 text-sm rounded-lg hover:bg-destructive/10 text-destructive focus:bg-destructive/10 focus:text-destructive"
                              >
                                <Trash2 className="h-4 w-4 text-destructive" />
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

            <Pagination
              currentPage={page}
              totalPages={Math.ceil((data?.total || 0) / limit)}
              totalItems={data?.total || 0}
              limit={limit}
              onPageChange={setPage}
              itemNamePlural="products"
            />
          </div>
        </div>
      </div>

      {/* Product Form Sheet (Add / Edit) */}
      <ProductFormSheet
        open={isProductSheetOpen}
        onOpenChange={setIsProductSheetOpen}
        productSlug={activeProductSlug}
        onSuccess={() => { refresh(); refreshStats(); }}
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
