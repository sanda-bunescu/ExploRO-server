package responses

import "time"

type TripPlanResponse struct {
	Id        uint      `json:"id"`
	TripName  string    `json:"trip_name"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	GroupName string    `json:"group_name"`
	CityName  string    `json:"city_name"`
	CityId    uint      `json:"city_id"`
}
