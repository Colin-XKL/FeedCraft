import axios from 'axios';

export interface DependencyStatus {
  name: string;
  configured: boolean;
  healthy: boolean;
  message: string;
}

export function getDependencyStatus() {
  return axios.get<DependencyStatus[]>('/api/admin/dependencies');
}
