<template>
  <div
    class="w-full min-h-[calc(100vh-60px)] flex justify-center items-start pt-20 bg-gray-50"
  >
    <div class="w-full max-w-3xl px-4 pb-20">
      <a-card class="rounded-3xl shadow-lg p-2" :bordered="false">
        <a-tabs type="capsule" size="large" animation>
          <!-- Tab 1: Generate -->
          <a-tab-pane key="generate" :title="t('urlGenerator.tabGenerate')">
            <div class="px-6 py-4 space-y-8">
              <h1 class="text-3xl font-bold text-center text-gray-900">{{
                t('urlGenerator.title')
              }}</h1>

              <div class="space-y-6">
                <div>
                  <label class="block text-sm font-bold text-gray-900 mb-2">{{
                    t('urlGenerator.selectCraft')
                  }}</label>
                  <CraftFlowSelect
                    v-model="selectedCraft"
                    mode="single"
                    class="w-full"
                  />
                </div>

                <div>
                  <label class="block text-sm font-bold text-gray-900 mb-2">{{
                    t('urlGenerator.inputOriginalUrl')
                  }}</label>
                  <a-input
                    v-model="inputUrl"
                    :placeholder="t('urlGenerator.inputUrlPlaceholder')"
                    class="!rounded-lg !bg-gray-50 !border-gray-200 !py-2.5 focus:!bg-white focus:!border-emerald-500"
                    size="large"
                  />
                </div>

                <a-button
                  type="primary"
                  size="large"
                  long
                  class="!rounded-lg !h-12 !text-base !font-bold"
                  :style="{ backgroundColor: 'rgb(16, 185, 129)' }"
                  @click="generateUrl"
                >
                  {{ t('urlGenerator.showCraftedUrl') }}
                </a-button>

                <!-- Result Area -->
                <div v-if="resultUrl" class="mt-8">
                  <div class="relative group">
                    <div
                      class="bg-gray-50 border border-gray-200 rounded-xl p-4 pr-14 break-all text-gray-600 font-mono text-sm leading-relaxed"
                    >
                      {{ resultUrl }}
                    </div>
                    <div class="absolute right-2 top-2 h-full">
                      <a-button
                        type="text"
                        shape="square"
                        size="large"
                        @click="copyUrl"
                      >
                        <template #icon>
                          <icon-check
                            v-if="copied"
                            class="text-green-500 text-lg"
                          />
                          <icon-copy v-else class="text-lg text-gray-500" />
                        </template>
                      </a-button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </a-tab-pane>

          <!-- Tab 2: Parse -->
          <a-tab-pane key="parse" :title="t('urlGenerator.tabParse')">
            <div class="px-6 py-4 space-y-8">
              <h1 class="text-3xl font-bold text-center text-gray-900">{{
                t('urlGenerator.tabParse')
              }}</h1>

              <div class="space-y-6">
                <div>
                  <label class="block text-sm font-bold text-gray-900 mb-2">{{
                    t('urlGenerator.inputFeedCraftUrl')
                  }}</label>
                  <a-input
                    v-model="parseInputUrl"
                    :placeholder="
                      t('urlGenerator.inputFeedCraftUrlPlaceholder')
                    "
                    allow-clear
                    class="!rounded-lg !bg-gray-50 !border-gray-200 !py-2.5 focus:!bg-white focus:!border-emerald-500"
                    size="large"
                  />
                </div>

                <a-button
                  type="primary"
                  size="large"
                  long
                  class="!rounded-lg !h-12 !text-base !font-bold"
                  :style="{ backgroundColor: 'rgb(16, 185, 129)' }"
                  @click="parseUrl"
                >
                  <template #icon>
                    <icon-thunderbolt />
                  </template>
                  {{ t('urlGenerator.parseButton') }}
                </a-button>

                <div
                  v-if="parsedResult"
                  class="mt-8 p-6 bg-gray-50 rounded-xl border border-gray-100"
                >
                  <!-- Display Craft -->
                  <div class="mb-5">
                    <span
                      class="text-xs font-bold text-gray-500 uppercase tracking-wider block mb-2"
                    >
                      {{ t('urlGenerator.parsedCraft') }}
                    </span>
                    <div class="flex items-center">
                      <a-tag color="arcoblue" size="large">
                        {{ parsedResult.craft }}
                      </a-tag>
                    </div>
                  </div>
                  <!-- Display Source URL -->
                  <div class="mb-5">
                    <span
                      class="text-xs font-bold text-gray-500 uppercase tracking-wider block mb-2"
                    >
                      {{ t('urlGenerator.parsedSourceUrl') }}
                    </span>
                    <div class="bg-white p-3 rounded-lg border border-gray-200">
                      <a
                        :href="parsedResult.sourceUrl"
                        target="_blank"
                        class="text-emerald-600 hover:text-emerald-700 break-all hover:underline"
                      >
                        {{ parsedResult.sourceUrl }}
                      </a>
                    </div>
                  </div>
                  <!-- Other Params -->
                  <div v-if="Object.keys(parsedResult.params).length > 0">
                    <span
                      class="text-xs font-bold text-gray-500 uppercase tracking-wider block mb-2"
                    >
                      {{ t('urlGenerator.parsedParams') }}
                    </span>
                    <div
                      class="bg-white rounded-lg border border-gray-200 overflow-hidden"
                    >
                      <div
                        v-for="(val, key) in parsedResult.params"
                        :key="key"
                        class="flex border-b border-gray-100 last:border-0 px-4 py-3 text-sm"
                      >
                        <span
                          class="text-gray-500 w-1/3 shrink-0 font-medium"
                          >{{ key }}</span
                        >
                        <span class="text-gray-900 font-mono break-all">{{
                          val
                        }}</span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </a-tab-pane>
        </a-tabs>
      </a-card>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref } from 'vue';
  import CraftFlowSelect from '@/views/dashboard/craft_flow/CraftFlowSelect.vue';
  import { useI18n } from 'vue-i18n';
  import { Message } from '@arco-design/web-vue';
  import {
    IconCopy,
    IconCheck,
    IconThunderbolt,
  } from '@arco-design/web-vue/es/icon';

  const { t } = useI18n();

  // Mode 1: Generate
  const selectedCraft = ref<string[]>([]);
  const customCraft = ref('');
  const inputUrl = ref('');
  const resultUrl = ref('');
  const copied = ref(false);

  const generateUrl = () => {
    const currentSelectedCraft = customCraft.value
      ? customCraft.value
      : selectedCraft.value;
    const baseUrl = import.meta.env.VITE_API_BASE_URL ?? window.location.origin;
    resultUrl.value = `${baseUrl}/craft/${currentSelectedCraft}?input_url=${encodeURIComponent(
      inputUrl.value,
    )}`;
    copied.value = false;
  };

  const copyUrl = () => {
    if (resultUrl.value) {
      navigator.clipboard
        .writeText(resultUrl.value)
        .then(() => {
          copied.value = true;
          setTimeout(() => {
            copied.value = false;
          }, 2000);
          Message.success(t('urlGenerator.copied'));
        })
        .catch((err) => {
          console.error('Failed to copy: ', err);
          Message.error('Failed to copy');
        });
    }
  };

  // Mode 2: Parse
  interface ParsedData {
    craft: string;
    sourceUrl: string;
    params: Record<string, string>;
  }

  const parseInputUrl = ref('');
  const parsedResult = ref<ParsedData | null>(null);

  const parseUrl = () => {
    if (!parseInputUrl.value) {
      return;
    }
    try {
      const urlObj = new URL(parseInputUrl.value);
      // Path structure: .../craft/{craftName}
      const pathParts = urlObj.pathname.split('/');
      const craftIndex = pathParts.indexOf('craft');
      let craft = '';
      if (craftIndex !== -1 && craftIndex + 1 < pathParts.length) {
        craft = pathParts[craftIndex + 1];
      }

      const sourceUrl = urlObj.searchParams.get('input_url') || '';
      const params: Record<string, string> = {};

      urlObj.searchParams.forEach((value, key) => {
        if (key !== 'input_url') {
          params[key] = value;
        }
      });

      if (!craft && !sourceUrl) {
        Message.error(t('urlGenerator.parseError'));
        return;
      }

      parsedResult.value = {
        craft,
        sourceUrl,
        params,
      };
    } catch (e) {
      Message.error(t('urlGenerator.parseError'));
      parsedResult.value = null;
    }
  };
</script>

<style scoped>
  /* Override styles for deep customization if needed */
  :deep(.arco-tabs-nav-tab) {
    justify-content: center;
  }
</style>
