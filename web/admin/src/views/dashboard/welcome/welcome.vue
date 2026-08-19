<template>
  <div class="welcome-page">
    <div class="welcome-content">
      <img :src="logo" alt="FeedCraft Logo" class="welcome-logo" />
      <h1 class="text-3xl my-4 font-light">
        Welcome To
        <span class="font-bold">
          <span class="welcome-brand-feed">Feed</span
          ><span class="welcome-brand-craft">Craft</span>
        </span>
      </h1>
      <p class="text-xl">{{ t('welcome.subtitle') }}</p>
      <a-row :gutter="20" class="my-8 welcome-link-row">
        <a-col :span="8">
          <a-card :title="t('welcome.card.home')" :bordered="false" hoverable>
            <a-link
              href="https://github.com/Colin-XKL/FeedCraft"
              target="_blank"
              rel="noopener noreferrer"
              icon
            >
              Colin-XKL/FeedCraft
            </a-link>
          </a-card>
        </a-col>
        <a-col :span="8">
          <a-card
            :title="t('welcome.card.moreRss')"
            :bordered="false"
            hoverable
          >
            <a-link
              href="https://docs.rsshub.app"
              target="_blank"
              rel="noopener noreferrer"
              icon
            >
              RSSHub Doc
            </a-link>
          </a-card>
        </a-col>
        <a-col :span="8">
          <a-card
            :title="t('welcome.card.learnRss')"
            :bordered="false"
            hoverable
          >
            <a-link
              href="https://sspai.com/post/56391"
              target="_blank"
              rel="noopener noreferrer"
              icon
            >
              RSS 入门指南 - sspai
            </a-link>
          </a-card>
        </a-col>
      </a-row>
      <a-row :gutter="20">
        <a-col :xs="24" :lg="16">
          <a-card :title="t('welcome.quickStart')" :bordered="false" hoverable>
            <p class="text-gray-600 mb-4">{{ t('welcome.quickStart.tip') }}</p>
            <div class="mb-4">
              <label class="block mb-2 font-medium" for="welcomeRssUrl">{{
                t('welcome.quickStart.rssUrl')
              }}</label>
              <div class="flex flex-wrap gap-2">
                <a-input
                  id="welcomeRssUrl"
                  v-model="inputUrl"
                  class="flex-1 min-w-[16rem]"
                  :placeholder="t('welcome.quickStart.rssPlaceholder')"
                  allow-clear
                  @keyup.enter="generateUrl"
                />
                <a-button @click="fillExample">{{
                  t('welcome.quickStart.useExample')
                }}</a-button>
              </div>
            </div>
            <div class="mb-4">
              <label class="block mb-2 font-medium">{{
                t('welcome.quickStart.craft')
              }}</label>
              <a-select
                v-model="selectedCraft"
                :placeholder="t('welcome.quickStart.craft')"
                allow-search
              >
                <a-optgroup
                  v-for="group in WELCOME_CRAFT_GROUPS"
                  :key="group.id"
                  :label="t(group.labelKey)"
                >
                  <a-option
                    v-for="craft in craftsInGroup(group.id)"
                    :key="craft.value"
                    :value="craft.value"
                    :label="`${craft.value} - ${t(craft.labelKey)}`"
                  />
                </a-optgroup>
              </a-select>
            </div>
            <a-space wrap>
              <a-button type="primary" @click="generateUrl">
                {{ t('welcome.quickStart.generate') }}
              </a-button>
              <router-link :to="{ name: 'QuickStartFeedCraftUrlGenerator' }">
                <a-button type="text">
                  {{ t('welcome.quickStart.fullGenerator') }}
                </a-button>
              </router-link>
            </a-space>
            <div v-if="resultUrl" class="result-panel mt-6">
              <div class="flex items-center justify-between gap-3 mb-3">
                <span class="font-medium">{{
                  t('welcome.quickStart.result')
                }}</span>
              </div>
              <a
                :href="resultUrl"
                class="result-url"
                target="_blank"
                rel="noopener noreferrer"
              >
                {{ resultUrl }}
              </a>
              <div class="mt-3 flex flex-wrap items-center gap-3">
                <a-button type="primary" @click="copyUrl">
                  {{ copyButtonText }}
                </a-button>
                <a-button @click="previewUrl">
                  {{ t('welcome.quickStart.preview') }}
                </a-button>
              </div>
            </div>
            <a-divider />
            <p class="mb-2 font-medium">{{
              t('welcome.quickStart.availableCrafts')
            }}</p>
            <div class="grid grid-cols-2 gap-2">
              <div v-for="group in WELCOME_CRAFT_GROUPS" :key="group.id">
                <h3 class="font-bold">{{ t(group.labelKey) }}</h3>
                <ul>
                  <li
                    v-for="craft in craftsInGroup(group.id)"
                    :key="craft.value"
                  >
                    <b>{{ craft.value }}</b> - {{ t(craft.labelKey) }}
                  </li>
                </ul>
              </div>
            </div>
          </a-card>
        </a-col>
        <a-col :xs="24" :lg="8">
          <a-card
            hoverable
            :title="t('welcome.feedback')"
            :bordered="false"
            :style="{ width: '100%' }"
          >
            <a-link
              href="https://github.com/Colin-XKL/FeedCraft/discussions"
              target="_blank"
              rel="noopener noreferrer"
              icon
            >
              Github Discussion
            </a-link>
          </a-card>
        </a-col>
      </a-row>
    </div>
  </div>
</template>

<script lang="ts" setup>
  import { ref, watch } from 'vue';
  import { useRouter } from 'vue-router';
  import { useI18n } from 'vue-i18n';
  import { Message } from '@arco-design/web-vue';
  import { useClipboard } from '@vueuse/core';
  import logo from '@/assets/logo.png';
  import { buildCraftFeedUrl, isHttpUrl } from '@/utils/publicFeedUrl';
  import {
    DEFAULT_WELCOME_CRAFT,
    EXAMPLE_RSS_URL,
    WELCOME_CRAFT_GROUPS,
    craftsInGroup,
  } from '@/views/dashboard/welcome/welcomeCrafts';

  const { t } = useI18n();
  const router = useRouter();
  const selectedCraft = ref(DEFAULT_WELCOME_CRAFT);
  const inputUrl = ref('');
  const resultUrl = ref('');
  const generatedCraft = ref('');
  const generatedInputUrl = ref('');
  const copyButtonText = ref(t('urlGenerator.copyUrl'));
  const { copy: copyResultUrl } = useClipboard({
    source: resultUrl,
    legacy: true,
    copiedDuring: 1500,
  });

  const fillExample = () => {
    inputUrl.value = EXAMPLE_RSS_URL;
  };

  const generateUrl = () => {
    const craft = selectedCraft.value.trim();
    const source = inputUrl.value.trim();
    if (!craft) {
      Message.warning(t('welcome.quickStart.validation.craft'));
      return;
    }
    if (!isHttpUrl(source)) {
      Message.warning(t('welcome.quickStart.validation.url'));
      return;
    }
    resultUrl.value = buildCraftFeedUrl(craft, source);
    generatedCraft.value = craft;
    generatedInputUrl.value = source;
    copyButtonText.value = t('urlGenerator.copyUrl');
  };

  watch([inputUrl, selectedCraft], () => {
    copyButtonText.value = t('urlGenerator.copyUrl');
    if (
      resultUrl.value &&
      selectedCraft.value.trim() &&
      isHttpUrl(inputUrl.value) &&
      (selectedCraft.value.trim() !== generatedCraft.value ||
        inputUrl.value.trim() !== generatedInputUrl.value)
    ) {
      generateUrl();
    }
  });

  const copyUrl = async () => {
    if (!resultUrl.value) {
      return;
    }
    try {
      await copyResultUrl();
      copyButtonText.value = t('urlGenerator.copied');
    } catch {
      Message.error(t('urlGenerator.copyError'));
    }
  };

  const previewUrl = () => {
    if (!resultUrl.value) {
      return;
    }
    router.push({
      name: 'FeedViewer',
      query: { url: resultUrl.value },
    });
  };
</script>

<script lang="ts">
  export default {
    name: 'WelcomePage',
  };
</script>

<style scoped>
  .welcome-page {
    width: 100%;
    min-height: calc(100vh - 80px);
    display: flex;
    align-items: flex-start;
    justify-content: center;
    padding: 2rem;
  }

  .welcome-content {
    width: min(72rem, 100%);
  }

  .welcome-logo {
    height: 8rem;
    width: 8rem;
  }

  .welcome-brand-feed {
    color: #4d4f5a;
  }

  .welcome-brand-craft {
    color: #0d9488;
  }

  .welcome-link-row {
    min-width: 0;
  }

  .result-panel {
    border: 1px solid var(--color-neutral-3);
    border-radius: 12px;
    background: rgba(var(--green-1), 0.45);
    padding: 16px;
  }

  .result-url {
    display: block;
    min-height: 40px;
    border-radius: 8px;
    background: var(--color-bg-2);
    color: rgb(var(--link-6));
    overflow-wrap: anywhere;
    padding: 10px 12px;
  }
</style>
