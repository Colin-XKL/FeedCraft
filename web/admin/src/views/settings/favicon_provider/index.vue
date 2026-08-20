<template>
  <div class="py-8 px-16">
    <a-card :title="$t('settings.faviconProvider.title')">
      <a-alert class="mb-6">
        {{ $t('settings.faviconProvider.alert') }}
      </a-alert>

      <a-spin :loading="loading" class="w-full">
        <a-form :model="form" layout="vertical">
          <a-form-item
            :label="$t('settings.faviconProvider.defaultProvider')"
            required
          >
            <a-select v-model="form.default_provider_id">
              <a-option
                v-for="provider in enabledProviderOptions"
                :key="provider.id"
                :value="provider.id"
              >
                {{ provider.name }} ({{ provider.id }})
              </a-option>
            </a-select>
          </a-form-item>

          <a-divider orientation="left">
            {{ $t('settings.faviconProvider.builtInProviders') }}
          </a-divider>
          <a-table
            :columns="builtInColumns"
            :data="builtInProviders"
            :pagination="false"
            row-key="id"
            class="mb-6"
          >
            <template #enabled="{ record }">
              <a-tag :color="record.enabled ? 'green' : 'gray'">
                {{
                  record.enabled
                    ? $t('settings.faviconProvider.enabled')
                    : $t('settings.faviconProvider.disabled')
                }}
              </a-tag>
            </template>
          </a-table>

          <div class="section-toolbar">
            <h3 class="section-toolbar__title">
              {{ $t('settings.faviconProvider.customProviders') }}
            </h3>
            <a-button type="outline" @click="addCustomProvider">
              <template #icon>
                <icon-plus />
              </template>
              {{ $t('settings.faviconProvider.add') }}
            </a-button>
          </div>

          <a-empty
            v-if="form.custom_providers.length === 0"
            :description="$t('settings.faviconProvider.empty')"
            class="mb-6"
          />

          <a-card
            v-for="(provider, index) in form.custom_providers"
            :key="index"
            size="small"
            class="mb-4"
          >
            <template #title>
              {{ provider.name || $t('settings.faviconProvider.unnamed') }}
            </template>
            <template #extra>
              <a-space>
                <span>{{ $t('settings.faviconProvider.enabled') }}</span>
                <a-switch v-model="provider.enabled" />
                <a-popconfirm
                  :content="$t('settings.faviconProvider.deleteConfirm')"
                  @ok="removeCustomProvider(index)"
                >
                  <a-button status="danger" size="small">
                    {{ $t('settings.faviconProvider.delete') }}
                  </a-button>
                </a-popconfirm>
              </a-space>
            </template>

            <a-row :gutter="16">
              <a-col :span="8">
                <a-form-item
                  :label="$t('settings.faviconProvider.id')"
                  required
                >
                  <a-input
                    v-model="provider.id"
                    placeholder="my_favicon_service"
                  />
                </a-form-item>
              </a-col>
              <a-col :span="16">
                <a-form-item
                  :label="$t('settings.faviconProvider.name')"
                  required
                >
                  <a-input v-model="provider.name" />
                </a-form-item>
              </a-col>
            </a-row>
            <a-form-item
              :label="$t('settings.faviconProvider.urlTemplate')"
              :help="$t('settings.faviconProvider.urlTemplate.help')"
              required
            >
              <a-textarea
                v-model="provider.url_template"
                :auto-size="{ minRows: 2, maxRows: 5 }"
                placeholder="https://icons.example.com/favicon?url={origin_query}&size={size}"
              />
            </a-form-item>
          </a-card>

          <a-divider orientation="left">
            {{ $t('settings.faviconProvider.preview') }}
          </a-divider>
          <a-row :gutter="16">
            <a-col :span="8">
              <a-form-item
                :label="$t('settings.faviconProvider.previewProvider')"
              >
                <a-select v-model="previewProviderId">
                  <a-option
                    v-for="provider in enabledProviderOptions"
                    :key="provider.id"
                    :value="provider.id"
                  >
                    {{ provider.name }}
                  </a-option>
                </a-select>
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item :label="$t('settings.faviconProvider.previewUrl')">
                <a-input v-model="previewPageUrl" />
              </a-form-item>
            </a-col>
            <a-col :span="4">
              <a-form-item label=" ">
                <a-button
                  type="outline"
                  :loading="previewing"
                  @click="handlePreview"
                >
                  {{ $t('settings.faviconProvider.previewAction') }}
                </a-button>
              </a-form-item>
            </a-col>
          </a-row>
          <a-card v-if="previewURL" size="small" class="mb-6">
            <div class="flex items-center gap-4">
              <img
                :src="previewURL"
                :alt="$t('settings.faviconProvider.preview')"
                class="favicon-preview"
              />
              <a-typography-text copyable>{{ previewURL }}</a-typography-text>
            </div>
          </a-card>

          <a-button type="primary" :loading="saving" @click="handleSave">
            {{ $t('settings.faviconProvider.save') }}
          </a-button>
        </a-form>
      </a-spin>
    </a-card>
  </div>
</template>

<script setup lang="ts">
  import { computed, onMounted, reactive, ref, watch } from 'vue';
  import { Message } from '@arco-design/web-vue';
  import { IconPlus } from '@arco-design/web-vue/es/icon';
  import { useI18n } from 'vue-i18n';
  import {
    FaviconProviderConfig,
    FaviconProviderDescriptor,
    FaviconSettings,
    getFaviconProviderConfig,
    previewFaviconProviderConfig,
    saveFaviconProviderConfig,
  } from '@/api/settings';

  const { t } = useI18n();
  const loading = ref(false);
  const saving = ref(false);
  const previewing = ref(false);
  const previewURL = ref('');
  const previewPageUrl = ref('https://github.com/Colin-XKL/FeedCraft');
  const previewProviderId = ref('gstatic_cn');
  const providers = ref<FaviconProviderDescriptor[]>([]);
  const form = reactive<FaviconSettings>({
    default_provider_id: 'gstatic_cn',
    custom_providers: [],
  });

  const builtInProviders = computed(() =>
    providers.value.filter((provider) => provider.built_in)
  );
  const enabledProviderOptions = computed(() => {
    const customOptions: FaviconProviderDescriptor[] = form.custom_providers
      .filter((provider) => provider.enabled)
      .map((provider) => ({ ...provider, built_in: false }));
    return [...builtInProviders.value, ...customOptions];
  });
  const builtInColumns = computed(() => [
    {
      title: t('settings.faviconProvider.name'),
      dataIndex: 'name',
    },
    {
      title: t('settings.faviconProvider.id'),
      dataIndex: 'id',
    },
    {
      title: t('settings.faviconProvider.urlTemplate'),
      dataIndex: 'url_template',
      ellipsis: true,
      tooltip: true,
    },
    {
      title: t('settings.faviconProvider.status'),
      slotName: 'enabled',
      width: 120,
    },
  ]);

  watch(
    enabledProviderOptions,
    (options) => {
      if (!options.some((item) => item.id === form.default_provider_id)) {
        form.default_provider_id = 'gstatic_cn';
      }
      if (!options.some((item) => item.id === previewProviderId.value)) {
        previewProviderId.value =
          form.default_provider_id || options[0]?.id || 'gstatic_cn';
      }
    },
    { deep: true }
  );

  const loadConfig = async () => {
    loading.value = true;
    try {
      const response = await getFaviconProviderConfig();
      const { data } = response.data;
      form.default_provider_id = data.default_provider_id || 'gstatic_cn';
      form.custom_providers = (data.custom_providers || []).map((item) => ({
        ...item,
      }));
      providers.value = data.providers || [];
      previewProviderId.value = form.default_provider_id;
    } catch (error: any) {
      Message.error(
        error.response?.data?.msg ||
          t('settings.faviconProvider.msg.loadFailed')
      );
    } finally {
      loading.value = false;
    }
  };

  const addCustomProvider = () => {
    const usedIDs = new Set([
      ...providers.value.map((provider) => provider.id),
      ...form.custom_providers.map((provider) => provider.id),
    ]);
    let suffix = 1;
    while (usedIDs.has(`custom_${suffix}`)) {
      suffix += 1;
    }
    const provider: FaviconProviderConfig = {
      id: `custom_${suffix}`,
      name: '',
      url_template: '',
      enabled: true,
    };
    form.custom_providers.push(provider);
  };

  const removeCustomProvider = (index: number) => {
    const [removed] = form.custom_providers.splice(index, 1);
    if (removed?.id === form.default_provider_id) {
      form.default_provider_id = 'gstatic_cn';
    }
  };

  const formPayload = (): FaviconSettings => ({
    default_provider_id: form.default_provider_id,
    custom_providers: form.custom_providers.map((provider) => ({
      ...provider,
    })),
  });

  const handleSave = async () => {
    saving.value = true;
    try {
      await saveFaviconProviderConfig(formPayload());
      Message.success(t('settings.faviconProvider.msg.saved'));
      await loadConfig();
    } catch (error: any) {
      Message.error(
        error.response?.data?.msg ||
          t('settings.faviconProvider.msg.saveFailed')
      );
    } finally {
      saving.value = false;
    }
  };

  const handlePreview = async () => {
    previewing.value = true;
    previewURL.value = '';
    try {
      const response = await previewFaviconProviderConfig(
        formPayload(),
        previewProviderId.value,
        previewPageUrl.value
      );
      previewURL.value = response.data.data.url;
    } catch (error: any) {
      Message.error(
        error.response?.data?.msg ||
          t('settings.faviconProvider.msg.previewFailed')
      );
    } finally {
      previewing.value = false;
    }
  };

  onMounted(loadConfig);
</script>

<style scoped>
  .section-toolbar {
    display: flex;
    align-items: center;
    justify-content: flex-start;
    gap: 12px;
    margin: 20px 0 16px;
    padding-bottom: 16px;
    border-bottom: 1px solid var(--color-border-2);
  }

  .section-toolbar__title {
    margin: 0;
    font-size: 16px;
    font-weight: 500;
    line-height: 32px;
    color: var(--color-text-1);
    white-space: nowrap;
  }

  .section-toolbar :deep(.arco-btn) {
    flex: none;
    white-space: nowrap;
  }

  .favicon-preview {
    width: 64px;
    height: 64px;
    object-fit: contain;
  }
</style>
