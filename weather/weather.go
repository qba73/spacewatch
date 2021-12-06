package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type response struct {
	Count int `json:"count"`
	Data  []struct {
		Rh           float64 `json:"rh"`
		Pod          string  `json:"pod"`
		Lon          float64 `json:"lon"`
		Pres         int     `json:"pres"`
		Timezone     string  `json:"timezone"`
		ObTime       string  `json:"ob_time"`
		CountryCode  string  `json:"country_code"`
		Clouds       int     `json:"clouds"`
		Ts           int64   `json:"ts"`
		SolarRad     float64 `json:"solar_rad"`
		StateCode    string  `json:"state_code"`
		CityName     string  `json:"city_name"`
		WindSpd      float64 `json:"wind_spd"`
		WindCdirFull string  `json:"wind_cdir_full"`
		WindCdir     string  `json:"wind_cdir"`
		Slp          int     `json:"slp"`
		Vis          int     `json:"vis"`
		HAngle       int     `json:"h_angle"`
		Sunset       string  `json:"sunset"`
		Dni          float64 `json:"dni"`
		Dewpt        float64 `json:"dewpt"`
		Snow         int     `json:"snow"`
		Uv           float64 `json:"uv"`
		Precip       int     `json:"precip"`
		WindDir      int     `json:"wind_dir"`
		Sunrise      string  `json:"sunrise"`
		Ghi          float64 `json:"ghi"`
		Dhi          float64 `json:"dhi"`
		Aqi          int     `json:"aqi"`
		Lat          float64 `json:"lat"`
		Weather      struct {
			Icon        string `json:"icon"`
			Code        int    `json:"code"`
			Description string `json:"description"`
		} `json:"weather"`
		Datetime  string  `json:"datetime"`
		Temp      float64 `json:"temp"`
		Station   string  `json:"station"`
		ElevAngle float64 `json:"elev_angle"`
		AppTemp   float64 `json:"app_temp"`
	} `json:"data"`
	Minutely []struct {
		TimestampUtc   string `json:"timestamp_utc"`
		Snow           int    `json:"snow"`
		Temp           int    `json:"temp"`
		TimestampLocal string `json:"timestamp_local"`
		Ts             int    `json:"ts"`
		Precip         int    `json:"precip"`
	} `json:"minutely"`
}

// Conditions represents information about
// weather status (cloud coverage) for the
// given lat and long coordinates in the
// given timezone.
type Condition struct {
	Lat       float64
	Long      float64
	Timezone  string
	LocalTime time.Time

	// Indicates if it's night or day
	// at the lat/long in local time.
	// 'd' day, 'n' night
	DayNight string

	// Cloud coverage measured in %
	Clouds int
}

// Position holds information about latitude and longitude.
type Location struct {
	Lat  float64
	Long float64
}

// Client is an API client for the weather api.
type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

// NewFromEnv knows how to create a new weather client.
// It requires a valid API KEY to be exported as
// the ENV Var: WEATHERBIT_KEY.
//
// To read more how to get the API KEY checkout
// waeatherbit docs: https://www.weatherbit.io/api
func NewFromEnv() *Client {
	return New(os.Getenv("WEATHERBIT_KEY"))
}

// New knows how to create a new weather client.
// A valid APIKEY for the weatherbit service is required.
//
// To read more how to get the API KEY checkout
// waeatherbit docs: https://www.weatherbit.io/api
func New(apikey string) *Client {
	return &Client{
		BaseURL: "https://api.weatherbit.io",
		APIKey:  apikey,
		HTTPClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

// Get returns weather condition for the given geographcal position.
func (c Client) Get(l Location) (Condition, error) {
	url := fmt.Sprintf("%s/v2.0/current?lat=%.4f&lon=%.4f&key=%s", c.BaseURL, l.Lat, l.Long, c.APIKey)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Condition{}, err
	}
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return Condition{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Condition{}, fmt.Errorf("error getting information from weather api: %v", err)
	}
	return parseResponse(res.Body)
}

// calculateLocalTime knows how to compute time of the day
// for the given location. The location represents a timezone,
// for example: "Africa/Johannesburg".
func calculateLocalTime(location string) (time.Time, error) {
	loc, err := time.LoadLocation(location)
	if err != nil {
		return time.Time{}, fmt.Errorf("loading location %s, %v", location, err)
	}
	return time.Now().In(loc), nil
}

// parseResponse takes response body received from the weather api
// and, if successful returns weather condition, error otherwise.
func parseResponse(r io.Reader) (Condition, error) {
	var res response
	if err := json.NewDecoder(r).Decode(&res); err != nil {
		return Condition{}, fmt.Errorf("decoding response: %v", err)
	}
	if len(res.Data) < 1 {
		return Condition{}, errors.New("missing data in response body received from weather service")
	}

	lt, err := calculateLocalTime(res.Data[0].Timezone)
	if err != nil {
		return Condition{}, fmt.Errorf("calculating local time: %v", err)
	}

	cc := Condition{
		Lat:       res.Data[0].Lat,
		Long:      res.Data[0].Lon,
		Timezone:  res.Data[0].Timezone,
		LocalTime: lt,
		DayNight:  res.Data[0].Pod,
		Clouds:    res.Data[0].Clouds,
	}
	return cc, nil
}

// GetCondition takes latitude and longitude and returns
// current weather condition for the provided location coords.
//
// GetCondition uses a default weather client that reads API KEY
// from the Env Var: WEATHERBIT_KEY.
//
// To read more how to get the API KEY checkout
// waeatherbit docs: https://www.weatherbit.io/api
func GetCondition(lat, long float64) (Condition, error) {
	l := Location{Lat: lat, Long: long}
	return NewFromEnv().Get(l)
}
