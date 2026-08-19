import { useOrdersViewModel } from './useOrdersViewModel';

export function useShopOrdersViewModel(shopId: string) {
  return useOrdersViewModel(shopId);
}
