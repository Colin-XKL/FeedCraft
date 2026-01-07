<template>
  <div class="py-8 px-16">
    <x-header
      :title="t('menu.channel')"
      :description="t('channel.description')"
    >
    </x-header>

    <a-space direction="horizontal" class="mb-4">
      <a-button type="primary" :loading="isLoading" @click="listChannels">
        {{ t('channel.query') }}
      </a-button>
      <a-button
        type="outline"
        @click="
          () => {
            resetForm();
            quickCreate = true;
            form.source_type = 'rss';
            showModal = true;
          }
        "
      >
        <template #icon><icon-plus /></template>
        {{ t('channel.quickCreateRSS') }}
      </a-button>
      <a-button
        type="outline"
        @click="
          () => {
            resetForm();
            showModal = true;
          }
        "
      >
        {{ t('channel.create') }}
      </a-button>
    </a-space>

    <a-table
      :data="channels"
      :columns="columns"
      :bordered="true"
      :loading="isLoading"
    >
      <template #status="{ record }">
        <a-tooltip
          v-if="record.is_active"
          :content="
            t('channel.status.activeTooltip', {
              time: dayjs(record.last_accessed_at).format(
                'YYYY-MM-DD HH:mm:ss',
              ),
            })
          "
        >
          <a-tag color="green" :default-checked="true">{{
            t('channel.status.active')
          }}</a-tag>
        </a-tooltip>
        <a-tooltip v-else :content="t('channel.status.inactiveTooltip')">
          <a-tag color="gray" :default-checked="true">{{
            t('channel.status.inactive')
          }}</a-tag>
        </a-tooltip>
      </template>
      <template #source_config="{ record }">
        <a-space>
          <span
            style="
              max-width: 300px;
              display: inline-block;
              overflow: hidden;
              text-overflow: ellipsis;
              white-space: nowrap;
            "
            :title="humanReadableConfig(record.source_config)"
          >
            {{ humanReadableConfig(record.source_config) }}
          </span>
          <a-tooltip :content="t('channel.viewConfig')">
            <a-button
              type="text"
              size="mini"
              :aria-label="t('channel.viewConfig')"
              @click="viewConfig(record.source_config)"
            >
              <template #icon>
                <icon-eye />
              </template>
            </a-button>
          </a-tooltip>
        </a-space>
      </template>
      <template #actions="{ record }">
        <a-space direction="horizontal">
          <a-button
            type="outline"
            @click="
              () => {
                isUpdating = true;
                showEditModal(record);
              }
            "
          >
            {{ t('channel.edit') }}
          </a-button>
          <a-popconfirm
            :content="t('channel.deleteConfirm')"
            @ok="deleteChannel(record.id)"
          >
            <a-button status="danger">{{ t('channel.delete') }}</a-button>
          </a-popconfirm>
          <a-link :href="`${baseUrl}/channel/${record?.id}`">{{
            t('channel.link')
          }}</a-link>
        </a-space>
      </template>
    </a-table>

    <!-- Create/Edit Modal -->
    <a-modal
      v-model:visible="showModal"
      :title="
        editing
          ? t('channel.editModalTitle.edit')
          : quickCreate
            ? t('channel.quickCreateRSS')
            : t('channel.editModalTitle.create')
      "
    >
      <a-form
        :model="form"
        :label-col="{ span: 6 }"
        :rules="rules"
        :wrapper-col="{ span: 18 }"
      >
        <a-form-item :label="t('channel.form.name')" field="id">
          <a-input v-model="form.id" :disabled="isUpdating" />
        </a-form-item>
        <a-form-item
          :label="t('channel.form.description')"
          field="description"
        >
          <a-input v-model="form.description" />
        </a-form-item>
        <a-form-item :label="t('channel.form.processor')" field="processor_name">
          <ProcessorSelector
            v-model="processorList"
            mode="multiple"
            :placeholder="t('channel.form.placeholder.processor')"
          />
        </a-form-item>

        <!-- Quick Create Fields -->
        <template v-if="quickCreate">
          <a-form-item
            :label="t('channel.form.feedURL')"
            field="feed_url"
            :rules="[
              {
                required: true,
                message: t('channel.form.rule.rssUrlRequired'),
              },
            ]"
          >
            <a-input
              v-model="rssUrl"
              :placeholder="t('channel.form.placeholder.rssUrl')"
            />
          </a-form-item>
        </template>

        <!-- Advanced Fields -->
        <template v-else>
          <a-form-item
            :label="t('channel.form.sourceType')"
            field="source_type"
          >
            <a-select v-model="form.source_type">
              <a-option value="rss">RSS</a-option>
              <a-option value="html">HTML</a-option>
              <a-option value="json">JSON</a-option>
            </a-select>
          </a-form-item>
          <a-form-item
            :label="t('channel.form.sourceConfig')"
            field="source_config"
          >
            <a-textarea
              v-model="form.source_config"
              :auto-size="{ minRows: 3, maxRows: 10 }"
              :placeholder="t('channel.form.placeholder.sourceConfig')"
            />
          </a-form-item>
        </template>
      </a-form>
      <template #footer>
        <a-button
          @click="
            () => {
              showModal = false;
              isUpdating = false;
            }
          "
          >{{ t('channel.form.cancel') }}
        </a-button>
        <a-button type="primary" :loading="saving" @click="saveChannel">{{
          t('channel.form.save')
        }}</a-button>
      </template>
    </a-modal>

    <!-- View Config Modal -->
    <a-modal
      v-model:visible="showConfigModal"
      :title="t('channel.viewConfigModalTitle')"
      :footer="false"
    >
      <pre
        style="
          background-color: #f5f5f5;
          padding: 10px;
          border-radius: 4px;
          overflow: auto;
          max-height: 400px;
        "
        >{{ currentConfig }}</pre
      >
    </a-modal>
  </div>
</template>

<script setup lang="ts">
  import { ref, onMounted, computed } from 'vue';
  import {
    createChannel,
    Channel,
    deleteChannel as deleteChannelApi,
    getChannels,
    updateChannel,
  } from '@/api/channel';
  import XHeader from '@/components/header/x-header.vue';
  import { namingValidator } from '@/utils/validator';
  import { IconEye, IconPlus } from '@arco-design/web-vue/es/icon';
  import { Message } from '@arco-design/web-vue';
  import dayjs from 'dayjs';
  import { useI18n } from 'vue-i18n';
  import ProcessorSelector from '../blueprint/ProcessorSelector.vue';

  const { t } = useI18n();

  const baseUrl = import.meta.env.VITE_API_BASE_URL ?? '';

  const channels = ref<Channel[]>([]);
  const showModal = ref(false);
  const showConfigModal = ref(false);
  const currentConfig = ref('');
  const quickCreate = ref(false);
  const rssUrl = ref('');

  const form = ref<Channel>({
    id: '',
    description: '',
    processor_name: '',
    source_type: 'rss',
    source_config: '',
  });

  const processorList = computed({
    get: () =>
      form.value.processor_name ? form.value.processor_name.split(',').filter(Boolean) : [],
    set: (val: string[] | string) => {
      if (Array.isArray(val)) {
        form.value.processor_name = val.join(',');
      } else if (typeof val === 'string') {
        form.value.processor_name = val;
      }
    },
  });

  const editing = ref(false);
  const selectedChannel = ref<Channel | null>(null);
  const isLoading = ref(false);
  const isUpdating = ref(false);
  const saving = ref(false);

  const columns = [
    { title: t('channel.form.name'), dataIndex: 'id' },
    { title: t('channel.form.description'), dataIndex: 'description' },
    { title: t('channel.form.processor'), dataIndex: 'processor_name' },
    { title: t('channel.status.active'), slotName: 'status' },
    { title: t('channel.form.sourceType'), dataIndex: 'source_type' },
    {
      title: t('channel.form.sourceConfig'),
      dataIndex: 'source_config',
      slotName: 'source_config',
    },
    { title: t('channel.edit'), slotName: 'actions' },
  ];

  async function listChannels() {
    isLoading.value = true;
    channels.value = (await getChannels()).data;
    isLoading.value = false;
  }

  onMounted(() => {
    listChannels();
  });
  const rules = {
    id: [
      {
        required: true,
        message: t('channel.form.rule.nameRequired'),
        trigger: 'blur',
      },
      namingValidator,
    ],
    processor_name: [
      {
        required: true,
        message: t('channel.form.rule.craftRequired'),
        trigger: 'blur',
      },
    ],
    source_type: [
      {
        required: true,
        message: t('channel.form.rule.sourceTypeRequired'),
        trigger: 'change',
      },
    ],
    source_config: [
      {
        required: true,
        message: t('channel.form.rule.sourceConfigRequired'),
        trigger: 'blur',
      },
    ],
  };

  const humanReadableConfig = (configStr: string) => {
    try {
      const config = JSON.parse(configStr);
      // Try to find the URL in common locations
      if (config.http_fetcher && config.http_fetcher.url) {
        return config.http_fetcher.url;
      }
      if (config.url) {
        return config.url;
      }
      return t('channel.jsonConfigFallback');
    } catch (e) {
      return configStr;
    }
  };

  const viewConfig = (configStr: string) => {
    try {
      const obj = JSON.parse(configStr);
      currentConfig.value = JSON.stringify(obj, null, 2);
    } catch (e) {
      currentConfig.value = configStr;
    }
    showConfigModal.value = true;
  };

  const showEditModal = (channel: Channel) => {
    editing.value = true;
    quickCreate.value = false; // Ensure we are not in quick create mode
    selectedChannel.value = channel;

    // Pretty print JSON for editing
    let prettyConfig = channel.source_config;
    try {
      const obj = JSON.parse(channel.source_config);
      prettyConfig = JSON.stringify(obj, null, 2);
    } catch (e) {
      // ignore error, keep original string
    }

    form.value = {
      id: channel.id,
      description: channel.description,
      processor_name: channel.processor_name,
      source_type: channel.source_type,
      source_config: prettyConfig,
    };
    showModal.value = true;
  };

  const saveChannel = async () => {
    if (quickCreate.value) {
      // Construct JSON for Quick Create
      const config = {
        http_fetcher: {
          url: rssUrl.value,
        },
      };
      form.value.source_config = JSON.stringify(config);
      form.value.source_type = 'rss';
    } else {
      // Validate JSON before saving in Advanced mode
      try {
        JSON.parse(form.value.source_config);
      } catch (e) {
        Message.error(t('channel.form.error.invalidJson'));
        return;
      }
    }

    saving.value = true;
    try {
      if (editing.value) {
        if (selectedChannel.value) {
          await updateChannel(form.value);
          selectedChannel.value.description = form.value.description;
          selectedChannel.value.processor_name = form.value.processor_name;
          selectedChannel.value.source_type = form.value.source_type;
          selectedChannel.value.source_config = form.value.source_config;
        }
      } else {
        await createChannel(form.value as Channel);
        await listChannels();
      }
      showModal.value = false;
      form.value = {
        id: '',
        description: '',
        processor_name: '',
        source_type: 'rss',
        source_config: '',
      };
      editing.value = false;
      isUpdating.value = false;
      selectedChannel.value = null;
      quickCreate.value = false;
      rssUrl.value = '';
    } catch (e) {
      Message.error(t('channel.form.error.saveFailed'));
    } finally {
      saving.value = false;
    }
  };

  const deleteChannel = async (id: string) => {
    await deleteChannelApi(id);
    await listChannels();
  };

  function resetForm() {
    form.value = {
      id: '',
      description: '',
      processor_name: '',
      source_type: 'rss',
      source_config: '',
    };
    quickCreate.value = false;
    rssUrl.value = '';
  }
</script>

<style scoped></style>
