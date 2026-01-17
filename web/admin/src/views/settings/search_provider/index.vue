<template>
  <div class="py-8 px-16">
    <a-card :title="$t('settings.searchProvider.title')">
      <a-alert class="mb-4">
        {{ $t('settings.searchProvider.alert') }}
      </a-alert>

      <a-form :model="form" layout="vertical" @submit="handleSave">
        <a-form-item
          :label="$t('settings.searchProvider.provider')"
          field="provider"
          required
          :tooltip="$t('settings.searchProvider.provider.tooltip')"
        >
          <a-select
            v-model="form.provider"
            :placeholder="
              $t('settings.searchProvider.placeholder.selectProvider')
            "
          >
            <a-option value="litellm">LiteLLM Proxy</a-option>
            <a-option value="searxng">SearXNG</a-option>
          </a-select>
        </a-form-item>

        <a-form-item
          :label="$t('settings.searchProvider.apiUrl')"
          field="api_url"
          required
          :tooltip="apiUrlTooltip"
        >
          <a-input v-model="form.api_url" :placeholder="apiUrlPlaceholder" />
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

        <!-- LiteLLM Specific -->
        <a-form-item
          v-if="form.provider === 'litellm'"
          :label="$t('settings.searchProvider.toolName')"
          field="litellm.search_tool_name"
          :tooltip="$t('settings.searchProvider.toolName.tooltip')"
        >
          <a-input
            v-model="form.litellm.search_tool_name"
            :placeholder="$t('settings.searchProvider.placeholder.toolName')"
          />
        </a-form-item>

        <!-- SearXNG Specific -->
        <a-form-item
          v-if="form.provider === 'searxng'"
          :label="$t('settings.searchProvider.engines')"
          field="searxng.engines"
          :tooltip="$t('settings.searchProvider.engines.tooltip')"
        >
          <a-input
            v-model="form.searxng.engines"
            :placeholder="$t('settings.searchProvider.placeholder.engines')"
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
  import { ref, reactive, onMounted, computed } from 'vue';
  import { Message } from '@arco-design/web-vue';
  import axios from 'axios';
  import { useI18n } from 'vue-i18n';

  const { t } = useI18n();

  const form = reactive({
    api_url: '',
    api_key: '',
    provider: 'litellm',
    litellm: {
      search_tool_name: '',
    },
    searxng: {
      engines: '',
    },
  });

  const apiUrlPlaceholder = computed(() => {
    if (form.provider === 'searxng') {
      return 'http://localhost:8080';
    }
    return t('settings.searchProvider.placeholder.apiUrl');
  });

  const apiUrlTooltip = computed(() => {
    if (form.provider === 'searxng') {
      return t('settings.searchProvider.apiUrl.tooltip.searxng');
    }
    return t('settings.searchProvider.apiUrl.tooltip.litellm');
  });

  const saving = ref(false);

  const loadConfig = async () => {
    try {
      const res = await axios.get('/api/admin/settings/search-provider');
      const data = res.data?.data || {};

      form.api_url = data.api_url || '';
      form.api_key = data.api_key || '';
      form.provider = data.provider || 'litellm';

      if (data.litellm) {
        form.litellm.search_tool_name = data.litellm.search_tool_name || '';
      }
      if (data.searxng) {
        form.searxng.engines = data.searxng.engines || '';
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
        litellm: form.litellm,
        searxng: form.searxng,
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
