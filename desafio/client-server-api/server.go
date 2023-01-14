package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
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
	http.HandleFunc("/cotacao", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	log.Println("Request started.")
	defer log.Println("Request ended.")

	err := getCotacao(w)
	if err != nil {
		log.Println(err)
		w.Write([]byte("{ \"error\": \"Dollar exchange request failed.\"}"))
	}

	select {
	case <-time.After(210 * time.Millisecond):
		log.Println("Request performed successfuly.")

	case <-ctx.Done():
		log.Println("Request canceled by client.")
		http.Error(w, "{ \"error\": \"Request canceled by client.\"}", http.StatusRequestTimeout)

	}

}

func getCotacao(w http.ResponseWriter) error {

	req, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if err != nil {
		log.Println("Dollar exchange request failed.")
		return err
	}
	defer req.Body.Close()

	res, err := io.ReadAll(req.Body)
	if err != nil {
		log.Println("Failed to read request body.")
		return err
	}

	var dataCotacao Cotacao
	err = json.Unmarshal(res, &dataCotacao)
	if err != nil {
		log.Println("Failed to parse Json.")
		return err
	}

	w.Write([]byte("{ \"dollar\": \"" + dataCotacao.Usdbrl.Bid + "\"}"))

	saveCotacao(dataCotacao.Usdbrl.Bid)

	return nil

}

func saveCotacao(cotacao string) {

}
