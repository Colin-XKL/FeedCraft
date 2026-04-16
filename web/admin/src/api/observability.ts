import axios from 'axios';
import { APIResponse } from '@/api/types';

export interface ObservableResource {
  resource_type: string;
  resource_id: string;
  resource_name: string;
  current_status: string;
  consecutive_failures: number;
  last_success_at?: string;
  last_failure_at?: string;
  last_error_kind?: string;
  last_error_message?: string;
  paused_at?: string;
  paused_reason?: string;
}

export interface ExecutionLog {
  id: number;
  resource_type: string;
  resource_id: string;
  resource_name: string;
  trigger: string;
  status: string;
  error_kind: string;
  message: string;
  details_json: string;
  details?: any;
  request_id: string;
  duration_ms: number;
  created_at: string;
}

export interface SystemNotification {
  id: number;
  resource_type: string;
  resource_id: string;
  event_type: string;
  title: string;
  content: string;
  status_after: string;
  created_at: string;
}

export function fetchObservableResources(params?: Record<string, any>) {
  return axios
    .get<APIResponse<ObservableResource[]>>(
      '/api/admin/observability/resources',
      {
        params,
      }
    )
    .then((res) => res.data);
}

export function fetchExecutionLogs(params?: Record<string, any>) {
  return axios
    .get<APIResponse<ExecutionLog[]>>('/api/admin/observability/executions', {
      params,
    })
    .then((res) => res.data);
}

export function fetchSystemNotifications() {
  return axios
    .get<APIResponse<SystemNotification[]>>('/api/admin/system-notifications')
    .then((res) => res.data);
}

export function resumeObservableResource(
  resourceType: string,
  resourceID: string
) {
  return axios
    .post<APIResponse>(
      `/api/admin/observability/resources/${resourceType}/${resourceID}/resume`
    )
    .then((res) => res.data);
}
