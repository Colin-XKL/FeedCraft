<template>
  <div class="preview-container">
    <iframe
      v-if="htmlContent"
      ref="iframeRef"
      class="html-preview"
      :srcdoc="htmlContent"
      @load="onIframeLoad"
    ></iframe>
    <a-empty v-else description="No content loaded" class="empty-state" />
  </div>
</template>

<script lang="ts" setup>
  import { ref, onBeforeUnmount, watch } from 'vue';

  const props = defineProps({
    htmlContent: {
      type: String,
      default: '',
    },
    isSelectionMode: {
      type: Boolean,
      default: true,
    },
  });

  const emit = defineEmits(['select']);

  const iframeRef = ref<HTMLIFrameElement | null>(null);
  const currentHoverEl = ref<HTMLElement | null>(null);
  const iframeDoc = ref<Document | null>(null);
  const injectedStyle = ref<HTMLStyleElement | null>(null);

  const updateHighlight = (target: HTMLElement) => {
    const doc = iframeRef.value?.contentDocument;
    if (doc) {
      const old = doc.querySelector('.fc-highlight');
      if (old) old.classList.remove('fc-highlight');
    }
    target.classList.add('fc-highlight');
    currentHoverEl.value = target;
  };

  const handleMouseOver = (e: Event) => {
    if (!props.isSelectionMode) return;
    let target = e.target as HTMLElement;
    // Handle text nodes
    if (target && target.nodeType === 3)
      target = target.parentElement as HTMLElement;

    if (
      target &&
      target.nodeType === Node.ELEMENT_NODE &&
      target !== currentHoverEl.value
    ) {
      updateHighlight(target);
    }
  };

  const handleKeyDown = (e: KeyboardEvent) => {
    if (!props.isSelectionMode || !currentHoverEl.value) return;

    if (e.key === 'ArrowUp') {
      e.preventDefault();
      const parent = currentHoverEl.value.parentElement;
      if (parent && parent.tagName !== 'BODY' && parent.tagName !== 'HTML') {
        updateHighlight(parent);
      }
    } else if (e.key === 'ArrowDown') {
      e.preventDefault();
      const child = currentHoverEl.value.firstElementChild as HTMLElement;
      if (child) {
        updateHighlight(child);
      }
    }
  };

  const handleClick = (e: Event) => {
    if (!props.isSelectionMode) return;
    e.preventDefault();
    e.stopPropagation();

    // Use currently highlighted element (allows keyboard navigation selection)
    let target = currentHoverEl.value || (e.target as HTMLElement);
    if (target && target.nodeType === 3)
      target = target.parentElement as HTMLElement;

    if (!target || target.nodeType !== Node.ELEMENT_NODE) return;

    // Emit the raw DOM element to the parent
    emit('select', target);
  };

  const onIframeLoad = () => {
    const iframe = iframeRef.value;
    if (iframe && iframe.contentDocument) {
      const doc = iframe.contentDocument;
      iframeDoc.value = doc;

      const style = doc.createElement('style');
      style.textContent = `
      .fc-highlight { outline: 2px dashed #165dff !important; background-color: rgba(22, 93, 255, 0.05) !important; cursor: pointer; }
    `;
      doc.head.appendChild(style);
      injectedStyle.value = style;

      doc.addEventListener('click', handleClick);
      doc.addEventListener('mouseover', handleMouseOver);
      doc.addEventListener('keydown', handleKeyDown);
      doc.body.setAttribute('tabindex', '0');
    }
  };

  onBeforeUnmount(() => {
    if (iframeDoc.value) {
      iframeDoc.value.removeEventListener('click', handleClick);
      iframeDoc.value.removeEventListener('mouseover', handleMouseOver);
      iframeDoc.value.removeEventListener('keydown', handleKeyDown);
    }
  });

  // Expose the document to parent if needed (e.g. for querySelectorAll validation)
  defineExpose({
    contentDocument: iframeDoc,
  });
</script>

<style scoped>
  .preview-container {
    overflow: hidden;
    position: relative;
    height: 100%;
    border: 1px solid #e5e6eb;
    background: white;
  }

  .html-preview {
    border: none;
    width: 100%;
    height: 100%;
  }

  .empty-state {
    margin-top: 20%;
  }
</style>
