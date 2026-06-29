<template>
  <div class="steps-container">
    <div
      v-for="(step, idx) in modelValue"
      :key="`step-${idx}`"
      class="step-card"
    >
      <div class="editor-row">
        <a-select
          v-model="step.type"
          style="width: 150px; flex-shrink: 0"
          @change="resetStepValue(idx)"
        >
          <a-option value="deduplicate">
            {{ t('topic.stepType.deduplicate') }}
          </a-option>
          <a-option value="sort">{{ t('topic.stepType.sort') }}</a-option>
          <a-option value="limit">{{ t('topic.stepType.limit') }}</a-option>
        </a-select>

        <template v-if="step.type === 'deduplicate'">
          <a-select
            v-model="step.value"
            style="width: 220px; flex-shrink: 0"
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
          <a-input-number
            v-if="step.value === 'by_simhash'"
            v-model="step.threshold"
            :min="0"
            :max="1"
            :step="0.01"
            :precision="2"
            :placeholder="t('topic.stepOption.threshold.simhash')"
            style="flex: 1; min-width: 120px"
          />
          <a-input-number
            v-else-if="step.value === 'by_embedding'"
            v-model="step.threshold"
            :min="0"
            :max="1"
            :step="0.01"
            :precision="2"
            :placeholder="t('topic.stepOption.threshold.embedding')"
            style="flex: 1; min-width: 120px"
          />
        </template>

        <a-select
          v-else-if="step.type === 'sort'"
          v-model="step.value"
          style="flex: 1; min-width: 150px"
        >
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

        <a-input-number
          v-else
          :model-value="toNumberValue(step.value)"
          :min="1"
          mode="button"
          style="flex: 1; min-width: 150px"
          @update:model-value="(value) => updateStepValue(idx, value)"
        />

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
      <p
        v-if="
          step.type === 'deduplicate' &&
          (step.value === 'by_simhash' || step.value === 'by_embedding')
        "
        class="step-hint"
      >
        {{ t(`topic.stepOption.strategy.${step.value}.hint`) }}
      </p>
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
  } from '../topicInputUtils';

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

  const toNumberValue = (value: string | number) => {
    if (typeof value === 'number') return value;
    const parsed = Number(value);
    return Number.isNaN(parsed) ? undefined : parsed;
  };

  const updateStepValue = (
    idx: number,
    value: string | number | undefined
  ) => {
    const next = [...props.modelValue];
    next[idx] = {
      ...next[idx],
      value: value ?? 1,
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
    gap: 8px;
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

  .editor-row {
    display: flex;
    gap: 12px;
    align-items: center;
    width: 100%;
  }

  .step-hint {
    margin: 4px 0 0 0;
    padding-left: 4px;
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
