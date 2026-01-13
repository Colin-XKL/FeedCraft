<template>
  <div class="py-8 px-16">
    <x-header title="HTML to RSS" description="HTML to RSS"></x-header>

    <div class="content-wrapper">
      <a-card class="wizard-card">
        <a-steps :current="currentStep" class="mb-8">
          <a-step
            :title="$t('htmlToRss.step.targetUrl')"
            :description="$t('htmlToRss.step.targetUrl.desc')"
          />
          <a-step
            :title="$t('htmlToRss.step.extractRules')"
            :description="$t('htmlToRss.step.extractRules.desc')"
          />
          <a-step
            :title="$t('htmlToRss.step.feedMetadata')"
            :description="$t('htmlToRss.step.feedMetadata.desc')"
          />
          <a-step
            :title="$t('htmlToRss.step.saveRecipe')"
            :description="$t('htmlToRss.step.saveRecipe.desc')"
          />
        </a-steps>

        <!-- STEP 1: URL Input -->
        <div v-show="currentStep === 1" class="step-content">
          <a-form layout="vertical" class="max-w-xl mx-auto">
            <a-form-item
              :label="$t('htmlToRss.step1.label')"
              :help="$t('htmlToRss.step1.help')"
            >
              <a-input
                v-model="url"
                :placeholder="$t('htmlToRss.step1.placeholder')"
                size="large"
                allow-clear
                @keyup.enter="fetchAndNext"
                @input="fetchError = ''"
              />
            </a-form-item>

            <div class="flex items-center gap-2 mb-6 ml-1">
              <span class="text-gray-600">{{
                $t('htmlToRss.step2.enhanceMode')
              }}</span>
              <a-tooltip :content="$t('htmlToRss.step2.enhanceMode.tooltip')">
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
                {{ $t('htmlToRss.step1.button') }} <icon-arrow-right />
              </a-button>
            </div>
          </a-form>
        </div>

        <!-- STEP 2: Extraction Rules -->
        <div v-show="currentStep === 2" class="step-content h-full">
          <a-row :gutter="16" class="h-full">
            <!-- Left: Preview Component -->
            <a-col :span="14" class="h-full flex flex-col">
              <div class="flex justify-between items-center mb-2">
                <div class="flex items-center gap-2">
                  <span class="font-bold">{{
                    $t('htmlToRss.step2.pagePreview')
                  }}</span>
                </div>
                <a-tag v-if="isSelectionMode" color="blue">{{
                  $t('htmlToRss.step2.selectionModeOn')
                }}</a-tag>
              </div>

              <!-- Use the extracted HtmlPreview component -->
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

            <!-- Right: Config & Preview -->
            <a-col :span="10" class="h-full flex flex-col">
              <div class="flex-1 overflow-y-auto pr-2">
                <a-alert type="info" class="mb-4">
                  {{ $t('htmlToRss.step2.alert.l1') }}<br />
                  {{ $t('htmlToRss.step2.alert.l2') }}
                </a-alert>

                <a-form :model="config" layout="vertical">
                  <a-card
                    :title="$t('htmlToRss.step2.card1.title')"
                    size="small"
                    class="mb-4 border-blue-100"
                  >
                    <a-form-item :label="$t('htmlToRss.step2.cssSelector')">
                      <a-input
                        v-model="config.item_selector"
                        placeholder=".article-card"
                        allow-clear
                      >
                        <template #suffix>
                          <a-button
                            size="mini"
                            :type="
                              currentTargetField === 'item_selector'
                                ? 'primary'
                                : 'primary'
                            "
                            :status="
                              currentTargetField === 'item_selector'
                                ? 'warning'
                                : 'success'
                            "
                            @click="setTargetField('item_selector')"
                          >
                            <icon-select-all />
                            {{
                              currentTargetField === 'item_selector'
                                ? $t('htmlToRss.step2.picking')
                                : $t('htmlToRss.step2.pick')
                            }}
                          </a-button>
                        </template>
                      </a-input>
                    </a-form-item>
                  </a-card>

                  <a-card
                    :title="$t('htmlToRss.step2.card2.title')"
                    size="small"
                    :class="{ 'opacity-50': !config.item_selector }"
                  >
                    <a-form-item :label="$t('htmlToRss.step2.title')">
                      <a-input
                        v-model="config.title_selector"
                        :disabled="!config.item_selector"
                        allow-clear
                      >
                        <template #suffix>
                          <a-button
                            size="mini"
                            :disabled="!config.item_selector"
                            :type="
                              currentTargetField === 'title_selector'
                                ? 'primary'
                                : 'secondary'
                            "
                            :status="
                              currentTargetField === 'title_selector'
                                ? 'warning'
                                : 'normal'
                            "
                            @click="setTargetField('title_selector')"
                          >
                            {{
                              currentTargetField === 'title_selector'
                                ? $t('htmlToRss.step2.picking')
                                : $t('htmlToRss.step2.pick')
                            }}
                          </a-button>
                        </template>
                      </a-input>
                    </a-form-item>
                    <a-form-item :label="$t('htmlToRss.step2.link')">
                      <a-input
                        v-model="config.link_selector"
                        :disabled="!config.item_selector"
                        allow-clear
                      >
                        <template #suffix>
                          <a-button
                            size="mini"
                            :disabled="!config.item_selector"
                            :type="
                              currentTargetField === 'link_selector'
                                ? 'primary'
                                : 'secondary'
                            "
                            :status="
                              currentTargetField === 'link_selector'
                                ? 'warning'
                                : 'normal'
                            "
                            @click="setTargetField('link_selector')"
                          >
                            {{
                              currentTargetField === 'link_selector'
                                ? $t('htmlToRss.step2.picking')
                                : $t('htmlToRss.step2.pick')
                            }}
                          </a-button>
                        </template>
                      </a-input>
                    </a-form-item>
                    <a-form-item :label="$t('htmlToRss.step2.date')">
                      <a-input
                        v-model="config.date_selector"
                        :disabled="!config.item_selector"
                        allow-clear
                      >
                        <template #suffix>
                          <a-button
                            size="mini"
                            :disabled="!config.item_selector"
                            :type="
                              currentTargetField === 'date_selector'
                                ? 'primary'
                                : 'secondary'
                            "
                            :status="
                              currentTargetField === 'date_selector'
                                ? 'warning'
                                : 'normal'
                            "
                            @click="setTargetField('date_selector')"
                          >
                            {{
                              currentTargetField === 'date_selector'
                                ? $t('htmlToRss.step2.picking')
                                : $t('htmlToRss.step2.pick')
                            }}
                          </a-button>
                        </template>
                      </a-input>
                    </a-form-item>
                    <a-form-item :label="$t('htmlToRss.step2.description')">
                      <a-input
                        v-model="config.description_selector"
                        :disabled="!config.item_selector"
                        allow-clear
                      >
                        <template #suffix>
                          <a-button
                            size="mini"
                            :disabled="!config.item_selector"
                            :type="
                              currentTargetField === 'description_selector'
                                ? 'primary'
                                : 'secondary'
                            "
                            :status="
                              currentTargetField === 'description_selector'
                                ? 'warning'
                                : 'normal'
                            "
                            @click="setTargetField('description_selector')"
                          >
                            {{
                              currentTargetField === 'description_selector'
                                ? $t('htmlToRss.step2.picking')
                                : $t('htmlToRss.step2.pick')
                            }}
                          </a-button>
                        </template>
                      </a-input>
                    </a-form-item>
                  </a-card>
                </a-form>

                <!-- Immediate Preview Results -->
                <div
                  v-if="parsedItems.length > 0"
                  ref="resultsRef"
                  class="mt-4"
                >
                  <a-divider orientation="left"
                    >{{ $t('htmlToRss.step2.previewResults') }} ({{
                      parsedItems.length
                    }})</a-divider
                  >
                  <a-collapse :default-active-key="[0]">
                    <a-collapse-item
                      v-for="(item, idx) in parsedItems"
                      :key="idx"
                      :header="item.title || $t('htmlToRss.step2.noTitle')"
                    >
                      <div class="text-xs text-gray-500 mb-1">{{
                        item.link
                      }}</div>
                      <div class="text-xs text-gray-400 mb-2">{{
                        item.date
                      }}</div>
                      <div class="text-sm text-gray-600 truncate">{{
                        item.content || item.description
                      }}</div>
                    </a-collapse-item>
                  </a-collapse>
                </div>
              </div>

              <!-- Actions Footer -->
              <div
                class="flex justify-between mt-4 pt-4 border-t border-gray-100 bg-white"
              >
                <a-button @click="prevStep">{{
                  $t('htmlToRss.common.back')
                }}</a-button>
                <a-space>
                  <a-button
                    type="outline"
                    :loading="parsing"
                    :disabled="!config.item_selector"
                    @click="runPreview"
                  >
                    {{ $t('htmlToRss.step2.runPreview') }}
                  </a-button>
                  <a-button
                    type="primary"
                    :disabled="parsedItems.length === 0"
                    @click="nextStep"
                  >
                    {{ $t('htmlToRss.step2.nextStep') }}
                  </a-button>
                </a-space>
              </div>
            </a-col>
          </a-row>
        </div>

        <!-- STEP 3: Feed Metadata -->
        <div v-show="currentStep === 3" class="step-content">
          <div class="max-w-2xl mx-auto">
            <a-alert
              :title="$t('htmlToRss.step3.alert.title')"
              type="success"
              class="mb-6"
            >
              {{
                $t('htmlToRss.step3.alert.desc', {
                  count: parsedItems.length,
                })
              }}
            </a-alert>

            <a-form :model="feedMeta" layout="vertical">
              <a-form-item :label="$t('htmlToRss.step3.feedTitle')" required>
                <a-input
                  v-model="feedMeta.title"
                  :placeholder="$t('htmlToRss.step3.feedTitle.placeholder')"
                  allow-clear
                />
              </a-form-item>
              <a-form-item :label="$t('htmlToRss.step3.feedDesc')">
                <a-textarea
                  v-model="feedMeta.description"
                  :placeholder="$t('htmlToRss.step3.feedDesc.placeholder')"
                  allow-clear
                />
              </a-form-item>
              <a-form-item :label="$t('htmlToRss.step3.siteLink')">
                <a-input
                  v-model="feedMeta.link"
                  :placeholder="$t('htmlToRss.step3.siteLink.placeholder')"
                  allow-clear
                />
              </a-form-item>
              <a-row :gutter="16">
                <a-col :span="12">
                  <a-form-item :label="$t('htmlToRss.step3.authorName')">
                    <a-input v-model="feedMeta.author_name" allow-clear />
                  </a-form-item>
                </a-col>
                <a-col :span="12">
                  <a-form-item :label="$t('htmlToRss.step3.authorEmail')">
                    <a-input v-model="feedMeta.author_email" allow-clear />
                  </a-form-item>
                </a-col>
              </a-row>
            </a-form>

            <div class="flex justify-between mt-8">
              <a-button @click="prevStep">{{
                $t('htmlToRss.common.back')
              }}</a-button>
              <a-button type="primary" @click="nextStep">{{
                $t('htmlToRss.common.next')
              }}</a-button>
            </div>
          </div>
        </div>

        <!-- STEP 4: Save -->
        <div v-show="currentStep === 4" class="step-content">
          <div class="max-w-xl mx-auto">
            <a-card
              :title="$t('htmlToRss.step4.card.title')"
              class="border-blue-100"
            >
              <a-descriptions
                :column="1"
                :title="$t('htmlToRss.step4.summary')"
                bordered
              >
                <a-descriptions-item :label="$t('htmlToRss.step4.sourceUrl')">{{
                  url
                }}</a-descriptions-item>
                <a-descriptions-item :label="$t('htmlToRss.step4.feedTitle')">{{
                  feedMeta.title
                }}</a-descriptions-item>
                <a-descriptions-item :label="$t('htmlToRss.step4.itemCount')">{{
                  $t('htmlToRss.step4.itemCount.value', {
                    count: parsedItems.length,
                  })
                }}</a-descriptions-item>
              </a-descriptions>

              <a-divider />

              <a-form :model="recipeMeta" layout="vertical" class="mt-6">
                <a-form-item
                  :label="$t('htmlToRss.step4.recipeId')"
                  required
                  :help="$t('htmlToRss.step4.recipeId.help')"
                >
                  <a-input
                    v-model="recipeMeta.id"
                    :placeholder="$t('htmlToRss.step4.recipeId.placeholder')"
                  />
                </a-form-item>
                <a-form-item :label="$t('htmlToRss.step4.internalDesc')">
                  <a-textarea
                    v-model="recipeMeta.description"
                    :placeholder="
                      $t('htmlToRss.step4.internalDesc.placeholder')
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
                    <icon-save /> {{ $t('htmlToRss.step4.confirmSave') }}
                  </a-button>
                </div>
              </a-form>
            </a-card>

            <div class="flex justify-start mt-8">
              <a-button @click="prevStep">{{
                $t('htmlToRss.common.back')
              }}</a-button>
            </div>
          </div>
        </div>
      </a-card>
    </div>
  </div>
</template>

<script lang="ts" setup>
  import { ref, reactive, nextTick } from 'vue';
  import axios from 'axios';
  import DOMPurify from 'dompurify';
  import { Message } from '@arco-design/web-vue';
  import {
    IconSelectAll,
    IconArrowRight,
    IconSave,
  } from '@arco-design/web-vue/es/icon';
  import XHeader from '@/components/header/x-header.vue';
  import { createCustomRecipe } from '@/api/custom_recipe';
  import { useRouter } from 'vue-router';
  import { useI18n } from 'vue-i18n';

  // Import extracted utils and components
  import { getCssSelector, IGNORED_CLASSES } from './utils/selector';
  import HtmlPreview from './components/HtmlPreview.vue';

  const router = useRouter();
  const { t } = useI18n();

  // --- State ---
  const currentStep = ref(1);
  const url = ref('');
  const enhancedMode = ref(false);
  const fetching = ref(false);
  const fetchError = ref('');
  const parsing = ref(false);
  const saving = ref(false);
  const htmlContent = ref('');
  const parsedItems = ref<any[]>([]);

  // Selection State
  const isSelectionMode = ref(true);
  const currentTargetField = ref<string>('');
  const previewRef = ref<InstanceType<typeof HtmlPreview> | null>(null);
  const resultsRef = ref<HTMLElement | null>(null);

  // Config State
  const config = reactive<{ [key: string]: string }>({
    item_selector: '',
    title_selector: '',
    link_selector: '',
    date_selector: '',
    description_selector: '',
  });

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

  const setTargetField = (field: string) => {
    currentTargetField.value = field;
    Message.info(t('htmlToRss.msg.pickInfo', { field }));
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

        // Auto-extract metadata for Step 3
        try {
          const doc = new DOMParser().parseFromString(raw, 'text/html');
          const title = doc.querySelector('title')?.innerText || '';
          const descMeta =
            doc.querySelector('meta[name="description"]') ||
            doc.querySelector('meta[property="og:description"]');
          const description = descMeta ? descMeta.getAttribute('content') : '';

          if (title) feedMeta.title = title.trim();
          if (description) feedMeta.description = description.trim();
        } catch (e) {
          // Ignore extraction errors
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

        feedMeta.link = url.value; // Auto-fill

        if (advanceStep) {
          nextStep();
          setTargetField('item_selector');
        }
      } else {
        const errorMsg = res.msg || t('htmlToRss.msg.fetchFailed');
        fetchError.value = errorMsg;
      }
    } catch (err: any) {
      // Axios interceptor throws an error with the backend message
      const errorMsg = err.message || t('htmlToRss.msg.errorFetching');
      fetchError.value = errorMsg;
    } finally {
      fetching.value = false;
    }
  };

  // Step 1 -> 2
  const fetchAndNext = async () => {
    await fetchContent(true);
  };

  // Step 2 -> 3
  const runPreview = async () => {
    if (!config.item_selector) return;
    parsing.value = true;
    try {
      const { data: res } = (await axios.post('/api/admin/tools/parse', {
        html: htmlContent.value,
        url: url.value,
        item_selector: config.item_selector,
        title_selector: config.title_selector,
        link_selector: config.link_selector,
        date_selector: config.date_selector,
        content_selector: config.description_selector,
      })) as any;

      if (res.code === 0) {
        parsedItems.value = res.data;
        if (parsedItems.value.length === 0) {
          Message.warning(t('htmlToRss.msg.noItems'));
          return;
        }
        Message.success(
          t('htmlToRss.msg.extracted', { count: parsedItems.value.length }),
        );
        // Do not auto-advance. Let user check preview first.
        nextTick(() => {
          resultsRef.value?.scrollIntoView({
            behavior: 'smooth',
            block: 'start',
          });
        });
      } else {
        Message.error(res.msg || t('htmlToRss.msg.parseFailed'));
      }
    } catch (err) {
      Message.error(t('htmlToRss.msg.errorParsing'));
    } finally {
      parsing.value = false;
    }
  };

  // Step 4: Save
  const handleSaveRecipe = async () => {
    if (!recipeMeta.id) {
      Message.warning(t('htmlToRss.msg.idRequired'));
      return;
    }

    saving.value = true;

    // Construct SourceConfig JSON based on Scenario D
    const sourceConfig = {
      type: 'html',
      http_fetcher: {
        url: url.value,
        use_browserless: enhancedMode.value,
      },
      html_parser: {
        item_selector: config.item_selector,
        title: config.title_selector,
        link: config.link_selector,
        date: config.date_selector,
        description: config.description_selector,
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
        craft: 'proxy', // Default craft flow
        source_type: 'html',
        source_config: JSON.stringify(sourceConfig),
      });

      Message.success(t('htmlToRss.msg.saveSuccess'));
      router.push({ name: 'CustomRecipe' }); // Navigate to custom recipe page
    } catch (err: any) {
      Message.error(t('htmlToRss.msg.saveFailed', { msg: err.message || err }));
    } finally {
      saving.value = false;
    }
  };

  // --- Selection Logic ---

  // This function now receives the RAW DOM element from HtmlPreview
  const handleElementSelect = (target: HTMLElement) => {
    if (!currentTargetField.value) {
      Message.warning(t('htmlToRss.msg.clickPickFirst'));
      return;
    }

    // Use the extracted utility, passing the iframe document for context checks
    const doc = previewRef.value?.contentDocument;
    const isItemSelector = currentTargetField.value === 'item_selector';

    // Calculate full absolute selector first
    const fullSelector = getCssSelector(
      target,
      doc || undefined,
      isItemSelector,
    );

    if (!doc) return;

    if (isItemSelector) {
      config.item_selector = fullSelector;
      try {
        const matches = doc.querySelectorAll(fullSelector);
        Message.success(
          t('htmlToRss.msg.matchedItems', {
            count: matches.length,
          }),
        );
      } catch {
        Message.success(
          t('htmlToRss.msg.setItemSelector', { selector: fullSelector }),
        );
      }
      currentTargetField.value = 'title_selector';
      Message.info(t('htmlToRss.msg.listSelectedNextTitle'));
    } else {
      // Relative selection logic
      if (!config.item_selector) {
        Message.warning(t('htmlToRss.msg.setListItemFirst'));
        return;
      }

      // Find which item container this target belongs to
      const items = doc.querySelectorAll(config.item_selector);
      let foundItem: HTMLElement | null = null;

      for (let i = 0; i < items.length; i += 1) {
        if (items[i].contains(target)) {
          foundItem = items[i] as HTMLElement;
          break;
        }
      }

      if (foundItem) {
        if (target === foundItem) {
          config[currentTargetField.value] = '.'; // "this" - explicit dot for backend
          Message.info(t('htmlToRss.msg.selectedContainer'));
        } else {
          // Calculate relative path
          const relPath: string[] = [];
          let curr: HTMLElement = target;

          while (curr && curr !== foundItem) {
            let selector = curr.tagName.toLowerCase();
            if (curr.classList.length > 0) {
              const validClasses = Array.from(curr.classList).filter(
                (c) => !IGNORED_CLASSES.includes(c),
              );
              if (validClasses.length > 0)
                selector += `.${CSS.escape(validClasses[0])}`;
            }
            relPath.unshift(selector);
            curr = curr.parentNode as HTMLElement;
          }

          config[currentTargetField.value] = relPath.join(' ');
          Message.success(
            t('htmlToRss.msg.setRelativePath', { path: relPath.join(' ') }),
          );
        }
      } else {
        Message.warning(t('htmlToRss.msg.insideListItem'));
      }
    }

    currentTargetField.value = ''; // Reset picker
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
</style>
