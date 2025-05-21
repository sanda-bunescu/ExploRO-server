package responses

type TouristicAttractionResponse struct {
	TouristicAttractionId uint    `json:"id"`
	Name                  string  `json:"attraction_name"`
	Description           string  `json:"attraction_description"`
	Category              string  `json:"category"`
	ImageUrl              string  `json:"image_url"`
	OpenHours             string  `json:"open_hours"`
	Fee                   float32 `json:"fee"`
	Link                  string  `json:"link"`
}
