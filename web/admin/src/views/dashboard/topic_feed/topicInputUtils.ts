import type { AggregatorStep, TopicFeed, TopicInput } from '@/api/topic';

export type StepType = 'deduplicate' | 'sort' | 'limit';
export type SourceType = 'external' | 'recipe' | 'topic';

export const STRATEGIES_WITH_THRESHOLD = [
  'by_simhash',
  'by_embedding',
] as const;

export interface InputSourceItem {
  sourceType: SourceType;
  externalUrl: string;
  resourceId: string;
  description: string;
  disabled: boolean;
}

export interface StepFormItem {
  type: StepType;
  value: string | number;
  threshold?: number;
}

export interface TopicFormData {
  id: string;
  title: string;
  description: string;
  inputSources: InputSourceItem[];
  aggregator_config: StepFormItem[];
}

export const defaultThreshold = (strategy: string): number | undefined => {
  if (strategy === 'by_simhash') return 0.05;
  if (strategy === 'by_embedding') return 0.1;
  return undefined;
};

export const createDefaultStep = (type: StepType = 'limit'): StepFormItem => {
  if (type === 'deduplicate') return { type, value: 'by_link' };
  if (type === 'sort') return { type, value: 'date_desc' };
  return { type, value: 50 };
};

export const defaultFormData = (): TopicFormData => ({
  id: '',
  title: '',
  description: '',
  inputSources: [
    {
      sourceType: 'external',
      externalUrl: '',
      resourceId: '',
      description: '',
      disabled: false,
    },
  ],
  aggregator_config: [],
});

export const parseUriToSource = (
  input: TopicInput | string
): InputSourceItem => {
  const uri = typeof input === 'string' ? input : input.uri;
  const description =
    typeof input === 'string' ? '' : input.description?.trim() || '';
  const disabled = typeof input === 'string' ? false : Boolean(input.disabled);

  if (uri.startsWith('feedcraft://recipe/')) {
    return {
      sourceType: 'recipe',
      externalUrl: '',
      resourceId: uri.slice('feedcraft://recipe/'.length),
      description,
      disabled,
    };
  }
  if (uri.startsWith('feedcraft://topic/')) {
    return {
      sourceType: 'topic',
      externalUrl: '',
      resourceId: uri.slice('feedcraft://topic/'.length),
      description,
      disabled,
    };
  }
  return {
    sourceType: 'external',
    externalUrl: uri,
    resourceId: '',
    description,
    disabled,
  };
};

export const sourceToUri = (source: InputSourceItem): string => {
  if (source.sourceType === 'recipe') {
    return `feedcraft://recipe/${source.resourceId}`;
  }
  if (source.sourceType === 'topic') {
    return `feedcraft://topic/${source.resourceId}`;
  }
  return source.externalUrl.trim();
};

export const countEnabledInputs = (sources: InputSourceItem[]): number =>
  sources.filter((source) => sourceToUri(source) !== '' && !source.disabled)
    .length;

export const topicFeedToFormData = (record: TopicFeed): TopicFormData => {
  const inputs =
    record.inputs && record.inputs.length > 0
      ? record.inputs
      : record.input_uris.map((uri) => ({
          uri,
          description: '',
          disabled: false,
        }));

  return {
    id: record.id,
    title: record.title || '',
    description: record.description || '',
    inputSources:
      inputs.length > 0
        ? inputs.map(parseUriToSource)
        : [
            {
              sourceType: 'external',
              externalUrl: '',
              resourceId: '',
              description: '',
              disabled: false,
            },
          ],
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
};

export const normalizeTopicPayload = (formData: TopicFormData): TopicFeed => {
  const inputs: TopicInput[] = formData.inputSources
    .map((source) => ({
      uri: sourceToUri(source),
      description: source.description.trim(),
      disabled: source.disabled,
    }))
    .filter((item) => item.uri !== '');

  return {
    id: formData.id.trim(),
    title: formData.title.trim(),
    description: formData.description.trim(),
    inputs,
    input_uris: inputs.filter((item) => !item.disabled).map((item) => item.uri),
    aggregator_config: formData.aggregator_config.map((step) => {
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
  };
};

export const formatAggregatorSummary = (
  steps: AggregatorStep[],
  t: (key: string) => string
): string => {
  if (!steps || steps.length === 0) return t('topic.noAggregator');
  return steps
    .map((step) => {
      if (step.type === 'deduplicate') {
        const strategy = step.option?.strategy || 'by_link';
        const label = t(`topic.stepOption.strategy.${strategy}`);
        if (
          (strategy === 'by_simhash' || strategy === 'by_embedding') &&
          step.option?.threshold
        ) {
          return `${t('topic.stepType.deduplicate')} · ${label} (${
            step.option.threshold
          })`;
        }
        return `${t('topic.stepType.deduplicate')} · ${label}`;
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
