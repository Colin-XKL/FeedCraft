import axios from 'axios';
import { md5 } from 'js-md5';
import type { RouteRecordNormalized } from 'vue-router';
import { UserState } from '@/store/modules/user/types';

export interface LoginData {
  username: string;
  password: string;
}

export interface LoginRes {
  token: string;
}

export function login(data: LoginData) {
  const md5Password = md5(data.password).toString();
  return axios.post<LoginRes>('/api/login', {
    username: data.username,
    md5_password: md5Password,
  });
}

export function logout() {
  return axios.post<LoginRes>('/api/user/logout');
}

export function getUserInfo() {
  return axios.post<UserState>('/api/admin/user/info');
}

export function getMenuList() {
  return axios.post<RouteRecordNormalized[]>('/api/user/menu');
}
export interface ChangePasswordData {
  username: string;
  currentPassword: string;
  newPassword: string;
}

import { md5 } from 'js-md5';

export function changePassword(data: ChangePasswordData) {
  const currentPasswordMd5 = md5(data.currentPassword).toString();
  const newPasswordMd5 = md5(data.newPassword).toString();
  return axios.post('/api/admin/user/change-password', {
    username: data.username,
    currentPassword: currentPasswordMd5,
    newPassword: newPasswordMd5,
  });
}
