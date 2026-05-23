<template>
  <div class="py-8 px-[clamp(20px,4vw,64px)] max-md:py-6 max-md:px-4">
    <x-header
      :title="t('menu.feedViewer')"
      :description="t('feedViewer.description')"
    >
    </x-header>

    <a-card class="my-2" :title="t('feedViewer.inputLink')">
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
        class="mt-3 grid grid-cols-[minmax(20rem,1fr)_auto] items-center gap-3 max-md:grid-cols-1"
      >
        <a-input
          v-if="previewMode === 'url'"
          v-model="feedUrl"
          type="text"
          class="w-full min-w-0"
          :placeholder="t('feedViewer.placeholder')"
          allow-clear
          @input="clearPreviewState"
          @keyup.enter="fetchFeed"
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
          @keyup.enter="fetchFeed"
        />
        <a-button
          :loading="isLoading"
          :disabled="!currentInputURI"
          @click="fetchFeed"
          >{{ t('feedViewer.preview') }}</a-button
        >
      </div>
      <a-alert v-if="currentInputURI" type="info" class="mt-3" show-icon>
        {{ t('feedViewer.currentInput') }}: {{ currentInputURI }}
      </a-alert>
    </a-card>
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
  </div>
</template>

<script lang="ts" setup>
  import { computed, onMounted, ref } from 'vue';
  import FeedViewContainer from '@/views/dashboard/feed_viewer/feed_view_container.vue';
  import XHeader from '@/components/header/x-header.vue';
  import { useI18n } from 'vue-i18n';
  import { useRoute } from 'vue-router';
  import { previewFeed, type FeedViewerPreview } from '@/api/feed_viewer';
  import { getCustomRecipes, type CustomRecipe } from '@/api/custom_recipe';
  import { listTopicFeeds, type TopicFeed } from '@/api/topic';
  import { listInboxes, type Inbox } from '@/api/inbox';
  import { Message } from '@arco-design/web-vue';

  type PreviewMode = 'url' | 'recipe' | 'topic' | 'inbox' | 'uri';

  const { t } = useI18n();
  const route = useRoute();

  const previewMode = ref<PreviewMode>('url');
  const feedUrl = ref('');
  const advancedURI = ref('');
  const selectedRecipeId = ref('');
  const selectedTopicId = ref('');
  const selectedInboxId = ref('');
  const recipes = ref<CustomRecipe[]>([]);
  const topics = ref<TopicFeed[]>([]);
  const inboxes = ref<Inbox[]>([]);
  const feedContent = ref<FeedViewerPreview | null>(null);
  const errorMessage = ref('');
  const isLoading = ref(false);
  const resourceLoading = ref(false);

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
    errorMessage.value = '';
    feedContent.value = null;
  }

  async function fetchFeed() {
    if (!currentInputURI.value) return;
    isLoading.value = true;
    errorMessage.value = '';
    try {
      const response = await previewFeed(currentInputURI.value);
      feedContent.value = response.data;
    } catch (error) {
      feedContent.value = null;
      errorMessage.value =
        error instanceof Error
          ? error.message
          : t('feedViewer.message.unknownError');
    } finally {
      isLoading.value = false;
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

    selectInputURI(
      firstQueryValue(
        route.query.input_uri || route.query.uri || route.query.url
      )
    );
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
      fetchFeed();
    }
  });
</script>

<script lang="ts">
  export default {
    name: 'FeedViewer',
  };
</script>
