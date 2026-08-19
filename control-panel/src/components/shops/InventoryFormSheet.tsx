import { useState, useEffect } from 'react';
import { Loader2, Save, Plus, Package, AlertCircle } from 'lucide-react';
import { Button } from '../ui/button';
import { Input } from '../ui/input';
import { Label } from '../ui/label';
import { Badge } from '../ui/badge';
import {
  Sheet,
  SheetContent,
} from '../ui/sheet';
import { fetchApi } from '../../lib/api';

interface InventoryFormSheetProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  shopId: string;
  shopName?: string;
  product?: any | null; // If passed -> Edit mode; if null/omitted -> Add mode
  existingProducts?: any[]; // List of products currently assigned to shop (for Add mode filtering)
  onSuccess: () => void;
}

export default function InventoryFormSheet({
  open,
  onOpenChange,
  shopId,
  shopName,
  product,
  existingProducts = [],
  onSuccess,
}: InventoryFormSheetProps) {
  const isEdit = Boolean(product);

  // Common form states
  const [stock, setStock] = useState<string>('');
  const [submitting, setSubmitting] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  // Add Mode states
  const [catalogProducts, setCatalogProducts] = useState<any[]>([]);
  const [selectedProductId, setSelectedProductId] = useState<string>('');
  const [loadingCatalog, setLoadingCatalog] = useState<boolean>(false);

  // Reset form / Sync states when sheet opens
  useEffect(() => {
    if (open) {
      setError(null);
      if (isEdit && product) {
        setStock(String(product.inventory?.total_stock ?? 0));
        setSelectedProductId('');
      } else {
        setStock('');
        setSelectedProductId('');
        // Load catalog products for Add Mode
        const loadCatalog = async () => {
          try {
            setLoadingCatalog(true);
            const res = await fetchApi('/products?page=1&limit=100');
            setCatalogProducts(res.products || []);
          } catch (err: any) {
            console.error('Failed to load catalog products', err);
            setError(err.message || 'Failed to load catalog products');
          } finally {
            setLoadingCatalog(false);
          }
        };
        loadCatalog();
      }
    }
  }, [open, isEdit, product]);

  // Available catalog products for assignment (excluding products already in shop)
  const availableCatalogProducts = catalogProducts.filter(
    (p) => !existingProducts.some((ep) => ep.id === p.id)
  );

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!shopId || stock === '') return;

    if (!isEdit && !selectedProductId) {
      setError('Please select a product to assign to this shop.');
      return;
    }

    setSubmitting(true);
    setError(null);

    try {
      if (isEdit && product) {
        // Edit Mode: PUT /shops/{shopID}/products/{productID}/inventories
        await fetchApi(`/shops/${shopId}/products/${product.id}/inventories`, {
          method: 'PUT',
          body: JSON.stringify({ stock: Number(stock) }),
        });
      } else {
        // Add Mode: POST /shops/{shopID}/products/{productID}/inventories
        await fetchApi(`/shops/${shopId}/products/${selectedProductId}/inventories`, {
          method: 'POST',
          body: JSON.stringify({ stock: Number(stock) }),
        });
      }
      onSuccess();
      onOpenChange(false);
    } catch (err: any) {
      console.error(err);
      setError(err.message || (isEdit ? 'Failed to update inventory' : 'Failed to add inventory'));
    } finally {
      setSubmitting(false);
    }
  };

  const formContent = (
    <>
      {/* Header */}
      <div className="px-6 py-5 border-b flex items-center justify-between shrink-0">
        <div>
          <h3 className="text-xl font-bold font-display text-foreground flex items-center gap-2">
            <Package className="h-5 w-5 text-primary" />
            {isEdit ? 'Update Product Inventory' : 'Add Product Inventory'}
          </h3>
          <p className="text-xs text-muted-foreground mt-1">
            {isEdit ? (
              <>Adjust stock levels for <strong className="text-foreground">{product?.name}</strong>{shopName ? ` at ${shopName}` : ''}.</>
            ) : (
              <>Assign stock for a catalog product{shopName ? ` at ${shopName}` : ''}.</>
            )}
          </p>
        </div>
      </div>

      {/* Scrollable Body */}
      <div className="flex-1 overflow-y-auto px-6 py-6 space-y-6 pb-24">
        {error && (
          <div className="flex items-start gap-2 p-3 text-xs text-destructive bg-destructive/10 rounded-xl border border-destructive/20 font-sans">
            <AlertCircle className="h-4 w-4 shrink-0 text-destructive mt-0.5" />
            <span>{error}</span>
          </div>
        )}

        {/* Product Section */}
        {isEdit ? (
          /* Edit Mode: Product Information Card */
          <div className="space-y-3">
            <h3 className="text-xs font-bold uppercase tracking-wider text-muted-foreground">Product Information</h3>
            <div className="p-4 rounded-xl border border-border bg-muted/20 space-y-3">
              <div className="flex items-center justify-between text-xs">
                <span className="text-muted-foreground font-medium">Product Name</span>
                <span className="font-semibold text-foreground">{product?.name || '-'}</span>
              </div>
              <div className="flex items-center justify-between text-xs">
                <span className="text-muted-foreground font-medium">SKU</span>
                <span className="font-mono text-foreground font-semibold">{product?.sku || '-'}</span>
              </div>
              <div className="flex items-center justify-between text-xs">
                <span className="text-muted-foreground font-medium">Price</span>
                <span className="font-bold text-primary">
                  {product?.price
                    ? new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(product.price)
                    : '-'}
                </span>
              </div>
              <div className="flex items-center justify-between text-xs">
                <span className="text-muted-foreground font-medium">Status</span>
                <Badge
                  variant={product?.status === 'active' ? 'default' : 'secondary'}
                  className={
                    product?.status === 'active'
                      ? 'bg-primary/10 text-primary border-0 rounded-lg scale-90 origin-right'
                      : 'bg-muted text-muted-foreground border-0 rounded-lg scale-90 origin-right'
                  }
                >
                  {product?.status || 'Active'}
                </Badge>
              </div>
            </div>
          </div>
        ) : (
          /* Add Mode: Select Product Dropdown */
          <div className="space-y-3">
            <h3 className="text-xs font-bold uppercase tracking-wider text-muted-foreground">Select Product</h3>
            <div className="space-y-2">
              <Label htmlFor="catalogProductSelect" className="text-xs font-semibold text-foreground">
                Catalog Product *
              </Label>
              <select
                id="catalogProductSelect"
                value={selectedProductId}
                onChange={(e) => setSelectedProductId(e.target.value)}
                disabled={loadingCatalog || availableCatalogProducts.length === 0}
                className="w-full rounded-xl border border-border bg-background p-2.5 text-xs text-foreground focus:outline-none focus:ring-1 focus:ring-primary/40 disabled:opacity-50"
                required
              >
                <option value="">
                  {loadingCatalog
                    ? 'Loading products catalog...'
                    : availableCatalogProducts.length === 0
                    ? 'No remaining catalog products to add'
                    : 'Select a product from catalog...'}
                </option>
                {availableCatalogProducts.map((p) => (
                  <option key={p.id} value={p.id}>
                    {p.name} ({p.sku || 'No SKU'}) - {new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(p.price || 0)}
                  </option>
                ))}
              </select>
              {availableCatalogProducts.length === 0 && !loadingCatalog && (
                <p className="text-[11px] text-muted-foreground italic">
                  All available products in your catalog are already listed for this shop location.
                </p>
              )}
            </div>
          </div>
        )}

        {/* Stock Allocation Section */}
        <div className="space-y-4 pt-2">
          <h3 className="text-xs font-bold uppercase tracking-wider text-muted-foreground">Stock Allocation</h3>
          
          {isEdit && (
            <div className="p-4 rounded-xl border border-border bg-muted/20 space-y-3">
              <div className="flex items-center justify-between text-xs">
                <span className="text-muted-foreground font-medium">Current Available Stock</span>
                <span className="font-bold text-foreground text-sm">{product?.inventory?.available ?? 0}</span>
              </div>
              <div className="flex items-center justify-between text-xs">
                <span className="text-muted-foreground font-medium">Reserved Stock</span>
                <span className="font-mono text-muted-foreground">{product?.inventory?.reserved_stock ?? 0}</span>
              </div>
            </div>
          )}

          <div className="space-y-2">
            <Label htmlFor="inventoryStockInput" className="text-xs font-semibold text-foreground">
              {isEdit ? 'Total Stock Quantity *' : 'Initial Stock Quantity *'}
            </Label>
            <Input
              id="inventoryStockInput"
              name="stock"
              type="number"
              min="0"
              placeholder="e.g. 50"
              value={stock}
              onChange={(e) => setStock(e.target.value)}
              className="mt-1 text-xs rounded-xl"
              required
            />
            <p className="text-[11px] text-muted-foreground">
              {isEdit
                ? `Available stock will automatically calculate as Total Stock minus Reserved Stock (${product?.inventory?.reserved_stock ?? 0}).`
                : 'Enter the initial units to assign to this shop location.'}
            </p>
          </div>
        </div>
      </div>

      {/* Action Footer */}
      <div className="px-6 py-4 border-t bg-muted/10 flex items-center justify-end gap-2 shrink-0">
        <Button
          type="button"
          variant="outline"
          size="sm"
          className="text-xs rounded-xl"
          onClick={() => onOpenChange(false)}
          disabled={submitting}
        >
          Cancel
        </Button>
        <Button
          type="submit"
          size="sm"
          className="text-xs rounded-xl font-medium flex items-center gap-1.5 bg-primary hover:bg-primary/90 text-primary-foreground"
          onClick={handleSubmit}
          disabled={submitting}
        >
          {submitting ? (
            <Loader2 className="h-3.5 w-3.5 animate-spin" />
          ) : isEdit ? (
            <Save className="h-3.5 w-3.5" />
          ) : (
            <Plus className="h-3.5 w-3.5" />
          )}
          {isEdit ? 'Save Inventory' : 'Add Stock'}
        </Button>
      </div>
    </>
  );

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent className="w-full sm:max-w-none md:w-[45vw] md:min-w-[45vw] p-0 flex flex-col h-full border-l border-border/60 bg-background shadow-2xl">
        {formContent}
      </SheetContent>
    </Sheet>
  );
}
