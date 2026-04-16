export const IGNORED_CLASSES = [
  'flex',
  'row',
  'col',
  'grid',
  'hidden',
  'block',
  'items-center',
  'justify-center',
  'fc-highlight',
];

/**
 * Generates a CSS selector path for a given element.
 * @param el The target HTML element
 * @param isItemSelector If true, avoids using unique IDs to allow selecting multiple similar items
 * @param rootDocument Optional document reference (used to check boundaries)
 */
export const getCssSelector = (
  el: HTMLElement,
  rootDocument?: Document, // Optional parameter moved before default parameter
  isItemSelector = false // Default parameter
): string => {
  if (!el || el.nodeType !== Node.ELEMENT_NODE) return '';

  // 1. If ID exists and we are not looking for a generic list item selector, use it.
  if (el.id && !isItemSelector) {
    return `#${CSS.escape(el.id)}`;
  }

  const path: string[] = [];
  let currentEl: HTMLElement | null = el;

  const body = rootDocument?.body;
  const html = rootDocument?.documentElement;

  while (currentEl && currentEl.nodeType === Node.ELEMENT_NODE) {
    let selector = currentEl.nodeName.toLowerCase();

    // Try to use class if available and meaningful
    if (currentEl.classList.length > 0) {
      const classes = Array.from(currentEl.classList).filter(
        (c) => !IGNORED_CLASSES.includes(c)
      );
      if (classes.length > 0) {
        selector += `.${classes.map((c) => CSS.escape(c)).join('.')}`;
      }
    }

    // Special handling for the target element when selecting a list item
    if (currentEl === el && isItemSelector) {
      // Do NOT add :nth-of-type to the target element itself to ensure we match all siblings
    } else if (currentEl.id) {
      // For parent path, or normal selection, ID is the best anchor
      selector = `#${CSS.escape(currentEl.id)}`;
      path.unshift(selector);
      break;
    } else {
      // Check siblings for nth-of-type
      let sib = currentEl;
      let nth = 1;
      // eslint-disable-next-line no-cond-assign
      while ((sib = sib.previousElementSibling as HTMLElement)) {
        if (sib.nodeName.toLowerCase() === currentEl.nodeName.toLowerCase()) {
          nth += 1;
        }
      }
      if (nth !== 1) selector += `:nth-of-type(${nth})`;
    }

    path.unshift(selector);
    currentEl = currentEl.parentNode as HTMLElement;

    // Stop if we hit the container boundaries
    if (!currentEl || currentEl === body || currentEl === html) break;
  }

  return path.join(' > ');
};
