import axios from 'axios';

export interface CustomRecipe {
  id: string;
  description?: string;
  craft: string;
  feed_url: string;
  is_active?: boolean;
  last_accessed_at?: string;
}

const adminApiBase = '/api/admin';

export function getCustomRecipes() {
  return axios.get<CustomRecipe[]>(`${adminApiBase}/recipes`);
}

export function getCustomRecipeById(id: string) {
  return axios.get<CustomRecipe>(`${adminApiBase}/recipes/${id}`);
}

export function createCustomRecipe(data: CustomRecipe) {
  return axios.post<CustomRecipe>(`${adminApiBase}/recipes`, data);
}

export function updateCustomRecipe(data: CustomRecipe) {
  return axios.put<CustomRecipe>(`${adminApiBase}/recipes/${data.id}`, data);
}

export function deleteCustomRecipe(id: string): Promise<void> {
  return axios.delete(`${adminApiBase}/recipes/${id}`);
}
