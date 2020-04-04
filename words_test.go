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
		if got := extractWords(tc.line); !reflect.DeepEqual(got, tc.want) {
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
		if got := sanitizeWord(tc.line); got != tc.want {
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
		if got := isTitle(tc.word); got != tc.want {
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
		if got := shouldIgnoreWord(tc.word); got != tc.want {
			t.Errorf("want=[%t], got=[%t]", tc.want, got)
		}
	}
}

func Test_sanitizeFileName(t *testing.T) {
	type test struct {
		fileName, want string
	}

	tests := []test{
		{
			fileName: "Hola mundo.txt",
			want:     "Holamundo.txt",
		},
	}

	for _, tt := range tests {
		if got := sanitizeFileName(tt.fileName); got != tt.want {
			t.Errorf("got=[%s], want=[%s]", got, tt.want)
		}
	}
}
