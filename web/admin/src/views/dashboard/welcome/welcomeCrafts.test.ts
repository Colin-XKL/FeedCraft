import { describe, expect, it } from 'vitest';
import {
  DEFAULT_WELCOME_CRAFT,
  WELCOME_CRAFT_OPTIONS,
  craftsInGroup,
} from '@/views/dashboard/welcome/welcomeCrafts';

describe('welcome craft catalog', () => {
  it('defaults to the proxy craft', () => {
    expect(DEFAULT_WELCOME_CRAFT).toBe('proxy');
    expect(
      WELCOME_CRAFT_OPTIONS.some(
        (craft) => craft.value === DEFAULT_WELCOME_CRAFT
      )
    ).toBe(true);
  });

  it('keeps grouped crafts unique and complete', () => {
    const values = WELCOME_CRAFT_OPTIONS.map((craft) => craft.value);
    expect(new Set(values).size).toBe(values.length);
    expect(craftsInGroup('basic').map((craft) => craft.value)).toEqual([
      'proxy',
      'limit',
      'keyword',
      'guid-fix',
      'relative-link-fix',
    ]);
    expect(craftsInGroup('extract').map((craft) => craft.value)).toEqual([
      'fulltext',
      'fulltext-plus',
      'cleanup',
    ]);
    expect(craftsInGroup('ai').map((craft) => craft.value)).toContain(
      'translate-title'
    );
  });
});
