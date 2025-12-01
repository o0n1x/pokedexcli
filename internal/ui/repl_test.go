package ui

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCleanInput(t *testing.T) {
	cases := map[string]struct {
		input    string
		expected []string
	}{
		"simple":        {"  hello  world  ", []string{"hello", "world"}},
		"nosplit":       {"interesting", []string{"interesting"}},
		"lowersimple":   {"  WELL WELL  ", []string{"well", "well"}},
		"lower_nosplit": {"WELlHeLlO", []string{"wellhello"}},
		"simple_edge":   {"  edge", []string{"edge"}},
		"split edge":    {"  a nice edge", []string{"a", "nice", "edge"}},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := CleanInput(tc.input)
			diff := cmp.Diff(tc.expected, got)
			if diff != "" {
				t.Fatal(diff)
			}
		})
	}

}
