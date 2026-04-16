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

export function previewFeed(
  inputUrl: string
): Promise<APIResponse<FeedViewerPreview>> {
  return axios
    .get<APIResponse<FeedViewerPreview>>('/api/admin/tools/feed/preview', {
      params: { input_url: inputUrl },
    })
    .then((res) => res.data);
}
