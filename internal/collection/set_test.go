package collection

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTreeSet_AddAndSlice(t *testing.T) {
	cases := []struct {
		items    []string
		expected []string
	}{
		{
			items:    []string{"foo", "bar", "baz"},
			expected: []string{"bar", "baz", "foo"},
		},
		{
			items:    []string{"foo", "bar", "foo"},
			expected: []string{"bar", "foo"},
		},
	}

	for _, tc := range cases {
		set := NewTreeSet()
		for _, item := range tc.items {
			set.Add(item)
		}

		actual := set.Slice()
		if diff := cmp.Diff(tc.expected, actual); diff != "" {
			t.Errorf("expected: %v, actual: %v", tc.expected, actual)
		}
	}
}
