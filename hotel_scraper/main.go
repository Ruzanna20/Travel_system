package main

import (
	"database/sql"
	"fmt"
	"hotel-scraper/repository"
	"hotel-scraper/services"
	"log"
	"os"
	_ "github.com/lib/pq" 

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Starting Hotel Data Scraper.")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Ping error", err)
	}

	fmt.Println("DB connected successfully.")

	hotelRepo := repository.NewHotelRepository(db)
	fmt.Println("Hotel repo initialized.")

	apiKey := os.Getenv("AMADEUS_API_KEY")
	apiSecret := os.Getenv("AMADEUS_SECRET")

	amadeusService := services.NewAmadeusService(apiKey, apiSecret)
	fmt.Println("Amadeus service initialized.")

	city := "Paris"
	radius := 300
	fmt.Printf("Searching hotels in %s \n", city)
	hotels, err := amadeusService.SearchHotels(city,radius)
	if err != nil {
		log.Fatal("Failed search", err)
	}

	fmt.Printf("Found hotels %d \n",len(hotels))

	for _, hotel := range hotels {
		err := hotelRepo.CreateHotel(hotel)
		if err != nil {
			log.Printf("Failed to save %s : %v", hotel.HotelName, err)
			continue
		}
	}
	if err != nil {
		log.Fatal("Hotels are not in db")
	} else {
			fmt.Println("All hotels saved successfully.")

	}

}
