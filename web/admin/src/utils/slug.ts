import limax from 'limax';

/**
 * Generates a recipe ID from a given text string.
 * The generated ID strictly contains only lowercase letters (a-z), numbers (0-9), and hyphens (-).
 * It uses limax for intelligent transliteration (e.g. Chinese Pinyin) before cleaning.
 *
 * @param text The source text (e.g. title)
 * @returns A safe, slugified string
 */
export default function generateRecipeId(text: string): string {
  if (!text) return '';

  // Use limax to generate a base slug.
  // tone: false ensures Pinyin doesn't have tone numbers (e.g. 'ni3-hao3' -> 'ni-hao')
  const baseSlug = limax(text, {
    tone: false,
    separateNumbers: false,
    maintainCase: false, // force lowercase
    separator: '-',
  });

  // Strict enforcement of [a-z0-9-]
  return baseSlug
    .toLowerCase()
    .replace(/[^a-z0-9-]/g, '-') // Replace any remaining invalid chars with hyphen
    .replace(/-+/g, '-') // Collapse multiple hyphens
    .replace(/^-+|-+$/g, ''); // Trim leading/trailing hyphens
}
