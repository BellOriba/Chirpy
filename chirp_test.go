package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCheckProfane(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"All lower single slur":                           {input: "you are such a kerfuffle", want: "you are such a ****"},
		"Title case single slur":                          {input: "LMAO, get reckt fornax", want: "LMAO, get reckt ****"},
		"Start of string Upper case":                      {input: "SHARBERT SHOULD NOT HAVE RIGHTS", want: "**** SHOULD NOT HAVE RIGHTS"},
		"Mixed case middle of string":                     {input: "You people are all forNaX do not speak with me", want: "You people are all **** do not speak with me"},
		"Multiple slurs, start, middle and end of string": {input: "kerfuffle country, should all sharbert themselves and your culture is all fornax", want: "**** country, should all **** themselves and your culture is all ****"},
		"Single slur":                                     {input: "fornax", want: "****"},
		"Only slurs":                                      {input: "SHARBERT kerfuffle Fornax", want: "**** **** ****"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := checkProfane(tc.input)
			diff := cmp.Diff(tc.want, got)
			if diff != "" {
				t.Fatal(diff)
			}
		})
	}
}
