export interface AuditLog {
  id: string;
  category: string;
  action: string;
  resource: string;
  resource_id: string;
  actor_id: string;
  outcome: string;
  request_id: string;
  client_ip: string;
  metadata: {
    user_agent?: string;
    [key: string]: any;
  };
  created_at: string;
}

export interface AuditLogsResponse {
  audit_logs: AuditLog[];
  page: number;
  limit: number;
  total: number;
}
