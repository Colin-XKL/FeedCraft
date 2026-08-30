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
  SessionExpiredError,
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
    content: `${SESSION_EXPIRED_MESSAGE}你可以先关闭弹窗留在当前页，复制或保存未完成的内容；也可以点击“前往登录”重新登录。`,
    closable: true,
    hideCancel: false,
    okText: '前往登录',
    cancelText: '留在当前页',
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
      return Promise.reject(new SessionExpiredError());
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
      return Promise.reject(new SessionExpiredError());
    }

    const respMsg = error?.response?.data?.msg;
    if (respMsg && typeof respMsg === 'string') {
      error.message = respMsg;
    }
    Message.error({
      content: respMsg || error.message || 'Request Error',
      duration: 5 * 1000,
    });
    return Promise.reject(error);
  }
);
