import { describe, expect, it } from 'vitest';
import { namingValidator } from '@/utils/validator';

describe('namingValidator', () => {
  it('skips empty values so required rules can own that case', () => {
    let called: string | undefined = 'unset';
    namingValidator.validator('', (err) => {
      called = err;
    });
    expect(called).toBeUndefined();
  });

  it('rejects uppercase and other invalid characters', () => {
    let message: string | undefined;
    namingValidator.validator('Bad Name', (err) => {
      message = err;
    });
    expect(message).toMatch(/lowercase letters/);
  });

  it('accepts lowercase ids', () => {
    let called = false;
    namingValidator.validator('e2e-ai-recipe-887', (err) => {
      called = true;
      expect(err).toBeUndefined();
    });
    expect(called).toBe(true);
  });
});
