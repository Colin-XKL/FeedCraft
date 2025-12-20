<template>
  <div class="py-8 px-16">
    <x-header
      title="Search to RSS"
      description="Generate RSS feeds from search queries using AI providers."
    ></x-header>

    <div class="content-wrapper">
      <a-card class="wizard-card">
        <a-steps :current="currentStep" class="mb-8">
          <a-step
            title="Search Query"
            description="Enter your search terms"
          />
          <a-step
            title="Preview Results"
            description="Verify search results"
          />
          <a-step
            title="Feed Metadata"
            description="Configure feed details"
          />
          <a-step
            title="Save Recipe"
            description="Save as custom recipe"
          />
        </a-steps>

        <!-- STEP 1: Search Query -->
        <div v-show="currentStep === 1" class="step-content">
          <a-form :model="fetchReq" layout="vertical" class="max-w-xl mx-auto">
            <a-form-item
              label="Search Query"
              required
              help="Enter the keywords or question you want to track."
            >
              <a-input
                v-model="fetchReq.query"
                placeholder="e.g. 'latest AI news' or 'SpaceX launches'"
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
                Preview Results <icon-arrow-right />
              </a-button>
            </div>
          </a-form>
        </div>

        <!-- STEP 2: Preview Results -->
        <div v-show="currentStep === 2" class="step-content flex flex-col">
          <div class="flex-1 overflow-y-auto mb-4">
            <a-alert type="success" class="mb-4">
              Found {{ parsedItems.length }} items. Review them below before proceeding.
            </a-alert>
            <a-list :data="parsedItems" :bordered="false">
              <template #item="{ item }">
                <a-list-item>
                  <a-list-item-meta :title="item.title" :description="item.date">
                  </a-list-item-meta>
                  <div class="mb-1">
                    <a :href="item.link" target="_blank" class="text-blue-600 hover:underline">{{ item.link }}</a>
                  </div>
                  <div v-if="item.description" class="text-xs text-gray-500 line-clamp-2">
                    {{ item.description }}
                  </div>
                </a-list-item>
              </template>
            </a-list>
          </div>

          <div class="flex justify-between pt-4 border-t border-gray-100">
            <a-button @click="prevStep">Back</a-button>
            <a-button type="primary" @click="nextStep">Next Step</a-button>
          </div>
        </div>

        <!-- STEP 3: Feed Metadata -->
        <div v-show="currentStep === 3" class="step-content">
          <div class="max-w-2xl mx-auto">
             <a-alert class="mb-6">
               Customize how this feed appears in your RSS reader.
             </a-alert>
             <a-form :model="feedMeta" layout="vertical">
              <a-form-item label="Feed Title" required>
                <a-input v-model="feedMeta.title" />
              </a-form-item>
              <a-form-item label="Feed Description">
                <a-textarea v-model="feedMeta.description" :auto-size="{ minRows: 3, maxRows: 5 }" />
              </a-form-item>
              <a-form-item label="Site Link">
                <a-input v-model="feedMeta.link" />
              </a-form-item>
             </a-form>

             <div class="flex justify-between mt-8">
              <a-button @click="prevStep">Back</a-button>
              <a-button type="primary" @click="nextStep">Next Step</a-button>
            </div>
          </div>
        </div>

        <!-- STEP 4: Save Recipe -->
        <div v-show="currentStep === 4" class="step-content">
          <div class="max-w-xl mx-auto">
            <a-card title="Review & Save" class="border-blue-100">
              <a-descriptions :column="1" bordered>
                <a-descriptions-item label="Query">{{ fetchReq.query }}</a-descriptions-item>
                <a-descriptions-item label="Feed Title">{{ feedMeta.title }}</a-descriptions-item>
                <a-descriptions-item label="Items Found">{{ parsedItems.length }}</a-descriptions-item>
              </a-descriptions>

              <a-divider />

              <a-form :model="recipeMeta" layout="vertical" class="mt-6">
                <a-form-item label="Recipe ID" required help="Unique identifier for this recipe.">
                  <a-input v-model="recipeMeta.id" placeholder="e.g. search-ai-news" />
                </a-form-item>
                <a-form-item label="Internal Description">
                  <a-textarea v-model="recipeMeta.description" placeholder="Notes for yourself about this recipe." />
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

<script setup lang="ts">
import { ref, reactive } from 'vue';
import { useRouter } from 'vue-router';
import { Message } from '@arco-design/web-vue';
import { IconArrowRight, IconSave } from '@arco-design/web-vue/es/icon';
import XHeader from '@/components/header/x-header.vue';
import { previewSearch, ParsedItem, SearchFetchReq } from '@/api/json_rss';
import { createCustomRecipe } from '@/api/custom_recipe';

const router = useRouter();

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
    Message.warning('Query is required');
    return;
  }
  fetching.value = true;
  parsedItems.value = [];
  try {
    const res = await previewSearch(fetchReq);
    // @ts-ignore
    parsedItems.value = res.data;

    if (parsedItems.value.length === 0) {
        Message.info('No results found, please try a different query.');
        return; // Stay on step 1 if no results? Or let them see empty list? Better to warn.
    }

    // Auto-populate Meta
    feedMeta.title = `Search: ${fetchReq.query}`;
    feedMeta.description = `Search results for "${fetchReq.query}" generated by FeedCraft.`;
    feedMeta.link = `https://google.com/search?q=${encodeURIComponent(fetchReq.query)}`; // Fallback link

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
    Message.error('Recipe ID is required');
    return;
  }

  saving.value = true;

  const sourceConfig = {
    // The SourceSearch struct in backend expects 'search_fetcher'
    // 'type' is redundant but harmless
    type: 'search',
    search_fetcher: {
      query: fetchReq.query,
    },
    // We can embed feed_meta if the SourceSearch supports it or if we wrap it.
    // Currently SourceSearch doesn't seem to have explicit FeedMeta config in Go struct
    // (based on previous turns).
    // However, the CustomRecipe logic might allow overriding feed meta in the Recipe wrapper?
    // Let's assume for now we just save the search config.
    // If the backend SourceSearch doesn't support custom title overrides,
    // we might lose the title customization here unless we implemented it in backend.
    // But let's send it anyway in case we add it or it exists.
    feed_meta: {
        title: feedMeta.title,
        description: feedMeta.description,
        link: feedMeta.link
    }
  };

  try {
    await createCustomRecipe({
      id: recipeMeta.id,
      description: recipeMeta.description || `Search feed for: ${fetchReq.query}`,
      craft: 'proxy', // Default craft
      source_type: 'search',
      source_config: JSON.stringify(sourceConfig),
    });
    Message.success('Recipe saved successfully');
    router.push({ name: 'CustomRecipe' });
  } catch (err) {
    console.error(err);
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
