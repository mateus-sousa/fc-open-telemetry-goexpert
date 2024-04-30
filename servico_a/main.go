package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/mateus-sousa/fc-open-telemetry-goexpert/servico_a/config"
	"io"
	"log"
	"net/http"
)

type CEP struct {
	Number string `json:"cep"`
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
	fmt.Println("listening in port :8080")
	http.ListenAndServe(":8080", r)
}

func getWeather(w http.ResponseWriter, r *http.Request) {
	log.Println("init request")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	var cep CEP
	err = json.Unmarshal(body, &cep)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	err = validateCEP(cep)
	if err != nil {
		log.Println(err)

		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("invalid zipcode"))
		return
	}
	log.Println("request to service B")
	req, err := http.NewRequestWithContext(context.Background(), "POST", fmt.Sprintf("%s/getweather", cfg.BaseUrl), bytes.NewBuffer(body))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	log.Println("request to service B susccesfully")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	var responseHTTP *ResponseHTTP
	err = json.Unmarshal(resBody, &responseHTTP)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	log.Println("go to handler response")
	handleResponse(w, res.StatusCode, responseHTTP)
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
func validateCEP(cep CEP) error {
	if len(cep.Number) != 8 {
		return errors.New("cep is invalid")
	}
	return nil
}
