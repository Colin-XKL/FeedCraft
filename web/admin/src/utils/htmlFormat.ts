import { html as beautifyHtml } from 'js-beautify';

export default function formatHTML(input: string): string {
  return beautifyHtml(input, { indent_size: 2, wrap_line_length: 120 });
}
