package model

// Ad represents an advertisement ready to be served
type Ad struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	AdvertiserID string  `json:"advertiser_id"`
	Bid          float64 `json:"bid"`
	Placement    string  `json:"placement"`
	ServeURL     string  `json:"serve_url"`
	Score        float64 `json:"-"`
}
