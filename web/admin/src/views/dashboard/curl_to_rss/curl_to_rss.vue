<template>
  <div class="py-8 px-16">
    <x-header
      :title="$t('curlToRss.title')"
      :description="$t('curlToRss.description')"
    ></x-header>

    <div class="content-wrapper">
      <a-card class="wizard-card">
        <a-steps :current="currentStep" class="mb-8">
          <a-step :title="$t('curlToRss.step.request')" :description="$t('curlToRss.step.request.desc')" />
          <a-step :title="$t('curlToRss.step.parsing')" :description="$t('curlToRss.step.parsing.desc')" />
          <a-step :title="$t('curlToRss.step.meta')" :description="$t('curlToRss.step.meta.desc')" />
          <a-step :title="$t('curlToRss.step.save')" :description="$t('curlToRss.step.save.desc')" />
        </a-steps>

        <!-- STEP 1: Request Configuration -->
        <div v-show="currentStep === 1" class="step-content">
          <a-space direction="vertical" size="large" fill>
            <a-alert
              >{{ $t('curlToRss.alert.request') }}</a-alert
            >

            <a-form :model="fetchReq" layout="vertical">
              <a-row :gutter="16">
                <a-col :span="24">
                  <a-form-item
                    :label="$t('curlToRss.label.curl')"
                  >
                    <div class="flex w-full gap-2">
                      <a-textarea
                        v-model="curlInput"
                        :auto-size="{ minRows: 2, maxRows: 6 }"
                        :placeholder="$t('curlToRss.placeholder.curl')"
                      />
                      <a-button
                        type="primary"
                        status="success"
                        @click="handleParseCurl"
                      >
                        <template #icon><icon-import /></template>
                        {{ $t('curlToRss.button.import') }}
                      </a-button>
                    </div>
                  </a-form-item>
                </a-col>
              </a-row>

              <a-divider />

              <a-row :gutter="16">
                <a-col :span="6">
                  <a-form-item :label="$t('curlToRss.label.method')" field="method" required>
                    <a-select v-model="fetchReq.method">
                      <a-option>GET</a-option>
                      <a-option>POST</a-option>
                    </a-select>
                  </a-form-item>
                </a-col>
                <a-col :span="18">
                  <a-form-item :label="$t('curlToRss.label.url')" field="url" required>
                    <a-input
                      v-model="fetchReq.url"
                      :placeholder="$t('curlToRss.placeholder.url')"
                      @keyup.enter="handleFetchAndNext"
                    />
                  </a-form-item>
                </a-col>
              </a-row>

              <a-form-item :label="$t('curlToRss.label.headers')" field="headers">
                <a-space direction="vertical" fill>
                  <div
                    v-for="(val, key) in fetchReq.headers"
                    :key="key"
                    class="header-row"
                  >
                    <a-input
                      :model-value="String(key)"
                      readonly
                      style="width: 30%"
                    />
                    <a-input
                      :model-value="String(val)"
                      readonly
                      style="width: 60%"
                    />
                    <a-button
                      type="text"
                      status="danger"
                      @click="removeHeader(String(key))"
                    >
                      <icon-delete />
                    </a-button>
                  </div>
                  <div class="flex gap-2">
                    <a-input
                      v-model="newHeaderKey"
                      :placeholder="$t('curlToRss.placeholder.headerKey')"
                      style="width: 30%"
                    />
                    <a-input
                      v-model="newHeaderVal"
                      :placeholder="$t('curlToRss.placeholder.headerVal')"
                      style="width: 60%"
                    />
                    <a-button @click="addHeader">{{ $t('curlToRss.button.addHeader') }}</a-button>
                  </div>
                </a-space>
              </a-form-item>

              <a-form-item :label="$t('curlToRss.label.body')" field="body">
                <a-textarea
                  v-model="fetchReq.body"
                  :auto-size="{ minRows: 3, maxRows: 10 }"
                  :placeholder="$t('curlToRss.placeholder.body')"
                />
              </a-form-item>

              <div class="text-center mt-8">
                <a-button
                  type="primary"
                  size="large"
                  :loading="fetching"
                  :disabled="!fetchReq.url"
                  @click="handleFetchAndNext"
                >
                  {{ $t('curlToRss.button.fetch') }} <icon-arrow-right />
                </a-button>
              </div>
            </a-form>
          </a-space>
        </div>

        <!-- STEP 2: Parsing Rules -->
        <div v-show="currentStep === 2" class="step-content h-full">
          <a-row :gutter="16" class="h-full">
            <!-- Left: JSON View -->
            <a-col :span="12" class="h-full flex flex-col">
              <div class="font-bold mb-2">{{ $t('curlToRss.label.jsonView') }}</div>
              <a-textarea
                v-model="jsonContent"
                class="flex-1 font-mono text-xs"
                :auto-size="false"
                style="height: 100%; font-family: monospace"
                readonly
              />
            </a-col>

            <!-- Right: Rules & Preview -->
            <a-col :span="12" class="h-full flex flex-col">
              <div class="flex-1 overflow-y-auto pr-2">
                <a-alert type="info" class="mb-4">
                  <span v-html="$t('curlToRss.alert.rules')"></span>
                </a-alert>

                <a-form :model="parseReq" layout="vertical">
                  <a-card :title="$t('curlToRss.group.iteration')" size="small" class="mb-4">
                    <a-form-item
                      :label="$t('curlToRss.label.iterator')"
                      required
                    >
                      <a-input
                        v-model="parseReq.list_selector"
                        :placeholder="$t('curlToRss.placeholder.iterator')"
                      />
                    </a-form-item>
                  </a-card>

                  <a-card
                    :title="$t('curlToRss.group.fields')"
                    size="small"
                  >
                    <a-form-item :label="$t('curlToRss.label.titleSelector')" required>
                      <a-input
                        v-model="parseReq.title_selector"
                        :placeholder="$t('curlToRss.placeholder.titleSelector')"
                      />
                    </a-form-item>
                    <a-form-item :label="$t('curlToRss.label.linkSelector')">
                      <a-input
                        v-model="parseReq.link_selector"
                        :placeholder="$t('curlToRss.placeholder.linkSelector')"
                      />
                    </a-form-item>
                    <a-form-item :label="$t('curlToRss.label.dateSelector')">
                      <a-input
                        v-model="parseReq.date_selector"
                        :placeholder="$t('curlToRss.placeholder.dateSelector')"
                      />
                    </a-form-item>
                    <a-form-item :label="$t('curlToRss.label.contentSelector')">
                      <a-input
                        v-model="parseReq.content_selector"
                        :placeholder="$t('curlToRss.placeholder.contentSelector')"
                      />
                    </a-form-item>
                  </a-card>
                </a-form>

                <!-- Preview Results -->
                <div v-if="parsedItems.length > 0" class="mt-4">
                  <a-divider orientation="left">
                    {{ $t('curlToRss.preview.title', { count: parsedItems.length }) }}
                  </a-divider>
                  <a-collapse :default-active-key="[0]">
                    <a-collapse-item
                      v-for="(item, idx) in parsedItems"
                      :key="idx"
                      :header="item.title || 'No Title'"
                    >
                      <div class="text-xs text-gray-500 mb-1">
                        {{ item.link }}
                      </div>
                      <div class="text-xs text-gray-400 mb-2">
                        {{ item.date }}
                      </div>
                      <div class="text-sm text-gray-600 truncate">
                        {{ item.content }}
                      </div>
                    </a-collapse-item>
                  </a-collapse>
                </div>
              </div>

              <!-- Footer -->
              <div
                class="flex justify-between mt-4 pt-4 border-t border-gray-100 bg-white"
              >
                <a-button @click="prevStep">{{ $t('searchToRss.button.back') }}</a-button>
                <a-space>
                  <a-button
                    type="outline"
                    :loading="parsing"
                    @click="handlePreview"
                  >
                    {{ $t('curlToRss.button.runPreview') }}
                  </a-button>
                  <a-button
                    type="primary"
                    :disabled="parsedItems.length === 0"
                    @click="nextStep"
                  >
                    {{ $t('searchToRss.button.next') }}
                  </a-button>
                </a-space>
              </div>
            </a-col>
          </a-row>
        </div>

        <!-- STEP 3: Feed Metadata -->
        <div v-show="currentStep === 3" class="step-content">
          <div class="max-w-2xl mx-auto">
            <a-alert type="success" class="mb-6">
              {{ $t('curlToRss.alert.meta', { count: parsedItems.length }) }}
            </a-alert>

            <a-form :model="feedMeta" layout="vertical">
              <a-form-item :label="$t('curlToRss.label.feedTitle')" required>
                <a-input
                  v-model="feedMeta.title"
                  :placeholder="$t('curlToRss.placeholder.feedTitle')"
                />
              </a-form-item>
              <a-form-item :label="$t('curlToRss.label.feedDesc')">
                <a-textarea
                  v-model="feedMeta.description"
                  :placeholder="$t('curlToRss.placeholder.feedDesc')"
                />
              </a-form-item>
              <a-form-item :label="$t('curlToRss.label.siteLink')">
                <a-input
                  v-model="feedMeta.link"
                  :placeholder="$t('curlToRss.placeholder.siteLink')"
                />
              </a-form-item>
              <a-row :gutter="16">
                <a-col :span="12">
                  <a-form-item :label="$t('curlToRss.label.authorName')">
                    <a-input v-model="feedMeta.author_name" />
                  </a-form-item>
                </a-col>
                <a-col :span="12">
                  <a-form-item :label="$t('curlToRss.label.authorEmail')">
                    <a-input v-model="feedMeta.author_email" />
                  </a-form-item>
                </a-col>
              </a-row>
            </a-form>

            <div class="flex justify-between mt-8">
              <a-button @click="prevStep">{{ $t('searchToRss.button.back') }}</a-button>
              <a-button type="primary" @click="handleStep3Next">{{ $t('searchToRss.button.next') }}</a-button>
            </div>
          </div>
        </div>

        <!-- STEP 4: Save -->
        <div v-show="currentStep === 4" class="step-content">
          <div class="max-w-xl mx-auto">
            <a-card :title="$t('curlToRss.card.review')" class="border-blue-100">
              <a-descriptions :column="1" title="Summary" bordered>
                <a-descriptions-item :label="$t('curlToRss.label.sourceUrl')">
                  {{ fetchReq.url }}
                </a-descriptions-item>
                <a-descriptions-item :label="$t('curlToRss.label.feedTitle')">
                  {{ feedMeta.title }}
                </a-descriptions-item>
                <a-descriptions-item :label="$t('curlToRss.label.itemCount')">
                  {{ parsedItems.length }}
                </a-descriptions-item>
              </a-descriptions>

              <a-divider />

              <a-form :model="recipeMeta" layout="vertical" class="mt-6">
                <a-form-item
                  :label="$t('curlToRss.label.recipeId')"
                  required
                  :help="$t('curlToRss.help.recipeId')"
                >
                  <a-input
                    v-model="recipeMeta.id"
                    :placeholder="$t('curlToRss.placeholder.recipeId')"
                  />
                </a-form-item>
                <a-form-item :label="$t('curlToRss.label.internalDesc')">
                  <a-textarea
                    v-model="recipeMeta.description"
                    :placeholder="$t('curlToRss.placeholder.internalDesc')"
                  />
                </a-form-item>

                <div class="mt-8 text-center">
                  <a-button
                    type="primary"
                    long
                    size="large"
                    status="success"
                    :loading="saving"
                    @click="handleSaveRecipe"
                  >
                    <icon-save /> {{ $t('curlToRss.button.save') }}
                  </a-button>
                </div>
              </a-form>
            </a-card>

            <div class="flex justify-start mt-8">
              <a-button @click="prevStep">{{ $t('searchToRss.button.back') }}</a-button>
            </div>
          </div>
        </div>
      </a-card>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, reactive } from 'vue';
  import { useRouter } from 'vue-router';
  import { Message } from '@arco-design/web-vue';
  import {
    IconImport,
    IconDelete,
    IconArrowRight,
    IconSave,
  } from '@arco-design/web-vue/es/icon';
  import { useI18n } from 'vue-i18n';
  import XHeader from '@/components/header/x-header.vue';
  import {
    parseCurl,
    fetchJson,
    parseJsonRss,
    JsonFetchReq,
    ParsedItem,
  } from '@/api/json_rss';
  import { createCustomRecipe } from '@/api/custom_recipe';

  const router = useRouter();
  const { t } = useI18n();

  // --- State ---
  const currentStep = ref(1);
  const fetching = ref(false);
  const parsing = ref(false);
  const saving = ref(false);

  // Step 1 State
  const curlInput = ref('');
  const fetchReq = reactive<JsonFetchReq>({
    method: 'GET',
    url: '',
    headers: {},
    body: '',
  });
  const newHeaderKey = ref('');
  const newHeaderVal = ref('');

  // Step 2 State
  const jsonContent = ref('');
  const parseReq = reactive({
    list_selector: '',
    title_selector: '',
    link_selector: '',
    date_selector: '',
    content_selector: '',
  });
  const parsedItems = ref<ParsedItem[]>([]);

  // Step 3 State
  const feedMeta = reactive({
    title: '',
    link: '',
    description: '',
    author_name: '',
    author_email: '',
  });

  // Step 4 State
  const recipeMeta = reactive({
    id: '',
    description: '',
  });

  // --- Actions ---

  const nextStep = () => {
    currentStep.value += 1;
  };

  const prevStep = () => {
    if (currentStep.value > 1) currentStep.value -= 1;
  };

  // Step 1 Logic
  const addHeader = () => {
    if (newHeaderKey.value && newHeaderVal.value) {
      fetchReq.headers[newHeaderKey.value] = newHeaderVal.value;
      newHeaderKey.value = '';
      newHeaderVal.value = '';
    }
  };

  const removeHeader = (key: string) => {
    delete fetchReq.headers[key];
  };

  const handleParseCurl = async () => {
    if (!curlInput.value) {
      Message.warning(t('curlToRss.validation.curl'));
      return;
    }
    try {
      const res = await parseCurl(curlInput.value);
      if (res.data) {
        fetchReq.method = res.data.method;
        fetchReq.url = res.data.url;
        fetchReq.headers = res.data.headers || {};
        fetchReq.body = res.data.body || '';
        Message.success(t('curlToRss.message.curlParsed'));
      }
    } catch (err) {
      // Error handled by interceptor usually
      console.error(err);
    }
  };

  const handleFetchAndNext = async () => {
    if (!fetchReq.url) {
      Message.warning(t('curlToRss.validation.url'));
      return;
    }
    fetching.value = true;
    try {
      const res = await fetchJson(fetchReq);
      jsonContent.value = res.data;
      if (jsonContent.value) {
        // Auto-fill link in meta if possible
        feedMeta.link = fetchReq.url;
        Message.success(t('curlToRss.message.fetched'));
        nextStep();
      } else {
        Message.warning(t('curlToRss.message.empty'));
      }
    } catch (err) {
      console.error(err);
    } finally {
      fetching.value = false;
    }
  };

  // Step 2 Logic
  const handlePreview = async () => {
    if (!jsonContent.value) return;
    if (!parseReq.list_selector) {
      Message.warning(t('curlToRss.validation.iterator'));
      return;
    }
    if (!parseReq.title_selector) {
      Message.warning(t('curlToRss.validation.titleSelector'));
      return;
    }
    parsing.value = true;
    try {
      const res = await parseJsonRss({
        json_content: jsonContent.value,
        ...parseReq,
      });
      parsedItems.value = res.data || [];
      if (parsedItems.value.length === 0) {
        Message.warning(t('curlToRss.message.noItems'));
      } else {
        Message.success(t('curlToRss.message.parsed', { count: parsedItems.value.length }));
      }
    } catch (err) {
      console.error(err);
    } finally {
      parsing.value = false;
    }
  };

  // Step 3 Logic
  const handleStep3Next = () => {
    if (!feedMeta.title.trim()) {
      Message.warning(t('curlToRss.validation.feedTitle'));
      return;
    }
    nextStep();
  };

  // Step 4 Logic
  const handleSaveRecipe = async () => {
    if (!recipeMeta.id) {
      Message.warning(t('curlToRss.validation.recipeId'));
      return;
    }

    saving.value = true;

    // Construct the SourceConfig object conforming to internal/config/source_config.go
    // Note: The backend expects snake_case json structure in the stringified config
    const sourceConfig = {
      type: 'json',
      http_fetcher: {
        url: fetchReq.url,
        headers: fetchReq.headers,
      },
      json_parser: {
        items_iterator: parseReq.list_selector,
        title: parseReq.title_selector,
        link: parseReq.link_selector,
        date: parseReq.date_selector,
        description: parseReq.content_selector,
      },
      feed_meta: {
        title: feedMeta.title,
        link: feedMeta.link,
        description: feedMeta.description,
        author_name: feedMeta.author_name,
        author_email: feedMeta.author_email,
      },
    };

    try {
      await createCustomRecipe({
        id: recipeMeta.id,
        description: recipeMeta.description,
        craft: 'proxy',
        source_type: 'json',
        source_config: JSON.stringify(sourceConfig),
      });

      Message.success(t('curlToRss.message.saved'));
      router.push({ name: 'CustomRecipe' });
    } catch (err: any) {
      Message.error(`Failed to save: ${err.message || err}`);
    } finally {
      saving.value = false;
    }
  };
</script>

<style scoped>
  .wizard-card {
    min-height: 700px;
  }

  .step-content {
    margin-top: 24px;
    height: 600px;
    /* Ensure content takes available space */
    display: flex;
    flex-direction: column;
  }

  .header-row {
    margin-bottom: 8px;
    display: flex;
    gap: 8px;
    align-items: center;
  }
</style>
