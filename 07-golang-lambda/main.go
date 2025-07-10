package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

func HandleLambdaEvent(event MyEvent) (MyResponse, error) {
	return MyResponse{Message: fmt.Sprintf("%s is %d years old", event.Name, event.Age)}, nil
}

type MyResponse struct {
	Message string `json:"Answer:"`
}

type MyEvent struct {
	Name string `json:"what is your name?"`
	Age  int    `json:"How old are you"`
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
