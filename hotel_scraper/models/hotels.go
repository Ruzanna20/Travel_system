package models

import "time"

type Hotel struct {
	HotelID        int       `json:"hotel_ID" db:"hotelID"`
	HotelName      string    `json:"hotel_name" db:"hotelName"`
	HotelAddress   string    `json:"hotel_address" db:"hotelAddress"`
	City           string    `json:"City" db:"city"`
	Price          float64   `json:"Price" db:"price"`
	Rating         float64   `json:"Rating" db:"rating"`
	Description    string    `json:"Description" db:"description"`
	HotelCreatedAt time.Time `json:"hotelCreatedAt" db:"hotelCreatedAt"`
}
