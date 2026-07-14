import { useState, Fragment } from 'react';
import { PackageOpen, ArrowUpDown, MoreHorizontal, ChevronDown, ChevronUp } from 'lucide-react';
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
import {
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuItem,
} from '../../components/ui/dropdown-menu';
import { useOrdersViewModel } from '../../viewmodels/useOrdersViewModel';
import EmptyState from '../../components/EmptyState';
import SearchInput from '../../components/SearchInput';
import StatusBadge from '../../components/StatusBadge';
import Pagination from '../../components/Pagination';
import { Skeleton } from '../../components/ui/skeleton';

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

  const [expandedOrderId, setExpandedOrderId] = useState<string | null>(null);

  const toggleExpandOrder = (id: string) => {
    setExpandedOrderId(prev => prev === id ? null : id);
  };

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
      <div className="flex-1 space-y-8 p-6 sm:p-8 lg:p-12 animate-in fade-in duration-300">
        <div className="flex items-center justify-between space-y-2">
          <div>
            <h2 className="text-3xl font-bold font-display tracking-tight text-foreground">Orders</h2>
            <p className="text-muted-foreground text-sm">
              Manage your orders and fulfillments
            </p>
          </div>
        </div>

        <Card className="border-0 shadow-none bg-zinc-50/40 dark:bg-slate-900/40">
          <CardHeader>
            <CardTitle className="font-bold font-display tracking-tight text-lg text-foreground">All Orders</CardTitle>
            <CardDescription className="text-muted-foreground text-sm">
              {data?.total ? `Found ${data.total} orders matching your criteria.` : 'View and manage orders.'}
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="mb-4 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
              {/* Left Side: Filter and Search */}
              <div className="flex flex-1 flex-col sm:flex-row items-center gap-4 w-full sm:w-auto">
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
                    <SelectTrigger className="w-full rounded-xl border border-border bg-background">
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
              </div>

              {/* Right Side: Refresh */}
              <div className="flex items-center gap-2 justify-end w-full sm:w-auto">
                <Button variant="outline" className="border-border text-foreground hover:text-primary hover:bg-primary/5 rounded-xl" onClick={() => refresh()}>
                  Refresh
                </Button>
              </div>
            </div>

            <div className="rounded-2xl border border-border overflow-hidden">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead 
                      className="cursor-pointer hover:bg-muted/50 text-foreground" 
                      onClick={() => handleSort('number')}
                    >
                      <div className="flex items-center">
                        Order Number
                        <ArrowUpDown className="ml-2 h-4 w-4" />
                      </div>
                    </TableHead>
                    <TableHead 
                      className="cursor-pointer hover:bg-muted/50 text-foreground"
                      onClick={() => handleSort('date')}
                    >
                      <div className="flex items-center">
                        Date
                        <ArrowUpDown className="ml-2 h-4 w-4" />
                      </div>
                    </TableHead>
                    <TableHead 
                      className="cursor-pointer hover:bg-muted/50 text-foreground"
                      onClick={() => handleSort('status')}
                    >
                      <div className="flex items-center">
                        Status
                        <ArrowUpDown className="ml-2 h-4 w-4" />
                      </div>
                    </TableHead>
                    <TableHead className="text-foreground">Payment Status</TableHead>
                    <TableHead className="text-foreground">Shipment Status</TableHead>
                    <TableHead 
                      className="text-right cursor-pointer hover:bg-muted/50 text-foreground"
                      onClick={() => handleSort('total')}
                    >
                      <div className="flex items-center justify-end">
                        Total
                        <ArrowUpDown className="ml-2 h-4 w-4" />
                      </div>
                    </TableHead>
                    <TableHead className="w-[80px]"></TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {loading ? (
                    Array.from({ length: 5 }).map((_, i) => (
                      <TableRow key={`skeleton-${i}`}>
                        <TableCell><Skeleton className="h-5 w-28 animate-pulse bg-muted" /></TableCell>
                        <TableCell><Skeleton className="h-5 w-20 animate-pulse bg-muted" /></TableCell>
                        <TableCell><Skeleton className="h-5 w-16 animate-pulse bg-muted" /></TableCell>
                        <TableCell><Skeleton className="h-5 w-16 animate-pulse bg-muted" /></TableCell>
                        <TableCell><Skeleton className="h-5 w-16 animate-pulse bg-muted" /></TableCell>
                        <TableCell className="text-right"><Skeleton className="h-5 w-20 ml-auto animate-pulse bg-muted" /></TableCell>
                        <TableCell><Skeleton className="h-8 w-8 rounded-xl ml-auto animate-pulse bg-muted" /></TableCell>
                      </TableRow>
                    ))
                  ) : error ? (
                    <TableRow>
                      <TableCell colSpan={7} className="p-0">
                        <EmptyState 
                          title="Failed to load orders" 
                          description={error} 
                          className="flex h-48 flex-col items-center justify-center text-center p-4 gap-2 border-0 bg-transparent text-destructive"
                        />
                      </TableCell>
                    </TableRow>
                  ) : !data?.orders || data.orders.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={7} className="p-0">
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
                      <Fragment key={order.id}>
                        <TableRow 
                          onClick={() => toggleExpandOrder(order.id)}
                          className={`cursor-pointer transition-colors ${
                            expandedOrderId === order.id 
                              ? 'bg-primary/5 hover:bg-primary/5 border-l-4 border-l-primary border-b-0' 
                              : 'hover:bg-muted/50'
                          }`}
                        >
                          <TableCell className="font-medium">
                            {order.number}
                          </TableCell>
                          <TableCell>
                            {new Date(order.created_at).toLocaleDateString()}
                          </TableCell>
                          <TableCell>
                            <StatusBadge status={order.status} />
                          </TableCell>
                          <TableCell>
                            <StatusBadge status={order.payment?.status || 'N/A'} />
                          </TableCell>
                          <TableCell>
                            <StatusBadge status={order.shipment?.status || 'N/A'} />
                          </TableCell>
                          <TableCell className="text-right font-medium">
                            {formatCurrency(order.total)}
                          </TableCell>
                          <TableCell className="text-right" onClick={(e) => e.stopPropagation()}>
                            <DropdownMenu>
                              <DropdownMenuTrigger asChild>
                                <Button variant="ghost" size="icon" className="h-8 w-8 rounded-xl hover:bg-muted/80">
                                  <MoreHorizontal className="h-4 w-4" />
                                </Button>
                              </DropdownMenuTrigger>
                              <DropdownMenuContent align="end" className="min-w-[150px] p-1 rounded-xl">
                                <DropdownMenuItem
                                  onClick={() => toggleExpandOrder(order.id)}
                                  className="cursor-pointer flex items-center gap-2 px-2.5 py-1.5 text-sm rounded-lg hover:bg-muted"
                                >
                                  {expandedOrderId === order.id ? (
                                    <>
                                      <ChevronUp className="h-4 w-4 text-muted-foreground" />
                                      <span>Collapse Details</span>
                                    </>
                                  ) : (
                                    <>
                                      <ChevronDown className="h-4 w-4 text-muted-foreground" />
                                      <span>View Details</span>
                                    </>
                                  )}
                                </DropdownMenuItem>
                              </DropdownMenuContent>
                            </DropdownMenu>
                          </TableCell>
                        </TableRow>
                        {expandedOrderId === order.id && (
                          <TableRow className="bg-zinc-100/55 dark:bg-slate-900/80 hover:bg-zinc-100/55 dark:hover:bg-slate-900/80 border-t-0 border-l-4 border-l-primary border-b border-border/80">
                            <TableCell colSpan={7} className="p-6">
                              <div className="space-y-6 animate-in slide-in-from-top-2 duration-200">
                                <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b pb-4">
                                  <div>
                                    <h4 className="text-lg font-bold text-foreground">Order Details</h4>
                                    <p className="text-xs text-muted-foreground">
                                      Placed on {new Date(order.created_at).toLocaleString('id-ID', { dateStyle: 'medium', timeStyle: 'short' })}
                                    </p>
                                  </div>
                                  <div className="flex gap-2">
                                    <StatusBadge status={order.status} className="h-6" />
                                    {order.payment && <StatusBadge status={`Payment: ${order.payment.status}`} className="h-6" />}
                                    {order.shipment && <StatusBadge status={`Shipment: ${order.shipment.status}`} className="h-6" />}
                                  </div>
                                </div>

                                <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                                  {/* Items Table */}
                                  <div className="md:col-span-2 space-y-3">
                                    <h5 className="text-sm font-semibold text-foreground">Order Items</h5>
                                    <div className="rounded-xl border border-border overflow-hidden bg-background">
                                      <Table>
                                        <TableHeader className="bg-muted/40">
                                          <TableRow>
                                            <TableHead>Item</TableHead>
                                            <TableHead className="text-right">Qty</TableHead>
                                            <TableHead className="text-right">Price</TableHead>
                                          </TableRow>
                                        </TableHeader>
                                        <TableBody>
                                          {order.items.map((item) => (
                                            <TableRow key={item.id}>
                                              <TableCell>
                                                <div className="font-semibold text-foreground text-sm">{item.product_name}</div>
                                                <div className="text-xs text-muted-foreground">Shop: {item.shop_name}</div>
                                                {item.courier_code && (
                                                  <div className="text-xs text-muted-foreground uppercase">
                                                    Courier: {item.courier_code} {item.courier_service}
                                                  </div>
                                                )}
                                              </TableCell>
                                              <TableCell className="text-right text-sm font-medium">{item.quantity}</TableCell>
                                              <TableCell className="text-right text-sm font-medium">{formatCurrency(item.subtotal)}</TableCell>
                                            </TableRow>
                                          ))}
                                        </TableBody>
                                      </Table>
                                    </div>
                                  </div>

                                  {/* Pricing and metadata */}
                                  <div className="space-y-6">
                                    <div className="bg-zinc-50 dark:bg-slate-900/60 p-5 rounded-2xl border border-border/60 space-y-3">
                                      <h5 className="text-sm font-bold text-foreground uppercase tracking-wider">Payment Summary</h5>
                                      <div className="space-y-2 text-sm">
                                        <div className="flex justify-between text-muted-foreground">
                                          <span>Subtotal</span>
                                          <span className="font-medium text-foreground">{formatCurrency(order.subtotal)}</span>
                                        </div>
                                        <div className="flex justify-between text-muted-foreground">
                                          <span>Shipping Fee</span>
                                          <span className="font-medium text-foreground">{formatCurrency(order.shipping_fee)}</span>
                                        </div>
                                        <div className="flex justify-between border-t border-border pt-3 mt-3 text-base font-bold text-foreground">
                                          <span>Total Amount</span>
                                          <span className="text-primary">{formatCurrency(order.total)}</span>
                                        </div>
                                      </div>
                                    </div>

                                    <div className="bg-zinc-50 dark:bg-slate-900/60 p-5 rounded-2xl border border-border/60 space-y-3">
                                      <h5 className="text-sm font-bold text-foreground uppercase tracking-wider">Customer & Shipping</h5>
                                      <div className="text-xs space-y-2 text-muted-foreground">
                                        <p><span className="font-semibold text-foreground">User ID:</span> <span className="font-mono">{order.user_id}</span></p>
                                        <p><span className="font-semibold text-foreground">Address ID:</span> <span className="font-mono">{order.address_id}</span></p>
                                        {order.shipment?.tracking_number && (
                                          <p><span className="font-semibold text-foreground">Tracking Number:</span> <span className="font-mono text-primary font-semibold">{order.shipment.tracking_number}</span></p>
                                        )}
                                      </div>
                                    </div>
                                  </div>
                                </div>
                              </div>
                            </TableCell>
                          </TableRow>
                        )}
                      </Fragment>
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
    </div>
  );
}
