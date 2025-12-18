import axios, { AxiosResponse } from 'axios';

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
const publicApiBase = '/api';

// Define the API base URL
const craftFlowApiBase = `${adminApiBase}/craft-flows`;

// Create a CraftFlow
export function createCraftFlow(
  craftFlow: CraftFlow
): Promise<AxiosResponse<CraftFlow>> {
  return axios.post<CraftFlow>(craftFlowApiBase, craftFlow);
}

// Get a CraftFlow by name
export function getCraftFlow(name: string): Promise<AxiosResponse<CraftFlow>> {
  return axios.get<CraftFlow>(`${craftFlowApiBase}/${name}`);
}

// Update a CraftFlow
export function updateCraftFlow(
  name: string,
  craftFlow: CraftFlow
): Promise<AxiosResponse<CraftFlow>> {
  return axios.put<CraftFlow>(`${craftFlowApiBase}/${name}`, craftFlow);
}

// Delete a CraftFlow
export function deleteCraftFlow(name: string): Promise<AxiosResponse<void>> {
  return axios.delete<void>(`${craftFlowApiBase}/${name}`);
}

// List all CraftFlows
export function listCraftFlows(): Promise<AxiosResponse<CraftFlow[]>> {
  return axios.get<CraftFlow[]>(`${publicApiBase}/craft-flows`);
}

export function listSysCraftAtoms(): Promise<
  AxiosResponse<{ name: string; description: string }[]>
> {
  return axios.get<{ name: string; description: string }[]>(
    `${publicApiBase}/sys-craft-atoms`
  );
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

export function listCraftTemplates(): Promise<AxiosResponse<CraftTemplate[]>> {
  return axios.get<CraftTemplate[]>(`${adminApiBase}/craft-templates`);
}
