package services

import (
	"encoding/json"
	"fmt"
	"hotel-scraper/models"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type AmadeusService struct {
	APIKey string
	Secret string
	URL    string
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func NewAmadeusService(apiKey, secret string) *AmadeusService {
	return &AmadeusService{
		APIKey: apiKey,
		Secret: secret,
	}
}

func (a *AmadeusService) GetAccessToken() (string, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", a.APIKey)
	data.Set("client_secret", a.Secret)

	url := "https://test.api.amadeus.com/v1/security/oauth2/token"
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	checkError(err)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	checkError(err)
	defer resp.Body.Close()

	var tokenResp TokenResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&tokenResp)
	checkError(err)

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	return tokenResp.AccessToken, nil
}

func (a *AmadeusService) SearchHotels(city string,radius int) ([]models.Hotel, error) {
	token, err := a.GetAccessToken()
	checkError(err)

	URL := "https://test.api.amadeus.com/v1/reference-data/locations/hotels/by-geocode"

	params := url.Values{}
    params.Set("latitude", "48.8566")   
    params.Set("longitude", "2.3522")   
    params.Set("radius", fmt.Sprintf("%d", radius))

	fullURL := URL + "?" + params.Encode()

	req, err := http.NewRequest("GET", fullURL, nil)
	checkError(err)

	req.Header.Set("Authorization", "Bearer "+token)

	client := http.Client{}
	resp, err := client.Do(req)
	checkError(err)
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var apiResponse struct {
		Data []struct {
			Name    string `json:"name"`
			Address struct {
				CityName    string `json:"cityName"`
				CountryName string `json:"countryName"`
				Street      string `json:"street"`
			} `json:"address"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&apiResponse)
	checkError(err)

	var hotels []models.Hotel
	for _, item := range apiResponse.Data {
		var address string
		if item.Address.Street != "" {
			address = item.Address.Street + ", " + item.Address.CityName
		}
		hotel := models.Hotel{
			HotelName:      item.Name,
			HotelAddress:   address,
			City:           item.Address.CityName,
			Price:          0.0,
			Rating:         0.0,
			Description:    "",
			HotelCreatedAt: time.Now(),
		}
		hotels = append(hotels, hotel)
	}
	return hotels, nil
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
