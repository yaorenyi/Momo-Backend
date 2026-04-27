package utils

import "regexp"

var (
	reScriptStyle = regexp.MustCompile(`(?i)<(?:script|style)[\s\S]*?</(?:script|style)>`)
	reEventHandler = regexp.MustCompile(`(?i)\s+on\w+\s*=\s*(?:"[^"]*"|'[^']*'|[^\s>]+)`)
	reJSLinkDQuote = regexp.MustCompile(`(?i)(?:href|src|action|formaction)\s*=\s*"(?:javascript|vbscript):[^"]*"`)
	reJSLinkSQuote = regexp.MustCompile(`(?i)(?:href|src|action|formaction)\s*=\s*'(?:javascript|vbscript):[^']*'`)
	reJSProtocol     = regexp.MustCompile(`(?i)^(?:javascript|vbscript):\s*`)
	reDangerousTags  = regexp.MustCompile(`(?i)</?(?:iframe|object|embed|frame|meta|link|base|form|input)\b[^>]*>`)
)

// CheckContent sanitizes raw user input by removing XSS attack vectors
func CheckContent(content string) string {
	content = reScriptStyle.ReplaceAllString(content, "")
	content = reEventHandler.ReplaceAllString(content, "")
	content = reJSLinkDQuote.ReplaceAllString(content, "")
	content = reJSLinkSQuote.ReplaceAllString(content, "")
	content = reJSProtocol.ReplaceAllString(content, "")
	content = reDangerousTags.ReplaceAllString(content, "")
	return content
}
