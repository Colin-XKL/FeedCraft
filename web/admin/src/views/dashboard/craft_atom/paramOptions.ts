export type CraftParamValue = number | string | string[];

export const aiFilterExtraPayloadOptions = [
  {
    label: 'article_summary',
    value: 'article_summary',
  },
  {
    label: 'article_content',
    value: 'article_content',
  },
  {
    label: 'article_date',
    value: 'article_date',
  },
  {
    label: 'raw_rss_item',
    value: 'raw_rss_item',
  },
];

const AI_FILTER_TEMPLATE = 'ai-filter';
const AI_CONTENT_PROCESS_TEMPLATE = 'ai-content-process';
const AI_FILTER_EXTRA_PAYLOAD_PARAM = 'extra-payload';
const AI_CONTENT_PROCESS_PLACEMENT_PARAM = 'placement';
const EMBEDDING_FILTER_TEMPLATE = 'embedding-filter';
const EMBEDDING_FILTER_ANCHORS_PARAM = 'anchors';
const EMBEDDING_FILTER_MODE_PARAM = 'mode';
const EMBEDDING_FILTER_THRESHOLD_PARAM = 'threshold';
const EMBEDDING_FILTER_MAX_CONTENT_LENGTH_PARAM = 'max_content_length';

export const embeddingFilterModeOptions = [
  {
    label: 'include',
    value: 'include',
  },
  {
    label: 'exclude',
    value: 'exclude',
  },
];

export const aiContentProcessPlacementOptions = [
  {
    label: 'prepend',
    value: 'prepend',
  },
  {
    label: 'replace',
    value: 'replace',
  },
  {
    label: 'append',
    value: 'append',
  },
];

export function isAIFilterExtraPayloadParam(
  templateName: string,
  paramKey: string
) {
  return (
    [AI_FILTER_TEMPLATE, AI_CONTENT_PROCESS_TEMPLATE].includes(templateName) &&
    paramKey === AI_FILTER_EXTRA_PAYLOAD_PARAM
  );
}

export function isAIContentProcessPlacementParam(
  templateName: string,
  paramKey: string
) {
  return (
    templateName === AI_CONTENT_PROCESS_TEMPLATE &&
    paramKey === AI_CONTENT_PROCESS_PLACEMENT_PARAM
  );
}

export function isEmbeddingFilterParam(templateName: string, paramKey: string) {
  return templateName === EMBEDDING_FILTER_TEMPLATE && Boolean(paramKey);
}

export function isEmbeddingFilterAnchorsParam(
  templateName: string,
  paramKey: string
) {
  return (
    templateName === EMBEDDING_FILTER_TEMPLATE &&
    paramKey === EMBEDDING_FILTER_ANCHORS_PARAM
  );
}

export function isEmbeddingFilterModeParam(
  templateName: string,
  paramKey: string
) {
  return (
    templateName === EMBEDDING_FILTER_TEMPLATE &&
    paramKey === EMBEDDING_FILTER_MODE_PARAM
  );
}

export function isEmbeddingFilterThresholdParam(
  templateName: string,
  paramKey: string
) {
  return (
    templateName === EMBEDDING_FILTER_TEMPLATE &&
    paramKey === EMBEDDING_FILTER_THRESHOLD_PARAM
  );
}

export function isEmbeddingFilterMaxContentLengthParam(
  templateName: string,
  paramKey: string
) {
  return (
    templateName === EMBEDDING_FILTER_TEMPLATE &&
    paramKey === EMBEDDING_FILTER_MAX_CONTENT_LENGTH_PARAM
  );
}

function parseNumberParam(value: string | undefined, fallback: number) {
  const parsed = Number(value);
  return Number.isFinite(parsed) ? parsed : fallback;
}

export function deserializeAIFilterExtraPayloadValue(value?: string) {
  return (value || '')
    .split(/[,\n\t|]+/)
    .map((item) => item.trim())
    .filter(Boolean);
}

export function toCraftParamFormValue(
  templateName: string,
  paramKey: string,
  value?: string
): CraftParamValue {
  if (isAIFilterExtraPayloadParam(templateName, paramKey)) {
    return deserializeAIFilterExtraPayloadValue(value);
  }
  if (isEmbeddingFilterThresholdParam(templateName, paramKey)) {
    return parseNumberParam(value, 0.6);
  }
  if (isEmbeddingFilterMaxContentLengthParam(templateName, paramKey)) {
    return parseNumberParam(value, 2000);
  }
  return value || '';
}

export function serializeCraftParamValue(value: CraftParamValue) {
  if (Array.isArray(value)) {
    return value.join(',');
  }
  if (typeof value === 'number') {
    return String(value);
  }
  return value;
}
