<template>
  <div class="py-8 px-16">
    <x-header
      title="LLM 调试"
      description="大模型API测试, 你可以快捷调试后端的LLM API配置是否正常工作"
    ></x-header>

    <a-space direction="vertical" class="w-full">
      <a-card class="my-2 w-full" title="输入参数">
        <a-space direction="horizontal" class="w-full">
          <a-input
            v-model:model-value="model"
            placeholder="模型名称"
            class="max-w-48"
          ></a-input>
          <a-textarea
            v-model:model-value="prompt"
            placeholder="在此输入你的提示词"
            allow-clear
            class="w-full min-w-96"
            auto-size
          ></a-textarea>
          <a-button :loading="isLoading" type="primary" @click="onSubmit"
            >提交
          </a-button>
        </a-space>
      </a-card>
      <a-card title="response" :loading="isLoading" class="w-full">
        <div v-if="response.length > 0">
          {{ response }}
        </div>
        <a-empty v-else></a-empty>
      </a-card>
    </a-space>
  </div>
</template>

<script lang="ts" setup>
  import { ref } from 'vue';
  import { Message } from '@arco-design/web-vue';
  import axios from 'axios';
  import XHeader from '@/components/header/x-header.vue';

  const model = ref('');
  const prompt = ref('what is rss and how to use it?');
  const response = ref('');
  const isLoading = ref(false);
  const baseUrl = import.meta.env.VITE_API_BASE_URL ?? '';

  async function onSubmit() {
    isLoading.value = true;
    try {
      const reqBody = {
        model: model.value,
        input: prompt.value,
      };
      const apiPath = `${baseUrl}/api/admin/craft-debug/common-llm-call-test`;
      const resp = await axios.post(apiPath, reqBody);
      response.value = resp.data.output;
    } catch (error) {
      Message.warning(error?.toString() ?? 'Unknown Error');
      console.error(error);
    } finally {
      isLoading.value = false;
    }
  }
</script>
