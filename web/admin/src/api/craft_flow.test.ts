import { describe, expect, it } from 'vitest';
import { normalizeCraftItems } from '@/api/craft_flow';

describe('craft flow API helpers', () => {
  it('keeps craft item arrays intact', () => {
    const crafts = [
      {
        name: 'proxy',
        description: 'Proxy feed',
        type: '@sys/atom',
        template_only: false,
      },
    ];

    expect(normalizeCraftItems(crafts)).toBe(crafts);
  });

  it('rejects malformed craft list responses', () => {
    expect(() => normalizeCraftItems(null)).toThrow(
      'Expected craft list response data to be an array.'
    );
    expect(() => normalizeCraftItems({ data: [] })).toThrow(
      'Expected craft list response data to be an array.'
    );
  });
});
