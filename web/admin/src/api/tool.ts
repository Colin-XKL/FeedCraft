import axios from 'axios';
import { APIResponse } from './types';

export interface Tool {
  name: string;
  description?: string;
  template_name: string;
  params: Record<string, string>;
}

const adminApiBase = '/api/admin';

// Define the API base URL
const toolApiBase = `${adminApiBase}/tools`;

// Create a Tool
export function createTool(tool: Tool): Promise<APIResponse<Tool>> {
  return axios
    .post<APIResponse<Tool>>(toolApiBase, tool)
    .then((res) => res.data);
}

// Get a Tool by name
export function getTool(name: string): Promise<APIResponse<Tool>> {
  return axios
    .get<APIResponse<Tool>>(`${toolApiBase}/${name}`)
    .then((res) => res.data);
}

// Update a Tool
export function updateTool(
  name: string,
  tool: Tool,
): Promise<APIResponse<Tool>> {
  return axios
    .put<APIResponse<Tool>>(`${toolApiBase}/${name}`, tool)
    .then((res) => res.data);
}

// Delete a Tool
export function deleteTool(name: string): Promise<APIResponse<void>> {
  return axios
    .delete<APIResponse<void>>(`${toolApiBase}/${name}`)
    .then((res) => res.data);
}

// List all Tools
export function listTools(): Promise<APIResponse<Tool[]>> {
  return axios.get<APIResponse<Tool[]>>(toolApiBase).then((res) => res.data);
}
