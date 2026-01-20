package entity

import (
	"errors"
	"strings"

	"github.com/gabriel-ulisses-andrade/goexpert-desafio-multithreading-api/internal/dto"
)

var (
	CEPNaoInformado     = "CEP não informado"
	CEPInvalido         = "CEP invalido"
	EstadoNaoInformado  = "Estado não informado"
	CidadeNaoInformada  = "Cidade não informada"
	BairroNaoInformado  = "Bairro não informado"
	RuaNaoInformada     = "Rua não informada"
	ServicoNaoInformado = "Serviço não informado"
)

type CEP struct {
	Cep     string `json:"cep"`
	Estado  string `json:"estado"`
	Cidade  string `json:"cidade"`
	Bairro  string `json:"bairro"`
	Rua     string `json:"rua"`
	Servico string `json:"servico"`
}

func NewCepViaCep(cep dto.ViaCepApiCEPResponse) (*CEP, error) {
	cepEntity := &CEP{
		Cep:     strings.Replace(cep.Cep, "-", "", -1),
		Estado:  cep.Uf,
		Cidade:  cep.Localidade,
		Bairro:  cep.Bairro,
		Rua:     cep.Logradouro,
		Servico: "ViaCep",
	}
	err := cepEntity.Validate()
	if err != nil {
		return nil, err
	}
	return cepEntity, nil
}

func NewCepBrasilApiCep(cep dto.BrasilApiCepResponse) (*CEP, error) {
	cepEntity := &CEP{
		Cep:     strings.Replace(cep.Cep, "-", "", -1),
		Estado:  cep.State,
		Cidade:  cep.City,
		Bairro:  cep.Neighborhood,
		Rua:     cep.Street,
		Servico: "BrasilApiCep",
	}
	err := cepEntity.Validate()
	if err != nil {
		return nil, err
	}
	return cepEntity, nil
}

func (cep *CEP) Validate() error {
	if cep.Cep == "" {
		return errors.New(CEPNaoInformado)
	}
	if len(cep.Cep) != 8 {
		return errors.New(CEPInvalido)
	}
	if cep.Estado == "" {
		return errors.New(EstadoNaoInformado)
	}
	if cep.Cidade == "" {
		return errors.New(CidadeNaoInformada)
	}
	if cep.Bairro == "" {
		return errors.New(BairroNaoInformado)
	}
	if cep.Rua == "" {
		return errors.New(RuaNaoInformada)
	}
	if cep.Servico == "" {
		return errors.New(ServicoNaoInformado)
	}
	return nil
}
