import { useState, useMemo } from 'react';
import {
  ResponsiveContainer,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip,
  PieChart,
  Pie,
  Cell,
  ScatterChart,
  Scatter,
  ZAxis,
} from 'recharts';
import type { ProductStat } from '../../models/Product';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '../ui/card';

interface ProductPerformanceChartsProps {
  stats: ProductStat[];
}

type TimeWindow = '7d' | '30d' | '90d';

export default function ProductPerformanceCharts({ stats }: ProductPerformanceChartsProps) {
  const [timeWindow, setTimeWindow] = useState<TimeWindow>('30d');

  // Colors based on index.css CSS variables
  const colors = [
    'hsl(var(--chart-1))',
    'hsl(var(--chart-2))',
    'hsl(var(--chart-3))',
    'hsl(var(--chart-4))',
    'hsl(var(--chart-5))',
    '#9ca3af', // Gray color for 'Others'
  ];

  // 1. Horizontal Bar Chart data: Top products by sales velocity
  const barChartData = useMemo(() => {
    const field =
      timeWindow === '7d'
        ? 'sales_velocity_7d'
        : timeWindow === '30d'
        ? 'sales_velocity_30d'
        : 'sales_velocity_90d';

    return [...stats]
      .sort((a, b) => b[field] - a[field])
      .slice(0, 8)
      .map((item) => ({
        name: item.name,
        sales: item[field],
      }));
  }, [stats, timeWindow]);

  // 2. Donut Chart data: Revenue contribution percentage
  const donutChartData = useMemo(() => {
    const sorted = [...stats].sort(
      (a, b) => b.revenue_contribution_percentage - a.revenue_contribution_percentage
    );
    const top5 = sorted.slice(0, 5);
    const othersPct = sorted
      .slice(5)
      .reduce((acc, item) => acc + item.revenue_contribution_percentage, 0);

    const data = top5.map((item) => ({
      name: item.name,
      value: parseFloat(item.revenue_contribution_percentage.toFixed(1)),
    }));

    if (othersPct > 0) {
      data.push({
        name: 'Others',
        value: parseFloat(othersPct.toFixed(1)),
      });
    }
    return data;
  }, [stats]);

  // 3. Scatter Chart data: Conversion Rate vs Views
  const scatterChartData = useMemo(() => {
    return stats.map((item) => ({
      name: item.name,
      views: item.view_count,
      conversion: item.conversion_rate,
      stock: item.stock,
    }));
  }, [stats]);

  // Custom tooltips for nice, borderless theme matching
  const BarTooltip = ({ active, payload }: any) => {
    if (active && payload && payload.length) {
      const data = payload[0].payload;
      return (
        <div className="bg-popover text-popover-foreground border border-border rounded-xl p-3 shadow-md text-xs font-sans">
          <p className="font-semibold font-display mb-1">{data.name}</p>
          <p className="text-muted-foreground">
            Units Sold: <span className="font-semibold text-primary">{data.sales}</span>
          </p>
        </div>
      );
    }
    return null;
  };

  const DonutTooltip = ({ active, payload }: any) => {
    if (active && payload && payload.length) {
      const data = payload[0].payload;
      return (
        <div className="bg-popover text-popover-foreground border border-border rounded-xl p-3 shadow-md text-xs font-sans">
          <p className="font-semibold font-display mb-1">{data.name}</p>
          <p className="text-muted-foreground">
            Revenue Share: <span className="font-semibold text-primary">{data.value}%</span>
          </p>
        </div>
      );
    }
    return null;
  };

  const ScatterTooltip = ({ active, payload }: any) => {
    if (active && payload && payload.length) {
      const data = payload[0].payload;
      return (
        <div className="bg-popover text-popover-foreground border border-border rounded-xl p-3 shadow-md text-xs font-sans">
          <p className="font-semibold font-display mb-1">{data.name}</p>
          <p className="text-muted-foreground">
            Views: <span className="font-medium text-foreground">{data.views}</span>
          </p>
          <p className="text-muted-foreground">
            Conversion Rate: <span className="font-semibold text-primary">{data.conversion}%</span>
          </p>
          <p className="text-muted-foreground">
            Stock: <span className="font-medium text-foreground">{data.stock} units</span>
          </p>
        </div>
      );
    }
    return null;
  };

  if (stats.length === 0) {
    return null;
  }

  return (
    <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 animate-in fade-in duration-300">
      {/* 1. Bar Chart Card */}
      <Card className="border-0 shadow-none bg-zinc-50/40 dark:bg-slate-900/40 flex flex-col justify-between">
        <CardHeader className="pb-2">
          <div className="flex items-center justify-between gap-2 flex-wrap">
            <div>
              <CardTitle className="font-bold font-display tracking-tight text-lg text-foreground">
                Sales Velocity
              </CardTitle>
              <CardDescription className="text-muted-foreground text-xs font-sans">
                Top products sold by volume
              </CardDescription>
            </div>
            {/* Pill Toggles */}
            <div className="flex items-center gap-1 bg-muted p-0.5 rounded-lg text-xs">
              {(['7d', '30d', '90d'] as TimeWindow[]).map((window) => (
                <button
                  key={window}
                  onClick={() => setTimeWindow(window)}
                  className={`px-2 py-1 rounded-md transition-colors font-medium font-sans ${
                    timeWindow === window
                      ? 'bg-primary text-primary-foreground shadow-sm'
                      : 'text-muted-foreground hover:text-foreground'
                  }`}
                >
                  {window}
                </button>
              ))}
            </div>
          </div>
        </CardHeader>
        <CardContent className="h-[240px] pt-4">
          <ResponsiveContainer width="100%" height="100%">
            <BarChart
              data={barChartData}
              layout="vertical"
              margin={{ top: 5, right: 10, left: -10, bottom: 5 }}
            >
              <XAxis type="number" hide />
              <YAxis
                dataKey="name"
                type="category"
                width={85}
                axisLine={false}
                tickLine={false}
                tick={{ fontSize: 10, fill: 'var(--muted-foreground)' }}
              />
              <Tooltip content={<BarTooltip />} cursor={{ fill: 'rgba(21, 94, 55, 0.05)' }} />
              <Bar dataKey="sales" radius={[0, 4, 4, 0]}>
                {barChartData.map((_, index) => (
                  <Cell key={`cell-${index}`} fill={colors[index % 5]} />
                ))}
              </Bar>
            </BarChart>
          </ResponsiveContainer>
        </CardContent>
      </Card>

      {/* 2. Donut Chart Card */}
      <Card className="border-0 shadow-none bg-zinc-50/40 dark:bg-slate-900/40 flex flex-col justify-between">
        <CardHeader className="pb-2">
          <CardTitle className="font-bold font-display tracking-tight text-lg text-foreground">
            Revenue Share
          </CardTitle>
          <CardDescription className="text-muted-foreground text-xs font-sans">
            Revenue contribution percentage breakdown
          </CardDescription>
        </CardHeader>
        <CardContent className="h-[240px] pt-2">
          <div className="flex flex-col items-center justify-center h-full w-full">
            <div className="w-full h-[140px] relative flex items-center justify-center">
              <ResponsiveContainer width="100%" height="100%">
                <PieChart>
                  <Pie
                    data={donutChartData}
                    cx="50%"
                    cy="50%"
                    innerRadius={42}
                    outerRadius={58}
                    paddingAngle={3}
                    dataKey="value"
                  >
                    {donutChartData.map((_, index) => (
                      <Cell key={`cell-${index}`} fill={colors[index % colors.length]} />
                    ))}
                  </Pie>
                  <Tooltip content={<DonutTooltip />} />
                </PieChart>
              </ResponsiveContainer>
              {/* Center Legend Overlay */}
              <div className="absolute text-center pointer-events-none">
                <span className="text-lg font-bold font-display text-primary block leading-none">
                  Top 5
                </span>
                <span className="text-[9px] text-muted-foreground uppercase tracking-wider font-sans font-medium">
                  Products
                </span>
              </div>
            </div>
            {/* Direct Legend Panel below the Donut */}
            <div className="w-full flex flex-wrap justify-center items-center gap-x-3 gap-y-1.5 font-sans text-[11px] mt-2 px-1">
              {donutChartData.map((item, index) => (
                <div key={item.name} className="flex items-center gap-1 text-foreground">
                  <span 
                    className="w-2.5 h-2.5 rounded-full shrink-0" 
                    style={{ backgroundColor: colors[index % colors.length] }}
                  />
                  <span className="font-medium max-w-[100px] truncate" title={item.name}>
                    {item.name}
                  </span>
                  <span className="font-semibold text-muted-foreground text-[10px]">
                    ({item.value}%)
                  </span>
                </div>
              ))}
            </div>
          </div>
        </CardContent>
      </Card>

      {/* 3. Scatter Chart Card */}
      <Card className="border-0 shadow-none bg-zinc-50/40 dark:bg-slate-900/40 flex flex-col justify-between">
        <CardHeader className="pb-2">
          <CardTitle className="font-bold font-display tracking-tight text-lg text-foreground">
            Conversion vs Views
          </CardTitle>
          <CardDescription className="text-muted-foreground text-xs font-sans">
            Identify high-interest vs high-conversion items
          </CardDescription>
        </CardHeader>
        <CardContent className="h-[240px] pt-4">
          <ResponsiveContainer width="100%" height="100%">
            <ScatterChart margin={{ top: 10, right: 10, left: -25, bottom: 5 }}>
              <XAxis
                type="number"
                dataKey="views"
                name="views"
                unit=" views"
                axisLine={false}
                tickLine={false}
                tick={{ fontSize: 9, fill: 'var(--muted-foreground)' }}
              />
              <YAxis
                type="number"
                dataKey="conversion"
                name="conversion"
                unit="%"
                axisLine={false}
                tickLine={false}
                tick={{ fontSize: 9, fill: 'var(--muted-foreground)' }}
              />
              <ZAxis type="number" dataKey="stock" range={[40, 300]} />
              <Tooltip cursor={{ strokeDasharray: '3 3' }} content={<ScatterTooltip />} />
              <Scatter name="Products" data={scatterChartData} fill="hsl(var(--primary))" opacity={0.7} />
            </ScatterChart>
          </ResponsiveContainer>
        </CardContent>
      </Card>
    </div>
  );
}
