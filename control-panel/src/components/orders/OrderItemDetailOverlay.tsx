import React from 'react';
import {
  Warehouse,
  Truck,
  ClipboardPenLine,
  Type,
  DollarSign,
  Info,
} from 'lucide-react';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '../ui/dialog';
import { Button } from '../ui/button';
import { Badge } from '../ui/badge';
import type { OrderItem, OrderShipment } from '../../models/Order';

export interface OrderItemDetailOverlayProps {
  item: OrderItem | null;
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
  shipment?: OrderShipment | null;
  packageIndex?: number;
}

const formatCurrency = (amount: number) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
  }).format(amount);
};

export const OrderItemDetailOverlay: React.FC<OrderItemDetailOverlayProps> = ({
  item,
  isOpen,
  onOpenChange,
  shipment,
  packageIndex,
}) => {
  if (!item) return null;

  const isCustom = item.is_custom || item.product_variant_type === 'custom' || !item.product_id;
  const customDesign = item.custom_design;

  const getSizeLabel = (size?: string) => {
    switch (size?.toLowerCase()) {
      case 'small':
        return 'Small (1.5 × 2.0m)';
      case 'medium':
        return 'Medium (1.8 × 2.5m)';
      case 'large':
        return 'Large (2.0 × 3.0m)';
      default:
        return size || 'Standard Dimensions';
    }
  };

  const getJambulLabel = (jambul?: string) => {
    switch (jambul?.toLowerCase()) {
      case 'none':
        return 'No Jambul (Clean Border)';
      case 'top':
        return 'Top Floral Jambul';
      case 'bottom':
        return 'Bottom Floral Jambul';
      case 'both':
        return 'Full Floral Jambul (Top & Bottom)';
      default:
        return jambul || 'Standard Border';
    }
  };

  return (
    <Dialog open={isOpen} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-2xl md:max-w-3xl max-h-[85vh] h-[85vh] p-0 flex flex-col rounded-2xl border border-border bg-card shadow-2xl overflow-hidden gap-0">
        {/* Fixed Opaque Header */}
        <DialogHeader className="p-6 pb-4 pr-14 border-b border-border bg-card shrink-0 space-y-1.5 z-10 text-left">
          <div className="flex items-center gap-2 flex-wrap">
            {isCustom && (
              <Badge
                variant="default"
                className="bg-purple-500/10 text-purple-700 dark:text-purple-300 border-purple-500/20 text-xs font-semibold"
              >
                Custom Flower Board
              </Badge>
            )}
            <Badge
              variant="secondary"
              className="bg-primary/10 text-primary border-primary/20 text-xs font-semibold flex items-center gap-1"
            >
              <Warehouse className="h-3 w-3" />
              {item.shop_name}
            </Badge>
            {item.product_variant_type && item.product_variant_type !== 'standard' && (
              <Badge variant="outline" className="text-xs uppercase font-mono">
                {item.product_variant_type}
              </Badge>
            )}
          </div>
          <DialogTitle className="text-2xl font-bold font-display text-foreground leading-tight">
            {item.product_name}
          </DialogTitle>
          <DialogDescription className="text-xs text-muted-foreground font-mono">
            Item ID: {item.id}
          </DialogDescription>
        </DialogHeader>

        {/* Scrollable Modal Content (Isolated scroll container) */}
        <div className="p-6 space-y-8 flex-1 overflow-y-auto bg-background">
          {/* Section 1: Financial & Quantity Breakdown */}
          <div className="space-y-3">
            <h4 className="text-xs font-bold uppercase tracking-wider text-muted-foreground flex items-center gap-2">
              <DollarSign className="h-4 w-4 text-primary" /> Pricing & Quantity
            </h4>
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
              <div className="p-4 rounded-xl border border-border/60 bg-card space-y-1">
                <span className="text-xs text-muted-foreground font-medium">Unit Price</span>
                <p className="text-base font-bold font-mono text-foreground">{formatCurrency(item.unit_price)}</p>
              </div>
              <div className="p-4 rounded-xl border border-border/60 bg-card space-y-1">
                <span className="text-xs text-muted-foreground font-medium">Quantity Ordered</span>
                <p className="text-base font-bold font-display text-foreground">{item.quantity} units</p>
              </div>
              <div className="p-4 rounded-xl border border-border/60 bg-card space-y-1">
                <span className="text-xs text-muted-foreground font-medium">Line Subtotal</span>
                <p className="text-base font-bold font-mono text-primary">{formatCurrency(item.subtotal)}</p>
              </div>
            </div>
          </div>

          {/* Section 2: Physical Options & Custom Design Specifications */}
          <div className="space-y-4">
            <h4 className="text-xs font-bold uppercase tracking-wider text-muted-foreground flex items-center gap-2">
              <ClipboardPenLine className="h-4 w-4 text-primary" /> Specifications & Customizations
            </h4>

            <div className="space-y-3 p-5 rounded-2xl border border-border/60 bg-card">
              {/* Size & Dimensions */}
              <div className="flex items-start justify-between gap-4 pb-3 border-b border-border/40 text-sm">
                <div className="flex items-center gap-2 text-muted-foreground font-medium">
                  Physical Size
                </div>
                <div className="font-semibold text-foreground text-right">
                  {getSizeLabel(item.item_options?.size)}
                </div>
              </div>

              {/* Jambul / Floral Border Accent */}
              <div className="flex items-start justify-between gap-4 pb-3 border-b border-border/40 text-sm">
                <div className="flex items-center gap-2 text-muted-foreground font-medium">
                  Jambul Floral Style
                </div>
                <div className="font-semibold text-foreground text-right">
                  {getJambulLabel(item.item_options?.jambul)}
                </div>
              </div>

              {/* Product Reference ID */}
              <div className="flex items-start justify-between gap-4 text-sm">
                <div className="flex items-center gap-2 text-muted-foreground font-medium">
                  Product Reference
                </div>
                <div className="font-mono text-xs text-foreground text-right">
                  {item.product_id || 'Custom Atelier Order'}
                </div>
              </div>
            </div>

            {/* Custom Inscription Texts (if present) */}
            {customDesign && (
              <div className="space-y-3 p-5 rounded-2xl border border-purple-500/30 bg-purple-500/5">
                <h5 className="text-xs font-bold uppercase tracking-wider text-purple-700 dark:text-purple-300 flex items-center gap-2">
                  <Type className="h-4 w-4" /> Board Inscriptions & Custom Greetings
                </h5>
                <div className="space-y-3 text-sm pt-1">
                  {(customDesign.header_text_upper || customDesign.body_text_upper) && (
                    <div className="p-3.5 rounded-xl bg-background border border-border/40 space-y-1">
                      <span className="text-xs font-bold text-muted-foreground uppercase">Top Section Inscription</span>
                      {customDesign.header_text_upper && (
                        <p className="font-bold text-foreground">{customDesign.header_text_upper}</p>
                      )}
                      {customDesign.body_text_upper && (
                        <p className="text-muted-foreground text-xs leading-relaxed">{customDesign.body_text_upper}</p>
                      )}
                    </div>
                  )}

                  {(customDesign.header_text_lower || customDesign.body_text_lower) && (
                    <div className="p-3.5 rounded-xl bg-background border border-border/40 space-y-1">
                      <span className="text-xs font-bold text-muted-foreground uppercase">Bottom Section Inscription</span>
                      {customDesign.header_text_lower && (
                        <p className="font-bold text-foreground">{customDesign.header_text_lower}</p>
                      )}
                      {customDesign.body_text_lower && (
                        <p className="text-muted-foreground text-xs leading-relaxed">{customDesign.body_text_lower}</p>
                      )}
                    </div>
                  )}

                  {customDesign.preview_url && (
                    <div className="space-y-1 pt-1">
                      <span className="text-xs font-semibold text-muted-foreground">Snapshot Render</span>
                      <div className="rounded-xl overflow-hidden border border-border/60 bg-background max-h-48 flex items-center justify-center p-2">
                        <img
                          src={customDesign.preview_url}
                          alt="Custom flower board preview"
                          className="max-h-44 object-contain rounded-lg"
                        />
                      </div>
                    </div>
                  )}
                </div>
              </div>
            )}
          </div>

          {/* Section 3: Shop Warehouse & Logistics Assignment */}
          <div className="space-y-4">
            <h4 className="text-xs font-bold uppercase tracking-wider text-muted-foreground flex items-center gap-2">
              <Warehouse className="h-4 w-4 text-primary" /> Fulfillment & Logistics Allocation
            </h4>

            <div className="space-y-3.5 p-5 rounded-2xl border border-border/60 bg-card">
              {/* Shop Branch */}
              <div className="flex items-start justify-between gap-4 pb-3 border-b border-border/40 text-sm">
                <span className="text-muted-foreground font-medium">Assigned Store Branch</span>
                <span className="font-semibold text-foreground text-right">{item.shop_name}</span>
              </div>

              {/* Requested Courier */}
              <div className="flex items-start justify-between gap-4 pb-3 border-b border-border/40 text-sm">
                <span className="text-muted-foreground font-medium">Customer Selected Courier</span>
                <span className="font-semibold uppercase text-foreground text-right">
                  {item.courier_code ? `${item.courier_code} (${item.courier_service || 'Standard'})` : 'Default Store Carrier'}
                </span>
              </div>

              {/* Package Assignment Status */}
              <div className="flex items-start justify-between gap-4 text-sm">
                <span className="text-muted-foreground font-medium">Shipping Package</span>
                <div className="text-right">
                  {shipment ? (
                    <div className="space-y-0.5">
                      <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-bold bg-primary/10 text-primary">
                        <Truck className="h-3.5 w-3.5" /> Package #{packageIndex !== undefined ? packageIndex + 1 : 1} ({shipment.courier.toUpperCase()})
                      </span>
                      <p className="text-xs font-mono text-muted-foreground">
                        Waybill: <strong className="text-foreground">{shipment.tracking_number || 'Self Delivery'}</strong>
                      </p>
                    </div>
                  ) : (
                    <span className="text-amber-600 dark:text-amber-400 font-semibold text-xs inline-flex items-center gap-1">
                      <Info className="h-3.5 w-3.5" /> Ready for Package Allocation
                    </span>
                  )}
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Fixed Opaque Footer */}
        <DialogFooter className="p-4 sm:p-5 border-t border-border bg-card shrink-0 flex flex-row items-center justify-end z-10">
          <Button
            type="button"
            variant="outline"
            onClick={() => onOpenChange(false)}
            className="rounded-xl text-sm font-semibold h-10 px-6 cursor-pointer hover:bg-muted"
          >
            Close Inspection
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};

export default OrderItemDetailOverlay;
