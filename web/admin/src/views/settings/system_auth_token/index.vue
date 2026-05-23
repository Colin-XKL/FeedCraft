<template>
  <div class="py-8 px-16">
    <x-header
      :title="t('systemAuthToken.title')"
      :description="t('systemAuthToken.desc')"
    />

    <a-card class="general-card" :title="t('systemAuthToken.title')">
      <template #extra>
        <a-button type="primary" @click="showCreateModal = true">
          <template #icon>
            <icon-plus />
          </template>
          {{ t('systemAuthToken.btn.create') }}
        </a-button>
      </template>

      <a-table
        :data="tokens"
        :columns="columns"
        :bordered="true"
        :loading="loading"
      >
        <template #token="{ record }">
          <a-space>
            <span class="monospace-text">
              {{
                isTokenVisible(record.id)
                  ? record.token
                  : maskToken(record.token)
              }}
            </span>
            <a-button
              type="text"
              size="mini"
              :aria-label="
                isTokenVisible(record.id) ? 'Hide Token' : 'Show Token'
              "
              @click="toggleTokenVisibility(record.id)"
            >
              <template #icon>
                <icon-eye-invisible v-if="isTokenVisible(record.id)" />
                <icon-eye v-else />
              </template>
            </a-button>
            <a-button
              type="text"
              size="mini"
              aria-label="Copy Token"
              @click="copyToken(record.token)"
            >
              <template #icon>
                <icon-copy />
              </template>
            </a-button>
          </a-space>
        </template>
        <template #created_at="{ record }">
          {{ formatTime(record.created_at) }}
        </template>
        <template #actions="{ record }">
          <a-popconfirm
            :content="t('systemAuthToken.deleteConfirm')"
            @ok="handleDelete(record.id)"
          >
            <a-button type="text" status="danger" size="small">
              {{ t('inbox.btn.delete') }}
            </a-button>
          </a-popconfirm>
        </template>
      </a-table>
    </a-card>

    <!-- Create Token Modal -->
    <a-modal
      v-model:visible="showCreateModal"
      :title="t('systemAuthToken.btn.create')"
      @ok="handleCreate"
      @cancel="resetForm"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item :label="t('systemAuthToken.label')" field="label" required>
          <a-input
            v-model="form.label"
            :placeholder="t('systemAuthToken.label.placeholder')"
          />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- Created Success Modal -->
    <a-modal
      v-model:visible="showSuccessModal"
      :title="t('systemAuthToken.createSuccess')"
      :footer="false"
      :mask-closable="false"
      :closable="true"
    >
      <div class="success-container">
        <a-alert
          type="warning"
          :title="t('systemAuthToken.createdAlert.title')"
          class="mb-4"
        >
          {{ t('systemAuthToken.createdAlert.desc') }}
        </a-alert>

        <div class="token-box">
          <span class="monospace-text display-token">{{ generatedToken }}</span>
          <a-button
            type="primary"
            size="small"
            class="ml-4"
            @click="copyToken(generatedToken)"
          >
            <template #icon><icon-copy /></template>
            Copy
          </a-button>
        </div>

        <div class="mt-6 flex justify-end">
          <a-button type="primary" @click="showSuccessModal = false">
            Close
          </a-button>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
  import { ref, reactive, onMounted } from 'vue';
  import { useI18n } from 'vue-i18n';
  import { Message } from '@arco-design/web-vue';
  import {
    listSystemAuthTokens,
    createSystemAuthToken,
    deleteSystemAuthToken,
    SystemAuthToken,
  } from '@/api/inbox';
  import XHeader from '@/components/header/x-header.vue';
  import {
    IconPlus,
    IconEye,
    IconEyeInvisible,
    IconCopy,
  } from '@arco-design/web-vue/es/icon';
  import dayjs from 'dayjs';

  const { t } = useI18n();
  const loading = ref(false);
  const tokens = ref<SystemAuthToken[]>([]);

  // Modals & Forms
  const showCreateModal = ref(false);
  const showSuccessModal = ref(false);
  const generatedToken = ref('');
  const form = reactive({
    label: '',
  });

  const columns = [
    {
      title: t('systemAuthToken.id'),
      dataIndex: 'id',
      width: 80,
    },
    {
      title: t('systemAuthToken.label'),
      dataIndex: 'label',
    },
    {
      title: t('systemAuthToken.value'),
      slotName: 'token',
    },
    {
      title: t('systemAuthToken.createdAt'),
      slotName: 'created_at',
      width: 180,
    },
    {
      title: t('systemAuthToken.actions'),
      slotName: 'actions',
      width: 100,
    },
  ];

  // Token Visibility State
  const visibleTokens = reactive<Record<number, boolean>>({});

  const isTokenVisible = (id: number) => {
    return !!visibleTokens[id];
  };

  const toggleTokenVisibility = (id: number) => {
    visibleTokens[id] = !visibleTokens[id];
  };

  const maskToken = (token: string) => {
    if (!token) return '';
    if (token.length < 8) return '****';
    return `${token.substring(0, 4)}****-****-****-****${token.substring(
      token.length - 4
    )}`;
  };

  const formatTime = (timeStr: string) => {
    return dayjs(timeStr).format('YYYY-MM-DD HH:mm:ss');
  };

  const fetchTokens = async () => {
    loading.value = true;
    try {
      const res = await listSystemAuthTokens();
      if (res.data) {
        tokens.value = res.data;
      }
    } catch (err: any) {
      Message.error(err.message || 'Failed to fetch tokens');
    } finally {
      loading.value = false;
    }
  };

  const resetForm = () => {
    form.label = '';
  };

  const handleCreate = async () => {
    if (!form.label.trim()) {
      Message.warning(t('systemAuthToken.label.placeholder'));
      return;
    }

    try {
      const res = await createSystemAuthToken({ label: form.label });
      if (res.data) {
        generatedToken.value = res.data.token;
        showCreateModal.value = false;
        showSuccessModal.value = true;
        resetForm();
        fetchTokens();
      }
    } catch (err: any) {
      Message.error(err.message || 'Failed to create token');
    }
  };

  const handleDelete = async (id: number) => {
    try {
      await deleteSystemAuthToken(id);
      Message.success(t('systemAuthToken.deleteSuccess'));
      fetchTokens();
    } catch (err: any) {
      Message.error(err.message || 'Failed to delete token');
    }
  };

  const copyToken = (text: string) => {
    navigator.clipboard.writeText(text);
    Message.success(t('systemAuthToken.copied'));
  };

  onMounted(() => {
    fetchTokens();
  });
</script>

<style scoped>
  .monospace-text {
    font-family: monospace;
    font-size: 13px;
  }
  .token-box {
    display: flex;
    align-items: center;
    justify-content: space-between;
    background-color: var(--color-fill-2);
    padding: 12px 16px;
    border-radius: 4px;
    border: 1px solid var(--color-border-2);
  }
  .display-token {
    font-weight: 600;
    color: rgb(var(--primary-6));
    font-size: 15px;
  }
</style>
