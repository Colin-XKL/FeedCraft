import type { LocationQueryRaw } from 'vue-router';
import type { APIResponse } from '@/api/types';

export const SESSION_EXPIRED_REASON = 'session-expired';
export const SESSION_EXPIRED_MESSAGE = '登录态已过期，请重新登录后继续使用。';

const SESSION_EXPIRED_CODES = [50008, 50012, 50014];

export class SessionExpiredError extends Error {
  constructor() {
    super(SESSION_EXPIRED_MESSAGE);
    this.name = 'SessionExpiredError';
  }
}

export function isSessionExpiredError(error: unknown) {
  return error instanceof SessionExpiredError;
}

export function isSessionExpiredHTTPStatus(status?: number) {
  return status === 401;
}

export function isSessionExpiredAPIResponse(
  response?: Pick<APIResponse, 'status' | 'code'>
) {
  const code = response?.code;
  return (
    isSessionExpiredHTTPStatus(response?.status) ||
    (typeof code === 'number' && SESSION_EXPIRED_CODES.includes(code))
  );
}

export function buildSessionExpiredRedirectQuery(
  currentQuery: LocationQueryRaw,
  redirect: string
): LocationQueryRaw {
  return {
    ...currentQuery,
    redirect,
    reason: SESSION_EXPIRED_REASON,
  };
}
