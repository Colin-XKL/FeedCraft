export const ALL_CRAFT_LIST_SCROLL_X = 900;

export function getAllCraftListColumns(t: (key: string) => string) {
  return [
    {
      title: t('allCraftList.table.name'),
      dataIndex: 'name',
      width: 280,
      ellipsis: true,
      tooltip: true,
    },
    {
      title: t('allCraftList.table.type'),
      dataIndex: 'type',
      width: 140,
      ellipsis: true,
      tooltip: true,
    },
    {
      title: t('allCraftList.table.templateOnly'),
      dataIndex: 'template_only',
      width: 120,
    },
    {
      title: t('allCraftList.table.description'),
      dataIndex: 'description',
    },
  ];
}
