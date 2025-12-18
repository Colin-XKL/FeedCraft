import axios, { AxiosResponse } from 'axios';
import { md5 } from 'js-md5';
import type { RouteRecordNormalized } from 'vue-router';
import { UserState } from '@/store/modules/user/types';
import { util } from './interceptor';

export interface LoginData {
  username: string;
  password: string;
}

export interface LoginRes {
  token: string;
}

export function login(data: LoginData): Promise<util.APIResponse<LoginRes>> {
  const md5Password = md5(data.password).toString();
  return axios.post<util.APIResponse<LoginRes>>('/api/login', {
    username: data.username,
    md5_password: md5Password,
  }) as unknown as Promise<util.APIResponse<LoginRes>>;
}

export function logout(): Promise<util.APIResponse<LoginRes>> {
  return axios.post<util.APIResponse<LoginRes>>('/api/user/logout') as unknown as Promise<util.APIResponse<LoginRes>>;
}

export function getUserInfo(): Promise<util.APIResponse<UserState>> {
  return axios.post<util.APIResponse<UserState>>('/api/admin/user/info') as unknown as Promise<util.APIResponse<UserState>>;
}

export function getMenuList(): Promise<util.APIResponse<RouteRecordNormalized[]>> {
  return axios.post<util.APIResponse<RouteRecordNormalized[]>>('/api/user/menu') as unknown as Promise<util.APIResponse<RouteRecordNormalized[]>>;
}

export interface ChangePasswordData {
  username: string;
  currentPassword: string;
  newPassword: string;
}

export function changePassword(data: ChangePasswordData): Promise<util.APIResponse<any>> {
  const currentPasswordMd5 = md5(data.currentPassword).toString();
  const newPasswordMd5 = md5(data.newPassword).toString();
  return axios.post<util.APIResponse<any>>('/api/admin/user/change-password', {
    username: data.username,
    currentPassword: currentPasswordMd5,
    newPassword: newPasswordMd5,
  }) as unknown as Promise<util.APIResponse<any>>;
}
