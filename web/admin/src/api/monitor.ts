import axios from 'axios';
import { APIResponse } from './types';

export interface DependencyStatus {
  name: string;
  status: string; // "Healthy", "Unhealthy", "Not Configured"
  details: string;
  error?: string;
  latency?: string;
}

export function fetchDependencyStatus() {
  // Use generic type to define the shape of APIResponse data
  return axios.get<APIResponse<DependencyStatus[]>>('/api/admin/dependencies');
}
