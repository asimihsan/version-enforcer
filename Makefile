APP_NAME := enforce-tool-versions

# depends on *.go, *.mod, *.sum
APP_FILES := $(shell find . -type f -name '*.go' -name '*.mod' -name '*.sum')

.PHONY: run
run: $(APP_FILES)
	cd src && go mod download && go build -o $(APP_NAME) && ./$(APP_NAME) --config tool-enforcer.hcl

.PHONY: test
test: $(APP_FILES)
	cd src && go test -v ./...