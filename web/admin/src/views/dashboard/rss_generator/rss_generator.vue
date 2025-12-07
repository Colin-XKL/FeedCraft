<template>
  <div class="container">
    <Breadcrumb :items="['menu.tools', 'menu.rssGenerator']" />
    <div class="content-wrapper">
      <!-- Top Bar: URL Input -->
      <a-card class="mb-4">
        <a-space>
          <a-input
            v-model="url"
            placeholder="Enter URL to fetch"
            style="width: 400px"
            @keyup.enter="handleFetch"
          />
          <a-button type="primary" :loading="fetching" @click="handleFetch">
            Fetch Page
          </a-button>
          <a-button @click="toggleMode">
            {{ isSelectionMode ? 'Selection Mode' : 'Preview Mode' }}
          </a-button>
        </a-space>
      </a-card>

      <a-row :gutter="16" class="main-area">
        <!-- Left: HTML Preview / Interaction Area -->
        <a-col :span="14" class="h-full">
          <a-card class="preview-card" title="Page Preview">
            <iframe
              v-if="htmlContent"
              ref="previewIframe"
              class="html-preview"
              :srcdoc="htmlContent"
              @load="onIframeLoad"
            ></iframe>
            <a-empty v-else description="No content loaded" />
          </a-card>
        </a-col>

        <!-- Right: Configuration & Result -->
        <a-col :span="10" class="h-full">
          <a-card class="config-card" title="RSS Configuration">
            <a-form :model="config" layout="vertical">
              <a-form-item label="List Item Selector">
                <a-input v-model="config.item_selector">
                  <template #suffix>
                    <a-button
                      size="mini"
                      type="text"
                      @click="setTargetField('item_selector')"
                    >
                      <icon-select-all /> Pick
                    </a-button>
                  </template>
                </a-input>
              </a-form-item>
              <a-form-item label="Title Selector (Relative)">
                <a-input v-model="config.title_selector">
                  <template #suffix>
                    <a-button
                      size="mini"
                      type="text"
                      @click="setTargetField('title_selector')"
                    >
                      <icon-select-all /> Pick
                    </a-button>
                  </template>
                </a-input>
              </a-form-item>
              <a-form-item label="Link Selector (Relative)">
                <a-input v-model="config.link_selector">
                  <template #suffix>
                    <a-button
                      size="mini"
                      type="text"
                      @click="setTargetField('link_selector')"
                    >
                      <icon-select-all /> Pick
                    </a-button>
                  </template>
                </a-input>
              </a-form-item>
              <a-form-item label="Date Selector (Relative)">
                <a-input v-model="config.date_selector">
                  <template #suffix>
                    <a-button
                      size="mini"
                      type="text"
                      @click="setTargetField('date_selector')"
                    >
                      <icon-select-all /> Pick
                    </a-button>
                  </template>
                </a-input>
              </a-form-item>
              <a-form-item label="Content Selector (Relative)">
                <a-input v-model="config.content_selector">
                  <template #suffix>
                    <a-button
                      size="mini"
                      type="text"
                      @click="setTargetField('content_selector')"
                    >
                      <icon-select-all /> Pick
                    </a-button>
                  </template>
                </a-input>
              </a-form-item>

              <a-space>
                <a-button
                  type="primary"
                  :loading="parsing"
                  @click="handlePreview"
                >
                  Preview RSS Items
                </a-button>
                <a-button @click="clearConfig">Clear</a-button>
              </a-space>
            </a-form>
          </a-card>

          <a-card class="result-card mt-4" title="Extracted Items">
            <div v-if="parsedItems.length > 0">
              <a-collapse>
                <a-collapse-item
                  v-for="(item, idx) in parsedItems"
                  :key="idx"
                  :header="item.title || 'No Title'"
                >
                  <p><strong>Link:</strong> {{ item.link }}</p>
                  <p><strong>Date:</strong> {{ item.date }}</p>
                  <div class="content-preview">
                    <strong>Content Preview:</strong>
                    <div
                      style="
                        max-height: 100px;
                        overflow: hidden;
                        text-overflow: ellipsis;
                      "
                    >
                      {{ item.content }}
                    </div>
                  </div>
                </a-collapse-item>
              </a-collapse>
            </div>
            <a-empty v-else description="No items extracted yet" />
          </a-card>
        </a-col>
      </a-row>
    </div>
  </div>
</template>

<script lang="ts" setup>
  import { ref, reactive } from 'vue';
  import axios from 'axios';
  import DOMPurify from 'dompurify';
  import { Message } from '@arco-design/web-vue';
  import { IconSelectAll } from '@arco-design/web-vue/es/icon';

  const url = ref('');
  const fetching = ref(false);
  const parsing = ref(false);
  const htmlContent = ref('');
  const isSelectionMode = ref(true);
  const previewIframe = ref<HTMLIFrameElement | null>(null);
  const currentHoverEl = ref<HTMLElement | null>(null);

  const config = reactive<{ [key: string]: string }>({
    item_selector: '',
    title_selector: '',
    link_selector: '',
    date_selector: '',
    content_selector: '',
  });

  const parsedItems = ref<any[]>([]);
  const currentTargetField = ref<string>(''); // Field currently being picked

  const toggleMode = () => {
    isSelectionMode.value = !isSelectionMode.value;
  };

  const setTargetField = (field: string) => {
    currentTargetField.value = field;
    Message.info(
      `Please click an element in the preview to select selector for ${field}`
    );
  };

  const handleFetch = async () => {
    if (!url.value) return;
    fetching.value = true;
    try {
      const res = (await axios.post('/api/admin/tools/fetch', {
        url: url.value,
      })) as any;
      if (res.code === 0) {
        let raw = res.data;
        // Inject base tag for relative links
        const baseTag = `<base href="${url.value}" />`;
        // Insert after <head> or at start
        if (raw.toLowerCase().includes('<head>')) {
          raw = raw.replace(/<head>/i, `<head>${baseTag}`);
        } else {
          raw = `${baseTag}${raw}`;
        }

        htmlContent.value = DOMPurify.sanitize(raw, {
          WHOLE_DOCUMENT: true,
          ADD_TAGS: ['link', 'style', 'head', 'meta', 'body', 'html', 'base'],
          ADD_ATTR: ['href', 'rel', 'src', 'type'],
        });
        Message.success('Page fetched successfully');
      } else {
        Message.error(res.msg || 'Fetch failed');
      }
    } catch (err) {
      Message.error('Error fetching page');
    } finally {
      fetching.value = false;
    }
  };

  const handlePreview = async () => {
    if (!config.item_selector) {
      Message.warning('Please set at least the Item Selector');
      return;
    }
    parsing.value = true;
    try {
      const res = (await axios.post('/api/admin/tools/parse', {
        html: htmlContent.value,
        url: url.value,
        ...config,
      })) as any;

      if (res.code === 0) {
        parsedItems.value = res.data;
        Message.success(`Extracted ${parsedItems.value.length} items`);
      } else {
        Message.error(res.msg || 'Parse failed');
      }
    } catch (err) {
      Message.error('Error parsing RSS');
    } finally {
      parsing.value = false;
    }
  };

  const clearConfig = () => {
    config.item_selector = '';
    config.title_selector = '';
    config.link_selector = '';
    config.date_selector = '';
    config.content_selector = '';
    parsedItems.value = [];
  };

  // --- Selector Generation Logic ---

  // Enhanced smart selector generator
  const getCssSelector = (
    el: HTMLElement,
    isItemSelector = false
  ): string => {
    if (!(el instanceof Element)) return '';

    // 1. If ID exists, use it.
    if (el.id) {
      // Unless it's the item selector and we want to select multiples?
      // Usually IDs are unique, so bad for list items.
      if (!isItemSelector) return `#${el.id}`;
    }

    const path: string[] = [];
    let currentEl: HTMLElement | null = el;

    const doc = previewIframe.value?.contentDocument;
    const body = doc?.body;
    const html = doc?.documentElement;

    while (currentEl && currentEl.nodeType === Node.ELEMENT_NODE) {
      let selector = currentEl.nodeName.toLowerCase();

      // Try to use class if available and meaningful
      if (currentEl.classList.length > 0) {
        // Filter out common utility classes (simplified heuristic)
        const classes = Array.from(currentEl.classList).filter(
          (c) =>
            ![
              'flex',
              'row',
              'col',
              'grid',
              'hidden',
              'block',
              'items-center',
              'justify-center',
            ].some((x) => c.includes(x))
        );
        if (classes.length > 0) {
          selector += `.${classes.join('.')}`;
        }
      }

      // If this is the *target* element (first in loop) and we are selecting List Item
      if (currentEl === el && isItemSelector) {
        // Do NOT add :nth-of-type to the target element itself
        // This ensures we match all siblings
      } else if (currentEl.id) {
        // For parent path, or normal selection, use nth-of-type to be precise
        selector = `#${currentEl.id}`;
        path.unshift(selector);
        break;
      } else {
        // check siblings
        let sib = currentEl;
        let nth = 1;
        // eslint-disable-next-line no-cond-assign
        while ((sib = sib.previousElementSibling as HTMLElement)) {
          if (
            sib.nodeName.toLowerCase() === currentEl.nodeName.toLowerCase()
          ) {
            nth += 1;
          }
        }
        if (nth !== 1) selector += `:nth-of-type(${nth})`;
      }

      path.unshift(selector);
      currentEl = currentEl.parentNode as HTMLElement;

      // Stop if we hit the container or root
      if (!currentEl || currentEl === body || currentEl === html) break;
    }

    return path.join(' > ');
  };

  const updateHighlight = (target: HTMLElement) => {
    // Remove old highlight
    const doc = previewIframe.value?.contentDocument;
    if (doc) {
      const old = doc.querySelector('.fc-highlight');
      if (old) old.classList.remove('fc-highlight');
    }

    // Add new highlight
    target.classList.add('fc-highlight');
    currentHoverEl.value = target;
  };

  const handleMouseOver = (e: Event) => {
    if (!isSelectionMode.value) return;
    const target = e.target as HTMLElement;
    if (target && target !== currentHoverEl.value) {
      updateHighlight(target);
    }
  };

  const handleKeyDown = (e: KeyboardEvent) => {
    if (!isSelectionMode.value || !currentHoverEl.value) return;

    if (e.key === 'ArrowUp') {
      e.preventDefault();
      const parent = currentHoverEl.value.parentElement;
      if (
        parent &&
        parent.tagName !== 'BODY' &&
        parent.tagName !== 'HTML'
      ) {
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
    if (!isSelectionMode.value) return;
    e.preventDefault();
    e.stopPropagation();

    // Use currentHighlighted element instead of event target directly
    // This allows picking parent via keyboard
    const target = currentHoverEl.value || (e.target as HTMLElement);
    if (!target) return;

    if (!currentTargetField.value) {
      Message.info(
        'Select a field (Pick button) first before clicking an element.'
      );
      return;
    }

    // Determine if we are picking item selector
    const isItemSelector = currentTargetField.value === 'item_selector';
    const fullSelector = getCssSelector(target, isItemSelector);

    // Need access to iframe body to query selector
    const doc = previewIframe.value?.contentDocument;
    if (!doc) return;

    if (isItemSelector) {
      config.item_selector = fullSelector;

      // Try to validate how many items matched
      try {
        const matches = doc.querySelectorAll(fullSelector);
        Message.success(
          `Set Item Selector: ${fullSelector} (${matches.length} matches)`
        );
      } catch (err) {
        Message.success(`Set Item Selector: ${fullSelector}`);
      }
    } else {
      // Relative selection
      if (!config.item_selector) {
        Message.warning(
          'Please set List Item Selector first to calculate relative path.'
        );
        return;
      }

      // Check if target is inside an element matching item_selector
      const items = doc.querySelectorAll(config.item_selector);
      let foundItem: HTMLElement | null = null;

      if (items) {
        for (let i = 0; i < items.length; i += 1) {
          if (items[i].contains(target)) {
            foundItem = items[i] as HTMLElement;
            break;
          }
        }
      }

      if (foundItem) {
        if (target === foundItem) {
          config[currentTargetField.value] = ''; // "this"
          Message.info(
            "Selected the item itself. Leave empty to use item's direct text/attr?"
          );
        } else {
          const relPath: string[] = [];
          let curr: HTMLElement = target;
          while (curr && curr !== foundItem) {
            let selector = curr.tagName.toLowerCase();
            if (curr.classList.length > 0) {
              selector += `.${Array.from(curr.classList)[0]}`; // take first class
            }
            relPath.unshift(selector);
            curr = curr.parentNode as HTMLElement;
          }
          config[currentTargetField.value] = relPath.join(' ');
          Message.success(
            `Set ${currentTargetField.value} (relative): ${relPath.join(' ')}`
          );
        }
      } else {
        Message.warning(
          'Clicked element is not inside any element matching the List Item Selector'
        );
      }
    }

    currentTargetField.value = ''; // Reset picker state
  };

  const onIframeLoad = () => {
    const iframe = previewIframe.value;
    if (iframe && iframe.contentDocument) {
      const doc = iframe.contentDocument;

      // Inject styles for hover effect
      const style = doc.createElement('style');
      style.textContent = `
            .fc-highlight { outline: 2px dashed #165dff !important; background-color: rgba(22, 93, 255, 0.05) !important; cursor: pointer; }
          `;
      doc.head.appendChild(style);

      // Attach listeners
      doc.addEventListener('click', handleClick);
      doc.addEventListener('mouseover', handleMouseOver);
      doc.addEventListener('keydown', handleKeyDown);
      // Ensure iframe can receive focus for key events
      doc.body.setAttribute('tabindex', '0');
    }
  };
</script>

<style scoped>
  .html-preview {
    border: 1px solid #e5e6eb;
    padding: 0; /* Remove padding for iframe */
    width: 100%;
    height: 600px;
    background: white;
  }
</style>
