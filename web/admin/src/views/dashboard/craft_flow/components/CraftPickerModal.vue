<template>
  <a-modal
    v-model:visible="internalVisible"
    :title="title || t('feedCompare.selectCraftFlow.placeholder')"
    width="800px"
    :footer="mode === 'multiple'"
    @ok="confirmSelection"
    @cancel="cancelSelection"
  >
    <CraftList
      v-model="localValue"
      :multiple="mode === 'multiple'"
      @change="handleListChange"
    />
  </a-modal>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue';
import { useI18n } from 'vue-i18n';
import CraftList from './CraftList.vue';

const { t } = useI18n();

const props = defineProps({
  visible: {
    type: Boolean,
    default: false,
  },
  modelValue: {
    type: [String, Array],
    default: () => [],
  },
  mode: {
    type: String as () => 'single' | 'multiple',
    default: 'single',
  },
  title: {
    type: String,
    default: '',
  },
});

const emit = defineEmits(['update:visible', 'update:modelValue', 'ok']);

// Internal visible state synced with prop
const internalVisible = computed({
  get: () => props.visible,
  set: (val) => emit('update:visible', val),
});

// Internal selection value
const localValue = ref<string[]>([]);

// Sync localValue with modelValue when modal opens
watch(
  () => props.visible,
  (newVal) => {
    if (newVal) {
      if (props.mode === 'single') {
        localValue.value = props.modelValue ? [props.modelValue as string] : [];
      } else {
        localValue.value = Array.isArray(props.modelValue)
          ? [...props.modelValue]
          : [];
      }
    }
  },
);

const confirmSelection = () => {
  let emitVal: string | string[] = localValue.value;
  if (props.mode === 'single') {
    emitVal = localValue.value.length > 0 ? localValue.value[0] : '';
  }
  emit('update:modelValue', emitVal);
  emit('ok', emitVal);
  internalVisible.value = false;
};

const handleListChange = (newVal: string[]) => {
  // If single mode, we might want to auto-confirm, but let's stick to the original behavior
  // Original behavior: handleSelect -> if single -> confirmSelection()
  if (props.mode === 'single' && newVal.length > 0) {
    localValue.value = newVal; // Update local value
    confirmSelection(); // Auto-confirm
  }
};

const cancelSelection = () => {
  internalVisible.value = false;
};
</script>
