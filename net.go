package main

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

type TokenizeFunc func(input string) ([]Token, error)

func NetTokenize(data string) ([]Token, error) {
	tokenizer := html.NewTokenizer(strings.NewReader(data))
	tokenizer.SetMaxBuf(100)
	tokens := []Token{}

	for {
		tt := tokenizer.Next()
		if tt == html.ErrorToken && tokenizer.Err() != io.EOF {
			return nil, tokenizer.Err()
		}

		token := tokenizer.Token()
		name := ""
		terminate := false

		switch tt {
		case html.StartTagToken, html.SelfClosingTagToken:
			name = token.Data
			tt = html.StartTagToken
		case html.EndTagToken:
			name = token.Data
		case html.TextToken:
			name = "#text"
		case html.CommentToken:
			name = "#text"
			tt = html.TextToken
		case html.DoctypeToken:
			name = "DOCTYPE"
		case html.ErrorToken:
			name = ""
			terminate = true
		}

		tokens = append(tokens, Token{
			Name: name,
			Type: tt,
		})
		if terminate {
			break
		}
	}

	return tokens, nil
}
