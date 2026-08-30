<template>
  <div class="list-empty-guide">
    <a-empty :description="description" />
    <div class="list-empty-guide__extra">
      <p v-if="hint" class="list-empty-guide__hint">{{ hint }}</p>
      <a-space>
        <a-button type="primary" @click="$emit('create')">
          <template #icon>
            <icon-plus />
          </template>
          {{ createLabel }}
        </a-button>
        <a-link :href="docsHref" target="_blank" rel="noopener noreferrer" icon>
          {{ docsLabel }}
        </a-link>
      </a-space>
    </div>
  </div>
</template>

<script lang="ts" setup>
  defineProps<{
    description: string;
    hint?: string;
    createLabel: string;
    docsLabel: string;
    docsHref: string;
  }>();

  defineEmits<{
    (e: 'create'): void;
  }>();
</script>

<script lang="ts">
  export default {
    name: 'ListEmptyGuide',
  };
</script>

<style scoped lang="less">
  .list-empty-guide {
    padding: 3rem 1rem;
  }

  .list-empty-guide__extra {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.75rem;
  }

  .list-empty-guide__hint {
    max-width: 28rem;
    margin: 0;
    color: var(--color-text-3);
    font-size: 0.875rem;
    line-height: 1.5;
  }
</style>
