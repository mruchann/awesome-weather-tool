// I/O part has been done, now fix the api http call part.

// execution
// go build -o awesome-weather-tool main.go
// ./awesome-weather-tool

// TO DO LIST

// do string manipulation, change the values in the string
// CLI tool makes an HTTP call to retrieve the info
// implement structure, map sth. like that
// parse the json output
// http headers, request
// Sprint function

package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func main() {
	var latitude, longtitude float64
	fmt.Printf("The weather for:\n")
	fmt.Printf("Latitude? ")
	fmt.Scanf("%f", &latitude)
	fmt.Printf("Longtitude? ")
	fmt.Scanf("%f", &longtitude)
	//var URL string = fmt.Sprint("https://api.openweathermap.org/data/3.0/onecall?lat=", latitude, "&lon=", longtitude, "&appid=f90f2252ffbb060ba8d8c3bd7e7e500d", "&units=metric")

	var toyBody string = `{"temp":25, "humidity":50, "feels_like":30}`

	/*response, err := http.Get(URL)
	if err != nil {
		log.Fatalln(err)
	}*/

	/*body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}*/

	type Current struct {
		Temp       float64 `json:"temp"`
		Feels_like float64 `json:"feels_like"`
		Humidity   float64 `json:"humidity"`
	}
	var result Current
	//body_string := string(body)
	err := json.Unmarshal([]byte(toyBody), &result)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("current temp: %v, ", result.Temp)
	fmt.Printf("feels like: %v, ", result.Feels_like)
	fmt.Printf("humidity: %v\n", result.Humidity)
	fmt.Println(result)
}

/*
	parsing JSON format
	json.Marshall
*/
