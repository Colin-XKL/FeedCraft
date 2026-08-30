/// <reference types="vite/client" />

declare module '*.vue' {
  import type { DefineComponent } from 'vue';

  const component: DefineComponent;
  export default component;
}
interface ImportMetaEnv {
  readonly VITE_API_BASE_URL: string;
  readonly VITE_APP_VERSION?: string;
}

// Vite compile-time define; conventional double-underscore name.
// eslint-disable-next-line no-underscore-dangle
declare const __APP_VERSION__: string;
