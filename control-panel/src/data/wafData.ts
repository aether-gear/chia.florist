// @ts-ignore
import wafLogs from '../../waf-logs.json';
import { format, parseISO } from 'date-fns';

export interface WafLog {
  id: string;
  timestamp: string;
  ip: string;
  method: string;
  url: string;
  status: 'Allowed' | 'Blocked';
  rule_id?: string;
  details: string;
  geo: string;
}

const logs = wafLogs as WafLog[];

export const getRecentLogs = (count: number = 20): WafLog[] => {
  return [...logs].sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime()).slice(0, count);
};

export const getThreatData = () => {
  // Aggregate by hour or by a specific interval to show in a chart
  const dataMap: Record<string, { time: string; blocked: number; allowed: number }> = {};
  
  logs.forEach(log => {
    try {
      const date = parseISO(log.timestamp);
      // Group by minute for a more detailed threat graph given the logs timeframe
      const key = format(date, 'HH:mm'); 
      if (!dataMap[key]) {
        dataMap[key] = { time: key, blocked: 0, allowed: 0 };
      }
      if (log.status === 'Blocked') {
        dataMap[key].blocked += 1;
      } else {
        dataMap[key].allowed += 1;
      }
    } catch (e) {
      console.error("Invalid date", log.timestamp);
    }
  });

  return Object.values(dataMap).sort((a, b) => a.time.localeCompare(b.time));
};

export const getWafSummary = () => {
  const total = logs.length;
  const blocked = logs.filter(l => l.status === 'Blocked').length;
  const allowed = total - blocked;
  const threatLevel = blocked > 50 ? 'High' : (blocked > 10 ? 'Medium' : 'Low');

  return {
    total,
    blocked,
    allowed,
    threatLevel,
    activeRules: 5, // Mock value
  };
};
