package main

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

func Sanitize(content io.Reader) string {
	var builder strings.Builder
	tok := html.NewTokenizer(content)
	for {
		tt := tok.Next()
		token := tok.Token()
		switch tt {
		case html.StartTagToken, html.SelfClosingTagToken, html.EndTagToken:
			name := token.Data
			token.Attr = nil
			if name != "strong" {
				continue
			}
			builder.WriteString(token.String())
		case html.ErrorToken:
			return builder.String()
		case html.TextToken:
			builder.WriteString(token.String())
		default:
			continue
		}
	}
}
