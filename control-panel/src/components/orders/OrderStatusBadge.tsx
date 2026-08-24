import React from 'react';
import StatusBadge from '../StatusBadge';

interface OrderStatusBadgeProps {
  status: string;
  className?: string;
}

export const OrderStatusBadge: React.FC<OrderStatusBadgeProps> = ({ status, className = '' }) => {
  return <StatusBadge status={status} className={className} />;
};

export default OrderStatusBadge;
