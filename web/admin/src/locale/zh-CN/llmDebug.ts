export default {
  llmDebug: {
    adCheck: {
      title: 'LLM 调试页面 - 广告检测',
      description:
        '利用大模型能力检测文章是否为广告或软文。对应 craft: ignore-advertorial',
      inputLink: '输入链接',
      inputTip: '输入文章链接进行预览',
      placeholder: '输入文章链接',
      enhanceMode: '开启增强模式获取文章内容',
      submit: '提交',
      resultPreview: '结果预览',
      isAdvertorial: '是否为广告：',
      articleContent: '文章内容：',
      message: {
        unknownError: '未知错误',
      },
    },
    llmTest: {
      title: 'LLM 调试',
      description: '大模型API测试, 你可以快捷调试后端的LLM API配置是否正常工作',
      inputParams: '输入参数',
      modelNamePlaceholder: '模型名称',
      promptPlaceholder: '在此输入你的提示词',
      submit: '提交',
      response: '响应',
      unknownError: '未知错误',
    },
  },
};
