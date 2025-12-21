import axios from 'axios';
import { APIResponse } from './types';

export interface DependencyStatus {
  name: string;
  status: string; // "Healthy", "Unhealthy", "Configured", "Not Configured"
  details: string;
  error?: string;
  latency?: string;
}

export function fetchDependencyStatus() {
  return axios.get<APIResponse<DependencyStatus[]>>('/api/admin/dependencies');
}

export function checkDependencyStatus() {
  return axios.post<APIResponse<DependencyStatus[]>>(
    '/api/admin/dependencies/check',
  );
}
