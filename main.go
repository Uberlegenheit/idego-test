package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"idego-test/models"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	weatherAPIURL     = "https://api.weather.gov"
	coordinatesAPIURL = "https://photon.komoot.io/api"
)

func getCoordinates(city string) (float64, float64, error) {
	params := url.Values{}
	params.Add("q", city)

	link := fmt.Sprintf("%s?%s", coordinatesAPIURL, params.Encode())
	resp, err := http.Get(link)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	var coordinates models.Coordinates
	err = json.NewDecoder(resp.Body).Decode(&coordinates)
	if err != nil {
		return 0, 0, err
	}

	if len(coordinates.Features) != 0 {
		for _, v := range coordinates.Features {
			if v.Properties.Type == "city" {
				if len(v.Geometry.Coordinates) > 1 {
					return v.Geometry.Coordinates[0], v.Geometry.Coordinates[1], nil
				}

				break
			}
		}
	}

	if len(coordinates.Features[0].Geometry.Coordinates) > 1 {
		return coordinates.Features[0].Geometry.Coordinates[0], coordinates.Features[0].Geometry.Coordinates[1], nil
	}

	return 0, 0, nil
}

func getWeather(city string) (*models.Forecast, string, error) {
	lon, lat, err := getCoordinates(city)
	if err != nil {
		return nil, "", err
	}

	link := fmt.Sprintf("%s/points/%f,%f", weatherAPIURL, lat, lon)
	resp, err := http.Get(link)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	var info models.PointInfo
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		return nil, "", err
	}

	link = fmt.Sprintf(info.Properties.Forecast)
	resp2, err := http.Get(link)
	if err != nil {
		return nil, "", err
	}
	defer resp2.Body.Close()

	var forecast models.Forecast
	err = json.NewDecoder(resp2.Body).Decode(&forecast)
	if err != nil {
		return nil, "", err
	}

	return &forecast, info.Properties.RelativeLocation.Properties.State, nil
}

func getAlerts(state string) (*models.AlertResponse, error) {
	link := fmt.Sprintf("%s/alerts/active?area=%s", weatherAPIURL, state)
	resp, err := http.Get(link)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var alertResp models.AlertResponse
	err = json.NewDecoder(resp.Body).Decode(&alertResp)
	if err != nil {
		return nil, err
	}

	return &alertResp, nil
}

func displayWeather(weather *models.Forecast) {
	if len(weather.Properties.Periods) > 0 {
		period := weather.Properties.Periods[0]
		fmt.Printf("Temperature: %d°F\n", period.Temperature)
		fmt.Printf("Wind Speed: %s\n", period.WindSpeed)
		fmt.Printf("Forecast: %s\n", period.ShortForecast)
	}
}

func displayAlerts(alerts *models.AlertResponse) {
	if len(alerts.Features) == 0 {
		fmt.Println("none")
	}
	for _, feature := range alerts.Features {
		fmt.Printf("Alert: %s - %s\n", feature.Properties.Event, feature.Properties.Headline)
	}
}

func main() {
	city := flag.String("city", "oregon", "Name of the city in the USA")
	updateInterval := flag.Duration("time", time.Minute, "Update interval for data in minutes")
	flag.Parse()

	if *city == "" {
		fmt.Println("Please provide a city name using the -city flag.")
		os.Exit(1)
	}

	for {
		weather, state, err := getWeather(*city)
		if err != nil {
			fmt.Printf("Error fetching weather data: %v\n", err)
			continue
		}
		fmt.Printf("Current Weather in %s:\n", *city)
		displayWeather(weather)

		alerts, err := getAlerts(state)
		if err != nil {
			fmt.Printf("Error fetching alert data: %v\n", err)
			continue
		}
		fmt.Printf("\n%s:\n", alerts.Title)
		displayAlerts(alerts)

		time.Sleep(*updateInterval)
		fmt.Println("\n---------- UPDATE ----------")
	}
}
