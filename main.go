package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type apiConfigData struct {
	OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey"`
}

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

func loadApiConfig(filename string) (apiConfigData, error) {
	bytes, error := ioutil.ReadFile(filename)

	if error != nil {
		return apiConfigData{}, error
	}

	var c apiConfigData

	error = json.Unmarshal(bytes, &c)
	if error != nil {
		return apiConfigData{}, error
	}
	return c, nil
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello from Go \n"))
}

func query(city string) (weatherData, error) {
	apiConfigData, error := loadApiConfig(".apiConfig")
	if error != nil {
		fmt.Println(error)
		return weatherData{}, error
	}
	fmt.Println(apiConfigData.OpenWeatherMapApiKey)

	resp, error := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + apiConfigData.OpenWeatherMapApiKey + "&q=" + city)

	if error != nil {
		return weatherData{}, error
	}

	defer resp.Body.Close()

	var d weatherData

	if error := json.NewDecoder(resp.Body).Decode(&d); error != nil {
		return weatherData{}, error
	}
	return d, nil

}

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/weather/",
		func(w http.ResponseWriter, r *http.Request) {
			city := strings.SplitN(r.URL.Path, "/", 3)[2]
			data, err := query(city)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(data)

		})

	http.ListenAndServe(":8082", nil)

}
