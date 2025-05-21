package requests

import "time"

type CreateTripPlanRequest struct {
	TripName  string    `json:"trip_name"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	GroupId   uint      `json:"group_id"`
	CityId    uint      `json:"city_id"`
}
