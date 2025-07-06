export default {
  llmDebug: {
    adCheck: {
      title: 'LLM Debug Page - Check for Advertisements',
      description:
        'Using large model capabilities to detect whether an article is an advertisement or a promotional soft article. Corresponding craft: ignore-advertorial',
      inputLink: 'Input Link',
      inputTip: 'Enter the article link to preview',
      placeholder: 'Enter article link',
      enhanceMode: 'Enable enhanced mode to get article content',
      submit: 'Submit',
      resultPreview: 'Result Preview',
      isAdvertorial: 'Is Advertorial:',
      articleContent: 'Article Content:',
      message: {
        unknownError: 'Unknown Error',
      },
    },
    llmTest: {
      title: 'LLM Debug',
      description:
        'LLM API test, you can quickly debug whether the backend LLM API configuration is working properly',
      inputParams: 'Input Parameters',
      modelNamePlaceholder: 'Model Name',
      promptPlaceholder: 'Enter your prompt here',
      submit: 'Submit',
      response: 'Response',
      unknownError: 'Unknown Error',
    },
  },
};
