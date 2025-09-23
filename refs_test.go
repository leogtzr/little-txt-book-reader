package main

import (
	"testing"
	"textreader/internal/references"
)

func TestIdentifyReferences(t *testing.T) {

	type test struct {
		line string
		want []string
	}

	tests := []test{
		{
			line: "Hola a Todos Ustedes, Maria, te odio OK bye No hola Hmm Byes",
			want: []string{"Hola", "Todos Ustedes", "Maria", "OK", "No", "Hmm Byes"},
		},
		{
			line: "La casa roja, Maria, Ok. No",
			want: []string{"La", "Maria", "Ok", "No"},
		},
		{
			line: "Aquel Mio, ok, bye. Final.",
			want: []string{"Aquel Mio", "Final"},
		},
		{
			line: "Cervantes, Sterne, Melville, Proust, Musil y Pynchon»",
			want: []string{"Cervantes", "Sterne", "Melville", "Proust", "Musil", "Pynchon»"},
		},
		{
			line: "Cervantes, Sterne, Melville, Proust, Musil y Pynchon Tamal Bye, ok. Si",
			want: []string{"Cervantes", "Sterne", "Melville", "Proust", "Musil", "Pynchon Tamal Bye", "Si"},
		},
		{
			line: "Cervantes, Sterne.",
			want: []string{"Cervantes", "Sterne"},
		},
		{
			line: "Cervantes, ok? Bolaño",
			want: []string{"Cervantes", "Bolaño"},
		},
		{
			line: "Bolaño En Sus Laureles, OK.",
			want: []string{"Bolaño En Sus Laureles", "OK"},
		},
		{
			line: "El secreto del mal, La Universidad Desconocida, Elcmn",
			want: []string{"El", "La Universidad Desconocida", "Elcmn"},
		},
		{
			line: "La Lola a Todos Ustedes, Maria, te odio OK bye No hola Hmm Byes",
			want: []string{"La Lola", "Todos Ustedes", "Maria", "OK", "No", "Hmm Byes"},
		},
		{
			line: "",
			want: []string{},
		},
	}

	for _, tc := range tests {
		got := references.ExtractReferences(tc.line)
		if len(got) != len(tc.want) {
			t.Errorf("got=[%d], want=[%d] References", len(got), len(tc.want))
		}

		for i, testRef := range tc.want {
			if got[i] != testRef {
				t.Errorf("got=[%s], expected=[%s]", got[i], testRef)
			}
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
