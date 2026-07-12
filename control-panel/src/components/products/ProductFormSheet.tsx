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
}

export default function ProductFormSheet({
  open,
  onOpenChange,
  productSlug,
  onSuccess,
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

    if (open) {
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
      }
    }
  }, [open, isEdit, productSlug, loadProduct, clearProduct]);

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

    const saved = await saveProduct(saveValues, localInventories);

    if (saved) {
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

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files.length > 0) {
      const file = e.target.files[0];
      await uploadImage(file);
      onSuccess();
    }
  };

  // Filter out shops already added to the local inventory list
  const availableShops = shops.filter(
    s => !localInventories.some(item => item.shopId === s.id && !item.isDeleted)
  );

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent className="sm:max-w-xl md:max-w-2xl p-0 flex flex-col h-full border-l bg-white shadow-2xl">
        {/* Header */}
        <div className="px-6 py-5 border-b flex items-center justify-between shrink-0">
          <div>
            <SheetTitle className="text-xl font-bold text-slate-900 flex items-center gap-2">
              <Layers className="h-5 w-5 text-indigo-600" />
              {isEdit ? 'Update product' : 'Add new product'}
            </SheetTitle>
            <SheetDescription className="text-xs text-slate-500 mt-1">
              {isEdit
                ? `Modify product database details and manage stock levels.`
                : 'Insert a new row to the products database.'}
            </SheetDescription>
          </div>
          <button
            onClick={() => onOpenChange(false)}
            className="rounded-md p-1.5 text-slate-400 hover:text-slate-500 hover:bg-slate-100 transition-colors"
          >
            <X className="h-5 w-5" />
          </button>
        </div>

        {/* Scrollable Body */}
        <div className="flex-1 overflow-y-auto px-6 py-6 space-y-6 pb-24">
          {error && (
            <div className="flex items-start gap-2 p-3 text-xs text-red-600 bg-red-50 rounded-lg border border-red-100">
              <AlertCircle className="h-4 w-4 shrink-0 text-red-500 mt-0.5" />
              <span>{error}</span>
            </div>
          )}

          {success && !isEdit && (
            <div className="p-3 text-xs text-emerald-600 bg-emerald-50 rounded-lg border border-emerald-100">
              Product successfully created.
            </div>
          )}

          {loading ? (
            <div className="flex h-48 items-center justify-center">
              <Loader2 className="h-8 w-8 animate-spin text-indigo-600" />
            </div>
          ) : (
            <form onSubmit={handleSave} className="space-y-4 division-y">
              {/* SKU */}
              <div className="grid grid-cols-3 gap-4 items-start py-4 border-b border-slate-100">
                <div className="col-span-1">
                  <Label htmlFor="sku" className="text-sm font-bold text-slate-800">SKU</Label>
                </div>
                <div className="col-span-2">
                  <Input
                    id="sku"
                    name="sku"
                    placeholder="e.g. PRO-ROSE-01"
                    value={formData.sku}
                    onChange={handleChange}
                    required
                    className="border-slate-200 focus-visible:ring-indigo-500"
                  />
                </div>
              </div>

              {/* NAME */}
              <div className="grid grid-cols-3 gap-4 items-start py-4 border-b border-slate-100">
                <div className="col-span-1">
                  <Label htmlFor="name" className="text-sm font-bold text-slate-800">Name</Label>
                </div>
                <div className="col-span-2">
                  <Input
                    id="name"
                    name="name"
                    placeholder="e.g. Red Rose Bouquet"
                    value={formData.name}
                    onChange={handleChange}
                    required
                    className="border-slate-200 focus-visible:ring-indigo-500"
                  />
                </div>
              </div>

              {/* DESCRIPTION */}
              <div className="grid grid-cols-3 gap-4 items-start py-4 border-b border-slate-100">
                <div className="col-span-1">
                  <Label htmlFor="description" className="text-sm font-bold text-slate-800">Description</Label>
                </div>
                <div className="col-span-2">
                  <textarea
                    id="description"
                    name="description"
                    rows={3}
                    placeholder="NULL (optional description...)"
                    value={formData.description}
                    onChange={handleChange}
                    className="flex min-h-[80px] w-full rounded-md border border-slate-200 bg-transparent px-3 py-2 text-sm shadow-sm placeholder:text-slate-400 focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-indigo-500"
                  />
                </div>
              </div>

              {/* PRICE */}
              <div className="grid grid-cols-3 gap-4 items-start py-4 border-b border-slate-100">
                <div className="col-span-1">
                  <Label htmlFor="price" className="text-sm font-bold text-slate-800">Price</Label>
                </div>
                <div className="col-span-2">
                  <div className="relative flex items-center">
                    <span className="absolute left-3 text-sm text-slate-400 font-semibold select-none">Rp.</span>
                    <Input
                      id="price"
                      name="price"
                      type="number"
                      min="0"
                      placeholder="0"
                      value={formData.price}
                      onChange={handleChange}
                      required
                      className="pl-9 pr-9 border-slate-200 focus-visible:ring-indigo-500"
                    />
                    <span className="absolute right-3 text-sm text-slate-400 font-semibold select-none">,00</span>
                  </div>
                </div>
              </div>

              {/* WEIGHT */}
              <div className="grid grid-cols-3 gap-4 items-start py-4 border-b border-slate-100">
                <div className="col-span-1">
                  <Label htmlFor="weight" className="text-sm font-bold text-slate-800">Weight</Label>
                </div>
                <div className="col-span-2">
                  <div className="relative flex items-center">
                    <Input
                      id="weight"
                      name="weight"
                      type="number"
                      step="0.01"
                      min="0"
                      placeholder="NULL"
                      value={formData.weight}
                      onChange={handleChange}
                      className="pr-14 border-slate-200 focus-visible:ring-indigo-500"
                    />
                    <span className="absolute right-3 text-xs text-slate-400 font-medium select-none">grams</span>
                  </div>
                </div>
              </div>

              {/* STATUS */}
              <div className="grid grid-cols-3 gap-4 items-start py-4 border-b border-slate-100">
                <div className="col-span-1">
                  <Label htmlFor="status" className="text-sm font-bold text-slate-800">Status</Label>
                </div>
                <div className="col-span-2">
                  <select
                    id="status"
                    name="status"
                    value={formData.status}
                    onChange={handleChange}
                    className="flex h-10 w-full rounded-md border border-slate-200 bg-white px-3 py-2 text-sm focus:outline-none focus:ring-1 focus:ring-indigo-500"
                  >
                    <option value="active">Active</option>
                    <option value="inactive">Inactive</option>
                    <option value="archived">Archived</option>
                  </select>
                </div>
              </div>

              {/* IS AVAILABLE */}
              <div className="grid grid-cols-3 gap-4 items-start py-4 border-b border-slate-100">
                <div className="col-span-1">
                  <Label htmlFor="is_available" className="text-sm font-bold text-slate-800">Available for sale</Label>
                </div>
                <div className="col-span-2 flex items-center h-10">
                  <input
                    type="checkbox"
                    id="is_available"
                    name="is_available"
                    checked={formData.is_available}
                    onChange={handleChange}
                    className="h-4 w-4 rounded border-slate-200 text-indigo-600 focus:ring-indigo-500 accent-indigo-600 cursor-pointer"
                  />
                </div>
              </div>

              {/* IMAGE MANAGER (Edit Mode Only) */}
              {isEdit && product && (
                <div className="grid grid-cols-3 gap-4 items-start py-4 border-b border-slate-100">
                  <div className="col-span-1">
                    <Label className="text-sm font-bold text-slate-800">Banner image</Label>
                  </div>
                  <div className="col-span-2 space-y-4">
                    {/* Image display */}
                    <div className="relative aspect-video w-full rounded-lg border bg-slate-50 overflow-hidden flex items-center justify-center">
                      {product.banner?.thumbnail ? (
                        <img
                          src={product.banner.thumbnail}
                          alt={product.name}
                          className="h-full w-full object-cover"
                        />
                      ) : (
                        <span className="text-xs text-slate-400">No image uploaded</span>
                      )}
                    </div>

                    {/* Gallery strip */}
                    {product.gallery && product.gallery.length > 0 && (
                      <div className="flex gap-2 overflow-x-auto pb-1">
                        {product.gallery.map((img, idx) => (
                          <div key={idx} className="h-12 w-12 rounded border overflow-hidden shrink-0">
                            <img
                              src={img.thumbnail || img.preview || img.detail || ''}
                              alt="gallery"
                              className="h-full w-full object-cover"
                            />
                          </div>
                        ))}
                      </div>
                    )}

                    {/* Upload button */}
                    <div>
                      <input
                        type="file"
                        id="sheet-image-upload"
                        accept="image/*"
                        onChange={handleFileChange}
                        className="hidden"
                        disabled={uploading}
                      />
                      <Label
                        htmlFor="sheet-image-upload"
                        className="flex flex-col items-center justify-center w-full h-20 border border-dashed border-slate-200 rounded-lg cursor-pointer hover:bg-slate-50 hover:border-indigo-400 transition-colors"
                      >
                        {uploading ? (
                          <Loader2 className="h-5 w-5 animate-spin text-indigo-600" />
                        ) : (
                          <div className="flex flex-col items-center justify-center text-slate-400">
                            <UploadCloud className="h-5 w-5 mb-0.5 text-slate-300" />
                            <span className="text-[10px] font-semibold">Upload product image</span>
                          </div>
                        )}
                      </Label>
                    </div>
                  </div>
                </div>
              )}

              {/* INVENTORY SECTION */}
              <div className="grid grid-cols-3 gap-4 items-start py-4 border-b border-slate-100">
                <div className="col-span-1">
                  <Label className="text-sm font-bold text-slate-800">Inventory Stock</Label>
                </div>
                <div className="col-span-2 space-y-4">
                  {/* Local Inventory List */}
                  <div className="space-y-3">
                    {localInventories.filter(item => !item.isDeleted).length > 0 ? (
                      <div className="space-y-2">
                        {localInventories
                          .filter(item => !item.isDeleted)
                          .map((item) => (
                            <div
                              key={item.shopId}
                              className="flex items-center justify-between p-2.5 bg-slate-50 border rounded-lg hover:border-indigo-200 transition-colors gap-2"
                            >
                              <div className="flex-1 min-w-0">
                                <span className="text-xs font-semibold text-slate-700 block truncate" title={item.shopName}>
                                  {item.shopName}
                                </span>
                                {item.isNew && (
                                  <span className="text-[9px] text-emerald-600 font-bold uppercase tracking-wider">New</span>
                                )}
                                {!item.isNew && item.isModified && (
                                  <span className="text-[9px] text-amber-600 font-bold uppercase tracking-wider">Modified</span>
                                )}
                              </div>
                              <div className="flex items-center gap-1.5 shrink-0">
                                <span className="text-[10px] text-slate-400 select-none">Qty:</span>
                                <Input
                                  type="number"
                                  min="0"
                                  value={item.stock}
                                  onChange={(e) => handleStockChange(item.shopId, e.target.value)}
                                  className="h-8 w-20 text-xs px-2 text-right border-slate-200 focus-visible:ring-indigo-500 font-semibold"
                                />
                                <button
                                  type="button"
                                  onClick={() => handleRemoveShop(item.shopId)}
                                  className="p-1 rounded text-slate-400 hover:text-red-500 hover:bg-red-50 transition-colors"
                                >
                                  <X className="h-4 w-4" />
                                </button>
                              </div>
                            </div>
                          ))}
                      </div>
                    ) : (
                      <div className="text-xs text-slate-400 italic py-2">
                        No inventories assigned to shops yet.
                      </div>
                    )}

                    {/* Add shop selector */}
                    {availableShops.length > 0 && (
                      <div className="flex items-center gap-2 pt-1 border-t border-dashed border-slate-100">
                        <select
                          value={selectedShopId}
                          onChange={(e) => {
                            setSelectedShopId(e.target.value);
                            handleAddShopToInventory(e.target.value);
                          }}
                          className="flex h-8 flex-1 rounded border border-slate-200 bg-white px-2.5 text-xs focus:outline-none focus:ring-1 focus:ring-indigo-500 text-slate-700"
                        >
                          <option value="">+ Assign shop location...</option>
                          {availableShops.map((s) => (
                            <option key={s.id} value={s.id}>
                              {s.name}
                            </option>
                          ))}
                        </select>
                      </div>
                    )}
                  </div>
                </div>
              </div>
            </form>
          )}
        </div>

        {/* Sticky Footer */}
        <div className="absolute bottom-0 left-0 right-0 h-16 border-t bg-white px-6 flex items-center justify-between shrink-0 z-10">
          <div className="flex items-center gap-2">
            {!isEdit && (
              <label className="flex items-center gap-2 cursor-pointer text-xs text-slate-500 hover:text-slate-600 select-none">
                <input
                  type="checkbox"
                  checked={createMore}
                  onChange={(e) => setCreateMore(e.target.checked)}
                  className="h-3.5 w-3.5 rounded border-slate-300 text-indigo-600 focus:ring-indigo-500 accent-indigo-600"
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
              className="text-xs h-8 px-3 border-slate-200 hover:bg-slate-50 text-slate-600"
              onClick={() => onOpenChange(false)}
              disabled={saving}
            >
              Cancel
            </Button>
            <Button
              type="submit"
              size="sm"
              className="text-xs h-8 px-3 bg-indigo-600 hover:bg-indigo-700 text-white font-medium shadow-sm flex items-center gap-1.5"
              onClick={handleSave}
              disabled={saving || loading}
            >
              {saving ? (
                <Loader2 className="h-3 w-3 animate-spin" />
              ) : (
                <Save className="h-3.5 w-3.5" />
              )}
              {isEdit ? 'Save Product' : 'Create Product'}
            </Button>
          </div>
        </div>
      </SheetContent>
    </Sheet>
  );
}
