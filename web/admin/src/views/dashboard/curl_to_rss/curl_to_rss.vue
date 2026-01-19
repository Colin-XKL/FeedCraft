<template>
  <div class="py-8 px-16">
    <x-header
      :title="$t('menu.curlToRss')"
      :description="$t('curlToRss.description')"
    ></x-header>

    <div class="content-wrapper">
      <a-card class="wizard-card">
        <a-steps
          :current="currentStep"
          changeable
          class="mb-8"
          @change="onStepChange"
        >
          <a-step
            :title="$t('curlToRss.step.requestConfig')"
            :description="$t('curlToRss.step.requestConfig.desc')"
          />
          <a-step
            :title="$t('curlToRss.step.parsingRules')"
            :description="$t('curlToRss.step.parsingRules.desc')"
          />
          <a-step
            :title="$t('curlToRss.step.feedMetadata')"
            :description="$t('curlToRss.step.feedMetadata.desc')"
          />
          <a-step
            :title="$t('curlToRss.step.saveRecipe')"
            :description="$t('curlToRss.step.saveRecipe.desc')"
          />
        </a-steps>

        <!-- STEP 1: Request Configuration -->
        <div v-show="currentStep === 1" class="step-content">
          <a-space direction="vertical" size="large" fill>
            <a-alert>{{ $t('curlToRss.step1.alert') }}</a-alert>
            <a-alert v-if="fetchError" type="error" show-icon>
              {{ fetchError }}
            </a-alert>

            <a-form :model="fetchReq" layout="vertical">
              <a-row :gutter="16">
                <a-col :span="24">
                  <a-form-item :label="$t('curlToRss.step1.curlCommand')">
                    <div class="flex w-full gap-2">
                      <a-textarea
                        v-model="curlInput"
                        :auto-size="{ minRows: 2, maxRows: 6 }"
                        :placeholder="$t('curlToRss.placeholder.curl')"
                      />
                      <a-button
                        type="primary"
                        status="success"
                        :loading="parsingCurl"
                        @click="handleParseCurl"
                      >
                        <template #icon><icon-import /></template>
                        {{ $t('curlToRss.step1.import') }}
                      </a-button>
                    </div>
                  </a-form-item>
                </a-col>
              </a-row>

              <a-divider />

              <a-row :gutter="16">
                <a-col :span="6">
                  <a-form-item
                    :label="$t('curlToRss.step1.method')"
                    field="method"
                    required
                  >
                    <a-select v-model="fetchReq.method">
                      <a-option>GET</a-option>
                      <a-option>POST</a-option>
                    </a-select>
                  </a-form-item>
                </a-col>
                <a-col :span="18">
                  <a-form-item
                    :label="$t('curlToRss.step1.url')"
                    field="url"
                    required
                  >
                    <a-input
                      v-model="fetchReq.url"
                      :placeholder="$t('curlToRss.placeholder.url')"
                      @keyup.enter="handleFetchAndNext"
                    />
                  </a-form-item>
                </a-col>
              </a-row>

              <a-form-item
                :label="$t('curlToRss.step1.headers')"
                field="headers"
              >
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
                      :placeholder="$t('curlToRss.placeholder.key')"
                      style="width: 30%"
                    />
                    <a-input
                      v-model="newHeaderVal"
                      :placeholder="$t('curlToRss.placeholder.value')"
                      style="width: 60%"
                    />
                    <a-button @click="addHeader">{{
                      $t('curlToRss.step1.add')
                    }}</a-button>
                  </div>
                </a-space>
              </a-form-item>

              <a-form-item
                :label="$t('curlToRss.step1.requestBody')"
                field="body"
              >
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
                  {{ $t('curlToRss.step1.fetchAndNext') }} <icon-arrow-right />
                </a-button>
              </div>
            </a-form>
          </a-space>
        </div>

        <!-- STEP 2: Parsing Rules -->
        <div
          v-show="currentStep === 2"
          class="step-content h-full step-fixed-height"
        >
          <a-row :gutter="16" class="h-full">
            <!-- Left: JSON View -->
            <a-col :span="12" class="h-full flex flex-col">
              <div class="font-bold mb-2">
                {{ $t('curlToRss.step2.responseJson') }}
              </div>
              <div
                class="flex-1 overflow-auto border border-gray-200 rounded p-2 bg-gray-50"
              >
                <a-tree
                  v-if="treeData.length"
                  :data="treeData"
                  :show-line="true"
                  block-node
                  @select="handleNodeSelect"
                />
                <div v-else class="text-gray-400 text-center mt-4">
                  {{ jsonContent ? 'Invalid JSON' : 'No Data' }}
                </div>
              </div>
            </a-col>

            <!-- Right: Rules & Preview -->
            <a-col :span="12" class="h-full flex flex-col">
              <div class="flex-1 overflow-y-auto pr-2">
                <a-alert type="info" class="mb-4">
                  <span v-html="$t('curlToRss.step2.alert')"></span>
                </a-alert>

                <a-form :model="parseReq" layout="vertical">
                  <a-card
                    :title="$t('curlToRss.step2.iteration')"
                    size="small"
                    class="mb-4"
                  >
                    <a-form-item
                      :label="$t('curlToRss.step2.itemsIterator')"
                      required
                    >
                      <template #label>
                        {{ $t('curlToRss.step2.itemsIterator') }}
                        <icon-edit
                          v-if="activeField === 'list_selector'"
                          class="ml-2 text-primary"
                        />
                      </template>
                      <a-input
                        v-model="parseReq.list_selector"
                        :placeholder="$t('curlToRss.placeholder.items')"
                        :class="{
                          'border-primary': activeField === 'list_selector',
                        }"
                        @focus="activeField = 'list_selector'"
                      />
                    </a-form-item>
                  </a-card>

                  <a-card
                    :title="$t('curlToRss.step2.itemFields')"
                    size="small"
                  >
                    <a-form-item
                      :label="$t('curlToRss.step2.titleSelector')"
                      required
                    >
                      <template #label>
                        {{ $t('curlToRss.step2.titleSelector') }}
                        <icon-edit
                          v-if="activeField === 'title_selector'"
                          class="ml-2 text-primary"
                        />
                      </template>
                      <a-input
                        v-model="parseReq.title_selector"
                        :placeholder="$t('curlToRss.placeholder.title')"
                        :class="{
                          'border-primary': activeField === 'title_selector',
                        }"
                        @focus="activeField = 'title_selector'"
                      />
                    </a-form-item>
                    <a-form-item :label="$t('curlToRss.step2.linkSelector')">
                      <template #label>
                        {{ $t('curlToRss.step2.linkSelector') }}
                        <icon-edit
                          v-if="activeField === 'link_selector'"
                          class="ml-2 text-primary"
                        />
                      </template>
                      <a-input
                        v-model="parseReq.link_selector"
                        :placeholder="$t('curlToRss.placeholder.link')"
                        :class="{
                          'border-primary': activeField === 'link_selector',
                        }"
                        @focus="activeField = 'link_selector'"
                      />
                    </a-form-item>
                    <a-form-item :label="$t('curlToRss.step2.dateSelector')">
                      <template #label>
                        {{ $t('curlToRss.step2.dateSelector') }}
                        <icon-edit
                          v-if="activeField === 'date_selector'"
                          class="ml-2 text-primary"
                        />
                      </template>
                      <a-input
                        v-model="parseReq.date_selector"
                        :placeholder="$t('curlToRss.placeholder.date')"
                        :class="{
                          'border-primary': activeField === 'date_selector',
                        }"
                        @focus="activeField = 'date_selector'"
                      />
                    </a-form-item>
                    <a-form-item :label="$t('curlToRss.step2.contentSelector')">
                      <template #label>
                        {{ $t('curlToRss.step2.contentSelector') }}
                        <icon-edit
                          v-if="activeField === 'content_selector'"
                          class="ml-2 text-primary"
                        />
                      </template>
                      <a-input
                        v-model="parseReq.content_selector"
                        :placeholder="$t('curlToRss.placeholder.content')"
                        :class="{
                          'border-primary': activeField === 'content_selector',
                        }"
                        @focus="activeField = 'content_selector'"
                      />
                    </a-form-item>
                  </a-card>
                </a-form>

                <!-- Preview Results -->
                <div v-if="parsedItems.length > 0" class="mt-4">
                  <a-divider orientation="left">
                    {{
                      $t('curlToRss.step2.previewResults', {
                        count: parsedItems.length,
                      })
                    }}
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
                <!-- Empty State -->
                <div v-else class="mt-8 text-center text-gray-400">
                  <a-empty
                    :description="$t('curlToRss.step2.previewPlaceholder')"
                  >
                    <template #extra>
                      {{ $t('curlToRss.step2.previewPlaceholder.help') }}
                    </template>
                  </a-empty>
                </div>
              </div>

              <!-- Footer -->
              <div
                class="flex justify-between mt-4 pt-4 border-t border-gray-100 bg-white"
              >
                <a-button @click="prevStep">{{
                  $t('curlToRss.common.back')
                }}</a-button>
                <a-space>
                  <a-button
                    type="outline"
                    :loading="parsing"
                    @click="handlePreview"
                  >
                    {{ $t('curlToRss.step2.runPreview') }}
                  </a-button>
                  <a-button
                    type="primary"
                    :disabled="parsedItems.length === 0"
                    @click="nextStep"
                  >
                    {{ $t('curlToRss.step2.nextStep') }}
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
              {{ $t('curlToRss.step3.alert', { count: parsedItems.length }) }}
            </a-alert>

            <a-form :model="feedMeta" layout="vertical">
              <a-form-item :label="$t('curlToRss.step3.feedTitle')" required>
                <a-input
                  v-model="feedMeta.title"
                  :placeholder="$t('curlToRss.placeholder.feedTitle')"
                />
              </a-form-item>
              <a-form-item :label="$t('curlToRss.step3.feedDescription')">
                <a-textarea
                  v-model="feedMeta.description"
                  :placeholder="$t('curlToRss.placeholder.feedDesc')"
                />
              </a-form-item>
              <a-form-item :label="$t('curlToRss.step3.siteLink')">
                <a-input
                  v-model="feedMeta.link"
                  :placeholder="$t('curlToRss.placeholder.siteLink')"
                />
              </a-form-item>
              <a-row :gutter="16">
                <a-col :span="12">
                  <a-form-item :label="$t('curlToRss.step3.authorName')">
                    <a-input v-model="feedMeta.author_name" />
                  </a-form-item>
                </a-col>
                <a-col :span="12">
                  <a-form-item :label="$t('curlToRss.step3.authorEmail')">
                    <a-input v-model="feedMeta.author_email" />
                  </a-form-item>
                </a-col>
              </a-row>
            </a-form>

            <div class="flex justify-between mt-8">
              <a-button @click="prevStep">{{
                $t('curlToRss.common.back')
              }}</a-button>
              <a-button type="primary" @click="handleStep3Next">{{
                $t('curlToRss.common.next')
              }}</a-button>
            </div>
          </div>
        </div>

        <!-- STEP 4: Save -->
        <div v-show="currentStep === 4" class="step-content">
          <div class="max-w-xl mx-auto">
            <a-card
              :title="$t('curlToRss.step4.reviewAndSave')"
              class="border-blue-100"
            >
              <a-descriptions
                :column="1"
                :title="$t('curlToRss.step4.summary')"
                bordered
              >
                <a-descriptions-item :label="$t('curlToRss.step4.sourceUrl')">
                  {{ fetchReq.url }}
                </a-descriptions-item>
                <a-descriptions-item :label="$t('curlToRss.step3.feedTitle')">
                  {{ feedMeta.title }}
                </a-descriptions-item>
                <a-descriptions-item :label="$t('curlToRss.step4.itemCount')">
                  {{ parsedItems.length }}
                </a-descriptions-item>
              </a-descriptions>

              <a-divider />

              <a-form :model="recipeMeta" layout="vertical" class="mt-6">
                <a-form-item
                  :label="$t('curlToRss.step4.recipeId')"
                  required
                  :help="$t('curlToRss.step4.recipeId.help')"
                >
                  <a-input
                    v-model="recipeMeta.id"
                    :placeholder="$t('curlToRss.placeholder.recipeId')"
                  />
                </a-form-item>
                <a-form-item :label="$t('curlToRss.step4.internalDescription')">
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
                    <icon-save /> {{ $t('curlToRss.step4.confirmAndSave') }}
                  </a-button>
                </div>
              </a-form>
            </a-card>

            <div class="flex justify-start mt-8">
              <a-button @click="prevStep">{{
                $t('curlToRss.common.back')
              }}</a-button>
            </div>
          </div>
        </div>
      </a-card>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, reactive, watch } from 'vue';
  import { useRouter } from 'vue-router';
  import { Message, Tree, TreeNodeData } from '@arco-design/web-vue';
  import {
    IconImport,
    IconDelete,
    IconArrowRight,
    IconSave,
    IconEdit,
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
  import { useI18n } from 'vue-i18n';
  import kebabCase from 'lodash/kebabCase';
  import isPlainObject from 'lodash/isPlainObject';
  import isArray from 'lodash/isArray';

  const router = useRouter();
  const { t } = useI18n();

  // --- State ---
  const currentStep = ref(1);
  const fetching = ref(false);
  const parsing = ref(false);
  const parsingCurl = ref(false);
  const saving = ref(false);

  // Step 1 State
  const curlInput = ref('');
  const fetchError = ref('');
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
  const treeData = ref<TreeNodeData[]>([]);
  const activeField = ref<string>(''); // Currently focused input
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

  // --- Watchers & Helpers ---

  const jsonToTree = (data: any, rootPath = ''): TreeNodeData[] => {
    const nodes: TreeNodeData[] = [];
    const getType = (val: any) => {
      if (isArray(val)) return 'array';
      if (isPlainObject(val)) return 'object';
      return 'primitive';
    };

    if (isPlainObject(data)) {
      Object.entries(data).forEach(([key, value]) => {
        const currentPath = rootPath ? `${rootPath}.${key}` : `.${key}`;
        const isObj = isPlainObject(value);
        const isArr = isArray(value);
        const isPrimitive = !isObj && !isArr;

        const node: TreeNodeData = {
          key: currentPath,
          title: isPrimitive ? `${key}: ${value}` : key,
          isLeaf: isPrimitive,
          // store meta data if needed, e.g. type
          data: { type: getType(value) },
        };

        if (!isPrimitive) {
          node.children = jsonToTree(value, currentPath);
        }

        nodes.push(node);
      });
    } else if (isArray(data)) {
      data.forEach((item: any, index: number) => {
        const currentPath = `${rootPath}[${index}]`;
        const isObj = isPlainObject(item);
        const isArr = isArray(item);
        const isPrimitive = !isObj && !isArr;

        const node: TreeNodeData = {
          key: currentPath,
          title: `[${index}]`,
          isLeaf: isPrimitive,
          data: { type: getType(item) },
        };

        if (!isPrimitive) {
          node.children = jsonToTree(item, currentPath);
        } else {
          node.title = `[${index}]: ${item}`;
        }

        nodes.push(node);
      });
    }
    return nodes;
  };

  watch(
    () => jsonContent.value,
    (val) => {
      if (!val) {
        treeData.value = [];
        return;
      }
      try {
        const data = JSON.parse(val);
        treeData.value = jsonToTree(data);
      } catch (e) {
        console.error('Invalid JSON content:', e);
        treeData.value = [];
      }
    },
  );

  const getRelativePath = (fullPath: string, listSel: string) => {
    if (!listSel) return fullPath;
    // Remove trailing [] from list selector to get base path
    const base = listSel.replace(/\[\]$/, '');
    if (fullPath.startsWith(base)) {
      let suffix = fullPath.slice(base.length);
      // Suffix should start with array index, e.g. [0].title or [0]
      // Remove the leading [digits]
      suffix = suffix.replace(/^\[\d+\]/, '');
      if (!suffix) return '.'; // It was the item itself
      return suffix;
    }
    return fullPath;
  };

  const handleNodeSelect = (
    selectedKeys: (string | number)[],
    { node }: { node: TreeNodeData },
  ) => {
    if (!activeField.value || !node.key) return;

    const path = node.key as string;

    if (activeField.value === 'list_selector') {
      // Suggest iterator
      if (node.data && node.data.type === 'array') {
        parseReq.list_selector = `${path}[]`;
      } else {
        parseReq.list_selector = path;
      }
    } else {
      // Relative path calculation
      const rel = getRelativePath(path, parseReq.list_selector);
      // @ts-ignore
      parseReq[activeField.value] = rel;
    }
  };

  // --- Actions ---

  const nextStep = () => {
    currentStep.value += 1;
  };

  const prevStep = () => {
    if (currentStep.value > 1) currentStep.value -= 1;
  };

  const onStepChange = (step: number) => {
    if (step <= currentStep.value) {
      currentStep.value = step;
    }
  };

  watch(
    () => currentStep.value,
    (val) => {
      if (val === 4 && !recipeMeta.id && feedMeta.title) {
        recipeMeta.id = kebabCase(feedMeta.title);
      }
    },
  );

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
      Message.warning(t('curlToRss.msg.enterCurl'));
      return;
    }
    parsingCurl.value = true;
    try {
      const res = await parseCurl(curlInput.value);
      if (res.data) {
        fetchReq.method = res.data.method;
        fetchReq.url = res.data.url;
        fetchReq.headers = res.data.headers || {};
        fetchReq.body = res.data.body || '';
        Message.success(t('curlToRss.msg.curlParsed'));
      }
    } catch (err) {
      // Error handled by interceptor usually
      console.error(err);
    } finally {
      parsingCurl.value = false;
    }
  };

  const handleFetchAndNext = async () => {
    if (!fetchReq.url) {
      Message.warning(t('curlToRss.msg.urlRequired'));
      return;
    }
    fetching.value = true;
    fetchError.value = '';

    try {
      const res = await fetchJson(fetchReq);
      const { data } = res;

      if (!data) {
        const msg = t('curlToRss.msg.emptyResponse');
        fetchError.value = msg;
        Message.warning(msg);
        return;
      }

      // Robust Validation
      let isValid = false;
      if (typeof data === 'string') {
        try {
          JSON.parse(data);
          isValid = true;
        } catch (e) {
          isValid = false;
        }
      } else if (typeof data === 'object') {
        // If data is already an object, it is valid JSON structure
        isValid = true;
      }

      if (!isValid) {
        const msg = t('curlToRss.msg.invalidJson');
        fetchError.value = msg;
        Message.error(msg);
        return;
      }

      // Assign to jsonContent
      // jsonContent expects a string for the watcher to parse it again
      if (typeof data === 'object') {
        jsonContent.value = JSON.stringify(data, null, 2);
      } else {
        jsonContent.value = data;
      }

      // Auto-fill link in meta if possible
      feedMeta.link = fetchReq.url;
      Message.success(t('curlToRss.msg.fetched'));
      nextStep();
    } catch (err: any) {
      console.error(err);
      fetchError.value = err.message || String(err);
    } finally {
      fetching.value = false;
    }
  };

  // Step 2 Logic
  const handlePreview = async () => {
    if (!jsonContent.value) return;
    if (!parseReq.list_selector) {
      Message.warning(t('curlToRss.msg.iteratorRequired'));
      return;
    }
    if (!parseReq.title_selector) {
      Message.warning(t('curlToRss.msg.titleRequired'));
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
        Message.warning(t('curlToRss.msg.noItems'));
      } else {
        Message.success(
          t('curlToRss.msg.parsedItems', { count: parsedItems.value.length }),
        );
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
      Message.warning(t('curlToRss.msg.feedTitleRequired'));
      return;
    }
    nextStep();
  };

  // Step 4 Logic
  const handleSaveRecipe = async () => {
    if (!recipeMeta.id) {
      Message.warning(t('curlToRss.msg.recipeIdRequired'));
      return;
    }

    saving.value = true;

    // Construct the SourceConfig object conforming to internal/config/source_config.go
    const sourceConfig = {
      type: 'json',
      http_fetcher: {
        url: fetchReq.url,
        method: fetchReq.method,
        headers: fetchReq.headers,
        body: fetchReq.body,
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

      Message.success(t('curlToRss.msg.saved'));
      router.push({ name: 'CustomRecipe' });
    } catch (err: any) {
      Message.error(t('curlToRss.msg.saveFailed', { msg: err.message || err }));
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
    min-height: 600px;
    /* Ensure content takes available space */
    display: flex;
    flex-direction: column;
  }

  .step-fixed-height {
    height: 600px;
  }

  .header-row {
    margin-bottom: 8px;
    display: flex;
    gap: 8px;
    align-items: center;
  }
</style>
