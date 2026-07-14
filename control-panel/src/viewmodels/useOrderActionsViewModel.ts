import { useState } from 'react';
import { fetchApi } from '../lib/api';
import { useToast } from '../hooks/use-toast';

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

  const confirmPayment = async (paymentId: string) => {
    return handleAction(
      `/payments/${paymentId}/action`,
      {
        method: 'POST',
        body: JSON.stringify({ action: 'confirm' }),
      },
      'Payment confirmed successfully.',
      'Failed to confirm payment.'
    );
  };

  const rejectPayment = async (paymentId: string) => {
    return handleAction(
      `/payments/${paymentId}/action`,
      {
        method: 'POST',
        body: JSON.stringify({ action: 'reject' }),
      },
      'Payment rejected successfully.',
      'Failed to reject payment.'
    );
  };

  const updateOrderStatus = async (
    orderId: string,
    status: string,
    trackingNumber?: string,
    fulfillmentMethod?: string
  ) => {
    const payload: Record<string, any> = { status };
    if (trackingNumber !== undefined) payload.tracking_number = trackingNumber;
    if (fulfillmentMethod !== undefined) payload.fulfillment_method = fulfillmentMethod;

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
    confirmPayment,
    rejectPayment,
    updateOrderStatus,
    updateShipmentStatus,
    updateShipmentDetails,
  };
}
