import { useState, useEffect } from 'react';
import { Loader2, Save, Plus, MapPin, AlertCircle } from 'lucide-react';
import { Button } from '../ui/button';
import { Input } from '../ui/input';
import { Label } from '../ui/label';
import { Checkbox } from '../ui/checkbox';
import {
  Sheet,
  SheetContent,
} from '../ui/sheet';
import { fetchApi } from '../../lib/api';

interface AddressFormSheetProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  shopId: string;
  shopName?: string;
  address?: any | null; // If passed -> Edit mode; if null/omitted -> Add mode
  onSuccess: () => void;
}

// Helper to extract array from API response
const extractList = (res: any, key: string): any[] => {
  if (!res) return [];
  if (Array.isArray(res)) return res;
  if (Array.isArray(res[key])) return res[key];
  if (Array.isArray(res.data)) return res.data;
  if (res.data && Array.isArray(res.data[key])) return res.data[key];
  for (const k of Object.keys(res)) {
    if (Array.isArray(res[k])) return res[k];
  }
  return [];
};

// Dedicated helpers to safely extract location IDs per level
const getProvinceId = (item: any): string => {
  if (!item) return '';
  const val = item.province_id ?? item.provinceId ?? item.id ?? item.code ?? '';
  return String(val);
};

const getCityId = (item: any): string => {
  if (!item) return '';
  const val = item.city_id ?? item.cityId ?? item.id ?? item.code ?? '';
  return String(val);
};

const getDistrictId = (item: any): string => {
  if (!item) return '';
  const val = item.district_id ?? item.districtId ?? item.id ?? item.code ?? '';
  return String(val);
};

const getVillageId = (item: any): string => {
  if (!item) return '';
  const val = item.village_id ?? item.villageId ?? item.id ?? item.code ?? '';
  return String(val);
};

export default function AddressFormSheet({
  open,
  onOpenChange,
  shopId,
  shopName,
  address,
  onSuccess,
}: AddressFormSheetProps) {
  const isEdit = Boolean(address);

  // Form states
  const [label, setLabel] = useState('');
  const [phone, setPhone] = useState('');
  const [fullAddress, setFullAddress] = useState('');
  const [postalCode, setPostalCode] = useState('');
  const [isActive, setIsActive] = useState(false);

  // Location dropdown lists
  const [provinces, setProvinces] = useState<any[]>([]);
  const [cities, setCities] = useState<any[]>([]);
  const [districts, setDistricts] = useState<any[]>([]);
  const [villages, setVillages] = useState<any[]>([]);

  // Selected location IDs
  const [provinceId, setProvinceId] = useState<string>('');
  const [cityId, setCityId] = useState<string>('');
  const [districtId, setDistrictId] = useState<string>('');
  const [villageId, setVillageId] = useState<string>('');

  const [loadingLocations, setLoadingLocations] = useState<boolean>(false);
  const [submitting, setSubmitting] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  // Sync form values & fetch location data sequentially FIRST before attaching values
  useEffect(() => {
    if (!open) return;

    setError(null);

    const initFormAndLocations = async () => {
      setLoadingLocations(true);

      if (isEdit && address) {
        // Set basic form text fields
        setLabel(address.label || '');
        setPhone(address.phone || '');
        setFullAddress(address.full_address || '');
        setPostalCode(address.postal_code || '');
        setIsActive(Boolean(address.is_active));

        try {
          // 1. Fetch Provinces
          const provRes = await fetchApi('/provinces');
          const provList = extractList(provRes, 'provinces');
          const targetProvId = String(address.province_id ?? '');
          const foundProv = provList.find((p: any) => String(p.id) === targetProvId);
          const matchedProvId = foundProv ? String(foundProv.id) : '';

          // 2. Fetch Cities
          let cityList: any[] = [];
          let matchedCityId = '';
          if (matchedProvId) {
            const cityRes = await fetchApi(`/provinces/${matchedProvId}/cities`);
            cityList = extractList(cityRes, 'cities');
            const targetCityId = String(address.city_id ?? '');
            const foundCity = cityList.find((c: any) => String(c.id) === targetCityId);
            matchedCityId = foundCity ? String(foundCity.id) : '';
          }

          // 3. Fetch Districts
          let distList: any[] = [];
          let matchedDistId = '';
          if (matchedCityId) {
            const distRes = await fetchApi(`/cities/${matchedCityId}/districts`);
            distList = extractList(distRes, 'districts');
            const targetDistId = String(address.district_id ?? '');
            const foundDist = distList.find((d: any) => String(d.id) === targetDistId);
            matchedDistId = foundDist ? String(foundDist.id) : '';
          }

          // 4. Fetch Villages
          let vilList: any[] = [];
          let matchedVilId = '';
          if (matchedDistId) {
            const vilRes = await fetchApi(`/districts/${matchedDistId}/villages`);
            vilList = extractList(vilRes, 'villages');
            const targetVilId = String(address.village_id ?? '');
            const foundVil = vilList.find((v: any) => String(v.id) === targetVilId);
            matchedVilId = foundVil ? String(foundVil.id) : '';
          }


          // Attach option lists first, then selected values
          setProvinces(provList);
          setCities(cityList);
          setDistricts(distList);
          setVillages(vilList);

          setProvinceId(matchedProvId);
          setCityId(matchedCityId);
          setDistrictId(matchedDistId);
          setVillageId(matchedVilId);
        } catch (err: any) {
          console.error('Failed to load edit location chain', err);
          setError(err.message || 'Failed to load location parameters');
        } finally {
          setLoadingLocations(false);
        }
      } else {
        // Reset form for Add mode
        setLabel('');
        setPhone('');
        setFullAddress('');
        setPostalCode('');
        setIsActive(false);
        setProvinceId('');
        setCityId('');
        setDistrictId('');
        setVillageId('');
        setCities([]);
        setDistricts([]);
        setVillages([]);

        try {
          const provRes = await fetchApi('/provinces');
          setProvinces(extractList(provRes, 'provinces'));
        } catch (err: any) {
          console.error('Failed to load provinces', err);
          setError(err.message || 'Failed to load provinces');
        } finally {
          setLoadingLocations(false);
        }
      }
    };

    initFormAndLocations();
  }, [open, isEdit, address]);

  // Location change handlers for user manual selection
  const handleProvinceChange = async (provId: string) => {
    setProvinceId(provId);
    setCityId('');
    setDistrictId('');
    setVillageId('');
    setCities([]);
    setDistricts([]);
    setVillages([]);
    if (!provId) return;
    try {
      setError(null);
      const res = await fetchApi(`/provinces/${provId}/cities`);
      setCities(extractList(res, 'cities'));
    } catch (err: any) {
      console.error(err);
      setError(err.message || 'Failed to load cities');
    }
  };

  const handleCityChange = async (cId: string) => {
    setCityId(cId);
    setDistrictId('');
    setVillageId('');
    setDistricts([]);
    setVillages([]);
    if (!cId) return;
    try {
      setError(null);
      const res = await fetchApi(`/cities/${cId}/districts`);
      setDistricts(extractList(res, 'districts'));
    } catch (err: any) {
      console.error(err);
      setError(err.message || 'Failed to load districts');
    }
  };

  const handleDistrictChange = async (distId: string) => {
    setDistrictId(distId);
    setVillageId('');
    setVillages([]);
    if (!distId) return;
    try {
      setError(null);
      const res = await fetchApi(`/districts/${distId}/villages`);
      setVillages(extractList(res, 'villages'));
    } catch (err: any) {
      console.error(err);
      setError(err.message || 'Failed to load villages');
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!shopId || !label || !provinceId || !cityId || !districtId || !villageId || !fullAddress || !postalCode) {
      setError('Please fill in all required location fields.');
      return;
    }

    setSubmitting(true);
    setError(null);

    const formData = {
      label,
      phone: phone || null,
      province_id: provinceId,
      city_id: cityId,
      district_id: districtId,
      village_id: villageId,
      full_address: fullAddress,
      postal_code: postalCode,
      is_active: isActive,
    };

    try {
      const isActiveStr = isActive ? 'true' : 'false';
      if (isEdit && address) {
        await fetchApi(`/shops/${shopId}/addresses/${address.id}`, {
          method: 'PUT',
          body: JSON.stringify({
            ...formData,
            shop_id: shopId,
            is_active: isActiveStr,
          }),
        });
      } else {
        await fetchApi(`/shops/${shopId}/addresses`, {
          method: 'POST',
          body: JSON.stringify({
            ...formData,
            is_active: isActiveStr,
          }),
        });
      }
      onSuccess();
      onOpenChange(false);
    } catch (err: any) {
      console.error(err);
      setError(err.message || (isEdit ? 'Failed to update address' : 'Failed to create address'));
    } finally {
      setSubmitting(false);
    }
  };

  const formContent = (
    <>
      {/* Header */}
      <div className="px-6 py-5 border-b flex items-center justify-between shrink-0">
        <div>
          <h3 className="text-xl font-bold font-display text-foreground flex items-center gap-2">
            <MapPin className="h-5 w-5 text-primary" />
            {isEdit ? 'Update Shop Address' : 'Add Shop Address'}
          </h3>
          <p className="text-xs text-muted-foreground mt-1">
            {isEdit ? (
              <>Modify location parameters for <strong className="text-foreground">{address?.label}</strong>{shopName ? ` at ${shopName}` : ''}.</>
            ) : (
              <>Create a physical pickup or delivery location{shopName ? ` for ${shopName}` : ''}.</>
            )}
          </p>
        </div>
      </div>

      {/* Scrollable Body */}
      <div className="flex-1 overflow-y-auto px-6 py-6 space-y-6 pb-24">
        {error && (
          <div className="flex items-start gap-2 p-3 text-xs text-destructive bg-destructive/10 rounded-xl border border-destructive/20 font-sans">
            <AlertCircle className="h-4 w-4 shrink-0 text-destructive mt-0.5" />
            <span>{error}</span>
          </div>
        )}

        {loadingLocations && (
          <div className="flex items-center gap-2 p-3 text-xs text-primary bg-primary/10 rounded-xl border border-primary/20 font-sans animate-pulse">
            <Loader2 className="h-3.5 w-3.5 animate-spin text-primary shrink-0" />
            <span>Loading location parameters and pre-populating address chain...</span>
          </div>
        )}

        {/* Basic Info Section */}
        <div className="space-y-4">
          <h3 className="text-xs font-bold uppercase tracking-wider text-muted-foreground">Basic Information</h3>
          <div className="grid grid-cols-2 gap-3">
            <div>
              <Label htmlFor="addressLabel" className="text-xs font-semibold text-foreground">
                Branch Label / Name *
              </Label>
              <Input
                id="addressLabel"
                placeholder="e.g. Main Warehouse"
                value={label}
                onChange={(e) => setLabel(e.target.value)}
                className="mt-1 text-xs rounded-xl"
                required
              />
            </div>
            <div>
              <Label htmlFor="addressPhone" className="text-xs font-semibold text-foreground">
                Contact Phone
              </Label>
              <Input
                id="addressPhone"
                placeholder="e.g. +62812345678"
                value={phone}
                onChange={(e) => setPhone(e.target.value)}
                className="mt-1 text-xs rounded-xl"
              />
            </div>
          </div>
        </div>

        {/* Location Details Section */}
        <div className="space-y-4 pt-2">
          <h3 className="text-xs font-bold uppercase tracking-wider text-muted-foreground">Location Details</h3>

          <div className="grid grid-cols-2 gap-3">
            <div>
              <Label htmlFor="province" className="text-xs font-semibold text-foreground">Province *</Label>
              <select
                id="province"
                value={provinceId}
                onChange={(e) => handleProvinceChange(e.target.value)}
                disabled={loadingLocations || provinces.length === 0}
                className="mt-1 w-full rounded-xl border border-border bg-background p-2.5 text-xs text-foreground focus:outline-none focus:ring-1 focus:ring-primary/40 disabled:opacity-50"
                required
              >
                <option value="">Select Province</option>
                {provinces.map((p) => {
                  const idVal = getProvinceId(p);
                  return (
                    <option key={idVal} value={idVal}>
                      {p.name}
                    </option>
                  );
                })}
              </select>
            </div>

            <div>
              <Label htmlFor="city" className="text-xs font-semibold text-foreground">City / Regency *</Label>
              <select
                id="city"
                value={cityId}
                onChange={(e) => handleCityChange(e.target.value)}
                disabled={loadingLocations || !provinceId || cities.length === 0}
                className="mt-1 w-full rounded-xl border border-border bg-background p-2.5 text-xs text-foreground focus:outline-none focus:ring-1 focus:ring-primary/40 disabled:opacity-50"
                required
              >
                <option value="">Select City</option>
                {cities.map((c) => {
                  const idVal = getCityId(c);
                  return (
                    <option key={idVal} value={idVal}>
                      {c.name}
                    </option>
                  );
                })}
              </select>
            </div>
          </div>

          <div className="grid grid-cols-2 gap-3">
            <div>
              <Label htmlFor="district" className="text-xs font-semibold text-foreground">District *</Label>
              <select
                id="district"
                value={districtId}
                onChange={(e) => handleDistrictChange(e.target.value)}
                disabled={loadingLocations || !cityId || districts.length === 0}
                className="mt-1 w-full rounded-xl border border-border bg-background p-2.5 text-xs text-foreground focus:outline-none focus:ring-1 focus:ring-primary/40 disabled:opacity-50"
                required
              >
                <option value="">Select District</option>
                {districts.map((d) => {
                  const idVal = getDistrictId(d);
                  return (
                    <option key={idVal} value={idVal}>
                      {d.name}
                    </option>
                  );
                })}
              </select>
            </div>

            <div>
              <Label htmlFor="village" className="text-xs font-semibold text-foreground">Village *</Label>
              <select
                id="village"
                value={villageId}
                onChange={(e) => setVillageId(e.target.value)}
                disabled={loadingLocations || !districtId || villages.length === 0}
                className="mt-1 w-full rounded-xl border border-border bg-background p-2.5 text-xs text-foreground focus:outline-none focus:ring-1 focus:ring-primary/40 disabled:opacity-50"
                required
              >
                <option value="">Select Village</option>
                {villages.map((v) => {
                  const idVal = getVillageId(v);
                  return (
                    <option key={idVal} value={idVal}>
                      {v.name}
                    </option>
                  );
                })}
              </select>
            </div>
          </div>

          <div>
            <Label htmlFor="fullAddress" className="text-xs font-semibold text-foreground">Full Street Address *</Label>
            <textarea
              id="fullAddress"
              rows={3}
              placeholder="Street name, building number, unit, landmark, etc."
              value={fullAddress}
              onChange={(e) => setFullAddress(e.target.value)}
              className="mt-1 w-full rounded-xl border border-border bg-background p-2.5 text-xs text-foreground focus:outline-none focus:ring-1 focus:ring-primary/40"
              required
            />
          </div>

          <div className="grid grid-cols-2 gap-3">
            <div>
              <Label htmlFor="postalCode" className="text-xs font-semibold text-foreground">Postal Code *</Label>
              <Input
                id="postalCode"
                placeholder="e.g. 17520"
                value={postalCode}
                onChange={(e) => setPostalCode(e.target.value)}
                className="mt-1 text-xs rounded-xl"
                required
              />
            </div>

            <div className="flex items-center pt-5">
              <label className="flex items-center gap-2 cursor-pointer text-xs font-medium text-foreground">
                <Checkbox
                  id="isActive"
                  checked={isActive}
                  onCheckedChange={(checked) => setIsActive(checked === true)}
                />
                Set as active address
              </label>
            </div>
          </div>
        </div>
      </div>

      {/* Action Footer */}
      <div className="px-6 py-4 border-t bg-muted/10 flex items-center justify-end gap-2 shrink-0">
        <Button
          type="button"
          variant="outline"
          size="sm"
          className="text-xs rounded-xl"
          onClick={() => onOpenChange(false)}
          disabled={submitting || loadingLocations}
        >
          Cancel
        </Button>
        <Button
          type="submit"
          size="sm"
          className="text-xs rounded-xl font-medium flex items-center gap-1.5 bg-primary hover:bg-primary/90 text-primary-foreground"
          onClick={handleSubmit}
          disabled={submitting || loadingLocations}
        >
          {submitting ? (
            <Loader2 className="h-3.5 w-3.5 animate-spin" />
          ) : isEdit ? (
            <Save className="h-3.5 w-3.5" />
          ) : (
            <Plus className="h-3.5 w-3.5" />
          )}
          {isEdit ? 'Save Address' : 'Add Address'}
        </Button>
      </div>
    </>
  );

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent className="w-full sm:max-w-none md:w-[45vw] md:min-w-[45vw] p-0 flex flex-col h-full border-l border-border/60 bg-background shadow-2xl">
        <form onSubmit={handleSubmit} className="flex flex-col h-full">
          {formContent}
        </form>
      </SheetContent>
    </Sheet>
  );
}
