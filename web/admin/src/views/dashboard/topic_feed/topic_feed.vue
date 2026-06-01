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
                <a-button
                  type="text"
                  size="small"
                  @click="previewTopic(record.id)"
                >
                  {{ t('topic.preview') }}
                </a-button>
                <a-button type="text" size="small" @click="handleEdit(record)">
                  {{ t('topic.editAction') }}
                </a-button>
                <a-link :href="buildTopicFeedUrl(record.id)" target="_blank">
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
      <a-form ref="formRef" :model="formData" layout="vertical">
        <a-form-item
          field="id"
          :label="t('topic.id')"
          required
          :rules="getRecipeIdRules(t('topic.idRequired'))"
        >
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

        <a-divider
          orientation="left"
          style="margin-top: 4px; margin-bottom: 4px"
        >
          {{ t('topic.sectionInputs') }}
        </a-divider>
        <a-form-item>
          <template #help>{{ t('topic.inputsHelp') }}</template>
          <div class="input-sources-container">
            <div
              v-for="(source, idx) in formData.inputSources"
              :key="`source-${idx}`"
              class="input-source-card"
            >
              <a-radio-group
                v-model="source.sourceType"
                type="button"
                @change="resetSourceValue(idx)"
              >
                <a-radio value="external">{{
                  t('topic.sourceType.external')
                }}</a-radio>
                <a-radio value="recipe">Recipe</a-radio>
                <a-radio value="topic">Topic</a-radio>
              </a-radio-group>
              <a-input
                v-if="source.sourceType === 'external'"
                v-model="source.externalUrl"
                :placeholder="t('topic.sourceUrl.placeholder')"
                allow-clear
              />
              <a-select
                v-else-if="source.sourceType === 'recipe'"
                v-model="source.resourceId"
                :loading="pickerLoading"
                allow-search
                allow-clear
                :placeholder="t('topic.sourceSelect.placeholder.recipe')"
              >
                <a-option
                  v-for="r in availableRecipes"
                  :key="r.id"
                  :value="r.id"
                  :label="r.description ? `${r.id} — ${r.description}` : r.id"
                >
                  <span class="option-id">{{ r.id }}</span>
                  <span v-if="r.description" class="option-desc">
                    — {{ r.description }}
                  </span>
                </a-option>
              </a-select>
              <a-select
                v-else
                v-model="source.resourceId"
                :loading="pickerLoading"
                allow-search
                allow-clear
                :placeholder="t('topic.sourceSelect.placeholder.topic')"
              >
                <a-option
                  v-for="tp in pickerTopics"
                  :key="tp.id"
                  :value="tp.id"
                  :label="tp.title ? `${tp.id} — ${tp.title}` : tp.id"
                >
                  <span class="option-id">{{ tp.id }}</span>
                  <span v-if="tp.title" class="option-desc">
                    — {{ tp.title }}</span
                  >
                </a-option>
              </a-select>
              <a-button
                type="text"
                status="danger"
                class="remove-btn"
                @click="removeSource(idx)"
              >
                <template #icon>
                  <icon-delete />
                </template>
                {{ t('topic.removeInput') }}
              </a-button>
            </div>
            <a-button type="dashed" long class="add-btn" @click="addSource">
              <template #icon>
                <icon-plus />
              </template>
              {{ t('topic.addInput') }}
            </a-button>
          </div>
        </a-form-item>

        <a-divider
          orientation="left"
          style="margin-top: 4px; margin-bottom: 4px"
        >
          {{ t('topic.sectionAggregator') }}
        </a-divider>
        <a-form-item>
          <template #help>{{ t('topic.aggregatorHelp') }}</template>
          <div class="steps-container">
            <div
              v-for="(step, idx) in formData.aggregator_config"
              :key="`step-${idx}`"
              class="step-card"
            >
              <div class="editor-row">
                <a-select
                  v-model="step.type"
                  style="width: 150px; flex-shrink: 0"
                  @change="resetStepValue(idx)"
                >
                  <a-option value="deduplicate">
                    {{ t('topic.stepType.deduplicate') }}
                  </a-option>
                  <a-option value="sort">{{
                    t('topic.stepType.sort')
                  }}</a-option>
                  <a-option value="limit">{{
                    t('topic.stepType.limit')
                  }}</a-option>
                </a-select>

                <template v-if="step.type === 'deduplicate'">
                  <a-select
                    v-model="step.value"
                    style="width: 220px; flex-shrink: 0"
                    @change="onDeduplicateStrategyChange(idx)"
                  >
                    <a-option value="by_link">
                      {{ t('topic.stepOption.strategy.by_link') }}
                    </a-option>
                    <a-option value="by_id">
                      {{ t('topic.stepOption.strategy.by_id') }}
                    </a-option>
                    <a-option value="by_title">
                      {{ t('topic.stepOption.strategy.by_title') }}
                    </a-option>
                    <a-option value="by_simhash">
                      {{ t('topic.stepOption.strategy.by_simhash') }}
                    </a-option>
                    <a-option value="by_embedding">
                      {{ t('topic.stepOption.strategy.by_embedding') }}
                    </a-option>
                  </a-select>
                  <a-input-number
                    v-if="step.value === 'by_simhash'"
                    v-model="step.threshold"
                    :min="0"
                    :max="1"
                    :step="0.01"
                    :precision="2"
                    :placeholder="t('topic.stepOption.threshold.simhash')"
                    style="flex: 1; min-width: 120px"
                  />
                  <a-input-number
                    v-else-if="step.value === 'by_embedding'"
                    v-model="step.threshold"
                    :min="0"
                    :max="1"
                    :step="0.01"
                    :precision="2"
                    :placeholder="t('topic.stepOption.threshold.embedding')"
                    style="flex: 1; min-width: 120px"
                  />
                </template>

                <a-select
                  v-else-if="step.type === 'sort'"
                  v-model="step.value"
                  style="flex: 1; min-width: 150px"
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
                  style="flex: 1; min-width: 150px"
                />

                <a-button
                  type="text"
                  status="danger"
                  class="remove-btn"
                  @click="removeStep(idx)"
                >
                  <template #icon>
                    <icon-delete />
                  </template>
                  {{ t('topic.removeStep') }}
                </a-button>
              </div>
              <p
                v-if="
                  step.type === 'deduplicate' &&
                  (step.value === 'by_simhash' || step.value === 'by_embedding')
                "
                class="step-hint"
              >
                {{ t(`topic.stepOption.strategy.${step.value}.hint`) }}
              </p>
            </div>

            <a-button type="dashed" long class="add-btn" @click="addStep">
              <template #icon>
                <icon-plus />
              </template>
              {{ t('topic.addStep') }}
            </a-button>
          </div>
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
        <a-alert
          v-if="validationWarnings.length > 0"
          type="warning"
          style="margin-top: 12px"
        >
          <template #title>{{ t('topic.validationWarnings') }}</template>
          <div
            v-for="issue in validationWarnings"
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
  import { computed, onMounted, ref } from 'vue';
  import { Message } from '@arco-design/web-vue';
  import { useI18n } from 'vue-i18n';
  import { useRouter } from 'vue-router';
  import { CustomRecipe, getCustomRecipes } from '@/api/custom_recipe';
  import XHeader from '@/components/header/x-header.vue';
  import buildPublicFeedUrl from '@/utils/publicFeedUrl';
  import { getRecipeIdRules } from '@/utils/slug';
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

  type SourceType = 'external' | 'recipe' | 'topic';

  const STRATEGIES_WITH_THRESHOLD = ['by_simhash', 'by_embedding'] as const;

  interface InputSourceItem {
    sourceType: SourceType;
    externalUrl: string;
    resourceId: string;
  }

  interface StepFormItem {
    type: StepType;
    value: string | number;
    /** Deduplication threshold: Hamming distance for by_simhash, cosine similarity for by_embedding */
    threshold?: number;
  }

  interface TopicFormData {
    id: string;
    title: string;
    description: string;
    inputSources: InputSourceItem[];
    aggregator_config: StepFormItem[];
  }

  const { t } = useI18n();
  const router = useRouter();
  const topics = ref<TopicFeed[]>([]);
  const loading = ref(false);
  const formRef = ref();
  const modalVisible = ref(false);
  const isEdit = ref(false);
  const submitting = ref(false);
  const validating = ref(false);
  const validationErrors = ref<TopicValidationIssue[]>([]);
  const validationWarnings = ref<TopicValidationIssue[]>([]);
  const availableRecipes = ref<CustomRecipe[]>([]);
  const availableTopics = ref<TopicFeed[]>([]);
  const pickerLoading = ref(false);

  const defaultThreshold = (strategy: string): number | undefined => {
    if (strategy === 'by_simhash') return 0.05;
    if (strategy === 'by_embedding') return 0.1;
    return undefined;
  };

  const createDefaultStep = (type: StepType = 'limit'): StepFormItem => {
    if (type === 'deduplicate') return { type, value: 'by_link' };
    if (type === 'sort') return { type, value: 'date_desc' };
    return { type, value: 50 };
  };

  const defaultFormData = (): TopicFormData => ({
    id: '',
    title: '',
    description: '',
    inputSources: [{ sourceType: 'external', externalUrl: '', resourceId: '' }],
    aggregator_config: [],
  });

  const parseUriToSource = (uri: string): InputSourceItem => {
    if (uri.startsWith('feedcraft://recipe/')) {
      return {
        sourceType: 'recipe',
        externalUrl: '',
        resourceId: uri.slice('feedcraft://recipe/'.length),
      };
    }
    if (uri.startsWith('feedcraft://topic/')) {
      return {
        sourceType: 'topic',
        externalUrl: '',
        resourceId: uri.slice('feedcraft://topic/'.length),
      };
    }
    return { sourceType: 'external', externalUrl: uri, resourceId: '' };
  };

  const sourceToUri = (source: InputSourceItem): string => {
    if (source.sourceType === 'recipe') {
      return `feedcraft://recipe/${source.resourceId}`;
    }
    if (source.sourceType === 'topic') {
      return `feedcraft://topic/${source.resourceId}`;
    }
    return source.externalUrl.trim();
  };

  const loadPickerData = async () => {
    pickerLoading.value = true;
    try {
      const [recipesRes, topicsRes] = await Promise.all([
        getCustomRecipes(),
        listTopicFeeds(),
      ]);
      availableRecipes.value = recipesRes.data ?? [];
      availableTopics.value = topicsRes.data ?? [];
    } catch {
      // non-fatal: picker continues with empty lists
    } finally {
      pickerLoading.value = false;
    }
  };

  const formData = ref<TopicFormData>(defaultFormData());

  const pickerTopics = computed(() => {
    if (!isEdit.value) return availableTopics.value;
    return availableTopics.value.filter((tp) => tp.id !== formData.value.id);
  });

  const normalizeTopicPayload = (): TopicFeed => ({
    id: formData.value.id.trim(),
    title: formData.value.title.trim(),
    description: formData.value.description.trim(),
    input_uris: formData.value.inputSources
      .map(sourceToUri)
      .filter((uri) => uri !== ''),
    aggregator_config: formData.value.aggregator_config.map((step) => {
      const option: Record<string, string> = {};
      if (step.type === 'deduplicate') {
        option.strategy = String(step.value);
        if (
          step.threshold !== undefined &&
          STRATEGIES_WITH_THRESHOLD.includes(
            step.value as (typeof STRATEGIES_WITH_THRESHOLD)[number]
          )
        ) {
          option.threshold = String(step.threshold);
        }
      }
      if (step.type === 'sort') option.by = String(step.value);
      if (step.type === 'limit') option.max = String(step.value);
      return {
        type: step.type,
        option,
      };
    }),
  });

  const formatDeduplicateStrategy = (step: AggregatorStep): string => {
    const strategy = step.option?.strategy || 'by_link';
    const label = t(`topic.stepOption.strategy.${strategy}`);
    if (strategy === 'by_simhash' && step.option?.threshold) {
      return `${label} (${step.option.threshold})`;
    }
    if (strategy === 'by_embedding' && step.option?.threshold) {
      return `${label} (${step.option.threshold})`;
    }
    return label;
  };

  const formatAggregatorSummary = (steps: AggregatorStep[]) => {
    if (!steps || steps.length === 0) return t('topic.noAggregator');
    return steps
      .map((step) => {
        if (step.type === 'deduplicate') {
          return `${t(
            'topic.stepType.deduplicate'
          )} · ${formatDeduplicateStrategy(step)}`;
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
    formRef.value?.resetFields();
    validationErrors.value = [];
    validationWarnings.value = [];
    loadPickerData();
    modalVisible.value = true;
  };

  const buildTopicFeedUrl = (id: string) => buildPublicFeedUrl(`/topic/${id}`);

  const previewTopic = (id: string) => {
    router.push({
      name: 'FeedViewer',
      query: { target: 'topic', id },
    });
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
      inputSources:
        record.input_uris.length > 0
          ? record.input_uris.map(parseUriToSource)
          : [{ sourceType: 'external', externalUrl: '', resourceId: '' }],
      aggregator_config: (record.aggregator_config || []).map((step) => {
        if (step.type === 'deduplicate') {
          const strategy = step.option?.strategy || 'by_link';
          const item: StepFormItem = { type: 'deduplicate', value: strategy };
          if (
            step.option?.threshold !== undefined &&
            STRATEGIES_WITH_THRESHOLD.includes(
              strategy as (typeof STRATEGIES_WITH_THRESHOLD)[number]
            )
          ) {
            item.threshold = Number(step.option.threshold);
          } else {
            item.threshold = defaultThreshold(strategy);
          }
          return item;
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

  const addSource = () => {
    formData.value.inputSources.push({
      sourceType: 'external',
      externalUrl: '',
      resourceId: '',
    });
  };

  const removeSource = (idx: number) => {
    formData.value.inputSources.splice(idx, 1);
    if (formData.value.inputSources.length === 0) {
      formData.value.inputSources.push({
        sourceType: 'external',
        externalUrl: '',
        resourceId: '',
      });
    }
  };

  const resetSourceValue = (idx: number) => {
    const { sourceType } = formData.value.inputSources[idx];
    formData.value.inputSources[idx] = {
      sourceType,
      externalUrl: '',
      resourceId: '',
    };
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

  const onDeduplicateStrategyChange = (idx: number) => {
    const strategy = String(formData.value.aggregator_config[idx].value);
    formData.value.aggregator_config[idx].threshold =
      defaultThreshold(strategy);
  };

  const runValidation = async () => {
    const payload = normalizeTopicPayload();
    const res = await validateTopicFeed(payload);
    validationErrors.value = res.data?.errors || [];
    validationWarnings.value = res.data?.warnings || [];
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
    const res = await formRef.value?.validate();
    if (res) {
      return;
    }

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
  .input-sources-container,
  .steps-container {
    display: flex;
    flex-direction: column;
    gap: 12px;
    width: 100%;
  }

  .input-source-card,
  .step-card {
    background-color: var(--color-fill-1);
    border: 1px solid var(--color-border-1);
    border-radius: 6px;
    padding: 12px 16px;
    transition: all 0.2s ease;
  }

  .input-source-card:hover,
  .step-card:hover {
    border-color: var(--color-border-3);
    background-color: var(--color-fill-2);
  }

  .input-source-card {
    display: grid;
    grid-template-columns: auto 1fr auto;
    gap: 12px;
    align-items: center;
  }

  .input-source-card :deep(.arco-radio-group) {
    white-space: nowrap;
    flex-shrink: 0;
  }

  .step-card {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .editor-row {
    display: flex;
    gap: 12px;
    align-items: center;
    width: 100%;
  }

  .step-hint {
    margin: 4px 0 0 0;
    padding-left: 4px;
    font-size: 12px;
    color: var(--color-text-3);
    line-height: 1.5;
  }

  .remove-btn {
    flex-shrink: 0;
  }

  .add-btn {
    margin-top: 4px;
  }

  .option-id {
    font-weight: 500;
  }

  .option-desc {
    color: var(--color-text-3);
    margin-left: 2px;
    font-size: 12px;
  }

  .validation-item {
    margin-top: 6px;
  }
</style>
