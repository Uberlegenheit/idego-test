package models

import "time"

type Forecast struct {
	Context       []interface{} `json:"@context"`
	RequestedCity string        `json:"-"`
	Type          string        `json:"type"`
	Geometry      struct {
		Type        string        `json:"type"`
		Coordinates [][][]float64 `json:"coordinates"`
	} `json:"geometry"`
	Properties struct {
		Updated           time.Time `json:"updated"`
		Units             string    `json:"units"`
		ForecastGenerator string    `json:"forecastGenerator"`
		GeneratedAt       time.Time `json:"generatedAt"`
		UpdateTime        time.Time `json:"updateTime"`
		ValidTimes        string    `json:"validTimes"`
		Elevation         struct {
			UnitCode string  `json:"unitCode"`
			Value    float64 `json:"value"`
		} `json:"elevation"`
		Periods []struct {
			Number                     int         `json:"number"`
			Name                       string      `json:"name"`
			StartTime                  time.Time   `json:"startTime"`
			EndTime                    time.Time   `json:"endTime"`
			IsDaytime                  bool        `json:"isDaytime"`
			Temperature                int         `json:"temperature"`
			TemperatureUnit            string      `json:"temperatureUnit"`
			TemperatureTrend           interface{} `json:"temperatureTrend"`
			ProbabilityOfPrecipitation struct {
				UnitCode string      `json:"unitCode"`
				Value    interface{} `json:"value"`
			} `json:"probabilityOfPrecipitation"`
			Dewpoint struct {
				UnitCode string  `json:"unitCode"`
				Value    float64 `json:"value"`
			} `json:"dewpoint"`
			RelativeHumidity struct {
				UnitCode string `json:"unitCode"`
				Value    int    `json:"value"`
			} `json:"relativeHumidity"`
			WindSpeed        string `json:"windSpeed"`
			WindDirection    string `json:"windDirection"`
			Icon             string `json:"icon"`
			ShortForecast    string `json:"shortForecast"`
			DetailedForecast string `json:"detailedForecast"`
		} `json:"periods"`
	} `json:"properties"`
}
