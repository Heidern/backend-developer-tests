package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/stackpath/backend-developer-tests/rest-service/pkg/models"
)

func TestGetHead(t *testing.T) {
	r := chi.NewRouter()

	r.Mount("/people", (&PersonController{}).Router())

	ts := httptest.NewServer(r)
	defer ts.Close()

	var p models.Person

	if response := request(t, ts, "GET", "/people/5b81b629-9026-450d-8e46-da4f8c7bd513", &p); response.StatusCode != 200 {
		t.Fatal(response)
	}

	if p.FirstName != "Jane" || p.LastName != "Doe" || p.PhoneNumber != "+1 (800) 555-1313" {
		t.Fatalf("invalid response: %+v", p)
	}
}

func request(t *testing.T, ts *httptest.Server, method string, path string, target interface{}) *http.Response {
	r, err := http.NewRequest(method, ts.URL+path, nil)
	if err != nil {
		t.Fatal(err)
		return nil
	}

	response, err := http.DefaultClient.Do(r)

	if err != nil {
		t.Fatal(err)
		return nil
	}

	defer response.Body.Close()

	json.NewDecoder(response.Body).Decode(target)

	return response
}
