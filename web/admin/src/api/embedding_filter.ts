import axios from 'axios';
import { APIResponse } from '@/api/types';
import { FeedViewerPreview } from '@/api/feed_viewer';

export interface EmbeddingFilterPreviewRequest {
  input_url: string;
  anchors: string;
  threshold?: number;
  mode: 'include' | 'exclude';
  max_content_length?: number;
  instruction?: string;
  atom_craft_name?: string;
  atom_craft_description?: string;
}

export function previewEmbeddingFilter(
  data: EmbeddingFilterPreviewRequest
): Promise<APIResponse<FeedViewerPreview>> {
  return axios
    .post<APIResponse<FeedViewerPreview>>(
      '/api/admin/tools/embedding-filter/preview',
      data
    )
    .then((res) => res.data);
}
