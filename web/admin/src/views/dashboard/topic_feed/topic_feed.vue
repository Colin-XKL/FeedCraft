<template>
  <div class="container">
    <Breadcrumb :items="['menu.worktable', 'menu.topicFeed']" />
    <a-card class="general-card" :title="$t('menu.topicFeed')">
      <template #extra>
        <a-button type="primary" @click="handleAdd">
          <template #icon>
            <icon-plus />
          </template>
          Create Topic
        </a-button>
      </template>

      <a-table :data="topics" :loading="loading" :pagination="false">
        <template #columns>
          <a-table-column title="ID" data-index="id" />
          <a-table-column title="Title" data-index="title" />
          <a-table-column title="Description" data-index="description" />
          <a-table-column title="Input URIs">
            <template #cell="{ record }">
              <div v-for="(uri, idx) in record.input_uris" :key="idx">
                {{ uri }}
              </div>
            </template>
          </a-table-column>
          <a-table-column title="Actions">
            <template #cell="{ record }">
              <a-space>
                <a-button type="text" size="small" @click="handleEdit(record)">
                  Edit
                </a-button>
                <a-popconfirm
                  content="Are you sure you want to delete this topic?"
                  @ok="handleDelete(record.id)"
                >
                  <a-button type="text" status="danger" size="small">
                    Delete
                  </a-button>
                </a-popconfirm>
              </a-space>
            </template>
          </a-table-column>
        </template>
      </a-table>
    </a-card>

    <a-modal
      v-model:visible="modalVisible"
      :title="isEdit ? 'Edit Topic' : 'Create Topic'"
      @ok="handleSubmit"
      @cancel="modalVisible = false"
    >
      <a-form :model="formData" layout="vertical">
        <a-form-item
          field="id"
          label="ID"
          :rules="[{ required: true, message: 'ID is required' }]"
        >
          <a-input v-model="formData.id" :disabled="isEdit" />
        </a-form-item>
        <a-form-item field="title" label="Title">
          <a-input v-model="formData.title" />
        </a-form-item>
        <a-form-item field="description" label="Description">
          <a-textarea v-model="formData.description" />
        </a-form-item>

        <a-form-item label="Input URIs">
          <div
            v-for="(uri, idx) in formData.input_uris"
            :key="idx"
            class="uri-input"
          >
            <a-input
              v-model="formData.input_uris[idx]"
              placeholder="e.g. feedcraft://recipe/my-recipe"
            />
            <a-button type="text" status="danger" @click="removeUri(idx)">
              <icon-delete />
            </a-button>
          </div>
          <a-button type="dashed" long @click="addUri">
            <icon-plus /> Add URI
          </a-button>
        </a-form-item>

        <a-form-item label="Aggregator Config">
          <div
            v-for="(step, idx) in formData.aggregator_config"
            :key="idx"
            class="aggregator-step"
          >
            <a-space>
              <a-select
                v-model="step.type"
                style="width: 120px"
                placeholder="Type"
              >
                <a-option value="deduplicate">Deduplicate</a-option>
                <a-option value="sort">Sort</a-option>
                <a-option value="limit">Limit</a-option>
              </a-select>
              <a-input
                v-model="step.optionKey"
                placeholder="Option Key"
                style="width: 100px"
              />
              <a-input
                v-model="step.optionValue"
                placeholder="Option Value"
                style="width: 100px"
              />
              <a-button type="text" status="danger" @click="removeStep(idx)">
                <icon-delete />
              </a-button>
            </a-space>
          </div>
          <a-button type="dashed" long @click="addStep" style="margin-top: 8px">
            <icon-plus /> Add Step
          </a-button>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
  import { ref, onMounted } from 'vue';
  import { Message } from '@arco-design/web-vue';
  import {
    listTopicFeeds,
    createTopicFeed,
    updateTopicFeed,
    deleteTopicFeed,
    TopicFeed,
  } from '@/api/topic';

  const topics = ref<TopicFeed[]>([]);
  const loading = ref(false);
  const modalVisible = ref(false);
  const isEdit = ref(false);

  const defaultFormData = {
    id: '',
    title: '',
    description: '',
    input_uris: [],
    aggregator_config: [],
  };

  const formData = ref<any>(JSON.parse(JSON.stringify(defaultFormData)));

  const fetchTopics = async () => {
    loading.value = true;
    try {
      const res = await listTopicFeeds();
      topics.value = res.data;
    } catch (err) {
      Message.error('Failed to fetch topics');
    } finally {
      loading.value = false;
    }
  };

  const handleAdd = () => {
    isEdit.value = false;
    formData.value = JSON.parse(JSON.stringify(defaultFormData));
    modalVisible.value = true;
  };

  const handleEdit = (record: TopicFeed) => {
    isEdit.value = true;
    const configWithKV = (record.aggregator_config || []).map((step) => {
      const key = Object.keys(step.option || {})[0] || '';
      const value = step.option?.[key] || '';
      return {
        type: step.type,
        optionKey: key,
        optionValue: value,
      };
    });

    formData.value = {
      ...record,
      aggregator_config: configWithKV,
    };
    modalVisible.value = true;
  };

  const handleDelete = async (id: string) => {
    try {
      await deleteTopicFeed(id);
      Message.success('Deleted successfully');
      fetchTopics();
    } catch (err) {
      Message.error('Failed to delete');
    }
  };

  const addUri = () => {
    formData.value.input_uris.push('');
  };

  const removeUri = (idx: number) => {
    formData.value.input_uris.splice(idx, 1);
  };

  const addStep = () => {
    formData.value.aggregator_config.push({
      type: 'limit',
      optionKey: 'max',
      optionValue: '50',
    });
  };

  const removeStep = (idx: number) => {
    formData.value.aggregator_config.splice(idx, 1);
  };

  const handleSubmit = async () => {
    try {
      const payload: TopicFeed = {
        id: formData.value.id,
        title: formData.value.title,
        description: formData.value.description,
        input_uris: formData.value.input_uris.filter(
          (u: string) => u.trim() !== '',
        ),
        aggregator_config: formData.value.aggregator_config.map((s: any) => ({
          type: s.type,
          option: s.optionKey ? { [s.optionKey]: s.optionValue } : {},
        })),
      };

      if (isEdit.value) {
        await updateTopicFeed(payload.id, payload);
        Message.success('Updated successfully');
      } else {
        await createTopicFeed(payload);
        Message.success('Created successfully');
      }
      modalVisible.value = false;
      fetchTopics();
    } catch (err) {
      Message.error(isEdit.value ? 'Failed to update' : 'Failed to create');
    }
  };

  onMounted(() => {
    fetchTopics();
  });
</script>

<script lang="ts">
  export default {
    name: 'TopicFeed',
  };
</script>

<style scoped>
  .container {
    padding: 0 20px 20px 20px;
  }
  .uri-input {
    display: flex;
    margin-bottom: 8px;
    gap: 8px;
  }
  .aggregator-step {
    margin-bottom: 8px;
  }
</style>
