<template>
  <div class="py-8 px-16">
    <x-header
      :title="$t('searchToRss.title')"
      :description="$t('searchToRss.description')"
    ></x-header>

    <div class="content-wrapper">
      <a-card class="wizard-card">
        <a-steps :current="currentStep" class="mb-8">
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
          <a-form :model="fetchReq" layout="vertical" class="max-w-xl mx-auto">
            <a-form-item
              :label="$t('searchToRss.step1.label')"
              required
              :help="$t('searchToRss.step1.help')"
            >
              <a-input
                v-model="fetchReq.query"
                :placeholder="$t('searchToRss.step1.placeholder')"
                size="large"
                allow-clear
                @press-enter="handlePreview"
              />
            </a-form-item>
            <div class="text-center mt-8">
              <a-button
                type="primary"
                size="large"
                :loading="fetching"
                :disabled="!fetchReq.query"
                @click="handlePreview"
              >
                {{ $t('searchToRss.step1.button') }} <icon-arrow-right />
              </a-button>
              <div class="mt-4 text-gray-500 text-sm">
                <router-link
                  :to="{ name: 'SearchProvider' }"
                  class="text-blue-500 hover:underline"
                >
                  {{ $t('settings.searchProvider.configure') }}
                </router-link>
              </div>
            </div>
          </a-form>
        </div>

        <!-- STEP 2: Preview Results -->
        <div v-show="currentStep === 2" class="step-content flex flex-col">
          <div class="flex-1 overflow-y-auto mb-4">
            <a-alert type="success" class="mb-4">
              {{ $t('searchToRss.step2.alert', { count: parsedItems.length }) }}
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
                  :help="$t('searchToRss.step4.recipeId.help')"
                >
                  <a-input
                    v-model="recipeMeta.id"
                    :placeholder="$t('searchToRss.placeholder.recipeId')"
                  />
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
  import { ref, reactive } from 'vue';
  import { useRouter } from 'vue-router';
  import { Message } from '@arco-design/web-vue';
  import { IconArrowRight, IconSave } from '@arco-design/web-vue/es/icon';
  import XHeader from '@/components/header/x-header.vue';
  import { previewSearch, ParsedItem, SearchFetchReq } from '@/api/json_rss';
  import { createChannel } from '@/api/channel';
  import { useI18n } from 'vue-i18n';

  const router = useRouter();
  const { t } = useI18n();

  // --- State ---
  const currentStep = ref(1);
  const fetching = ref(false);
  const saving = ref(false);
  const parsedItems = ref<ParsedItem[]>([]);

  // Step 1: Query
  const fetchReq = reactive<SearchFetchReq>({
    query: '',
  });

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

  // Step 1 -> 2
  const handlePreview = async () => {
    if (!fetchReq.query) {
      Message.warning(t('searchToRss.msg.queryRequired'));
      return;
    }
    fetching.value = true;
    parsedItems.value = [];
    try {
      const res = await previewSearch(fetchReq);
      if (res.data) {
        parsedItems.value = res.data.map((item: any) => ({
          title: item.title || item.Title,
          link: item.link || item.Link,
          date: item.published || item.Published || item.date || item.Date,
          description: item.description || item.Description,
          content: item.content || item.Content,
        }));
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
        fetchReq.query,
      )}`; // Fallback link

      nextStep();
    } catch (err) {
      // handled by interceptor usually, but log here
      console.error(err);
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

    const sourceConfig = {
      type: 'search',
      search_fetcher: {
        query: fetchReq.query,
      },
      feed_meta: {
        title: feedMeta.title,
        description: feedMeta.description,
        link: feedMeta.link,
      },
    };

    try {
      await createChannel({
        id: recipeMeta.id,
        description:
          recipeMeta.description || `Search feed for: ${fetchReq.query}`,
        craft: 'proxy', // Default craft
        source_type: 'search',
        source_config: JSON.stringify(sourceConfig),
      });
      Message.success(t('searchToRss.msg.saved'));
      router.push({ name: 'Channel' });
    } catch (err: any) {
      console.error(err);
      Message.error(
        t('searchToRss.msg.saveFailed', { msg: err.message || err }),
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
    height: 500px; /* fixed height for consistent layout */
  }
</style>
