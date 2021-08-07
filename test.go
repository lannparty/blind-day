package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	location := r.URL.Query()["location"][0]
	destinationType := r.URL.Query()["type"][0]
	apiKey := ""
	resp, err := http.Get("https://maps.googleapis.com/maps/api/place/nearbysearch/json?location=" + location + "&radius=10000&type=" + destinationType + "&key=" + apiKey)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(body)
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
	handleRequests()
}
