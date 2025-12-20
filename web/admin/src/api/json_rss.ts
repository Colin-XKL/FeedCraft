import axios from 'axios';

export interface JsonFetchReq {
  method: string;
  url: string;
  headers: Record<string, string>;
  body: string;
}

export interface JsonParseReq {
  json_content: string;
  list_selector: string;
  title_selector: string;
  link_selector: string;
  date_selector: string;
  content_selector: string;
}

export interface ParsedItem {
  title: string;
  link: string;
  date: string;
  content: string;
  description: string;
}

export interface SearchFetchReq {
  query: string;
}

export function parseCurl(curlCommand: string) {
  return axios.post<JsonFetchReq>('/api/admin/tools/json/parse_curl', {
    curl_command: curlCommand,
  });
}

export function fetchJson(req: JsonFetchReq) {
  return axios.post<string>('/api/admin/tools/json/fetch', req);
}

export function previewSearch(req: SearchFetchReq) {
  return axios.post<ParsedItem[]>('/api/admin/tools/search/preview', req);
}

export function parseJsonRss(req: JsonParseReq) {
  return axios.post<ParsedItem[]>('/api/admin/tools/json/parse', req);
}
