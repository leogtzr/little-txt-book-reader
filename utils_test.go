package main

import (
	"testing"
)

func Test_getNumberLineGoto(t *testing.T) {
	line := "hola12_mundo_345"
	lineWithOnlyNumbers := getNumberLineGoto(line)
	expectedText := "12345"
	if lineWithOnlyNumbers != expectedText {
		t.Errorf("expected: %s, got: %s", expectedText, lineWithOnlyNumbers)
	}
}

func Test_percent(t *testing.T) {
	totalLines := 150
	i := 30

	progress := percent(i, totalLines)
	expectedPercentageProgress := 20.0
	if progress != expectedPercentageProgress {
		t.Errorf("expected: %f, got: %f", expectedPercentageProgress, progress)
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
