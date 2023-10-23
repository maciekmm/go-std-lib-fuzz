package rerender

import (
	"strings"
	"testing"
)

var tests = []string{
	"</a/>",
	"<xml/><a>",
	"<script><script>text</script/>test",
	"<p>text</p>",
	"<img src=test>",
	"<!-- comment --><p>more</p>",
	"<!DO",
	"<!DOCTYPE html><p>html</p>",
	"<!--!> <h1 value=\"--><a href=\"javascript:alert(document.domain)\">link",
	"0</",
	"</<a href=\"javascript:alert(document.domain)\"><a link>",
	"</<a href=\"javascript:alert(document.domain)\"><a href=\"javascript:alert(document.domain)\">link",
	"<xmp>test</xmp>",
	"<xmp><a>wtf</a></xmp>",
	"</xmp><a>wtf</a>",
	// "<script =\">alert(document.domain)</script>",
}

func FuzzParseRender(f *testing.F) {
	for _, test := range tests {
		f.Add(test)
	}

	f.Fuzz(func(t *testing.T, html string) {
		lex, err := LexborParseRender(html)
		if err != nil {
			t.Logf("Lexbor: could not parse %s", html)
			t.Skip()
		}
		net, err := NetParseRender(html)
		if err != nil {
			t.Logf("Net: could not parse %s", html)
			t.Skip()
		}
		// some massaging
		net = strings.ReplaceAll(net, "/>", ">")
		net = strings.ReplaceAll(net, `=""`, "")
		if lex != net {
			t.Fatalf("\nnet: %s\nlex: %s\noriginal: %s", net, lex, html)
		}
	})
}
