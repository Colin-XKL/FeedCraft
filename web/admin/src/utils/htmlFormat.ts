/**
 * Format HTML string with proper indentation.
 * Uses a simple token-based approach for RSS content display.
 */
const INDENT = '  ';

const VOID_TAGS = new Set([
  'area',
  'base',
  'br',
  'col',
  'embed',
  'hr',
  'img',
  'input',
  'link',
  'meta',
  'param',
  'source',
  'track',
  'wbr',
]);

const INLINE_TAGS = new Set([
  'a',
  'abbr',
  'acronym',
  'b',
  'bdo',
  'big',
  'br',
  'button',
  'cite',
  'code',
  'dfn',
  'em',
  'i',
  'img',
  'input',
  'kbd',
  'label',
  'map',
  'object',
  'output',
  'q',
  's',
  'samp',
  'select',
  'small',
  'span',
  'strong',
  'sub',
  'sup',
  'textarea',
  'time',
  'tt',
  'u',
  'var',
]);

export default function formatHTML(html: string): string {
  if (!html || !html.trim()) return '';

  const result: string[] = [];
  let depth = 0;

  const tokens = html.split(/(<[^>]+>)/g);

  tokens.forEach((token) => {
    const trimmed = token.trim();
    if (!trimmed) return;

    if (trimmed.startsWith('</')) {
      const tagName = trimmed
        .slice(2, trimmed.length - 1)
        .split(/\s/)[0]
        .toLowerCase();
      if (!INLINE_TAGS.has(tagName)) {
        depth = Math.max(0, depth - 1);
      }
      result.push(INDENT.repeat(depth) + trimmed);
    } else if (trimmed.startsWith('<!--')) {
      result.push(INDENT.repeat(depth) + trimmed);
    } else if (trimmed.startsWith('<')) {
      const tagMatch = trimmed.match(/^<([a-zA-Z][a-zA-Z0-9-]*)/);
      const tagName = tagMatch ? tagMatch[1].toLowerCase() : '';
      const selfClosing = trimmed.endsWith('/>') || VOID_TAGS.has(tagName);

      result.push(INDENT.repeat(depth) + trimmed);

      if (!selfClosing && tagName && !INLINE_TAGS.has(tagName)) {
        depth += 1;
      }
    } else {
      const text = trimmed;
      if (text) {
        result.push(INDENT.repeat(depth) + text);
      }
    }
  });

  return result.join('\n');
}
