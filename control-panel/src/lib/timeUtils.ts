/**
 * Formats a last login timestamp into a human-friendly relative time text string.
 *
 * Rules:
 * - Under 1 minute: "Active now"
 * - 1 minute to 59 minutes: "X minute(s) ago"
 * - 1 hour to 12 hours: "X hour(s) ago"
 * - More than 12 hours up to 7 days: "X day(s) ago"
 * - Over 7 days: Standard date string (e.g. "15 Aug 2026")
 * - Null / undefined: "Never logged in" (or custom fallback)
 */
export function formatRelativeLastLogin(
  dateStr?: string | null,
  fallback = 'Never logged in'
): string {
  if (!dateStr) return fallback;

  const date = new Date(dateStr);
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();

  if (isNaN(diffMs) || diffMs < 0) {
    return 'Active now';
  }

  const diffSec = Math.floor(diffMs / 1000);
  const diffMin = Math.floor(diffSec / 60);
  const diffHours = Math.floor(diffMin / 60);
  const diffDays = Math.floor(diffHours / 24);

  // Under 1 minute -> "Active now"
  if (diffMin < 1) {
    return 'Active now';
  }

  // 1 to 59 minutes -> "X minute(s) ago"
  if (diffMin < 60) {
    return diffMin === 1 ? '1 minute ago' : `${diffMin} minutes ago`;
  }

  // 1 to 12 hours -> "X hour(s) ago"
  if (diffHours <= 12) {
    return diffHours === 1 ? '1 hour ago' : `${diffHours} hours ago`;
  }

  // More than 12 hours up to 7 days -> "X day(s) ago"
  if (diffDays <= 7) {
    const days = Math.max(1, diffDays);
    return days === 1 ? '1 day ago' : `${days} days ago`;
  }

  // Over 7 days -> Formatted date
  return date.toLocaleDateString('en-GB', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
  });
}
