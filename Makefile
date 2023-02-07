APP_NAME := version-enforcer

# depends on *.go, *.mod, *.sum
APP_FILES := $(shell find . -type f -name '*.go' -name '*.mod' -name '*.sum')

build: $(APP_FILES)
	go mod download && go build -o $(APP_NAME)

.PHONY: run
run: $(APP_FILES)
	go mod download && go build -o $(APP_NAME) && ./$(APP_NAME) --config version-enforcer.hcl

.PHONY: test
test: $(APP_FILES)
	go test -v ./...
	go test -v ./identifier -fuzz=FuzzDoesSemverMatch -fuzztime=10s -fuzzminimizetime=10s -parallel=8

.PHONY: clean
clean:
	go clean && rm -f $(APP_NAME)