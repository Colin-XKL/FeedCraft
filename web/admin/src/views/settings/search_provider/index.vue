<template>
  <div class="container">
    <a-card title="Search Provider Configuration">
      <a-form :model="form" layout="vertical" @submit="handleSubmit">
        <a-form-item
          field="api_url"
          label="API URL"
          :rules="[{ required: true, message: 'API URL is required' }]"
        >
          <a-input
            v-model="form.api_url"
            placeholder="e.g., https://api.litellm.ai/search"
          />
        </a-form-item>
        <a-form-item field="api_key" label="API Key">
          <a-input-password
            v-model="form.api_key"
            placeholder="Enter API Key"
          />
        </a-form-item>
        <a-form-item field="provider" label="Provider Type">
          <a-select v-model="form.provider" placeholder="Select Provider">
            <a-option value="litellm">LiteLLM Proxy</a-option>
            <a-option value="other">Other (Generic)</a-option>
          </a-select>
        </a-form-item>
        <a-form-item field="search_tool_name" label="Search Tool Name">
          <a-input
            v-model="form.search_tool_name"
            placeholder="e.g. perplexity-search (Optional for LiteLLM)"
          />
        </a-form-item>
        <a-form-item>
          <a-button type="primary" html-type="submit" :loading="loading"
            >Save Configuration</a-button
          >
        </a-form-item>
      </a-form>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
  import { ref, onMounted } from 'vue';
  import {
    getSearchProviderConfig,
    saveSearchProviderConfig,
    SearchProviderConfig,
  } from '@/api/settings';
  import { Message } from '@arco-design/web-vue';

  const form = ref<SearchProviderConfig>({
    api_url: '',
    api_key: '',
    provider: 'litellm',
    search_tool_name: '',
  });
  const loading = ref(false);

  const loadConfig = async () => {
    try {
      const { data } = await getSearchProviderConfig();
      if (data) {
        form.value = data;
      }
    } catch (err) {
      // Message.error('Failed to load configuration');
    }
  };

  const handleSubmit = async () => {
    loading.value = true;
    try {
      await saveSearchProviderConfig(form.value);
      Message.success('Configuration saved successfully');
    } catch (err) {
      Message.error('Failed to save configuration');
    } finally {
      loading.value = false;
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
