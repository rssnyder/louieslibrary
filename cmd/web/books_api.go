package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// A Response struct to map the Entire Response

type ISBNResponse struct {
	Type    		string  `json:"type"`
	Identifier 	string	`json:"identifier"`
}

type ImageResponse struct {
	SmallThumbnail	string	`json:"smallThumbnail"`
	Thumbnail 			string	`json:"thumbnail"`
	Small 					string	`json:"small"`
	Medium 					string	`json:"medium"`
	Large 					string	`json:"large"`
}

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

type RetailResponse struct {
	Amount 				float64	`json:"amount"`
	CurrencyCode	string 	`json:"currencyCode"`
}
type SaleResponse struct {
	Retail	RetailResponse	`json:"retailPrice"`
}

type VolumeResponse struct {
	Id				string    		`json:"id"`
	Data 			DataResponse 	`json:"volumeInfo"`
	SaleInfo	SaleResponse	`json:"saleInfo"`
}

func GetBookInfo(volume_id, api_key string) VolumeResponse {
	var response VolumeResponse

	http_response, err := http.Get(fmt.Sprintf("https://www.googleapis.com/books/v1/volumes/%s?key=%s", volume_id, api_key))
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	response_data, err := ioutil.ReadAll(http_response.Body)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(response_data, &response)

	return response
}