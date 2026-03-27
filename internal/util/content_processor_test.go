package util

import (
	"strings"
	"testing"
)

func TestProcessContent(t *testing.T) {
	htmlContent := `
<html>
<body>
<h1>Title</h1>
<p>Hello <a href="http://example.com">link</a></p>
<img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNkYAAAAAYAAjCB0C8AAAAASUVORK5CYII=" alt="image" />
<img src="http://example.com/image.jpg" alt="image" />
<script>alert("hello");</script>
<style>.css { color: red; }</style>
<div aria-label="hello" data-custom="value">world</div>
<custom-tag>custom</custom-tag>

<!-- Some other image tags -->
<img
  src="data:image/jpeg;base64,/9j/4AAQSkZJRgABAQEASABIAAD/2wBD..."
  alt="Base64 Image"
/>
<IMG SRC="data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7">

<a href="data:image/jpeg;base64,/9j/4AAQSkZJRgABAQEASABIAAD/2wBD...">base64 link</a>

<svg width="100" height="100">
  <circle cx="50" cy="50" r="40" stroke="green" stroke-width="4" fill="yellow" />
</svg>
</body>
</html>`

	option := ContentProcessOption{
		RemoveLinks: true,
		RemoveImage: true,
		ConvertToMd: true,
	}

	result := ProcessContent(htmlContent, option)
	t.Logf("Result:\n%s", result)

	if strings.Contains(result, "alert(\"hello\")") {
		t.Errorf("script tag not removed")
	}
	if strings.Contains(result, ".css") {
		t.Errorf("style tag not removed")
	}
	if strings.Contains(result, "base64") {
		t.Errorf("base64 not removed")
	}
}
