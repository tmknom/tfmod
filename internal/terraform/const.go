package terraform

import (
	"os"
)

const (
	maxConcurrentJobs = 10
	tfExt             = ".tf"
	RootModuleDir     = "."
	ModulesDir        = ".terraform" + string(os.PathSeparator) + "modules"
	ModulesJsonPath   = ModulesDir + string(os.PathSeparator) + "modules.json"
)
