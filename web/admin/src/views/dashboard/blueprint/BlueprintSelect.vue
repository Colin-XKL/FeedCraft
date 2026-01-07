<template>
  <ProcessorSelector
    :model-value="modelValue"
    :mode="mode"
    :placeholder="t('feedCompare.selectBlueprint.placeholder')"
    allow-clear
    @update:model-value="onUpdate"
  />
</template>

<script setup lang="ts">
  import { useI18n } from 'vue-i18n';
  import ProcessorSelector from './ProcessorSelector.vue';

  const { t } = useI18n();
  const props = defineProps<{
    modelValue: string[] | string;
    mode: 'single' | 'multiple';
  }>();

  const emit = defineEmits<{
    (event: 'update:modelValue', value: string | string[]): void;
  }>();

  const onUpdate = (val: string | string[]) => {
    emit('update:modelValue', val);
  };
</script>
