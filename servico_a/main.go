package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"io"
	"net/http"
)

type CEP struct {
	number string
}

func main() {
	r := chi.NewRouter()
	r.Post("/getweather", getWeather)
	http.ListenAndServe(":8080", r)
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
	err = validateCEP(cep)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("invalid zipcode"))
		return
	}
	req, err := http.NewRequestWithContext(context.Background(), "POST", "localhost:8081/getweather", bytes.NewBuffer(body))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	handleResponse(w, res.StatusCode, resBody)
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
