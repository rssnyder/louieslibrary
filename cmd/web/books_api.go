package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// ISBNResponse
// Structure for industryIdentifiers json field
type ISBNResponse struct {
	Type    		string  `json:"type"`
	Identifier 	string	`json:"identifier"`
}

// ImageResponse
// Structure for imageLinks json field
type ImageResponse struct {
	SmallThumbnail	string	`json:"smallThumbnail"`
	Thumbnail 			string	`json:"thumbnail"`
	Small 					string	`json:"small"`
	Medium 					string	`json:"medium"`
	Large 					string	`json:"large"`
}

// DataResponse
// Structure for volumeInfo json field
type DataResponse struct {
	Title 							string 					`json:"title"`
	Subtitle 						string 					`json:"subtitle"`
	Publisher 					string 					`json:"publisher"`
	PublishedDate 			string 					`json:"publishedDate"`
	Description 				string 					`json:"description"`
	PageCount 					int				 			`json:"pageCount"`
	MaturityRating 			string 					`json:"maturityRating"`
	Authors 						[]string 				`json:"authors"`
	IndustryIdentifiers	[]ISBNResponse	`json:"industryIdentifiers"`
	Categories 					[]string 				`json:"categories"`
	ImageLinks 					ImageResponse 	`json:"imageLinks"`
}

// RetailResponse
// Structure for retailPrice json field
type RetailResponse struct {
	Amount 				float64	`json:"amount"`
	CurrencyCode	string 	`json:"currencyCode"`
}

// SaleResponse
// Structure for saleInfo json field
type SaleResponse struct {
	Retail	RetailResponse	`json:"retailPrice"`
}

// VolumeResponse
// Structure for response json
type VolumeResponse struct {
	Id				string    		`json:"id"`
	Data 			DataResponse 	`json:"volumeInfo"`
	SaleInfo	SaleResponse	`json:"saleInfo"`
}

// GetBookInfo
// Retrive info from books.google api on book
// Requires volumeid of book and gcloud api key
func GetBookInfo(volume_id, api_key string) VolumeResponse {

	// Empty response struct
	var response VolumeResponse

	// Make books.google api call
	http_response, err := http.Get(fmt.Sprintf("https://www.googleapis.com/books/v1/volumes/%s?key=%s", volume_id, api_key))
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	// Read in the response data
	response_data, err := ioutil.ReadAll(http_response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Fill volumeresponse with json data from api
	json.Unmarshal(response_data, &response)

	// Return the data from the books.google api
	return response
}