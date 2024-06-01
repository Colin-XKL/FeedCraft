import axios, { AxiosResponse } from 'axios';

export interface CraftFlowItem {
  craftName: string;
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
  return axios.get<CraftFlow[]>(craftFlowApiBase);
}
