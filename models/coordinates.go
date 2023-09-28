package models

type Coordinates struct {
	Features []struct {
		Geometry struct {
			Coordinates []float64 `json:"coordinates"`
			Type        string    `json:"type"`
		} `json:"geometry"`
		Type       string `json:"type"`
		Properties struct {
			OsmType     string    `json:"osm_type"`
			OsmId       int64     `json:"osm_id"`
			Extent      []float64 `json:"extent,omitempty"`
			Country     string    `json:"country"`
			OsmKey      string    `json:"osm_key"`
			Countrycode string    `json:"countrycode"`
			OsmValue    string    `json:"osm_value"`
			Name        string    `json:"name"`
			Type        string    `json:"type"`
			City        string    `json:"city,omitempty"`
			Postcode    string    `json:"postcode,omitempty"`
			Locality    string    `json:"locality,omitempty"`
			Street      string    `json:"street,omitempty"`
			District    string    `json:"district,omitempty"`
			Housenumber string    `json:"housenumber,omitempty"`
			County      string    `json:"county,omitempty"`
			State       string    `json:"state,omitempty"`
		} `json:"properties"`
	} `json:"features"`
	Type string `json:"type"`
}
