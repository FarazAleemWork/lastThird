package geocode

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

func buildUrl(city string, state string, country string) string {
	replaceWhiteSpace := regexp.MustCompile(`\s+`)
	var formattedCity string = replaceWhiteSpace.ReplaceAllString(city, "%")
	var formattedState string = replaceWhiteSpace.ReplaceAllString(state, "%")
	var formattedCountry string = replaceWhiteSpace.ReplaceAllString(country, "%")

	godotenv.Load("c:/lastThird/environmentvar.env")
	geoCodeApiKey := os.Getenv("GEOCODE_API_KEY")
	fmt.Println(geoCodeApiKey)
	var geoCodeGetUrl string = "https://geocode.maps.co/search?q=" + formattedCity + "%" + formattedState + "%" + formattedCountry + "&api_key=" + geoCodeApiKey
	println(geoCodeGetUrl)
	return geoCodeGetUrl
}

func GetCoordinates(city string, state string, country string) ([]byte, error) {

	url := buildUrl(city, state, country)
	fmt.Printf(url)
	method := "GET"

	client := &http.Client{}

	geoCodeRequest, error := http.NewRequest(method, url, nil)

	if error != nil {
		fmt.Println(error)
		return nil, fmt.Errorf("Request Creation Failed: ", error)
	}

	geoCodeResult, error := client.Do(geoCodeRequest)
	if error != nil {
		fmt.Println(error)
		return nil, fmt.Errorf("API request failed: ", error)

	}
	defer geoCodeResult.Body.Close()

	if geoCodeResult.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %s", geoCodeResult.Status, geoCodeResult)
	}

	geoCodeResponseBody, error := io.ReadAll(geoCodeResult.Body)
	if error != nil {
		fmt.Println(error)
		return nil, fmt.Errorf("Failed to read response body: ", error)

	}

	fmt.Println(string(geoCodeResponseBody))
	return geoCodeResponseBody, nil
}

type GeocodeResponse []struct {
	Latitude  string `json:"lat"`
	Longitude string `json:"lon"`
}

func ProcessGeoData(city, state, country string) {
	jsonBody, error := GetCoordinates(city, state, country)
	if error != nil {
		log.Printf("Error getting user's coordinates", error)
		return
	}

	var response GeocodeResponse
	if err := json.Unmarshal(jsonBody, &response); err != nil {
		log.Printf("Error unmarshalling JSON: %v", err)
		return
	}

	if len(response) > 0 {
		getFirstResult := response[0]
		fmt.Println(getFirstResult.Latitude)
		fmt.Println(getFirstResult.Longitude)
	}
}
