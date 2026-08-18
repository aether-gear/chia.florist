import { useState } from 'react';
import { fetchApi } from '../lib/api';
import { useToast } from '../hooks/use-toast';

export interface ShipmentDispatchPayload {
  fulfillment_method: string;
  courier: string;
  service: string;
  tracking_number?: string;
  item_ids: string[];
}

export interface TrackingTimelineEvent {
  status: string;
  description: string;
  location: string;
  timestamp: string;
}

export interface OrderTrackingResponse {
  order_id: string;
  shipment_id: string;
  courier: string;
  tracking_number?: string;
  warning?: string;
  timeline: TrackingTimelineEvent[];
}

export function useOrderActionsViewModel() {
  const [submitting, setSubmitting] = useState<boolean>(false);
  const { toast } = useToast();

  const handleAction = async (
    endpoint: string,
    options: RequestInit,
    successMessage: string,
    errorMessage: string
  ) => {
    setSubmitting(true);
    try {
      const response = await fetchApi(endpoint, options);
      toast({
        title: 'Success',
        description: successMessage,
      });
      return response;
    } catch (err: any) {
      console.error(err);
      toast({
        title: 'Error',
        description: err.message || errorMessage,
        variant: 'destructive',
      });
      throw err;
    } finally {
      setSubmitting(false);
    }
  };

  const fetchOrderTracking = async (orderId: string): Promise<OrderTrackingResponse | null> => {
    try {
      const response = await fetchApi(`/orders/${orderId}/tracking`);
      return response;
    } catch (err: any) {
      console.error('Failed to fetch order tracking', err);
      return null;
    }
  };

  const updateOrderStatus = async (
    orderId: string,
    status: string,
    trackingNumber?: string,
    fulfillmentMethod?: string,
    shipments?: ShipmentDispatchPayload[]
  ) => {
    const payload: Record<string, any> = { status };
    if (trackingNumber !== undefined) payload.tracking_number = trackingNumber;
    if (fulfillmentMethod !== undefined) payload.fulfillment_method = fulfillmentMethod;
    if (shipments !== undefined && shipments.length > 0) payload.shipments = shipments;

    return handleAction(
      `/orders/${orderId}/status`,
      {
        method: 'PATCH',
        body: JSON.stringify(payload),
      },
      `Order status updated to ${status}.`,
      'Failed to update order status.'
    );
  };

  const updateShipmentStatus = async (shipmentId: string, status: string) => {
    return handleAction(
      `/shipments/${shipmentId}/status`,
      {
        method: 'PATCH',
        body: JSON.stringify({ status }),
      },
      `Shipment status updated to ${status}.`,
      'Failed to update shipment status.'
    );
  };

  const updateShipmentDetails = async (
    shipmentId: string,
    details: { tracking_number?: string; courier?: string; service?: string }
  ) => {
    return handleAction(
      `/shipments/${shipmentId}`,
      {
        method: 'PATCH',
        body: JSON.stringify(details),
      },
      'Shipment details updated successfully.',
      'Failed to update shipment details.'
    );
  };

  return {
    submitting,
    fetchOrderTracking,
    updateOrderStatus,
    updateShipmentStatus,
    updateShipmentDetails,
  };
}
