<template>
  <div class="py-8 px-16">
    <Breadcrumb
      :items="['menu.worktable', 'menu.topicFeed', breadcrumbActionLabel]"
    />
    <x-header :title="pageTitle" :description="t('topic.wizard.description')" />

    <a-card class="wizard-card">
      <a-steps :current="currentStep" class="wizard-steps">
        <a-step
          :title="t('topic.wizard.step.basic')"
          :description="t('topic.wizard.step.basic.desc')"
        />
        <a-step
          :title="t('topic.wizard.step.inputs')"
          :description="t('topic.wizard.step.inputs.desc')"
        />
        <a-step
          :title="t('topic.wizard.step.aggregator')"
          :description="t('topic.wizard.step.aggregator.desc')"
        />
        <a-step
          :title="t('topic.wizard.step.review')"
          :description="t('topic.wizard.step.review.desc')"
        />
      </a-steps>

      <a-spin :loading="pageLoading" class="wizard-spin">
        <a-form
          ref="formRef"
          :model="formData"
          layout="vertical"
          label-align="left"
          class="wizard-form"
        >
          <section
            v-show="currentStep === 1"
            class="step-panel step-panel--basic"
          >
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
              <a-input
                v-model="formData.title"
                :placeholder="t('topic.title')"
              />
            </a-form-item>
            <a-form-item
              field="description"
              :label="t('topic.descriptionLabel')"
            >
              <a-textarea
                v-model="formData.description"
                :placeholder="t('topic.descriptionLabel')"
                :auto-size="{ minRows: 3, maxRows: 6 }"
              />
            </a-form-item>
          </section>

          <section v-show="currentStep === 2" class="step-panel">
            <a-alert type="info" class="mb-4" show-icon>
              {{ t('topic.inputsHelp') }}
            </a-alert>
            <TopicInputSourcesEditor
              v-model="formData.inputSources"
              :available-recipes="availableRecipes"
              :available-topics="availableTopics"
              :picker-loading="pickerLoading"
              :exclude-topic-id="isEdit ? formData.id : undefined"
            />
          </section>

          <section v-show="currentStep === 3" class="step-panel">
            <a-alert type="info" class="mb-4" show-icon>
              {{ t('topic.aggregatorHelp') }}
            </a-alert>
            <TopicAggregatorEditor v-model="formData.aggregator_config" />
          </section>

          <section v-show="currentStep === 4" class="step-panel">
            <a-descriptions :column="1" bordered>
              <a-descriptions-item :label="t('topic.id')">
                {{ formData.id || '-' }}
              </a-descriptions-item>
              <a-descriptions-item :label="t('topic.title')">
                {{ formData.title || '-' }}
              </a-descriptions-item>
              <a-descriptions-item :label="t('topic.descriptionLabel')">
                {{ formData.description || '-' }}
              </a-descriptions-item>
              <a-descriptions-item :label="t('topic.inputCount')">
                {{ configuredInputCount }}
              </a-descriptions-item>
              <a-descriptions-item :label="t('topic.aggregator')">
                {{
                  formatAggregatorSummary(
                    normalizeTopicPayload(formData).aggregator_config,
                    t
                  )
                }}
              </a-descriptions-item>
            </a-descriptions>

            <div class="section-label mt-4">{{ t('topic.inputs') }}</div>
            <a-list bordered class="mb-4">
              <a-list-item
                v-for="(source, idx) in reviewInputs"
                :key="`review-input-${idx}`"
              >
                <a-list-item-meta
                  :title="source.description || source.uri"
                  :description="source.description ? source.uri : undefined"
                />
                <template #actions>
                  <a-tag v-if="source.disabled" color="gray">
                    {{ t('topic.inputDisabled.badge') }}
                  </a-tag>
                </template>
              </a-list-item>
            </a-list>

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
          </section>
        </a-form>
      </a-spin>

      <div class="wizard-footer">
        <a-space>
          <a-button @click="goBackToList">{{ t('topic.cancel') }}</a-button>
          <a-button v-if="currentStep > 1" @click="prevStep">
            {{ t('topic.wizard.prev') }}
          </a-button>
          <a-button
            v-if="currentStep < 4"
            type="primary"
            :loading="stepValidating"
            @click="nextStep"
          >
            {{ t('topic.wizard.next') }}
          </a-button>
          <template v-else>
            <a-button :loading="validating" @click="handleValidate">
              {{ t('topic.validate') }}
            </a-button>
            <a-button
              type="primary"
              :loading="submitting"
              @click="handleSubmit"
            >
              {{ t('topic.save') }}
            </a-button>
          </template>
        </a-space>
      </div>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
  import { computed, onMounted, ref } from 'vue';
  import { Message } from '@arco-design/web-vue';
  import { useI18n } from 'vue-i18n';
  import { useRoute, useRouter } from 'vue-router';
  import { CustomRecipe, getCustomRecipes } from '@/api/custom_recipe';
  import {
    TopicFeed,
    TopicValidationIssue,
    createTopicFeed,
    getTopicFeed,
    listTopicFeeds,
    updateTopicFeed,
    validateTopicFeed,
  } from '@/api/topic';
  import XHeader from '@/components/header/x-header.vue';
  import { getRecipeIdRules } from '@/utils/slug';
  import TopicAggregatorEditor from '@/views/dashboard/topic_feed/components/TopicAggregatorEditor.vue';
  import TopicInputSourcesEditor from '@/views/dashboard/topic_feed/components/TopicInputSourcesEditor.vue';
  import {
    defaultFormData,
    formatAggregatorSummary,
    countEnabledInputs,
    normalizeTopicPayload,
    topicFeedToFormData,
    type TopicFormData,
  } from '@/views/dashboard/topic_feed/topicInputUtils';

  const { t } = useI18n();
  const route = useRoute();
  const router = useRouter();

  const isEdit = computed(() => Boolean(route.params.id));
  const pageTitle = computed(() =>
    isEdit.value ? t('topic.edit') : t('topic.create')
  );
  const breadcrumbActionLabel = computed(() =>
    isEdit.value ? t('topic.edit') : t('topic.create')
  );

  const currentStep = ref(1);
  const pageLoading = ref(false);
  const stepValidating = ref(false);
  const submitting = ref(false);
  const validating = ref(false);
  const formRef = ref();
  const formData = ref<TopicFormData>(defaultFormData());
  const validationErrors = ref<TopicValidationIssue[]>([]);
  const validationWarnings = ref<TopicValidationIssue[]>([]);
  const availableRecipes = ref<CustomRecipe[]>([]);
  const availableTopics = ref<TopicFeed[]>([]);
  const pickerLoading = ref(false);

  const configuredInputCount = computed(() =>
    countEnabledInputs(formData.value.inputSources)
  );

  const reviewInputs = computed(
    () => normalizeTopicPayload(formData.value).inputs || []
  );

  const loadPickerData = async () => {
    pickerLoading.value = true;
    try {
      const [recipesRes, topicsRes] = await Promise.all([
        getCustomRecipes(),
        listTopicFeeds(),
      ]);
      availableRecipes.value = recipesRes.data ?? [];
      availableTopics.value = topicsRes.data ?? [];
    } finally {
      pickerLoading.value = false;
    }
  };

  const loadTopic = async () => {
    if (!isEdit.value) return;
    pageLoading.value = true;
    try {
      const res = await getTopicFeed(String(route.params.id));
      if (!res.data) {
        Message.error(t('topic.detail.notFound'));
        goBackToList();
        return;
      }
      formData.value = topicFeedToFormData(res.data);
    } catch (err: any) {
      Message.error(err.message || t('topic.fetchFailed'));
      goBackToList();
    } finally {
      pageLoading.value = false;
    }
  };

  const goBackToList = () => {
    router.push({ name: 'TopicFeed' });
  };

  const prevStep = () => {
    currentStep.value = Math.max(1, currentStep.value - 1);
  };

  const validateCurrentStep = async () => {
    if (currentStep.value === 1) {
      const result = await formRef.value?.validateField('id');
      return !result;
    }
    if (currentStep.value === 2) {
      if (configuredInputCount.value === 0) {
        Message.warning(t('topic.wizard.inputsRequired'));
        return false;
      }
      return true;
    }
    return true;
  };

  const nextStep = async () => {
    stepValidating.value = true;
    try {
      const ok = await validateCurrentStep();
      if (!ok) return;
      currentStep.value = Math.min(4, currentStep.value + 1);
    } finally {
      stepValidating.value = false;
    }
  };

  const runValidation = async () => {
    const payload = normalizeTopicPayload(formData.value);
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
    const formResult = await formRef.value?.validate();
    if (formResult) return;

    if (configuredInputCount.value === 0) {
      Message.warning(t('topic.wizard.inputsRequired'));
      currentStep.value = 2;
      return;
    }

    submitting.value = true;
    try {
      const result = await runValidation();
      if (!result?.valid) {
        Message.error(t('topic.validateFailed'));
        return;
      }

      const payload = normalizeTopicPayload(formData.value);
      if (isEdit.value) {
        await updateTopicFeed(payload.id, payload);
        Message.success(t('topic.updateSuccess'));
      } else {
        await createTopicFeed(payload);
        Message.success(t('topic.createSuccess'));
      }
      goBackToList();
    } catch (err: any) {
      Message.error(err.message || t('topic.saveFailed'));
    } finally {
      submitting.value = false;
    }
  };

  onMounted(async () => {
    await Promise.all([loadPickerData(), loadTopic()]);
  });
</script>

<script lang="ts">
  export default {
    name: 'TopicFeedEditor',
  };
</script>

<style scoped>
  .wizard-card {
    margin-top: 8px;
  }

  .wizard-card :deep(.arco-card-body) {
    display: flex;
    flex-direction: column;
    min-height: 560px;
  }

  .wizard-steps {
    margin-bottom: 24px;
    padding-bottom: 20px;
    border-bottom: 1px solid var(--color-border-2);
  }

  .wizard-spin {
    display: block;
    width: 100%;
    flex: 1;
  }

  .wizard-form {
    width: 100%;
  }

  .wizard-form :deep(.arco-form-item) {
    display: block;
    width: 100%;
  }

  .wizard-form :deep(.arco-form-item-label-col),
  .wizard-form :deep(.arco-form-item-wrapper-col) {
    display: block;
    width: 100%;
    max-width: 100%;
    flex: none;
    justify-content: flex-start;
    align-items: stretch;
    padding-right: 0;
  }

  .wizard-form :deep(.arco-form-item-content-wrapper),
  .wizard-form :deep(.arco-form-item-content),
  .wizard-form :deep(.arco-input-wrapper),
  .wizard-form :deep(.arco-textarea-wrapper) {
    width: 100%;
  }

  .step-panel {
    width: 100%;
    max-width: 720px;
    margin: 0;
  }

  .step-panel--basic {
    max-width: 560px;
  }

  .wizard-footer {
    margin-top: auto;
    padding-top: 16px;
    border-top: 1px solid var(--color-border-2);
    display: flex;
    justify-content: flex-end;
  }

  .section-label {
    margin-bottom: 12px;
    font-weight: 600;
  }

  .validation-item {
    margin-top: 6px;
  }
</style>
