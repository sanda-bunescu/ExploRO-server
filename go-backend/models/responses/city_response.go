package responses

type CityResponse struct {
	CityId          uint   `json:"id"`
	CityName        string `json:"city_name"`
	CityDescription string `json:"city_description"`
	ImageURL        string `json:"image_url"`
}
