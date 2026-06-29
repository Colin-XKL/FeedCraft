<template>
  <div class="py-8 px-16">
    <x-header
      :title="$t('searchToRss.title')"
      :description="$t('searchToRss.description')"
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
            :title="$t('searchToRss.step.searchQuery')"
            :description="$t('searchToRss.step.searchQuery.desc')"
          />
          <a-step
            :title="$t('searchToRss.step.previewResults')"
            :description="$t('searchToRss.step.previewResults.desc')"
          />
          <a-step
            :title="$t('searchToRss.step.feedMetadata')"
            :description="$t('searchToRss.step.feedMetadata.desc')"
          />
          <a-step
            :title="$t('searchToRss.step.saveRecipe')"
            :description="$t('searchToRss.step.saveRecipe.desc')"
          />
        </a-steps>

        <!-- STEP 1: Search Query -->
        <div v-show="currentStep === 1" class="step-content">
          <a-form :model="fetchReq" layout="vertical" class="max-w-2xl mx-auto">
            <a-form-item :label="$t('searchToRss.step1.modeLabel')" required>
              <div class="mode-grid">
                <button
                  v-for="option in searchModeOptions"
                  :key="option.value"
                  type="button"
                  class="mode-card"
                  :class="{ 'mode-card--active': searchMode === option.value }"
                  @click="searchMode = option.value"
                >
                  <div class="mode-card__header">
                    <span class="mode-card__icon">
                      <icon-search v-if="option.value === 'keyword'" />
                      <icon-robot v-else />
                    </span>
                    <span class="mode-card__title">
                      {{ $t(option.titleKey) }}
                    </span>
                    <a-tag size="small" color="arcoblue">
                      {{ $t(option.badgeKey) }}
                    </a-tag>
                  </div>
                  <p class="mode-card__description">
                    {{ $t(option.descriptionKey) }}
                  </p>
                </button>
              </div>
            </a-form-item>

            <a-form-item
              :label="currentSearchModeTitle"
              required
              :help="currentSearchModeHelp"
            >
              <a-input
                v-model="fetchReq.query"
                :placeholder="currentSearchModePlaceholder"
                size="large"
                allow-clear
                @press-enter="handlePreview"
              />
            </a-form-item>

            <a-alert class="mode-tip" type="info">
              <div class="flex flex-wrap items-center gap-2">
                <span>{{ $t('searchToRss.step1.providerTip') }}</span>
                <router-link
                  :to="{ name: 'SearchProvider' }"
                  class="text-blue-600 hover:underline"
                >
                  {{ $t('settings.searchProvider.configure') }}
                </router-link>
              </div>
            </a-alert>

            <div class="text-center mt-12">
              <a-button
                type="primary"
                size="large"
                :loading="fetching"
                :disabled="!fetchReq.query"
                class="w-full sm:w-64"
                @click="handlePreview"
              >
                {{ $t('searchToRss.step1.button') }} <icon-arrow-right />
              </a-button>
            </div>
          </a-form>
        </div>

        <!-- STEP 2: Preview Results -->
        <div v-show="currentStep === 2" class="step-content flex flex-col">
          <div class="flex-1 overflow-y-auto mb-4">
            <a-alert type="success" class="mb-4">
              <div class="flex flex-wrap items-center gap-2">
                <span>
                  {{
                    $t('searchToRss.step2.alert', { count: parsedItems.length })
                  }}
                </span>
                <a-tag color="green">{{ currentSearchModeTitle }}</a-tag>
              </div>
            </a-alert>
            <a-list :data="parsedItems" :bordered="false">
              <template #item="{ item }">
                <a-list-item>
                  <a-list-item-meta
                    :title="item.title"
                    :description="item.date"
                  >
                  </a-list-item-meta>
                  <div class="mb-1">
                    <a
                      :href="item.link"
                      target="_blank"
                      class="text-blue-600 hover:underline"
                      >{{ item.link }}</a
                    >
                  </div>
                  <div
                    v-if="item.description"
                    class="text-xs text-gray-500 line-clamp-2"
                  >
                    {{ item.description }}
                  </div>
                </a-list-item>
              </template>
            </a-list>
          </div>

          <div class="flex justify-between pt-4 border-t border-gray-100">
            <a-button @click="prevStep">{{
              $t('searchToRss.common.back')
            }}</a-button>
            <a-button type="primary" @click="nextStep">{{
              $t('searchToRss.common.next')
            }}</a-button>
          </div>
        </div>

        <!-- STEP 3: Feed Metadata -->
        <div v-show="currentStep === 3" class="step-content">
          <div class="max-w-2xl mx-auto">
            <a-alert class="mb-6">
              {{ $t('searchToRss.step3.alert') }}
            </a-alert>
            <a-form :model="feedMeta" layout="vertical">
              <a-form-item :label="$t('searchToRss.step3.feedTitle')" required>
                <a-input v-model="feedMeta.title" allow-clear />
              </a-form-item>
              <a-form-item :label="$t('searchToRss.step3.feedDescription')">
                <a-textarea
                  v-model="feedMeta.description"
                  :auto-size="{ minRows: 3, maxRows: 5 }"
                  allow-clear
                />
              </a-form-item>
              <a-form-item :label="$t('searchToRss.step3.siteLink')">
                <a-input v-model="feedMeta.link" allow-clear />
              </a-form-item>
            </a-form>

            <div class="flex justify-between mt-8">
              <a-button @click="prevStep">{{
                $t('searchToRss.common.back')
              }}</a-button>
              <a-button type="primary" @click="nextStep">{{
                $t('searchToRss.common.next')
              }}</a-button>
            </div>
          </div>
        </div>

        <!-- STEP 4: Save Recipe -->
        <div v-show="currentStep === 4" class="step-content">
          <div class="max-w-xl mx-auto">
            <a-card
              :title="$t('searchToRss.step4.reviewAndSave')"
              class="border-blue-100"
            >
              <a-descriptions :column="1" bordered>
                <a-descriptions-item :label="$t('searchToRss.step4.mode')">
                  {{ currentSearchModeTitle }}
                </a-descriptions-item>
                <a-descriptions-item :label="$t('searchToRss.step4.query')">{{
                  fetchReq.query
                }}</a-descriptions-item>
                <a-descriptions-item
                  :label="$t('searchToRss.step3.feedTitle')"
                  >{{ feedMeta.title }}</a-descriptions-item
                >
                <a-descriptions-item
                  :label="$t('searchToRss.step4.itemsFound')"
                  >{{ parsedItems.length }}</a-descriptions-item
                >
              </a-descriptions>

              <a-divider />

              <a-form :model="recipeMeta" layout="vertical" class="mt-6">
                <a-form-item
                  :label="$t('searchToRss.step4.recipeId')"
                  required
                  field="id"
                  :rules="
                    getRecipeIdRules($t('searchToRss.step4.recipeId.help'))
                  "
                  :help="$t('searchToRss.step4.recipeId.help')"
                >
                  <a-input
                    v-model="recipeMeta.id"
                    :placeholder="$t('searchToRss.placeholder.recipeId')"
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
                <a-form-item
                  :label="$t('searchToRss.step4.internalDescription')"
                >
                  <a-textarea
                    v-model="recipeMeta.description"
                    :placeholder="$t('searchToRss.placeholder.internalDesc')"
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
                    <icon-save /> {{ $t('searchToRss.step4.confirmAndSave') }}
                  </a-button>
                </div>
              </a-form>
            </a-card>

            <div class="flex justify-start mt-8">
              <a-button @click="prevStep">{{
                $t('searchToRss.common.back')
              }}</a-button>
            </div>
          </div>
        </div>
      </a-card>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { computed, ref, reactive, watch } from 'vue';
  import { useRouter } from 'vue-router';
  import { Message } from '@arco-design/web-vue';
  import {
    IconArrowRight,
    IconSearch,
    IconRobot,
    IconSave,
    IconRefresh,
  } from '@arco-design/web-vue/es/icon';
  import XHeader from '@/components/header/x-header.vue';
  import {
    previewSearch,
    SearchFetchReq,
    SearchPreviewItem,
  } from '@/api/json_rss';
  import {
    buildSearchFetchReq,
    buildSearchSourceConfig,
    SearchMode,
    searchModeOptions,
  } from '@/views/dashboard/search_to_rss/searchMode';
  import { createCustomRecipe } from '@/api/custom_recipe';
  import { useI18n } from 'vue-i18n';
  import generateRecipeId, { getRecipeIdRules } from '@/utils/slug';

  const router = useRouter();
  const { t } = useI18n();

  // --- State ---
  const currentStep = ref(1);
  const fetching = ref(false);
  const saving = ref(false);
  const parsedItems = ref<SearchPreviewItem[]>([]);

  // Step 1: Query
  const searchMode = ref<SearchMode>('keyword');
  const fetchReq = reactive<SearchFetchReq>({
    query: '',
  });

  const currentSearchModeOption = computed(
    () =>
      searchModeOptions.find((option) => option.value === searchMode.value) ||
      searchModeOptions[0]
  );

  const currentSearchModeTitle = computed(() =>
    t(currentSearchModeOption.value.titleKey)
  );
  const currentSearchModeHelp = computed(() =>
    t(currentSearchModeOption.value.helpKey)
  );
  const currentSearchModePlaceholder = computed(() =>
    t(currentSearchModeOption.value.placeholderKey)
  );

  // Step 3: Feed Meta
  const feedMeta = reactive({
    title: '',
    description: '',
    link: '',
  });

  // Step 4: Recipe Meta
  const recipeMeta = reactive({
    id: '',
    description: '',
  });

  // --- Actions ---

  const nextStep = () => {
    if (currentStep.value < 4) currentStep.value += 1;
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
        recipeMeta.id = generateRecipeId(feedMeta.title);
      }
    }
  );

  // Step 1 -> 2
  const handlePreview = async () => {
    if (!fetchReq.query) {
      Message.warning(t('searchToRss.msg.queryRequired'));
      return;
    }
    fetching.value = true;
    parsedItems.value = [];
    try {
      const res = await previewSearch(
        buildSearchFetchReq(fetchReq.query, searchMode.value)
      );
      if (res.data) {
        parsedItems.value = res.data;
      }

      if (parsedItems.value.length === 0) {
        Message.info(t('searchToRss.msg.noResults'));
        return;
      }

      // Auto-populate Meta
      feedMeta.title = t('searchToRss.default.title', {
        query: fetchReq.query,
      });
      feedMeta.description = t('searchToRss.default.description', {
        query: fetchReq.query,
      });
      feedMeta.link = `https://google.com/search?q=${encodeURIComponent(
        fetchReq.query
      )}`; // Fallback link

      nextStep();
    } catch (err) {
      Message.error(err instanceof Error ? err.message : String(err));
    } finally {
      fetching.value = false;
    }
  };

  // Step 4: Save
  const handleSaveRecipe = async () => {
    if (!recipeMeta.id) {
      Message.error(t('searchToRss.msg.idRequired'));
      return;
    }

    saving.value = true;

    const sourceConfig = buildSearchSourceConfig(
      fetchReq.query,
      searchMode.value,
      feedMeta
    );

    try {
      await createCustomRecipe({
        id: recipeMeta.id,
        description:
          recipeMeta.description || `Search feed for: ${fetchReq.query}`,
        craft: 'proxy', // Default craft
        source_type: 'search',
        source_config: JSON.stringify(sourceConfig),
      });
      Message.success(t('searchToRss.msg.saved'));
      router.push({ name: 'CustomRecipe' });
    } catch (err: any) {
      Message.error(
        t('searchToRss.msg.saveFailed', { msg: err.message || err })
      );
    } finally {
      saving.value = false;
    }
  };
</script>

<style scoped>
  .wizard-card {
    min-height: 600px;
  }
  .step-content {
    margin-top: 24px;
    min-height: 450px;
  }

  .mode-grid {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 16px;
    width: 100%;
  }

  .mode-card {
    display: flex;
    flex-direction: column;
    gap: 12px;
    padding: 18px;
    text-align: left;
    background: #fff;
    border: 1px solid var(--color-border-2);
    border-radius: 14px;
    cursor: pointer;
    transition: border-color 0.2s ease, box-shadow 0.2s ease,
      transform 0.2s ease;
  }

  .mode-card:hover,
  .mode-card--active {
    border-color: rgb(var(--primary-6));
    box-shadow: 0 8px 24px rgba(22, 93, 255, 0.12);
    transform: translateY(-1px);
  }

  .mode-card--active {
    background: linear-gradient(180deg, #f7fbff 0%, #fff 100%);
  }

  .mode-card__header {
    display: flex;
    align-items: center;
    gap: 10px;
  }

  .mode-card__icon {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    color: rgb(var(--primary-6));
    background: rgb(var(--primary-1));
    border-radius: 10px;
    font-size: 18px;
  }

  .mode-card__title {
    flex: 1;
    color: var(--color-text-1);
    font-size: 16px;
    font-weight: 600;
  }

  .mode-card__description {
    margin: 0;
    color: var(--color-text-3);
    font-size: 13px;
    line-height: 1.6;
  }

  .mode-tip {
    margin-top: 4px;
  }

  @media (max-width: 640px) {
    .mode-grid {
      grid-template-columns: 1fr;
    }
  }
</style>
