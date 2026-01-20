package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gabriel-ulisses-andrade/goexpert-desafio-multithreading-api/internal/dto"
	"github.com/gabriel-ulisses-andrade/goexpert-desafio-multithreading-api/internal/entity"
)

type CEPHandler struct {
	URLBrasilApi string
	URLViaCep    string
}

func NewCEPHandler(urlBrasilApi, urlViaCep string) *CEPHandler {
	return &CEPHandler{
		URLBrasilApi: urlBrasilApi,
		URLViaCep:    urlViaCep,
	}
}

func (cepHandler *CEPHandler) ConsultarCep(w http.ResponseWriter, r *http.Request) {
	cep := r.PathValue("cep")
	ch := make(chan entity.CEP)

	go cepHandler.ConsultarCepViaCep(ch, cep)
	go cepHandler.ConsultarCepBrasilApi(ch, cep)

	select {
	case result := <-ch:
		response, err := json.Marshal(result)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Erro ao converter resposta"))
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	case <-time.After(time.Second * 1):
		w.WriteHeader(http.StatusRequestTimeout)
		w.Write([]byte("Timeout: não foi possível obter o CEP no tempo esperado"))
		return
	}
}

func (cepHandler *CEPHandler) ConsultarCepBrasilApi(ch chan entity.CEP, cep string) {
	resp, err := http.Get(fmt.Sprintf(cepHandler.URLBrasilApi, cep))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var cepResponse dto.BrasilApiCepResponse
	err = json.Unmarshal(content, &cepResponse)
	cepEntity, err := entity.NewCepBrasilApiCep(cepResponse)
	if err != nil {
		panic(err)
	}
	ch <- *cepEntity
}

func (cepHandler *CEPHandler) ConsultarCepViaCep(ch chan entity.CEP, cep string) {
	resp, err := http.Get(fmt.Sprintf(cepHandler.URLViaCep, cep))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var cepResponse dto.ViaCepApiCEPResponse
	err = json.Unmarshal(content, &cepResponse)
	if err != nil {
		panic(err)
	}
	cepEntity, err := entity.NewCepViaCep(cepResponse)
	if err != nil {
		panic(err)
	}
	ch <- *cepEntity
}
