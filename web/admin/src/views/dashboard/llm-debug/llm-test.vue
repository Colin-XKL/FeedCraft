<template>
  <div class="py-8 px-16">
    <x-header title="LLM Debug" description="大模型测试"></x-header>

    <a-space direction="vertical" class="w-full">
      <a-card class="my-2 w-full" title="输入参数">
        <a-space direction="horizontal" class="w-full">
          <a-input
            v-model:model-value="model"
            placeholder="model name"
            class="max-w-48"
          ></a-input>
          <a-textarea
            v-model:model-value="prompt"
            placeholder="prompt input"
            allow-clear
            class="w-full min-w-96"
            auto-size
          ></a-textarea>
          <a-button :loading="isLoading" type="primary" @click="onSubmit"
            >Submit
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

  const model = ref('chatgpt-3.5');
  const prompt = ref('what is rss and how to use it?');
  const response = ref('');
  const isLoading = ref(false);

  async function onSubmit() {
    isLoading.value = true;
    try {
      const reqBody = {
        model: model.value,
        input: prompt.value,
      };
      const apiPath = '/api/admin/craft-debug/common-llm-call-test';
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
