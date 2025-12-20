<template>
  <div class="container">
    <a-space direction="vertical" size="large" fill>
      <a-card title="Step 1: Search Query">
        <template #extra>
          <a-button type="primary" :loading="fetching" @click="handlePreview">
            <template #icon><icon-search /></template>
            Preview RSS
          </a-button>
        </template>
        <a-form :model="fetchReq" layout="vertical">
          <a-form-item label="Search Query" field="query" required>
            <a-input v-model="fetchReq.query" placeholder="e.g. 'latest AI news'" @press-enter="handlePreview" />
          </a-form-item>
        </a-form>
      </a-card>

      <a-card v-if="parsedItems.length > 0" title="Step 2: Preview & Save">
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
              <div v-if="item.description" class="content-preview">
                {{ item.description.substring(0, 150) }}...
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
import { previewSearch, ParsedItem, SearchFetchReq } from '@/api/json_rss';
import { createCustomRecipe } from '@/api/custom_recipe';
import { Message } from '@arco-design/web-vue';
import { useRouter } from 'vue-router';

const router = useRouter();

const fetchReq = reactive<SearchFetchReq>({
  query: '',
});

const fetching = ref(false);
const parsedItems = ref<ParsedItem[]>([]);

const saveModalVisible = ref(false);
const recipeForm = reactive({
  name: '',
  description: '',
});

const handlePreview = async () => {
  if (!fetchReq.query) {
    Message.warning('Query is required');
    return;
  }
  fetching.value = true;
  parsedItems.value = [];
  try {
    const res = await previewSearch(fetchReq);
    // @ts-ignore
    parsedItems.value = res.data;
    if (parsedItems.value.length === 0) {
        Message.info('No results found');
    } else {
        Message.success(`Found ${parsedItems.value.length} items`);
    }
  } catch (err) {
    // handled by interceptor
  } finally {
    fetching.value = false;
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
    // No json_parser config needed
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
  } catch (err) {
    console.error(err);
  }
};
</script>

<style scoped>
.container {
  padding: 20px;
}
.content-preview {
  color: var(--color-text-3);
  font-size: 12px;
  margin-top: 4px;
}
</style>
