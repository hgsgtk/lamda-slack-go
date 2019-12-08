build:
	GOOS=linux GOARCH=amd64 go build -o handler ./handler

deploy:
	build-lambda-zip --output handlerFunc.zip handler/handler

.PHONY: build
