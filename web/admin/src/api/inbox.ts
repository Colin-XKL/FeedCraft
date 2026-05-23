import axios from 'axios';
import { APIResponse } from '@/api/types';

export interface Inbox {
  id: string;
  title?: string;
  description?: string;
  max_items: number;
  is_public: boolean;
}

export interface SystemAuthToken {
  id: number;
  token: string;
  label?: string;
  created_at: string;
}

const adminApiBase = '/api/admin';

// Inboxes CRUD
export function createInbox(data: Inbox): Promise<APIResponse<Inbox>> {
  return axios
    .post<APIResponse<Inbox>>(`${adminApiBase}/inboxes`, data)
    .then((res) => res.data);
}

export function listInboxes(): Promise<APIResponse<Inbox[]>> {
  return axios
    .get<APIResponse<Inbox[]>>(`${adminApiBase}/inboxes`)
    .then((res) => res.data);
}

export function getInbox(id: string): Promise<APIResponse<Inbox>> {
  return axios
    .get<APIResponse<Inbox>>(`${adminApiBase}/inboxes/${id}`)
    .then((res) => res.data);
}

export function updateInbox(
  id: string,
  data: Inbox
): Promise<APIResponse<Inbox>> {
  return axios
    .put<APIResponse<Inbox>>(`${adminApiBase}/inboxes/${id}`, data)
    .then((res) => res.data);
}

export function deleteInbox(id: string): Promise<APIResponse<void>> {
  return axios
    .delete<APIResponse<void>>(`${adminApiBase}/inboxes/${id}`)
    .then((res) => res.data);
}

// System Auth Tokens CRUD
export function createSystemAuthToken(data: {
  label: string;
}): Promise<APIResponse<SystemAuthToken>> {
  return axios
    .post<APIResponse<SystemAuthToken>>(
      `${adminApiBase}/system-auth-tokens`,
      data
    )
    .then((res) => res.data);
}

export function listSystemAuthTokens(): Promise<
  APIResponse<SystemAuthToken[]>
> {
  return axios
    .get<APIResponse<SystemAuthToken[]>>(`${adminApiBase}/system-auth-tokens`)
    .then((res) => res.data);
}

export function deleteSystemAuthToken(id: number): Promise<APIResponse<void>> {
  return axios
    .delete<APIResponse<void>>(`${adminApiBase}/system-auth-tokens/${id}`)
    .then((res) => res.data);
}
