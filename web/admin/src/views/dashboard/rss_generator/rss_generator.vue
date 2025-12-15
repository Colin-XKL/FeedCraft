<template>
  <div class="py-8 px-16">
    <x-header
      title="RSS Generator Wizard"
      description="Create a custom RSS feed from any webpage in 4 simple steps."
    ></x-header>

    <div class="content-wrapper">
      <a-card class="wizard-card">
        <a-steps :current="currentStep" class="mb-8">
          <a-step title="Target URL" description="Enter webpage URL" />
          <a-step title="Extract Rules" description="Select content" />
          <a-step title="Feed Metadata" description="Define feed info" />
          <a-step title="Save Recipe" description="Review and save" />
        </a-steps>

        <!-- STEP 1: URL Input -->
        <div v-show="currentStep === 1" class="step-content">
          <a-form layout="vertical" class="max-w-xl mx-auto">
            <a-form-item
              label="Target Webpage URL"
              help="Enter the full URL of the page you want to turn into an RSS feed."
            >
              <a-input
                v-model="url"
                placeholder="https://example.com/blog"
                size="large"
                allow-clear
                @keyup.enter="fetchAndNext"
              />
            </a-form-item>
            <div class="text-center mt-8">
              <a-button
                type="primary"
                size="large"
                :loading="fetching"
                :disabled="!url"
                @click="fetchAndNext"
              >
                Fetch & Next <icon-arrow-right />
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
                <span class="font-bold">Page Preview</span>
                <a-tag v-if="isSelectionMode" color="blue"
                  >Selection Mode On</a-tag
                >
              </div>

              <!-- Use the extracted HtmlPreview component -->
              <HtmlPreview
                ref="previewRef"
                class="flex-1"
                :html-content="htmlContent"
                :is-selection-mode="isSelectionMode"
                @select="handleElementSelect"
              />
            </a-col>

            <!-- Right: Config & Preview -->
            <a-col :span="10" class="h-full flex flex-col">
              <div class="flex-1 overflow-y-auto pr-2">
                <a-alert type="info" class="mb-4">
                  1. Click "Pick" and select an element in the preview.<br />
                  2. Click "Run Preview" to verify extracted data.
                </a-alert>

                <a-form :model="config" layout="vertical">
                  <a-card
                    title="1. List Item (Required)"
                    size="small"
                    class="mb-4 border-blue-100"
                  >
                    <a-form-item label="CSS Selector">
                      <a-input
                        v-model="config.item_selector"
                        placeholder=".article-card"
                      >
                        <template #suffix>
                          <a-button
                            size="mini"
                            type="primary"
                            status="success"
                            @click="setTargetField('item_selector')"
                          >
                            <icon-select-all /> Pick
                          </a-button>
                        </template>
                      </a-input>
                    </a-form-item>
                  </a-card>

                  <a-card
                    title="2. Fields (Relative)"
                    size="small"
                    :class="{ 'opacity-50': !config.item_selector }"
                  >
                    <a-form-item label="Title">
                      <a-input
                        v-model="config.title_selector"
                        :disabled="!config.item_selector"
                      >
                        <template #suffix>
                          <a-button
                            size="mini"
                            @click="setTargetField('title_selector')"
                            :disabled="!config.item_selector"
                            >Pick</a-button
                          >
                        </template>
                      </a-input>
                    </a-form-item>
                    <a-form-item label="Link">
                      <a-input
                        v-model="config.link_selector"
                        :disabled="!config.item_selector"
                      >
                        <template #suffix>
                          <a-button
                            size="mini"
                            @click="setTargetField('link_selector')"
                            :disabled="!config.item_selector"
                            >Pick</a-button
                          >
                        </template>
                      </a-input>
                    </a-form-item>
                    <a-form-item label="Date (Optional)">
                      <a-input
                        v-model="config.date_selector"
                        :disabled="!config.item_selector"
                      >
                        <template #suffix>
                          <a-button
                            size="mini"
                            @click="setTargetField('date_selector')"
                            :disabled="!config.item_selector"
                            >Pick</a-button
                          >
                        </template>
                      </a-input>
                    </a-form-item>
                    <a-form-item label="Description (Optional)">
                      <a-input
                        v-model="config.description_selector"
                        :disabled="!config.item_selector"
                      >
                        <template #suffix>
                          <a-button
                            size="mini"
                            @click="setTargetField('description_selector')"
                            :disabled="!config.item_selector"
                            >Pick</a-button
                          >
                        </template>
                      </a-input>
                    </a-form-item>
                  </a-card>
                </a-form>

                <!-- Immediate Preview Results -->
                <div v-if="parsedItems.length > 0" class="mt-4">
                  <a-divider orientation="left"
                    >Preview Results ({{ parsedItems.length }})</a-divider
                  >
                  <a-collapse :default-active-key="[0]">
                    <a-collapse-item
                      v-for="(item, idx) in parsedItems"
                      :key="idx"
                      :header="item.title || '(No Title)'"
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
                <a-button @click="prevStep">Back</a-button>
                <a-space>
                  <a-button
                    type="outline"
                    :loading="parsing"
                    :disabled="!config.item_selector"
                    @click="runPreview"
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
            <a-alert title="Extraction Configured" type="success" class="mb-6">
              Rules defined! Extracted {{ parsedItems.length }} items. Now set
              the feed metadata.
            </a-alert>

            <a-form :model="feedMeta" layout="vertical">
              <a-form-item label="Feed Title" required>
                <a-input
                  v-model="feedMeta.title"
                  placeholder="e.g. My Tech Blog RSS"
                />
              </a-form-item>
              <a-form-item label="Feed Description">
                <a-textarea
                  v-model="feedMeta.description"
                  placeholder="A brief description of this feed..."
                />
              </a-form-item>
              <a-form-item label="Site Link">
                <a-input
                  v-model="feedMeta.link"
                  placeholder="Original website URL"
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
              <a-button type="primary" @click="nextStep">Next</a-button>
            </div>
          </div>
        </div>

        <!-- STEP 4: Save -->
        <div v-show="currentStep === 4" class="step-content">
          <div class="max-w-xl mx-auto">
            <a-card title="Review & Save Recipe" class="border-blue-100">
              <a-descriptions :column="1" title="Summary" bordered>
                <a-descriptions-item label="Source URL">{{
                  url
                }}</a-descriptions-item>
                <a-descriptions-item label="Feed Title">{{
                  feedMeta.title
                }}</a-descriptions-item>
                <a-descriptions-item label="Item Count"
                  >{{ parsedItems.length }} items detected</a-descriptions-item
                >
              </a-descriptions>

              <a-divider />

              <a-form :model="recipeMeta" layout="vertical" class="mt-6">
                <a-form-item
                  label="Recipe Unique ID"
                  required
                  help="e.g., 'tech-news-daily'"
                >
                  <a-input v-model="recipeMeta.id" placeholder="my-recipe-id" />
                </a-form-item>
                <a-form-item label="Internal Description">
                  <a-textarea
                    v-model="recipeMeta.description"
                    placeholder="Notes for yourself..."
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
                    <icon-save /> Confirm & Save
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

<script lang="ts" setup>
  import { ref, reactive } from 'vue';
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

  // Import extracted utils and components
  import { getCssSelector, IGNORED_CLASSES } from './utils/selector';
  import HtmlPreview from './components/HtmlPreview.vue';

  const router = useRouter();

  // --- State ---
  const currentStep = ref(1);
  const url = ref('');
  const fetching = ref(false);
  const parsing = ref(false);
  const saving = ref(false);
  const htmlContent = ref('');
  const parsedItems = ref<any[]>([]);

  // Selection State
  const isSelectionMode = ref(true);
  const currentTargetField = ref<string>('');
  const previewRef = ref<InstanceType<typeof HtmlPreview> | null>(null);

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

  // Step 1 -> 2
  const fetchAndNext = async () => {
    if (!url.value) return;
    fetching.value = true;
    try {
      const res = (await axios.post('/api/admin/tools/fetch', {
        url: url.value,
      })) as any;
      if (res.code === 0) {
        let raw = res.data;
        const baseTag = `<base href="${url.value}" />`;
        if (raw.toLowerCase().includes('<head>')) {
          raw = raw.replace(/<head>/i, `<head>${baseTag}`);
        } else {
          raw = `${baseTag}${raw}`;
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
        nextStep();
      } else {
        Message.error(res.msg || 'Fetch failed');
      }
    } catch (err) {
      Message.error('Error fetching page');
    } finally {
      fetching.value = false;
    }
  };

  // Step 2 -> 3
  const runPreview = async () => {
    if (!config.item_selector) return;
    parsing.value = true;
    try {
      const res = (await axios.post('/api/admin/tools/parse', {
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
          Message.warning('No items matched. Please check your selectors.');
          return;
        }
        Message.success(`Extracted ${parsedItems.value.length} items`);
        // Do not auto-advance. Let user check preview first.
      } else {
        Message.error(res.msg || 'Parse failed');
      }
    } catch (err) {
      Message.error('Error parsing content');
    } finally {
      parsing.value = false;
    }
  };

  // Step 4: Save
  const handleSaveRecipe = async () => {
    if (!recipeMeta.id) {
      Message.warning('Recipe Name (ID) is required');
      return;
    }

    saving.value = true;

    // Construct SourceConfig JSON based on Scenario D
    const sourceConfig = {
      type: 'html',
      http_fetcher: {
        url: url.value,
      },
      html_parser: {
        item_selector: config.item_selector,
        title: config.title_selector,
        link: config.link_selector,
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

      Message.success('Recipe saved successfully!');
      router.push({ name: 'CustomRecipes' }); // Assuming route exists or will exist
    } catch (err: any) {
      Message.error(`Failed to save recipe: ${err.message || err}`);
    } finally {
      saving.value = false;
    }
  };

  // --- Selection Logic ---

  const setTargetField = (field: string) => {
    currentTargetField.value = field;
    Message.info(
      `Click an element in the preview to pick selector for: ${field}`
    );
  };

  // This function now receives the RAW DOM element from HtmlPreview
  const handleElementSelect = (target: HTMLElement) => {
    if (!currentTargetField.value) {
      Message.warning('Please click a "Pick" button on the right first.');
      return;
    }

    // Use the extracted utility, passing the iframe document for context checks
    const doc = previewRef.value?.contentDocument;
    const isItemSelector = currentTargetField.value === 'item_selector';

    // Calculate full absolute selector first
    const fullSelector = getCssSelector(
      target,
      doc || undefined,
      isItemSelector
    );

    if (!doc) return;

    if (isItemSelector) {
      config.item_selector = fullSelector;
      try {
        const matches = doc.querySelectorAll(fullSelector);
        Message.success(
          `Matched ${matches.length} items with: ${fullSelector}`
        );
      } catch {
        Message.success(`Set Item Selector: ${fullSelector}`);
      }
    } else {
      // Relative selection logic
      if (!config.item_selector) {
        Message.warning('Set List Item Selector first!');
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
          Message.info('Selected the item container itself.');
        } else {
          // Calculate relative path
          const relPath: string[] = [];
          let curr: HTMLElement = target;

          while (curr && curr !== foundItem) {
            let selector = curr.tagName.toLowerCase();
            if (curr.classList.length > 0) {
              const validClasses = Array.from(curr.classList).filter(
                (c) => !IGNORED_CLASSES.includes(c)
              );
              if (validClasses.length > 0)
                selector += `.${CSS.escape(validClasses[0])}`;
            }
            relPath.unshift(selector);
            curr = curr.parentNode as HTMLElement;
          }

          config[currentTargetField.value] = relPath.join(' ');
          Message.success(`Set relative path: ${relPath.join(' ')}`);
        }
      } else {
        Message.warning('Selection must be inside a matched List Item!');
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
