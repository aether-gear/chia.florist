import React from 'react';
import { AlertTriangle, ShieldCheck, Sparkles, TrendingUp } from 'lucide-react';
import type { DynamicAIInsight } from '../../../viewmodels/useDashboardViewModel';
import { Skeleton } from '../../../components/ui/skeleton';

interface DashboardAIInsightsProps {
  insights: DynamicAIInsight[];
  loading: boolean;
}

export const DashboardAIInsights: React.FC<DashboardAIInsightsProps> = ({ insights, loading }) => {
  const getIcon = (category: string, severity: string) => {
    if (category === 'cyber') return <ShieldCheck className="w-4 h-4 text-rose-500 shrink-0 mt-0.5" />;
    if (severity === 'critical') return <AlertTriangle className="w-4 h-4 text-rose-500 shrink-0 mt-0.5" />;
    if (severity === 'warning') return <AlertTriangle className="w-4 h-4 text-amber-500 shrink-0 mt-0.5" />;
    if (category === 'ecommerce') return <TrendingUp className="w-4 h-4 text-primary shrink-0 mt-0.5" />;
    return <Sparkles className="w-4 h-4 text-primary shrink-0 mt-0.5" />;
  };

  const getCardStyle = (severity: string) => {
    switch (severity) {
      case 'critical':
        return 'bg-rose-500/5 border-rose-500/15 text-rose-950 dark:text-rose-200';
      case 'warning':
        return 'bg-amber-500/5 border-amber-500/15 text-amber-950 dark:text-amber-200';
      case 'success':
        return 'bg-primary/5 border-primary/15 text-primary-950 dark:text-primary-200';
      default:
        return 'bg-muted/40 border-border/40 text-foreground';
    }
  };

  return (
    <div className="flex flex-col space-y-4">
      <div className="pb-2 border-b border-border/60 flex items-center justify-between">
        <div>
          <h3 className="flex items-center font-bold font-display tracking-tight text-lg text-foreground">
            <Sparkles className="w-4 h-4 mr-2 text-amber-500" />
            Live AI Advisory Feed
          </h3>
          <p className="text-muted-foreground text-xs font-sans">
            Real-time automated signals across e-commerce, predictive stocking, and cyber telemetry.
          </p>
        </div>
      </div>

      <div className="space-y-3">
        {loading ? (
          Array.from({ length: 3 }).map((_, i) => (
            <div key={i} className="p-4 rounded-2xl bg-muted/30 border border-border/40 space-y-2">
              <Skeleton className="h-4 w-32 bg-muted" />
              <Skeleton className="h-3 w-full bg-muted" />
              <Skeleton className="h-3 w-4/5 bg-muted" />
            </div>
          ))
        ) : insights.length === 0 ? (
          <div className="p-6 rounded-2xl bg-muted/20 border border-border/40 text-center">
            <Sparkles className="w-6 h-6 text-primary mx-auto mb-2 opacity-60" />
            <p className="text-xs font-semibold text-foreground">All systems operating nominally</p>
            <p className="text-[11px] text-muted-foreground mt-0.5">
              No anomalies, critical stockouts, or urgent cyber threats detected at this time.
            </p>
          </div>
        ) : (
          insights.map((insight) => (
            <div
              key={insight.id}
              className={`p-4 rounded-2xl border transition-all duration-200 ${getCardStyle(insight.severity)}`}
            >
              <div className="flex items-start gap-2.5">
                {getIcon(insight.category, insight.severity)}
                <div className="min-w-0 flex-1">
                  <h4 className="text-xs font-bold font-display tracking-tight mb-1 text-foreground">
                    {insight.title}
                  </h4>
                  <p className="text-xs text-muted-foreground font-sans leading-relaxed">
                    {insight.description}
                  </p>
                </div>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
};

export default DashboardAIInsights;
