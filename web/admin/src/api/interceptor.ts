import axios from 'axios';
import type { AxiosRequestConfig, AxiosResponse } from 'axios';
import { Message, Modal } from '@arco-design/web-vue';
import { useUserStore } from '@/store';
import { getToken } from '@/utils/auth';
import router from '@/router';
import { APIResponse } from '@/api/types';
import {
  buildSessionExpiredRedirectQuery,
  isSessionExpiredAPIResponse,
  isSessionExpiredHTTPStatus,
  SESSION_EXPIRED_MESSAGE,
} from '@/api/auth-expired';

if (import.meta.env.VITE_API_BASE_URL) {
  axios.defaults.baseURL = import.meta.env.VITE_API_BASE_URL;
}

let sessionExpiredModalVisible = false;

const redirectToLogin = () => {
  const currentRoute = router.currentRoute.value;
  if (currentRoute.name === 'login') return Promise.resolve();

  return router.push({
    name: 'login',
    query: buildSessionExpiredRedirectQuery(
      currentRoute.query,
      currentRoute.fullPath
    ),
  });
};

const showSessionExpiredModal = () => {
  const userStore = useUserStore();
  userStore.logoutCallBack();

  if (
    sessionExpiredModalVisible ||
    router.currentRoute.value.name === 'login'
  ) {
    return;
  }

  sessionExpiredModalVisible = true;
  Modal.warning({
    title: '登录态已过期',
    content: `${SESSION_EXPIRED_MESSAGE}点击“前往登录”可立即重新登录。`,
    okText: '前往登录',
    cancelText: '稍后处理',
    onOk() {
      sessionExpiredModalVisible = false;
      return redirectToLogin();
    },
    onCancel() {
      sessionExpiredModalVisible = false;
    },
    onClose() {
      sessionExpiredModalVisible = false;
    },
  });
};

axios.interceptors.request.use(
  (config: AxiosRequestConfig) => {
    // let each request carry token
    // this example using the JWT token
    // Authorization is a custom headers key
    // please modify it according to the actual situation
    const token = getToken();
    if (token) {
      if (!config.headers) {
        config.headers = {};
      }
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    // do something
    return Promise.reject(error);
  }
);
// add response interceptors
axios.interceptors.response.use(
  (response: AxiosResponse<APIResponse>) => {
    const res = response.data;
    if (isSessionExpiredAPIResponse(res)) {
      showSessionExpiredModal();
      return Promise.reject(new Error(SESSION_EXPIRED_MESSAGE));
    }
    if (res.code !== 0) {
      Message.error({
        content: res.msg || 'Error',
        duration: 5 * 1000,
      });
      return Promise.reject(new Error(res.msg || 'Error'));
    }
    return response;
  },
  (error) => {
    if (isSessionExpiredHTTPStatus(error?.response?.status)) {
      showSessionExpiredModal();
      return Promise.reject(new Error(SESSION_EXPIRED_MESSAGE));
    }

    const respMsg = error?.response?.data?.msg;
    Message.error({
      content: respMsg || error.message || 'Request Error',
      duration: 5 * 1000,
    });
    return Promise.reject(error);
  }
);
