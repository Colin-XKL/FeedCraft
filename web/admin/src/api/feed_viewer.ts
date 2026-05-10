import axios from 'axios';
import { APIResponse } from '@/api/types';

export interface FeedViewerPreviewImage {
  url: string;
  title: string;
}

export interface FeedViewerPreviewItem {
  guid: string;
  title: string;
  link: string;
  pubDate: string;
  isoDate: string;
  content: string;
  contentSnippet: string;
}

export interface FeedViewerPreview {
  title: string;
  description: string;
  link: string;
  feedUrl: string;
  copyright: string;
  image?: FeedViewerPreviewImage;
  items: FeedViewerPreviewItem[];
}

export interface PreviewFeedOptions {
  craftName?: string;
}

export function previewFeed(
  inputUrl: string,
  options: PreviewFeedOptions = {}
): Promise<APIResponse<FeedViewerPreview>> {
  return axios
    .get<APIResponse<FeedViewerPreview>>('/api/admin/tools/feed/preview', {
      params: {
        input_url: inputUrl,
        craft_name: options.craftName,
      },
    })
    .then((res) => res.data);
}
