import axios from 'axios';
import { APIResponse } from './types';

export interface DependencyStatus {
  name: string;
  status: string; // "Healthy", "Unhealthy", "Configured", "Not Configured"
  details: string;
  error?: string;
  latency?: string;
}

export function fetchDependencyStatus(): Promise<APIResponse<DependencyStatus[]>> {
  return axios
    .get<APIResponse<DependencyStatus[]>>('/api/admin/dependencies')
    .then((res) => res.data);
}

export function checkDependencyStatus(): Promise<APIResponse<DependencyStatus[]>> {
  return axios
    .post<APIResponse<DependencyStatus[]>>('/api/admin/dependencies/check')
    .then((res) => res.data);
}
