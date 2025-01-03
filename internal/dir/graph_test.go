package dir

import (
	"testing"

	"github.com/tmknom/tfmod/internal/testlib"
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
		actual := sut.ListDst(NewFakeSrcDir(tc.input))

		if len(tc.expected) != len(actual) {
			t.Fatalf("error length: %s", testlib.Format(sut, tc.expected, actual, tc.input))
		}
		for i, item := range tc.expected {
			if item != actual[i].Rel() {
				t.Errorf("error index: %d %s", i, testlib.Format(sut, tc.expected, actual, tc.input))
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
		actual := sut.Include(createNewDir(tc.input))
		if actual != tc.expected {
			t.Errorf(testlib.Format(sut, tc.expected, actual, tc.input))
		}
	}
}

type FakeSrcDir struct {
	*Dir
}

func NewFakeSrcDir(raw string) *FakeSrcDir {
	return &FakeSrcDir{
		Dir: createNewDir(raw),
	}
}

type FakeDstDir struct {
	*Dir
}

func NewFakeDstDir(raw string) *FakeDstDir {
	return &FakeDstDir{
		Dir: createNewDir(raw),
	}
}

func createNewDir(raw string) *Dir {
	return NewDir(raw, NewBaseDir("."))
}
