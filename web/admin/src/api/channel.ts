import axios from 'axios';
import { APIResponse } from './types';

export interface Channel {
  id: string;
  description?: string;
  processor_name: string;
  source_type: string;
  source_config: string;
  is_active?: boolean;
  last_accessed_at?: string;
}

const adminApiBase = '/api/admin';

export function getChannels(): Promise<APIResponse<Channel[]>> {
  return axios
    .get<APIResponse<Channel[]>>(`${adminApiBase}/channels`)
    .then((res) => res.data);
}

export function getChannelById(id: string): Promise<APIResponse<Channel>> {
  return axios
    .get<APIResponse<Channel>>(`${adminApiBase}/channels/${id}`)
    .then((res) => res.data);
}

export function createChannel(data: Channel): Promise<APIResponse<Channel>> {
  return axios
    .post<APIResponse<Channel>>(`${adminApiBase}/channels`, data)
    .then((res) => res.data);
}

export function updateChannel(data: Channel): Promise<APIResponse<Channel>> {
  return axios
    .put<APIResponse<Channel>>(`${adminApiBase}/channels/${data.id}`, data)
    .then((res) => res.data);
}

export function deleteChannel(id: string): Promise<APIResponse<void>> {
  return axios
    .delete<APIResponse<void>>(`${adminApiBase}/channels/${id}`)
    .then((res) => res.data);
}
