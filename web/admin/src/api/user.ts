import axios from 'axios';
import { md5 } from 'js-md5';
import type { RouteRecordNormalized } from 'vue-router';
import { UserState } from '@/store/modules/user/types';
import { APIResponse } from './types';

export interface LoginData {
  username: string;
  password: string;
}

export interface LoginRes {
  token: string;
}

export function login(data: LoginData): Promise<APIResponse<LoginRes>> {
  const md5Password = md5(data.password).toString();
  return axios.post<APIResponse<LoginRes>>('/api/login', {
    username: data.username,
    md5_password: md5Password,
  }).then(res => res.data);
}

export function logout(): Promise<APIResponse<LoginRes>> {
  return axios.post<APIResponse<LoginRes>>('/api/user/logout').then(res => res.data);
}

export function getUserInfo(): Promise<APIResponse<UserState>> {
  return axios.post<APIResponse<UserState>>('/api/admin/user/info').then(res => res.data);
}

export function getMenuList(): Promise<APIResponse<RouteRecordNormalized[]>> {
  return axios.post<APIResponse<RouteRecordNormalized[]>>('/api/user/menu').then(res => res.data);
}

export interface ChangePasswordData {
  username: string;
  currentPassword: string;
  newPassword: string;
}

export function changePassword(data: ChangePasswordData): Promise<APIResponse<any>> {
  const currentPasswordMd5 = md5(data.currentPassword).toString();
  const newPasswordMd5 = md5(data.newPassword).toString();
  return axios.post<APIResponse<any>>('/api/admin/user/change-password', {
    username: data.username,
    currentPassword: currentPasswordMd5,
    newPassword: newPasswordMd5,
  }).then(res => res.data);
}
