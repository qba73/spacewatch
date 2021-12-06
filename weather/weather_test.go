package weather_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/qba73/spacewatch/weather"
)

// newTestServer
func newTestServer(testFile, wantURI string, t *testing.T) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		gotReqURI := r.RequestURI

		// Bail if URIs don't match.
		verifyURIs(wantURI, gotReqURI, t)
		f, err := os.Open(testFile)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		_, err = io.Copy(rw, f)
		if err != nil {
			t.Fatal(err)
		}
	}))
	return ts
}

// verifyURIs is a test helper function that verifies if provided URIs are the same.
func verifyURIs(wantURI, gotURI string, t *testing.T) {
	wantUri, err := url.Parse(wantURI)
	if err != nil {
		t.Fatalf("error parsing URI %q, %v", wantURI, err)
	}
	gotUri, err := url.Parse(gotURI)
	if err != nil {
		t.Fatalf("error parsing URI %q, %v", wantURI, err)
	}
	// Verify if paths of both URIs are the same.
	if wantUri.Path != gotUri.Path {
		t.Fatalf("want %q, got %q", wantUri.Path, gotUri.Path)
	}

	wantQuery, err := url.ParseQuery(wantUri.RawQuery)
	if err != nil {
		t.Fatal(err)
	}
	gotQuery, err := url.ParseQuery(gotUri.RawQuery)
	if err != nil {
		t.Fatal(err)
	}

	// Verify if query parameters match.
	if !cmp.Equal(wantQuery, gotQuery) {
		t.Fatalf("URIs are not equal, \n%s\n", cmp.Diff(wantQuery, gotQuery))
	}
}

func TestWeatherCloudCoverage(t *testing.T) {
	t.Parallel()

	ts := newTestServer("testdata/response-weatherbit.json", "/v2/current?lat=-12.3387&lon=-111.7409&key=APIKEY123", t)
	defer ts.Close()

	client := weather.New("APIKEY123")
	client.BaseURL = ts.URL

	got, err := client.Get(weather.Location{-12.3387, -111.7409})
	if err != nil {
		t.Fatalf("client.Get(weather.Location{-12.3387, -111.7409}) got: %v", err)
	}

	loc, err := time.LoadLocation("Pacific/Easter")
	if err != nil {
		t.Fatal(err)
	}

	want := weather.Condition{
		Lat:       -12.34,
		Long:      -111.74,
		LocalTime: time.Date(2021, 12, 04, 18, 03, 12, 884451, loc),
		Timezone:  "Pacific/Easter",
		DayNight:  "d",
		Clouds:    83,
	}

	if !cmp.Equal(got, want, cmpopts.IgnoreFields(weather.Condition{}, "LocalTime")) {
		t.Errorf("client.Get(weather.Location{-12.3387, -111.7409}) got:\n%s\n", cmp.Diff(got, want))
	}
}
