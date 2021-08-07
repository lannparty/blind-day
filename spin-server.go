package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func getFirstPage(location string, destinationType string, radius string, apiKey string) (interface{}, []interface{}) {
	var nearbyDestinationCache []interface{}
	resp, err := http.Get("https://maps.googleapis.com/maps/api/place/nearbysearch/json?location=" + location + "&radius=" + radius + "&type=" + destinationType + "&key=" + apiKey)
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
		nearbyDestinationCache = append(nearbyDestinationCache, bodyMap["results"].([]interface{})[index].(map[string]interface{})["name"])
	}
	return bodyMap["next_page_token"], nearbyDestinationCache
}

func getNextPage(location string, destinationType string, radius string, apiKey string, nextPageToken interface{}) (interface{}, []interface{}) {
	var nearbyDestinationCache []interface{}
	resp, err := http.Get("https://maps.googleapis.com/maps/api/place/nearbysearch/json?location=" + location + "&radius=" + radius + "&type=" + destinationType + "&key=" + apiKey + "&pagetoken=" + nextPageToken.(string))
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
	fmt.Println(bodyMap)
	for index, _ := range bodyMap["results"].([]interface{}) {
		nearbyDestinationCache = append(nearbyDestinationCache, bodyMap["results"].([]interface{})[index].(map[string]interface{})["name"])
	}
	return bodyMap["next_page_token"], nearbyDestinationCache
}

func homePage(w http.ResponseWriter, r *http.Request) {
	var nearbyDestinations []interface{}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	location := r.URL.Query()["location"][0]
	destinationType := r.URL.Query()["type"][0]
	radius := r.URL.Query()["radius"][0]
	apiKey := "AIzaSyCoKlHo--DXlETVKWKQlW1zt6MJRIutb0c"
	nextPageToken, nearbyDestinationsCache := getFirstPage(location, destinationType, radius, apiKey)
	for _, j := range nearbyDestinationsCache {
		nearbyDestinations = append(nearbyDestinations, j)
	}
	for nextPageToken != nil {
		time.Sleep(2 * time.Second)
		nextPageToken, nearbyDestinationsCache = getNextPage(location, destinationType, radius, apiKey, nextPageToken)
		for _, j := range nearbyDestinationsCache {
			nearbyDestinations = append(nearbyDestinations, j)
		}
	}
	for _, j := range nearbyDestinations {
		fmt.Fprintf(w, j.(string)+"\n")
	}
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
	handleRequests()
}
