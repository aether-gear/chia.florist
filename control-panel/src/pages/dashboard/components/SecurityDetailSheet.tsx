import React from 'react';
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetDescription } from '../../../components/ui/sheet';
import { Badge } from '../../../components/ui/badge';
import { Button } from '../../../components/ui/button';
import { ShieldAlert, ShieldCheck, Globe, Terminal, Copy, Check } from 'lucide-react';
import type { SecurityEventLog } from '../../../viewmodels/useDashboardViewModel';

interface SecurityDetailSheetProps {
  log: SecurityEventLog | null;
  isOpen: boolean;
  onOpenChange: (open: boolean) => void;
}

export const SecurityDetailSheet: React.FC<SecurityDetailSheetProps> = ({
  log,
  isOpen,
  onOpenChange,
}) => {
  const [copied, setCopied] = React.useState(false);

  if (!log) return null;

  const handleCopyIp = () => {
    navigator.clipboard.writeText(log.ip);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const isBlocked = log.status === 'Blocked';

  return (
    <Sheet open={isOpen} onOpenChange={onOpenChange}>
      <SheetContent side="right" className="w-full sm:max-w-none md:w-[48vw] md:min-w-[480px] overflow-y-auto p-6 space-y-6">
        <SheetHeader className="pb-4 border-b border-border/60">
          <div className="flex items-center gap-2">
            {isBlocked ? (
              <div className="w-8 h-8 rounded-lg bg-rose-500/10 text-rose-600 dark:text-rose-400 flex items-center justify-center">
                <ShieldAlert className="w-4 h-4" />
              </div>
            ) : (
              <div className="w-8 h-8 rounded-lg bg-primary/10 text-primary flex items-center justify-center">
                <ShieldCheck className="w-4 h-4" />
              </div>
            )}
            <div>
              <SheetTitle className="font-display text-lg font-bold text-foreground">
                Security Incident Telemetry
              </SheetTitle>
              <SheetDescription className="text-xs text-muted-foreground font-sans">
                Full forensic log breakdown evaluated by the Web Application Firewall.
              </SheetDescription>
            </div>
          </div>
        </SheetHeader>

        {/* Status Banner */}
        <div className={`p-4 rounded-xl border flex items-center justify-between ${
          isBlocked
            ? 'bg-rose-500/5 border-rose-500/20 text-rose-700 dark:text-rose-300'
            : 'bg-primary/5 border-primary/20 text-primary'
        }`}>
          <div>
            <p className="text-xs font-semibold uppercase tracking-wider">Evaluation Outcome</p>
            <p className="text-sm font-bold mt-0.5">{isBlocked ? 'Request Blocked & Dropped' : 'Request Allowed'}</p>
          </div>
          <Badge variant={isBlocked ? 'destructive' : 'secondary'} className="rounded-md uppercase">
            {log.status}
          </Badge>
        </div>

        {/* Forensic Metadata Grid */}
        <div className="space-y-4 text-xs font-sans">
          <div className="grid grid-cols-2 gap-4 p-4 rounded-xl bg-muted/40 border border-border/40">
            <div>
              <span className="text-muted-foreground font-medium">Timestamp</span>
              <p className="font-mono font-semibold text-foreground mt-0.5">
                {new Date(log.timestamp).toLocaleString()}
              </p>
            </div>
            <div>
              <span className="text-muted-foreground font-medium">Rule Triggered</span>
              <p className="font-mono font-semibold text-foreground mt-0.5">
                {log.ruleId || 'N/A'}
              </p>
            </div>
            <div>
              <span className="text-muted-foreground font-medium">HTTP Method</span>
              <p className="font-mono font-bold text-foreground mt-0.5">{log.method}</p>
            </div>
            <div>
              <span className="text-muted-foreground font-medium">Target URL / Route</span>
              <p className="font-mono text-foreground mt-0.5 truncate" title={log.url}>
                {log.url}
              </p>
            </div>
          </div>

          {/* Client IP Section */}
          <div className="p-4 rounded-xl bg-muted/40 border border-border/40 space-y-2">
            <div className="flex items-center justify-between">
              <span className="text-muted-foreground font-medium flex items-center gap-1.5">
                <Globe className="w-3.5 h-3.5" />
                Origin Client IP
              </span>
              <Button
                variant="ghost"
                size="sm"
                onClick={handleCopyIp}
                className="h-7 text-xs text-muted-foreground hover:text-foreground gap-1"
              >
                {copied ? <Check className="w-3 h-3 text-primary" /> : <Copy className="w-3 h-3" />}
                {copied ? 'Copied' : 'Copy IP'}
              </Button>
            </div>
            <p className="font-mono font-bold text-sm text-foreground">{log.ip}</p>
          </div>

          {/* Trigger Reason */}
          {log.reason && log.reason !== '-' && (
            <div className="p-4 rounded-xl bg-muted/40 border border-border/40 space-y-1.5">
              <span className="text-muted-foreground font-medium flex items-center gap-1.5">
                <Terminal className="w-3.5 h-3.5" />
                Detection Signature / Reason
              </span>
              <p className="font-mono text-xs text-foreground bg-background p-2.5 rounded-lg border border-border/40 break-words">
                {log.reason}
              </p>
            </div>
          )}

          {/* Payload if present */}
          {log.payload && log.payload !== '-' && (
            <div className="p-4 rounded-xl bg-muted/40 border border-border/40 space-y-1.5">
              <span className="text-muted-foreground font-medium">Intercepted Payload Sample</span>
              <pre className="font-mono text-[11px] text-foreground bg-background p-2.5 rounded-lg border border-border/40 overflow-x-auto whitespace-pre-wrap break-all">
                {log.payload}
              </pre>
            </div>
          )}

          {/* User Agent */}
          {log.userAgent && log.userAgent !== '-' && (
            <div className="p-4 rounded-xl bg-muted/40 border border-border/40 space-y-1.5">
              <span className="text-muted-foreground font-medium">User-Agent Header</span>
              <p className="font-mono text-[11px] text-muted-foreground bg-background p-2.5 rounded-lg border border-border/40 break-words">
                {log.userAgent}
              </p>
            </div>
          )}
        </div>
      </SheetContent>
    </Sheet>
  );
};

export default SecurityDetailSheet;
