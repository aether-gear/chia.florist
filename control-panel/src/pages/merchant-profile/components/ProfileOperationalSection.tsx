import { useFormContext, Controller } from 'react-hook-form';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';

const DAYS_OF_WEEK = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'];

export function ProfileOperationalSection() {
  const { register, control } = useFormContext();

  return (
    <Card className="border-0 shadow-sm">
      <CardHeader>
        <CardTitle>Operational Settings</CardTitle>
        <CardDescription>Configure your business hours and delivery capabilities.</CardDescription>
      </CardHeader>
      <CardContent className="space-y-6">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="space-y-2">
            <label className="text-sm font-medium">Opening Hours</label>
            <Input type="time" {...register('operational.openingHours')} />
          </div>
          <div className="space-y-2">
            <label className="text-sm font-medium">Closing Hours</label>
            <Input type="time" {...register('operational.closingHours')} />
          </div>
        </div>

        <div className="space-y-3">
          <label className="text-sm font-medium">Business Days</label>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
            <Controller
              name="operational.businessDays"
              control={control}
              render={({ field }) => (
                <>
                  {DAYS_OF_WEEK.map((day) => {
                    const isChecked = field.value?.includes(day);
                    return (
                      <div key={day} className="flex items-center space-x-2">
                        <input
                          type="checkbox"
                          id={`day-${day}`}
                          className="h-4 w-4 rounded border-slate-300 text-indigo-600 focus:ring-indigo-600"
                          checked={isChecked}
                          onChange={(e) => {
                            const newValue = e.target.checked
                              ? [...(field.value || []), day]
                              : (field.value || []).filter((d: string) => d !== day);
                            field.onChange(newValue);
                          }}
                        />
                        <label htmlFor={`day-${day}`} className="text-sm cursor-pointer select-none">
                          {day}
                        </label>
                      </div>
                    );
                  })}
                </>
              )}
            />
          </div>
        </div>

        <div className="pt-4 border-t border-slate-100">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <label className="text-sm font-medium">Delivery Radius (km)</label>
              <Input type="number" {...register('operational.deliveryRadius', { valueAsNumber: true })} />
            </div>
            <div className="space-y-2 flex flex-col justify-center pt-6">
              <div className="flex items-center space-x-2">
                <input 
                  type="checkbox" 
                  id="pickup" 
                  {...register('operational.pickupAvailable')}
                  className="h-4 w-4 rounded border-slate-300 text-indigo-600 focus:ring-indigo-600"
                />
                <label htmlFor="pickup" className="text-sm font-medium cursor-pointer select-none">
                  Store Pickup Available
                </label>
              </div>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
