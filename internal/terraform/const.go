package terraform

import (
	"os"
)

const (
	maxConcurrentJobs = 8
	tfExt             = ".tf"
	RootModuleDir     = "."
	ModulesDir        = ".terraform" + string(os.PathSeparator) + "modules"
	ModulesJsonPath   = ModulesDir + string(os.PathSeparator) + "modules.json"
)
