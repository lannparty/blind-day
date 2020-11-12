package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

//TO DO: Move this into a function to prevent crosstalk
var nearbyDestinations []interface{}

func getFirstPage(location string, destinationType string, apiKey string) interface{} {
	resp, err := http.Get("https://maps.googleapis.com/maps/api/place/nearbysearch/json?location=" + location + "&radius=10000&type=" + destinationType + "&key=" + apiKey)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var bodyMap map[string]interface{}
	err = json.Unmarshal(body, &bodyMap)
	if err != nil {
		log.Fatal(err)
	}
	for index, _ := range bodyMap["results"].([]interface{}) {
		nearbyDestinations = append(nearbyDestinations, bodyMap["results"].([]interface{})[index].(map[string]interface{})["name"])
	}
	return bodyMap["next_page_token"]
}

func getNextPage(location string, destinationType string, apiKey string, nextPageToken interface{}) interface{} {
	resp, err := http.Get("https://maps.googleapis.com/maps/api/place/nearbysearch/json?location=" + location + "&radius=10000&type=" + destinationType + "&key=" + apiKey + "&pagetoken=" + nextPageToken.(string))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var bodyMap map[string]interface{}
	err = json.Unmarshal(body, &bodyMap)
	if err != nil {
		log.Fatal(err)
	}
	for index, _ := range bodyMap["results"].([]interface{}) {
		nearbyDestinations = append(nearbyDestinations, bodyMap["results"].([]interface{})[index].(map[string]interface{})["name"])
	}
	return bodyMap["next_page_token"]
}

func homePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	location := r.URL.Query()["location"][0]
	destinationType := r.URL.Query()["type"][0]
	apiKey := "AIzaSyCoKlHo--DXlETVKWKQlW1zt6MJRIutb0c"
	nextPageToken := getFirstPage(location, destinationType, apiKey)
	for nextPageToken != nil {
		time.Sleep(2 * time.Second)
		nextPageToken = getNextPage(location, destinationType, apiKey, nextPageToken)
	}
	for _, j := range nearbyDestinations {
		fmt.Fprintf(w, j.(string)+"\n")
	}
	fmt.Println(nearbyDestinations)
	nearbyDestinations = nil
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
	handleRequests()
}
