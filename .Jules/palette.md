## 2025-12-18 - Icon-Only Tooltips Accessibility
**Learning:** Icon-only tooltips (like `icon-question-circle` inside `a-tooltip`) are often implemented as non-focusable SVGs, making them inaccessible to keyboard users who cannot trigger the tooltip or understand the icon's purpose.
**Action:** Always wrap interactive icons in a focusable element (like `<button>` or `<span tabindex="0" role="button">`) with a descriptive `aria-label` to ensure they are discoverable and usable via keyboard navigation.
