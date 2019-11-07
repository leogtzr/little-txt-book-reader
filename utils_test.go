package main

import (
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
		got := getNumberLineGoto(tc.line)

		if got != tc.want {
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
		got := percent(tc.currentIndex, tc.totalLines)
		if got != tc.want {
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
		got := countWithoutWhitespaces(tc.strs)
		if got != tc.want {
			t.Errorf("got=%d, expected=%d", got, tc.want)
		}
	}
}
