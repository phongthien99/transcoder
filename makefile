SHELL=/bin/bash -e -o pipefail
PWD := $(shell pwd)
LINTER_VERSION := v1.55.1

run:
	@go run apps/$(app)/cli/main.go -c apps/$(app)/src/config/default.yaml

gc $(app):
	./scripts/gen_type_config.sh apps/example/src/config
