package responses

type StopPointResponse struct {
	Id                  uint                        `json:"id"`
	ItineraryId         uint                        `json:"itinerary_id"`
	TouristicAttraction TouristicAttractionResponse `json:"touristic_attraction"`
}
