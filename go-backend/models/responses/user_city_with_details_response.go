package responses

type UserCityWithDetails struct {
	CityID          uint   `json:"id" gorm:"column:id"`
	CityName        string `json:"city_name"`
	CityDescription string `json:"city_description"`
	ImageURL        string `json:"image_url"`
}
