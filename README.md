## Apresentação

Este projeto visa resolver o desafio de desenvolver uma API integrada a dois serviços de consulta de CEP, devolvendo a consulta mais rápida usando técnicas de multithreading. Desafio proposto pela fullcycle no curso de Go Expert.

#### Os requisitos da solução são:
- Realizar duas requisições serão feitas simultaneamente para as seguintes APIs:
  - https://brasilapi.com.br/api/cep/v1/01153000 + cep
  - http://viacep.com.br/ws/" + cep + "/json/
- Acatar a API que entregar a resposta mais rápida e descartar a resposta mais lenta.
- O resultado da request deverá ser exibido no command line com os dados do endereço, bem como qual API a enviou.
- Limitar o tempo de resposta em 1 segundo. Caso contrário, o erro de timeout deve ser exibido.

## Como Executar o Projeto

1. Clone o repositório

```shell 
git clone https://github.com/gabriel-ulisses-andrade/goexpert-desafio-multithreading-api-/tree/main
cd goexpert-desafio-multithreading-api-
```

2. Execute o servidor

```shell
cd cmd/server
go mod tity
go run main.go
```
