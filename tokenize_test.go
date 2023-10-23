package main

import (
	"reflect"
	"testing"

	"golang.org/x/net/html"
)

var tests = []string{
	"</a/>",
	"<Xmp/><A>",
	"<script><script>text</script/>test",
	"<p>text</p>",
	"<img src=test>",
	"<!-- comment --><p>more</p>",
	"<!DOCTYPE html><p>html</p>",
	"<!--!> <h1 value=\"--><a href=\"javascript:alert(document.domain)\">link",
	"0</",
	"</<a href=\"javascript:alert(document.domain)\"><a link>",
	"</<a href=\"javascript:alert(document.domain)\"><a href=\"javascript:alert(document.domain)\">link",
}

func FuzzTokenize(f *testing.F) {
	for _, test := range tests {
		f.Add(test)
	}

	f.Fuzz(func(t *testing.T, input string) {
		netTokens, err := NetTokenize(input)
		if err != nil {
			t.SkipNow()
		}

		lexborTokens, err := LexborTokenize(input)
		if err != nil {
			t.SkipNow()
		}

		tokenOfInterest := func(token Token) bool {
			return token.Type == html.StartTagToken || token.Type == html.CommentToken || token.Type == html.SelfClosingTagToken
		}

		netStartTags := []Token{}
		for _, token := range netTokens {
			if !tokenOfInterest(token) {
				continue
			}

			// I couldn't get lexbor to properly interpret content inside raw nodes, we're skipping them for now
			switch s := token.Name; s {
			case "iframe", "noembed", "noframes", "noscript", "plaintext", "script", "style", "title", "textarea", "xmp":
				t.SkipNow()
			}
		}

		lexborStartTags := []Token{}
		for _, token := range lexborTokens {
			if !tokenOfInterest(token) {
				continue
			}
			lexborStartTags = append(lexborStartTags, token)
		}

		if !reflect.DeepEqual(netStartTags, lexborStartTags) {
			t.Errorf("Tokenization mismatch:\nlexbor\t=%+v, \nnet\t=%+v\n not equal, input: %s", lexborStartTags, netStartTags, input)
			return
		}
	})
}
