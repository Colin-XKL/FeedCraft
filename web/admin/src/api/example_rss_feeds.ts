import axios from 'axios';
import { APIResponse } from '@/api/types';

export interface ExampleRssFeed {
  slug: string;
  title: string;
  description: string;
  path: string;
}

export function listExampleRssFeeds(): Promise<APIResponse<ExampleRssFeed[]>> {
  return axios
    .get<APIResponse<ExampleRssFeed[]>>('/api/example-rss-feeds')
    .then((res) => res.data);
}
