package models

type TripRoute struct {
	Distance float64
	Duration float64
	Geometry struct {
		Coordinates [][]float64
	}
}
