package terraform

import (
	"testing"

	"github.com/tmknom/tfmod/internal/dir"
)

func TestTerraform_GetAll(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cases := []struct {
		inputs []string
	}{
		{
			inputs: []string{"dev", "prd", "stg"},
		},
	}

	baseDir := dir.NewBaseDir("testdata/terraform/env")
	for _, tc := range cases {
		sut := NewCommand()

		err := sut.GetAll(baseDir.ConvertDirs(tc.inputs))
		if err != nil {
			t.Errorf("%v: unexpected error: %e", tc.inputs, err)
		}
	}
}
