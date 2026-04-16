import { App } from 'vue';
import permission from '@/directive/permission';

export default {
  install(Vue: App) {
    Vue.directive('permission', permission);
  },
};
