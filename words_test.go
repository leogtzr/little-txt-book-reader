package main

import (
	"reflect"
	"testing"
)

func TestExtractWords(t *testing.T) {
	type test struct {
		line string
		want []string
	}

	tests := []test{
		{line: "anita lava la tina", want: []string{"anita", "lava", "la", "tina"}},
	}

	for _, tc := range tests {
		got := extractWords(tc.line)
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("got=[%s], want=[%s]", got, tc.want)
		}
	}

}

func Test_SanitizeWord(t *testing.T) {
	type test struct {
		line string
		want string
	}

	tests := []test{
		{line: "el mejor, ok.", want: "el mejor ok"},
		{line: "\"hola\" pendejos", want: "hola pendejos"},
	}

	for _, tc := range tests {
		got := sanitizeWord(tc.line)
		if got != tc.want {
			t.Errorf("got=[%s], want=[%s]", got, tc.want)
		}
	}
}

func TestIsTitle(t *testing.T) {
	type test struct {
		word string
		want bool
	}

	tests := []test{
		{word: "ok", want: false},
		{word: "Leo", want: true},
		{word: "HOLA", want: true},
	}

	for _, tc := range tests {
		got := isTitle(tc.word)
		if got != tc.want {
			t.Errorf("got=[%t], want=[%t]", got, tc.want)
		}
	}
}

func TestWordIsInBannedWords(t *testing.T) {
	type test struct {
		word string
		want bool
	}

	tests := []test{
		{word: "su", want: false},
		{word: "Yo", want: true},
	}

	for _, tc := range tests {
		got := shouldIgnoreWord(tc.word)
		if got != tc.want {
			t.Errorf("want=[%t], got=[%t]", tc.want, got)
		}
	}
}
