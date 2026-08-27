import { describe, expect, it } from 'vitest';
import {
  ALL_CRAFT_LIST_SCROLL_X,
  getAllCraftListColumns,
} from './all_craft_list.columns';

const t = (key: string) => key;

describe('all craft list table columns', () => {
  it('reserves enough width so long English craft names stay on one line', () => {
    const name = getAllCraftListColumns(t).find(
      (column) => column.dataIndex === 'name'
    );

    expect(name?.width).toBeGreaterThanOrEqual(240);
    expect(name?.ellipsis).toBe(true);
  });

  it('keeps type identifiers like @sys/atom from wrapping mid-token', () => {
    const type = getAllCraftListColumns(t).find(
      (column) => column.dataIndex === 'type'
    );

    expect(type?.width).toBeGreaterThanOrEqual(120);
    expect(type?.ellipsis).toBe(true);
  });

  it('enables horizontal scroll so fixed identifier columns are not crushed', () => {
    expect(ALL_CRAFT_LIST_SCROLL_X).toBeGreaterThanOrEqual(800);
  });
});
