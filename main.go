package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

type Request struct {
	User     string `json:"user"`
	Location string `json:"location"`
}

type Response struct {
	StatusCode int    `json:"statusCode"`
	Msg        string `json:"msg"`
}

func HandleRequest(ctx context.Context, req Request) (Response, error) {
	log.Printf("User: %s \n Location: %s", req.User, req.Location)

	res := Response{
		StatusCode: 200,
		Msg:        fmt.Sprintf("User %s location %s is stored", req.User, req.Location),
	}

	return res, nil
}

func main() {
	lambda.Start(HandleRequest)
}
