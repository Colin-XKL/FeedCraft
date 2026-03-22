type TranslateFn = (key: string) => string;

function withFallback(t: TranslateFn, key: string, fallback: string): string {
  const translated = t(key);
  return translated === key ? fallback : translated;
}

export function formatObservabilityStatus(
  t: TranslateFn,
  status?: string
): string {
  if (!status) return '-';
  return withFallback(t, `observability.statusValue.${status}`, status);
}

export function formatObservabilityTrigger(
  t: TranslateFn,
  trigger?: string
): string {
  if (!trigger) return '-';
  return withFallback(t, `observability.triggerValue.${trigger}`, trigger);
}

export function formatObservabilityErrorKind(
  t: TranslateFn,
  kind?: string
): string {
  if (!kind) return '-';
  return withFallback(t, `observability.errorKind.${kind}`, kind);
}

export function formatObservabilityResourceType(
  t: TranslateFn,
  type?: string
): string {
  if (!type) return '-';
  return withFallback(t, `observability.resourceType.${type}`, type);
}
