import { createPinia } from 'pinia';
import useAppStore from '@/store/modules/app';
import useUserStore from '@/store/modules/user';
import useTabBarStore from '@/store/modules/tab-bar';

const pinia = createPinia();

export { useAppStore, useUserStore, useTabBarStore };
export default pinia;
