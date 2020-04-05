package main

import (
	"path/filepath"
	"testing"
)

func Test_getNumberLineGoto(t *testing.T) {

	type test struct {
		line string
		want string
	}

	tests := []test{
		{line: "hola12_mundo_345", want: "12345"},
	}

	for _, tc := range tests {
		if got := getNumberLineGoto(tc.line); got != tc.want {
			t.Errorf("expected: %s, got: %s", tc.want, got)
		}
	}
}

func Test_percent(t *testing.T) {

	type test struct {
		totalLines   int
		currentIndex int
		want         float64
	}

	tests := []test{
		{
			totalLines:   150,
			currentIndex: 30,
			want:         20.0,
		},
	}

	for _, tc := range tests {
		if got := percent(tc.currentIndex, tc.totalLines); got != tc.want {
			t.Errorf("expected: %f, got: %f", tc.want, got)
		}
	}
}

func Test_linesToChangePercentagePoint(t *testing.T) {
	currentLine := 100
	totalLines := 1000
	expectedLinesToChangePercentagePoint := 10

	nextPercentagePoint := linesToChangePercentagePoint(currentLine, totalLines)

	if nextPercentagePoint != expectedLinesToChangePercentagePoint {
		t.Errorf("expected: %d, got: %d", expectedLinesToChangePercentagePoint, nextPercentagePoint)
	}
}

func Test_needsSemiWrap(t *testing.T) {
	type test struct {
		line string
		want bool
	}

	tests := []test{
		{
			line: "1234567890 1234567890 1234567890 1234567890 1234567890 12345678901234567890 12345678901234567890",
			want: false,
		},
		{
			line: "1234567890 1234567890",
			want: false,
		},
	}

	for _, tt := range tests {
		if got := needsSemiWrap(tt.line); got != tt.want {
			t.Errorf("got=[%t], want=[%t]", got, tt.want)
		}
	}
}

func Test_countWithoutWhitespaces(t *testing.T) {
	type test struct {
		strs []string
		want int
	}

	tests := []test{
		{strs: []string{}, want: 0},
		{strs: []string{"Hola", "Mundo", "Cruel"}, want: 14},
	}

	for _, tc := range tests {
		if got := countWithoutWhitespaces(tc.strs); got != tc.want {
			t.Errorf("got=%d, expected=%d", got, tc.want)
		}
	}
}

func TestGetFileToSaveName(t *testing.T) {
	name := "/home/leo/code/little-txt-book-reader/refs.go"
	if baseFileName := filepath.Base(name); baseFileName != "refs.go" {
		t.Errorf("got=[%s], want=[%s]", baseFileName, "refs.go")
	}
}

func TestGetChunk(t *testing.T) {
	content := []string{
		"hola",
		"mundo",
		"cruel",
		"que onda",
		"ok",
		"bye",
	}

	type test struct {
		content  []string
		from, to int
		want     []string
	}

	tests := []test{
		{
			content: content,
			from:    0,
			to:      3,
			want: []string{
				"hola",
				"mundo",
				"cruel",
			},
		},
	}

	for _, tt := range tests {
		if got := getChunk(&tt.content, tt.from, tt.to); !equal(got, tt.want) {
			t.Errorf("got=[%s], want=[%s]", got, tt.want)
		}
	}
}

func Test_findAndRemove(t *testing.T) {
	strs := []string{
		"hola",
		"a",
		"todos",
		"bye",
	}

	type test struct {
		strs     []string
		toRemove string
		want     []string
	}

	tests := []test{
		{
			strs:     strs,
			toRemove: "todos",
			want: []string{
				"hola",
				"a",
				"bye",
			},
		},

		{
			strs:     strs,
			toRemove: "bye",
			want: []string{
				"hola",
				"a",
			},
		},
	}

	for _, tt := range tests {
		findAndRemove(&strs, tt.toRemove)
		if !equal(strs, tt.want) {
			t.Errorf("got=[%s], want=[%s]", strs, tt.want)
		}
	}
}

func Test_remove(t *testing.T) {
	type test struct {
		s     []string
		index int
		want  []string
	}

	tests := []test{
		{
			s: []string{
				"hola",
				"mundo",
				"cruel",
			},
			index: 1,
			want: []string{
				"hola",
				"cruel",
			},
		},
	}

	for _, tt := range tests {
		if got := remove(tt.s, tt.index); !equal(got, tt.want) {
			t.Errorf("got=[%s], want=[%s]", got, tt.want)
		}
	}
}

func Test_longestLineLength(t *testing.T) {
	type test struct {
		text string
		len  int
	}

	tests := []test{
		{
			text: `
La navaja
se fue
a su cuna
`,
			len: len("La navaja"),
		},

		{
			text: "",
			len:  0,
		},
	}

	for _, tt := range tests {
		if got := longestLineLength(tt.text); got != tt.len {
			t.Errorf("got=[%d], want=[%d]", got, tt.len)
		}
	}
}

func Test_equal(t *testing.T) {
	type test struct {
		a, b          []string
		shouldBeEqual bool
	}

	tests := []test{
		{
			a:             []string{},
			b:             []string{},
			shouldBeEqual: true,
		},

		{
			a:             []string{"hola", "ok"},
			b:             []string{"no"},
			shouldBeEqual: false,
		},

		{
			a:             []string{"a", "b"},
			b:             []string{"no", "hmm"},
			shouldBeEqual: false,
		},
	}

	for _, tt := range tests {
		got := equal(tt.a, tt.b)
		if got != tt.shouldBeEqual {
			t.Errorf("got=[%t], want=[%t]", got, tt.shouldBeEqual)
		}
	}
}
