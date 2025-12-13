import axios, { AxiosResponse } from 'axios';

export interface CraftAtom {
  name: string;
  description?: string;
  template_name: string;
  params: Record<string, string>;
}

const adminApiBase = '/api/admin';
const publicApiBase = '/api';

// Define the API base URL
const craftAtomApiBase = `${adminApiBase}/craft-atoms`;

// Create a CraftAtom
export function createCraftAtom(
  craftAtom: CraftAtom
): Promise<AxiosResponse<CraftAtom>> {
  return axios.post<CraftAtom>(craftAtomApiBase, craftAtom);
}

// Get a CraftAtom by name
export function getCraftAtom(name: string): Promise<AxiosResponse<CraftAtom>> {
  return axios.get<CraftAtom>(`${craftAtomApiBase}/${name}`);
}

// Update a CraftAtom
export function updateCraftAtom(
  name: string,
  craftAtom: CraftAtom
): Promise<AxiosResponse<CraftAtom>> {
  return axios.put<CraftAtom>(`${craftAtomApiBase}/${name}`, craftAtom);
}

// Delete a CraftAtom
export function deleteCraftAtom(name: string): Promise<AxiosResponse<void>> {
  return axios.delete<void>(`${craftAtomApiBase}/${name}`);
}

// List all CraftAtoms
export function listCraftAtoms(): Promise<AxiosResponse<CraftAtom[]>> {
  return axios.get<CraftAtom[]>(`${publicApiBase}/craft-atoms`);
}
