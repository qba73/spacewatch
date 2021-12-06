package iss

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

// Response represents response data received
// upon successful call to the ISS API.
type Response struct {
	Timestamp   int64
	Message     string
	ISSPosition struct {
		Lat  string `json:"latitude"`
		Long string `json:"longitude"`
	} `json:"iss_position"`
}

// Position represents geographical coordinates
// of the International Space Station for the given time.
type Position struct {
	Lat  float64
	Long float64
}

// Client is the International Space Station (ISS) client.
//
// More information about the ISS API can be obtained on
// the website: http://open-notify.org/Open-Notify-API/
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// New returns a default ISS Client API.
func New() *Client {
	return &Client{
		BaseURL: "http://api.open-notify.org/iss-now.json",
		HTTPClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

// Get returns International Space Station coordinates
// (latitude/longitude) at the time of the request.
func (c Client) Get() (Position, error) {
	req, err := http.NewRequest(http.MethodGet, c.BaseURL, nil)
	if err != nil {
		return Position{}, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return Position{}, fmt.Errorf("error contacting iss api: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Position{}, fmt.Errorf("unexpected response status code: %v", res.StatusCode)
	}
	return parseResponse(res.Body)
}

func parseResponse(r io.Reader) (Position, error) {
	var res Response
	if err := json.NewDecoder(r).Decode(&res); err != nil {
		return Position{}, err
	}
	lat, long, err := toFloat64(res.ISSPosition.Lat, res.ISSPosition.Long)
	if err != nil {
		return Position{}, err
	}
	return Position{Lat: lat, Long: long}, nil
}

func toFloat64(lat, long string) (float64, float64, error) {
	lat64, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("converting latitude %s to float64 %v", lat, err)
	}
	long64, err := strconv.ParseFloat(long, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("converting longitude %s to float64 %v", long, err)
	}
	return lat64, long64, nil
}

// GetPosition returns current position of the International Space Station
// It returns latitude and longitude coordinates or error if the position
// cannot be determined.
func GetPosition() (float64, float64, error) {
	pos, err := New().Get()
	if err != nil {
		return 0, 0, errors.New("unable to retrieve position coordinates")
	}
	return pos.Lat, pos.Long, nil
}
