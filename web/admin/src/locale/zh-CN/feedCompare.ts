export default {
  feedCompare: {
    description:
      '指定feed源url与要应用的craft(craft flow 或 atom都可), 对比处理前和处理后的结果',
    inputLink: '输入链接',
    placeholder: '输入 RSS feed URL',
    compare: '对比',
    originalFeed: '原始 Feed',
    craftAppliedFeed: 'Craft 处理后的 Feed',
    selectCraftFlow: {
      placeholder: '选择 CraftFlow',
      tabs: {
        system: '系统 Craft Templates',
        user: '用户 Craft Atoms',
        flow: 'Craft Flows',
      },
    },
    message: {
      inputRequired: '请输入 feed URL 并选择一个 craft',
      unknownError: '未知错误',
    },
  },
};
