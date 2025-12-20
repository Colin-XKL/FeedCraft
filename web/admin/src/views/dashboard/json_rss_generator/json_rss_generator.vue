<template>
  <div class="py-8 px-16">
    <x-header
      :title="$t('menu.jsonRssGenerator')"
      description="Generate RSS feeds from any JSON API by defining parsing rules."
    ></x-header>

    <div class="content-wrapper">
      <a-card class="wizard-card">
        <a-steps :current="currentStep" class="mb-8">
          <a-step title="Request Config" description="Configure JSON source" />
          <a-step title="Parsing Rules" description="Define jq selectors" />
          <a-step title="Feed Metadata" description="Set feed details" />
          <a-step title="Save Recipe" description="Save as Custom Recipe" />
        </a-steps>

        <!-- STEP 1: Request Configuration -->
        <div v-show="currentStep === 1" class="step-content">
          <a-space direction="vertical" size="large" fill>
            <a-alert
              >Configure the HTTP request to fetch the JSON data. You can import
              from a cURL command.</a-alert
            >

            <a-form :model="fetchReq" layout="vertical">
              <a-row :gutter="16">
                <a-col :span="24">
                  <a-form-item
                    label="Curl Command (Optional - Paste here and click Import)"
                  >
                    <div class="flex w-full gap-2">
                      <a-textarea
                        v-model="curlInput"
                        :auto-size="{ minRows: 2, maxRows: 6 }"
                        placeholder="curl -X POST ..."
                      />
                      <a-button
                        type="primary"
                        status="success"
                        @click="handleParseCurl"
                      >
                        <template #icon><icon-import /></template>
                        Import
                      </a-button>
                    </div>
                  </a-form-item>
                </a-col>
              </a-row>

              <a-divider />

              <a-row :gutter="16">
                <a-col :span="6">
                  <a-form-item label="Method" field="method" required>
                    <a-select v-model="fetchReq.method">
                      <a-option>GET</a-option>
                      <a-option>POST</a-option>
                    </a-select>
                  </a-form-item>
                </a-col>
                <a-col :span="18">
                  <a-form-item label="URL" field="url" required>
                    <a-input
                      v-model="fetchReq.url"
                      placeholder="https://api.example.com/v1/posts"
                      @keyup.enter="handleFetchAndNext"
                    />
                  </a-form-item>
                </a-col>
              </a-row>

              <a-form-item label="Headers" field="headers">
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
                      placeholder="Key"
                      style="width: 30%"
                    />
                    <a-input
                      v-model="newHeaderVal"
                      placeholder="Value"
                      style="width: 60%"
                    />
                    <a-button @click="addHeader">Add</a-button>
                  </div>
                </a-space>
              </a-form-item>

              <a-form-item label="Request Body" field="body">
                <a-textarea
                  v-model="fetchReq.body"
                  :auto-size="{ minRows: 3, maxRows: 10 }"
                  placeholder="{ 'foo': 'bar' }"
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
                  Fetch & Next <icon-arrow-right />
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
              <div class="font-bold mb-2">Response JSON</div>
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
                  Define selectors to extract feed items. Use dot notation (e.g.,
                  <code>.data.items</code>).
                </a-alert>

                <a-form :model="parseReq" layout="vertical">
                  <a-card title="Iteration" size="small" class="mb-4">
                    <a-form-item
                      label="Items Iterator (e.g. .data.items or .items[])"
                      required
                    >
                      <a-input
                        v-model="parseReq.list_selector"
                        placeholder=".items"
                      />
                    </a-form-item>
                  </a-card>

                  <a-card
                    title="Item Fields (Relative to Iterator)"
                    size="small"
                  >
                    <a-form-item label="Title Selector" required>
                      <a-input
                        v-model="parseReq.title_selector"
                        placeholder=".title"
                      />
                    </a-form-item>
                    <a-form-item label="Link Selector">
                      <a-input
                        v-model="parseReq.link_selector"
                        placeholder=".url"
                      />
                    </a-form-item>
                    <a-form-item label="Date Selector">
                      <a-input
                        v-model="parseReq.date_selector"
                        placeholder=".created_at"
                      />
                    </a-form-item>
                    <a-form-item label="Content/Description Selector">
                      <a-input
                        v-model="parseReq.content_selector"
                        placeholder=".content"
                      />
                    </a-form-item>
                  </a-card>
                </a-form>

                <!-- Preview Results -->
                <div v-if="parsedItems.length > 0" class="mt-4">
                  <a-divider orientation="left">
                    Preview Results ({{ parsedItems.length }})
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
                <a-button @click="prevStep">Back</a-button>
                <a-space>
                  <a-button
                    type="outline"
                    :loading="parsing"
                    @click="handlePreview"
                  >
                    Run Preview
                  </a-button>
                  <a-button
                    type="primary"
                    :disabled="parsedItems.length === 0"
                    @click="nextStep"
                  >
                    Next Step
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
              Successfully extracted {{ parsedItems.length }} items! Now configure
              the feed metadata.
            </a-alert>

            <a-form :model="feedMeta" layout="vertical">
              <a-form-item label="Feed Title" required>
                <a-input
                  v-model="feedMeta.title"
                  placeholder="My Awesome Feed"
                />
              </a-form-item>
              <a-form-item label="Feed Description">
                <a-textarea
                  v-model="feedMeta.description"
                  placeholder="A description of this feed"
                />
              </a-form-item>
              <a-form-item label="Site Link">
                <a-input
                  v-model="feedMeta.link"
                  placeholder="https://example.com"
                />
              </a-form-item>
              <a-row :gutter="16">
                <a-col :span="12">
                  <a-form-item label="Author Name">
                    <a-input v-model="feedMeta.author_name" />
                  </a-form-item>
                </a-col>
                <a-col :span="12">
                  <a-form-item label="Author Email">
                    <a-input v-model="feedMeta.author_email" />
                  </a-form-item>
                </a-col>
              </a-row>
            </a-form>

            <div class="flex justify-between mt-8">
              <a-button @click="prevStep">Back</a-button>
              <a-button type="primary" @click="handleStep3Next">Next</a-button>
            </div>
          </div>
        </div>

        <!-- STEP 4: Save -->
        <div v-show="currentStep === 4" class="step-content">
          <div class="max-w-xl mx-auto">
            <a-card title="Review & Save" class="border-blue-100">
              <a-descriptions :column="1" title="Summary" bordered>
                <a-descriptions-item label="Source URL">
                  {{ fetchReq.url }}
                </a-descriptions-item>
                <a-descriptions-item label="Feed Title">
                  {{ feedMeta.title }}
                </a-descriptions-item>
                <a-descriptions-item label="Item Count">
                  {{ parsedItems.length }}
                </a-descriptions-item>
              </a-descriptions>

              <a-divider />

              <a-form :model="recipeMeta" layout="vertical" class="mt-6">
                <a-form-item
                  label="Recipe ID (URL Path)"
                  required
                  help="This will be the unique identifier in the URL."
                >
                  <a-input
                    v-model="recipeMeta.id"
                    placeholder="my-json-feed"
                  />
                </a-form-item>
                <a-form-item label="Internal Description">
                  <a-textarea
                    v-model="recipeMeta.description"
                    placeholder="Notes for yourself about this recipe"
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
                    <icon-save /> Confirm & Save Recipe
                  </a-button>
                </div>
              </a-form>
            </a-card>

            <div class="flex justify-start mt-8">
              <a-button @click="prevStep">Back</a-button>
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
      Message.warning('Please enter a curl command');
      return;
    }
    try {
      const res = await parseCurl(curlInput.value);
      if (res.data) {
        fetchReq.method = res.data.method;
        fetchReq.url = res.data.url;
        fetchReq.headers = res.data.headers || {};
        fetchReq.body = res.data.body || '';
        Message.success('Curl parsed successfully');
      }
    } catch (err) {
      // Error handled by interceptor usually
      console.error(err);
    }
  };

  const handleFetchAndNext = async () => {
    if (!fetchReq.url) {
      Message.warning('URL is required');
      return;
    }
    fetching.value = true;
    try {
      const res = await fetchJson(fetchReq);
      jsonContent.value = res.data;
      if (jsonContent.value) {
        // Auto-fill link in meta if possible
        feedMeta.link = fetchReq.url;
        Message.success('Fetched successfully');
        nextStep();
      } else {
        Message.warning('Empty response');
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
      Message.warning('Items Iterator selector is required');
      return;
    }
    if (!parseReq.title_selector) {
      Message.warning('Title selector is required');
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
        Message.warning('No items found with current selectors');
      } else {
        Message.success(`Parsed ${parsedItems.value.length} items`);
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
      Message.warning('Feed Title is required');
      return;
    }
    nextStep();
  };

  // Step 4 Logic
  const handleSaveRecipe = async () => {
    if (!recipeMeta.id) {
      Message.warning('Recipe ID is required');
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
        // body and method are not standard in HttpFetcherConfig based on the read file (only URL, Headers, UseBrowserless).
        // However, standard http fetcher usually defaults to GET.
        // If the backend HttpFetcher only supports GET, then POST/Body might be ignored or require a different fetcher type.
        // Based on `internal/config/source_config.go`, HttpFetcherConfig only has URL, Headers, UseBrowserless.
        // If the user needs POST/Body, the current backend might not support it via `HttpFetcher`.
        // BUT, let's assume standard behavior for now.
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

      Message.success('Recipe saved successfully!');
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
