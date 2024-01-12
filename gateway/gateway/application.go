package gateway

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Application struct {
	// in here we will be able to add "objects" will need to be shared across our application
	// for instance shared database connection, or similar things;
	// Having it here rather than it in main make it more easy to test, for instance in a unit test.
	configuration *Configuration
}

func NewApplication() Application {
	config, err := ReadConfigFromFile("configuration.yaml")
	if err != nil {
		log.Fatal("[NewApplication] Loaded configuration file", err)
	}
	log.Println("[NewApplication] Loaded configuration file")

	return Application{
		configuration: config,
	}
}

func (a *Application) Run() error {
	var port = "3333"
	if configuredPort := os.Getenv("PORT"); configuredPort != "" {
		port = configuredPort
	}

	http.HandleFunc("/", withErrorHandling(a.forwardRequest))

	log.Println("[Application] Starting Gateway Application at port", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil && err != http.ErrServerClosed {
		// if the server is closed in not the "normal way" we return the error
		return fmt.Errorf("server closed with an unexpected error: %w", err)
	}
	return nil
}

// withErrorHandling wrap a function that handle HTTP requests so that it's
// easier to handle error in a "Go" idiomatic way - just rerturn the error
func withErrorHandling(handleFun func(w http.ResponseWriter, r *http.Request) error) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// try to execute the request
		err := handleFun(w, r)
		// if there was an error we check, if it's of type `StatusError` we are
		// able to return a custom status code. Otherwise we have to return 500 - Internal Server Error
		switch err := err.(type) {
		case StatusError:
			w.WriteHeader(err.StatusCode)
			json.NewEncoder(w).Encode(err)
		default:
			// in this case we log since it's a sort of unexpected scenario
			log.Println("[withErrorHandling] - encountered unexpected error while processing request", r.Method, r.URL, err)
			var errorResponse = StatusError{StatusCode: http.StatusInternalServerError, Message: err.Error()}
			w.WriteHeader(errorResponse.StatusCode)
			json.NewEncoder(w).Encode(errorResponse)
		}
	}
}

func (a Application) forwardRequest(w http.ResponseWriter, r *http.Request) error {
	// find which route should be applied
	destHost, ok := a.identifyDestHost(r)
	if !ok {
		return StatusError{StatusCode: 404, Message: "route is not configured"}
	}

	var requestToBeForwarded, err = http.NewRequest(r.Method, replaceHost(r.URL, destHost), r.Body)
	if err != nil {
		return fmt.Errorf("failed to create the request to be forwarded: %w", err)
	}

	log.Println("[forwardRequest] - forwarding the request to", destHost)
	response, err := http.DefaultClient.Do(requestToBeForwarded)
	if err != nil {
		return fmt.Errorf("forwarded request failed: %w", err)
	}

	// forward the response that we got from the actual server
	// copy headers, status code and the whole body
	for headerKey, headerValues := range response.Header {
		for _, v := range headerValues {
			w.Header().Add(headerKey, v)
		}
	}
	w.WriteHeader(response.StatusCode)
	if _, err := io.Copy(w, response.Body); err != nil {
		return fmt.Errorf("failed to send response data: %w", err)
	}

	return nil

}

func (a Application) identifyDestHost(r *http.Request) (string, bool) {
	requestPath := r.URL.Path
	for _, route := range a.configuration.Routes {
		if strings.Contains(requestPath, route.Prefix) {
			return route.Host, true
		}
	}
	return "", false
}

type StatusError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func (e StatusError) Error() string {
	return e.Message
}

// since URL is a struct in Go it's not so easy to replace only the
// host. We extract it into a function so that it's easier to use.
func replaceHost(old *url.URL, newHost string) string {
	var new = url.URL{
		Opaque:      old.Opaque,
		Scheme:      "http", // only support HTTP, no HTTPs
		User:        old.User,
		Host:        newHost,
		Path:        old.Path,
		RawPath:     old.RawPath,
		OmitHost:    false,
		ForceQuery:  old.ForceQuery,
		RawQuery:    old.RawQuery,
		Fragment:    old.Fragment,
		RawFragment: old.RawFragment,
	}
	return new.String()
}
