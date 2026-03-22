<template>
  <div class="py-8 px-16">
    <Breadcrumb :items="['menu.worktable', 'menu.topicFeed']" />
    <x-header
      :title="$t('menu.topicFeed')"
      :description="t('topic.description')"
    />

    <a-card class="general-card" :title="$t('menu.topicFeed')">
      <template #extra>
        <a-space>
          <a-button :loading="loading" @click="fetchTopics">
            {{ t('topic.refresh') }}
          </a-button>
          <a-button type="primary" @click="handleAdd">
            <template #icon>
              <icon-plus />
            </template>
            {{ t('topic.create') }}
          </a-button>
        </a-space>
      </template>

      <a-table
        :data="topics"
        :loading="loading"
        :pagination="false"
        row-key="id"
      >
        <template #columns>
          <a-table-column :title="t('topic.id')" data-index="id" />
          <a-table-column :title="t('topic.title')" data-index="title">
            <template #cell="{ record }">
              {{ record.title || record.id }}
            </template>
          </a-table-column>
          <a-table-column
            :title="t('topic.descriptionLabel')"
            data-index="description"
            :ellipsis="true"
          />
          <a-table-column :title="t('topic.inputCount')">
            <template #cell="{ record }">
              <a-tag color="arcoblue">{{ record.input_uris.length }}</a-tag>
            </template>
          </a-table-column>
          <a-table-column :title="t('topic.aggregator')">
            <template #cell="{ record }">
              <span>{{
                formatAggregatorSummary(record.aggregator_config)
              }}</span>
            </template>
          </a-table-column>
          <a-table-column :title="t('observability.actions')">
            <template #cell="{ record }">
              <a-space wrap>
                <a-button
                  type="text"
                  size="small"
                  @click="goToDetail(record.id)"
                >
                  {{ t('topic.viewDetails') }}
                </a-button>
                <a-button type="text" size="small" @click="handleEdit(record)">
                  {{ t('topic.editAction') }}
                </a-button>
                <a-link :href="`/topic/${record.id}`" target="_blank">
                  {{ t('topic.viewFeed') }}
                </a-link>
                <a-popconfirm
                  :content="t('topic.deleteConfirm')"
                  @ok="handleDelete(record.id)"
                >
                  <a-button type="text" status="danger" size="small">
                    {{ t('topic.deleteAction') }}
                  </a-button>
                </a-popconfirm>
              </a-space>
            </template>
          </a-table-column>
        </template>
      </a-table>

      <a-empty
        v-if="!loading && topics.length === 0"
        :description="t('topic.noTopics')"
      />
    </a-card>

    <a-modal
      v-model:visible="modalVisible"
      :title="isEdit ? t('topic.edit') : t('topic.create')"
      width="860px"
      :mask-closable="false"
      :ok-button-props="{ disabled: submitting }"
      :cancel-button-props="{ disabled: submitting || validating }"
      @cancel="modalVisible = false"
    >
      <a-form :model="formData" layout="vertical">
        <a-form-item field="id" :label="t('topic.id')">
          <a-input
            v-model="formData.id"
            :disabled="isEdit"
            :placeholder="t('topic.id')"
          />
        </a-form-item>
        <a-form-item field="title" :label="t('topic.title')">
          <a-input v-model="formData.title" :placeholder="t('topic.title')" />
        </a-form-item>
        <a-form-item field="description" :label="t('topic.descriptionLabel')">
          <a-textarea
            v-model="formData.description"
            :placeholder="t('topic.descriptionLabel')"
          />
        </a-form-item>

        <a-form-item :label="t('topic.inputs')">
          <template #help>{{ t('topic.inputsHelp') }}</template>
          <div
            v-for="(uri, idx) in formData.input_uris"
            :key="`uri-${idx}`"
            class="editor-row"
          >
            <a-input
              v-model="formData.input_uris[idx]"
              :placeholder="t('topic.inputPlaceholder')"
            />
            <a-button type="text" status="danger" @click="removeUri(idx)">
              {{ t('topic.removeInput') }}
            </a-button>
          </div>
          <a-button type="dashed" long @click="addUri">
            <icon-plus />
            {{ t('topic.addInput') }}
          </a-button>
        </a-form-item>

        <a-form-item :label="t('topic.aggregatorConfig')">
          <template #help>{{ t('topic.aggregatorHelp') }}</template>
          <div
            v-for="(step, idx) in formData.aggregator_config"
            :key="`step-${idx}`"
            class="editor-row"
          >
            <a-select
              v-model="step.type"
              style="width: 180px"
              @change="resetStepValue(idx)"
            >
              <a-option value="deduplicate">
                {{ t('topic.stepType.deduplicate') }}
              </a-option>
              <a-option value="sort">{{ t('topic.stepType.sort') }}</a-option>
              <a-option value="limit">{{ t('topic.stepType.limit') }}</a-option>
            </a-select>

            <a-select
              v-if="step.type === 'deduplicate'"
              v-model="step.value"
              style="width: 220px"
            >
              <a-option value="by_link">
                {{ t('topic.stepOption.strategy.by_link') }}
              </a-option>
              <a-option value="by_id">
                {{ t('topic.stepOption.strategy.by_id') }}
              </a-option>
            </a-select>

            <a-select
              v-else-if="step.type === 'sort'"
              v-model="step.value"
              style="width: 220px"
            >
              <a-option value="date_desc">
                {{ t('topic.stepOption.sort.date_desc') }}
              </a-option>
              <a-option value="date_asc">
                {{ t('topic.stepOption.sort.date_asc') }}
              </a-option>
              <a-option value="quality_desc">
                {{ t('topic.stepOption.sort.quality_desc') }}
              </a-option>
              <a-option value="quality_asc">
                {{ t('topic.stepOption.sort.quality_asc') }}
              </a-option>
            </a-select>

            <a-input-number
              v-else
              v-model="step.value"
              :min="1"
              mode="button"
              style="width: 220px"
            />

            <a-button type="text" status="danger" @click="removeStep(idx)">
              {{ t('topic.removeInput') }}
            </a-button>
          </div>

          <a-button type="dashed" long @click="addStep">
            <icon-plus />
            {{ t('topic.addStep') }}
          </a-button>
        </a-form-item>

        <a-alert v-if="validationErrors.length > 0" type="error">
          <template #title>{{ t('topic.validationSummary') }}</template>
          <div
            v-for="issue in validationErrors"
            :key="`${issue.field}-${issue.message}`"
            class="validation-item"
          >
            <strong>{{ issue.field }}</strong
            >: {{ issue.message }}
          </div>
        </a-alert>
      </a-form>

      <template #footer>
        <a-space>
          <a-button @click="modalVisible = false">{{
            t('topic.cancel')
          }}</a-button>
          <a-button :loading="validating" @click="handleValidate">
            {{ t('topic.validate') }}
          </a-button>
          <a-button type="primary" :loading="submitting" @click="handleSubmit">
            {{ t('topic.save') }}
          </a-button>
        </a-space>
      </template>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
  import { onMounted, ref } from 'vue';
  import { Message } from '@arco-design/web-vue';
  import { useI18n } from 'vue-i18n';
  import { useRouter } from 'vue-router';
  import XHeader from '@/components/header/x-header.vue';
  import {
    AggregatorStep,
    TopicFeed,
    TopicValidationIssue,
    createTopicFeed,
    deleteTopicFeed,
    listTopicFeeds,
    updateTopicFeed,
    validateTopicFeed,
  } from '@/api/topic';

  type StepType = 'deduplicate' | 'sort' | 'limit';

  interface StepFormItem {
    type: StepType;
    value: string | number;
  }

  interface TopicFormData {
    id: string;
    title: string;
    description: string;
    input_uris: string[];
    aggregator_config: StepFormItem[];
  }

  const { t } = useI18n();
  const router = useRouter();
  const topics = ref<TopicFeed[]>([]);
  const loading = ref(false);
  const modalVisible = ref(false);
  const isEdit = ref(false);
  const submitting = ref(false);
  const validating = ref(false);
  const validationErrors = ref<TopicValidationIssue[]>([]);

  const createDefaultStep = (type: StepType = 'limit'): StepFormItem => {
    if (type === 'deduplicate') return { type, value: 'by_link' };
    if (type === 'sort') return { type, value: 'date_desc' };
    return { type, value: 50 };
  };

  const defaultFormData = (): TopicFormData => ({
    id: '',
    title: '',
    description: '',
    input_uris: [''],
    aggregator_config: [],
  });

  const formData = ref<TopicFormData>(defaultFormData());

  const normalizeTopicPayload = (): TopicFeed => ({
    id: formData.value.id.trim(),
    title: formData.value.title.trim(),
    description: formData.value.description.trim(),
    input_uris: formData.value.input_uris
      .map((item) => item.trim())
      .filter((item) => item !== ''),
    aggregator_config: formData.value.aggregator_config.map((step) => {
      const option: Record<string, string> = {};
      if (step.type === 'deduplicate') option.strategy = String(step.value);
      if (step.type === 'sort') option.by = String(step.value);
      if (step.type === 'limit') option.max = String(step.value);
      return {
        type: step.type,
        option,
      };
    }),
  });

  const formatAggregatorSummary = (steps: AggregatorStep[]) => {
    if (!steps || steps.length === 0) return t('topic.noAggregator');
    return steps
      .map((step) => {
        if (step.type === 'deduplicate') {
          return `${t('topic.stepType.deduplicate')} · ${t(
            `topic.stepOption.strategy.${step.option?.strategy || 'by_link'}`
          )}`;
        }
        if (step.type === 'sort') {
          return `${t('topic.stepType.sort')} · ${t(
            `topic.stepOption.sort.${step.option?.by || 'date_desc'}`
          )}`;
        }
        if (step.type === 'limit') {
          return `${t('topic.stepType.limit')} · ${step.option?.max || '-'}`;
        }
        return step.type;
      })
      .join(' / ');
  };

  const fetchTopics = async () => {
    loading.value = true;
    try {
      const res = await listTopicFeeds();
      topics.value = res.data ?? [];
    } catch (err: any) {
      Message.error(err.message || t('topic.fetchFailed'));
    } finally {
      loading.value = false;
    }
  };

  const openModal = () => {
    validationErrors.value = [];
    modalVisible.value = true;
  };

  const handleAdd = () => {
    isEdit.value = false;
    formData.value = defaultFormData();
    openModal();
  };

  const handleEdit = (record: TopicFeed) => {
    isEdit.value = true;
    formData.value = {
      id: record.id,
      title: record.title || '',
      description: record.description || '',
      input_uris: record.input_uris.length > 0 ? [...record.input_uris] : [''],
      aggregator_config: (record.aggregator_config || []).map((step) => {
        if (step.type === 'deduplicate') {
          return {
            type: 'deduplicate',
            value: step.option?.strategy || 'by_link',
          };
        }
        if (step.type === 'sort') {
          return { type: 'sort', value: step.option?.by || 'date_desc' };
        }
        return { type: 'limit', value: Number(step.option?.max || 50) };
      }),
    };
    openModal();
  };

  const handleDelete = async (id: string) => {
    try {
      await deleteTopicFeed(id);
      Message.success(t('topic.deleteSuccess'));
      await fetchTopics();
    } catch (err: any) {
      Message.error(err.message || t('topic.deleteFailed'));
    }
  };

  const addUri = () => {
    formData.value.input_uris.push('');
  };

  const removeUri = (idx: number) => {
    formData.value.input_uris.splice(idx, 1);
    if (formData.value.input_uris.length === 0) {
      formData.value.input_uris.push('');
    }
  };

  const addStep = () => {
    formData.value.aggregator_config.push(createDefaultStep());
  };

  const removeStep = (idx: number) => {
    formData.value.aggregator_config.splice(idx, 1);
  };

  const resetStepValue = (idx: number) => {
    const currentType = formData.value.aggregator_config[idx].type;
    formData.value.aggregator_config[idx] = createDefaultStep(currentType);
  };

  const runValidation = async () => {
    const payload = normalizeTopicPayload();
    const res = await validateTopicFeed(payload);
    validationErrors.value = res.data?.errors || [];
    return res.data;
  };

  const handleValidate = async () => {
    validating.value = true;
    try {
      const result = await runValidation();
      if (result?.valid) {
        Message.success(t('topic.validateSuccess'));
      } else {
        Message.error(t('topic.validateFailed'));
      }
    } catch (err: any) {
      Message.error(err.message || t('topic.validateFailed'));
    } finally {
      validating.value = false;
    }
  };

  const handleSubmit = async () => {
    submitting.value = true;
    try {
      const result = await runValidation();
      if (!result?.valid) {
        Message.error(t('topic.validateFailed'));
        return;
      }

      const payload = normalizeTopicPayload();
      if (isEdit.value) {
        await updateTopicFeed(payload.id, payload);
        Message.success(t('topic.updateSuccess'));
      } else {
        await createTopicFeed(payload);
        Message.success(t('topic.createSuccess'));
      }
      modalVisible.value = false;
      await fetchTopics();
    } catch (err: any) {
      Message.error(err.message || t('topic.saveFailed'));
    } finally {
      submitting.value = false;
    }
  };

  const goToDetail = (id: string) => {
    router.push({ name: 'TopicFeedDetail', params: { id } });
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
  .editor-row {
    display: flex;
    gap: 12px;
    align-items: center;
    margin-bottom: 12px;
  }

  .validation-item {
    margin-top: 6px;
  }
</style>
