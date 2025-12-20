<template>
  <div class="container">
    <a-space direction="vertical" size="large" fill>
      <a-card title="Step 1: Search Query">
        <template #extra>
          <a-button type="primary" :loading="fetching" @click="handleFetch">
            <template #icon><icon-search /></template>
            Search & Fetch JSON
          </a-button>
        </template>
        <a-form :model="fetchReq" layout="vertical">
          <a-form-item label="Search Query" field="query" required>
            <a-input v-model="fetchReq.query" placeholder="e.g. 'latest AI news'" @press-enter="handleFetch" />
          </a-form-item>
        </a-form>
      </a-card>

      <a-card v-if="jsonContent" title="Step 2: Parsing Rules (JSON to RSS)">
        <template #extra>
          <a-button type="primary" :loading="parsing" @click="handlePreview">
            <template #icon><icon-eye /></template>
            Preview RSS
          </a-button>
        </template>
        <div class="split-view">
          <div class="json-view">
            <h3>Response JSON</h3>
            <a-textarea
              v-model="jsonContent"
              :auto-size="{ minRows: 10, maxRows: 30 }"
              style="font-family: monospace"
            />
          </div>
          <div class="rules-view">
            <a-alert type="info" style="margin-bottom: 10px">
              Configure how to map the JSON search results to RSS fields.
              Use 'jq' syntax (e.g. .data[] or .items[]).
            </a-alert>
            <a-form :model="parseReq" layout="vertical">
              <a-form-item label="List Selector" required>
                <a-input v-model="parseReq.list_selector" placeholder=".data or .items" />
              </a-form-item>
              <a-form-item label="Title Selector" required>
                <a-input v-model="parseReq.title_selector" placeholder=".title" />
              </a-form-item>
              <a-form-item label="Link Selector">
                <a-input v-model="parseReq.link_selector" placeholder=".url" />
              </a-form-item>
              <a-form-item label="Date Selector">
                <a-input v-model="parseReq.date_selector" placeholder=".published_at" />
              </a-form-item>
              <a-form-item label="Description/Content Selector">
                <a-input v-model="parseReq.content_selector" placeholder=".snippet or .content" />
              </a-form-item>
            </a-form>
          </div>
        </div>
      </a-card>

      <a-card v-if="parsedItems.length > 0" title="Step 3: Preview & Save">
        <template #extra>
          <a-button type="primary" status="success" @click="handleShowSaveModal">
            <template #icon><icon-save /></template>
            Save Recipe
          </a-button>
        </template>
        <a-list :data="parsedItems">
          <template #item="{ item }">
            <a-list-item>
              <a-list-item-meta :title="item.title" :description="item.date">
              </a-list-item-meta>
              <div>
                <a :href="item.link" target="_blank">{{ item.link }}</a>
              </div>
              <div v-if="item.content" class="content-preview">
                {{ item.content.substring(0, 150) }}...
              </div>
            </a-list-item>
          </template>
        </a-list>
      </a-card>
    </a-space>

    <!-- Save Recipe Modal -->
    <a-modal v-model:visible="saveModalVisible" title="Save Recipe" @ok="handleSaveRecipe">
      <a-form :model="recipeForm">
        <a-form-item label="Recipe Name (ID)" required>
          <a-input v-model="recipeForm.name" />
        </a-form-item>
        <a-form-item label="Description">
           <a-input v-model="recipeForm.description" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import { fetchSearch, parseJsonRss, ParsedItem, SearchFetchReq } from '@/api/json_rss';
import { createCustomRecipe } from '@/api/custom_recipe';
import { Message } from '@arco-design/web-vue';
import { useRouter } from 'vue-router';

const router = useRouter();

const fetchReq = reactive<SearchFetchReq>({
  query: '',
});

const fetching = ref(false);
const parsing = ref(false);
const jsonContent = ref('');
const parsedItems = ref<ParsedItem[]>([]);

const parseReq = reactive({
  list_selector: '.',
  title_selector: '.title',
  link_selector: '.url',
  date_selector: '',
  content_selector: '.content',
});

const saveModalVisible = ref(false);
const recipeForm = reactive({
  name: '',
  description: '',
});

const handleFetch = async () => {
  if (!fetchReq.query) {
    Message.warning('Query is required');
    return;
  }
  fetching.value = true;
  try {
    const res = await fetchSearch(fetchReq);
    jsonContent.value = res.data;
    Message.success('Search results fetched');
  } catch (err: any) {
    Message.error(`Failed to fetch search results: ${err.message || err}`);
  } finally {
    fetching.value = false;
  }
};

const handlePreview = async () => {
  if (!jsonContent.value) return;
  parsing.value = true;
  try {
    const res = await parseJsonRss({
      json_content: jsonContent.value,
      ...parseReq,
    });
    parsedItems.value = res.data;
    Message.success(`Parsed ${res.data.length} items`);
  } catch (err: any) {
    Message.error(`Failed to parse RSS: ${err.message || err}`);
  } finally {
    parsing.value = false;
  }
};

const handleShowSaveModal = () => {
  if (!recipeForm.name) {
    recipeForm.name = `Search_${Date.now()}`;
  }
  saveModalVisible.value = true;
};

const handleSaveRecipe = async () => {
  if (!recipeForm.name) {
    Message.error('Recipe name (ID) is required');
    return;
  }

  const sourceConfig = {
    type: 'search',
    search_fetcher: {
      query: fetchReq.query,
    },
    json_parser: {
      items_iterator: parseReq.list_selector,
      title: parseReq.title_selector,
      link: parseReq.link_selector,
      date: parseReq.date_selector,
      description: parseReq.content_selector,
    },
  };

  try {
    await createCustomRecipe({
      id: recipeForm.name,
      description: recipeForm.description || `Search feed for: ${fetchReq.query}`,
      craft: 'proxy', // Default craft
      source_type: 'search',
      source_config: JSON.stringify(sourceConfig),
    });
    Message.success('Recipe saved successfully');
    saveModalVisible.value = false;
    router.push({ name: 'CustomRecipe' });
  } catch (err: any) {
    Message.error(`Failed to save recipe: ${err.message || err}`);
  }
};
</script>

<style scoped>
.container {
  padding: 20px;
}
.split-view {
  display: flex;
  gap: 20px;
}
.json-view {
  flex: 1;
}
.rules-view {
  flex: 1;
}
.content-preview {
  color: var(--color-text-3);
  font-size: 12px;
  margin-top: 4px;
}
</style>
