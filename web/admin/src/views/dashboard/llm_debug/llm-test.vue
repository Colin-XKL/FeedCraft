<template>
  <div class="py-8 px-16">
    <x-header
      :title="t('llmDebug.llmTest.title')"
      :description="t('llmDebug.llmTest.description')"
    ></x-header>

    <a-space direction="vertical" class="w-full">
      <a-card class="my-2 w-full" :title="t('llmDebug.llmTest.inputParams')">
        <a-space direction="horizontal" class="w-full">
          <a-input
            v-model:model-value="model"
            :placeholder="t('llmDebug.llmTest.modelNamePlaceholder')"
            class="max-w-48"
          ></a-input>
          <a-textarea
            v-model:model-value="prompt"
            :placeholder="t('llmDebug.llmTest.promptPlaceholder')"
            allow-clear
            class="w-full min-w-96"
            auto-size
          ></a-textarea>
          <a-button :loading="isLoading" type="primary" @click="onSubmit"
            >{{ t('llmDebug.llmTest.submit') }}
          </a-button>
        </a-space>
      </a-card>
      <a-card
        :title="t('llmDebug.llmTest.response')"
        :loading="isLoading"
        class="w-full"
      >
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
  import { useI18n } from 'vue-i18n';

  const { t } = useI18n();

  const model = ref('');
  const prompt = ref('');
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
      Message.warning(error?.toString() ?? t('llmDebug.llmTest.unknownError'));
      console.error(error);
    } finally {
      isLoading.value = false;
    }
  }
</script>
