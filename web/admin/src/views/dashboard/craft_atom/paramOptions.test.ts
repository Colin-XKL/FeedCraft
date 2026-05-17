import { describe, expect, it } from 'vitest';
import {
  aiContentProcessPlacementOptions,
  isAIFilterExtraPayloadParam,
  isAIContentProcessPlacementParam,
  serializeCraftParamValue,
  toCraftParamFormValue,
} from './paramOptions';

describe('craft atom param options', () => {
  it('treats ai-content-process extra-payload like ai-filter extra-payload', () => {
    expect(isAIFilterExtraPayloadParam('ai-filter', 'extra-payload')).toBe(
      true
    );
    expect(
      isAIFilterExtraPayloadParam('ai-content-process', 'extra-payload')
    ).toBe(true);
    expect(
      toCraftParamFormValue(
        'ai-content-process',
        'extra-payload',
        'article_content|article_date'
      )
    ).toEqual(['article_content', 'article_date']);
    expect(serializeCraftParamValue(['article_content', 'article_date'])).toBe(
      'article_content,article_date'
    );
  });

  it('exposes ai-content-process placement as a single-select option set', () => {
    expect(
      isAIContentProcessPlacementParam('ai-content-process', 'placement')
    ).toBe(true);
    expect(isAIContentProcessPlacementParam('ai-filter', 'placement')).toBe(
      false
    );
    expect(
      aiContentProcessPlacementOptions.map((option) => option.value)
    ).toEqual(['prepend', 'replace', 'append']);
  });
});
