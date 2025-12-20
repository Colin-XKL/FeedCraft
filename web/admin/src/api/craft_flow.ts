import axios from 'axios';
import { APIResponse } from './types';

export interface CraftFlowItem {
  craft_name: string;
  options?: Map<string, string>;
}

// Define the CraftFlow type
export interface CraftFlow {
  craftList?: string[]; // 额外添加的字段, 暂时存储craft 的数组,而不是flow config中的嵌套结构
  name: string;
  description?: string;
  craft_flow_config?: CraftFlowItem[];
}

const adminApiBase = '/api/admin';

// Define the API base URL
const craftFlowApiBase = `${adminApiBase}/craft-flows`;

// Create a CraftFlow
export function createCraftFlow(
  craftFlow: CraftFlow
): Promise<APIResponse<CraftFlow>> {
  return axios.post<APIResponse<CraftFlow>>(craftFlowApiBase, craftFlow).then(res => res.data);
}

// Get a CraftFlow by name
export function getCraftFlow(
  name: string
): Promise<APIResponse<CraftFlow>> {
  return axios.get<APIResponse<CraftFlow>>(`${craftFlowApiBase}/${name}`).then(res => res.data);
}

// Update a CraftFlow
export function updateCraftFlow(
  name: string,
  craftFlow: CraftFlow
): Promise<APIResponse<CraftFlow>> {
  return axios.put<APIResponse<CraftFlow>>(
    `${craftFlowApiBase}/${name}`,
    craftFlow
  ).then(res => res.data);
}

// Delete a CraftFlow
export function deleteCraftFlow(
  name: string
): Promise<APIResponse<void>> {
  return axios.delete<APIResponse<void>>(`${craftFlowApiBase}/${name}`).then(res => res.data);
}

// List all CraftFlows
export function listCraftFlows(): Promise<
  APIResponse<CraftFlow[]>
> {
  return axios.get<APIResponse<CraftFlow[]>>(craftFlowApiBase).then(res => res.data);
}

export function listSysCraftAtoms(): Promise<
  APIResponse<{ name: string; description: string }[]>
> {
  return axios.get<
    APIResponse<{ name: string; description: string }[]>
  >(`${adminApiBase}/sys-craft-atoms`).then(res => res.data);
}

interface ParamTemplate {
  key: string;
  description: string;
  default: string;
}

interface CraftTemplate {
  name: string;
  description?: string;
  param_template_define: ParamTemplate[];
}

export function listCraftTemplates(): Promise<
  APIResponse<CraftTemplate[]>
> {
  return axios.get<APIResponse<CraftTemplate[]>>(
    `${adminApiBase}/craft-templates`
  ).then(res => res.data);
}
