import { useState } from 'react';
import { PackageOpen, Search, Loader2, ChevronLeft, ChevronRight, ArrowUpDown, Eye } from 'lucide-react';
import { Button } from '../../components/ui/button';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '../../components/ui/table';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '../../components/ui/card';
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetDescription } from '../../components/ui/sheet';
import { useOrdersViewModel } from '../../viewmodels/useOrdersViewModel';
import type { Order } from '../../models/Order';
import LoadingState from '../../components/LoadingState';
import EmptyState from '../../components/EmptyState';
import SearchInput from '../../components/SearchInput';
import StatusBadge from '../../components/StatusBadge';
import Pagination from '../../components/Pagination';

import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '../../components/ui/select';

export default function OrdersPage() {
  const {
    data,
    loading,
    error,
    page,
    limit,
    sort,
    searchNumber,
    statusFilter,
    setPage,
    setSort,
    setSearchNumber,
    setStatusFilter,
    refresh
  } = useOrdersViewModel();

  const [selectedOrder, setSelectedOrder] = useState<Order | null>(null);

  const handleSort = (field: string) => {
    let newDirection = 'desc';
    const [currentField, currentDirection] = sort.split(':');
    if (currentField === field && currentDirection === 'desc') {
      newDirection = 'asc';
    }
    setSort(`${field}:${newDirection}`);
    setPage(1);
  };

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('id-ID', {
      style: 'currency',
      currency: 'IDR',
      minimumFractionDigits: 0,
    }).format(amount);
  };

  return (
    <div className="flex-col md:flex">
      <div className="flex-1 space-y-4 p-8 pt-6">
        <div className="flex items-center justify-between space-y-2">
          <div>
            <h2 className="text-3xl font-bold tracking-tight">Orders</h2>
            <p className="text-muted-foreground">
              Manage your orders and fulfillments
            </p>
          </div>
        </div>

        <Card>
          <CardHeader>
            <CardTitle>All Orders</CardTitle>
            <CardDescription>
              {data?.total ? `Found ${data.total} orders matching your criteria.` : 'View and manage orders.'}
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="mb-4 flex flex-col sm:flex-row items-center gap-4">
              <SearchInput
                value={searchNumber}
                onChange={(val) => {
                  setSearchNumber(val);
                  setPage(1);
                }}
                placeholder="Search by Order Number..."
              />
              <div className="w-full sm:w-[180px]">
                <Select
                  value={statusFilter || "all"}
                  onValueChange={(val) => {
                    setStatusFilter(val === "all" ? "" : val);
                    setPage(1);
                  }}
                >
                  <SelectTrigger className="w-full">
                    <SelectValue placeholder="All Statuses" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">All Statuses</SelectItem>
                    <SelectItem value="pending">Pending</SelectItem>
                    <SelectItem value="confirmed">Confirmed</SelectItem>
                    <SelectItem value="processing">Processing</SelectItem>
                    <SelectItem value="shipped">Shipped</SelectItem>
                    <SelectItem value="delivered">Delivered</SelectItem>
                    <SelectItem value="cancelled">Cancelled</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <Button variant="outline" onClick={() => refresh()}>
                Refresh
              </Button>
            </div>

            <div className="rounded-md border">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead 
                      className="cursor-pointer hover:bg-slate-50" 
                      onClick={() => handleSort('number')}
                    >
                      <div className="flex items-center">
                        Order Number
                        <ArrowUpDown className="ml-2 h-4 w-4" />
                      </div>
                    </TableHead>
                    <TableHead 
                      className="cursor-pointer hover:bg-slate-50"
                      onClick={() => handleSort('date')}
                    >
                      <div className="flex items-center">
                        Date
                        <ArrowUpDown className="ml-2 h-4 w-4" />
                      </div>
                    </TableHead>
                    <TableHead 
                      className="cursor-pointer hover:bg-slate-50"
                      onClick={() => handleSort('status')}
                    >
                      <div className="flex items-center">
                        Status
                        <ArrowUpDown className="ml-2 h-4 w-4" />
                      </div>
                    </TableHead>
                    <TableHead 
                      className="text-right cursor-pointer hover:bg-slate-50"
                      onClick={() => handleSort('total')}
                    >
                      <div className="flex items-center justify-end">
                        Total
                        <ArrowUpDown className="ml-2 h-4 w-4" />
                      </div>
                    </TableHead>
                    <TableHead className="w-[100px]"></TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {loading ? (
                    <TableRow>
                      <TableCell colSpan={5} className="p-0">
                        <LoadingState message="Loading orders..." className="flex h-48 flex-col items-center justify-center gap-2" />
                      </TableCell>
                    </TableRow>
                  ) : error ? (
                    <TableRow>
                      <TableCell colSpan={5} className="p-0">
                        <EmptyState 
                          title="Failed to load orders" 
                          description={error} 
                          className="flex h-48 flex-col items-center justify-center text-center p-4 gap-2 border-0 bg-transparent text-destructive"
                        />
                      </TableCell>
                    </TableRow>
                  ) : !data?.orders || data.orders.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={5} className="p-0">
                        <EmptyState 
                          icon={<PackageOpen className="h-8 w-8 mb-2 mx-auto text-slate-400" />} 
                          title="No orders found" 
                          description="No orders found matching your criteria."
                          className="flex h-48 flex-col items-center justify-center text-center p-4 gap-1.5 border-0 bg-transparent"
                        />
                      </TableCell>
                    </TableRow>
                  ) : (
                    data.orders.map((order) => (
                      <TableRow key={order.id}>
                        <TableCell className="font-medium">
                          {order.number}
                        </TableCell>
                        <TableCell>
                          {new Date(order.created_at).toLocaleDateString()}
                        </TableCell>
                        <TableCell>
                          <StatusBadge status={order.status} />
                        </TableCell>
                        <TableCell className="text-right font-medium">
                          {formatCurrency(order.total)}
                        </TableCell>
                        <TableCell>
                          <Button 
                            variant="ghost" 
                            size="sm" 
                            onClick={() => setSelectedOrder(order)}
                          >
                            <Eye className="h-4 w-4 mr-2" />
                            View
                          </Button>
                        </TableCell>
                      </TableRow>
                    ))
                  )}
                </TableBody>
              </Table>
            </div>

            <Pagination
              currentPage={page}
              totalPages={Math.ceil((data?.total || 0) / limit)}
              totalItems={data?.total || 0}
              limit={limit}
              onPageChange={setPage}
              itemNamePlural="orders"
            />
          </CardContent>
        </Card>
      </div>

      {/* Order Details Sheet/Modal */}
      <Sheet open={!!selectedOrder} onOpenChange={(open) => !open && setSelectedOrder(null)}>
        <SheetContent side="right" className="w-[400px] sm:w-[540px] overflow-y-auto">
          {selectedOrder && (
            <>
              <SheetHeader className="mb-6">
                <SheetTitle className="flex items-center justify-between">
                  <span>Order {selectedOrder.number}</span>
                  <StatusBadge status={selectedOrder.status} />
                </SheetTitle>
                <SheetDescription>
                  Placed on {new Date(selectedOrder.created_at).toLocaleString()}
                </SheetDescription>
              </SheetHeader>

              <div className="space-y-6">
                <div>
                  <h4 className="text-sm font-semibold mb-3">Order Items</h4>
                  <div className="rounded-md border overflow-hidden">
                    <Table>
                      <TableHeader className="bg-slate-50">
                        <TableRow>
                          <TableHead>Item</TableHead>
                          <TableHead className="text-right">Qty</TableHead>
                          <TableHead className="text-right">Price</TableHead>
                        </TableRow>
                      </TableHeader>
                      <TableBody>
                        {selectedOrder.items.map((item) => (
                          <TableRow key={item.id}>
                            <TableCell>
                              <div className="font-medium text-sm">{item.product_name}</div>
                              <div className="text-xs text-muted-foreground">Shop: {item.shop_name}</div>
                              {item.courier_code && (
                                <div className="text-xs text-muted-foreground uppercase">
                                  Courier: {item.courier_code} {item.courier_service}
                                </div>
                              )}
                            </TableCell>
                            <TableCell className="text-right text-sm">{item.quantity}</TableCell>
                            <TableCell className="text-right text-sm">{formatCurrency(item.subtotal)}</TableCell>
                          </TableRow>
                        ))}
                      </TableBody>
                    </Table>
                  </div>
                </div>

                <div className="bg-slate-50 p-4 rounded-lg space-y-2">
                  <div className="flex justify-between text-sm">
                    <span className="text-muted-foreground">Subtotal</span>
                    <span>{formatCurrency(selectedOrder.subtotal)}</span>
                  </div>
                  <div className="flex justify-between text-sm">
                    <span className="text-muted-foreground">Shipping Fee</span>
                    <span>{formatCurrency(selectedOrder.shipping_fee)}</span>
                  </div>
                  <div className="flex justify-between font-semibold border-t pt-2 mt-2">
                    <span>Total</span>
                    <span>{formatCurrency(selectedOrder.total)}</span>
                  </div>
                </div>

                <div>
                  <h4 className="text-sm font-semibold mb-2">Customer & Delivery</h4>
                  <div className="text-sm space-y-1 text-slate-600 bg-slate-50 p-4 rounded-lg">
                    <p><span className="font-medium">User ID:</span> {selectedOrder.user_id}</p>
                    <p><span className="font-medium">Address ID:</span> {selectedOrder.address_id}</p>
                  </div>
                </div>
              </div>
            </>
          )}
        </SheetContent>
      </Sheet>
    </div>
  );
}
