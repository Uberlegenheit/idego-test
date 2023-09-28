package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"idego-test/models"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	weatherAPIURL     = "https://api.weather.gov"
	coordinatesAPIURL = "https://photon.komoot.io/api"
)

func getCoordinates(city string, coordCh chan<- models.Coordinates, errCh chan<- error) {
	fmt.Println(city)
	link := fmt.Sprintf("%s?q=%s", coordinatesAPIURL, city)
	resp, err := http.Get(link)
	if err != nil {
		errCh <- err
		return
	}
	defer resp.Body.Close()

	var coordinates models.Coordinates
	err = json.NewDecoder(resp.Body).Decode(&coordinates)
	if err != nil {
		errCh <- err
		return
	}

	coordCh <- coordinates
}

func getWeather(lon, lat float64, areaName string, weatherCh chan<- *models.Forecast, stateCh chan<- string, errCh chan<- error) {
	link := fmt.Sprintf("%s/points/%f,%f", weatherAPIURL, lat, lon)
	resp, err := http.Get(link)
	if err != nil {
		errCh <- err
		return
	}
	defer resp.Body.Close()

	var info models.PointInfo
	err = json.NewDecoder(resp.Body).Decode(&info)
	if err != nil {
		errCh <- err
		return
	}

	link = fmt.Sprintf(info.Properties.Forecast)
	resp2, err := http.Get(link)
	if err != nil {
		errCh <- err
		return
	}
	defer resp2.Body.Close()

	var forecast models.Forecast
	err = json.NewDecoder(resp2.Body).Decode(&forecast)
	if err != nil {
		errCh <- err
		return
	}
	forecast.RequestedCity = areaName

	weatherCh <- &forecast
	stateCh <- info.Properties.RelativeLocation.Properties.State
}

func getAlerts(state string, alertsCh chan<- *models.AlertResponse, errCh chan<- error) {
	link := fmt.Sprintf("%s/alerts/active?area=%s", weatherAPIURL, state)
	resp, err := http.Get(link)
	if err != nil {
		errCh <- err
		return
	}
	defer resp.Body.Close()

	var alertResp models.AlertResponse
	err = json.NewDecoder(resp.Body).Decode(&alertResp)
	if err != nil {
		errCh <- err
		return
	}

	alertsCh <- &alertResp
}

func displayWeather(weather *models.Forecast) {
	fmt.Printf("\nWeather forecast for %s:\n", weather.RequestedCity)
	if len(weather.Properties.Periods) > 0 {
		period := weather.Properties.Periods[0]
		fmt.Printf("Temperature: %d°F\n", period.Temperature)
		fmt.Printf("Wind Speed: %s\n", period.WindSpeed)
		fmt.Printf("Forecast: %s\n", period.ShortForecast)
	}
}

func displayAlerts(alerts *models.AlertResponse) {
	fmt.Printf("\n%s:\n", alerts.Title)
	if len(alerts.Features) == 0 {
		fmt.Println("none")
	}
	for _, feature := range alerts.Features {
		fmt.Printf("Alert: %s - %s\n", feature.Properties.Event, feature.Properties.Headline)
	}
}

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		fmt.Println("Please provide at least a city name.")
		os.Exit(1)
	}

	citiesStr := args[0]
	cities := strings.Split(citiesStr, ",")

	updateIntervalStr := "1m"
	if len(args) > 1 {
		updateIntervalStr = args[1]
	}
	updateInterval, err := time.ParseDuration(updateIntervalStr)
	if err != nil {
		fmt.Println("Error parsing duration:", err)
		return
	}
	if updateInterval < (time.Second * 10) {
		updateInterval = time.Second * 10
	}
	fmt.Printf("Updating time is set to: %v\n", updateInterval)

	coordCh := make(chan models.Coordinates)
	weatherCh := make(chan *models.Forecast)
	alertsCh := make(chan *models.AlertResponse)
	stateCh := make(chan string)
	errCh := make(chan error)

	for _, city := range cities {
		go func(city string) {
			getCoordinates(city, coordCh, errCh)
		}(city)
		go func(city string) {
			for {
				select {
				case <-time.Tick(updateInterval):
					getCoordinates(city, coordCh, errCh)
				}
			}

		}(city)
	}

	go func() {
		for {
			select {
			case coordinates := <-coordCh:
				if len(coordinates.Features) > 0 {
					for _, v := range coordinates.Features {
						if v.Properties.Type == "city" && len(v.Geometry.Coordinates) > 1 {
							go getWeather(
								v.Geometry.Coordinates[0],
								v.Geometry.Coordinates[1],
								v.Properties.Name,
								weatherCh,
								stateCh,
								errCh,
							)
							break
						} else {
							go getWeather(
								coordinates.Features[0].Geometry.Coordinates[0],
								coordinates.Features[0].Geometry.Coordinates[1],
								v.Properties.Name,
								weatherCh,
								stateCh,
								errCh,
							)
							break
						}
					}
				}
			case weather := <-weatherCh:
				go displayWeather(weather)
			case alert := <-alertsCh:
				go displayAlerts(alert)
			case state := <-stateCh:
				go getAlerts(state, alertsCh, errCh)
			case err := <-errCh:
				fmt.Printf("Error: %v\n", err)
			}
		}
	}()

	select {}
}
