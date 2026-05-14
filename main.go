package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// openweather api key,found in your env, rename it if need be
var apiKey = os.Getenv("WEATHER_API_KEY")

type Location struct {
	City    string
	Country string
}

type WeatherData struct {
	Weather []struct {
		Main string `json:"main"`
	} `json:"weather"`
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	}
}

func (l Location) Info() string {
	return strings.Title(l.City) + "," + strings.ToTitle(l.Country)
}

func prompt(scanner *bufio.Scanner, question string) string {
	fmt.Print(question)
	scanner.Scan()
	return scanner.Text()
}

const urlScrub = "https://api.openweathermap.org/data/2.5/weather?"

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var units string
	cityIn := prompt(scanner, "City? ")
	countryIn := prompt(scanner, "Country? ")
	unitsIn := strings.ToLower(prompt(scanner, "F or C? "))
	if unitsIn == "f" {
		units = "units=imperial"
	} else {
		units = "units=metric"
	}
	location := Location{City: cityIn, Country: countryIn}

	fullURL := fmt.Sprintf("%sq=%s&%s&appid=%s", urlScrub, url.QueryEscape(location.Info()),
		units, apiKey)

	data, err := http.Get(fullURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer data.Body.Close()

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

	fmt.Printf("\nCurrent temp in %s is %.1f, and '%s'.", location.Info(),
		weather.Main.Temp, weather.Weather[0].Main)
	fmt.Printf("\nHumidity is %d%%", weather.Main.Humidity)
	fmt.Println()
}
