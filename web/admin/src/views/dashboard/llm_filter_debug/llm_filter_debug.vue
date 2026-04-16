<template>
  <div class="py-8 px-16">
    <x-header
      :title="t('llmDebug.llmFilter.title')"
      :description="t('llmDebug.llmFilter.description')"
    ></x-header>

    <a-tabs type="rounded" class="mt-4">
      <!-- Feed Mode Tab -->
      <a-tab-pane key="feed" :title="t('llmDebug.llmFilter.modeFeed')">
        <a-card class="my-2" :title="t('llmDebug.llmFilter.inputFeed')">
          <a-space direction="vertical" fill>
            <a-input
              v-model="feedModeUrl"
              type="text"
              class="w-full"
              :placeholder="t('llmDebug.llmFilter.feedPlaceholder')"
            />

            <a-textarea
              v-model="filterCondition"
              :placeholder="t('llmDebug.llmFilter.conditionPlaceholder')"
              :auto-size="{ minRows: 3, maxRows: 6 }"
            />

            <a-button
              :loading="isFeedLoading"
              type="primary"
              @click="testFeedFilter"
            >
              {{ t('llmDebug.llmFilter.submitFeed') }}
            </a-button>
          </a-space>
        </a-card>

        <a-card
          class="my-2"
          :title="t('llmDebug.llmFilter.resultPreview')"
          :loading="isFeedLoading"
        >
          <div v-if="feedResult && feedResult.length > 0">
            <a-tabs type="line">
              <a-tab-pane
                key="all"
                :title="`${t('llmDebug.llmFilter.tabAll')} (${
                  feedResult.length
                })`"
              >
                <a-list :data="feedResult" :bordered="false">
                  <template #item="{ item }">
                    <a-list-item class="mb-4 bg-gray-50 rounded">
                      <a-list-item-meta :title="item.title">
                        <template #description>
                          <a
                            :href="item.link"
                            target="_blank"
                            class="text-blue-500 hover:underline mb-2 block text-xs"
                            >{{ item.link }}</a
                          >
                          <div
                            class="line-clamp-3 text-sm text-gray-500 mb-2"
                            >{{ item.content }}</div
                          >
                          <a-space>
                            <a-tag :color="item.is_filtered ? 'red' : 'green'">
                              {{
                                item.is_filtered
                                  ? t('llmDebug.llmFilter.statusFiltered')
                                  : t('llmDebug.llmFilter.statusKept')
                              }}
                            </a-tag>
                            <a-button
                              size="mini"
                              type="outline"
                              @click="testSingleItem(item.content)"
                              >{{
                                t('llmDebug.llmFilter.testSingle')
                              }}</a-button
                            >
                          </a-space>
                        </template>
                      </a-list-item-meta>
                    </a-list-item>
                  </template>
                </a-list>
              </a-tab-pane>
              <a-tab-pane
                key="kept"
                :title="`${t('llmDebug.llmFilter.tabKept')} (${
                  keptItems.length
                })`"
              >
                <a-list :data="keptItems" :bordered="false">
                  <template #item="{ item }">
                    <a-list-item class="mb-4 bg-green-50 rounded">
                      <a-list-item-meta :title="item.title">
                        <template #description>
                          <a
                            :href="item.link"
                            target="_blank"
                            class="text-blue-500 hover:underline mb-2 block text-xs"
                            >{{ item.link }}</a
                          >
                          <div
                            class="line-clamp-3 text-sm text-gray-500 mb-2"
                            >{{ item.content }}</div
                          >
                        </template>
                      </a-list-item-meta>
                    </a-list-item>
                  </template>
                </a-list>
              </a-tab-pane>
              <a-tab-pane
                key="filtered"
                :title="`${t('llmDebug.llmFilter.tabFiltered')} (${
                  filteredItems.length
                })`"
              >
                <a-list :data="filteredItems" :bordered="false">
                  <template #item="{ item }">
                    <a-list-item class="mb-4 bg-red-50 rounded">
                      <a-list-item-meta :title="item.title">
                        <template #description>
                          <a
                            :href="item.link"
                            target="_blank"
                            class="text-blue-500 hover:underline mb-2 block text-xs"
                            >{{ item.link }}</a
                          >
                          <div
                            class="line-clamp-3 text-sm text-gray-500 mb-2"
                            >{{ item.content }}</div
                          >
                        </template>
                      </a-list-item-meta>
                    </a-list-item>
                  </template>
                </a-list>
              </a-tab-pane>
            </a-tabs>
          </div>
          <a-empty v-else />
        </a-card>
      </a-tab-pane>

      <!-- URL Mode Tab -->
      <a-tab-pane key="url" :title="t('llmDebug.llmFilter.modeUrl')">
        <a-row :gutter="24">
          <a-col :span="12">
            <a-card class="my-2" :title="t('llmDebug.llmFilter.inputLink')">
              <a-space direction="vertical" fill>
                <a-input
                  v-model="urlModeUrl"
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
                  :loading="isUrlLoading"
                  type="primary"
                  @click="testUrlFilter"
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
              :loading="isUrlLoading"
            >
              <a-space v-if="urlResult" direction="vertical" fill>
                <div>
                  <span class="font-bold mr-2">{{
                    t('llmDebug.llmFilter.isFiltered')
                  }}</span>
                  <a-tag :color="urlResult.is_filtered ? 'red' : 'green'">
                    {{
                      urlResult.is_filtered
                        ? t('llmDebug.llmFilter.statusFiltered')
                        : t('llmDebug.llmFilter.statusKept')
                    }}
                  </a-tag>
                </div>
                <div>
                  <div class="font-bold mb-2">{{
                    t('llmDebug.llmFilter.articleContent')
                  }}</div>
                  <div
                    class="bg-gray-50 p-4 rounded max-h-[500px] overflow-y-auto whitespace-pre-wrap"
                  >
                    {{ urlResult.article_content }}
                  </div>
                </div>
              </a-space>
              <a-empty v-else />
            </a-card>
          </a-col>
        </a-row>
      </a-tab-pane>
    </a-tabs>

    <!-- Single Text Test Modal -->
    <a-modal
      v-model:visible="isSingleModalVisible"
      :title="t('llmDebug.llmFilter.testSingleModalTitle')"
      :ok-loading="isSingleTesting"
      @ok="runSingleTextTest"
    >
      <a-space direction="vertical" fill>
        <div>{{ t('llmDebug.llmFilter.conditionPlaceholder') }}</div>
        <a-textarea
          v-model="filterCondition"
          :auto-size="{ minRows: 2, maxRows: 4 }"
        />
        <div>{{ t('llmDebug.llmFilter.articleContent') }}</div>
        <div
          class="bg-gray-50 p-2 rounded max-h-[200px] overflow-y-auto whitespace-pre-wrap text-sm border"
        >
          {{ singleTestContent }}
        </div>
        <div v-if="singleTestResult !== null">
          <span class="font-bold mr-2">{{
            t('llmDebug.llmFilter.isFiltered')
          }}</span>
          <a-tag :color="singleTestResult ? 'red' : 'green'">
            {{
              singleTestResult
                ? t('llmDebug.llmFilter.statusFiltered')
                : t('llmDebug.llmFilter.statusKept')
            }}
          </a-tag>
        </div>
      </a-space>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
  import { ref, computed } from 'vue';
  import { Message } from '@arco-design/web-vue';
  import XHeader from '@/components/header/x-header.vue';
  import { useI18n } from 'vue-i18n';
  import axios from 'axios';

  const { t } = useI18n();

  const filterCondition = ref('Is this content spam or low quality?');

  // --- URL Mode State ---
  const urlModeUrl = ref('');
  const enhanceMode = ref(false);
  const isUrlLoading = ref(false);
  const urlResult = ref<any>(null);

  // --- Feed Mode State ---
  const feedModeUrl = ref('');
  const isFeedLoading = ref(false);
  const feedResult = ref<any[]>([]);

  const keptItems = computed(() =>
    feedResult.value.filter((i) => !i.is_filtered)
  );
  const filteredItems = computed(() =>
    feedResult.value.filter((i) => i.is_filtered)
  );

  // --- Single Test State ---
  const isSingleModalVisible = ref(false);
  const singleTestContent = ref('');
  const singleTestResult = ref<boolean | null>(null);
  const isSingleTesting = ref(false);

  // --- Handlers ---
  async function testUrlFilter() {
    if (!urlModeUrl.value || !filterCondition.value) {
      Message.warning(t('llmDebug.llmFilter.message.inputRequired'));
      return;
    }

    isUrlLoading.value = true;
    urlResult.value = null;
    try {
      const response = await axios.post(
        `/api/admin/craft-debug/llm-filter/url`,
        {
          url: urlModeUrl.value,
          enhance_mode: enhanceMode.value,
          filter_condition: filterCondition.value,
        }
      );
      if (response.data.code === 0 || response.data.code === 200) {
        urlResult.value = response.data.data;
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
      isUrlLoading.value = false;
    }
  }

  async function testFeedFilter() {
    if (!feedModeUrl.value || !filterCondition.value) {
      Message.warning(t('llmDebug.llmFilter.message.inputRequired'));
      return;
    }

    isFeedLoading.value = true;
    feedResult.value = [];
    try {
      const response = await axios.post(
        `/api/admin/craft-debug/llm-filter/feed`,
        {
          feed_url: feedModeUrl.value,
          filter_condition: filterCondition.value,
        }
      );
      if (response.data.code === 0 || response.data.code === 200) {
        feedResult.value = response.data.data || [];
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
      isFeedLoading.value = false;
    }
  }

  function testSingleItem(content: string) {
    singleTestContent.value = content;
    singleTestResult.value = null;
    isSingleModalVisible.value = true;
  }

  async function runSingleTextTest() {
    if (!singleTestContent.value || !filterCondition.value) {
      Message.warning(t('llmDebug.llmFilter.message.inputRequired'));
      return false; // prevent modal close
    }

    isSingleTesting.value = true;
    try {
      const response = await axios.post(
        `/api/admin/craft-debug/llm-filter/text`,
        {
          text: singleTestContent.value,
          filter_condition: filterCondition.value,
        }
      );
      if (response.data.code === 0 || response.data.code === 200) {
        singleTestResult.value = response.data.data.is_filtered;
        return false; // prevent modal close so user can see result
      }
      Message.error(
        response.data.msg || t('llmDebug.llmFilter.message.unknownError')
      );
    } catch (error: any) {
      Message.error(
        error?.response?.data?.msg ||
          error?.message ||
          t('llmDebug.llmFilter.message.unknownError')
      );
    } finally {
      isSingleTesting.value = false;
    }
    return true;
  }
</script>

<script lang="ts">
  export default {
    name: 'LlmFilterDebug',
  };
</script>
