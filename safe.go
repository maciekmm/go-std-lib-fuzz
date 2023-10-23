package main

import (
	"io"

	"golang.org/x/net/html"
)

func IsSafeTokenizer(content io.Reader) bool {
	tok := html.NewTokenizer(content)
	for {
		tt := tok.Next()
		switch tt {
		case html.StartTagToken:
			name, hasAttr := tok.TagName()
			if hasAttr || string(name) != "strong" {
				return false
			}
		case html.ErrorToken:
			if tok.Err() == io.EOF {
				return true
			}
			return false
		case html.TextToken, html.EndTagToken:
		default:
			return false
		}
	}
}

func IsSafeParser(content io.Reader) bool {
	parsed, err := html.ParseFragment(content, nil)
	if err != nil {
		return false
	}
	for _, el := range parsed {
		if !isNodeSafe(el) {
			return false
		}
	}
	return true
}

func isNodeSafe(node *html.Node) bool {
	if node == nil {
		return true
	}
	if len(node.Attr) != 0 {
		return false
	}
	if node.Type == html.ElementNode {
		if node.Data != "strong" && node.Data != "html" && node.Data != "head" && node.Data != "body" {
			return false
		}
	}
	return isNodeSafe(node.NextSibling) && isNodeSafe(node.FirstChild)
}
