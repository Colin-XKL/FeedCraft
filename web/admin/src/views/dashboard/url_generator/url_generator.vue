<template>
  <div class="container p-4 w-full">
    <div>
      <div class="mb-8 text-2xl">
        <p class="text-gray-700"
          >{{ t('urlGenerator.welcome') }}<br />
          <span
            class="text-4xl font-bold text-sky-700 underline decoration-sky-500 decoration-wavy hover:underline-offset-2 hover:decoration-4"
            >Feed Craft</span
          ></p
        >
      </div>
      <a-card class="rounded-lg w-full min-w-8xl">
        <h1 class="text-2xl mb-2">{{ t('urlGenerator.title') }}</h1>
        <p class="mb-8">
          {{ t('urlGenerator.description') }}
        </p>
        <div class="mb-4">
          <label for="siteSelector" class="mr-2">{{
            t('urlGenerator.selectCraft')
          }}</label>
          <CraftFlowSelect v-model="selectedCraft" mode="single" />
        </div>
        <div class="mb-4">
          <label for="inputUrl" class="mr-2">{{
            t('urlGenerator.inputOriginalUrl')
          }}</label>
          <a-input
            id="inputUrl"
            v-model="inputUrl"
            type="text"
            :placeholder="t('urlGenerator.inputUrlPlaceholder')"
          />
        </div>
        <a-button :loading="isLoading" type="primary" @click="generateUrl"
          >{{ t('urlGenerator.showCraftedUrl') }}
        </a-button>
        <div class="mt-8">
          <label for="resultUrl" class="mr-2">{{
            t('urlGenerator.resultUrl')
          }}</label>
          <span id="resultUrl">{{ resultUrl }}</span>
          <a-button
            v-if="resultUrl"
            id="copyButton"
            class="px-2 py-0.5 rounded ml-0.5"
            @click="copyUrl"
            >{{ copyButtonText }}
          </a-button>
        </div>
      </a-card>

      <!-- Preview Section -->
      <a-card
        v-if="previewFeedContent"
        class="mt-8 rounded-lg w-full min-w-8xl"
        :title="t('menu.feedViewer')"
        :loading="isLoading"
      >
        <FeedViewContainer :feed-data="previewFeedContent" />
      </a-card>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref } from 'vue';
  import Parser from 'rss-parser';
  import { Message } from '@arco-design/web-vue';
  import CraftFlowSelect from '@/views/dashboard/craft_flow/CraftFlowSelect.vue';
  import FeedViewContainer from '@/views/dashboard/feed_viewer/feed_view_container.vue';
  import { useI18n } from 'vue-i18n';

  const { t } = useI18n();

  const selectedCraft = ref<string[]>([]);
  const customCraft = ref('');
  const inputUrl = ref('');
  const resultUrl = ref('');
  const copyButtonText = ref(t('urlGenerator.copyUrl'));
  const isLoading = ref(false);
  const previewFeedContent = ref<any>(null);

  async function fetchFeed(url: string) {
    const parser = new Parser();
    const resp = await fetch(url);
    if (!resp.ok) {
      throw new Error(`HTTP error! status: ${resp.status}`);
    }
    return parser.parseString(await resp.text());
  }

  const generateUrl = async () => {
    const currentSelectedCraft = customCraft.value
      ? customCraft.value
      : selectedCraft.value;

    if (!currentSelectedCraft || currentSelectedCraft.length === 0) {
      Message.warning(t('urlGenerator.selectCraft'));
      return;
    }

    if (!inputUrl.value) {
      Message.warning(t('urlGenerator.inputOriginalUrl'));
      return;
    }

    const baseUrl = import.meta.env.VITE_API_BASE_URL ?? window.location.origin;
    resultUrl.value = `${baseUrl}/craft/${currentSelectedCraft}?input_url=${encodeURIComponent(
      inputUrl.value
    )}`;
    copyButtonText.value = t('urlGenerator.copyUrl');

    // Fetch preview
    isLoading.value = true;
    previewFeedContent.value = null;
    try {
      previewFeedContent.value = await fetchFeed(resultUrl.value);
    } catch (error) {
      console.error(error);
      Message.warning(t('feedCompare.message.unknownError') || 'Error fetching feed');
    } finally {
      isLoading.value = false;
    }
  };

  const copyUrl = () => {
    if (resultUrl.value) {
      navigator.clipboard
        .writeText(resultUrl.value)
        .then(() => {
          copyButtonText.value = t('urlGenerator.copied');
        })
        .catch((err) => {
          console.error('无法复制文本: ', err);
        });
    }
  };
</script>

<style scoped>
  .container {
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 90vh; /* changed from height to min-height to allow scrolling */
    flex-direction: column; /* Ensure vertical stacking if needed, though the inner div handles it */
  }
</style>
