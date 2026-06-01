<template>
  <!-- Normal text mode -->
  <div v-if="viewMode === 'normal'" class="feed-content-normal">
    <p
      v-if="snippetText"
      class="text-sm leading-relaxed whitespace-pre-line m-0"
    >
      {{ snippetText }}
    </p>
    <p v-else class="text-sm text-gray-400 m-0">
      {{ t('feedViewer.item.noContent') }}
    </p>
  </div>

  <!-- Rich text mode -->
  <!-- eslint-disable-next-line vue/no-v-html -->
  <div v-else-if="viewMode === 'rich'" class="feed-content-rich">
    <!-- eslint-disable-next-line vue/no-v-html -->
    <div
      v-if="sanitizedContent"
      class="rich-text-body"
      v-html="sanitizedContent"
    />
    <p v-else class="text-sm text-gray-400 m-0">
      {{ t('feedViewer.item.noContent') }}
    </p>
  </div>

  <!-- Raw HTML mode -->
  <div v-else class="feed-content-html">
    <div v-if="rawHtml" class="html-code-wrapper">
      <!-- eslint-disable-next-line vue/no-v-html -->
      <pre
        class="html-code-block"
      ><code class="language-html" v-html="highlightedHtml" /></pre>
    </div>
    <p v-else class="text-sm text-gray-400 m-0">
      {{ t('feedViewer.item.noContent') }}
    </p>
  </div>
</template>

<script lang="ts" setup>
  import { computed } from 'vue';
  import DOMPurify from 'dompurify';
  import hljs from 'highlight.js/lib/core';
  import xml from 'highlight.js/lib/languages/xml';
  import { useI18n } from 'vue-i18n';
  import type { FeedViewerPreviewItem } from '@/api/feed_viewer';
  import formatHTML from '@/utils/htmlFormat';

  hljs.registerLanguage('xml', xml);

  export type ViewMode = 'normal' | 'rich' | 'html';

  interface Props {
    item: FeedViewerPreviewItem;
    viewMode: ViewMode;
  }

  const props = defineProps<Props>();
  const { t } = useI18n();

  function escapeHtml(str: string) {
    return str
      .replace(/&/g, '&amp;')
      .replace(/</g, '&lt;')
      .replace(/>/g, '&gt;')
      .replace(/"/g, '&quot;')
      .replace(/'/g, '&#039;');
  }

  const snippetText = computed(() => props.item.contentSnippet?.trim() || '');

  const sanitizedContent = computed(() => {
    const raw = props.item.content || props.item.contentSnippet || '';
    if (!raw.trim()) return '';
    return DOMPurify.sanitize(raw);
  });

  const rawHtml = computed(
    () => props.item.content?.trim() || props.item.contentSnippet?.trim() || ''
  );

  const highlightedHtml = computed(() => {
    if (!rawHtml.value) return '';
    try {
      const formatted = formatHTML(rawHtml.value);
      return hljs.highlight(formatted, { language: 'xml' }).value;
    } catch {
      return escapeHtml(rawHtml.value);
    }
  });
</script>

<style scoped>
  .feed-content-rich .rich-text-body {
    font-size: 14px;
    line-height: 1.7;
    color: var(--color-text-1);
    overflow-wrap: break-word;
    word-break: break-word;
  }

  .feed-content-rich .rich-text-body :deep(img) {
    max-width: 100%;
    height: auto;
    display: block;
    margin: 8px 0;
  }

  .feed-content-rich .rich-text-body :deep(a) {
    color: var(--color-primary-6);
    text-decoration: underline;
  }

  .feed-content-rich .rich-text-body :deep(p) {
    margin: 0.5em 0;
  }

  .feed-content-rich .rich-text-body :deep(h1),
  .feed-content-rich .rich-text-body :deep(h2),
  .feed-content-rich .rich-text-body :deep(h3),
  .feed-content-rich .rich-text-body :deep(h4) {
    margin: 0.75em 0 0.25em;
    font-weight: 600;
  }

  .feed-content-rich .rich-text-body :deep(pre) {
    background: var(--color-fill-2);
    border-radius: 4px;
    padding: 12px;
    overflow-x: auto;
    white-space: pre-wrap;
    word-wrap: break-word;
    font-size: 13px;
  }

  .feed-content-rich .rich-text-body :deep(blockquote) {
    border-left: 3px solid var(--color-border-2);
    margin: 0.5em 0;
    padding: 0 12px;
    color: var(--color-text-3);
  }

  .feed-content-rich .rich-text-body :deep(ul),
  .feed-content-rich .rich-text-body :deep(ol) {
    padding-left: 1.5em;
    margin: 0.5em 0;
  }

  .feed-content-rich .rich-text-body :deep(table) {
    border-collapse: collapse;
    width: 100%;
    font-size: 13px;
  }

  .feed-content-rich .rich-text-body :deep(th),
  .feed-content-rich .rich-text-body :deep(td) {
    border: 1px solid var(--color-border-1);
    padding: 6px 10px;
  }

  .html-code-wrapper {
    border: 1px solid var(--color-border-1);
    border-radius: 6px;
    overflow: hidden;
  }

  .html-code-block {
    margin: 0;
    padding: 12px 16px;
    font-size: 12.5px;
    line-height: 1.6;
    font-family: 'Fira Code', 'Cascadia Code', 'JetBrains Mono', Consolas,
      'Courier New', monospace;
    overflow-x: auto;
    overflow-y: auto;
    max-height: 500px;
    white-space: pre;
    background: var(--color-fill-1);
    color: var(--color-text-1);
  }

  .html-code-block :deep(.hljs-tag) {
    color: #569cd6;
  }

  .html-code-block :deep(.hljs-attr) {
    color: #9cdcfe;
  }

  .html-code-block :deep(.hljs-string) {
    color: #ce9178;
  }

  .html-code-block :deep(.hljs-comment) {
    color: #6a9955;
    font-style: italic;
  }

  .html-code-block :deep(.hljs-name) {
    color: #4ec9b0;
  }

  .html-code-block :deep(.hljs-keyword) {
    color: #c586c0;
  }

  .html-code-block :deep(.hljs-meta) {
    color: #9b9b9b;
  }
</style>
