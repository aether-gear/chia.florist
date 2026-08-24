import { useEffect } from 'react';
import { Edit, X, Package, CheckCircle2, XCircle, Loader2, Weight, BarChart2, Sparkles, TrendingUp, AlertTriangle } from 'lucide-react';
import { Button } from '../ui/button';
import { Badge } from '../ui/badge';
import StatusBadge from '../StatusBadge';
import { useProductFormViewModel } from '../../viewmodels/useProductFormViewModel';
import { useDemandForecastViewModel } from '../../viewmodels/useDemandForecastViewModel';
import { Skeleton } from '../ui/skeleton';

interface ProductDetailsViewProps {
  productSlug: string;
  onClose: () => void;
  onEditProduct: (slug: string) => void;
}

export default function ProductDetailsView({
  productSlug,
  onClose,
  onEditProduct,
}: ProductDetailsViewProps) {
  const { product, loading, error, loadProduct } = useProductFormViewModel();
  const { forecast, loading: forecastLoading, fetchForecast } = useDemandForecastViewModel();

  useEffect(() => {
    if (productSlug) {
      loadProduct(productSlug);
    }
  }, [productSlug, loadProduct]);

  useEffect(() => {
    if (product?.id) {
      fetchForecast(product.id);
    }
  }, [product?.id, fetchForecast]);

  if (loading) {
    return (
      <div className="border border-border/60 rounded-2xl p-8 bg-background flex flex-col items-center justify-center min-h-[350px]">
        <Loader2 className="h-8 w-8 animate-spin text-primary mb-3" />
        <p className="text-sm text-muted-foreground font-sans">Loading product details...</p>
      </div>
    );
  }

  if (error || !product) {
    return (
      <div className="border border-border/60 rounded-2xl p-6 bg-background space-y-4">
        <div className="flex items-center justify-between pb-4 border-b border-border/60">
          <h3 className="text-xl font-bold font-display text-destructive">Error Loading Product</h3>
          <Button variant="ghost" size="icon" className="h-8 w-8 rounded-xl" onClick={onClose}>
            <X className="h-4 w-4" />
          </Button>
        </div>
        <p className="text-sm text-destructive bg-destructive/10 p-3 rounded-xl border border-destructive/20">
          {error || 'Product not found.'}
        </p>
      </div>
    );
  }

  const formattedPrice = new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
  }).format(product.price || 0);

  return (
    <div className="border border-border/60 rounded-2xl p-6 bg-background space-y-6 shadow-none animate-in fade-in slide-in-from-right-2 duration-200">
      {/* Header Bar */}
      <div className="flex items-start justify-between pb-4 border-b border-border/60 gap-4">
        <div className="min-w-0 flex-1">
          <h3 className="text-2xl font-bold font-display tracking-tight text-foreground line-clamp-2 leading-tight">
            {product.name}
          </h3>
          <p className="text-xs text-muted-foreground font-mono mt-0.5">SKU: {product.sku || 'N/A'}</p>
        </div>

        <div className="flex items-center gap-2 shrink-0">
          <Button
            onClick={() => onEditProduct(product.slug)}
            className="flex items-center gap-1.5 bg-primary hover:bg-primary/90 text-primary-foreground rounded-xl text-xs"
            size="sm"
          >
            <Edit className="h-3.5 w-3.5" /> Edit Product
          </Button>
          <Button
            variant="ghost"
            size="icon"
            className="h-8 w-8 rounded-xl text-muted-foreground hover:text-foreground"
            onClick={onClose}
          >
            <X className="h-4 w-4" />
          </Button>
        </div>
      </div>

      {/* Main Details Grid */}
      <div className="grid grid-cols-1 md:grid-cols-12 gap-6">
        {/* Left Column: Image Banner */}
        <div className="md:col-span-5 space-y-3">
          <div className="aspect-square w-full rounded-2xl border border-border/60 overflow-hidden bg-muted flex items-center justify-center relative group">
            {product.banner?.preview || product.banner?.thumbnail ? (
              <img
                src={product.banner.preview || product.banner.thumbnail!}
                alt={product.name}
                className="h-full w-full object-cover"
              />
            ) : (
              <div className="flex flex-col items-center justify-center text-muted-foreground p-4 text-center">
                <Package className="h-12 w-12 stroke-[1.5] mb-2" />
                <span className="text-xs">No image banner</span>
              </div>
            )}
          </div>
        </div>

        {/* Right Column: Key Attributes */}
        <div className="md:col-span-7 space-y-4">
          {/* Status Section */}
          <div className="space-y-1">
            <span className="text-xs text-muted-foreground uppercase font-medium tracking-wider">Status & Availability</span>
            <div className="flex items-center gap-2 pt-0.5">
              <StatusBadge status={product.status} />
              {product.is_available ? (
                <Badge variant="outline" className="text-[10px] bg-primary/10 text-primary border-primary/20 gap-1 rounded-md">
                  <CheckCircle2 className="h-3 w-3" /> Listed & Active
                </Badge>
              ) : (
                <Badge variant="outline" className="text-[10px] bg-muted text-muted-foreground border-border gap-1 rounded-md">
                  <XCircle className="h-3 w-3" /> Hidden / Disabled
                </Badge>
              )}
            </div>
          </div>

          {/* Unit Price Section */}
          <div className="space-y-1 pt-2 border-t border-border/60">
            <span className="text-xs text-muted-foreground uppercase font-medium tracking-wider">Unit Price</span>
            <div className="text-2xl font-bold font-display text-primary">{formattedPrice}</div>
          </div>

          {/* Weight & Total Stock Section */}
          <div className="grid grid-cols-2 gap-4 pt-2 border-t border-border/60 text-xs">
            <div>
              <span className="text-xs text-muted-foreground uppercase font-medium tracking-wider flex items-center gap-1 mb-1">
                <Weight className="h-3.5 w-3.5 text-muted-foreground" /> Weight
              </span>
              <p className="font-semibold text-foreground font-mono">
                {product.weight ? `${product.weight} g` : 'Not specified'}
              </p>
            </div>
            <div>
              <span className="text-xs text-muted-foreground uppercase font-medium tracking-wider flex items-center gap-1 mb-1">
                <BarChart2 className="h-3.5 w-3.5 text-muted-foreground" /> Total Stock
              </span>
              <p className={`font-semibold font-mono ${product.stock < 10 ? 'text-destructive' : 'text-foreground'}`}>
                {product.stock ?? 0} units
              </p>
            </div>
          </div>

          {/* Description Section */}
          <div className="space-y-1.5 pt-2 border-t border-border/60">
            <span className="text-xs text-muted-foreground uppercase font-medium tracking-wider">Description</span>
            <p className="text-xs text-foreground/80 leading-relaxed font-sans whitespace-pre-wrap">
              {product.description || 'No product description provided.'}
            </p>
          </div>
        </div>
      </div>

      {/* AI Demand & Stockout Intelligence Card */}
      <div className="space-y-3 pt-4 border-t border-border/60">
        <div className="flex items-center justify-between">
          <span className="text-xs font-bold font-display tracking-tight text-foreground flex items-center gap-1.5">
            <Sparkles className="h-4 w-4 text-primary" />
            AI Demand & Stockout Intelligence
          </span>
          {forecast && (
            <Badge variant="outline" className="text-[10px] bg-primary/10 text-primary border-primary/20 uppercase font-mono">
              {forecast.confidence_tier} Confidence
            </Badge>
          )}
        </div>

        {forecastLoading ? (
          <div className="p-4 rounded-xl border border-border/60 bg-muted/20 space-y-2">
            <Skeleton className="h-4 w-48 bg-muted" />
            <Skeleton className="h-8 w-full bg-muted" />
          </div>
        ) : forecast ? (
          <div className="p-4 rounded-xl border border-border/60 bg-muted/10 space-y-3">
            <div className="grid grid-cols-2 gap-3 text-xs">
              <div className="p-2.5 rounded-lg bg-background border border-border/60">
                <span className="text-[10px] text-muted-foreground uppercase font-medium tracking-wider flex items-center gap-1">
                  <TrendingUp className="h-3 w-3 text-primary" /> 7d Forecast
                </span>
                <div className="text-base font-bold font-display text-primary mt-0.5">
                  {forecast.predicted_units_sold_7d.toFixed(1)} <span className="text-[10px] font-normal text-muted-foreground">units</span>
                </div>
              </div>

              <div className="p-2.5 rounded-lg bg-background border border-border/60">
                <span className="text-[10px] text-muted-foreground uppercase font-medium tracking-wider flex items-center gap-1">
                  <BarChart2 className="h-3 w-3 text-muted-foreground" /> Historical 7d
                </span>
                <div className="text-base font-bold font-display text-foreground mt-0.5">
                  {forecast.historical_velocity_7d} <span className="text-[10px] font-normal text-muted-foreground">units</span>
                </div>
              </div>
            </div>

            <div className="flex items-center justify-between text-xs pt-1 border-t border-border/40">
              <span className="text-muted-foreground">Depletion Runway:</span>
              <span className="font-semibold text-foreground">
                {forecast.current_stock < forecast.predicted_units_sold_7d ? (
                  <span className="text-destructive font-bold flex items-center gap-1">
                    <AlertTriangle className="h-3.5 w-3.5" /> High Risk (Stockout Deficit)
                  </span>
                ) : (
                  <span className="text-emerald-600 font-medium">
                    ✓ Stock sufficient for ~{Math.max(1, Math.round(forecast.current_stock / Math.max(0.1, forecast.predicted_units_sold_7d / 7)))} days
                  </span>
                )}
              </span>
            </div>
          </div>
        ) : null}
      </div>

      {/* Shop Stock Breakdown Section */}
      <div className="space-y-3 pt-4 border-t border-border/60">
        <div className="flex items-center justify-between">
          <span className="text-xs text-muted-foreground uppercase font-medium tracking-wider">
            Store Availability & Inventory
          </span>
          <span className="text-xs text-muted-foreground font-mono">
            {product.availability?.length || 0} locations registered
          </span>
        </div>

        {(!product.availability || product.availability.length === 0) ? (
          <div className="py-6 text-center text-xs text-muted-foreground border border-dashed border-border/80 rounded-xl bg-muted/20">
            This product is not assigned to any shop inventory.
          </div>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-2.5">
            {product.availability.map((avail, idx) => (
              <div
                key={idx}
                className="
                  flex flex-col gap-2
                  p-3 rounded-xl border border-border/60 bg-muted/20
                  lg:justify-between xl:flex-row xl:gap-0
                "
              >
                <div className="min-w-0 pr-2">
                  <p className="font-semibold text-xs text-foreground truncate">{avail.slug || avail.name}</p>
                  <p className="text-[10px] text-muted-foreground font-mono">Store Branch</p>
                </div>
                <Badge variant="secondary" className="font-mono text-xs shrink-0 w-fit rounded-lg bg-background text-primary border border-border/60">
                  {avail.stock} in stock
                </Badge>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
