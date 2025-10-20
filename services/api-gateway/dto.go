package main

import "github.com/atcheri/ride-booking-go/shared/types"

type tripPreviewRequest struct {
	UserId      string           `json:"userID"`
	PickUp      types.Coordinate `json:"pickup"`
	Destination types.Coordinate `json:"destination"`
}
