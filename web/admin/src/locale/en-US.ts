import localeMessageBox from '@/components/message-box/locale/en-US';
import localeLogin from '@/views/login/locale/en-US';

import localeWorkplace from '@/views/dashboard/workplace/locale/en-US';

import localeSettings from './en-US/settings';
import localeMenu from './en-US/menu';
import localeCraftAtom from './en-US/craftAtom';
import localeCraftFlow from './en-US/craftFlow';
import localeCustomRecipe from './en-US/customRecipe';
import localeAllCraftList from './en-US/allCraftList';
import localeFeedCompare from './en-US/feedCompare';
import localeFeedViewer from './en-US/feedViewer';
import localeLlmDebug from './en-US/llmDebug';
import localeUrlGenerator from './en-US/urlGenerator';
import localeRssGenerator from './en-US/rssGenerator';

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
  ...localeRssGenerator,
  ...localeWorkplace,
  'menu.jsonRssGenerator': 'JSON RSS Generator',
};
