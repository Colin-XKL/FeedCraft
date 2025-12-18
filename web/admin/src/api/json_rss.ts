import axios from 'axios';
import { util } from './interceptor';

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
}

export function parseCurl(curlCommand: string): Promise<util.APIResponse<JsonFetchReq>> {
  return axios.post<util.APIResponse<JsonFetchReq>>(
    '/api/admin/tools/json/parse_curl',
    {
      curl_command: curlCommand,
    }
  ) as unknown as Promise<util.APIResponse<JsonFetchReq>>;
}

export function fetchJson(req: JsonFetchReq): Promise<util.APIResponse<string>> {
  return axios.post<util.APIResponse<string>>('/api/admin/tools/json/fetch', req) as unknown as Promise<util.APIResponse<string>>;
}

export function parseJsonRss(req: JsonParseReq): Promise<util.APIResponse<ParsedItem[]>> {
  return axios.post<util.APIResponse<ParsedItem[]>>(
    '/api/admin/tools/json/parse',
    req
  ) as unknown as Promise<util.APIResponse<ParsedItem[]>>;
}
