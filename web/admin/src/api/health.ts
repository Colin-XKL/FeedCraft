import axios from 'axios';
import { APIResponse } from '@/api/types';

export interface DependencyNode {
  name: string;
  type: string; // recipe, flow, atom, built-in, missing, cycle
  exists: boolean;
  children?: DependencyNode[];
  details?: string;
  key: string;
}

export interface InboxGCStats {
  total_items: number;
  orphaned_count: number;
  overflow_count: number;
}

export interface GCCleanupResult {
  orphaned_deleted: number;
  overflow_deleted: number;
}

export function fetchDependencyHealth() {
  return axios
    .get<APIResponse<DependencyNode[]>>('/api/admin/dependencies/health')
    .then((res) => res.data);
}

export function fetchInboxGCStats() {
  return axios
    .get<APIResponse<InboxGCStats>>('/api/admin/inboxes/gc/stats')
    .then((res) => res.data);
}

export function triggerInboxGCCleanup() {
  return axios
    .post<APIResponse<GCCleanupResult>>('/api/admin/inboxes/gc/cleanup')
    .then((res) => res.data);
}
