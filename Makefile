build:
	GOOS=linux GOARCH=amd64 go build -o handler main.go

deploy: build
	build-lambda-zip --output handlerFunc.zip handler

fmt:
	go fmt ./...

.PHONY: build deploy fmt
