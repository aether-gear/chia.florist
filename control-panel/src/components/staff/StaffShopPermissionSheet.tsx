import React, { useState, useEffect } from 'react';
import {
  Store,
  Shield,
  Loader2,
  AlertCircle,
  Package,
  Boxes,
  Truck,
  MapPin,
  FileText,
  Trash2,
  Edit3,
  Save,
  ChevronDown,
  ChevronRight,
} from 'lucide-react';
import {
  Sheet,
  SheetContent,
} from '../ui/sheet';
import { Button } from '../ui/button';
import { Label } from '../ui/label';
import { Switch } from '../ui/switch';

import { fetchApi } from '@/lib/api';
import type { StaffShopPermission, SaveStaffShopPermissionPayload } from '@/models/Staff';

interface StaffShopPermissionSheetProps {
  isOpen: boolean;
  onClose: () => void;
  staffName: string;
  existingPermission?: StaffShopPermission | null;
  onSave: (payload: SaveStaffShopPermissionPayload) => Promise<void>;
}

interface PermissionItem {
  key: string;
  label: string;
  description: string;
  icon: React.ComponentType<{ className?: string }>;
}

interface PermissionCategory {
  id: string;
  title: string;
  description: string;
  icon: React.ComponentType<{ className?: string }>;
  items: PermissionItem[];
}

const CATEGORIZED_PERMISSIONS: PermissionCategory[] = [
  {
    id: 'shop',
    title: 'Shop Settings',
    description: 'Operating parameters and shop profile details',
    icon: Store,
    items: [
      {
        key: 'shop:update',
        label: 'Edit Shop Settings',
        description: 'Update shop name, description, active status, and preferences.',
        icon: Edit3,
      },
    ],
  },
  {
    id: 'catalog',
    title: 'Product Catalog',
    description: 'Adding, editing, and deleting items in the shop catalog',
    icon: Package,
    items: [
      {
        key: 'product:create',
        label: 'Create Products',
        description: 'Add new items and listings to this shop.',
        icon: Package,
      },
      {
        key: 'product:update',
        label: 'Edit Product Details & Prices',
        description: 'Modify product prices, descriptions, and SKU information.',
        icon: Edit3,
      },
      {
        key: 'product:delete',
        label: 'Delete Products',
        description: 'Remove products from this shop catalog.',
        icon: Trash2,
      },
    ],
  },
  {
    id: 'inventory',
    title: 'Stock & Inventory',
    description: 'Warehouse quantities and stock adjustments',
    icon: Boxes,
    items: [
      {
        key: 'inventory:manage',
        label: 'Manage Stock & Inventory',
        description: 'Update available stock levels and warehouse inventory.',
        icon: Boxes,
      },
    ],
  },
  {
    id: 'orders',
    title: 'Orders & Sales',
    description: 'Order fulfillment and transaction status updates',
    icon: FileText,
    items: [
      {
        key: 'order:read',
        label: 'View Shop Orders',
        description: 'Read customer orders and purchase details.',
        icon: FileText,
      },
      {
        key: 'order:update_status',
        label: 'Process & Update Order Status',
        description: 'Change order status (e.g. processing, shipping, fulfilled).',
        icon: Truck,
      },
    ],
  },
  {
    id: 'shipping',
    title: 'Couriers & Delivery',
    description: 'Shipping methods and courier configuration',
    icon: Truck,
    items: [
      {
        key: 'courier:manage',
        label: 'Configure Couriers',
        description: 'Enable or adjust shipping couriers for order delivery.',
        icon: Truck,
      },
    ],
  },
  {
    id: 'locations',
    title: 'Shop Addresses',
    description: 'Physical address and location management',
    icon: MapPin,
    items: [
      {
        key: 'address:manage',
        label: 'Manage Shop Addresses',
        description: 'Create, update, or remove physical shop addresses.',
        icon: MapPin,
      },
    ],
  },
];

const ALL_PERMISSION_KEYS = CATEGORIZED_PERMISSIONS.flatMap((cat) => cat.items.map((i) => i.key));

export default function StaffShopPermissionSheet({
  isOpen,
  onClose,
  staffName,
  existingPermission,
  onSave,
}: StaffShopPermissionSheetProps) {
  const [shops, setShops] = useState<Array<{ id: string; name: string; slug: string }>>([]);
  const [selectedShopId, setSelectedShopId] = useState<string>('');
  const [selectedPermissions, setSelectedPermissions] = useState<string[]>([]);
  const [loadingShops, setLoadingShops] = useState(false);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Expandable Category Sections (all open by default)
  const [openCategories, setOpenCategories] = useState<Record<string, boolean>>(() => {
    const initial: Record<string, boolean> = {};
    CATEGORIZED_PERMISSIONS.forEach((cat) => {
      initial[cat.id] = true;
    });
    return initial;
  });

  useEffect(() => {
    if (isOpen) {
      setError(null);

      // Ensure all categories start open when sheet opens
      const initialOpen: Record<string, boolean> = {};
      CATEGORIZED_PERMISSIONS.forEach((cat) => {
        initialOpen[cat.id] = true;
      });
      setOpenCategories(initialOpen);

      if (existingPermission) {
        setSelectedShopId(existingPermission.shop_id);
        setSelectedPermissions(existingPermission.permissions || []);
      } else {
        setSelectedShopId('');
        setSelectedPermissions(ALL_PERMISSION_KEYS); // default select all for new shop assignment
      }

      // Fetch available shops
      const loadShops = async () => {
        setLoadingShops(true);
        try {
          const res = await fetchApi('/shops');
          setShops(res.shops || []);
          if (!existingPermission && res.shops && res.shops.length > 0) {
            setSelectedShopId(res.shops[0].id);
          }
        } catch (err) {
          console.error('Failed to fetch shops for assignment', err);
        } finally {
          setLoadingShops(false);
        }
      };
      loadShops();
    }
  }, [isOpen, existingPermission]);

  const toggleCategoryOpen = (catId: string) => {
    setOpenCategories((prev) => ({
      ...prev,
      [catId]: !prev[catId],
    }));
  };

  const handleExpandAllCategories = () => {
    const next: Record<string, boolean> = {};
    CATEGORIZED_PERMISSIONS.forEach((cat) => {
      next[cat.id] = true;
    });
    setOpenCategories(next);
  };

  const handleCollapseAllCategories = () => {
    const next: Record<string, boolean> = {};
    CATEGORIZED_PERMISSIONS.forEach((cat) => {
      next[cat.id] = false;
    });
    setOpenCategories(next);
  };

  const togglePermission = (key: string) => {
    setSelectedPermissions((prev) =>
      prev.includes(key) ? prev.filter((p) => p !== key) : [...prev, key]
    );
  };

  const isCategoryFullySelected = (cat: PermissionCategory) => {
    return cat.items.every((item) => selectedPermissions.includes(item.key));
  };

  const getCategoryEnabledCount = (cat: PermissionCategory) => {
    return cat.items.filter((item) => selectedPermissions.includes(item.key)).length;
  };

  const toggleCategory = (cat: PermissionCategory) => {
    const catKeys = cat.items.map((i) => i.key);
    if (isCategoryFullySelected(cat)) {
      setSelectedPermissions((prev) => prev.filter((k) => !catKeys.includes(k)));
    } else {
      setSelectedPermissions((prev) => Array.from(new Set([...prev, ...catKeys])));
    }
  };

  const handleSelectAllGlobal = () => {
    setSelectedPermissions(ALL_PERMISSION_KEYS);
  };

  const handleClearAllGlobal = () => {
    setSelectedPermissions([]);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedShopId) {
      setError('Please select a shop to assign.');
      return;
    }

    setSaving(true);
    setError(null);
    try {
      await onSave({
        shop_id: selectedShopId,
        permissions: selectedPermissions,
      });
      onClose();
    } catch (err: any) {
      console.error('Failed to save staff shop permission', err);
      setError(err.message || 'Failed to save shop permissions');
    } finally {
      setSaving(false);
    }
  };

  return (
    <Sheet open={isOpen} onOpenChange={(open) => !open && onClose()}>
      <SheetContent className="w-full sm:max-w-none md:w-[48vw] md:min-w-[480px] p-0 flex flex-col h-full border-l border-border/60 bg-background shadow-2xl">
        {/* Header - Matches StaffFormSheet */}
        <div className="flex items-center justify-between px-6 py-5 border-b border-border/60 pr-12">
          <div>
            <h2 className="text-xl font-bold font-display tracking-tight text-foreground flex items-center gap-2">
              <Shield className="h-5 w-5 text-primary" />
              {existingPermission ? 'Edit Shop Access Rights' : 'Assign Shop & Permissions'}
            </h2>
            <p className="text-xs text-muted-foreground mt-0.5">
              Configure shop management capabilities for <span className="font-semibold text-foreground">{staffName}</span>
            </p>
          </div>
        </div>

        {/* Scrollable Content Body - Flat Minimalist UI */}
        <form id="permission-sheet-form" onSubmit={handleSubmit} className="flex-1 overflow-y-auto px-6 py-6 space-y-6">
          {error && (
            <div className="flex items-center gap-2.5 p-3.5 text-xs text-rose-600 bg-rose-500/10 rounded-xl border border-rose-500/20 animate-in fade-in">
              <AlertCircle className="h-4 w-4 shrink-0" />
              <span>{error}</span>
            </div>
          )}

          {/* Section 1: Target Shop */}
          <div className="space-y-3">
            <div className="pb-2 border-b border-border/40">
              <h3 className="text-xs font-bold font-display uppercase tracking-wider text-muted-foreground">
                Target Shop
              </h3>
            </div>
            <div className="space-y-2">
              <Label htmlFor="target-shop-select" className="text-xs font-semibold text-foreground">
                Select Shop <span className="text-rose-500">*</span>
              </Label>
              {loadingShops ? (
                <div className="flex items-center gap-2 text-xs text-muted-foreground p-3 bg-muted/30 rounded-xl animate-pulse">
                  <Loader2 className="h-4 w-4 animate-spin text-primary" /> Loading shops list...
                </div>
              ) : (
                <select
                  id="target-shop-select"
                  value={selectedShopId}
                  onChange={(e) => setSelectedShopId(e.target.value)}
                  disabled={Boolean(existingPermission) || shops.length === 0}
                  className="w-full h-11 px-3.5 rounded-xl border border-border bg-background text-sm font-medium text-foreground focus:ring-2 focus:ring-primary/20 focus:border-primary disabled:opacity-60 transition"
                >
                  {shops.length === 0 ? (
                    <option value="">No active shops available</option>
                  ) : (
                    shops.map((shop) => (
                      <option key={shop.id} value={shop.id}>
                        {shop.name} ({shop.slug})
                      </option>
                    ))
                  )}
                </select>
              )}
              <p className="text-[11px] text-muted-foreground">
                Select which shop this staff member is authorized to manage.
              </p>
            </div>
          </div>

          {/* Section 2: Permissions List with Expandable Category Sections */}
          <div className="space-y-5 pt-2">
            <div className="flex flex-col sm:flex-row sm:items-center justify-between pb-2 border-b border-border/40 gap-2">
              <h3 className="text-xs font-bold font-display uppercase tracking-wider text-muted-foreground">
                Shop Capabilities ({selectedPermissions.length}/{ALL_PERMISSION_KEYS.length})
              </h3>

              {/* Global Quick Action Controls */}
              <div className="flex items-center gap-3 text-xs flex-wrap">
                <button
                  type="button"
                  onClick={handleExpandAllCategories}
                  className="text-[11px] font-medium text-muted-foreground hover:text-foreground"
                >
                  Expand All
                </button>
                <span className="text-muted-foreground/30">•</span>
                <button
                  type="button"
                  onClick={handleCollapseAllCategories}
                  className="text-[11px] font-medium text-muted-foreground hover:text-foreground"
                >
                  Collapse All
                </button>
                <span className="text-muted-foreground/30">•</span>
                <button
                  type="button"
                  onClick={handleSelectAllGlobal}
                  className="text-xs font-semibold text-primary hover:underline"
                >
                  Enable All
                </button>
                <span className="text-muted-foreground/30">•</span>
                <button
                  type="button"
                  onClick={handleClearAllGlobal}
                  className="text-xs font-medium text-muted-foreground hover:text-foreground"
                >
                  Disable All
                </button>
              </div>
            </div>

            {/* Expandable Category Sections */}
            <div className="space-y-8 py-2">
              {CATEGORIZED_PERMISSIONS.map((cat) => {
                const fullySelected = isCategoryFullySelected(cat);
                const enabledCount = getCategoryEnabledCount(cat);
                const isCategoryOpen = Boolean(openCategories[cat.id]);

                return (
                  <div key={cat.id} className="space-y-2">
                    {/* Expandable Section Header */}
                    <div className="flex items-center justify-between">
                      <button
                        type="button"
                        onClick={() => toggleCategoryOpen(cat.id)}
                        className="flex items-center gap-2 text-left group cursor-pointer"
                        title={isCategoryOpen ? 'Collapse section' : 'Expand section'}
                      >
                        <div className="p-0.5 rounded text-muted-foreground group-hover:text-foreground transition">
                          {isCategoryOpen ? (
                            <ChevronDown className="h-4 w-4 text-primary" />
                          ) : (
                            <ChevronRight className="h-4 w-4" />
                          )}
                        </div>
                        <h4 className="text-xs font-bold text-foreground font-display uppercase tracking-wider group-hover:text-primary transition">
                          {cat.title}
                        </h4>
                        <span className="text-[10px] font-mono font-medium px-2 py-0.5 rounded-full bg-muted text-muted-foreground">
                          {enabledCount}/{cat.items.length}
                        </span>
                      </button>

                      <button
                        type="button"
                        onClick={(e) => {
                          e.stopPropagation();
                          toggleCategory(cat);
                        }}
                        className="text-[11px] font-medium text-primary hover:underline"
                      >
                        {fullySelected ? 'Disable Category' : 'Enable Category'}
                      </button>
                    </div>

                    {/* Expandable Capability Rows */}
                    {isCategoryOpen && (
                      <div className="divide-y divide-border/20 animate-in fade-in duration-150">
                        {cat.items.map((item) => {
                          const isChecked = selectedPermissions.includes(item.key);
                          const ItemIcon = item.icon;

                          return (
                            <div
                              key={item.key}
                              onClick={() => togglePermission(item.key)}
                              className="py-3 flex items-center justify-between gap-4 cursor-pointer hover:bg-muted/10 px-1.5 rounded-lg transition"
                            >
                              <div className="flex items-center gap-2.5 min-w-0 pr-2">
                                <div className={`p-1.5 rounded-lg shrink-0 transition-colors ${isChecked ? 'bg-primary/10 text-primary' : 'bg-muted text-muted-foreground'}`}>
                                  <ItemIcon className="h-3.5 w-3.5" />
                                </div>
                                <div className="space-y-0.5 min-w-0">
                                  <Label className="text-xs font-semibold text-foreground cursor-pointer">
                                    {item.label}
                                  </Label>
                                  <p className="text-[11px] text-muted-foreground leading-normal">
                                    {item.description}
                                  </p>
                                </div>
                              </div>
                              <div onClick={(e) => e.stopPropagation()}>
                                <Switch
                                  checked={isChecked}
                                  onCheckedChange={() => togglePermission(item.key)}
                                  className="shrink-0"
                                />
                              </div>
                            </div>
                          );
                        })}
                      </div>
                    )}
                  </div>
                );
              })}
            </div>
          </div>

        </form>

        {/* Footer - Matches StaffFormSheet */}
        <div className="px-6 py-4 border-t border-border/60 bg-background shrink-0 flex items-center justify-between">
          <div className="text-xs text-muted-foreground">
            <span className="font-bold text-foreground">{selectedPermissions.length}</span> capabilities enabled
          </div>
          <div className="flex items-center gap-2">
            <Button type="button" variant="outline" onClick={onClose} disabled={saving} className="rounded-xl text-xs h-9 px-4">
              Cancel
            </Button>
            <Button
              type="submit"
              form="permission-sheet-form"
              disabled={saving || !selectedShopId}
              className="rounded-xl text-xs font-semibold h-9 px-4 gap-2 bg-primary text-primary-foreground"
            >
              {saving ? <Loader2 className="h-3.5 w-3.5 animate-spin" /> : <Save className="h-3.5 w-3.5" />}
              {existingPermission ? 'Save Changes' : 'Assign Access'}
            </Button>
          </div>
        </div>
      </SheetContent>
    </Sheet>
  );
}
