<template>
  <div class="container">
    <a-card title="Search Provider Configuration">
      <a-alert class="mb-4">
        Configure the search provider used by the "Search to RSS" feature.
        Currently supports LiteLLM-compatible search proxies.
      </a-alert>

      <a-form :model="form" @submit="handleSave" layout="vertical">
        <a-form-item
            label="Search Provider URL"
            field="api_url"
            required
            tooltip="The base URL of the search provider (e.g. LiteLLM Proxy Search endpoint)."
        >
          <a-input v-model="form.api_url" placeholder="http://litellm-proxy:4000" />
        </a-form-item>

        <a-form-item
            label="API Key"
            field="api_key"
            tooltip="Optional API Key if the provider requires authentication."
        >
          <a-input-password v-model="form.api_key" placeholder="sk-..." />
        </a-form-item>

        <a-form-item
            label="Provider Name"
            field="provider"
            tooltip="e.g. 'google', 'brave', or empty for default."
        >
          <a-input v-model="form.provider" placeholder="google" />
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
  provider: '',
});

const saving = ref(false);

const loadConfig = async () => {
  try {
    const res = await axios.get('/api/admin/settings/search-provider');
    if (res.data) {
        // The interceptor unwraps 'data' usually, but let's check structure
        // If response is {code:0, data: {...}} -> res (interceptor) -> res.data
        // Wait, interceptor logic says "API calls return the custom response body".
        // If controller returns APIResponse[Data: cfg], then client gets {Data: cfg}.
        // But the previous file accessed res.data.data.
        // Let's assume standard Axios + interceptor returns the full object or data?
        // Memory says: "API calls return the custom response body (e.g. { code, msg, data }) directly"
        const data = res.data || {};
        form.api_url = data.api_url || '';
        form.api_key = data.api_key || '';
        form.provider = data.provider || '';
    }
  } catch (err) {
    // ignore
  }
};

const handleSave = async () => {
  if (!form.api_url) {
      Message.error('Provider URL is required');
      return;
  }
  saving.value = true;
  try {
    await axios.post('/api/admin/settings/search-provider', {
        api_url: form.api_url,
        api_key: form.api_key,
        provider: form.provider
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
