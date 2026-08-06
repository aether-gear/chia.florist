import { useState, useEffect } from 'react';
import { Loader2, Save, UploadCloud, X, AlertCircle, Layers } from 'lucide-react';
import { Button } from '../ui/button';
import { Input } from '../ui/input';
import { Label } from '../ui/label';
import {
  Sheet,
  SheetContent,
  SheetTitle,
  SheetDescription,
} from '../ui/sheet';
import { useProductFormViewModel, type InventorySyncItem } from '../../viewmodels/useProductFormViewModel';
import { fetchApi } from '../../lib/api';

interface ProductFormSheetProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  productSlug?: string;
  onSuccess: () => void;
  inline?: boolean;
  onClose?: () => void;
}

export default function ProductFormSheet({
  open,
  onOpenChange,
  productSlug,
  onSuccess,
  inline = false,
  onClose,
}: ProductFormSheetProps) {
  const isEdit = !!productSlug;
  const {
    product,
    loading,
    saving,
    uploading,
    error,
    success,
    loadProduct,
    clearProduct,
    saveProduct,
    uploadImage,
  } = useProductFormViewModel();

  // Form states
  const [formData, setFormData] = useState({
    sku: '',
    name: '',
    description: '',
    is_available: true,
    status: 'active' as 'active' | 'inactive' | 'archived',
    price: 0,
    weight: '' as string | number,
  });

  // Local inventory list state
  const [localInventories, setLocalInventories] = useState<InventorySyncItem[]>([]);
  const [selectedShopId, setSelectedShopId] = useState('');
  const [createMore, setCreateMore] = useState(false);

  // Pending image for product creation
  const [pendingImageFile, setPendingImageFile] = useState<File | null>(null);
  const [pendingImagePreview, setPendingImagePreview] = useState<string | null>(null);

  // Shops list
  const [shops, setShops] = useState<any[]>([]);

  // Load shops and product details
  useEffect(() => {
    const loadShops = async () => {
      try {
        const res = await fetchApi('/shops');
        setShops(res.shops || []);
      } catch (err) {
        console.error('Failed to load shops', err);
      }
    };

    if (open || inline) {
      loadShops();
      if (isEdit && productSlug) {
        loadProduct(productSlug);
      } else {
        clearProduct();
        setFormData({
          sku: '',
          name: '',
          description: '',
          is_available: true,
          status: 'active',
          price: 0,
          weight: '',
        });
        setLocalInventories([]);
        setSelectedShopId('');
        setPendingImageFile(null);
        setPendingImagePreview(null);
      }
    }
  }, [open, inline, isEdit, productSlug, loadProduct, clearProduct]);

  // Sync form values on edit
  useEffect(() => {
    if (isEdit && product) {
      setFormData({
        sku: product.sku || '',
        name: product.name || '',
        description: product.description || '',
        is_available: product.is_available,
        status: product.status || 'active',
        price: product.price || 0,
        weight: product.weight !== undefined && product.weight !== null ? product.weight : '',
      });
    }
  }, [product, isEdit]);

  // Sync inventories when product and shops are available
  useEffect(() => {
    if (isEdit && product && shops.length > 0) {
      const resolved = (product.availability || []).map(avail => {
        const matchingShop = shops.find(s => s.slug === avail.name || s.name === avail.slug);
        return {
          shopId: matchingShop?.id || '',
          shopName: avail.slug || matchingShop?.name || '',
          stock: avail.stock,
          isNew: false,
          isModified: false,
          isDeleted: false,
          originalStock: avail.stock
        };
      });
      setLocalInventories(resolved);
    } else if (!isEdit) {
      setLocalInventories([]);
    }
  }, [product, shops, isEdit]);

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>
  ) => {
    const { name, value, type } = e.target as HTMLInputElement;
    const isCheckbox = type === 'checkbox';

    setFormData((prev) => ({
      ...prev,
      [name]: isCheckbox ? (e.target as HTMLInputElement).checked : value,
    }));
  };

  const handleAddShopToInventory = (shopId: string) => {
    if (!shopId) return;
    if (localInventories.some(item => item.shopId === shopId && !item.isDeleted)) return;

    const shop = shops.find(s => s.id === shopId);
    if (shop) {
      setLocalInventories(prev => {
        const existing = prev.find(item => item.shopId === shopId);
        if (existing) {
          return prev.map(item =>
            item.shopId === shopId ? { ...item, isDeleted: false, isModified: item.stock !== item.originalStock } : item
          );
        }
        return [...prev, {
          shopId: shop.id,
          shopName: shop.name,
          stock: 0,
          isNew: true,
          isModified: false,
          isDeleted: false,
          originalStock: 0
        }];
      });
      setSelectedShopId('');
    }
  };

  const handleStockChange = (shopId: string, value: string) => {
    setLocalInventories(prev =>
      prev.map(item => {
        if (item.shopId === shopId) {
          const numStock = Number(value);
          return {
            ...item,
            stock: isNaN(numStock) ? 0 : numStock,
            isModified: item.isNew ? false : numStock !== item.originalStock
          };
        }
        return item;
      })
    );
  };

  const handleRemoveShop = (shopId: string) => {
    setLocalInventories(prev =>
      prev.map(item => {
        if (item.shopId === shopId) {
          if (item.isNew) {
            return null; // delete local item
          }
          return { ...item, isDeleted: true };
        }
        return item;
      }).filter(Boolean) as InventorySyncItem[]
    );
  };

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();

    const saveValues = {
      id: product?.id,
      sku: formData.sku,
      name: formData.name,
      description: formData.description,
      is_available: formData.is_available,
      status: formData.status,
      price: Number(formData.price),
      weight: formData.weight !== '' ? Number(formData.weight) : null,
    };

    const saved = await saveProduct(saveValues, localInventories, pendingImageFile);

    if (saved) {
      setPendingImageFile(null);
      setPendingImagePreview(null);
      onSuccess();
      if (!isEdit && createMore) {
        setFormData({
          sku: '',
          name: '',
          description: '',
          is_available: true,
          status: 'active',
          price: 0,
          weight: '',
        });
        setLocalInventories([]);
        setSelectedShopId('');
      } else {
        onOpenChange(false);
      }
    }
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files.length > 0) {
      const file = e.target.files[0];
      setPendingImageFile(file);
      const objectUrl = URL.createObjectURL(file);
      setPendingImagePreview(objectUrl);
    }
  };

  // Filter out shops already added to the local inventory list
  const availableShops = shops.filter(
    s => !localInventories.some(item => item.shopId === s.id && !item.isDeleted)
  );

  const formContent = (
    <>
      {/* Header */}
      <div className={inline ? "pb-4 border-b border-border/60 flex items-center justify-between" : "px-6 py-5 border-b flex items-center justify-between shrink-0"}>
        <div>
          <h3 className="text-xl font-bold font-display text-foreground flex items-center gap-2">
            <Layers className="h-5 w-5 text-primary" />
            {isEdit ? 'Update product' : 'Add new product'}
          </h3>
          <p className="text-xs text-muted-foreground mt-1">
            {isEdit
              ? `Modify product database details and manage stock levels.`
              : 'Insert a new row to the products database.'}
          </p>
        </div>
        {inline && (
          <Button
            variant="ghost"
            size="icon"
            className="h-8 w-8 rounded-xl text-muted-foreground hover:text-foreground"
            onClick={() => (onClose ? onClose() : onOpenChange(false))}
          >
            <X className="h-4 w-4" />
          </Button>
        )}
      </div>

      {/* Scrollable Body */}
      <div className={inline ? "space-y-6" : "flex-1 overflow-y-auto px-6 py-6 space-y-6 pb-24"}>
        {error && (
          <div className="flex items-start gap-2 p-3 text-xs text-destructive bg-destructive/10 rounded-xl border border-destructive/20">
            <AlertCircle className="h-4 w-4 shrink-0 text-destructive mt-0.5" />
            <span>{error}</span>
          </div>
        )}

        {success && !isEdit && (
          <div className="flex items-start gap-2 p-3 text-xs text-primary bg-primary/10 rounded-xl border border-primary/20">
            <Layers className="h-4 w-4 shrink-0 text-primary mt-0.5" />
            <span>Product created successfully!</span>
          </div>
        )}

        {/* Basic Details Section */}
        <div className="space-y-4">
          <h3 className="text-xs font-bold uppercase tracking-wider text-muted-foreground">Basic Info</h3>

          <div className="grid grid-cols-2 gap-3">
            <div>
              <Label htmlFor="sku" className="text-xs font-semibold text-foreground">SKU / Item Code *</Label>
              <Input
                id="sku"
                name="sku"
                placeholder="e.g. ROSE-RED-01"
                value={formData.sku}
                onChange={handleChange}
                className="mt-1 text-xs rounded-xl"
                required
              />
            </div>
            <div>
              <Label htmlFor="name" className="text-xs font-semibold text-foreground">Product Name *</Label>
              <Input
                id="name"
                name="name"
                placeholder="e.g. Red Rose Bouquet"
                value={formData.name}
                onChange={handleChange}
                className="mt-1 text-xs rounded-xl"
                required
              />
            </div>
          </div>

          <div>
            <Label htmlFor="description" className="text-xs font-semibold text-foreground">Description</Label>
            <textarea
              id="description"
              name="description"
              rows={3}
              placeholder="Detailed description of the product..."
              value={formData.description}
              onChange={handleChange}
              className="mt-1 w-full rounded-xl border border-border bg-background p-2.5 text-xs text-foreground focus:outline-none focus:ring-1 focus:ring-primary/40"
            />
          </div>

          <div className="grid grid-cols-2 gap-3">
            <div>
              <Label htmlFor="price" className="text-xs font-semibold text-foreground">Price (IDR) *</Label>
              <Input
                id="price"
                name="price"
                type="number"
                min="0"
                placeholder="0"
                value={formData.price}
                onChange={handleChange}
                className="mt-1 text-xs rounded-xl"
                required
              />
            </div>
            <div>
              <Label htmlFor="weight" className="text-xs font-semibold text-foreground">Weight (grams)</Label>
              <Input
                id="weight"
                name="weight"
                type="number"
                min="0"
                placeholder="e.g. 500"
                value={formData.weight}
                onChange={handleChange}
                className="mt-1 text-xs rounded-xl"
              />
            </div>
          </div>

          <div className="grid grid-cols-2 gap-3">
            <div>
              <Label htmlFor="status" className="text-xs font-semibold text-foreground">Catalog Status</Label>
              <select
                id="status"
                name="status"
                value={formData.status}
                onChange={handleChange}
                className="mt-1 w-full rounded-xl border border-border bg-background p-2 text-xs text-foreground focus:outline-none focus:ring-1 focus:ring-primary/40"
              >
                <option value="active">Active</option>
                <option value="inactive">Inactive</option>
                <option value="archived">Archived</option>
              </select>
            </div>
            <div className="flex items-center pt-5">
              <label className="flex items-center gap-2 cursor-pointer text-xs font-medium text-foreground">
                <input
                  type="checkbox"
                  id="is_available"
                  name="is_available"
                  checked={formData.is_available}
                  onChange={handleChange}
                  className="h-4 w-4 rounded border-border text-primary accent-primary"
                />
                Available for Purchase
              </label>
            </div>
          </div>
        </div>

        {/* Media / Banner Upload Section */}
        <div className="space-y-3 pt-2">
          <h3 className="text-xs font-bold uppercase tracking-wider text-muted-foreground">Product Image</h3>
          <div className="border border-dashed border-border/80 rounded-2xl p-4 bg-muted/20 text-center">
            {isEdit && product?.banner?.thumbnail ? (
              <div className="relative inline-block">
                <img
                  src={product.banner.thumbnail}
                  alt={formData.name}
                  className="h-32 w-32 object-cover rounded-xl border border-border"
                />
                <label className="absolute bottom-1 right-1 bg-background hover:bg-muted p-1.5 rounded-lg border shadow-sm cursor-pointer">
                  <UploadCloud className="h-3.5 w-3.5 text-foreground" />
                  <input
                    type="file"
                    accept="image/*"
                    className="hidden"
                    onChange={handleFileChange}
                    disabled={uploading}
                  />
                </label>
              </div>
            ) : pendingImagePreview ? (
              <div className="relative inline-block">
                <img
                  src={pendingImagePreview}
                  alt="Pending Product Preview"
                  className="h-32 w-32 object-cover rounded-xl border border-border"
                />
                <div className="absolute -top-2 -right-2 flex gap-1">
                  <Button
                    type="button"
                    variant="destructive"
                    size="icon"
                    className="h-6 w-6 rounded-full"
                    onClick={() => {
                      setPendingImageFile(null);
                      setPendingImagePreview(null);
                    }}
                  >
                    <X className="h-3 w-3" />
                  </Button>
                </div>
                <label className="absolute bottom-1 right-1 bg-background hover:bg-muted p-1.5 rounded-lg border shadow-sm cursor-pointer">
                  <UploadCloud className="h-3.5 w-3.5 text-foreground" />
                  <input
                    type="file"
                    accept="image/*"
                    className="hidden"
                    onChange={handleFileChange}
                  />
                </label>
              </div>
            ) : (
              <label className="flex flex-col items-center justify-center cursor-pointer py-4">
                <UploadCloud className="h-8 w-8 text-muted-foreground mb-2" />
                <span className="text-xs font-medium text-foreground">Click to upload product image</span>
                <span className="text-[10px] text-muted-foreground">PNG, JPG up to 5MB</span>
                <input
                  type="file"
                  accept="image/*"
                  className="hidden"
                  onChange={handleFileChange}
                  disabled={uploading}
                />
              </label>
            )}
            {uploading && <p className="text-xs text-primary mt-2">Uploading image...</p>}
          </div>
        </div>

        {/* Inventory Stock Assignment */}
        <div className="space-y-4 pt-2">
          <h3 className="text-xs font-bold uppercase tracking-wider text-muted-foreground">Shop Inventories</h3>

          {availableShops.length > 0 && (
            <div className="flex gap-2">
              <select
                value={selectedShopId}
                onChange={(e) => {
                  setSelectedShopId(e.target.value);
                  handleAddShopToInventory(e.target.value);
                }}
                className="flex-1 rounded-xl border border-border bg-background p-2 text-xs text-foreground focus:outline-none focus:ring-1 focus:ring-primary/40"
              >
                <option value="">+ Assign shop location...</option>
                {availableShops.map((s) => (
                  <option key={s.id} value={s.id}>{s.name}</option>
                ))}
              </select>
            </div>
          )}

          <div className="space-y-2">
            {localInventories.filter(i => !i.isDeleted).length === 0 ? (
              <p className="text-xs text-muted-foreground italic">No stock assigned to any shop branch.</p>
            ) : (
              localInventories
                .filter(i => !i.isDeleted)
                .map((item) => (
                  <div key={item.shopId} className="flex items-center justify-between gap-3 p-2.5 rounded-xl border border-border bg-muted/20">
                    <div className="min-w-0">
                      <span className="text-xs font-semibold text-foreground truncate block">{item.shopName}</span>
                      {item.isNew && <span className="text-[9px] text-emerald-600 font-bold uppercase">New</span>}
                      {!item.isNew && item.isModified && <span className="text-[9px] text-amber-600 font-bold uppercase">Modified</span>}
                    </div>
                    <div className="flex items-center gap-2 shrink-0">
                      <Input
                        type="number"
                        min="0"
                        value={item.stock}
                        onChange={(e) => handleStockChange(item.shopId, e.target.value)}
                        className="w-20 h-8 text-xs rounded-lg text-right"
                      />
                      <Button
                        type="button"
                        variant="ghost"
                        size="icon"
                        onClick={() => handleRemoveShop(item.shopId)}
                        className="h-7 w-7 text-destructive hover:text-destructive hover:bg-destructive/10 rounded-lg"
                      >
                        <X className="h-3.5 w-3.5" />
                      </Button>
                    </div>
                  </div>
                ))
            )}
          </div>
        </div>
      </div>

      {/* Action Footer */}
      <div className={inline ? "pt-4 border-t border-border/60 flex items-center justify-between" : "px-6 py-4 border-t bg-muted/10 flex items-center justify-between shrink-0"}>
        <div className="flex items-center gap-2">
          {!isEdit && (
            <label className="flex items-center gap-2 cursor-pointer text-xs text-muted-foreground select-none">
              <input
                type="checkbox"
                checked={createMore}
                onChange={(e) => setCreateMore(e.target.checked)}
                className="h-3.5 w-3.5 rounded border-border text-primary accent-primary"
              />
              <span>Create more</span>
            </label>
          )}
        </div>
        <div className="flex items-center gap-2">
          <Button
            type="button"
            variant="outline"
            size="sm"
            className="text-xs rounded-xl"
            onClick={() => (onClose ? onClose() : onOpenChange(false))}
            disabled={saving}
          >
            Cancel
          </Button>
          <Button
            type="submit"
            size="sm"
            className="text-xs rounded-xl font-medium flex items-center gap-1.5 bg-primary hover:bg-primary/90 text-primary-foreground"
            onClick={handleSave}
            disabled={saving || loading}
          >
            {saving ? <Loader2 className="h-3 w-3 animate-spin" /> : <Save className="h-3.5 w-3.5" />}
            {isEdit ? 'Save Product' : 'Create Product'}
          </Button>
        </div>
      </div>
    </>
  );

  if (inline) {
    return (
      <div className="border border-border/60 rounded-2xl p-6 bg-background space-y-6 shadow-none animate-in fade-in slide-in-from-right-2 duration-200">
        {formContent}
      </div>
    );
  }

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent className="w-full sm:max-w-none md:w-[45vw] md:min-w-[45vw] p-0 flex flex-col h-full border-l border-border/60 bg-background shadow-2xl">
        {formContent}
      </SheetContent>
    </Sheet>
  );
}
