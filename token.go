package main

import "golang.org/x/net/html"

type Token struct {
	Name string
	Type html.TokenType
}
