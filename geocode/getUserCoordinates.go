package geocode

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

/*
// TODO: use url encoder or something to write this, formatting yourself is not normal iggg
// this function builds the url to call the get function for the geocoder api
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
}*/

// Building url the correct way with net/url
func urlBuilder(city string, state string, country string) string {
	//pass env variable for api key at docker image run time
	//docker run -e GEOCODE_API_KEY=<apikay> testenv
	//geoCodeApiKey := //"INSERT LOCAL API KEY"

	geoCodeApiKey := os.Getenv("GEOCODE_API_KEY")
	if geoCodeApiKey == "" {
		log.Fatalln("GEOCODE API KEY NOT SET")
	}
	log.Println("GEOCODE_API_KEY", geoCodeApiKey)

	baseUrl := "https://geocode.maps.co/search"

	url, err := url.Parse(baseUrl)
	if err != nil {
		log.Fatalf("failed to parse URL: %v", err)
	}

	queryParams := url.Query()
	queryParams.Set("q", city+" "+state+" "+country)
	queryParams.Set("api_key", geoCodeApiKey)

	url.RawQuery = queryParams.Encode()

	geoCodeUrl := url.String()
	fmt.Println(geoCodeUrl)
	return geoCodeUrl
}

// this function uses the url and actually sends the request and returns the respone
func GetCoordinates(city string, state string, country string) ([]byte, error) {

	url := urlBuilder(city, state, country)
	//fmt.Printf(url)
	method := "GET"

	client := &http.Client{}

	geoCodeRequest, error := http.NewRequest(method, url, nil)

	if error != nil {
		log.Println(error)
		return nil, fmt.Errorf("request Creation Failed: %v", error)
	}

	geoCodeResult, error := client.Do(geoCodeRequest)
	if error != nil {
		log.Println(error)
		return nil, fmt.Errorf("API request failed: %v", error)

	}
	defer geoCodeResult.Body.Close()

	if geoCodeResult.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api request failed with status: %v", geoCodeResult.Status)
	}

	geoCodeResponseBody, error := io.ReadAll(geoCodeResult.Body)
	if error != nil {
		log.Println(error)
		return nil, fmt.Errorf("failed to read response body: %v", error)
	}

	fmt.Println(string(geoCodeResponseBody))
	return geoCodeResponseBody, nil
}

// this is how you take responses in GO and make parsing easier
type GeocodeResponse []struct {
	Latitude  string `json:"lat"`
	Longitude string `json:"lon"`
}

// this function is parsing the response data to return the lattitude and longitude
func ProcessGeoData(city, state, country string) (float64, float64, error) {
	jsonBody, err := GetCoordinates(city, state, country)
	if err != nil {
		log.Printf("Error getting user's coordinates %v", err)
		return 0, 0, err
	}

	var response GeocodeResponse
	if err := json.Unmarshal(jsonBody, &response); err != nil {
		log.Printf("Error unmarshalling JSON: %v", err)
		return 0, 0, err

	}

	if len(response) == 0 {
		log.Printf("No results found for the given location")
	}
	getFirstResult := response[0]
	fmt.Println(getFirstResult.Latitude)
	fmt.Println(getFirstResult.Longitude)

	lattitudeFloat, err := strconv.ParseFloat(getFirstResult.Latitude, 64)
	if err != nil {
		fmt.Println("problem converting lattitude string to float")
		return 0, 0, err

	}

	longitudeFloat, err := strconv.ParseFloat(getFirstResult.Longitude, 64)
	if err != nil {
		fmt.Println("problem converting longitute string to float")
		return 0, 0, err
	}

	return lattitudeFloat, longitudeFloat, nil
}
