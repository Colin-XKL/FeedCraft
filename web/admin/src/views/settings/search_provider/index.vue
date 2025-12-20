<template>
  <div  class="py-8 px-16">
    <a-card :title="$t('settings.searchProvider.title')">
      <a-alert class="mb-4">
        {{ $t('settings.searchProvider.alert') }}
      </a-alert>

      <a-form :model="form" @submit="handleSave" layout="vertical">
        <a-form-item
          :label="$t('settings.searchProvider.provider')"
          field="provider"
          required
          :tooltip="$t('settings.searchProvider.provider.tooltip')"
        >
          <a-select
            v-model="form.provider"
            :placeholder="$t('settings.searchProvider.placeholder.selectProvider')"
          >
            <a-option value="litellm">LiteLLM Proxy</a-option>
          </a-select>
        </a-form-item>

        <a-form-item
          :label="$t('settings.searchProvider.apiUrl')"
          field="api_url"
          required
          :tooltip="$t('settings.searchProvider.apiUrl.tooltip')"
        >
          <a-input
            v-model="form.api_url"
            :placeholder="$t('settings.searchProvider.placeholder.apiUrl')"
          />
        </a-form-item>

        <a-form-item
          :label="$t('settings.searchProvider.apiKey')"
          field="api_key"
          :tooltip="$t('settings.searchProvider.apiKey.tooltip')"
        >
          <a-input-password
            v-model="form.api_key"
            :placeholder="$t('settings.searchProvider.placeholder.apiKey')"
          />
        </a-form-item>

        <a-form-item
          :label="$t('settings.searchProvider.toolName')"
          field="search_tool_name"
          :tooltip="$t('settings.searchProvider.toolName.tooltip')"
        >
          <a-input
            v-model="form.search_tool_name"
            :placeholder="$t('settings.searchProvider.placeholder.toolName')"
          />
        </a-form-item>

        <a-form-item>
          <a-button type="primary" html-type="submit" :loading="saving">{{
            $t('settings.searchProvider.save')
          }}</a-button>
        </a-form-item>
      </a-form>
    </a-card>
  </div>
</template>

<script setup lang="ts">
  import { ref, reactive, onMounted } from 'vue';
  import { Message } from '@arco-design/web-vue';
  import axios from 'axios';
  import { useI18n } from 'vue-i18n';

  const { t } = useI18n();

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
      Message.error(t('settings.searchProvider.msg.urlRequired'));
      return;
    }
    saving.value = true;
    try {
      await axios.post('/api/admin/settings/search-provider', {
        api_url: form.api_url,
        api_key: form.api_key,
        provider: form.provider,
        search_tool_name: form.search_tool_name,
      });
      Message.success(t('settings.searchProvider.msg.saved'));
    } catch (err) {
      Message.error(t('settings.searchProvider.msg.failed'));
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
