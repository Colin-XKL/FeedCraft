<template>
  <div class="py-8 px-[clamp(20px,4vw,64px)] max-md:py-6 max-md:px-4">
    <x-header
      :title="t('menu.feedViewer')"
      :description="t('feedViewer.description')"
    >
    </x-header>

    <a-card class="my-2">
      <template #title>
        <div class="flex items-center justify-between">
          <span>{{ t('feedViewer.inputLink') }}</span>
          <a-radio-group
            v-model="pageMode"
            type="button"
            size="small"
            @change="onPageModeChange"
          >
            <a-radio value="preview">{{
              t('feedViewer.pageMode.preview')
            }}</a-radio>
            <a-radio value="compare">{{
              t('feedViewer.pageMode.compare')
            }}</a-radio>
          </a-radio-group>
        </div>
      </template>

      <p>{{ t('feedViewer.inputTip') }}</p>
      <a-radio-group
        v-model="previewMode"
        type="button"
        class="mt-3"
        @change="clearPreviewState"
      >
        <a-radio value="url">{{ t('feedViewer.mode.url') }}</a-radio>
        <a-radio value="recipe">{{ t('feedViewer.mode.recipe') }}</a-radio>
        <a-radio value="topic">{{ t('feedViewer.mode.topic') }}</a-radio>
        <a-radio value="inbox">{{ t('feedViewer.mode.inbox') }}</a-radio>
        <a-radio value="uri">{{ t('feedViewer.mode.uri') }}</a-radio>
      </a-radio-group>

      <div
        class="mt-3 grid items-center gap-3 max-md:grid-cols-1"
        :class="
          pageMode === 'compare'
            ? 'grid-cols-[minmax(18rem,1fr)_minmax(14rem,20rem)_auto]'
            : 'grid-cols-[minmax(20rem,1fr)_auto]'
        "
      >
        <a-input
          v-if="previewMode === 'url'"
          v-model="feedUrl"
          type="text"
          class="w-full min-w-0"
          :placeholder="t('feedViewer.placeholder')"
          allow-clear
          @input="clearPreviewState"
          @keyup.enter="handleAction"
        />
        <a-select
          v-else-if="previewMode === 'recipe'"
          v-model="selectedRecipeId"
          class="w-full min-w-0"
          :placeholder="t('feedViewer.placeholder.recipe')"
          :loading="resourceLoading"
          allow-search
          allow-clear
          @change="clearPreviewState"
        >
          <a-option
            v-for="recipe in recipes"
            :key="recipe.id"
            :value="recipe.id"
          >
            {{ formatRecipeOption(recipe) }}
          </a-option>
        </a-select>
        <a-select
          v-else-if="previewMode === 'topic'"
          v-model="selectedTopicId"
          class="w-full min-w-0"
          :placeholder="t('feedViewer.placeholder.topic')"
          :loading="resourceLoading"
          allow-search
          allow-clear
          @change="clearPreviewState"
        >
          <a-option v-for="topic in topics" :key="topic.id" :value="topic.id">
            {{ formatTopicOption(topic) }}
          </a-option>
        </a-select>
        <a-select
          v-else-if="previewMode === 'inbox'"
          v-model="selectedInboxId"
          class="w-full min-w-0"
          :placeholder="t('feedViewer.placeholder.inbox')"
          :loading="resourceLoading"
          allow-search
          allow-clear
          @change="clearPreviewState"
        >
          <a-option v-for="inbox in inboxes" :key="inbox.id" :value="inbox.id">
            {{ formatInboxOption(inbox) }}
          </a-option>
        </a-select>
        <a-input
          v-else
          v-model="advancedURI"
          type="text"
          class="w-full min-w-0"
          :placeholder="t('feedViewer.placeholder.uri')"
          allow-clear
          @input="clearPreviewState"
          @keyup.enter="handleAction"
        />

        <CraftFlowSelect
          v-if="pageMode === 'compare'"
          v-model="selectedCraft"
          mode="single"
          class="w-full min-w-0"
          @change="clearPreviewState"
        />

        <a-button
          :loading="isLoading"
          :disabled="
            pageMode === 'compare'
              ? !currentInputURI || !selectedCraft
              : !currentInputURI
          "
          :type="pageMode === 'compare' ? 'primary' : undefined"
          @click="handleAction"
        >
          {{
            pageMode === 'compare'
              ? t('feedViewer.compare')
              : t('feedViewer.preview')
          }}
        </a-button>
      </div>

      <a-alert v-if="currentInputURI" type="info" class="mt-3" show-icon>
        {{ t('feedViewer.currentInput') }}: {{ currentInputURI }}
      </a-alert>
    </a-card>

    <!-- 预览模式：单栏结果 -->
    <template v-if="pageMode === 'preview'">
      <a-card
        :title="t('feedViewer.resultPreview')"
        class="my-4"
        :loading="isLoading"
      >
        <a-alert v-if="errorMessage" type="error" class="mb-4" show-icon>
          {{ errorMessage }}
        </a-alert>
        <div v-if="feedContent">
          <FeedViewContainer :feed-data="feedContent" />
        </div>
        <a-empty v-else-if="!errorMessage" />
      </a-card>
    </template>

    <!-- 对比模式：双栏结果 -->
    <template v-else>
      <a-row :gutter="[24, 24]" class="my-4">
        <a-col :span="12" :xs="24" :lg="12">
          <a-card
            :title="t('feedViewer.compare.originalFeed')"
            :loading="isLoading"
          >
            <a-alert
              v-if="compareOriginalError"
              type="error"
              class="mb-4"
              show-icon
            >
              {{ compareOriginalError }}
            </a-alert>
            <div v-if="compareOriginalContent">
              <FeedViewContainer :feed-data="compareOriginalContent" />
            </div>
            <a-empty v-else-if="!compareOriginalError" />
          </a-card>
        </a-col>
        <a-col :span="12" :xs="24" :lg="12">
          <a-card
            :title="t('feedViewer.compare.craftAppliedFeed')"
            :loading="isLoading"
          >
            <a-alert
              v-if="compareCraftError"
              type="error"
              class="mb-4"
              show-icon
            >
              {{ compareCraftError }}
            </a-alert>
            <div v-if="compareCraftContent">
              <FeedViewContainer :feed-data="compareCraftContent" />
            </div>
            <a-empty v-else-if="!compareCraftError" />
          </a-card>
        </a-col>
      </a-row>
    </template>
  </div>
</template>

<script lang="ts" setup>
  import { computed, onMounted, ref, watch } from 'vue';
  import FeedViewContainer from '@/views/dashboard/feed_viewer/feed_view_container.vue';
  import CraftFlowSelect from '@/views/dashboard/craft_flow/CraftFlowSelect.vue';
  import XHeader from '@/components/header/x-header.vue';
  import { useI18n } from 'vue-i18n';
  import { useRoute } from 'vue-router';
  import { previewFeed, type FeedViewerPreview } from '@/api/feed_viewer';
  import { getCustomRecipes, type CustomRecipe } from '@/api/custom_recipe';
  import { listTopicFeeds, type TopicFeed } from '@/api/topic';
  import { listInboxes, type Inbox } from '@/api/inbox';
  import { Message } from '@arco-design/web-vue';

  type PreviewMode = 'url' | 'recipe' | 'topic' | 'inbox' | 'uri';
  type PageMode = 'preview' | 'compare';

  const { t } = useI18n();
  const route = useRoute();

  const pageMode = ref<PageMode>('preview');
  const previewMode = ref<PreviewMode>('url');
  const feedUrl = ref('');
  const advancedURI = ref('');
  const selectedRecipeId = ref('');
  const selectedTopicId = ref('');
  const selectedInboxId = ref('');
  const selectedCraft = ref('');
  const recipes = ref<CustomRecipe[]>([]);
  const topics = ref<TopicFeed[]>([]);
  const inboxes = ref<Inbox[]>([]);

  // Preview mode state
  const feedContent = ref<FeedViewerPreview | null>(null);
  const errorMessage = ref('');

  // Compare mode state
  const compareOriginalContent = ref<FeedViewerPreview | null>(null);
  const compareCraftContent = ref<FeedViewerPreview | null>(null);
  const compareOriginalError = ref('');
  const compareCraftError = ref('');

  const isLoading = ref(false);
  const resourceLoading = ref(false);
  let previewRequestSeq = 0;

  const currentInputURI = computed(() => {
    if (previewMode.value === 'recipe' && selectedRecipeId.value) {
      return `feedcraft://recipe/${selectedRecipeId.value}`;
    }
    if (previewMode.value === 'topic' && selectedTopicId.value) {
      return `feedcraft://topic/${selectedTopicId.value}`;
    }
    if (previewMode.value === 'inbox' && selectedInboxId.value) {
      return `feedcraft://inbox/${selectedInboxId.value}`;
    }
    if (previewMode.value === 'uri') {
      return advancedURI.value.trim();
    }
    return feedUrl.value.trim();
  });

  function clearPreviewState() {
    previewRequestSeq += 1;
    errorMessage.value = '';
    feedContent.value = null;
    compareOriginalContent.value = null;
    compareCraftContent.value = null;
    compareOriginalError.value = '';
    compareCraftError.value = '';
  }

  function onPageModeChange() {
    clearPreviewState();
  }

  async function fetchPreview() {
    const inputURI = currentInputURI.value;
    if (!inputURI) return;
    const requestSeq = previewRequestSeq + 1;
    previewRequestSeq = requestSeq;
    isLoading.value = true;
    errorMessage.value = '';
    feedContent.value = null;
    try {
      const response = await previewFeed(inputURI);
      if (requestSeq !== previewRequestSeq) return;
      feedContent.value = response.data;
    } catch (error) {
      if (requestSeq !== previewRequestSeq) return;
      feedContent.value = null;
      errorMessage.value =
        error instanceof Error
          ? error.message
          : t('feedViewer.message.unknownError');
    } finally {
      if (requestSeq === previewRequestSeq) {
        isLoading.value = false;
      }
    }
  }

  async function fetchCompare() {
    const inputURI = currentInputURI.value;
    if (!inputURI || !selectedCraft.value) {
      Message.warning(t('feedViewer.compare.message.inputRequired'));
      return;
    }
    const requestSeq = previewRequestSeq + 1;
    previewRequestSeq = requestSeq;
    isLoading.value = true;
    compareOriginalError.value = '';
    compareCraftError.value = '';
    compareOriginalContent.value = null;
    compareCraftContent.value = null;

    const [originalResult, craftedResult] = await Promise.allSettled([
      previewFeed(inputURI),
      previewFeed(inputURI, { craftName: selectedCraft.value }),
    ]);

    if (requestSeq !== previewRequestSeq) return;

    if (originalResult.status === 'fulfilled') {
      compareOriginalContent.value = originalResult.value.data;
    } else {
      compareOriginalError.value =
        originalResult.reason instanceof Error
          ? originalResult.reason.message
          : t('feedViewer.message.unknownError');
    }

    if (craftedResult.status === 'fulfilled') {
      compareCraftContent.value = craftedResult.value.data;
    } else {
      compareCraftError.value =
        craftedResult.reason instanceof Error
          ? craftedResult.reason.message
          : t('feedViewer.message.unknownError');
    }

    isLoading.value = false;
  }

  function handleAction() {
    if (pageMode.value === 'compare') {
      fetchCompare();
    } else {
      fetchPreview();
    }
  }

  function firstQueryValue(value: unknown): string {
    if (Array.isArray(value)) return String(value[0] || '');
    return typeof value === 'string' ? value : '';
  }

  function selectInputURI(inputURI: string) {
    if (!inputURI) return;
    try {
      const parsed = new URL(inputURI);
      if (parsed.protocol === 'feedcraft:') {
        const id = parsed.pathname.replace(/^\/+/, '');
        if (parsed.hostname === 'recipe') {
          previewMode.value = 'recipe';
          selectedRecipeId.value = id;
          return;
        }
        if (parsed.hostname === 'topic') {
          previewMode.value = 'topic';
          selectedTopicId.value = id;
          return;
        }
        if (parsed.hostname === 'inbox') {
          previewMode.value = 'inbox';
          selectedInboxId.value = id;
          return;
        }
      }
      if (parsed.protocol === 'http:' || parsed.protocol === 'https:') {
        previewMode.value = 'url';
        feedUrl.value = inputURI;
        return;
      }
    } catch {
      // Fall through to advanced URI mode.
    }
    previewMode.value = 'uri';
    advancedURI.value = inputURI;
  }

  function resetPreviewTarget() {
    previewMode.value = 'url';
    feedUrl.value = '';
    advancedURI.value = '';
    selectedRecipeId.value = '';
    selectedTopicId.value = '';
    selectedInboxId.value = '';
    clearPreviewState();
  }

  function applyRouteQuery() {
    const target = firstQueryValue(route.query.target || route.query.mode);
    const id = firstQueryValue(route.query.id);
    if (target === 'recipe') {
      previewMode.value = 'recipe';
      selectedRecipeId.value = id || firstQueryValue(route.query.recipe_id);
      return;
    }
    if (target === 'topic') {
      previewMode.value = 'topic';
      selectedTopicId.value = id || firstQueryValue(route.query.topic_id);
      return;
    }
    if (target === 'inbox') {
      previewMode.value = 'inbox';
      selectedInboxId.value = id || firstQueryValue(route.query.inbox_id);
      return;
    }

    const inputURI = firstQueryValue(
      route.query.input_uri || route.query.uri || route.query.url
    );
    if (inputURI) {
      selectInputURI(inputURI);
      return;
    }
    resetPreviewTarget();
  }

  async function loadPreviewResources() {
    resourceLoading.value = true;
    const [recipeResult, topicResult, inboxResult] = await Promise.allSettled([
      getCustomRecipes(),
      listTopicFeeds(),
      listInboxes(),
    ]);
    if (recipeResult.status === 'fulfilled') {
      recipes.value = recipeResult.value.data ?? [];
    }
    if (topicResult.status === 'fulfilled') {
      topics.value = topicResult.value.data ?? [];
    }
    if (inboxResult.status === 'fulfilled') {
      inboxes.value = inboxResult.value.data ?? [];
    }
    if (
      recipeResult.status === 'rejected' ||
      topicResult.status === 'rejected' ||
      inboxResult.status === 'rejected'
    ) {
      Message.warning(t('feedViewer.message.resourceLoadFailed'));
    }
    resourceLoading.value = false;
  }

  function formatRecipeOption(recipe: CustomRecipe) {
    return recipe.description
      ? `${recipe.id} · ${recipe.description}`
      : recipe.id;
  }

  function formatTopicOption(topic: TopicFeed) {
    return topic.title ? `${topic.id} · ${topic.title}` : topic.id;
  }

  function formatInboxOption(inbox: Inbox) {
    return inbox.title ? `${inbox.id} · ${inbox.title}` : inbox.id;
  }

  onMounted(async () => {
    applyRouteQuery();
    await loadPreviewResources();
    if (currentInputURI.value) {
      handleAction();
    }
  });

  watch(
    () => route.fullPath,
    () => {
      applyRouteQuery();
      if (currentInputURI.value) {
        handleAction();
      } else {
        clearPreviewState();
      }
    }
  );
</script>

<script lang="ts">
  export default {
    name: 'FeedViewer',
  };
</script>
