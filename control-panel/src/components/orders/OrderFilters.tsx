import React from 'react';
import { Warehouse, RefreshCw, Calendar, X } from 'lucide-react';
import { Button } from '../ui/button';
import SearchInput from '../SearchInput';

export interface ShopOption {
  id: string;
  name: string;
  slug?: string;
}

export interface OrderFiltersProps {
  searchNumber: string;
  statusFilter: string;
  shopFilter?: string;
  fromDate?: string;
  toDate?: string;
  shopsList?: ShopOption[];
  isAdmin?: boolean;
  isShopLocked?: boolean;
  hideShopFilter?: boolean;
  loading?: boolean;
  isSwitchingCategory?: boolean;
  onSearchChange: (value: string) => void;
  onStatusChange: (status: string) => void;
  onShopChange?: (shop: string) => void;
  onFromDateChange?: (date: string) => void;
  onToDateChange?: (date: string) => void;
  onRefresh: () => void;
}

const statusTabs = [
  { value: 'all', label: 'All Orders' },
  { value: 'pending', label: 'Awaiting Payment' },
  { value: 'confirmed', label: 'To Process' },
  { value: 'processing', label: 'In Packaging' },
  { value: 'shipped', label: 'Shipped' },
  { value: 'delivered', label: 'Delivered' },
  { value: 'cancelled', label: 'Cancelled' },
];

export const OrderFilters: React.FC<OrderFiltersProps> = ({
  searchNumber,
  statusFilter,
  shopFilter = '',
  fromDate = '',
  toDate = '',
  shopsList = [],
  isAdmin = false,
  isShopLocked = false,
  hideShopFilter = false,
  loading = false,
  isSwitchingCategory = false,
  onSearchChange,
  onStatusChange,
  onShopChange,
  onFromDateChange,
  onToDateChange,
  onRefresh,
}) => {
  const showShopFilter = !hideShopFilter && !isShopLocked && shopsList.length > 0;

  return (
    <div className="space-y-4 w-full">
      {/* Assigned Shops Segmented Switcher */}
      {showShopFilter && (
        <div className="flex flex-col gap-2.5 p-4 rounded-2xl border border-border/60 bg-muted/20">
          <div className="flex items-center justify-between">
            <span className="text-xs font-semibold uppercase tracking-wider text-muted-foreground flex items-center gap-1.5">
              <Warehouse className="w-3.5 h-3.5 text-primary" />
              {isAdmin ? 'Shops Overview' : 'Your Assigned Shops'}
            </span>
            <span className="text-xs text-muted-foreground font-medium">
              {shopsList.length} {shopsList.length === 1 ? 'Shop Assigned' : 'Shops Assigned'}
            </span>
          </div>
          <div className="flex items-center gap-2 overflow-x-auto pt-0.5 pb-0.5 scrollbar-none">
            {isAdmin && (
              <button
                type="button"
                onClick={() => onShopChange?.('all')}
                className={`px-3.5 py-1.5 rounded-xl text-xs font-medium transition-all shrink-0 border cursor-pointer ${
                  (shopFilter || 'all') === 'all'
                    ? 'bg-primary text-primary-foreground border-primary shadow-sm'
                    : 'bg-background hover:bg-muted text-muted-foreground border-border/80'
                }`}
              >
                All Shops
              </button>
            )}
            {shopsList.map((shop) => {
              const shopKey = shop.slug || shop.id;
              const isSelected =
                shopFilter === shopKey ||
                (!shopFilter && shopsList.length === 1 && shopKey === (shopsList[0].slug || shopsList[0].id));
              return (
                <button
                  key={shop.id}
                  type="button"
                  onClick={() => onShopChange?.(shopKey)}
                  className={`px-3.5 py-1.5 rounded-xl text-xs font-medium transition-all shrink-0 border flex items-center gap-2 cursor-pointer ${
                    isSelected
                      ? 'bg-primary text-primary-foreground border-primary shadow-sm'
                      : 'bg-background hover:bg-muted text-muted-foreground border-border/80'
                  }`}
                >
                  <Warehouse className="w-3.5 h-3.5 text-current shrink-0" />
                  <span>{shop.name}</span>
                  {shop.slug && (
                    <span
                      className={`text-[10px] px-1.5 py-0.5 rounded-md font-mono ${
                        isSelected
                          ? 'bg-primary-foreground/20 text-primary-foreground'
                          : 'bg-muted text-muted-foreground'
                      }`}
                    >
                      {shop.slug}
                    </span>
                  )}
                </button>
              );
            })}
          </div>
        </div>
      )}

      {/* Search, Date Range and Refresh Controls */}
      <div className="flex flex-col md:flex-row items-stretch md:items-center justify-between gap-3 w-full">
        {/* Search Bar */}
        <div className="flex-1 max-w-sm">
          <SearchInput
            value={searchNumber}
            onChange={onSearchChange}
            placeholder="Search by Order Number..."
          />
        </div>

        {/* Date Range Inputs and Refresh Button */}
        <div className="flex flex-wrap items-center gap-2 justify-end">
          <div className="flex items-center gap-1.5 bg-background border border-border/70 rounded-xl px-2 py-1 text-xs text-muted-foreground shadow-sm">
            <Calendar className="h-3.5 w-3.5 text-muted-foreground shrink-0" />
            <span className="text-[11px] font-medium">From</span>
            <input
              type="date"
              value={fromDate}
              onChange={(e) => onFromDateChange?.(e.target.value)}
              className="bg-transparent text-foreground text-xs focus:outline-none cursor-pointer"
            />
            {fromDate && (
              <button
                type="button"
                onClick={() => onFromDateChange?.('')}
                className="text-muted-foreground hover:text-foreground cursor-pointer"
                title="Clear from date"
              >
                <X className="h-3 w-3" />
              </button>
            )}
          </div>

          <div className="flex items-center gap-1.5 bg-background border border-border/70 rounded-xl px-2 py-1 text-xs text-muted-foreground shadow-sm">
            <Calendar className="h-3.5 w-3.5 text-muted-foreground shrink-0" />
            <span className="text-[11px] font-medium">To</span>
            <input
              type="date"
              value={toDate}
              onChange={(e) => onToDateChange?.(e.target.value)}
              className="bg-transparent text-foreground text-xs focus:outline-none cursor-pointer"
            />
            {toDate && (
              <button
                type="button"
                onClick={() => onToDateChange?.('')}
                className="text-muted-foreground hover:text-foreground cursor-pointer"
                title="Clear to date"
              >
                <X className="h-3 w-3" />
              </button>
            )}
          </div>

          <Button
            variant="outline"
            onClick={onRefresh}
            disabled={loading}
            size="sm"
            className="rounded-xl h-8 text-xs cursor-pointer"
          >
            <RefreshCw className={`h-3.5 w-3.5 mr-1.5 ${loading ? 'animate-spin' : ''}`} />
            Refresh
          </Button>
        </div>
      </div>

      {/* Status Filter Tabs */}
      <div className="flex gap-2 overflow-x-auto pb-1 w-full border-b border-border/40 scrollbar-none">
        {statusTabs.map((tab) => {
          const isActive = (statusFilter || 'all') === tab.value;
          return (
            <button
              key={tab.value}
              disabled={isSwitchingCategory || loading}
              onClick={() => onStatusChange(tab.value === 'all' ? '' : tab.value)}
              className={`px-3.5 py-1.5 rounded-lg text-xs font-semibold whitespace-nowrap transition-all ${
                isActive
                  ? 'bg-primary text-primary-foreground shadow-sm'
                  : 'text-muted-foreground hover:text-foreground hover:bg-muted/40'
              } ${
                isSwitchingCategory || loading
                  ? 'opacity-60 cursor-not-allowed pointer-events-none'
                  : 'cursor-pointer'
              }`}
            >
              {tab.label}
            </button>
          );
        })}
      </div>
    </div>
  );
};

export default OrderFilters;
