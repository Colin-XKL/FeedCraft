<template>
  <div class="py-8 px-[clamp(20px,4vw,64px)] max-md:py-6 max-md:px-4">
    <x-header
      :title="t('embeddingFilterDebug.title')"
      :description="t('embeddingFilterDebug.description')"
    ></x-header>

    <a-card class="my-2" :title="t('embeddingFilterDebug.config')">
      <a-form :model="form" layout="vertical">
        <a-row :gutter="[16, 12]">
          <a-col :span="12" :xs="24" :lg="12">
            <a-form-item :label="t('embeddingFilterDebug.feedUrl')" required>
              <a-input
                v-model="form.inputUrl"
                :placeholder="t('embeddingFilterDebug.feedUrlPlaceholder')"
                allow-clear
              />
            </a-form-item>
          </a-col>
          <a-col :span="12" :xs="24" :lg="12">
            <a-form-item :label="t('embeddingFilterDebug.atomName')">
              <a-input
                v-model="form.atomName"
                :placeholder="t('embeddingFilterDebug.atomNamePlaceholder')"
                allow-clear
              />
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item :label="t('embeddingFilterDebug.anchors')" required>
              <a-textarea
                v-model="form.anchors"
                :placeholder="t('embeddingFilterDebug.anchorsPlaceholder')"
                :auto-size="{ minRows: 4, maxRows: 8 }"
              />
            </a-form-item>
          </a-col>
          <a-col :span="6" :xs="24" :md="8" :lg="6">
            <a-form-item :label="t('embeddingFilterDebug.mode')">
              <a-select v-model="form.mode" :options="modeOptions" />
            </a-form-item>
          </a-col>
          <a-col :span="6" :xs="24" :md="8" :lg="6">
            <a-form-item :label="t('embeddingFilterDebug.threshold')">
              <a-input-number
                v-model="form.threshold"
                :min="0"
                :max="1"
                :step="0.05"
                :precision="2"
                class="w-full"
              />
            </a-form-item>
          </a-col>
          <a-col :span="6" :xs="24" :md="8" :lg="6">
            <a-form-item :label="t('embeddingFilterDebug.maxContentLength')">
              <a-input-number
                v-model="form.maxContentLength"
                :min="1"
                :step="100"
                class="w-full"
              />
            </a-form-item>
          </a-col>
          <a-col :span="6" :xs="24" :lg="6">
            <a-form-item :label="t('embeddingFilterDebug.atomDescription')">
              <a-input
                v-model="form.atomDescription"
                :placeholder="
                  t('embeddingFilterDebug.atomDescriptionPlaceholder')
                "
                allow-clear
              />
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item :label="t('embeddingFilterDebug.instruction')">
              <a-input
                v-model="form.instruction"
                :placeholder="t('embeddingFilterDebug.instructionPlaceholder')"
                allow-clear
              />
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
      <a-space wrap>
        <a-button type="primary" :loading="isTesting" @click="testFilter">
          {{ t('embeddingFilterDebug.test') }}
        </a-button>
        <a-button :loading="isSaving" @click="saveAtomCraft">
          {{ t('embeddingFilterDebug.save') }}
        </a-button>
        <a-tag v-if="hasCompared" color="green">
          {{
            t('embeddingFilterDebug.summary', {
              kept: filteredFeedContent?.items?.length ?? 0,
              total: originalFeedContent?.items?.length ?? 0,
            })
          }}
        </a-tag>
      </a-space>
    </a-card>

    <a-row :gutter="[24, 24]">
      <a-col :span="12" :xs="24" :lg="12">
        <a-card
          :title="t('embeddingFilterDebug.originalFeed')"
          :loading="isTesting"
        >
          <a-alert v-if="originalFeedError" type="error" class="mb-4" show-icon>
            {{ originalFeedError }}
          </a-alert>
          <FeedViewContainer
            v-if="originalFeedContent"
            :feed-data="originalFeedContent"
          />
          <a-empty v-else-if="!originalFeedError" />
        </a-card>
      </a-col>
      <a-col :span="12" :xs="24" :lg="12">
        <a-card
          :title="t('embeddingFilterDebug.filteredFeed')"
          :loading="isTesting"
        >
          <a-alert v-if="filteredFeedError" type="error" class="mb-4" show-icon>
            {{ filteredFeedError }}
          </a-alert>
          <FeedViewContainer
            v-if="filteredFeedContent"
            :feed-data="filteredFeedContent"
          />
          <a-empty v-else-if="!filteredFeedError" />
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script lang="ts" setup>
  import { computed, reactive, ref } from 'vue';
  import { Message } from '@arco-design/web-vue';
  import XHeader from '@/components/header/x-header.vue';
  import FeedViewContainer from '@/views/dashboard/feed_viewer/feed_view_container.vue';
  import { useI18n } from 'vue-i18n';
  import { previewFeed, type FeedViewerPreview } from '@/api/feed_viewer';
  import {
    previewEmbeddingFilter,
    type EmbeddingFilterPreviewRequest,
  } from '@/api/embedding_filter';
  import { createCraftAtom } from '@/api/craft_atom';

  const { t } = useI18n();

  const modeOptions = [
    { label: 'include', value: 'include' },
    { label: 'exclude', value: 'exclude' },
  ];
  const maxContentLengthLimit = 8000;

  const form = reactive({
    inputUrl: '',
    anchors: '',
    mode: 'include' as 'include' | 'exclude',
    threshold: 0.6,
    maxContentLength: 2000,
    instruction: '',
    atomName: '',
    atomDescription: '',
  });

  const originalFeedContent = ref<FeedViewerPreview | null>(null);
  const filteredFeedContent = ref<FeedViewerPreview | null>(null);
  const originalFeedError = ref('');
  const filteredFeedError = ref('');
  const isTesting = ref(false);
  const isSaving = ref(false);
  const hasCompared = computed(
    () =>
      Boolean(originalFeedContent.value) || Boolean(filteredFeedContent.value)
  );

  function hasRequiredInput() {
    return Boolean(form.inputUrl.trim() && form.anchors.trim());
  }

  function normalizeAnchors() {
    return form.anchors
      .split('\n')
      .map((anchor) => anchor.trim())
      .filter(Boolean)
      .join('\n');
  }

  function buildPreviewRequest(): EmbeddingFilterPreviewRequest | null {
    const threshold = Number(form.threshold);
    const maxContentLength = Number(form.maxContentLength);
    if (!Number.isFinite(threshold) || threshold < 0 || threshold > 1) {
      Message.warning('Threshold must be between 0 and 1.');
      return null;
    }
    if (
      !Number.isFinite(maxContentLength) ||
      maxContentLength < 1 ||
      maxContentLength > maxContentLengthLimit
    ) {
      Message.warning(
        `Max content length must be between 1 and ${maxContentLengthLimit}.`
      );
      return null;
    }

    return {
      input_url: form.inputUrl.trim(),
      anchors: normalizeAnchors(),
      threshold,
      mode: form.mode.trim().toLowerCase() as 'include' | 'exclude',
      max_content_length: maxContentLength,
      instruction: form.instruction.trim(),
    };
  }

  function clearResults() {
    originalFeedError.value = '';
    filteredFeedError.value = '';
    originalFeedContent.value = null;
    filteredFeedContent.value = null;
  }

  function getErrorMessage(error: unknown) {
    return error instanceof Error
      ? error.message
      : t('embeddingFilterDebug.message.unknownError');
  }

  async function testFilter() {
    if (isTesting.value) return;
    if (!hasRequiredInput()) {
      Message.warning(t('embeddingFilterDebug.message.inputRequired'));
      return;
    }
    const previewRequest = buildPreviewRequest();
    if (!previewRequest) return;

    isTesting.value = true;
    clearResults();
    const [originalResult, filteredResult] = await Promise.allSettled([
      previewFeed(previewRequest.input_url),
      previewEmbeddingFilter(previewRequest),
    ]);

    if (originalResult.status === 'fulfilled') {
      originalFeedContent.value = originalResult.value.data;
    } else {
      originalFeedError.value = getErrorMessage(originalResult.reason);
    }

    if (filteredResult.status === 'fulfilled') {
      filteredFeedContent.value = filteredResult.value.data;
    } else {
      filteredFeedError.value = getErrorMessage(filteredResult.reason);
    }
    isTesting.value = false;
  }

  async function saveAtomCraft() {
    if (isSaving.value) return;
    if (!hasRequiredInput()) {
      Message.warning(t('embeddingFilterDebug.message.inputRequired'));
      return;
    }
    if (!form.atomName.trim()) {
      Message.warning(t('embeddingFilterDebug.message.atomNameRequired'));
      return;
    }
    const previewRequest = buildPreviewRequest();
    if (!previewRequest) return;

    isSaving.value = true;
    try {
      await createCraftAtom({
        name: form.atomName.trim(),
        description: form.atomDescription.trim(),
        template_name: 'embedding-filter',
        params: {
          anchors: previewRequest.anchors,
          threshold: String(previewRequest.threshold),
          mode: previewRequest.mode,
          max_content_length: String(previewRequest.max_content_length),
          instruction: previewRequest.instruction || '',
        },
      });
      Message.success(t('embeddingFilterDebug.message.saved'));
    } catch (error) {
      Message.error(getErrorMessage(error));
    } finally {
      isSaving.value = false;
    }
  }
</script>

<script lang="ts">
  export default {
    name: 'EmbeddingFilterDebug',
  };
</script>
