import localeMessageBox from '@/components/message-box/locale/zh-CN';
import localeLogin from '@/views/login/locale/zh-CN';

import localeWorkplace from '@/views/dashboard/workplace/locale/zh-CN';

import localeSettings from './zh-CN/settings';
import localeMenu from './zh-CN/menu';
import localeCraftAtom from './zh-CN/craftAtom';
import localeCraftFlow from './zh-CN/craftFlow';
import localeCustomRecipe from './zh-CN/customRecipe';
import localeAllCraftList from './zh-CN/allCraftList';
import localeFeedCompare from './zh-CN/feedCompare';
import localeFeedViewer from './zh-CN/feedViewer';
import localeLlmDebug from './zh-CN/llmDebug';
import localeUrlGenerator from './zh-CN/urlGenerator';

export default {
  ...localeSettings,
  ...localeMessageBox,
  ...localeLogin,
  ...localeMenu,
  ...localeCraftAtom,
  ...localeCraftFlow,
  ...localeCustomRecipe,
  ...localeAllCraftList,
  ...localeFeedCompare,
  ...localeFeedViewer,
  ...localeLlmDebug,
  ...localeUrlGenerator,
};
