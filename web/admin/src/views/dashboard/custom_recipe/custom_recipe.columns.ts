export const CUSTOM_RECIPE_SCROLL_X = 1410;
export const CUSTOM_RECIPE_ACTIONS_CELL_CLASS = 'custom-recipe-actions-cell';

export function getCustomRecipeColumns(t: (key: string) => string) {
  return [
    {
      title: t('customRecipe.form.name'),
      dataIndex: 'id',
      width: 200,
      ellipsis: true,
      tooltip: true,
    },
    {
      title: t('customRecipe.form.description'),
      dataIndex: 'description',
      width: 180,
      ellipsis: true,
      tooltip: true,
    },
    {
      title: t('customRecipe.form.craft'),
      dataIndex: 'craft',
      slotName: 'craft',
      width: 220,
      ellipsis: true,
      tooltip: true,
    },
    {
      title: t('customRecipe.status.active'),
      slotName: 'status',
      width: 100,
    },
    {
      title: t('customRecipe.form.sourceType'),
      dataIndex: 'source_type',
      width: 110,
    },
    {
      title: t('customRecipe.form.sourceConfig'),
      dataIndex: 'source_config',
      slotName: 'source_config',
      width: 240,
    },
    {
      title: t('customRecipe.actions'),
      slotName: 'actions',
      width: 360,
      fixed: 'right' as const,
      bodyCellClass: CUSTOM_RECIPE_ACTIONS_CELL_CLASS,
    },
  ];
}
