import axios from 'axios';
import { util } from './interceptor';

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

export function getCustomRecipes(): Promise<util.APIResponse<CustomRecipe[]>> {
  return axios.get<util.APIResponse<CustomRecipe[]>>(`${adminApiBase}/recipes`) as unknown as Promise<util.APIResponse<CustomRecipe[]>>;
}

export function getCustomRecipeById(id: string): Promise<util.APIResponse<CustomRecipe>> {
  return axios.get<util.APIResponse<CustomRecipe>>(
    `${adminApiBase}/recipes/${id}`
  ) as unknown as Promise<util.APIResponse<CustomRecipe>>;
}

export function createCustomRecipe(data: CustomRecipe): Promise<util.APIResponse<CustomRecipe>> {
  return axios.post<util.APIResponse<CustomRecipe>>(
    `${adminApiBase}/recipes`,
    data
  ) as unknown as Promise<util.APIResponse<CustomRecipe>>;
}

export function updateCustomRecipe(data: CustomRecipe): Promise<util.APIResponse<CustomRecipe>> {
  return axios.put<util.APIResponse<CustomRecipe>>(
    `${adminApiBase}/recipes/${data.id}`,
    data
  ) as unknown as Promise<util.APIResponse<CustomRecipe>>;
}

export function deleteCustomRecipe(id: string): Promise<util.APIResponse<void>> {
  return axios.delete<util.APIResponse<void>>(`${adminApiBase}/recipes/${id}`) as unknown as Promise<util.APIResponse<void>>;
}
