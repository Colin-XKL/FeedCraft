import { describe, expect, it } from 'vitest';
import {
  craftAtomModalBodyStyle,
  craftAtomModalLayout,
  modalHeightWithinViewport,
} from '@/views/dashboard/craft_atom/modalLayout';

const reservedPx =
  craftAtomModalLayout.headerPx +
  craftAtomModalLayout.footerPx +
  craftAtomModalLayout.verticalInsetPx;

describe('craft atom modal layout', () => {
  it('caps the body so header, footer, and inset stay inside a short viewport', () => {
    expect(craftAtomModalBodyStyle()).toEqual({
      maxHeight: `calc(100vh - ${reservedPx}px)`,
      overflowY: 'auto',
    });

    const shortViewport = 583;
    const tallBody = 822;
    const fitted = modalHeightWithinViewport(shortViewport, tallBody);

    expect(fitted).toBeLessThanOrEqual(
      shortViewport - craftAtomModalLayout.verticalInsetPx
    );
    expect(fitted).toBe(
      craftAtomModalLayout.headerPx +
        (shortViewport - reservedPx) +
        craftAtomModalLayout.footerPx
    );
  });

  it('keeps a short form at its content height', () => {
    expect(modalHeightWithinViewport(900, 200)).toBe(
      craftAtomModalLayout.headerPx + 200 + craftAtomModalLayout.footerPx
    );
  });
});
