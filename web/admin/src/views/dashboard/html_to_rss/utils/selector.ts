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
 * @param rootDocument Optional document reference (used to check boundaries)
 * @param isItemSelector If true, avoids using unique IDs to allow selecting multiple similar items
 * @param stopAtElement Optional element to stop traversal at (exclusive). The path will be relative to this element.
 */
export const getCssSelector = (
  el: HTMLElement,
  rootDocument?: Document,
  isItemSelector = false,
  stopAtElement?: HTMLElement,
): string => {
  if (!el || el.nodeType !== Node.ELEMENT_NODE) return '';

  // 1. If ID exists and we are not looking for a generic list item selector, use it.
  // But ONLY if we are not restricted by stopAtElement (relative path shouldn't use ID if it's inside the item usually, unless we want strict selection)
  // Actually, for relative selection within an item, ID might be too specific if it's unique per item.
  // But usually items don't have unique IDs inside repetitive structures.
  // If they do, isItemSelector should be true? No, relative selector is for single fields.
  // Let's keep it simple: use ID if available and not forbidden.
  if (el.id && !isItemSelector) {
    return `#${CSS.escape(el.id)}`;
  }

  const path: string[] = [];
  let currentEl: HTMLElement | null = el;

  const body = rootDocument?.body;
  const html = rootDocument?.documentElement;

  while (currentEl && currentEl.nodeType === Node.ELEMENT_NODE) {
    // Check stop condition at the BEGINNING of the loop
    // If currentEl is the stopAtElement, we are done. We don't include it in the path.
    if (stopAtElement && currentEl === stopAtElement) {
      break;
    }

    let selector = currentEl.nodeName.toLowerCase();

    // Try to use class if available and meaningful
    if (currentEl.classList.length > 0) {
      const classes = Array.from(currentEl.classList).filter(
        (c) => !IGNORED_CLASSES.includes(c),
      );
      if (classes.length > 0) {
        selector += `.${classes.map((c) => CSS.escape(c)).join('.')}`;
      }
    }

    // Special handling for the target element when selecting a list item
    if (currentEl === el && isItemSelector) {
      // Do NOT add :nth-of-type to the target element itself to ensure we match all siblings
    } else if (currentEl.id && !stopAtElement) {
      // For parent path (absolute), ID is the best anchor
      // For relative path (stopAtElement provided), ID might be okay if it exists,
      // but usually we rely on structure within the item.
      // If we encounter an ID during traversal up to stopAtElement, we can use it?
      // Yes, if there is an ID inside the item, it's specific.
      selector = `#${CSS.escape(currentEl.id)}`;
      path.unshift(selector);
      break; // Absolute path found
    } else {
      // Check siblings for nth-of-type
      // We need to count siblings that match the SAME tag
      let sib = currentEl;
      let nth = 1;
      let hasSibling = false;

      // Check previous siblings
      // eslint-disable-next-line no-cond-assign
      while ((sib = sib.previousElementSibling as HTMLElement)) {
        if (sib.nodeName.toLowerCase() === currentEl.nodeName.toLowerCase()) {
          nth += 1;
          hasSibling = true;
        }
      }

      // We should also check next siblings to decide if nth-of-type is needed at all?
      // If it's the only child of that type, we don't strictly need nth-of-type, but adding it is safer.
      // However, usually we only add it if nth > 1 OR there are other siblings of same type.
      // To keep it simple and consistent with previous logic:
      if (nth !== 1) {
        selector += `:nth-of-type(${nth})`;
      } else {
        // If it's the first one, check if there are others
        let next = currentEl;
        // eslint-disable-next-line no-cond-assign
        while ((next = next.nextElementSibling as HTMLElement)) {
          if (
            next.nodeName.toLowerCase() === currentEl.nodeName.toLowerCase()
          ) {
            hasSibling = true;
            break;
          }
        }
        if (hasSibling) {
          // If there are siblings of same type, specify nth-of-type(1) for clarity?
          // Or just leave it. CSS :nth-of-type(1) is implied if we don't specify? No.
          // div means all divs. div:nth-of-type(1) means first.
          // If we want a specific path, we should probably always include it if there are siblings?
          // The previous code only added if nth !== 1. Let's stick to that to minimize changes,
          // but arguably for strict selection (like relative fields), specificity is good.
          // Let's add it if we are doing relative selection (stopAtElement defined) or if strictness is needed?
          // The previous code was: if (nth !== 1) selector += ...
          // This implies "first child" doesn't get :nth-of-type(1).
          // This might be ambiguous if we want the first one specifically among many.
          // But `path.join(' > ')` means `parent > child`.
          // If parent has multiple `child` elements, `parent > child` selects ALL of them.
          // So if we want a specific one (e.g. the first link vs second link), we MUST usage :nth-of-type(1).
          // So the previous logic was potentially buggy for the first element if siblings exist!
          if (hasSibling) {
            selector += `:nth-of-type(${nth})`;
          }
        }
      }
    }

    path.unshift(selector);
    currentEl = currentEl.parentNode as HTMLElement;

    // Stop if we hit the container boundaries (if stopAtElement not matched yet)
    if (!currentEl || currentEl === body || currentEl === html) break;
  }

  return path.join(' > ');
};
