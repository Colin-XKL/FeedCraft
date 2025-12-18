import axios, { AxiosResponse } from 'axios';
import { APIResponse } from './types';

export interface CraftAtom {
  name: string;
  description?: string;
  template_name: string;
  params: Record<string, string>;
}

const adminApiBase = '/api/admin';

// Define the API base URL
const craftAtomApiBase = `${adminApiBase}/craft-atoms`;

// Create a CraftAtom
export function createCraftAtom(
  craftAtom: CraftAtom
): Promise<APIResponse<CraftAtom>> {
  return axios.post<APIResponse<CraftAtom>>(craftAtomApiBase, craftAtom) as unknown as Promise<APIResponse<CraftAtom>>;
}

// Get a CraftAtom by name
export function getCraftAtom(
  name: string
): Promise<APIResponse<CraftAtom>> {
  return axios.get<APIResponse<CraftAtom>>(`${craftAtomApiBase}/${name}`) as unknown as Promise<APIResponse<CraftAtom>>;
}

// Update a CraftAtom
export function updateCraftAtom(
  name: string,
  craftAtom: CraftAtom
): Promise<APIResponse<CraftAtom>> {
  return axios.put<APIResponse<CraftAtom>>(
    `${craftAtomApiBase}/${name}`,
    craftAtom
  ) as unknown as Promise<APIResponse<CraftAtom>>;
}

// Delete a CraftAtom
export function deleteCraftAtom(
  name: string
): Promise<APIResponse<void>> {
  return axios.delete<APIResponse<void>>(`${craftAtomApiBase}/${name}`) as unknown as Promise<APIResponse<void>>;
}

// List all CraftAtoms
export function listCraftAtoms(): Promise<
  APIResponse<CraftAtom[]>
> {
  return axios.get<APIResponse<CraftAtom[]>>(craftAtomApiBase) as unknown as Promise<APIResponse<CraftAtom[]>>;
}
