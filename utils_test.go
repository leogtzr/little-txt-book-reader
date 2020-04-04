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
	line := "1234567890 1234567890 1234567890 1234567890 1234567890 12345678901234567890 12345678901234567890"
	if needsSemiWrap(line) {
		t.Errorf("'%s' should not be wrapped", line)
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
	baseFileName := filepath.Base(name)

	if baseFileName != "refs.go" {
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
		got := getChunk(&tt.content, tt.from, tt.to)
		if !equal(got, tt.want) {
			t.Errorf("got=[%s], want=[%s]", got, tt.want)
		}
	}
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
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
