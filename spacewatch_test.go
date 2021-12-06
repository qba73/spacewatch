package spacewatch_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/qba73/spacewatch"
)

func GetISSStatus(apikey string) (spacewatch.Status, error) {
	return spacewatch.Status{
		Lat:           47.2,
		Long:          -131.29,
		Timezone:      "America/vancouver",
		CloudCoverage: 27,
		DayPart:       "day",
		IsVisible:     true,
	}, nil
}

func TestISSHandler_Get(t *testing.T) {
	handler := spacewatch.ISSStatusHandler{
		ApiKey:        "123",
		StatusChecker: GetISSStatus,
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	handler.Get(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("got %d, want %d", res.Code, http.StatusOK)
	}
}
