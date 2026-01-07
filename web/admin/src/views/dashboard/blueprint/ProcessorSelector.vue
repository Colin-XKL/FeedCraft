<template>
  <div class="processor-selector-wrapper">
    <!-- Single Mode -->
    <SingleProcessorSelector
      v-if="mode === 'single' && !$slots.default"
      :model-value="modelValue as string"
      :placeholder="placeholder"
      :allow-clear="allowClear"
      @update:model-value="onUpdate"
    />

    <!-- Multiple Mode -->
    <MultiProcessorSelector
      v-else-if="mode === 'multiple' && !$slots.default"
      :model-value="modelValue as string[]"
      :placeholder="placeholder"
      @update:model-value="onUpdate"
    />

    <!-- Custom Slot Trigger -->
    <div v-else class="custom-trigger" @click="openModal">
      <slot></slot>
      <ProcessorPickerModal
        v-model:visible="visible"
        :model-value="modelValue"
        :mode="mode"
        :title="placeholder"
        @update:model-value="onUpdate"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref } from 'vue';
  import SingleProcessorSelector from './components/SingleProcessorSelector.vue';
  import MultiProcessorSelector from './components/MultiProcessorSelector.vue';
  import ProcessorPickerModal from './components/ProcessorPickerModal.vue';

  const props = defineProps({
    modelValue: {
      type: [String, Array],
      default: () => [],
    },
    mode: {
      type: String as () => 'single' | 'multiple',
      default: 'single',
    },
    placeholder: {
      type: String,
      default: '',
    },
    allowClear: {
      type: Boolean,
      default: false,
    },
  });

  const emit = defineEmits(['update:modelValue', 'change']);

  const visible = ref(false);

  const openModal = () => {
    visible.value = true;
  };

  const onUpdate = (val: string | string[]) => {
    emit('update:modelValue', val);
    emit('change', val);
  };
</script>

<style scoped>
  .processor-selector-wrapper {
    width: 100%;
  }
  .custom-trigger {
    width: 100%;
    cursor: pointer;
  }
</style>
