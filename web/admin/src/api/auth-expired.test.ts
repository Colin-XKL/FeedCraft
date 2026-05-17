import { describe, expect, it } from 'vitest';
import {
  buildSessionExpiredRedirectQuery,
  isSessionExpiredAPIResponse,
  isSessionExpiredHTTPStatus,
} from './auth-expired';

describe('auth expired helpers', () => {
  it('identifies unauthorized responses as session expiration', () => {
    expect(isSessionExpiredHTTPStatus(401)).toBe(true);
    expect(isSessionExpiredHTTPStatus(403)).toBe(false);
    expect(isSessionExpiredAPIResponse({ status: 401, code: 1 })).toBe(true);
    expect(isSessionExpiredAPIResponse({ status: 500, code: 50014 })).toBe(
      true
    );
    expect(isSessionExpiredAPIResponse({ status: 500, code: 10001 })).toBe(
      false
    );
  });

  it('preserves the current page as the login redirect target', () => {
    expect(
      buildSessionExpiredRedirectQuery(
        { page: '2', keyword: 'rss' },
        '/worktable/craft_atom?page=2&keyword=rss'
      )
    ).toEqual({
      page: '2',
      keyword: 'rss',
      redirect: '/worktable/craft_atom?page=2&keyword=rss',
      reason: 'session-expired',
    });
  });
});
