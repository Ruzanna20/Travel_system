package repository

import (
	"database/sql"
	"hotel-scraper/models"
)

type HotelRepository struct {
	DB *sql.DB
}

func NewHotelRepository(db *sql.DB) *HotelRepository {
	return &HotelRepository{
		DB: db,
	}
}

func (r *HotelRepository) CreateHotel(hotel models.Hotel) error {
	query := `INSERT INTO hotels("hotelName", "hotelAddress", "city", "price", "rating", "description", "hotelCreatedAt")
			VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := r.DB.Exec(query,
		hotel.HotelName,
		hotel.HotelAddress,
		hotel.City,
		hotel.Price,
		hotel.Rating,
		hotel.Description,
		hotel.HotelCreatedAt)

	return err
}

