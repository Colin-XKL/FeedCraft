import { describe, expect, it } from 'vitest';
import {
  CUSTOM_RECIPE_ACTIONS_CELL_CLASS,
  CUSTOM_RECIPE_SCROLL_X,
  getCustomRecipeColumns,
} from './custom_recipe.columns';

const t = (key: string) => key;

describe('custom recipe table columns', () => {
  const columns = getCustomRecipeColumns(t);
  const actions = columns.find((column) => column.slotName === 'actions');

  it('pins the action column and reserves room for Edit/Preview/Copy Link/Delete', () => {
    expect(actions?.fixed).toBe('right');
    expect(actions?.width).toBeGreaterThanOrEqual(360);
    expect(actions?.bodyCellClass).toBe(CUSTOM_RECIPE_ACTIONS_CELL_CLASS);
  });

  it('gives every column an explicit width so Arco enables horizontal scroll and sticky actions', () => {
    const widths = columns.map((column) => column.width ?? 0);
    expect(widths.every((width) => width > 0)).toBe(true);

    const totalWidth = widths.reduce((sum, width) => sum + width, 0);
    expect(CUSTOM_RECIPE_SCROLL_X).toBeGreaterThanOrEqual(totalWidth);
  });
});
