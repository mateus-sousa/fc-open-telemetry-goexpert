package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/mateus-sousa/fc-open-telemetry-goexpert/servico_b/config"
	"io"
	"log"
	"net/http"
)

type CEP struct {
	Number string `json:"cep"`
}

type ResponseViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
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

type ResponseHTTP struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

var cfg *config.Conf

func main() {
	var err error
	cfg, err = config.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}
	r := chi.NewRouter()
	r.Post("/getweather", getWeather)
	fmt.Println("listening in port :8081")
	http.ListenAndServe(":8081", r)
}

func getWeather(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	var cep CEP
	err = json.Unmarshal(body, &cep)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	req, err := http.NewRequestWithContext(r.Context(), "GET", fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep.Number), nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("deu ruim 1")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("deu ruim 2")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	var responseViaCEP ResponseViaCEP
	err = json.Unmarshal(resBody, &responseViaCEP)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	req, err = http.NewRequestWithContext(
		r.Context(),
		"GET",
		fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s",
			cfg.WeatherToken,
			responseViaCEP.Localidade,
		), nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	resBody, err = io.ReadAll(res.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	var responseWeather ResponseWeather
	err = json.Unmarshal(resBody, &responseWeather)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	tempF := responseWeather.Current.TempC*1.8 + 32
	tempK := responseWeather.Current.TempC + 273
	response := ResponseHTTP{
		City:  responseWeather.Location.Name,
		TempC: responseWeather.Current.TempC,
		TempF: tempF,
		TempK: tempK,
	}
	handleResponse(w, res.StatusCode, &response)
}

func handleResponse(w http.ResponseWriter, statusCode int, response *ResponseHTTP) {
	if statusCode == http.StatusOK {
		responseBytes, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseBytes)
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
	w.WriteHeader(http.StatusInternalServerError)
}
