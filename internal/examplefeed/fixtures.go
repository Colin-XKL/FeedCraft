package examplefeed

const htmlElementsFixture = `
<article>
  <header>
    <h1>Heading level 1</h1>
    <h2>Heading level 2</h2>
    <p><strong>Strong text</strong>, <em>emphasis</em>, <mark>highlight</mark>, <small>small print</small>, and <time datetime="2026-05-16T08:00:00Z">machine-readable time</time>.</p>
  </header>
  <section>
    <h3>Lists</h3>
    <ul>
      <li>Unordered item</li>
      <li>Nested list <ol><li>First ordered item</li><li>Second ordered item</li></ol></li>
    </ul>
  </section>
  <section>
    <h3>Table</h3>
    <table>
      <caption>Reader support matrix</caption>
      <thead><tr><th>Feature</th><th>Expected display</th></tr></thead>
      <tbody><tr><td>Table cells</td><td>Aligned rows and columns</td></tr></tbody>
    </table>
  </section>
  <blockquote cite="https://example.com">Blockquote content with a cite attribute.</blockquote>
  <pre><code>const reader = "rss";</code></pre>
  <details>
    <summary>Expandable summary</summary>
    <p>This paragraph should appear when details content is expanded.</p>
  </details>
  <figure>
    <div role="img" aria-label="ASCII art">[ FeedCraft HTML fixture ]</div>
    <figcaption>Figure caption support.</figcaption>
  </figure>
  <p>Window UUID: <code>{{WINDOW_UUID}}</code></p>
</article>`

const htmlStylingFixture = `
<article>
  <h1 style="color: #0f766e; background: #ccfbf1; padding: 12px; border-radius: 8px;">Inline color, background, padding, and border radius</h1>
  <p style="font-size: 18px; line-height: 1.7; letter-spacing: 0.04em;">Typography sample with font-size, line-height, and letter spacing.</p>
  <p style="border: 2px solid #2563eb; border-left: 8px solid #1d4ed8; padding: 10px; margin: 12px 0;">Border, margin, and padding sample.</p>
  <div style="display: flex; gap: 8px; align-items: center; flex-wrap: wrap;">
    <span style="background: #fee2e2; color: #991b1b; padding: 4px 8px;">flex item A</span>
    <span style="background: #dbeafe; color: #1e3a8a; padding: 4px 8px;">flex item B</span>
    <span style="background: #dcfce7; color: #14532d; padding: 4px 8px;">flex item C</span>
  </div>
  <div style="display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 8px; margin-top: 12px;">
    <div style="background: #fef3c7; padding: 8px;">grid cell 1</div>
    <div style="background: #ede9fe; padding: 8px;">grid cell 2</div>
  </div>
  <p style="text-align: center; text-decoration: underline; color: rgb(79, 70, 229);">Centered underlined RGB color sample.</p>
  <p>Window UUID: <code>{{WINDOW_UUID}}</code></p>
</article>`

const mediaPictureFixture = `
<article>
  <h1>Picture and source fixture</h1>
  <picture>
    <source media="(min-width: 900px)" srcset="{{BASE_URL}}/example-rss-feeds/assets/picture-wide.svg">
    <source media="(min-width: 480px)" srcset="{{BASE_URL}}/example-rss-feeds/assets/picture-medium.svg">
    <img src="{{BASE_URL}}/example-rss-feeds/assets/picture-fallback.svg" width="640" height="320" alt="FeedCraft picture fallback fixture">
  </picture>
  <figure>
    <img src="{{BASE_URL}}/example-rss-feeds/assets/picture-fallback.svg" srcset="{{BASE_URL}}/example-rss-feeds/assets/picture-medium.svg 480w, {{BASE_URL}}/example-rss-feeds/assets/picture-wide.svg 900w" sizes="(min-width: 900px) 900px, 100vw" alt="FeedCraft srcset fixture">
    <figcaption>Picture, source, img fallback, srcset, sizes, width, height, alt, and caption support.</figcaption>
  </figure>
  <p>Window UUID: <code>{{WINDOW_UUID}}</code></p>
</article>`

var (
	htmlElementsSection = contentSection{
		key:         "html-elements",
		title:       "HTML5 elements fixture",
		description: "Common semantic HTML5 elements embedded in RSS content.",
		html:        htmlElementsFixture,
	}
	htmlStylingSection = contentSection{
		key:         "html-styling",
		title:       "Inline CSS fixture",
		description: "Common inline CSS declarations embedded in RSS content.",
		html:        htmlStylingFixture,
	}
	mediaPictureSection = contentSection{
		key:         "media-picture",
		title:       "Picture source fixture",
		description: "Picture, source, srcset, sizes, and fallback image support.",
		html:        mediaPictureFixture,
	}
)
