package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	var latitude, longitude float64
	fmt.Printf("The weather for:\n")
	fmt.Printf("Latitude? ")
	fmt.Scanf("%f", &latitude)
	fmt.Printf("Longitude? ")
	fmt.Scanf("%f", &longitude)
	var URL string = fmt.Sprintf("https://api.openweathermap.org/data/3.0/onecall?lat=%f&lon=%f&appid=f90f2252ffbb060ba8d8c3bd7e7e500d&units=metric", latitude, longitude)

	response, err := http.Get(URL)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	type Current struct {
		Temp       float64 `json:"temp"`
		Feels_like float64 `json:"feels_like"`
		Humidity   float64 `json:"humidity"`
	}

	type Weather struct {
		Current Current `json:"current"`
	}

	var result Weather
	text := string(body)
	err = json.Unmarshal([]byte(text), &result)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("current temp: %v, ", result.Current.Temp)
	fmt.Printf("feels like: %v, ", result.Current.Feels_like)
	fmt.Printf("humidity: %v\n", result.Current.Humidity)
}
