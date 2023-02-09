// Desafio FullCycle Multithreading
//
// Referência: https://blog.devgenius.io/building-a-go-rest-client-in-2022-1ba8bb4c2201
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	urlApiCep string = "https://cdn.apicep.com/file/apicep/%s.json"
	urlViaCep string = "https://viacep.com.br/ws/%s/json/"
)

type CEPResp struct {
	URL  string
	Body string
}

func main() {
	chanApiCep := make(chan CEPResp)
	chanViaCep := make(chan CEPResp)

	go Worker(urlApiCep, "01310-200", chanApiCep, 0)
	go Worker(urlViaCep, "01310200", chanViaCep, 0)

	select {
	case apiCep := <-chanApiCep:
		fmt.Printf("URL: %s\nResposta: %s\n", apiCep.URL, apiCep.Body)
	case viaCep := <-chanViaCep:
		fmt.Printf("URL: %s\nResposta: %s\n", viaCep.URL, viaCep.Body)
	case <-time.After(time.Second):
		log.Fatalln("Tempo de resposta excedido")
	}
}

func Worker(url string, cep string, bodyCannel chan<- CEPResp, delay time.Duration) {

	time.Sleep(delay)

	cr := CEPResp{URL: fmt.Sprintf(url, cep)}
	r, err := http.NewRequest("GET", cr.URL, nil)
	if err != nil {
		close(bodyCannel)
	}
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		close(bodyCannel)
	}

	body, err := io.ReadAll(res.Body)
	_ = res.Body.Close()

	if err != nil {
		close(bodyCannel)
	}

	cr.Body = string(body)

	bodyCannel <- cr
}
