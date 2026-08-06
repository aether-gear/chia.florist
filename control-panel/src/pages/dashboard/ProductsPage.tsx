import { useState, useMemo } from 'react';
import { Link } from 'react-router-dom';
import { Package, Plus, Loader2, RefreshCw, MoreHorizontal, Edit, Trash2, AlertTriangle, BarChart3 } from 'lucide-react';
import { Button } from '../../components/ui/button';
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
import { useAuthMeViewModel } from '../../viewmodels/useAuthMeViewModel';
import { fetchApi } from '../../lib/api';
import ProductFormSheet from '../../components/products/ProductFormSheet';
import ProductPerformanceCharts from '../../components/products/ProductPerformanceCharts';
import { Skeleton } from '../../components/ui/skeleton';
import EmptyState from '../../components/EmptyState';
import SearchInput from '../../components/SearchInput';
import StatusBadge from '../../components/StatusBadge';
import Pagination from '../../components/Pagination';
import { DataCard, DataCardGridHeader, DataCardList } from '../../components/DataCard';

export default function ProductsPage() {
  const { isAdmin } = useAuthMeViewModel();
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
          {isAdmin && (
            <Button asChild variant="outline" size="sm" className="rounded-xl text-xs gap-1.5 border-primary/30 text-primary hover:bg-primary/5 self-start sm:self-auto">
              <Link to="/admin/analytics?tab=products">
                <BarChart3 className="h-3.5 w-3.5" /> View Inventory Analytics →
              </Link>
            </Button>
          )}
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
            <h3 className="text-xl font-bold font-display tracking-tight text-foreground">All Products</h3>
            <p className="text-muted-foreground text-sm">Manage product catalog, inventory levels, and pricing.</p>
          </div>

          <div className="flex flex-col sm:flex-row items-center justify-between gap-4 w-full">
            <SearchInput
              value={searchQuery}
              onChange={setSearchQuery}
              placeholder="Search products..."
              className="relative flex-1 max-w-sm w-full"
            />
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

          <div className="flex items-center justify-between px-1 text-xs text-muted-foreground">
            <span>Found {filteredProducts.length} products</span>
          </div>

          {/* Content */}
          <div className="flex flex-col gap-2">
            <DataCardGridHeader>
              <span className="col-span-4">Product</span>
              <span className="col-span-2">SKU</span>
              <span className="col-span-2">Status</span>
              <span className="col-span-2">Price</span>
              <span className="col-span-1">Stock</span>
              <span className="col-span-1 text-right">Actions</span>
            </DataCardGridHeader>

            <DataCardList>
              {loading ? (
                Array.from({ length: 4 }).map((_, i) => (
                  <DataCard key={`skeleton-${i}`}>
                    <div className="col-span-4 flex items-center gap-3">
                      <Skeleton className="h-9 w-9 rounded-md bg-muted animate-pulse" />
                      <Skeleton className="h-4 w-32 bg-muted animate-pulse" />
                    </div>
                    <div className="col-span-2"><Skeleton className="h-4 w-20 bg-muted animate-pulse" /></div>
                    <div className="col-span-2"><Skeleton className="h-5 w-16 bg-muted animate-pulse rounded-full" /></div>
                    <div className="col-span-2"><Skeleton className="h-4 w-24 bg-muted animate-pulse" /></div>
                    <div className="col-span-1"><Skeleton className="h-4 w-8 bg-muted animate-pulse" /></div>
                    <div className="col-span-1 text-right"><Skeleton className="h-8 w-8 ml-auto bg-muted animate-pulse rounded-xl" /></div>
                  </DataCard>
                ))
              ) : error ? (
                <EmptyState title="Failed to load products" description={error} className="py-12 border-0 bg-transparent text-destructive" />
              ) : filteredProducts.length === 0 ? (
                <EmptyState icon={<Package className="h-8 w-8 text-slate-400 mb-2 mx-auto" />} title="No products found" description="No products match your current search criteria." className="py-12 border border-dashed border-border/80 rounded-2xl bg-zinc-50/10" />
              ) : (
                filteredProducts.map((product) => (
                  <DataCard key={product.id} onClick={() => { setActiveProductSlug(product.slug); setIsProductSheetOpen(true); }}>
                    <div className="col-span-1 md:col-span-4 flex items-center gap-3">
                      <div className="h-9 w-9 overflow-hidden rounded-md border shrink-0 bg-muted flex items-center justify-center">
                        {product.banner?.thumbnail ? (
                          <img src={product.banner.thumbnail} alt={product.name} className="h-full w-full object-cover" />
                        ) : (
                          <Package className="h-4 w-4 text-muted-foreground" />
                        )}
                      </div>
                      <div className="min-w-0">
                        <h4 className="font-semibold font-display text-sm text-foreground truncate">{product.name}</h4>
                        <p className="text-xs text-muted-foreground font-mono truncate">{product.slug}</p>
                      </div>
                    </div>

                    <div className="col-span-1 md:col-span-2 font-mono text-xs text-muted-foreground">
                      <span className="md:hidden font-sans text-muted-foreground mr-1">SKU:</span>
                      {product.sku || '-'}
                    </div>

                    <div className="col-span-1 md:col-span-2">
                      <StatusBadge status={product.is_active ? 'active' : 'inactive'} className="scale-90 origin-left" />
                    </div>

                    <div className="col-span-1 md:col-span-2 font-bold text-primary">
                      {new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(product.price)}
                    </div>

                    <div className="col-span-1 md:col-span-1 text-xs">
                      <span className="md:hidden font-sans text-muted-foreground mr-1">Stock:</span>
                      <span className={product.stock < 10 ? "text-destructive font-bold" : "font-medium"}>{product.stock}</span>
                    </div>

                    <div className="col-span-1 md:col-span-1 text-right" onClick={(e) => e.stopPropagation()}>
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button variant="ghost" size="icon" className="h-8 w-8 rounded-xl">
                            <MoreHorizontal className="h-4 w-4" />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end" className="w-40 rounded-xl">
                          <DropdownMenuItem onClick={() => { setActiveProductSlug(product.slug); setIsProductSheetOpen(true); }}>
                            <Edit className="h-4 w-4 mr-2 text-muted-foreground" /> Edit Product
                          </DropdownMenuItem>
                          <DropdownMenuSeparator className="my-1" />
                          <DropdownMenuItem onClick={() => setProductToDelete(product)} className="text-destructive focus:text-destructive">
                            <Trash2 className="h-4 w-4 mr-2 text-destructive" /> Delete Product
                          </DropdownMenuItem>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </div>
                  </DataCard>
                ))
              )}
            </DataCardList>

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
