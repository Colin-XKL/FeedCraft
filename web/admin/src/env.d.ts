/// <reference types="vite/client" />

declare module '*.vue' {
  import type { DefineComponent } from 'vue';

  const component: DefineComponent;
  export default component;
}
interface ImportMetaEnv {
  readonly VITE_API_BASE_URL: string;
}

// Build-time define holding the resolved display version (e.g. v3.1.0, dev-becb6a3).
// eslint-disable-next-line no-underscore-dangle
declare const __APP_VERSION__: string;
