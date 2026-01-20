import axios from 'axios';
import { APIResponse } from './types';

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
  enhanced_mode?: boolean;
}

export function parseCurl(
  curlCommand: string,
): Promise<APIResponse<JsonFetchReq>> {
  return axios
    .post<APIResponse<JsonFetchReq>>('/api/admin/tools/json/parse_curl', {
      curl_command: curlCommand,
    })
    .then((res) => res.data);
}

export function fetchJson(req: JsonFetchReq): Promise<APIResponse<string>> {
  return axios
    .post<APIResponse<string>>('/api/admin/tools/json/fetch', req)
    .then((res) => res.data);
}

export function parseJsonRss(
  req: JsonParseReq,
): Promise<APIResponse<ParsedItem[]>> {
  return axios
    .post<APIResponse<ParsedItem[]>>('/api/admin/tools/json/parse', req)
    .then((res) => res.data);
}
export function previewSearch(
  req: SearchFetchReq,
): Promise<APIResponse<any[]>> {
  return axios
    .post<APIResponse<any[]>>('/api/admin/tools/search/preview', req)
    .then((res) => res.data);
}
