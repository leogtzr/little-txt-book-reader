package main

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"testing"
	"textreader/internal/file"
	"textreader/internal/progress"
	"textreader/internal/text"
	"textreader/internal/utils"
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
		if got := progress.GetNumberLineGoto(tc.line); got != tc.want {
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
		if got := progress.Percent(tc.currentIndex, tc.totalLines); got != tc.want {
			t.Errorf("expected: %f, got: %f", tc.want, got)
		}
	}
}

func Test_linesToChangePercentagePoint(t *testing.T) {
	currentLine := 100
	totalLines := 1000
	expectedLinesToChangePercentagePoint := 10

	nextPercentagePoint := progress.LinesToChangePercentagePoint(currentLine, totalLines)

	if nextPercentagePoint != expectedLinesToChangePercentagePoint {
		t.Errorf("expected: %d, got: %d", expectedLinesToChangePercentagePoint, nextPercentagePoint)
	}
}

func TestGetFileToSaveName(t *testing.T) {
	name := "/getHomeDirectoryPath/leo/code/little-txt-book-reader/lala.xts"
	if baseFileName := filepath.Base(name); baseFileName != "lala.xts" {
		t.Errorf("got=[%s], want=[%s]", baseFileName, "lala.xts")
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
		if got := text.GetChunk(&tt.content, tt.from, tt.to); !listsAreEqual(got, tt.want) {
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
		text.FindAndRemove(&strs, tt.toRemove)
		if !listsAreEqual(strs, tt.want) {
			t.Errorf("got=[%s], want=[%s]", strs, tt.want)
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
		if got := listsAreEqual(tt.a, tt.b); got != tt.shouldBeEqual {
			t.Errorf("got=[%t], want=[%t]", got, tt.shouldBeEqual)
		}
	}
}

func Test_removeTrailingSpaces(t *testing.T) {
	type test struct {
		s, want string
	}

	tests := []test{
		test{
			s: `abc 
holis 
`, want: `abc
holis`,
		},
		test{
			s: `mayor de un señor ordenado que administra con prudencia su trapillo. Algunas                                                                                                                                                                 
			anotaciones sobre la compra de libros a anticuarios parisinos. Ahora lo veía todo      `,
			want: `mayor de un señor ordenado que administra con prudencia su trapillo. Algunas
anotaciones sobre la compra de libros a anticuarios parisinos. Ahora lo veía todo`,
		},
	}

	for _, tt := range tests {
		if got := text.RemoveTrailingSpaces(tt.s); got != tt.want {
			t.Errorf("got=[%s], want=[%s]", strings.ReplaceAll(got, "\n", "@"), strings.ReplaceAll(tt.want, "\n", "@"))
		}
	}

}

func Test_readLines(t *testing.T) {
	type test struct {
		file io.Reader
		want []string
	}

	tests := []test{
		test{
			file: strings.NewReader(`a
b`), want: []string{"a", "b"},
		},
	}

	for _, tt := range tests {
		if got, _ := file.ReadLines(tt.file); !listsAreEqual(got, tt.want) {
			t.Errorf("got=[%s], want=[%s]", got, tt.want)
		}
	}
}

func Test_check(t *testing.T) {
	type test struct {
		err         error
		shouldPanic bool
	}

	tests := []test{
		test{
			err:         fmt.Errorf("bye"),
			shouldPanic: true,
		},

		test{
			err:         nil,
			shouldPanic: false,
		},
	}

	for _, tt := range tests {
		defer func() {
			if r := recover(); r == nil && !tt.shouldPanic {
				t.Errorf("The code did not panic")
			}
		}()
		utils.Check(tt.err)
	}
}

func Test_paginate(t *testing.T) {
	type test struct {
		elements []string
		skip     int
		size     int
		want     []string
	}

	elements := []string{
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n",
		"o", "p", "q", "r", "s", "t", "u", "v", "w",
	}

	tests := []test{
		test{
			elements: elements,
			skip:     3,
			size:     3,
			want: []string{
				"d", "e", "f",
			},
		},

		test{
			elements: elements,
			skip:     len(elements) - 3,
			size:     3,
			want: []string{
				"u", "v", "w",
			},
		},
	}

	for _, tt := range tests {
		got := utils.Paginate(tt.elements, tt.skip, tt.size)
		if !listsAreEqual(got, tt.want) {
			t.Errorf("got=[%s], want=[%s]", got, tt.want)
		}
	}

}

func Test_removeWhiteSpaces(t *testing.T) {
	type test struct {
		input string
		want  string
	}

	tests := []test{
		test{
			input: `He considerado que Dios, siendo                                                                   improbable`,
			want:  `He considerado que Dios, siendo improbable`,
		},
		test{
			input: `  He considerado que Dios, siendo                                                                   improbable`,
			want:  ` He considerado que Dios, siendo improbable`,
		},
	}

	for _, tc := range tests {
		got := text.RemoveWhiteSpaces(tc.input)
		if got != tc.want {
			t.Errorf("got=[%s], want=[%s]", got, tc.want)
		}
	}

}

func Test_getPercentage(t *testing.T) {
	type test struct {
		position    int
		fileContent []string
		want        float64
	}

	tests := []test{
		test{
			position:    4,
			fileContent: []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"},
			want:        40.0,
		},
		test{
			position:    9,
			fileContent: []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o"},
			want:        60.0,
		},
	}

	for _, tc := range tests {
		got := progress.GetPercentage(tc.position, &tc.fileContent)
		if got != tc.want {
			t.Errorf("got=[%f], want=[%f]", got, tc.want)
		}
	}

}

func listsAreEqual(a, b []string) bool {
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
