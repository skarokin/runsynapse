package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"context"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"

	"github.com/skarokin/runsynapse/go/inits"
	"github.com/skarokin/runsynapse/go/handlers"
)

func main() {
	secrets, err := inits.InitSecrets()
	if err != nil {
		log.Fatalf("Failed to initialize secrets: %v", err)
	}

	supabaseClient, err := inits.NewSupabaseClient(secrets.DatabaseURL) 
	if err != nil {
		log.Fatalf("Failed to create Supabase client: %v", err)
	}
	defer supabaseClient.Close()

	geminiClient, err := inits.NewGeminiClient(secrets.GeminiAPIKey)
	if err != nil {
		log.Fatalf("Failed to create Gemini client: %v", err)
	}

	s3Client, err := inits.NewS3Client(secrets.S3Region)
	if err != nil {
		log.Fatalf("Failed to create S3 client: %v", err)
	}

	handler := handlers.NewHandler(supabaseClient, geminiClient, s3Client, secrets.S3Bucket)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		log.Println("Running in AWS Lambda environment")
		lambda.Start(createLambdaHandler(handler))
	} else {
		log.Println("Starting HTTP server (development mode)")
		startHTTPServer(handler, port)
	}
}

func startHTTPServer(handler *handlers.Handler, port string) {
    log.Printf("Server running on port %s", port)
    log.Printf("Health check: http://localhost:%s/health", port)

	// handlers.NewHandler sets up the routes, no need to do it again here
    if err := http.ListenAndServe(":"+port, handler); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}

func createLambdaHandler(handler *handlers.Handler) func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
        // create path with query string
        path := request.Path
        if request.QueryStringParameters != nil && len(request.QueryStringParameters) > 0 {
            path += "?"
            for k, v := range request.QueryStringParameters {
                path += k + "=" + v + "&"
            }
            path = path[:len(path)-1] // remove trailing &
        }

        // initialize a new HTTP request from API Gateway event
        req, err := http.NewRequest(request.HTTPMethod, path, strings.NewReader(request.Body))
        if err != nil {
            return events.APIGatewayProxyResponse{
                StatusCode: http.StatusInternalServerError,
                Body:       "Error creating request: " + err.Error(),
            }, nil
        }

        // set headers from API Gateway event
        for k, v := range request.Headers {
            req.Header.Add(k, v)
        }

        // create response recorder (captures response from the handler)
        w := httptest.NewRecorder()

        // pass in the response recorder and request to the HTTP handler
        handler.ServeHTTP(w, req)
        
        // convert response from recorder to API Gateway response format
        res := events.APIGatewayProxyResponse{
            StatusCode: w.Code,
            Headers:    make(map[string]string),
            Body:       w.Body.String(),
        }
        
        // copy headers
        for k, v := range w.Header() {
            res.Headers[k] = v[0]
        }
        
        return res, nil
    }
}