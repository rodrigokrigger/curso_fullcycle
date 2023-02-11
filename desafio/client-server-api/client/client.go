package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type Cotacao struct {
	Dolar string
	Error string
}

const URL_COTACAO = "http://localhost:8080/cotacao"
const REQUEST_TIMEOUT = 3000 * time.Millisecond

func main() {

	log.Println("Client started.")
	defer log.Println("Client ended.")

	// cria contexto
	log.Println("Creating Context.")
	ctx, cancelCtx := context.WithTimeout(context.Background(), REQUEST_TIMEOUT)
	defer cancelCtx()

	select {
	case <-ctx.Done():
		log.Println("Get request cancelled, timeout reached.")
		return
	case <-time.After(REQUEST_TIMEOUT):
		log.Println("Request processed.")
	default:
	}

	// cria request
	log.Println("Sending request.")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, URL_COTACAO, nil)
	if err != nil {
		log.Printf("Send Error: %v\n", err)
		return
	}

	// executa request
	log.Println("Receiving response.")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Receive Error: %v\n", err)
		return
	}
	defer res.Body.Close()

	// le o body do response
	log.Println("Reading response.")
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Read Error: %v\n", err)
		return
	}

	// parse do response json
	log.Println("Parsing response.")
	var cotacao Cotacao
	err = json.Unmarshal(body, &cotacao)
	if err != nil {
		log.Printf("Parse Error: %v\n", err)
		return
	}

	// salva em arquivo
	log.Println("Saving Cotacao.")
	if cotacao.Dolar != "" {
		fmt.Printf("Dólar: %v\n", cotacao.Dolar)
	} else {
		fmt.Printf("Erro: %v\n", cotacao.Error)
	}
	saveCotacao(&cotacao)
}

func saveCotacao(c *Cotacao) {

	file, err := os.Create("cotacao.txt")
	if err != nil {
		log.Printf("Create File Error: %v\n", err)
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("Dólar: %v", c.Dolar))
	if err != nil {
		log.Printf("Write File Error: %v\n", err)
	}
}
