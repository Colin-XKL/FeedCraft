import { describe, expect, it } from 'vitest';
import {
  normalizeTopicPayload,
  parseUriToSource,
  sourceToUri,
  topicFeedToFormData,
} from '@/views/dashboard/topic_feed/topicInputUtils';

describe('topic input URI mapping', () => {
  it('parses feedcraft inbox URIs as inbox sources', () => {
    expect(
      parseUriToSource({
        uri: 'feedcraft://inbox/alerts',
        description: 'Ops alerts',
        disabled: false,
      })
    ).toEqual({
      sourceType: 'inbox',
      externalUrl: '',
      resourceId: 'alerts',
      description: 'Ops alerts',
      disabled: false,
    });
  });

  it('serializes inbox sources to feedcraft inbox URIs', () => {
    expect(
      sourceToUri({
        sourceType: 'inbox',
        externalUrl: '',
        resourceId: 'alerts',
        description: 'Ops alerts',
        disabled: false,
      })
    ).toBe('feedcraft://inbox/alerts');
  });

  it('round-trips inbox inputs through form data', () => {
    const formData = topicFeedToFormData({
      id: 'daily',
      title: 'Daily',
      description: '',
      inputs: [{ uri: 'feedcraft://inbox/alerts', description: 'Inbox' }],
      aggregator_config: [],
    });

    expect(formData.inputSources[0].sourceType).toBe('inbox');
    expect(formData.inputSources[0].resourceId).toBe('alerts');
    expect(normalizeTopicPayload(formData).inputs).toEqual([
      {
        uri: 'feedcraft://inbox/alerts',
        description: 'Inbox',
        disabled: false,
      },
    ]);
  });
});
