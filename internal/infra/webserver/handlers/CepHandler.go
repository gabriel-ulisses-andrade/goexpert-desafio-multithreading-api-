package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gabriel-ulisses-andrade/goexpert-desafio-multithreading-api/internal/dto"
	"github.com/gabriel-ulisses-andrade/goexpert-desafio-multithreading-api/internal/entity"
)

func AddError(chErr chan error, err error) {
	if err != nil {
		chErr <- err
	}
}

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

func RegistrarLog(message string) {
	log.Println(" | " + message)
}

func (cepHandler *CEPHandler) ConsultarCep(w http.ResponseWriter, r *http.Request) {
	cep := r.PathValue("cep")
	ch := make(chan entity.CEP)
	chErr := make(chan error)

	go cepHandler.ConsultarCepViaCep(ch, chErr, cep)
	go cepHandler.ConsultarCepBrasilApi(ch, chErr, cep)

	select {
	case result := <-ch:
		response, err := json.Marshal(result)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Erro ao converter resposta"))
			return
		}
		RegistrarLog(result.ToString())
		RegistrarLog(result.GetServico())
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)

	case err := <-chErr:
		RegistrarLog(fmt.Sprintf("Erro ao consultar CEP: %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Erro ao consultar CEP: %s", err.Error())))

	case <-time.After(time.Second * 1):
		RegistrarLog("Timeout: não foi possível obter o CEP no tempo esperado")
		w.WriteHeader(http.StatusRequestTimeout)
		w.Write([]byte("Timeout: não foi possível obter o CEP no tempo esperado"))
		return
	}
}

func (cepHandler *CEPHandler) ConsultarCepBrasilApi(ch chan entity.CEP, chErr chan error, cep string) {
	resp, err := http.Get(fmt.Sprintf(cepHandler.URLBrasilApi, cep))
	AddError(chErr, err)

	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	AddError(chErr, err)

	var cepResponse dto.BrasilApiCepResponse
	err = json.Unmarshal(content, &cepResponse)
	cepEntity, err := entity.NewCepBrasilApiCep(cepResponse)
	AddError(chErr, err)

	ch <- *cepEntity
}

func (cepHandler *CEPHandler) ConsultarCepViaCep(ch chan entity.CEP, chErr chan error, cep string) {
	resp, err := http.Get(fmt.Sprintf(cepHandler.URLViaCep, cep))
	AddError(chErr, err)

	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	AddError(chErr, err)

	var cepResponse dto.ViaCepApiCEPResponse
	err = json.Unmarshal(content, &cepResponse)
	AddError(chErr, err)

	cepEntity, err := entity.NewCepViaCep(cepResponse)
	AddError(chErr, err)

	ch <- *cepEntity
}
