import axios from 'axios';
import { APIResponse } from './types';

export interface DependencyNode {
  name: string;
  type: string; // recipe, flow, atom, built-in, missing, cycle
  exists: boolean;
  children?: DependencyNode[];
  details?: string;
  key: string;
}

export function fetchDependencyHealth() {
  return axios
    .get<APIResponse<DependencyNode[]>>('/api/admin/dependencies/health')
    .then((res) => res.data);
}
