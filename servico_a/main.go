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
	r.Post("/users", Xpto)
	http.ListenAndServe(":8000", r)
}

func Xpto(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	var cep CEP
	err = json.Unmarshal(body, &cep)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
	err = validateCEP(cep)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("{'message': 'invalid zipcode'}"))
	}
	req, err := http.NewRequestWithContext(context.Background(), "POST", "localhost:8081", bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	if res.StatusCode != http.StatusOK {
		panic(err)
	}
}

func validateCEP(cep CEP) error {
	if len(cep.number) != 8 {
		return errors.New("cep is invalid")
	}
	return nil
}
