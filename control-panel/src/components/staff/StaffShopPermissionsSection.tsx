import { useState, useEffect, useCallback } from 'react';

import {
  Store,
  Plus,
  Shield,
  Trash2,
  Edit2,
  Loader2,
  AlertCircle,
  Package,
  Boxes,
  Truck,
  MapPin,
  FileText,
  Building,
} from 'lucide-react';
import { Button } from '../ui/button';
import type { StaffShopPermission, SaveStaffShopPermissionPayload } from '@/models/Staff';
import StaffShopPermissionSheet from './StaffShopPermissionSheet';


interface StaffShopPermissionsSectionProps {
  staffId: string;
  staffName: string;
  fetchPermissions: (staffId: string) => Promise<StaffShopPermission[]>;
  savePermission: (staffId: string, payload: SaveStaffShopPermissionPayload) => Promise<void>;
  deletePermission: (staffId: string, shopId: string) => Promise<void>;
}

export default function StaffShopPermissionsSection({
  staffId,
  staffName,
  fetchPermissions,
  savePermission,
  deletePermission,
}: StaffShopPermissionsSectionProps) {
  const [permissions, setPermissions] = useState<StaffShopPermission[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  // Modal State
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingPermission, setEditingPermission] = useState<StaffShopPermission | null>(null);
  const [deletingShopId, setDeletingShopId] = useState<string | null>(null);

  const loadPermissions = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await fetchPermissions(staffId);
      setPermissions(data);
    } catch (err: any) {
      console.error('Failed to load staff shop permissions', err);
      setError(err.message || 'Failed to load shop access permissions');
    } finally {
      setLoading(false);
    }
  }, [staffId, fetchPermissions]);

  useEffect(() => {
    loadPermissions();
  }, [loadPermissions]);

  const handleOpenAssignModal = () => {
    setEditingPermission(null);
    setIsModalOpen(true);
  };

  const handleOpenEditModal = (perm: StaffShopPermission) => {
    setEditingPermission(perm);
    setIsModalOpen(true);
  };

  const handleSavePermission = async (payload: SaveStaffShopPermissionPayload) => {
    await savePermission(staffId, payload);
    await loadPermissions();
  };

  const handleDeletePermission = async (shopId: string) => {
    setDeletingShopId(shopId);
    try {
      await deletePermission(staffId, shopId);
      await loadPermissions();
    } catch (err: any) {
      console.error('Failed to delete shop permission', err);
    } finally {
      setDeletingShopId(null);
    }
  };

  const renderPermissionBadge = (permKey: string) => {
    switch (permKey) {
      case 'shop:update':
        return (
          <span key={permKey} className="inline-flex items-center gap-1 text-[11px] font-medium px-2.5 py-0.5 rounded-full bg-blue-500/10 text-blue-400 border border-blue-500/20">
            <Store className="h-3 w-3" /> Shop Settings
          </span>
        );
      case 'product:create':
      case 'product:update':
      case 'product:delete':
        return (
          <span key={permKey} className="inline-flex items-center gap-1 text-[11px] font-medium px-2.5 py-0.5 rounded-full bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
            <Package className="h-3 w-3" /> {permKey.replace('product:', '')} product
          </span>
        );
      case 'inventory:manage':
        return (
          <span key={permKey} className="inline-flex items-center gap-1 text-[11px] font-medium px-2.5 py-0.5 rounded-full bg-amber-500/10 text-amber-400 border border-amber-500/20">
            <Boxes className="h-3 w-3" /> Inventory
          </span>
        );
      case 'order:read':
      case 'order:update_status':
        return (
          <span key={permKey} className="inline-flex items-center gap-1 text-[11px] font-medium px-2.5 py-0.5 rounded-full bg-purple-500/10 text-purple-400 border border-purple-500/20">
            <FileText className="h-3 w-3" /> Orders
          </span>
        );
      case 'courier:manage':
        return (
          <span key={permKey} className="inline-flex items-center gap-1 text-[11px] font-medium px-2.5 py-0.5 rounded-full bg-cyan-500/10 text-cyan-400 border border-cyan-500/20">
            <Truck className="h-3 w-3" /> Couriers
          </span>
        );
      case 'address:manage':
        return (
          <span key={permKey} className="inline-flex items-center gap-1 text-[11px] font-medium px-2.5 py-0.5 rounded-full bg-rose-500/10 text-rose-400 border border-rose-500/20">
            <MapPin className="h-3 w-3" /> Address
          </span>
        );
      default:
        return (
          <span key={permKey} className="inline-flex items-center gap-1 text-[11px] font-medium px-2 py-0.5 rounded-full bg-muted text-muted-foreground">
            {permKey}
          </span>
        );
    }
  };

  return (
    <div className="space-y-4 pt-4 border-t border-border/50">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Shield className="h-4 w-4 text-primary" />
          <h4 className="text-sm font-semibold text-foreground">Assigned Shop Permissions</h4>
          <span className="text-xs text-muted-foreground">({permissions.length})</span>
        </div>
        <Button
          type="button"
          variant="outline"
          size="sm"
          onClick={handleOpenAssignModal}
          className="rounded-2xl h-8 px-3 text-xs gap-1.5 border-primary/30 text-primary hover:bg-primary/10"
        >
          <Plus className="h-3.5 w-3.5" /> Assign Shop Access
        </Button>
      </div>

      {/* Loading state */}
      {loading ? (
        <div className="flex items-center justify-center p-6 text-sm text-muted-foreground gap-2 bg-muted/20 rounded-2xl animate-pulse">
          <Loader2 className="h-4 w-4 animate-spin" /> Loading shop permissions...
        </div>
      ) : error ? (
        <div className="p-4 rounded-2xl bg-destructive/10 text-destructive text-sm flex items-center gap-2">
          <AlertCircle className="h-4 w-4 shrink-0" />
          <span>{error}</span>
        </div>
      ) : permissions.length === 0 ? (
        /* Empty State */
        <div className="p-6 rounded-2xl border border-dashed border-border/80 bg-muted/10 text-center space-y-2">
          <div className="h-10 w-10 rounded-full bg-muted flex items-center justify-center mx-auto text-muted-foreground">
            <Building className="h-5 w-5" />
          </div>
          <p className="text-xs font-medium text-foreground">No Shops Assigned</p>
          <p className="text-xs text-muted-foreground max-w-sm mx-auto">
            This staff member currently has no rights to manage any shop. Click "Assign Shop Access" to grant permissions for a specific shop.
          </p>
        </div>
      ) : (
        /* List of Assigned Shops */
        <div className="space-y-3">
          {permissions.map((perm) => (
            <div
              key={perm.shop_id}
              className="p-4 rounded-2xl border border-border/60 bg-muted/10 hover:border-border transition space-y-3"
            >
              <div className="flex items-center justify-between gap-2">
                <div className="flex items-center gap-2.5 min-w-0">
                  <div className="p-2 rounded-xl bg-primary/10 text-primary shrink-0">
                    <Store className="h-4 w-4" />
                  </div>
                  <div className="min-w-0">
                    <h5 className="text-sm font-semibold text-foreground truncate">
                      {perm.shop_name || 'Shop ID: ' + perm.shop_id.substring(0, 8)}
                    </h5>
                    <p className="text-[11px] text-muted-foreground font-mono truncate">
                      {perm.shop_id}
                    </p>
                  </div>
                </div>

                {/* Actions */}
                <div className="flex items-center gap-1.5 shrink-0">
                  <Button
                    type="button"
                    variant="ghost"
                    size="icon"
                    onClick={() => handleOpenEditModal(perm)}
                    className="h-7 w-7 rounded-xl text-muted-foreground hover:text-foreground"
                    title="Edit permissions"
                  >
                    <Edit2 className="h-3.5 w-3.5" />
                  </Button>
                  <Button
                    type="button"
                    variant="ghost"
                    size="icon"
                    onClick={() => handleDeletePermission(perm.shop_id)}
                    disabled={deletingShopId === perm.shop_id}
                    className="h-7 w-7 rounded-xl text-destructive hover:bg-destructive/10"
                    title="Revoke shop access"
                  >
                    {deletingShopId === perm.shop_id ? (
                      <Loader2 className="h-3.5 w-3.5 animate-spin" />
                    ) : (
                      <Trash2 className="h-3.5 w-3.5" />
                    )}
                  </Button>
                </div>
              </div>

              {/* Permissions Badges */}
              <div className="flex flex-wrap items-center gap-1.5 pt-1">
                {perm.permissions && perm.permissions.length > 0 ? (
                  perm.permissions.map((p) => renderPermissionBadge(p))
                ) : (
                  <span className="text-xs text-muted-foreground italic">No specific permissions granted</span>
                )}
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Permission Overlay Sheet */}
      <StaffShopPermissionSheet
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        staffName={staffName}
        existingPermission={editingPermission}
        onSave={handleSavePermission}
      />
    </div>
  );
}

