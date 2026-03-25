<template>
  <div class="py-8 px-16">
    <x-header
      :title="t('llmDebug.adCheck.title')"
      :description="t('llmDebug.adCheck.description')"
    >
    </x-header>

    <a-card class="my-2" :title="t('llmDebug.adCheck.inputLink')">
      <p>{{ t('llmDebug.adCheck.inputTip') }}</p>
      <a-space>
        <a-input
          v-model="articleUrl"
          type="text"
          :placeholder="t('llmDebug.adCheck.placeholder')"
          class="w-full"
        />
        <a-checkbox v-model="useEnhanceMode">{{
          t('llmDebug.adCheck.enhanceMode')
        }}</a-checkbox>
        <a-button :loading="isLoading" @click="fetchArticle">{{
          t('llmDebug.adCheck.submit')
        }}</a-button>
      </a-space>
    </a-card>
    <a-card
      :title="t('llmDebug.adCheck.resultPreview')"
      class="my-4"
      :loading="isLoading"
    >
      <div v-if="articleContent?.length !== 0">
        <p class="mb-2"
          ><b>{{ t('llmDebug.adCheck.isAdvertorial') }}</b
          >{{ isAdvertorial }}</p
        >
        <hr class="my-1" />
        <p class="my-1"
          ><b>{{ t('llmDebug.adCheck.articleContent') }}</b></p
        >
        <p>{{ articleContent }}</p>
      </div>
      <a-empty v-else />
    </a-card>
  </div>
</template>

<script lang="ts" setup>
  import { ref } from 'vue';
  import { Message } from '@arco-design/web-vue';

  import axios from 'axios';
  import XHeader from '@/components/header/x-header.vue';
  import { useI18n } from 'vue-i18n';

  const { t } = useI18n();

  const articleUrl = ref('');
  const useEnhanceMode = ref(false);
  const articleContent = ref('');
  const isAdvertorial = ref(false);
  const isLoading = ref(false);
  const baseUrl = import.meta.env.VITE_API_BASE_URL ?? '';

  async function fetchArticle() {
    isLoading.value = true;
    try {
      const resp = await axios.post(
        `${baseUrl}/api/admin/craft-debug/advertorial`,
        {
          url: articleUrl.value,
          enhance_mode: useEnhanceMode.value,
        },
      );
      articleContent.value = resp.data.article_content;
      isAdvertorial.value = resp.data.is_advertorial;
    } catch (error) {
      Message.warning(
        error?.toString() ?? t('llmDebug.adCheck.message.unknownError'),
      );
      console.error(error);
    } finally {
      isLoading.value = false;
    }
  }
</script>

<script lang="ts">
  export default {
    name: 'LLMDebug',
  };
</script>
