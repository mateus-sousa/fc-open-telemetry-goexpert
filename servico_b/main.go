package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type CEP struct {
	number string
}

type ResponseWeather struct {
	Location struct {
		Name           string  `json:"name"`
		Region         string  `json:"region"`
		Country        string  `json:"country"`
		Lat            float64 `json:"lat"`
		Lon            float64 `json:"lon"`
		TzId           string  `json:"tz_id"`
		LocaltimeEpoch int     `json:"localtime_epoch"`
		Localtime      string  `json:"localtime"`
	} `json:"location"`
	Current struct {
		LastUpdatedEpoch int     `json:"last_updated_epoch"`
		LastUpdated      string  `json:"last_updated"`
		TempC            float64 `json:"temp_c"`
		TempF            float64 `json:"temp_f"`
		IsDay            int     `json:"is_day"`
		Condition        struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
			Code int    `json:"code"`
		} `json:"condition"`
		WindMph    float64 `json:"wind_mph"`
		WindKph    float64 `json:"wind_kph"`
		WindDegree int     `json:"wind_degree"`
		WindDir    string  `json:"wind_dir"`
		PressureMb float64 `json:"pressure_mb"`
		PressureIn float64 `json:"pressure_in"`
		PrecipMm   float64 `json:"precip_mm"`
		PrecipIn   float64 `json:"precip_in"`
		Humidity   int     `json:"humidity"`
		Cloud      int     `json:"cloud"`
		FeelslikeC float64 `json:"feelslike_c"`
		FeelslikeF float64 `json:"feelslike_f"`
		VisKm      float64 `json:"vis_km"`
		VisMiles   float64 `json:"vis_miles"`
		Uv         float64 `json:"uv"`
		GustMph    float64 `json:"gust_mph"`
		GustKph    float64 `json:"gust_kph"`
	} `json:"current"`
}

func main() {
	getWeather()
	//r := chi.NewRouter()
	//r.Post("/getweather", getWeather)
	//http.ListenAndServe(":8081", r)
}

func getWeather() {
	//body, err := io.ReadAll(r.Body)
	//if err != nil {
	//	w.WriteHeader(http.StatusBadRequest)
	//	w.Write([]byte(err.Error()))
	//	return
	//}
	//var cep CEP
	//err = json.Unmarshal(body, &cep)
	//if err != nil {
	//	w.WriteHeader(http.StatusBadRequest)
	//	w.Write([]byte(err.Error()))
	//	return
	//}
	//err = validateCEP(cep)
	//if err != nil {
	//	w.WriteHeader(http.StatusUnprocessableEntity)
	//	w.Write([]byte("invalid zipcode"))
	//	return
	//}
	cep := CEP{number: "73403303"}
	req, err := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep.number), nil)
	if err != nil {
		//w.WriteHeader(http.StatusInternalServerError)
		//w.Write([]byte(err.Error()))
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("deu ruim 1")
		//w.WriteHeader(http.StatusInternalServerError)
		//w.Write([]byte(err.Error()))
		return
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("deu ruim 2")

		//w.WriteHeader(http.StatusInternalServerError)
		//w.Write([]byte(err.Error()))
		return
	}
	fmt.Println(string(resBody))
	req, err = http.NewRequestWithContext(
		context.Background(),
		"GET",
		fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s",
			"b31e34e7d8d842b4849182050243004",
			"Brasília",
		), nil)

	if err != nil {
		//w.WriteHeader(http.StatusInternalServerError)
		//w.Write([]byte(err.Error()))
		return
	}
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("deu ruim 1")
		//w.WriteHeader(http.StatusInternalServerError)
		//w.Write([]byte(err.Error()))
		return
	}
	resBody, err = io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("deu ruim 2")
		//w.WriteHeader(http.StatusInternalServerError)
		//w.Write([]byte(err.Error()))
		return
	}

	var responseWeather ResponseWeather
	err = json.Unmarshal(resBody, &responseWeather)
	if err != nil {
		fmt.Println(err)
		//w.WriteHeader(http.StatusInternalServerError)
		//w.Write([]byte(err.Error()))
		return
	}

	fmt.Printf("%+v", responseWeather)
	//handleResponse(w, res.StatusCode, resBody)
}

func handleResponse(w http.ResponseWriter, statusCode int, resBody []byte) {
	if statusCode == http.StatusOK {
		w.WriteHeader(http.StatusOK)
		w.Write(resBody)
		return
	} else if statusCode == http.StatusUnprocessableEntity {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("invalid zipcode"))
		return
	} else if statusCode == http.StatusNotFound {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("can not find zipcode"))
		return
	}
}

func validateCEP(cep CEP) error {
	if len(cep.number) != 8 {
		return errors.New("cep is invalid")
	}
	return nil
}
