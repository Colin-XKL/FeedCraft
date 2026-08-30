<template>
  <div class="steps-container">
    <div
      v-for="(step, idx) in modelValue"
      :key="`step-${idx}`"
      class="step-card"
    >
      <div class="card-header">
        <span class="card-title">{{
          t('topic.stepIndex', { index: idx + 1 })
        }}</span>
        <a-button
          type="text"
          status="danger"
          class="remove-btn"
          @click="removeStep(idx)"
        >
          <template #icon>
            <icon-delete />
          </template>
          {{ t('topic.removeStep') }}
        </a-button>
      </div>

      <div class="field">
        <label class="field-label">{{ t('topic.stepType') }}</label>
        <a-select
          v-model="step.type"
          class="field-control"
          @change="resetStepValue(idx)"
        >
          <a-option value="deduplicate">
            {{ t('topic.stepType.deduplicate') }}
          </a-option>
          <a-option value="sort">{{ t('topic.stepType.sort') }}</a-option>
          <a-option value="limit">{{ t('topic.stepType.limit') }}</a-option>
        </a-select>
      </div>

      <div v-if="step.type === 'deduplicate'" class="field">
        <label class="field-label">{{ t('topic.stepOption.strategy') }}</label>
        <a-select
          v-model="step.value"
          class="field-control"
          @change="onDeduplicateStrategyChange(idx)"
        >
          <a-option value="by_link">
            {{ t('topic.stepOption.strategy.by_link') }}
          </a-option>
          <a-option value="by_id">
            {{ t('topic.stepOption.strategy.by_id') }}
          </a-option>
          <a-option value="by_title">
            {{ t('topic.stepOption.strategy.by_title') }}
          </a-option>
          <a-option value="by_simhash">
            {{ t('topic.stepOption.strategy.by_simhash') }}
          </a-option>
          <a-option value="by_embedding">
            {{ t('topic.stepOption.strategy.by_embedding') }}
          </a-option>
        </a-select>
      </div>

      <div v-else-if="step.type === 'sort'" class="field">
        <label class="field-label">{{ t('topic.stepOption.sort') }}</label>
        <a-select v-model="step.value" class="field-control">
          <a-option value="date_desc">
            {{ t('topic.stepOption.sort.date_desc') }}
          </a-option>
          <a-option value="date_asc">
            {{ t('topic.stepOption.sort.date_asc') }}
          </a-option>
          <a-option value="quality_desc">
            {{ t('topic.stepOption.sort.quality_desc') }}
          </a-option>
          <a-option value="quality_asc">
            {{ t('topic.stepOption.sort.quality_asc') }}
          </a-option>
        </a-select>
      </div>

      <div v-else class="field">
        <label class="field-label">{{ t('topic.stepOption.max') }}</label>
        <a-input-number
          v-model="step.value"
          :min="1"
          mode="button"
          class="field-control"
        />
      </div>

      <div
        v-if="
          step.type === 'deduplicate' &&
          (step.value === 'by_simhash' || step.value === 'by_embedding')
        "
        class="field"
      >
        <label class="field-label">{{
          t('topic.stepOption.threshold.label')
        }}</label>
        <a-input-number
          v-model="step.threshold"
          :min="0"
          :max="1"
          :step="0.01"
          :precision="2"
          :placeholder="
            step.value === 'by_simhash'
              ? t('topic.stepOption.threshold.simhash')
              : t('topic.stepOption.threshold.embedding')
          "
          class="field-control"
        />
        <p class="step-hint">
          {{ t(`topic.stepOption.strategy.${step.value}.hint`) }}
        </p>
      </div>
    </div>

    <a-button type="dashed" long class="add-btn" @click="addStep">
      <template #icon>
        <icon-plus />
      </template>
      {{ t('topic.addStep') }}
    </a-button>
  </div>
</template>

<script lang="ts" setup>
  import { useI18n } from 'vue-i18n';
  import {
    type StepFormItem,
    type StepType,
    createDefaultStep,
    defaultThreshold,
  } from '@/views/dashboard/topic_feed/topicInputUtils';

  const props = defineProps<{
    modelValue: StepFormItem[];
  }>();

  const emit = defineEmits<{
    (event: 'update:modelValue', value: StepFormItem[]): void;
  }>();

  const { t } = useI18n();

  const updateSteps = (steps: StepFormItem[]) => {
    emit('update:modelValue', steps);
  };

  const addStep = () => {
    updateSteps([...props.modelValue, createDefaultStep()]);
  };

  const removeStep = (idx: number) => {
    const next = [...props.modelValue];
    next.splice(idx, 1);
    updateSteps(next);
  };

  const resetStepValue = (idx: number) => {
    const currentType = props.modelValue[idx].type as StepType;
    const next = [...props.modelValue];
    next[idx] = createDefaultStep(currentType);
    updateSteps(next);
  };

  const onDeduplicateStrategyChange = (idx: number) => {
    const next = [...props.modelValue];
    const strategy = String(next[idx].value);
    next[idx] = {
      ...next[idx],
      threshold: defaultThreshold(strategy),
    };
    updateSteps(next);
  };
</script>

<script lang="ts">
  export default {
    name: 'TopicAggregatorEditor',
  };
</script>

<style scoped>
  .steps-container {
    display: flex;
    flex-direction: column;
    gap: 12px;
    width: 100%;
  }

  .step-card {
    display: flex;
    flex-direction: column;
    gap: 12px;
    background-color: var(--color-fill-1);
    border: 1px solid var(--color-border-1);
    border-radius: 6px;
    padding: 12px 16px;
    transition: all 0.2s ease;
  }

  .step-card:hover {
    border-color: var(--color-border-3);
    background-color: var(--color-fill-2);
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 12px;
  }

  .card-title {
    font-weight: 600;
    color: var(--color-text-1);
  }

  .field {
    display: flex;
    flex-direction: column;
    gap: 6px;
    width: 100%;
  }

  .field-label {
    font-size: 13px;
    color: var(--color-text-2);
    line-height: 1.4;
  }

  .field-control,
  .field-control :deep(.arco-select-view),
  .field-control :deep(.arco-input-number) {
    width: 100%;
  }

  .step-hint {
    margin: 0;
    font-size: 12px;
    color: var(--color-text-3);
    line-height: 1.5;
  }

  .remove-btn {
    flex-shrink: 0;
  }

  .add-btn {
    margin-top: 4px;
  }
</style>
