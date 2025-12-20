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
