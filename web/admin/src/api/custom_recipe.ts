import axios from 'axios';
import { APIResponse } from './types';

export interface CustomRecipe {
  id: string;
  description?: string;
  craft: string;
  // feed_url: string; // Deprecated
  source_type: string;
  source_config: string;
  is_active?: boolean;
  last_accessed_at?: string;
}

const adminApiBase = '/api/admin';

export function getCustomRecipes(): Promise<APIResponse<CustomRecipe[]>> {
  return axios.get<APIResponse<CustomRecipe[]>>(`${adminApiBase}/recipes`) as unknown as Promise<APIResponse<CustomRecipe[]>>;
}

export function getCustomRecipeById(id: string): Promise<APIResponse<CustomRecipe>> {
  return axios.get<APIResponse<CustomRecipe>>(
    `${adminApiBase}/recipes/${id}`
  ) as unknown as Promise<APIResponse<CustomRecipe>>;
}

export function createCustomRecipe(data: CustomRecipe): Promise<APIResponse<CustomRecipe>> {
  return axios.post<APIResponse<CustomRecipe>>(
    `${adminApiBase}/recipes`,
    data
  ) as unknown as Promise<APIResponse<CustomRecipe>>;
}

export function updateCustomRecipe(data: CustomRecipe): Promise<APIResponse<CustomRecipe>> {
  return axios.put<APIResponse<CustomRecipe>>(
    `${adminApiBase}/recipes/${data.id}`,
    data
  ) as unknown as Promise<APIResponse<CustomRecipe>>;
}

export function deleteCustomRecipe(id: string): Promise<APIResponse<void>> {
  return axios.delete<APIResponse<void>>(`${adminApiBase}/recipes/${id}`) as unknown as Promise<APIResponse<void>>;
}
