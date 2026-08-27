export const ALL_CRAFT_LIST_SCROLL_X = 1000;
export const ALL_CRAFT_ID_CELL_CLASS = 'all-craft-id-cell';

export function getAllCraftListColumns(t: (key: string) => string) {
  return [
    {
      title: t('allCraftList.table.name'),
      dataIndex: 'name',
      width: 360,
      ellipsis: true,
      tooltip: true,
      bodyCellClass: ALL_CRAFT_ID_CELL_CLASS,
    },
    {
      title: t('allCraftList.table.type'),
      dataIndex: 'type',
      width: 140,
      ellipsis: true,
      tooltip: true,
      bodyCellClass: ALL_CRAFT_ID_CELL_CLASS,
    },
    {
      title: t('allCraftList.table.templateOnly'),
      dataIndex: 'template_only',
      width: 120,
    },
    {
      title: t('allCraftList.table.description'),
      dataIndex: 'description',
      bodyCellClass: 'all-craft-desc-cell',
    },
  ];
}
