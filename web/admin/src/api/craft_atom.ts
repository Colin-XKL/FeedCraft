import axios from 'axios';
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
  return axios
    .post<APIResponse<CraftAtom>>(craftAtomApiBase, craftAtom)
    .then((res) => res.data);
}

// Get a CraftAtom by name
export function getCraftAtom(name: string): Promise<APIResponse<CraftAtom>> {
  return axios
    .get<APIResponse<CraftAtom>>(`${craftAtomApiBase}/${name}`)
    .then((res) => res.data);
}

// Update a CraftAtom
export function updateCraftAtom(
  name: string,
  craftAtom: CraftAtom
): Promise<APIResponse<CraftAtom>> {
  return axios
    .put<APIResponse<CraftAtom>>(`${craftAtomApiBase}/${name}`, craftAtom)
    .then((res) => res.data);
}

// Delete a CraftAtom
export function deleteCraftAtom(name: string): Promise<APIResponse<void>> {
  return axios
    .delete<APIResponse<void>>(`${craftAtomApiBase}/${name}`)
    .then((res) => res.data);
}

// List all CraftAtoms
export function listCraftAtoms(): Promise<APIResponse<CraftAtom[]>> {
  return axios
    .get<APIResponse<CraftAtom[]>>(craftAtomApiBase)
    .then((res) => res.data);
}
