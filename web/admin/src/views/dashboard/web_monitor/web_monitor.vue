<template>
  <div class="py-8 px-16">
    <x-header
      :title="t('webMonitor.header.title')"
      :description="t('webMonitor.header.description')"
    />

    <div class="content-wrapper">
      <a-card class="wizard-card">
        <a-steps
          :current="currentStep"
          changeable
          class="mb-8"
          @change="onStepChange"
        >
          <a-step
            :title="t('webMonitor.step.targetUrl')"
            :description="t('webMonitor.step.targetUrl.desc')"
          />
          <a-step
            :title="t('webMonitor.step.extractRules')"
            :description="t('webMonitor.step.extractRules.desc')"
          />
          <a-step
            :title="t('webMonitor.step.feedMetadata')"
            :description="t('webMonitor.step.feedMetadata.desc')"
          />
          <a-step
            :title="t('webMonitor.step.saveRecipe')"
            :description="t('webMonitor.step.saveRecipe.desc')"
          />
        </a-steps>

        <div v-show="currentStep === 1" class="step-content">
          <a-form layout="vertical" class="max-w-xl mx-auto">
            <a-form-item
              :label="t('webMonitor.step1.label')"
              :help="t('webMonitor.step1.help')"
            >
              <a-input
                v-model="url"
                :placeholder="t('webMonitor.step1.placeholder')"
                size="large"
                allow-clear
                @keyup.enter="fetchAndNext"
                @input="fetchError = ''"
              />
            </a-form-item>

            <div class="flex items-center gap-2 mb-6 ml-1">
              <span class="text-gray-600">
                {{ t('webMonitor.step2.enhanceMode') }}
              </span>
              <a-tooltip :content="t('webMonitor.step2.enhanceMode.tooltip')">
                <a-switch v-model="enhancedMode" />
              </a-tooltip>
            </div>

            <a-alert v-if="fetchError" type="error" class="mb-4" show-icon>
              {{ fetchError }}
            </a-alert>
            <div class="text-center mt-8">
              <a-button
                type="primary"
                size="large"
                :loading="fetching"
                :disabled="!url"
                @click="fetchAndNext"
              >
                {{ t('webMonitor.step1.button') }} <icon-arrow-right />
              </a-button>
            </div>
          </a-form>
        </div>

        <div v-show="currentStep === 2" class="step-content h-full">
          <a-row :gutter="16" class="h-full">
            <a-col :span="14" class="h-full flex flex-col">
              <div class="flex justify-between items-center mb-2">
                <span class="font-bold">
                  {{ t('webMonitor.step2.pagePreview') }}
                </span>
                <a-tag v-if="isSelectionMode" color="blue">
                  {{ t('webMonitor.step2.selectionModeOn') }}
                </a-tag>
              </div>
              <a-spin :loading="fetching" class="h-full flex-1 flex flex-col">
                <HtmlPreview
                  ref="previewRef"
                  class="flex-1"
                  :html-content="htmlContent"
                  :is-selection-mode="isSelectionMode"
                  @select="handleElementSelect"
                />
              </a-spin>
            </a-col>

            <a-col :span="10" class="h-full flex flex-col">
              <div class="flex-1 overflow-y-auto pr-2">
                <a-alert type="info" class="mb-4">
                  {{ t('webMonitor.step2.alert.l1') }}<br />
                  {{ t('webMonitor.step2.alert.l2') }}
                </a-alert>

                <a-card
                  :title="t('webMonitor.step2.variables')"
                  size="small"
                  class="mb-4 border-blue-100"
                >
                  <div class="mb-3 text-xs text-gray-500">
                    {{
                      t('webMonitor.step2.availableVars', {
                        vars: availableVariables,
                      })
                    }}
                  </div>
                  <div
                    v-for="(field, index) in fields"
                    :key="field.id"
                    class="field-row mb-3"
                  >
                    <a-space direction="vertical" fill>
                      <a-input
                        v-model="field.name"
                        :placeholder="t('webMonitor.step2.variableName')"
                      />
                      <a-input
                        v-model="field.selector"
                        :placeholder="t('webMonitor.step2.selector')"
                      >
                        <template #suffix>
                          <a-button
                            size="mini"
                            :type="
                              currentFieldId === field.id
                                ? 'primary'
                                : 'secondary'
                            "
                            :status="
                              currentFieldId === field.id ? 'warning' : 'normal'
                            "
                            @click="setTargetField(field.id)"
                          >
                            {{
                              currentFieldId === field.id
                                ? t('webMonitor.step2.picking')
                                : t('webMonitor.step2.pick')
                            }}
                          </a-button>
                        </template>
                      </a-input>
                      <div class="flex items-center justify-between">
                        <a-checkbox v-model="field.isKey">
                          {{ t('webMonitor.step2.keyField') }}
                        </a-checkbox>
                        <a-button
                          type="text"
                          status="danger"
                          size="mini"
                          @click="removeField(index)"
                        >
                          {{ t('webMonitor.step2.remove') }}
                        </a-button>
                      </div>
                    </a-space>
                  </div>
                  <a-button type="outline" long @click="addField">
                    {{ t('webMonitor.step2.addVariable') }}
                  </a-button>
                </a-card>

                <a-card size="small" class="mb-4 border-purple-100">
                  <template #title>
                    <div class="flex items-center gap-2">
                      <span>{{ t('webMonitor.step2.aiJudge') }}</span>
                      <a-tooltip
                        :content="t('webMonitor.step2.aiJudge.tooltip')"
                      >
                        <icon-info-circle class="text-gray-400" />
                      </a-tooltip>
                    </div>
                  </template>
                  <template #extra>
                    <a-switch v-model="aiJudge.enabled" size="small" />
                  </template>
                  <a-form v-show="aiJudge.enabled" layout="vertical">
                    <a-form-item :label="t('webMonitor.step2.aiJudge.prompt')">
                      <a-textarea
                        v-model="aiJudge.prompt"
                        :placeholder="
                          t('webMonitor.step2.aiJudge.prompt.placeholder')
                        "
                        :auto-size="{ minRows: 2, maxRows: 4 }"
                        allow-clear
                      />
                    </a-form-item>
                    <a-row :gutter="12">
                      <a-col :span="12">
                        <a-form-item
                          :label="t('webMonitor.step2.aiJudge.outputField')"
                          :help="t('webMonitor.step2.aiJudge.outputField.help')"
                        >
                          <a-input
                            v-model="aiJudge.outputField"
                            placeholder="ai_verdict"
                            allow-clear
                          />
                        </a-form-item>
                      </a-col>
                      <a-col :span="12">
                        <a-form-item
                          :label="t('webMonitor.step2.aiJudge.model')"
                        >
                          <a-input
                            v-model="aiJudge.model"
                            :placeholder="
                              t('webMonitor.step2.aiJudge.model.placeholder')
                            "
                            allow-clear
                          />
                        </a-form-item>
                      </a-col>
                    </a-row>
                    <a-checkbox v-model="aiJudge.useAsKey">
                      {{ t('webMonitor.step2.aiJudge.useAsKey') }}
                    </a-checkbox>
                  </a-form>
                </a-card>

                <a-card :title="t('webMonitor.step2.preview')" size="small">
                  <a-form layout="vertical">
                    <a-form-item :label="t('webMonitor.step2.previewTitle')">
                      <a-input v-model="templates.title" allow-clear />
                    </a-form-item>
                    <a-form-item
                      :label="t('webMonitor.step2.previewDescription')"
                    >
                      <a-textarea
                        v-model="templates.description"
                        :auto-size="{ minRows: 2, maxRows: 4 }"
                        allow-clear
                      />
                    </a-form-item>
                    <a-form-item :label="t('webMonitor.step2.previewContent')">
                      <a-textarea
                        v-model="templates.content"
                        :auto-size="{ minRows: 4, maxRows: 8 }"
                        allow-clear
                      />
                    </a-form-item>
                  </a-form>

                  <template
                    v-if="
                      preview.title || preview.description || preview.content
                    "
                  >
                    <a-divider />
                    <div v-if="aiVerdict" class="mb-2">
                      <div class="text-xs text-gray-500 mb-1">
                        {{ t('webMonitor.step2.aiJudge.verdict') }}
                      </div>
                      <a-tag color="purple">{{ aiVerdict }}</a-tag>
                    </div>
                    <div class="mb-2">
                      <div class="text-xs text-gray-500 mb-1">
                        {{ t('webMonitor.step2.previewResultTitle') }}
                      </div>
                      <div>{{ preview.title }}</div>
                    </div>
                    <div class="mb-2">
                      <div class="text-xs text-gray-500 mb-1">
                        {{ t('webMonitor.step2.previewResultDescription') }}
                      </div>
                      <div>{{ preview.description }}</div>
                    </div>
                    <div class="mb-2">
                      <div class="text-xs text-gray-500 mb-1">GUID</div>
                      <div class="break-all">{{ preview.guid }}</div>
                    </div>
                    <div>
                      <div class="text-xs text-gray-500 mb-1">
                        {{ t('webMonitor.step2.previewResultContent') }}
                      </div>
                      <pre class="preview-content">{{ preview.content }}</pre>
                    </div>
                  </template>
                </a-card>
              </div>

              <div
                class="flex justify-between mt-4 pt-4 border-t border-gray-100 bg-white"
              >
                <a-button @click="prevStep">
                  {{ t('webMonitor.common.back') }}
                </a-button>
                <a-space>
                  <a-button
                    type="outline"
                    :loading="parsing"
                    @click="runPreview"
                  >
                    {{ t('webMonitor.step2.runPreview') }}
                  </a-button>
                  <a-button
                    type="primary"
                    :disabled="!canProceedFromStep2"
                    @click="nextStep"
                  >
                    {{ t('webMonitor.step2.nextStep') }}
                  </a-button>
                </a-space>
              </div>
            </a-col>
          </a-row>
        </div>

        <div v-show="currentStep === 3" class="step-content">
          <div class="max-w-2xl mx-auto">
            <a-form :model="feedMeta" layout="vertical">
              <a-form-item :label="t('webMonitor.step3.feedTitle')" required>
                <a-input
                  v-model="feedMeta.title"
                  :placeholder="t('webMonitor.step3.feedTitle.placeholder')"
                  allow-clear
                />
              </a-form-item>
              <a-form-item :label="t('webMonitor.step3.feedDesc')">
                <a-textarea
                  v-model="feedMeta.description"
                  :placeholder="t('webMonitor.step3.feedDesc.placeholder')"
                  allow-clear
                />
              </a-form-item>
              <a-form-item :label="t('webMonitor.step3.siteLink')">
                <a-input
                  v-model="feedMeta.link"
                  :placeholder="t('webMonitor.step3.siteLink.placeholder')"
                  allow-clear
                />
              </a-form-item>
              <a-row :gutter="16">
                <a-col :span="12">
                  <a-form-item :label="t('webMonitor.step3.authorName')">
                    <a-input v-model="feedMeta.author_name" allow-clear />
                  </a-form-item>
                </a-col>
                <a-col :span="12">
                  <a-form-item :label="t('webMonitor.step3.authorEmail')">
                    <a-input v-model="feedMeta.author_email" allow-clear />
                  </a-form-item>
                </a-col>
              </a-row>
            </a-form>

            <div class="flex justify-between mt-8">
              <a-button @click="prevStep">
                {{ t('webMonitor.common.back') }}
              </a-button>
              <a-button type="primary" @click="nextStep">
                {{ t('webMonitor.common.next') }}
              </a-button>
            </div>
          </div>
        </div>

        <div v-show="currentStep === 4" class="step-content">
          <div class="max-w-xl mx-auto">
            <a-card
              :title="t('webMonitor.step4.card.title')"
              class="border-blue-100"
            >
              <a-descriptions
                :column="1"
                :title="t('webMonitor.step4.summary')"
                bordered
              >
                <a-descriptions-item :label="t('webMonitor.step4.sourceUrl')">
                  {{ url }}
                </a-descriptions-item>
                <a-descriptions-item :label="t('webMonitor.step4.feedTitle')">
                  {{ feedMeta.title }}
                </a-descriptions-item>
                <a-descriptions-item :label="t('webMonitor.step4.keyFields')">
                  {{ selectedKeyFields.join(', ') }}
                </a-descriptions-item>
              </a-descriptions>

              <a-divider />

              <a-form :model="recipeMeta" layout="vertical" class="mt-6">
                <a-form-item
                  :label="t('webMonitor.step4.recipeId')"
                  required
                  field="id"
                  :rules="getRecipeIdRules(t('webMonitor.msg.idRequired'))"
                  :help="t('webMonitor.step4.recipeId.help')"
                >
                  <a-input
                    v-model="recipeMeta.id"
                    :placeholder="t('webMonitor.step4.recipeId.placeholder')"
                    allow-clear
                  >
                    <template #append>
                      <a-tooltip content="Generate ID from Title">
                        <a-button
                          @click="
                            recipeMeta.id = generateRecipeId(feedMeta.title)
                          "
                        >
                          <template #icon><icon-refresh /></template>
                        </a-button>
                      </a-tooltip>
                    </template>
                  </a-input>
                </a-form-item>
                <a-form-item :label="t('webMonitor.step4.internalDesc')">
                  <a-textarea
                    v-model="recipeMeta.description"
                    :placeholder="
                      t('webMonitor.step4.internalDesc.placeholder')
                    "
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
                    <icon-save /> {{ t('webMonitor.step4.confirmSave') }}
                  </a-button>
                </div>
              </a-form>
            </a-card>

            <div class="flex justify-start mt-8">
              <a-button @click="prevStep">
                {{ t('webMonitor.common.back') }}
              </a-button>
            </div>
          </div>
        </div>
      </a-card>
    </div>
  </div>
</template>

<script lang="ts" setup>
  import { computed, reactive, ref, watch } from 'vue';
  import axios from 'axios';
  import DOMPurify from 'dompurify';
  import { Message } from '@arco-design/web-vue';
  import {
    IconArrowRight,
    IconInfoCircle,
    IconRefresh,
    IconSave,
  } from '@arco-design/web-vue/es/icon';
  import { useRouter } from 'vue-router';
  import { useI18n } from 'vue-i18n';
  import XHeader from '@/components/header/x-header.vue';
  import HtmlPreview from '@/views/dashboard/html_to_rss/components/HtmlPreview.vue';
  import { createCustomRecipe } from '@/api/custom_recipe';
  import generateRecipeId, { getRecipeIdRules } from '@/utils/slug';
  import { getCssSelector } from '@/views/dashboard/html_to_rss/utils/selector';

  const router = useRouter();
  const { t } = useI18n();

  type MonitorField = {
    id: string;
    name: string;
    selector: string;
    isKey: boolean;
  };

  type WebMonitorPreview = {
    values: Record<string, string>;
    key_fields: string[];
    guid: string;
    title: string;
    description: string;
    content: string;
    feed_title: string;
    feed_link: string;
  };

  const currentStep = ref(1);
  const url = ref('');
  const enhancedMode = ref(false);
  const fetching = ref(false);
  const fetchError = ref('');
  const parsing = ref(false);
  const saving = ref(false);
  const htmlContent = ref('');
  const currentFieldId = ref('');
  const previewRef = ref<InstanceType<typeof HtmlPreview> | null>(null);

  const fields = ref<MonitorField[]>([
    { id: crypto.randomUUID(), name: 'price', selector: '', isKey: true },
    { id: crypto.randomUUID(), name: 'stock', selector: '', isKey: false },
  ]);

  const templates = reactive({
    title: '【监控更新】{{.price}}',
    description: '价格 {{.price}}，状态 {{.stock}}',
    content: '价格：{{.price}}\n状态：{{.stock}}\n链接：{{.url}}',
  });

  const aiJudge = reactive({
    enabled: false,
    prompt: '',
    outputField: 'ai_verdict',
    model: '',
    useAsKey: true,
  });

  const aiJudgeFieldName = computed(
    () => aiJudge.outputField.trim() || 'ai_verdict'
  );

  const preview = reactive<WebMonitorPreview>({
    values: {},
    key_fields: [],
    guid: '',
    title: '',
    description: '',
    content: '',
    feed_title: '',
    feed_link: '',
  });

  const feedMeta = reactive({
    title: '',
    link: '',
    description: '',
    author_name: '',
    author_email: '',
  });

  const recipeMeta = reactive({
    id: '',
    description: '',
  });

  const isSelectionMode = computed(() => Boolean(currentFieldId.value));
  const availableVariables = computed(() => {
    const vars = [
      'url',
      ...fields.value.map((field) => field.name).filter(Boolean),
    ];
    if (aiJudge.enabled) {
      vars.push(aiJudgeFieldName.value);
    }
    return vars.map((item) => `{{.${item}}}`).join(', ');
  });
  const selectedKeyFields = computed(() => {
    const keys = fields.value
      .filter((field) => field.isKey && field.name.trim())
      .map((field) => field.name.trim());
    if (aiJudge.enabled && aiJudge.useAsKey) {
      keys.push(aiJudgeFieldName.value);
    }
    return keys;
  });
  const aiVerdict = computed(
    () => preview.values?.[aiJudgeFieldName.value] || ''
  );
  const canProceedFromStep2 = computed(
    () =>
      fields.value.some(
        (field) => field.name.trim() && field.selector.trim()
      ) && selectedKeyFields.value.length > 0
  );

  const addField = () => {
    fields.value.push({
      id: crypto.randomUUID(),
      name: `${t('webMonitor.step2.defaultVarName')}${fields.value.length + 1}`,
      selector: '',
      isKey: false,
    });
  };

  const removeField = (index: number) => {
    if (fields.value.length === 1) return;
    const removed = fields.value[index];
    fields.value = fields.value.filter((_, i) => i !== index);
    if (currentFieldId.value === removed.id) {
      currentFieldId.value = '';
    }
  };

  const nextStep = () => {
    currentStep.value += 1;
  };

  const prevStep = () => {
    if (currentStep.value > 1) currentStep.value -= 1;
  };

  const onStepChange = (step: number) => {
    if (step <= currentStep.value) currentStep.value = step;
  };

  watch(
    () => currentStep.value,
    (val) => {
      if (val === 4 && !recipeMeta.id && feedMeta.title) {
        recipeMeta.id = generateRecipeId(feedMeta.title);
      }
    }
  );

  const setTargetField = (fieldId: string) => {
    currentFieldId.value = fieldId;
    const field = fields.value.find((item) => item.id === fieldId);
    Message.info(
      t('webMonitor.msg.pickInfo', { field: field?.name || fieldId })
    );
  };

  const fetchContent = async (advanceStep = false) => {
    if (!url.value) return;
    fetching.value = true;
    fetchError.value = '';

    try {
      const { data: res } = (await axios.post('/api/admin/tools/fetch', {
        url: url.value,
        use_browserless: enhancedMode.value,
      })) as any;
      if (res.code === 0) {
        let raw = res.data;
        const baseTag = `<base href="${url.value}" />`;
        if (raw.toLowerCase().includes('<head>')) {
          raw = raw.replace(/<head>/i, `<head>${baseTag}`);
        } else {
          raw = `${baseTag}${raw}`;
        }

        try {
          const doc = new DOMParser().parseFromString(raw, 'text/html');
          const title = doc.querySelector('title')?.innerText || '';
          const descMeta =
            doc.querySelector('meta[name="description"]') ||
            doc.querySelector('meta[property="og:description"]');
          const description = descMeta ? descMeta.getAttribute('content') : '';

          if (title) feedMeta.title = title.trim();
          if (description) feedMeta.description = description.trim();
        } catch {
          // ignore
        }

        htmlContent.value = DOMPurify.sanitize(raw, {
          WHOLE_DOCUMENT: true,
          ADD_TAGS: ['link', 'style', 'head', 'meta', 'body', 'html', 'base'],
          ADD_ATTR: [
            'href',
            'rel',
            'src',
            'type',
            'class',
            'id',
            'style',
            'title',
            'alt',
            'target',
            'width',
            'height',
          ],
        });

        feedMeta.link = url.value;
        if (advanceStep) nextStep();
      } else {
        fetchError.value = res.msg || t('webMonitor.msg.fetchFailed');
      }
    } catch (err: any) {
      fetchError.value = err.message || t('webMonitor.msg.errorFetching');
    } finally {
      fetching.value = false;
    }
  };

  const fetchAndNext = async () => {
    await fetchContent(true);
  };

  const validateFields = () => {
    const names = new Set<string>();

    const invalidField = fields.value.find((field) => {
      const name = field.name.trim();
      const selector = field.selector.trim();

      if (!name) {
        Message.warning(t('webMonitor.msg.variableNameRequired'));
        return true;
      }
      if (!selector) {
        Message.warning(t('webMonitor.msg.selectorRequired'));
        return true;
      }
      if (names.has(name)) {
        Message.warning(t('webMonitor.msg.duplicateVariable'));
        return true;
      }

      names.add(name);
      return false;
    });

    if (invalidField) {
      return false;
    }

    if (aiJudge.enabled && !aiJudge.prompt.trim()) {
      Message.warning(t('webMonitor.msg.aiPromptRequired'));
      return false;
    }

    if (selectedKeyFields.value.length === 0) {
      Message.warning(t('webMonitor.msg.keyRequired'));
      return false;
    }
    return true;
  };

  const buildParserConfig = () => {
    const extractors = Object.fromEntries(
      fields.value.map((field) => [field.name.trim(), field.selector.trim()])
    );
    const parser: Record<string, any> = {
      extractors,
      key_fields: selectedKeyFields.value,
      title_template: templates.title,
      description_template: templates.description,
      content_template: templates.content,
    };
    if (aiJudge.enabled) {
      parser.ai_judge = {
        enabled: true,
        prompt: aiJudge.prompt.trim(),
        output_field: aiJudgeFieldName.value,
        ...(aiJudge.model.trim() ? { model: aiJudge.model.trim() } : {}),
      };
    }
    return parser;
  };

  const runPreview = async () => {
    if (!validateFields()) return;
    parsing.value = true;
    try {
      const { data: res } = (await axios.post(
        '/api/admin/tools/web-monitor/preview',
        {
          html: htmlContent.value,
          url: url.value,
          use_browserless: enhancedMode.value,
          web_monitor_parser: buildParserConfig(),
        }
      )) as any;

      if (res.code === 0) {
        Object.assign(preview, res.data);
        Message.success(t('webMonitor.msg.previewSuccess'));
      } else {
        Message.error(res.msg || t('webMonitor.msg.previewFailed'));
      }
    } catch (err: any) {
      Message.error(err.message || t('webMonitor.msg.errorPreview'));
    } finally {
      parsing.value = false;
    }
  };

  const handleSaveRecipe = async () => {
    if (!recipeMeta.id) {
      Message.warning(t('webMonitor.msg.idRequired'));
      return;
    }
    if (!validateFields()) return;

    saving.value = true;

    const sourceConfig = {
      type: 'web_monitor',
      http_fetcher: {
        url: url.value,
        use_browserless: enhancedMode.value,
      },
      web_monitor_parser: buildParserConfig(),
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
        source_type: 'web_monitor',
        source_config: JSON.stringify(sourceConfig),
      });
      Message.success(t('webMonitor.msg.saveSuccess'));
      router.push({ name: 'CustomRecipe' });
    } catch (err: any) {
      Message.error(
        t('webMonitor.msg.saveFailed', { msg: err.message || err })
      );
    } finally {
      saving.value = false;
    }
  };

  const handleElementSelect = (target: HTMLElement) => {
    if (!currentFieldId.value) {
      Message.warning(t('webMonitor.msg.clickPickFirst'));
      return;
    }
    const doc = previewRef.value?.contentDocument;
    if (!doc) return;
    const selector = getCssSelector(target, doc, false);
    const field = fields.value.find((item) => item.id === currentFieldId.value);
    if (!field) return;
    field.selector = selector;
    currentFieldId.value = '';
  };
</script>

<style scoped>
  .wizard-card {
    min-height: 700px;
  }

  .step-content {
    margin-top: 24px;
    height: 600px;
  }

  .preview-content {
    white-space: pre-wrap;
    background: #f7f8fa;
    padding: 12px;
    border-radius: 4px;
  }
</style>
