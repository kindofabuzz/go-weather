package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// found in your env, rename it if need be
var apiKey = os.Getenv("WEATHER_API_KEY")

type Location struct {
	City    string
	Country string
}

type WeatherData struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

func (l Location) Info() string {
	return l.City + "," + l.Country
}

func prompt(scanner *bufio.Scanner, question string) string {
	fmt.Print(question)
	scanner.Scan()
	return scanner.Text()
}

const url = "https://api.openweathermap.org/data/2.5/weather?"

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cityIn := prompt(scanner, "City?")
	countryIn := prompt(scanner, "Country?")
	location := Location{City: cityIn, Country: countryIn}

	fullURL := fmt.Sprintf("%sq=%s&units=imperial&appid=%s", url, location.Info(), apiKey)

	data, err := http.Get(fullURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(data.Body)

	body, err := io.ReadAll(data.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	var weather WeatherData
	err = json.Unmarshal(body, &weather)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("\nCurrent temp in %s is %.1f\n", location.Info(), weather.Main.Temp)

}
