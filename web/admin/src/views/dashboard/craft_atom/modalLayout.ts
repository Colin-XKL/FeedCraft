import type { CSSProperties } from 'vue';

export const craftAtomModalLayout = {
  headerPx: 48,
  footerPx: 65,
  verticalInsetPx: 64,
};

const reservedPx =
  craftAtomModalLayout.headerPx +
  craftAtomModalLayout.footerPx +
  craftAtomModalLayout.verticalInsetPx;

export function craftAtomModalBodyStyle(): CSSProperties {
  return {
    maxHeight: `calc(100vh - ${reservedPx}px)`,
    overflowY: 'auto',
  };
}

export function modalHeightWithinViewport(
  viewportHeight: number,
  contentBodyHeight: number
) {
  const maxBody = Math.max(0, viewportHeight - reservedPx);
  const body = Math.min(Math.max(0, contentBodyHeight), maxBody);
  return craftAtomModalLayout.headerPx + body + craftAtomModalLayout.footerPx;
}
