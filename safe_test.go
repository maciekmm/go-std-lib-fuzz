package main

import (
	"strings"
	"testing"
)

var safeContentTests = map[string]bool{
	"<strong>test</strong>":            true,
	"<script>test</script>":            false,
	"<script/>test":                    false,
	"<b>test</b>":                      false,
	`<strong onclick="">test</strong>`: false,
}

func TestIsSafeTokenizer(t *testing.T) {
	for payload, safe := range safeContentTests {
		payload := payload
		safe := safe
		t.Run(payload, func(t *testing.T) {
			if IsSafeTokenizer(strings.NewReader(payload)) != safe {
				t.Fatalf("expected %v for %s", safe, payload)
			}
		})
	}
}

func TestIsSafeParser(t *testing.T) {
	for payload, safe := range safeContentTests {
		payload := payload
		safe := safe
		t.Run(payload, func(t *testing.T) {
			if IsSafeParser(strings.NewReader(payload)) != safe {
				t.Fatalf("expected %v for %s", safe, payload)
			}
		})
	}
}
