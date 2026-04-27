package utils

import (
	"regexp"

	"github.com/gomarkdown/markdown"
)

var (
	reScript = regexp.MustCompile(`(?i)<script[\s\S]*?</script>`)
	reEvent  = regexp.MustCompile(`(?i)\s+on\w+\s*=\s*(?:"[^"]*"|'[^']*'|[^\s>]+)`)
	reJSLink = regexp.MustCompile(`(?i)href\s*=\s*"(?:javascript|vbscript):`)
)

func ParseMarkdown(content string) string {
	if content == "" {
		return ""
	}
	// Parse markdown to HTML
	result := markdown.ToHTML([]byte(content), nil, nil)
	html := string(result)

	// Sanitize: remove script tags
	html = reScript.ReplaceAllString(html, "")
	// Remove event handlers
	html = reEvent.ReplaceAllString(html, "")
	// Remove javascript: links
	html = reJSLink.ReplaceAllString(html, `href="#"`)

	return html
}
