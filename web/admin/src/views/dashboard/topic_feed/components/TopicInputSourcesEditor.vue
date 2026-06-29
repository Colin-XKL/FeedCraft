<template>
  <div class="input-sources-container">
    <div
      v-for="(source, idx) in modelValue"
      :key="`source-${idx}`"
      class="input-source-card"
      :class="{ 'input-source-card--disabled': source.disabled }"
    >
      <div class="source-header">
        <a-radio-group
          :model-value="source.sourceType"
          type="button"
          @change="(value) => onSourceTypeChange(idx, value as SourceType)"
        >
          <a-radio value="external">{{
            t('topic.sourceType.external')
          }}</a-radio>
          <a-radio value="recipe">Recipe</a-radio>
          <a-radio value="topic">Topic</a-radio>
        </a-radio-group>
        <a-space wrap>
          <div v-if="showDisabledToggle" class="disable-toggle">
            <span class="disable-toggle-label">{{
              t('topic.inputDisabled.label')
            }}</span>
            <a-switch
              :model-value="source.disabled"
              size="small"
              @change="(value) => setDisabled(idx, Boolean(value))"
            />
          </div>
          <a-button
            size="small"
            :disabled="!canPreviewSource(source) || source.disabled"
            :loading="previewingIndex === idx"
            @click="openPreview(idx)"
          >
            {{ t('topic.inputPreview') }}
          </a-button>
          <a-button
            type="text"
            status="danger"
            class="remove-btn"
            @click="removeSource(idx)"
          >
            <template #icon>
              <icon-delete />
            </template>
            {{ t('topic.removeInput') }}
          </a-button>
        </a-space>
      </div>

      <a-input
        v-if="source.sourceType === 'external'"
        v-model="source.externalUrl"
        :placeholder="t('topic.sourceUrl.placeholder')"
        allow-clear
      />
      <a-select
        v-else-if="source.sourceType === 'recipe'"
        v-model="source.resourceId"
        :loading="pickerLoading"
        allow-search
        allow-clear
        :placeholder="t('topic.sourceSelect.placeholder.recipe')"
      >
        <a-option
          v-for="r in availableRecipes"
          :key="r.id"
          :value="r.id"
          :label="r.description ? `${r.id} — ${r.description}` : r.id"
        >
          <span class="option-id">{{ r.id }}</span>
          <span v-if="r.description" class="option-desc">
            — {{ r.description }}
          </span>
        </a-option>
      </a-select>
      <a-select
        v-else
        v-model="source.resourceId"
        :loading="pickerLoading"
        allow-search
        allow-clear
        :placeholder="t('topic.sourceSelect.placeholder.topic')"
      >
        <a-option
          v-for="tp in pickerTopics"
          :key="tp.id"
          :value="tp.id"
          :label="tp.title ? `${tp.id} — ${tp.title}` : tp.id"
        >
          <span class="option-id">{{ tp.id }}</span>
          <span v-if="tp.title" class="option-desc"> — {{ tp.title }}</span>
        </a-option>
      </a-select>

      <a-input
        v-model="source.description"
        :placeholder="t('topic.inputDescription.placeholder')"
        allow-clear
      />
    </div>

    <a-button type="dashed" long class="add-btn" @click="addSource">
      <template #icon>
        <icon-plus />
      </template>
      {{ t('topic.addInput') }}
    </a-button>

    <a-drawer
      v-model:visible="previewDrawerVisible"
      :title="previewDrawerTitle"
      :width="720"
      unmount-on-close
    >
      <a-alert
        v-if="previewInputUri"
        type="info"
        class="mb-4"
        show-icon
        :title="t('topic.inputPreview.currentUri')"
      >
        {{ previewInputUri }}
      </a-alert>
      <FeedPreviewPanel
        ref="previewPanelRef"
        :empty-description="t('topic.inputPreview.empty')"
      />
    </a-drawer>
  </div>
</template>

<script lang="ts" setup>
  import { computed, nextTick, ref } from 'vue';
  import { useI18n } from 'vue-i18n';
  import type { CustomRecipe } from '@/api/custom_recipe';
  import type { TopicFeed } from '@/api/topic';
  import FeedPreviewPanel from '@/components/feed-preview/FeedPreviewPanel.vue';
  import {
    type InputSourceItem,
    type SourceType,
    sourceToUri,
  } from '../topicInputUtils';

  const props = withDefaults(
    defineProps<{
      modelValue: InputSourceItem[];
      availableRecipes: CustomRecipe[];
      availableTopics: TopicFeed[];
      pickerLoading?: boolean;
      excludeTopicId?: string;
      showDisabledToggle?: boolean;
    }>(),
    {
      showDisabledToggle: true,
    }
  );

  const emit = defineEmits<{
    (event: 'update:modelValue', value: InputSourceItem[]): void;
  }>();

  const { t } = useI18n();
  const previewDrawerVisible = ref(false);
  const previewInputUri = ref('');
  const previewDrawerTitle = ref('');
  const previewingIndex = ref<number | null>(null);
  const previewPanelRef = ref<InstanceType<typeof FeedPreviewPanel> | null>(
    null
  );

  const pickerTopics = computed(() => {
    if (!props.excludeTopicId) return props.availableTopics;
    return props.availableTopics.filter((tp) => tp.id !== props.excludeTopicId);
  });

  const updateSources = (sources: InputSourceItem[]) => {
    emit('update:modelValue', sources);
  };

  const canPreviewSource = (source: InputSourceItem) => {
    return sourceToUri(source) !== '';
  };

  const addSource = () => {
    updateSources([
      ...props.modelValue,
      {
        sourceType: 'external',
        externalUrl: '',
        resourceId: '',
        description: '',
        disabled: false,
      },
    ]);
  };

  const setDisabled = (idx: number, disabled: boolean) => {
    const next = [...props.modelValue];
    next[idx] = { ...next[idx], disabled };
    updateSources(next);
  };

  const removeSource = (idx: number) => {
    const next = [...props.modelValue];
    next.splice(idx, 1);
    if (next.length === 0) {
      next.push({
        sourceType: 'external',
        externalUrl: '',
        resourceId: '',
        description: '',
        disabled: false,
      });
    }
    updateSources(next);
  };

  const onSourceTypeChange = (idx: number, sourceType: SourceType) => {
    const next = [...props.modelValue];
    const current = next[idx];
    next[idx] = {
      sourceType,
      externalUrl: '',
      resourceId: '',
      description: current?.description || '',
      disabled: current?.disabled || false,
    };
    updateSources(next);
  };

  const openPreview = async (idx: number) => {
    const source = props.modelValue[idx];
    const uri = sourceToUri(source);
    if (!uri) return;

    previewingIndex.value = idx;
    previewInputUri.value = uri;
    previewDrawerTitle.value =
      source.description.trim() ||
      t('topic.inputPreview.title', { index: idx + 1 });
    previewDrawerVisible.value = true;

    await nextTick();
    try {
      await previewPanelRef.value?.loadPreview(uri);
    } finally {
      previewingIndex.value = null;
    }
  };
</script>

<script lang="ts">
  export default {
    name: 'TopicInputSourcesEditor',
  };
</script>

<style scoped>
  .input-sources-container {
    display: flex;
    flex-direction: column;
    gap: 12px;
    width: 100%;
  }

  .input-source-card {
    display: flex;
    flex-direction: column;
    gap: 12px;
    background-color: var(--color-fill-1);
    border: 1px solid var(--color-border-1);
    border-radius: 6px;
    padding: 12px 16px;
    transition: all 0.2s ease;
  }

  .input-source-card:hover {
    border-color: var(--color-border-3);
    background-color: var(--color-fill-2);
  }

  .input-source-card--disabled {
    opacity: 0.72;
  }

  .disable-toggle {
    display: inline-flex;
    align-items: center;
    gap: 8px;
  }

  .disable-toggle-label {
    font-size: 12px;
    color: var(--color-text-3);
  }

  .source-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 12px;
    flex-wrap: wrap;
  }

  .remove-btn {
    flex-shrink: 0;
  }

  .add-btn {
    margin-top: 4px;
  }

  .option-id {
    font-weight: 500;
  }

  .option-desc {
    color: var(--color-text-3);
    margin-left: 2px;
    font-size: 12px;
  }
</style>
