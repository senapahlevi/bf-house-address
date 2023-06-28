package address

type ResponseCalculate struct {
	Haversine       string  `json:"haversine"`
	Euclidean       string  `json:"euclidean"`
	OriginLong      float64 `json:"origin_long"`
	OriginLat       float64 `json:"origin_lat"`
	DestinationLat  float64 `json:"destination_lat"`
	DestinationLong float64 `json:"destination_long"`
	OtherLat        float64 `json:"other_lat"`
	OtherLong       float64 `json:"other_long"`
}
