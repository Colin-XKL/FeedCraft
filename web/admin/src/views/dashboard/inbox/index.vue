<template>
  <div class="py-8 px-16">
    <x-header :title="t('inbox.title')" :description="t('inbox.desc')" />

    <a-card class="general-card" :title="t('inbox.title')">
      <template #extra>
        <a-button type="primary" @click="handleShowCreate">
          <template #icon>
            <icon-plus />
          </template>
          {{ t('inbox.btn.create') }}
        </a-button>
      </template>

      <a-table
        :data="inboxes"
        :columns="columns"
        :bordered="true"
        :loading="loading"
      >
        <template #is_public="{ record }">
          <a-tag v-if="record.is_public" color="green">
            {{ t('inbox.status.public') }}
          </a-tag>
          <a-tag v-else color="arcoblue">
            {{ t('inbox.status.private') }}
          </a-tag>
        </template>
        <template #actions="{ record }">
          <a-space wrap>
            <a-button
              type="primary"
              size="small"
              @click="handleShowEdit(record)"
            >
              {{ t('inbox.btn.edit') }}
            </a-button>
            <a-button
              type="outline"
              size="small"
              @click="handleShowGuide(record)"
            >
              {{ t('inbox.btn.guide') }}
            </a-button>
            <a-popconfirm
              :content="t('inbox.deleteConfirm')"
              @ok="handleDelete(record.id)"
            >
              <a-button type="text" status="danger" size="small">
                {{ t('inbox.btn.delete') }}
              </a-button>
            </a-popconfirm>
          </a-space>
        </template>
      </a-table>
    </a-card>

    <!-- Create/Edit Modal -->
    <a-modal
      v-model:visible="showModal"
      :title="isEditing ? t('inbox.btn.edit') : t('inbox.btn.create')"
      @ok="handleSubmit"
      @cancel="resetForm"
    >
      <a-form :model="form" layout="vertical">
        <a-form-item
          :label="t('inbox.id')"
          field="id"
          required
          :rules="[
            { required: true, message: t('inbox.id.placeholder') },
            namingValidator,
          ]"
        >
          <a-input
            v-model="form.id"
            :disabled="isEditing"
            :placeholder="t('inbox.id.placeholder')"
          />
        </a-form-item>
        <a-form-item
          :label="t('inbox.name')"
          field="title"
          required
          :rules="[{ required: true, message: t('inbox.name.placeholder') }]"
        >
          <a-input
            v-model="form.title"
            :placeholder="t('inbox.name.placeholder')"
          />
        </a-form-item>
        <a-form-item :label="t('inbox.description')" field="description">
          <a-textarea
            v-model="form.description"
            :placeholder="t('inbox.description.placeholder')"
          />
        </a-form-item>
        <a-form-item :label="t('inbox.maxItems')" field="max_items" required>
          <a-input-number v-model="form.max_items" :min="1" :max="1000" />
        </a-form-item>
        <a-form-item :label="t('inbox.isPublic')" field="is_public">
          <a-switch v-model="form.is_public">
            <template #checked>
              {{ t('inbox.isPublic.true') }}
            </template>
            <template #unchecked>
              {{ t('inbox.isPublic.false') }}
            </template>
          </a-switch>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- Integration Guide Modal -->
    <a-modal
      v-model:visible="showGuideModal"
      :title="t('inbox.guide.title')"
      :footer="false"
      width="700px"
    >
      <div v-if="selectedInbox" class="guide-container">
        <h3>1. {{ t('inbox.guide.pushUrl') }}</h3>
        <p class="desc-text">{{ t('inbox.guide.pushUrl.desc') }}</p>
        <div class="code-box">
          <span class="monospace-text">{{ pushUrl }}</span>
          <a-button type="text" size="mini" @click="copyText(pushUrl)">
            <template #icon><icon-copy /></template>
          </a-button>
        </div>

        <h3 class="mt-4">2. {{ t('inbox.guide.headers') }}</h3>
        <pre class="monospace-text headers-block">
Authorization: Bearer &lt;YOUR_SYSTEM_AUTH_TOKEN&gt;
Content-Type: application/json</pre
        >

        <h3 class="mt-4">3. {{ t('inbox.guide.body') }}</h3>
        <pre class="monospace-text body-block">
[
  {
    "id": "optional-custom-unique-id",
    "title": "Article Title",
    "url": "https://example.com/article",
    "content": "&lt;p&gt;Article body HTML content&lt;/p&gt;",
    "summary": "Short description...",
    "author": "Author Name",
    "timestamp": 1716470400
  }
]</pre
        >

        <h3 class="mt-4">4. {{ t('inbox.guide.example') }}</h3>
        <div class="code-box bash-example">
          <pre class="monospace-text">
curl -X POST "{{ pushUrl }}" \
  -H "Authorization: Bearer &lt;YOUR_SYSTEM_AUTH_TOKEN&gt;" \
  -H "Content-Type: application/json" \
  -d '[{"title": "Test Push", "content": "&lt;p&gt;Hello World!&lt;/p&gt;"}]'</pre
          >
        </div>

        <h3 class="mt-6">5. {{ t('inbox.guide.query') }}</h3>
        <p v-if="selectedInbox.is_public" class="success-text">
          <icon-check-circle-fill /> {{ t('inbox.guide.query.public') }}
        </p>
        <p v-else class="warning-text">
          <icon-exclamation-circle-fill /> {{ t('inbox.guide.query.private') }}
        </p>

        <div class="recipe-step-box">
          <p
            ><strong>{{ t('inbox.guide.recipe.heading') }}:</strong></p
          >
          <ol>
            <li>
              {{ t('inbox.guide.recipe.step1.pre') }}
              <strong>{{ t('inbox.guide.recipe.step1.link') }}</strong>
              {{ t('inbox.guide.recipe.step1.post') }}
            </li>
            <li>
              {{ t('inbox.guide.recipe.step2') }}
              <strong>inbox</strong>。
            </li>
            <li>
              {{ t('inbox.guide.recipe.step3') }}
              <pre class="monospace-text mt-2 text-xs">
{ "inbox_source": { "inbox_id": "{{ selectedInbox.id }}" } }</pre
              >
            </li>
            <li>
              {{ t('inbox.guide.recipe.step4.pre') }}
              <strong>{{ t('inbox.guide.recipe.step4.link') }}</strong>
              {{ t('inbox.guide.recipe.step4.post') }}
            </li>
          </ol>
        </div>
      </div>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
  import { ref, reactive, onMounted, computed } from 'vue';
  import { useI18n } from 'vue-i18n';
  import { Message } from '@arco-design/web-vue';
  import {
    listInboxes,
    createInbox,
    updateInbox,
    deleteInbox,
    Inbox,
  } from '@/api/inbox';
  import { namingValidator } from '@/utils/validator';
  import XHeader from '@/components/header/x-header.vue';
  import {
    IconPlus,
    IconCopy,
    IconCheckCircleFill,
    IconExclamationCircleFill,
  } from '@arco-design/web-vue/es/icon';

  const { t } = useI18n();
  const loading = ref(false);
  const inboxes = ref<Inbox[]>([]);

  // Modals & States
  const showModal = ref(false);
  const isEditing = ref(false);
  const showGuideModal = ref(false);
  const selectedInbox = ref<Inbox | null>(null);

  const form = reactive({
    id: '',
    title: '',
    description: '',
    max_items: 100,
    is_public: true,
  });

  const columns = [
    {
      title: t('inbox.id'),
      dataIndex: 'id',
      width: 150,
    },
    {
      title: t('inbox.name'),
      dataIndex: 'title',
      width: 200,
    },
    {
      title: t('inbox.description'),
      dataIndex: 'description',
    },
    {
      title: t('inbox.maxItems'),
      dataIndex: 'max_items',
      width: 120,
    },
    {
      title: t('inbox.isPublic'),
      slotName: 'is_public',
      width: 120,
    },
    {
      title: t('inbox.actions'),
      slotName: 'actions',
      width: 240,
    },
  ];

  const fetchInboxes = async () => {
    loading.value = true;
    try {
      const res = await listInboxes();
      if (res.data) {
        inboxes.value = res.data;
      }
    } catch (err: any) {
      Message.error(err.message || 'Failed to fetch inboxes');
    } finally {
      loading.value = false;
    }
  };

  const resetForm = () => {
    form.id = '';
    form.title = '';
    form.description = '';
    form.max_items = 100;
    form.is_public = true;
  };

  const handleShowCreate = () => {
    isEditing.value = false;
    resetForm();
    showModal.value = true;
  };

  const handleShowEdit = (record: Inbox) => {
    isEditing.value = true;
    form.id = record.id;
    form.title = record.title || '';
    form.description = record.description || '';
    form.max_items = record.max_items;
    form.is_public = record.is_public;
    showModal.value = true;
  };

  const handleShowGuide = (record: Inbox) => {
    selectedInbox.value = record;
    showGuideModal.value = true;
  };

  const handleSubmit = async () => {
    if (!form.id.trim() || !form.title.trim()) {
      Message.warning('ID and Title are required');
      return;
    }

    try {
      if (isEditing.value) {
        await updateInbox(form.id, form as Inbox);
        Message.success(t('inbox.updateSuccess'));
      } else {
        await createInbox(form as Inbox);
        Message.success(t('inbox.createSuccess'));
      }
      showModal.value = false;
      resetForm();
      fetchInboxes();
    } catch (err: any) {
      Message.error(err.message || 'Failed to submit inbox');
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await deleteInbox(id);
      Message.success(t('inbox.deleteSuccess'));
      fetchInboxes();
    } catch (err: any) {
      Message.error(err.message || 'Failed to delete inbox');
    }
  };

  const copyText = (text: string) => {
    navigator.clipboard.writeText(text);
    Message.success(t('systemAuthToken.copied'));
  };

  const apiBaseUrl = computed(() => {
    return window.location.origin;
  });

  const pushUrl = computed(() => {
    if (!selectedInbox.value) return '';
    return `${apiBaseUrl.value}/api/inbox/${selectedInbox.value.id}/items`;
  });

  onMounted(() => {
    fetchInboxes();
  });
</script>

<style scoped>
  .monospace-text {
    font-family: monospace;
    font-size: 13px;
  }
  .guide-container h3 {
    margin-bottom: 8px;
    font-size: 15px;
    color: var(--color-text-1);
  }
  .desc-text {
    color: var(--color-text-3);
    font-size: 12px;
    margin-bottom: 6px;
  }
  .code-box {
    display: flex;
    align-items: center;
    justify-content: space-between;
    background-color: var(--color-fill-2);
    padding: 8px 12px;
    border-radius: 4px;
    border: 1px solid var(--color-border-2);
    margin-bottom: 16px;
  }
  .headers-block,
  .body-block {
    background-color: var(--color-fill-2);
    padding: 12px;
    border-radius: 4px;
    border: 1px solid var(--color-border-2);
    margin-bottom: 16px;
    overflow-x: auto;
  }
  .bash-example {
    padding: 12px;
    overflow-x: auto;
  }
  .success-text {
    color: rgb(var(--success-6));
    font-weight: 500;
  }
  .warning-text {
    color: rgb(var(--warning-6));
    font-weight: 500;
  }
  .recipe-step-box {
    background-color: var(--color-fill-1);
    border-left: 4px solid rgb(var(--primary-6));
    padding: 12px 16px;
    margin-top: 12px;
    border-radius: 0 4px 4px 0;
  }
  .recipe-step-box ol {
    margin: 4px 0 0 0;
    padding-left: 20px;
    font-size: 13px;
    line-height: 1.8;
  }
</style>
