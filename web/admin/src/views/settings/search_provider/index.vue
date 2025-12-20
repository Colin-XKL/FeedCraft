<template>
  <div class="container">
    <a-card title="Search Provider Configuration">
      <a-alert class="mb-4">
        Configure the search provider used by the "Search to RSS" feature.
        Currently supports LiteLLM-compatible search proxies.
      </a-alert>

      <a-form :model="form" @submit="handleSave" layout="vertical">
        <a-form-item
            label="Provider Implementation"
            field="provider"
            required
            tooltip="The internal provider logic to use."
        >
           <a-select v-model="form.provider" placeholder="Select Provider">
             <a-option value="litellm">LiteLLM Proxy</a-option>
             <a-option value="searxng">SearXNG</a-option>
           </a-select>
        </a-form-item>

        <a-form-item
            label="API URL"
            field="api_url"
            required
            tooltip="The base URL of the search provider. For LiteLLM, it's typically 'http://litellm:4000'. For SearXNG, it's the instance URL (e.g. 'https://searx.be')."
        >
          <a-input v-model="form.api_url" placeholder="http://litellm-proxy:4000 or https://searx.be" />
        </a-form-item>

        <a-form-item
            label="API Key"
            field="api_key"
            tooltip="Optional API Key if the provider requires authentication. For SearXNG, this is added as a Bearer token in the Authorization header if configured."
        >
          <a-input-password v-model="form.api_key" placeholder="sk-..." />
        </a-form-item>

        <a-form-item
            label="Tool Name"
            field="search_tool_name"
            tooltip="For LiteLLM: The specific tool name (e.g. 'google-search'). For SearXNG: The specific engines to use (e.g. 'google', 'bing'). Can be comma-separated."
        >
          <a-input v-model="form.search_tool_name" placeholder="google-search" />
        </a-form-item>

        <a-form-item>
           <a-button type="primary" html-type="submit" :loading="saving">Save Configuration</a-button>
        </a-form-item>
      </a-form>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue';
import { Message } from '@arco-design/web-vue';
import axios from 'axios';

const form = reactive({
  api_url: '',
  api_key: '',
  provider: 'litellm',
  search_tool_name: '',
});

const saving = ref(false);

const loadConfig = async () => {
  try {
    const res = await axios.get('/api/admin/settings/search-provider');
    if (res.data) {
        const data = res.data || {};
        form.api_url = data.api_url || '';
        form.api_key = data.api_key || '';
        form.provider = data.provider || 'litellm';
        form.search_tool_name = data.search_tool_name || '';
    }
  } catch (err) {
    // ignore
  }
};

const handleSave = async () => {
  if (!form.api_url) {
      Message.error('API URL is required');
      return;
  }
  saving.value = true;
  try {
    await axios.post('/api/admin/settings/search-provider', {
        api_url: form.api_url,
        api_key: form.api_key,
        provider: form.provider,
        search_tool_name: form.search_tool_name
    });
    Message.success('Configuration saved');
  } catch (err) {
    Message.error('Failed to save configuration');
  } finally {
    saving.value = false;
  }
};

onMounted(() => {
    loadConfig();
});
</script>

<style scoped>
.container {
  padding: 20px;
}
</style>
