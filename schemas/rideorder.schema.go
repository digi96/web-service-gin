package schemas

type CreateRideOrder struct {
	ContactId  string `json:"contact_id" binding:"required"`
	RiderName  string `json:"rider_name" binding:"required"`
	RiderPhone string `json:"rider_phone" binding:"required"`
	Destinaion string `json:"destination" binding:"required"`
}

type UpdateRideOrder struct {
	PickUpAt string `json:"pickup_at" binding:"required"`
}
