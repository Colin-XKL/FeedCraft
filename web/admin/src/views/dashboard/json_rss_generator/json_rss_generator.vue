<template>
  <div class="container">
    <Breadcrumb :items="['menu.tools', 'menu.jsonRssGenerator']" />

    <a-space direction="vertical" size="large" fill>
      <a-card title="Step 1: Request Configuration">
        <template #extra>
          <a-space>
            <a-button type="primary" status="success" @click="handleParseCurl">
              <template #icon><icon-import /></template>
              Import from Curl
            </a-button>
            <a-button type="primary" :loading="fetching" @click="handleFetch">
              <template #icon><icon-send /></template>
              Fetch JSON
            </a-button>
          </a-space>
        </template>

        <a-form :model="fetchReq" layout="vertical">
          <a-form-item
            v-if="showCurlInput"
            label="Curl Command (Optional - Paste here and click Import)"
          >
            <a-textarea
              v-model="curlInput"
              :auto-size="{ minRows: 3, maxRows: 6 }"
              placeholder="curl -X POST ..."
            />
          </a-form-item>

          <a-row :gutter="16">
            <a-col :span="6">
              <a-form-item label="Method" field="method">
                <a-select v-model="fetchReq.method">
                  <a-option>GET</a-option>
                  <a-option>POST</a-option>
                  <a-option>PUT</a-option>
                  <a-option>DELETE</a-option>
                </a-select>
              </a-form-item>
            </a-col>
            <a-col :span="18">
              <a-form-item label="URL" field="url">
                <a-input
                  v-model="fetchReq.url"
                  placeholder="https://api.example.com/v1/posts"
                />
              </a-form-item>
            </a-col>
          </a-row>

          <a-form-item label="Headers (JSON format)" field="headers">
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
                  ><icon-delete
                /></a-button>
              </div>
              <a-space>
                <a-input v-model="newHeaderKey" placeholder="Key" />
                <a-input v-model="newHeaderVal" placeholder="Value" />
                <a-button @click="addHeader">Add Header</a-button>
              </a-space>
            </a-space>
          </a-form-item>

          <a-form-item label="Request Body" field="body">
            <a-textarea
              v-model="fetchReq.body"
              :auto-size="{ minRows: 3, maxRows: 10 }"
              placeholder="{ 'foo': 'bar' }"
            />
          </a-form-item>
        </a-form>
      </a-card>

      <a-card v-if="jsonContent" title="Step 2: JQ Parsing Rules">
        <template #extra>
          <a-button type="primary" :loading="parsing" @click="handlePreview">
            <template #icon><icon-eye /></template>
            Preview RSS
          </a-button>
        </template>
        <div class="split-view">
          <div class="json-view">
            <h3>Response JSON</h3>
            <a-textarea
              v-model="jsonContent"
              :auto-size="{ minRows: 10, maxRows: 30 }"
              style="font-family: monospace"
            />
          </div>
          <div class="rules-view">
            <a-form :model="parseReq" layout="vertical">
              <a-form-item
                label="List Selector (e.g. .data.items or .items[])"
                required
              >
                <a-input v-model="parseReq.list_selector" />
              </a-form-item>
              <a-form-item label="Title Selector (e.g. .title)" required>
                <a-input v-model="parseReq.title_selector" />
              </a-form-item>
              <a-form-item label="Link Selector (e.g. .url)">
                <a-input v-model="parseReq.link_selector" />
              </a-form-item>
              <a-form-item label="Date Selector (e.g. .created_at)">
                <a-input v-model="parseReq.date_selector" />
              </a-form-item>
              <a-form-item label="Content Selector (e.g. .body)">
                <a-input v-model="parseReq.content_selector" />
              </a-form-item>
            </a-form>
          </div>
        </div>
      </a-card>

      <a-card v-if="parsedItems.length > 0" title="Step 3: Preview">
        <a-list :data="parsedItems">
          <template #item="{ item }">
            <a-list-item>
              <a-list-item-meta :title="item.title" :description="item.date">
              </a-list-item-meta>
              <div>
                <a :href="item.link" target="_blank">{{ item.link }}</a>
              </div>
              <div v-if="item.content" class="content-preview">
                {{ item.content.substring(0, 100) }}...
              </div>
            </a-list-item>
          </template>
        </a-list>
      </a-card>
    </a-space>
  </div>
</template>

<script setup lang="ts">
  import { ref, reactive } from 'vue';
  import {
    parseCurl,
    fetchJson,
    parseJsonRss,
    JsonFetchReq,
    ParsedItem,
  } from '@/api/json_rss';
  import { Message } from '@arco-design/web-vue';

  const showCurlInput = ref(true);
  const curlInput = ref('');
  const fetching = ref(false);
  const parsing = ref(false);

  const fetchReq = reactive<JsonFetchReq>({
    method: 'GET',
    url: '',
    headers: {},
    body: '',
  });

  const jsonContent = ref('');
  const parseReq = reactive({
    list_selector: '.',
    title_selector: '.title',
    link_selector: '.url',
    date_selector: '',
    content_selector: '',
  });

  const parsedItems = ref<ParsedItem[]>([]);

  const newHeaderKey = ref('');
  const newHeaderVal = ref('');

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
      // Error handled by interceptor usually, but safe to log
      console.error(err);
    }
  };

  const handleFetch = async () => {
    if (!fetchReq.url) {
      Message.warning('URL is required');
      return;
    }
    fetching.value = true;
    try {
      const res = await fetchJson(fetchReq);
      // The interceptor might unwrap 'data', but our backend returns string in data.
      // If the backend returns APIResponse<string>, the interceptor typically returns T.
      // Let's assume res is string based on api definition
      jsonContent.value = res as unknown as string;
      Message.success('Fetched successfully');
    } finally {
      fetching.value = false;
    }
  };

  const handlePreview = async () => {
    if (!jsonContent.value) return;
    parsing.value = true;
    try {
      const res = await parseJsonRss({
        json_content: jsonContent.value,
        ...parseReq,
      });
      parsedItems.value = res;
      Message.success(`Parsed ${res.length} items`);
    } finally {
      parsing.value = false;
    }
  };
</script>

<style scoped>
  .container {
    padding: 20px;
  }
  .header-row {
    margin-bottom: 8px;
    display: flex;
    gap: 8px;
    align-items: center;
  }
  .split-view {
    display: flex;
    gap: 20px;
  }
  .json-view {
    flex: 1;
  }
  .rules-view {
    flex: 1;
  }
  .content-preview {
    color: var(--color-text-3);
    font-size: 12px;
    margin-top: 4px;
  }
</style>
