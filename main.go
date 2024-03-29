package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/a-h/awsapigatewayv2handler"
	"github.com/aws/aws-lambda-go/events"
)

func handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

	var event Event
	err := json.Unmarshal([]byte(request.Body), &event)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{}, err
	}

	msg := fmt.Sprintf("Hello %v %v", *event.User, *event.Location)
	responseBody := ResponseBody{
		Message: &msg,
	}
	jbytes, _ := json.Marshal(responseBody)

	jstr := string(jbytes)

	response := events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Body:       jstr,
	}

	log.Printf(`{%+v}`, response)

	return response, nil
}

type ResponseBody struct {
	Message *string `json:"message"`
}

type Event struct {
	User     *string `json:"user"`
	Location *string `json:"location"`
}

func NewServer() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Write([]byte("Hello"))
	})

	mux.HandleFunc("/locations", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		jsonData := fmt.Sprintf(`{"path": "%s","isEmployed": false}`, "request")
		log.Println(jsonData)

		fmt.Println("Received new POST")

		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		for k, v := range r.Form {
			fmt.Printf("%s=%s\n", k, v)
		}
	})

	return mux
}

func main() {
	http.Handle("/tester/locations", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		io.WriteString(w, "Hello")
	}))
	// Serve the static directory that has been embedded into the binary.
	// The static directory contains a mix of binary and text files, for testing.
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		jsonData := fmt.Sprintf(`{"path": "%+v","isEmployed": false}`, r)
		log.Println(jsonData)
		io.WriteString(w, "Index")
	}))

	// This handler can work as a Lambda, or a local web server.
	if os.Getenv("RUN_WEBSERVER") != "" {
		fmt.Println("Listening on port 8000")
		http.ListenAndServe("localhost:8000", http.DefaultServeMux)
		return
	}

	// Start the Lambda handler.
	awsapigatewayv2handler.ListenAndServe(http.DefaultServeMux)

}
