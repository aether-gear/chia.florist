import React from 'react';
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetDescription } from '../ui/sheet';
import type { Order } from '../../models/Order';
import type { ShipmentDispatchPayload, OrderTrackingResponse } from '../../viewmodels/useOrderActionsViewModel';
import OrderDetailInspector from './OrderDetailInspector';

export interface OrderDetailSheetProps {
  order: Order | null;
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
  submitting: boolean;
  onStartProcessing: (orderId: string) => Promise<void>;
  onDispatchOrder: (orderId: string, shipments: ShipmentDispatchPayload[]) => Promise<void>;
  onUpdateShipmentStatus: (shipmentId: string, status: string) => Promise<void>;
  onUpdateWaybill: (
    shipmentId: string,
    details: { tracking_number?: string; courier?: string; service?: string }
  ) => Promise<void>;
  fetchOrderTracking: (orderId: string) => Promise<OrderTrackingResponse | null>;
}

export const OrderDetailSheet: React.FC<OrderDetailSheetProps> = ({
  order,
  isOpen,
  onOpenChange,
  submitting,
  onStartProcessing,
  onDispatchOrder,
  onUpdateShipmentStatus,
  onUpdateWaybill,
  fetchOrderTracking,
}) => {
  return (
    <Sheet open={isOpen} onOpenChange={onOpenChange}>
      <SheetContent side="right" className="w-full sm:max-w-2xl md:max-w-3xl overflow-y-auto p-0 border-l border-border">
        <SheetHeader className="sr-only">
          <SheetTitle>Order Details</SheetTitle>
          <SheetDescription>Inspect and manage order fulfillment details</SheetDescription>
        </SheetHeader>
        {order && (
          <OrderDetailInspector
            order={order}
            submitting={submitting}
            onClose={() => onOpenChange(false)}
            onStartProcessing={onStartProcessing}
            onDispatchOrder={onDispatchOrder}
            onUpdateShipmentStatus={onUpdateShipmentStatus}
            onUpdateWaybill={onUpdateWaybill}
            fetchOrderTracking={fetchOrderTracking}
          />
        )}
      </SheetContent>
    </Sheet>
  );
};

export default OrderDetailSheet;
