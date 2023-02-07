APP_NAME := enforce-tool-versions

run:
	cd src && go mod download && go build -o $(APP_NAME) && ./$(APP_NAME)