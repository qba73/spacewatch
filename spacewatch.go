package spacewatch

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/qba73/spacewatch/iss"
	"github.com/qba73/spacewatch/weather"
)

// Status represents the ISS status report.
type Status struct {
	// Lat/Long are coordinates of the ISS.
	Lat           float64 `json:"lat"`
	Long          float64 `json:"long"`
	Timezone      string  `json:"timezone"`
	CloudCoverage int     `json:"cloud_coverage"`
	DayPart       string  `json:"day_part"`

	// IsVisible is calculated based on
	// cloud coverage and part of the day.
	IsVisible bool `json:"is_visible"`
}

// isVisible holds logic for calculating visibility of the ISS.
// It is assumed that the Space Station is visible during
// the night and when the cloud coverage is less than 30%.
func isVisible(cloudCoverage int, dayPart string) bool {
	return cloudCoverage <= 30 && dayPart == "d"
}

// Location holds coordinates information.
type Location struct {
	Lat  float64
	Long float64
}

// GetISSLocation is a high level function that knows how to return
// current location (lat/long) of the International Space Station.
//
// GetISSLocation uses default implementation of the iss client.
func GetISSLocation() (Location, error) {
	pos, err := iss.New().Get()
	if err != nil {
		return Location{}, err
	}
	p := Location{Lat: pos.Lat, Long: pos.Long}
	return p, nil
}

// GetISSStatus holds the core logic used for generating
// ISS  status update. It takes APIKEY required by the undelying
// weather service and return the Status struct, or error
// if any of the internal operation fail.
//
// GetISSStatus leverages deafult clients for ISS location
// and weather status.
func GetISSStatus(apikey string) (Status, error) {
	loc, err := GetISSLocation()
	if err != nil {
		return Status{}, err
	}
	weather, err := weather.New(apikey).Get(weather.Location{Lat: loc.Lat, Long: loc.Long})
	if err != nil {
		return Status{}, nil
	}
	dayPart := map[string]string{"d": "day", "n": "night"}
	return Status{
		Lat:           weather.Lat,
		Long:          weather.Long,
		Timezone:      weather.Timezone,
		CloudCoverage: weather.Clouds,
		DayPart:       dayPart[weather.DayNight],
		IsVisible:     isVisible(weather.Clouds, weather.DayNight),
	}, nil
}

// =======================================================
// ISS Handlers

type ISSStatusHandler struct {
	ApiKey        string
	Log           *log.Logger
	StatusChecker func(apikey string) (Status, error)
}

func (s *ISSStatusHandler) Get(w http.ResponseWriter, r *http.Request) {
	issStatus, err := s.StatusChecker(s.ApiKey)
	if err != nil {
		s.Log.Printf("error: getting weather report %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(issStatus)
	if err != nil {
		s.Log.Printf("error: marshalling weather report %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		s.Log.Printf("error writing weather report: %v", err)
	}
}
