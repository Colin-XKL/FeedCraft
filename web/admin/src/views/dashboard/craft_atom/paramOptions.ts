export type CraftParamValue = string | string[];

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
const AI_FILTER_EXTRA_PAYLOAD_PARAM = 'extra-payload';

export function isAIFilterExtraPayloadParam(
  templateName: string,
  paramKey: string
) {
  return (
    templateName === AI_FILTER_TEMPLATE &&
    paramKey === AI_FILTER_EXTRA_PAYLOAD_PARAM
  );
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
  return value || '';
}

export function serializeCraftParamValue(value: CraftParamValue) {
  if (Array.isArray(value)) {
    return value.join(',');
  }
  return value;
}
