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
import localeHtmlToRss from './zh-CN/htmlToRss';
import localeCurlToRss from './zh-CN/curlToRss';
import localeSearchToRss from './zh-CN/searchToRss';
import localeDependencyService from './zh-CN/dependencyService';

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
  ...localeHtmlToRss,
  ...localeCurlToRss,
  ...localeSearchToRss,
  ...localeDependencyService,
  ...localeWorkplace,
};
