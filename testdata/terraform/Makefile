# Include: minimum
-include .makefiles/minimum/Makefile
.makefiles/minimum/Makefile:
	@git clone https://github.com/tmknom/makefiles.git .makefiles >/dev/null 2>&1

.PHONY: get
get: ## Get modules
	find . -name '*.tf' | xargs -I {} dirname {} | xargs -I {} terraform -chdir={} get

.PHONY: fmt
fmt: ## Format
	terraform fmt -recursive
