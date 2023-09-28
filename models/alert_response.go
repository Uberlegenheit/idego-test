package models

type AlertResponse struct {
	Title    string `json:"title"`
	Features []struct {
		Properties struct {
			Event    string `json:"event"`
			Headline string `json:"headline"`
		} `json:"properties"`
	} `json:"features"`
}
