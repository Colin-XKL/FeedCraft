<template>
  <div class="py-8 px-16">
    <x-header
      :title="t('llmDebug.llmFilter.title')"
      :description="t('llmDebug.llmFilter.description')"
    ></x-header>

    <a-row :gutter="24">
      <a-col :span="12">
        <a-card class="my-2" :title="t('llmDebug.llmFilter.inputLink')">
          <a-space direction="vertical" fill>
            <a-input
              v-model="feedUrl"
              type="text"
              class="w-full"
              :placeholder="t('llmDebug.llmFilter.placeholder')"
            />

            <a-textarea
              v-model="filterCondition"
              :placeholder="t('llmDebug.llmFilter.conditionPlaceholder')"
              :auto-size="{ minRows: 3, maxRows: 6 }"
            />

            <a-space>
              <a-checkbox v-model="enhanceMode">{{
                t('llmDebug.llmFilter.enhanceMode')
              }}</a-checkbox>
            </a-space>

            <a-button
              :loading="isLoading"
              type="primary"
              @click="testLLMFilter"
            >
              {{ t('llmDebug.llmFilter.submit') }}
            </a-button>
          </a-space>
        </a-card>
      </a-col>
      <a-col :span="12">
        <a-card
          class="my-2"
          :title="t('llmDebug.llmFilter.resultPreview')"
          :loading="isLoading"
        >
          <a-space v-if="testResult" direction="vertical" fill>
            <div>
              <span class="font-bold mr-2">{{
                t('llmDebug.llmFilter.isFiltered')
              }}</span>
              <a-tag :color="testResult.is_filtered ? 'red' : 'green'">
                {{ testResult.is_filtered ? 'Yes (Filtered)' : 'No (Kept)' }}
              </a-tag>
            </div>
            <div>
              <div class="font-bold mb-2">{{
                t('llmDebug.llmFilter.articleContent')
              }}</div>
              <div
                class="bg-gray-50 p-4 rounded max-h-[500px] overflow-y-auto whitespace-pre-wrap"
              >
                {{ testResult.article_content }}
              </div>
            </div>
          </a-space>
          <a-empty v-else />
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script lang="ts" setup>
  import { ref } from 'vue';
  import { Message } from '@arco-design/web-vue';
  import XHeader from '@/components/header/x-header.vue';
  import { useI18n } from 'vue-i18n';
  import axios from 'axios';
  import { getToken } from '@/utils/auth';

  const { t } = useI18n();

  const feedUrl = ref('');
  const filterCondition = ref('Is this content spam or low quality?');
  const enhanceMode = ref(false);
  const isLoading = ref(false);
  const testResult = ref<any>(null);

  const baseUrl = import.meta.env.VITE_API_BASE_URL ?? '';

  async function testLLMFilter() {
    if (!feedUrl.value || !filterCondition.value) {
      Message.warning(t('llmDebug.llmFilter.message.inputRequired'));
      return;
    }

    isLoading.value = true;
    testResult.value = null;
    try {
      const response = await axios.post(
        `${baseUrl}/api/admin/craft-debug/llm-filter`,
        {
          url: feedUrl.value,
          enhance_mode: enhanceMode.value,
          filter_condition: filterCondition.value,
        },
        {
          headers: {
            Authorization: `Bearer ${getToken()}`,
          },
        }
      );
      if (response.data.code === 200) {
        testResult.value = response.data.data;
      } else {
        Message.error(
          response.data.msg || t('llmDebug.llmFilter.message.unknownError')
        );
      }
    } catch (error: any) {
      Message.error(
        error?.response?.data?.msg ||
          error?.message ||
          t('llmDebug.llmFilter.message.unknownError')
      );
    } finally {
      isLoading.value = false;
    }
  }
</script>

<script lang="ts">
  export default {
    name: 'LlmFilterDebug',
  };
</script>
