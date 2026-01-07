import axios from 'axios';
import { APIResponse } from './types';

export interface BlueprintItem {
  processor_name: string;
  options?: Map<string, string>;
}

// Define the Blueprint type
export interface Blueprint {
  processorList?: string[]; // 额外添加的字段, 暂时存储craft 的数组,而不是flow config中的嵌套结构
  name: string;
  description?: string;
  blueprint_config?: BlueprintItem[];
}

const adminApiBase = '/api/admin';

// Define the API base URL
const blueprintApiBase = `${adminApiBase}/blueprints`;

// Create a Blueprint
export function createBlueprint(
  blueprint: Blueprint,
): Promise<APIResponse<Blueprint>> {
  return axios
    .post<APIResponse<Blueprint>>(blueprintApiBase, blueprint)
    .then((res) => res.data);
}

// Get a Blueprint by name
export function getBlueprint(name: string): Promise<APIResponse<Blueprint>> {
  return axios
    .get<APIResponse<Blueprint>>(`${blueprintApiBase}/${name}`)
    .then((res) => res.data);
}

// Update a Blueprint
export function updateBlueprint(
  name: string,
  blueprint: Blueprint,
): Promise<APIResponse<Blueprint>> {
  return axios
    .put<APIResponse<Blueprint>>(`${blueprintApiBase}/${name}`, blueprint)
    .then((res) => res.data);
}

// Delete a Blueprint
export function deleteBlueprint(name: string): Promise<APIResponse<void>> {
  return axios
    .delete<APIResponse<void>>(`${blueprintApiBase}/${name}`)
    .then((res) => res.data);
}

// List all Blueprints
export function listBlueprints(): Promise<APIResponse<Blueprint[]>> {
  return axios
    .get<APIResponse<Blueprint[]>>(blueprintApiBase)
    .then((res) => res.data);
}

export function listSysTools(): Promise<
  APIResponse<{ name: string; description: string }[]>
> {
  return axios
    .get<
      APIResponse<{ name: string; description: string }[]>
    >(`${adminApiBase}/sys-tools`)
    .then((res) => res.data);
}

interface ParamTemplate {
  key: string;
  description: string;
  default: string;
}

interface ToolTemplate {
  name: string;
  description?: string;
  param_template_define: ParamTemplate[];
}

export function listToolTemplates(): Promise<APIResponse<ToolTemplate[]>> {
  return axios
    .get<APIResponse<ToolTemplate[]>>(`${adminApiBase}/tool-templates`)
    .then((res) => res.data);
}
