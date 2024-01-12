package gateway

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

func TestForwardRequest(t *testing.T) {

	t.Run("should forward requests and pass the body back", func(t *testing.T) {
		defer gock.Off() // Flush pending mocks after test execution

		gock.New("http://api.example.com").
			Get("/api/resource").
			Reply(200).
			JSON(map[string]string{"foo": "bar"})

		// Creating a mock Application
		app := Application{
			configuration: &Configuration{
				Routes: []Route{
					{Prefix: "/api", Host: "api.example.com"},
					{Prefix: "/app", Host: "app.example.com"},
				},
			},
		}

		// Creating a mock HTTP request
		req, err := http.NewRequest("GET", "http://localhost:3333/api/resource", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Creating a mock HTTP response recorder
		w := httptest.NewRecorder()

		// Calling the forwardRequest function with the mock Application, request, and response recorder
		err = app.forwardRequest(w, req)

		assert.NoError(t, err)
		assert.Equal(t, w.Code, 200)
		assert.JSONEq(t, `{"foo": "bar"}`, w.Body.String())
	})

	t.Run("should forward requests with query params and pass the body back", func(t *testing.T) {
		defer gock.Off() // Flush pending mocks after test execution

		gock.New("http://api.example.com").
			Get("/api/resource").
			MatchParam("name", "blablabla").
			Reply(200).
			JSON(map[string]string{"foo": "bar"})

		// Creating a mock Application
		app := Application{
			configuration: &Configuration{
				Routes: []Route{
					{Prefix: "/api", Host: "api.example.com"},
					{Prefix: "/app", Host: "app.example.com"},
				},
			},
		}

		// Creating a mock HTTP request
		req, err := http.NewRequest("GET", "http://localhost:3333/api/resource?name=blablabla", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Creating a mock HTTP response recorder
		w := httptest.NewRecorder()

		// Calling the forwardRequest function with the mock Application, request, and response recorder
		err = app.forwardRequest(w, req)

		assert.NoError(t, err)
		assert.Equal(t, w.Code, 200)
		assert.JSONEq(t, `{"foo": "bar"}`, w.Body.String())
	})

	t.Run("should forward post with query params and pass the body back", func(t *testing.T) {
		defer gock.Off() // Flush pending mocks after test execution

		gock.New("http://api.example.com").
			Post("/api/resource").
			MatchParam("name", "blablabla").
			Reply(200).
			JSON(map[string]string{"foo": "bar"})

		// Creating a mock Application
		app := Application{
			configuration: &Configuration{
				Routes: []Route{
					{Prefix: "/api", Host: "api.example.com"},
					{Prefix: "/app", Host: "app.example.com"},
				},
			},
		}

		// Creating a mock HTTP request
		req, err := http.NewRequest("POST", "http://localhost:3333/api/resource?name=blablabla", nil)
		if err != nil {
			t.Fatal(err)
		}

		// Creating a mock HTTP response recorder
		w := httptest.NewRecorder()

		// Calling the forwardRequest function with the mock Application, request, and response recorder
		err = app.forwardRequest(w, req)

		assert.NoError(t, err)
		assert.Equal(t, w.Code, 200)
		assert.JSONEq(t, `{"foo": "bar"}`, w.Body.String())
	})

}
