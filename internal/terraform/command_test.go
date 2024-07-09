package terraform

import (
	"testing"

	"github.com/tmknom/tfmod/internal/testlib"

	"github.com/tmknom/tfmod/internal/dir"
)

func TestTerraform_GetAll(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cases := []struct {
		input []string
	}{
		{
			input: []string{"dev", "prd", "stg"},
		},
	}

	baseDir := dir.NewBaseDir("testdata/terraform/env")
	for _, tc := range cases {
		sut := NewCommand()

		err := sut.GetAll(baseDir.ConvertDirs(tc.input))
		if err != nil {
			t.Errorf(testlib.FormatError(err, sut, tc.input))
		}
	}
}
