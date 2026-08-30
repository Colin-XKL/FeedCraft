<template>
  <div class="container p-4 w-full">
    <div>
      <div class="mb-8 text-2xl">
        <p class="text-gray-700"
          >{{ t('urlGenerator.welcome') }}<br />
          <span
            class="text-4xl font-bold text-sky-700 underline decoration-sky-500 decoration-wavy hover:underline-offset-2 hover:decoration-4"
            >FeedCraft</span
          ></p
        >
      </div>

      <a-tabs type="card-gutter" size="large">
        <!-- Tab 1: Generate -->
        <a-tab-pane key="generate" :title="t('urlGenerator.tabGenerate')">
          <a-card class="rounded-lg w-full min-w-8xl">
            <h1 class="text-2xl mb-2">{{ t('urlGenerator.title') }}</h1>
            <p class="mb-8">
              {{ t('urlGenerator.description') }}
            </p>
            <div class="mb-4">
              <label for="siteSelector" class="mr-2">{{
                t('urlGenerator.selectCraft')
              }}</label>
              <CraftFlowSelect v-model="selectedCraft" mode="single" />
            </div>
            <div class="mb-4">
              <label for="inputUrl" class="mr-2">{{
                t('urlGenerator.inputOriginalUrl')
              }}</label>
              <a-input
                id="inputUrl"
                v-model="inputUrl"
                type="text"
                :placeholder="t('urlGenerator.inputUrlPlaceholder')"
              />
            </div>
            <a-button type="primary" @click="generateUrl">
              {{
                resultUrl
                  ? t('urlGenerator.regenerateUrl')
                  : t('urlGenerator.showCraftedUrl')
              }}
            </a-button>
            <div
              class="result-panel mt-8"
              :class="{
                'result-panel--fresh': resultUrl && !isResultStale,
                'result-panel--stale': isResultStale,
              }"
            >
              <div class="flex items-center justify-between gap-3 mb-3">
                <label for="resultUrl" class="font-medium text-gray-700">{{
                  t('urlGenerator.resultUrl')
                }}</label>
                <a-tag
                  v-if="resultUrl"
                  :color="isResultStale ? 'orange' : 'green'"
                >
                  {{
                    isResultStale
                      ? t('urlGenerator.staleResult')
                      : t('urlGenerator.freshResult')
                  }}
                </a-tag>
              </div>
              <div id="resultUrl" class="result-url">
                {{ resultUrl || t('urlGenerator.noResultPlaceholder') }}
              </div>
              <div class="mt-3 flex flex-wrap items-center gap-3">
                <a-button
                  id="copyButton"
                  type="primary"
                  :disabled="!canCopyResult"
                  @click="copyUrl"
                >
                  {{
                    isResultStale
                      ? t('urlGenerator.regenerateToCopy')
                      : copyButtonText
                  }}
                </a-button>
                <span v-if="isResultStale" class="text-sm text-orange-600">
                  {{ t('urlGenerator.staleResultTip') }}
                </span>
                <span v-else-if="resultUrl" class="text-sm text-green-600">
                  {{ t('urlGenerator.freshResultTip') }}
                </span>
              </div>
            </div>
          </a-card>
        </a-tab-pane>

        <!-- Tab 2: Parse -->
        <a-tab-pane key="parse" :title="t('urlGenerator.tabParse')">
          <a-card class="rounded-lg w-full min-w-8xl">
            <div class="mb-4">
              <label class="mr-2">{{
                t('urlGenerator.inputFeedCraftUrl')
              }}</label>
              <a-input
                v-model="parseInputUrl"
                :placeholder="t('urlGenerator.inputFeedCraftUrlPlaceholder')"
                allow-clear
              />
            </div>
            <a-button type="primary" @click="parseUrl">{{
              t('urlGenerator.parseButton')
            }}</a-button>

            <div v-if="parsedResult" class="mt-8">
              <!-- Display Craft -->
              <div class="mb-4">
                <span class="font-bold block mb-2"
                  >{{ t('urlGenerator.parsedCraft') }}:</span
                >
                <a-tag color="blue" size="large">{{
                  parsedResult.craft
                }}</a-tag>
              </div>
              <!-- Display Source URL -->
              <div class="mb-4">
                <span class="font-bold block mb-2"
                  >{{ t('urlGenerator.parsedSourceUrl') }}:</span
                >
                <div class="p-2 bg-gray-100 rounded break-all">
                  <a :href="parsedResult.sourceUrl" target="_blank">{{
                    parsedResult.sourceUrl
                  }}</a>
                </div>
              </div>
              <!-- Other Params -->
              <div
                v-if="Object.keys(parsedResult.params).length > 0"
                class="mb-4"
              >
                <span class="font-bold block mb-2"
                  >{{ t('urlGenerator.parsedParams') }}:</span
                >
                <a-descriptions bordered :column="1" class="mt-2">
                  <a-descriptions-item
                    v-for="(val, key) in parsedResult.params"
                    :key="key"
                    :label="key"
                  >
                    {{ val }}
                  </a-descriptions-item>
                </a-descriptions>
              </div>
            </div>
          </a-card>
        </a-tab-pane>
      </a-tabs>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { computed, ref, watch } from 'vue';
  import CraftFlowSelect from '@/views/dashboard/craft_flow/CraftFlowSelect.vue';
  import { useI18n } from 'vue-i18n';
  import { Message } from '@arco-design/web-vue';
  import { useClipboard } from '@vueuse/core';
  import { buildCraftFeedUrl } from '@/utils/publicFeedUrl';

  const { t } = useI18n();

  // Mode 1: Generate
  const selectedCraft = ref('');
  const customCraft = ref('');
  const inputUrl = ref('');
  const resultUrl = ref('');
  const generatedInputUrl = ref('');
  const generatedCraft = ref('');
  const copyButtonText = ref(t('urlGenerator.copyUrl'));
  const { copy: copyResultUrl } = useClipboard({
    source: resultUrl,
    legacy: true,
    copiedDuring: 1500,
  });
  const currentCraft = computed(() =>
    customCraft.value ? customCraft.value : selectedCraft.value
  );
  const isResultStale = computed(
    () =>
      !!resultUrl.value &&
      (inputUrl.value !== generatedInputUrl.value ||
        currentCraft.value !== generatedCraft.value)
  );
  const canCopyResult = computed(
    () => !!resultUrl.value && !isResultStale.value
  );

  const generateUrl = () => {
    const currentSelectedCraft = currentCraft.value;
    resultUrl.value = buildCraftFeedUrl(currentSelectedCraft, inputUrl.value);
    generatedInputUrl.value = inputUrl.value;
    generatedCraft.value = currentSelectedCraft;
    copyButtonText.value = t('urlGenerator.copyUrl');
  };

  watch([inputUrl, selectedCraft, customCraft], () => {
    if (resultUrl.value) {
      copyButtonText.value = t('urlGenerator.copyUrl');
    }
  });

  const copyUrl = async () => {
    if (!canCopyResult.value) {
      return;
    }

    try {
      await copyResultUrl();
      copyButtonText.value = t('urlGenerator.copied');
    } catch {
      Message.error(t('urlGenerator.copyError'));
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
  .container {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 90vh;
  }

  .result-panel {
    border: 1px solid var(--color-neutral-3);
    border-radius: 12px;
    background: var(--color-fill-1);
    padding: 16px;
    transition: border-color 0.2s ease, background-color 0.2s ease,
      box-shadow 0.2s ease;
  }

  .result-panel--fresh {
    border-color: rgb(var(--green-5));
    background: rgba(var(--green-1), 0.45);
  }

  .result-panel--stale {
    border-color: rgb(var(--orange-5));
    background: rgba(var(--orange-1), 0.55);
    box-shadow: 0 0 0 3px rgba(var(--orange-6), 0.08);
  }

  .result-url {
    min-height: 40px;
    border-radius: 8px;
    background: var(--color-bg-2);
    color: var(--color-text-1);
    overflow-wrap: anywhere;
    padding: 10px 12px;
  }
</style>
