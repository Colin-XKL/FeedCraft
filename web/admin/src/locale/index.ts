import { createI18n } from 'vue-i18n';
import en from '@/locale/en-US';
import cn from '@/locale/zh-CN';
import tw from '@/locale/zh-TW';

export const LOCALE_OPTIONS = [
  { label: '中文', value: 'zh-CN' },
  { label: '繁體中文', value: 'zh-TW' },
  { label: 'English', value: 'en-US' },
];
const defaultLocale = localStorage.getItem('arco-locale') || 'zh-CN';

const i18n = createI18n({
  locale: defaultLocale,
  fallbackLocale: 'en-US',
  legacy: false,
  allowComposition: true,
  messages: {
    'en-US': en,
    'zh-CN': cn,
    'zh-TW': tw,
  },
});

export default i18n;
