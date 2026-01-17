<template>
  <div
    class="w-full min-h-[calc(100vh-60px)] flex justify-center items-start pt-20 bg-gray-50"
  >
    <div class="w-full max-w-3xl px-4 pb-20">
      <!-- Custom Tabs -->
      <div class="flex justify-center gap-4 mb-8">
        <button
          class="px-8 py-3 rounded-lg text-base font-bold transition-all duration-200 flex items-center gap-2 shadow-sm"
          :class="
            activeTab === 'generate'
              ? 'bg-emerald-500 text-white hover:bg-emerald-600 shadow-emerald-200'
              : 'bg-white text-gray-600 hover:bg-gray-50 border border-gray-200'
          "
          @click="activeTab = 'generate'"
        >
          <icon-link class="text-lg" />
          {{ t('urlGenerator.tabGenerate') }}
        </button>
        <button
          class="px-8 py-3 rounded-lg text-base font-bold transition-all duration-200 flex items-center gap-2 shadow-sm"
          :class="
            activeTab === 'parse'
              ? 'bg-emerald-500 text-white hover:bg-emerald-600 shadow-emerald-200'
              : 'bg-white text-gray-600 hover:bg-gray-50 border border-gray-200'
          "
          @click="activeTab = 'parse'"
        >
          <icon-thunderbolt class="text-lg" />
          {{ t('urlGenerator.tabParse') }}
        </button>
      </div>

      <!-- Main Card -->
      <div
        class="bg-white rounded-3xl shadow-lg p-10 transition-all duration-300"
      >
        <!-- Tab 1: Generate -->
        <div v-if="activeTab === 'generate'" class="space-y-8">
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

            <button
              class="w-full bg-emerald-500 hover:bg-emerald-600 text-white font-bold py-3.5 px-4 rounded-lg transition-colors duration-200 shadow-md hover:shadow-lg active:scale-[0.99]"
              @click="generateUrl"
            >
              {{ t('urlGenerator.showCraftedUrl') }}
            </button>

            <!-- Result Area -->
            <div v-if="resultUrl" class="mt-8">
              <div class="relative group">
                <div
                  class="bg-gray-50 border border-gray-200 rounded-xl p-4 pr-14 break-all text-gray-600 font-mono text-sm leading-relaxed"
                >
                  {{ resultUrl }}
                </div>
                <div class="absolute right-2 top-2 h-full">
                  <button
                    class="p-2 text-gray-500 hover:text-gray-700 bg-white border border-gray-200 rounded-lg shadow-sm transition-all hover:scale-105 active:scale-95"
                    :title="t('urlGenerator.copyUrl')"
                    @click="copyUrl"
                  >
                    <icon-check v-if="copied" class="text-green-500 text-lg" />
                    <icon-copy v-else class="text-lg" />
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Tab 2: Parse -->
        <div v-if="activeTab === 'parse'" class="space-y-8">
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
                :placeholder="t('urlGenerator.inputFeedCraftUrlPlaceholder')"
                allow-clear
                class="!rounded-lg !bg-gray-50 !border-gray-200 !py-2.5 focus:!bg-white focus:!border-emerald-500"
                size="large"
              />
            </div>

            <button
              class="w-full bg-emerald-500 hover:bg-emerald-600 text-white font-bold py-3.5 px-4 rounded-lg transition-colors duration-200 shadow-md hover:shadow-lg active:scale-[0.99]"
              @click="parseUrl"
            >
              <icon-thunderbolt class="mr-2" />
              {{ t('urlGenerator.parseButton') }}
            </button>

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
                  <span
                    class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-emerald-100 text-emerald-800"
                  >
                    {{ parsedResult.craft }}
                  </span>
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
                    <span class="text-gray-500 w-1/3 shrink-0 font-medium">{{
                      key
                    }}</span>
                    <span class="text-gray-900 font-mono break-all">{{
                      val
                    }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
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
    IconLink,
    IconThunderbolt,
  } from '@arco-design/web-vue/es/icon';

  const { t } = useI18n();

  // Tabs
  const activeTab = ref<'generate' | 'parse'>('generate');

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
  /* Scoped styles can be used if Tailwind is not enough */
</style>
