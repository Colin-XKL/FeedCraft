import axios, { AxiosResponse } from 'axios';
import { util } from './interceptor';

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
): Promise<util.APIResponse<CraftAtom>> {
  return axios.post<util.APIResponse<CraftAtom>>(craftAtomApiBase, craftAtom) as unknown as Promise<util.APIResponse<CraftAtom>>;
}

// Get a CraftAtom by name
export function getCraftAtom(
  name: string
): Promise<util.APIResponse<CraftAtom>> {
  return axios.get<util.APIResponse<CraftAtom>>(`${craftAtomApiBase}/${name}`) as unknown as Promise<util.APIResponse<CraftAtom>>;
}

// Update a CraftAtom
export function updateCraftAtom(
  name: string,
  craftAtom: CraftAtom
): Promise<util.APIResponse<CraftAtom>> {
  return axios.put<util.APIResponse<CraftAtom>>(
    `${craftAtomApiBase}/${name}`,
    craftAtom
  ) as unknown as Promise<util.APIResponse<CraftAtom>>;
}

// Delete a CraftAtom
export function deleteCraftAtom(
  name: string
): Promise<util.APIResponse<void>> {
  return axios.delete<util.APIResponse<void>>(`${craftAtomApiBase}/${name}`) as unknown as Promise<util.APIResponse<void>>;
}

// List all CraftAtoms
export function listCraftAtoms(): Promise<
  util.APIResponse<CraftAtom[]>
> {
  return axios.get<util.APIResponse<CraftAtom[]>>(craftAtomApiBase) as unknown as Promise<util.APIResponse<CraftAtom[]>>;
}
