import React, { useState, useEffect } from 'react';
import { Truck, Plus, Trash2, Boxes, Store } from 'lucide-react';
import { Button } from '../ui/button';
import { Input } from '../ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '../ui/select';
import type { Order } from '../../models/Order';
import type { ShipmentDispatchPayload } from '../../viewmodels/useOrderActionsViewModel';

export interface ShipmentGroupForm {
  id: string;
  shop_id?: string;
  shop_name?: string;
  fulfillment_method: string;
  courier: string;
  service: string;
  tracking_number: string;
  item_ids: string[];
}

export interface ShipmentDispatchSectionProps {
  order: Order;
  submitting: boolean;
  onDispatchOrder: (orderId: string, shipments: ShipmentDispatchPayload[]) => Promise<void>;
}

export const ShipmentDispatchSection: React.FC<ShipmentDispatchSectionProps> = ({
  order,
  submitting,
  onDispatchOrder,
}) => {
  const [shipmentGroups, setShipmentGroups] = useState<ShipmentGroupForm[]>([]);

  // Initialize shipment groups from order items
  useEffect(() => {
    if (order && order.status === 'processing') {
      const shopMap = new Map<string, { shop_name: string; items: typeof order.items }>();
      order.items.forEach((item) => {
        const shopId = item.shop_id || 'default-shop';
        const existing = shopMap.get(shopId) || { shop_name: item.shop_name || 'Shop', items: [] };
        existing.items.push(item);
        shopMap.set(shopId, existing);
      });

      const initialGroups: ShipmentGroupForm[] = [];
      let groupIndex = 1;
      shopMap.forEach((shopData, shopId) => {
        const firstItem = shopData.items[0];
        initialGroups.push({
          id: `shipment-shop-${shopId}-${Date.now()}-${groupIndex}`,
          shop_id: shopId,
          shop_name: shopData.shop_name,
          fulfillment_method: 'courier',
          courier: firstItem?.courier_code || 'JNE',
          service: firstItem?.courier_service || 'REG',
          tracking_number: '',
          item_ids: shopData.items.map((i) => i.id),
        });
        groupIndex++;
      });

      setShipmentGroups(initialGroups);
    } else {
      setShipmentGroups([]);
    }
  }, [order.id, order.status]);

  const handleAddShipmentGroup = (targetShopId?: string) => {
    const shopItems = targetShopId ? order.items.filter((i) => i.shop_id === targetShopId) : order.items;
    const firstItem = shopItems[0] || order.items[0];

    setShipmentGroups((prev) => [
      ...prev,
      {
        id: `shipment-${Date.now()}-${prev.length + 1}`,
        shop_id: targetShopId || firstItem?.shop_id,
        shop_name: firstItem?.shop_name || 'Shop',
        fulfillment_method: 'courier',
        courier: firstItem?.courier_code || 'JNE',
        service: firstItem?.courier_service || 'REG',
        tracking_number: '',
        item_ids: [],
      },
    ]);
  };

  const handleRemoveShipmentGroup = (groupId: string) => {
    if (shipmentGroups.length <= 1) return;
    const removedGroup = shipmentGroups.find((g) => g.id === groupId);
    const remainingGroups = shipmentGroups.filter((g) => g.id !== groupId);

    // Reassign unallocated items from removed group back to matching shop shipment
    if (removedGroup && removedGroup.item_ids.length > 0 && remainingGroups.length > 0) {
      const matchShopGroup = remainingGroups.find((g) => g.shop_id === removedGroup.shop_id) || remainingGroups[0];
      matchShopGroup.item_ids = [...matchShopGroup.item_ids, ...removedGroup.item_ids];
    }
    setShipmentGroups([...remainingGroups]);
  };

  const handleToggleItemInGroup = (groupId: string, itemId: string) => {
    setShipmentGroups((prev) =>
      prev.map((g) => {
        if (g.id === groupId) {
          const exists = g.item_ids.includes(itemId);
          return {
            ...g,
            item_ids: exists ? g.item_ids.filter((id) => id !== itemId) : [...g.item_ids, itemId],
          };
        } else {
          return {
            ...g,
            item_ids: g.item_ids.filter((id) => id !== itemId),
          };
        }
      })
    );
  };

  const handleUpdateGroupField = (groupId: string, field: keyof ShipmentGroupForm, value: any) => {
    setShipmentGroups((prev) => prev.map((g) => (g.id === groupId ? { ...g, [field]: value } : g)));
  };

  const handleDispatch = async () => {
    if (shipmentGroups.length === 0) {
      alert('At least one shipment is required to dispatch the order.');
      return;
    }

    const allAssignedItemIds = shipmentGroups.flatMap((g) => g.item_ids);
    const unassignedItems = order.items.filter((item) => !allAssignedItemIds.includes(item.id));
    if (unassignedItems.length > 0) {
      alert(`All order items must be assigned to a shipment. ${unassignedItems.length} item(s) unassigned.`);
      return;
    }

    for (let i = 0; i < shipmentGroups.length; i++) {
      const group = shipmentGroups[i];
      if (group.item_ids.length === 0) {
        alert(`Shipment #${i + 1} (${group.shop_name}) has no items assigned. Please assign items or remove it.`);
        return;
      }
      if (group.fulfillment_method === 'courier' && !group.tracking_number.trim()) {
        alert(`Shipment #${i + 1} (${group.shop_name}) is missing a waybill tracking number.`);
        return;
      }
    }

    const payload: ShipmentDispatchPayload[] = shipmentGroups.map((g) => ({
      fulfillment_method: g.fulfillment_method,
      courier: g.courier,
      service: g.service,
      tracking_number: g.tracking_number.trim() || undefined,
      item_ids: g.item_ids,
    }));

    await onDispatchOrder(order.id, payload);
  };

  return (
    <div className="border border-border/80 bg-background p-5 rounded-xl space-y-5">
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2">
        <div className="space-y-0.5">
          <h5 className="text-sm font-semibold text-foreground flex items-center gap-1.5">
            <Truck className="h-4 w-4 text-primary" /> Dispatch Shipping & Assign Courier
          </h5>
          <p className="text-xs text-muted-foreground">
            Shipments are grouped by shop warehouse by default. You can split products into additional shipments if needed.
          </p>
        </div>
        <Button
          type="button"
          variant="outline"
          size="sm"
          onClick={() => handleAddShipmentGroup()}
          className="rounded-xl text-xs h-8 gap-1 self-start sm:self-auto border-primary/40 text-primary hover:bg-primary/5 cursor-pointer"
        >
          <Plus className="h-3.5 w-3.5" /> Add Split Shipment
        </Button>
      </div>

      {/* Shipment Groups List */}
      <div className="space-y-4">
        {shipmentGroups.map((group, groupIdx) => {
          const shopItems = group.shop_id ? order.items.filter((i) => i.shop_id === group.shop_id) : order.items;

          return (
            <div key={group.id} className="border border-border/80 rounded-xl p-4 bg-muted/10 space-y-4 relative">
              <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2 border-b border-border/40 pb-2.5">
                <div className="flex items-center gap-2 flex-wrap">
                  <span className="h-5 w-5 rounded-full bg-primary/10 text-primary text-[11px] font-bold flex items-center justify-center">
                    {groupIdx + 1}
                  </span>
                  <span className="text-xs font-bold text-foreground">Shipment #{groupIdx + 1}</span>
                  {group.shop_name && (
                    <span className="inline-flex items-center gap-1 text-[11px] bg-secondary text-secondary-foreground px-2 py-0.5 rounded-md font-medium">
                      <Store className="h-3 w-3" /> {group.shop_name}
                    </span>
                  )}
                  <span className="text-[10px] text-muted-foreground font-mono">
                    ({group.item_ids.length} item{group.item_ids.length !== 1 ? 's' : ''})
                  </span>
                </div>

                <div className="flex items-center gap-2">
                  {shipmentGroups.length > 1 && (
                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      onClick={() => handleRemoveShipmentGroup(group.id)}
                      className="h-6 px-2 text-destructive hover:bg-destructive/10 text-xs rounded-lg cursor-pointer"
                    >
                      <Trash2 className="h-3.5 w-3.5 mr-1" /> Remove
                    </Button>
                  )}
                </div>
              </div>

              {/* Items Selector in this group */}
              <div className="space-y-2">
                <div className="flex justify-between items-center">
                  <label className="text-[11px] font-semibold text-muted-foreground uppercase flex items-center gap-1">
                    <Boxes className="h-3.5 w-3.5" /> Assigned Products ({group.shop_name || 'All Shops'})
                  </label>
                  <span className="text-[10px] text-muted-foreground">
                    {group.item_ids.length} of {shopItems.length} selected
                  </span>
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
                  {shopItems.map((item) => {
                    const isChecked = group.item_ids.includes(item.id);
                    const otherGroupIdx = shipmentGroups.findIndex(
                      (g) => g.id !== group.id && g.item_ids.includes(item.id)
                    );

                    return (
                      <label
                        key={item.id}
                        className={`flex items-start gap-2.5 p-2.5 rounded-lg border text-xs cursor-pointer transition-all ${
                          isChecked
                            ? 'border-primary/60 bg-primary/5 text-foreground font-semibold ring-1 ring-primary/20'
                            : otherGroupIdx >= 0
                            ? 'border-border/40 opacity-55 hover:opacity-80 bg-muted/20 text-muted-foreground'
                            : 'border-border/60 hover:bg-muted/30 text-muted-foreground'
                        }`}
                      >
                        <input
                          type="checkbox"
                          checked={isChecked}
                          onChange={() => handleToggleItemInGroup(group.id, item.id)}
                          className="mt-0.5 rounded text-primary focus:ring-primary h-3.5 w-3.5 cursor-pointer"
                        />
                        <div className="flex-1 min-w-0">
                          <div className="truncate text-xs">{item.product_name}</div>
                          <div className="text-[10px] text-muted-foreground font-normal flex items-center justify-between">
                            <span>Qty: {item.quantity}</span>
                            {otherGroupIdx >= 0 && (
                              <span className="text-amber-600 dark:text-amber-400 font-medium">
                                (In Pkg #{otherGroupIdx + 1})
                              </span>
                            )}
                          </div>
                        </div>
                      </label>
                    );
                  })}
                </div>
              </div>

              {/* Logistics Fields for this shipment */}
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 pt-2">
                <div className="space-y-1.5">
                  <label className="text-[11px] font-semibold text-muted-foreground">Fulfillment Method</label>
                  <Select
                    value={group.fulfillment_method}
                    onValueChange={(val) => handleUpdateGroupField(group.id, 'fulfillment_method', val)}
                  >
                    <SelectTrigger className="rounded-xl h-8 text-xs bg-background">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="courier">Courier Logistics</SelectItem>
                      <SelectItem value="self_delivery">Self Delivery (Shop Rider)</SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                {group.fulfillment_method === 'courier' && (
                  <>
                    <div className="space-y-1.5">
                      <label className="text-[11px] font-semibold text-muted-foreground">Courier Brand</label>
                      <Input
                        value={group.courier}
                        onChange={(e) => handleUpdateGroupField(group.id, 'courier', e.target.value)}
                        placeholder="JNE, POS, SiCepat..."
                        className="rounded-xl h-8 text-xs bg-background"
                      />
                    </div>
                    <div className="space-y-1.5">
                      <label className="text-[11px] font-semibold text-muted-foreground">Service Level</label>
                      <Input
                        value={group.service}
                        onChange={(e) => handleUpdateGroupField(group.id, 'service', e.target.value)}
                        placeholder="REG, YES, etc."
                        className="rounded-xl h-8 text-xs bg-background"
                      />
                    </div>
                    <div className="space-y-1.5">
                      <label className="text-[11px] font-semibold text-muted-foreground">Waybill Tracking Number</label>
                      <Input
                        value={group.tracking_number}
                        onChange={(e) => handleUpdateGroupField(group.id, 'tracking_number', e.target.value)}
                        placeholder="Insert shipping receipt ID"
                        className="rounded-xl h-8 text-xs bg-background"
                      />
                    </div>
                  </>
                )}
              </div>
            </div>
          );
        })}
      </div>

      <Button
        disabled={submitting}
        onClick={handleDispatch}
        className="rounded-xl bg-primary text-primary-foreground hover:bg-primary/90 text-xs font-semibold h-9 w-full shadow-sm cursor-pointer"
      >
        <Truck className="h-4 w-4 mr-1.5" /> Dispatch {shipmentGroups.length > 1 ? `${shipmentGroups.length} Shipments` : 'Shipment'} & Mark Shipped
      </Button>
    </div>
  );
};

export default ShipmentDispatchSection;
