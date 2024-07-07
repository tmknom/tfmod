package dir

import (
	"testing"
)

func TestGraph_AddAndListDst(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "src/foo",
			expected: []string{"dst/foo", "dst/bar", "dst/baz"},
		},
		{
			input:    "src/baz",
			expected: []string{"dst/foo"},
		},
		{
			input:    "not/found",
			expected: []string{},
		},
	}

	sut := NewGraph[FakeSrcDir, FakeDstDir]()
	sut.Add(NewFakeSrcDir("src/foo"), NewFakeDstDir("dst/foo"))
	sut.Add(NewFakeSrcDir("src/foo"), NewFakeDstDir("dst/bar"))
	sut.Add(NewFakeSrcDir("src/foo"), NewFakeDstDir("dst/baz"))
	sut.Add(NewFakeSrcDir("src/bar"), NewFakeDstDir("dst/foo"))
	sut.Add(NewFakeSrcDir("src/bar"), NewFakeDstDir("dst/bar"))
	sut.Add(NewFakeSrcDir("src/baz"), NewFakeDstDir("dst/foo"))

	for _, tc := range cases {
		actual := sut.ListDst(tc.input)

		if len(tc.expected) != len(actual) {
			t.Fatalf("error length:\n input: %v, expected: %v, actual: %v\n graph: %v", tc.input, tc.expected, actual, sut)
		}
		for i, item := range tc.expected {
			if item != actual[i].Rel() {
				t.Errorf("error index:\n %d, input: %v, expected: %v, actual: %v\n graph: %v", i, tc.input, tc.expected, actual, sut)
			}
		}
	}
}

func TestGraph_Include(t *testing.T) {
	cases := []struct {
		input    string
		expected bool
	}{
		{
			input:    "src/bar",
			expected: true,
		},
		{
			input:    "not/found",
			expected: false,
		},
	}

	sut := NewGraph[FakeSrcDir, FakeDstDir]()
	sut.Add(NewFakeSrcDir("src/foo"), NewFakeDstDir("dst/foo"))
	sut.Add(NewFakeSrcDir("src/bar"), NewFakeDstDir("dst/foo"))
	sut.Add(NewFakeSrcDir("src/baz"), NewFakeDstDir("dst/foo"))

	for _, tc := range cases {
		actual := sut.Include(tc.input)
		if actual != tc.expected {
			t.Errorf("input: %v, expected: %v, actual: %v\n graph: %v", tc.input, tc.expected, actual, sut)
		}
	}
}

type FakeSrcDir struct {
	*Dir
}

func NewFakeSrcDir(raw string) *FakeSrcDir {
	return &FakeSrcDir{
		Dir: NewDir(raw, NewBaseDir(".")),
	}
}

type FakeDstDir struct {
	*Dir
}

func NewFakeDstDir(raw string) *FakeDstDir {
	return &FakeDstDir{
		Dir: NewDir(raw, NewBaseDir(".")),
	}
}
