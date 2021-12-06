package iss_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/qba73/spacewatch/iss"
)

func TestClient(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"timestamp": 1638559834, "message": "success", "iss_position": {"latitude": "29.9314", "longitude": "11.3786"}}`)
	}))
	defer ts.Close()

	issClient := iss.New()
	issClient.BaseURL = ts.URL

	got, err := issClient.Get()
	if err != nil {
		t.Fatal(err)
	}

	want := iss.Position{
		Lat:  29.9314,
		Long: 11.3786,
	}

	if !cmp.Equal(got, want) {
		t.Errorf("%s\n", cmp.Diff(got, want))
	}
}

func TestClientEmptyResponseShouldError(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{}`)
	}))
	defer ts.Close()

	issClient := iss.New()
	issClient.BaseURL = ts.URL

	_, err := issClient.Get()
	if err == nil {
		t.Fatal(err)
	}
}

func TestClientMissingCoordinates_Latitude(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"timestamp": 1638559834, "message": "success", "iss_position": {"latitude": "", "longitude": "11.3786"}}`)
	}))
	defer ts.Close()

	issClient := iss.New()
	issClient.BaseURL = ts.URL

	_, err := issClient.Get()
	if err == nil {
		t.Fatal(err)
	}
}

func TestClientMissingCoordinates_Longitude(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"timestamp": 1638559834, "message": "success", "iss_position": {"latitude": "29.9314", "longitude": ""}}`)
	}))
	defer ts.Close()

	issClient := iss.New()
	issClient.BaseURL = ts.URL

	_, err := issClient.Get()
	if err == nil {
		t.Fatal(err)
	}
}
