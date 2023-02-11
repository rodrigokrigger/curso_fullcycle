package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Cotacao struct {
	Usdbrl struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

func main() {

	var cotacao Cotacao

	log.Println("Test started.")
	defer log.Println("Test ended.")

	// cria request
	log.Println("Sending request.")
	req, err := http.NewRequest(http.MethodGet, "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
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

	log.Printf("body: %s\n", body)

	// parse do response json
	log.Println("Parsing Response URL Cotacao.")
	err = json.Unmarshal(body, &cotacao)
	if err != nil {
		log.Printf("Failed to parse Json: %v", err)
	}

}
