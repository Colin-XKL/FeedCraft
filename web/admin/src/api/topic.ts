import axios from 'axios';
import { APIResponse } from './types';
import {
  ExecutionLog,
  ObservableResource,
  SystemNotification,
} from './observability';

export interface AggregatorStep {
  type: string;
  option: Record<string, string>;
}

export interface TopicFeed {
  id: string;
  title?: string;
  description?: string;
  input_uris: string[];
  aggregator_config: AggregatorStep[];
}

export interface TopicValidationIssue {
  field: string;
  message: string;
}

export interface TopicValidationResult {
  valid: boolean;
  errors: TopicValidationIssue[];
  warnings: TopicValidationIssue[];
}

export interface TopicDetail {
  topic: TopicFeed;
  public_url: string;
  health: ObservableResource;
  recent_executions: ExecutionLog[];
  related_notifications: SystemNotification[];
}

const adminApiBase = '/api/admin';

export function createTopicFeed(
  data: TopicFeed
): Promise<APIResponse<TopicFeed>> {
  return axios
    .post<APIResponse<TopicFeed>>(`${adminApiBase}/topics`, data)
    .then((res) => res.data);
}

export function listTopicFeeds(): Promise<APIResponse<TopicFeed[]>> {
  return axios
    .get<APIResponse<TopicFeed[]>>(`${adminApiBase}/topics`)
    .then((res) => res.data);
}

export function getTopicFeed(id: string): Promise<APIResponse<TopicFeed>> {
  return axios
    .get<APIResponse<TopicFeed>>(`${adminApiBase}/topics/${id}`)
    .then((res) => res.data);
}

export function updateTopicFeed(
  id: string,
  data: TopicFeed
): Promise<APIResponse<TopicFeed>> {
  return axios
    .put<APIResponse<TopicFeed>>(`${adminApiBase}/topics/${id}`, data)
    .then((res) => res.data);
}

export function deleteTopicFeed(id: string): Promise<APIResponse<void>> {
  return axios
    .delete<APIResponse<void>>(`${adminApiBase}/topics/${id}`)
    .then((res) => res.data);
}

export function validateTopicFeed(
  data: TopicFeed
): Promise<APIResponse<TopicValidationResult>> {
  return axios
    .post<APIResponse<TopicValidationResult>>(
      `${adminApiBase}/topics/validate`,
      data
    )
    .then((res) => res.data);
}

export function getTopicFeedDetail(
  id: string
): Promise<APIResponse<TopicDetail>> {
  return axios
    .get<APIResponse<TopicDetail>>(`${adminApiBase}/topics/${id}/detail`)
    .then((res) => res.data);
}
