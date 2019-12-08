package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/hgsgtk/lamda-slack-go/lambdahandler"
)

func main() {
	lambda.Start(lambdahandler.Handler)
}
