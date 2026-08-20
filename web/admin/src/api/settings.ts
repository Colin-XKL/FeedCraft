import axios from 'axios';

export interface SearchProviderConfig {
  api_url: string;
  api_key: string;
  provider: string;
  search_tool_name?: string;
}

export function getSearchProviderConfig() {
  return axios.get<SearchProviderConfig>('/api/admin/settings/search-provider');
}

export function saveSearchProviderConfig(data: SearchProviderConfig) {
  return axios.post('/api/admin/settings/search-provider', data);
}

export interface FaviconProviderConfig {
  id: string;
  name: string;
  url_template: string;
  enabled: boolean;
}

export interface FaviconProviderDescriptor extends FaviconProviderConfig {
  built_in: boolean;
}

export interface FaviconSettings {
  default_provider_id: string;
  custom_providers: FaviconProviderConfig[];
}

export interface FaviconSettingsResponse extends FaviconSettings {
  providers: FaviconProviderDescriptor[];
}

export interface FaviconPreviewResponse {
  url: string;
  provider_id: string;
}

interface APIResponse<T> {
  data: T;
  code: number;
  msg: string;
}

export function getFaviconProviderConfig() {
  return axios.get<APIResponse<FaviconSettingsResponse>>(
    '/api/admin/settings/favicon-provider'
  );
}

export function saveFaviconProviderConfig(data: FaviconSettings) {
  return axios.post<APIResponse<FaviconSettingsResponse>>(
    '/api/admin/settings/favicon-provider',
    data
  );
}

export function previewFaviconProviderConfig(
  settings: FaviconSettings,
  providerId: string,
  pageUrl: string,
  size = 64
) {
  return axios.post<APIResponse<FaviconPreviewResponse>>(
    '/api/admin/settings/favicon-provider/preview',
    {
      settings,
      provider_id: providerId,
      page_url: pageUrl,
      size,
    }
  );
}
