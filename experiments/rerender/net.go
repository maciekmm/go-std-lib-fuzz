package rerender

import (
	"strings"

	"golang.org/x/net/html"
)

func NetParseRender(in string) (string, error) {
	parsed, err := html.Parse(strings.NewReader(in))
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	if err := html.Render(&builder, parsed); err != nil {
		return "", err
	}
	return builder.String(), nil
}
