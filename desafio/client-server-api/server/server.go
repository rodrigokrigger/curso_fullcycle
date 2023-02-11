package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
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

type Response struct {
	Dolar string
	Error string
}

type CotacaoDb struct {
	ID    int
	Dolar float64
	Date  string
}

const URL_COTACAO = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
const SRV_TIMEOUT = 2000 * time.Millisecond
const DB_TIMEOUT = 100 * time.Millisecond
const SRV_PORT = ":8080"

func main() {
	// criar server
	log.Printf("Server started (%v).\n", SRV_PORT)
	http.HandleFunc("/cotacao", handlerCotacao)
	http.ListenAndServe(SRV_PORT, nil)
}

func handlerCotacao(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// cria contexto do consumo da API
	//ctxApi, cancelCtx := context.WithTimeout(r.Context(), SRV_TIMEOUT)
	//defer cancelCtx()

	ctxApi := r.Context()

	log.Println("Request started.")
	defer log.Println("Request ended.")

	select {
	case <-ctxApi.Done():
		log.Println("Request canceled by client.")
		returnResponse(w, "0", http.StatusGone, "Request canceled by client.")

	case <-time.After(SRV_TIMEOUT):
		log.Println("Request timeout.")
		returnResponse(w, "0", http.StatusRequestTimeout, "Request timeout.")

	default:
	}

	// busca cotacao
	log.Println("Initialize getCotacao.")
	cotacao, err := getCotacao()
	if err != nil {
		log.Println(err)
		returnResponse(w, "0", http.StatusInternalServerError, "Dollar Exchange request failed.")
		return
	}

	// salva no banco a cotação
	log.Println("Initialize saveCotacao.")
	fDolar, err := strconv.ParseFloat(cotacao.Usdbrl.Bid, 64)
	if err != nil {
		log.Println(err)
		returnResponse(w, "0", http.StatusInternalServerError, "Problem to convert value.")
		return
	}
	log.Printf("US$ value = %v.\n", fDolar)
	newCotacao := CotacaoDb{
		Dolar: fDolar,
		Date:  cotacao.Usdbrl.CreateDate,
	}
	err = saveCotacao(newCotacao)
	if err != nil {
		log.Println(err)
		http.Error(w, "Saving db failed.", http.StatusInternalServerError)
		return
	}

	// response da cotação
	log.Println("Initialize return JSON Response.")
	returnResponse(w, cotacao.Usdbrl.Bid, http.StatusOK, "")
	log.Println("Request performed successfuly.")

}

func getCotacao() (Cotacao, error) {

	var cotacao Cotacao

	ctx, cancelCtx := context.WithTimeout(context.Background(), SRV_TIMEOUT)
	defer cancelCtx()

	// cria request da url de cotação com contexto
	log.Println("Initialize URL Cotacao.")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	//req, err := http.NewRequest(http.MethodGet, "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		log.Println("Dollar exchange request failed.")
		return cotacao, err
	}

	// executa a request
	log.Println("Execute URL Cotacao.")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Dollar exchange request failed.")
		return cotacao, err
	}
	defer res.Body.Close()

	// le o body do response
	log.Println("Reading Response URL Cotacao.")
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Failed to read request body.")
		return cotacao, err
	}

	log.Printf("body: %s\n", body)

	// parse do response json
	log.Println("Parsing Response URL Cotacao.")
	err = json.Unmarshal(body, &cotacao)
	if err != nil {
		log.Println("Failed to parse Json.")
		return cotacao, err
	}

	return cotacao, nil
}

func saveCotacao(cotacao CotacaoDb) error {

	ctxDb, cancelCtxDb := context.WithTimeout(context.Background(), DB_TIMEOUT)
	defer cancelCtxDb()

	select {
	case <-ctxDb.Done():
		log.Println("database timeout exceeded")

	case <-time.After(DB_TIMEOUT):
		log.Println("Saving DB performed successfuly.")

	default:
	}

	// conecta com o banco SQLITE
	db, err := sql.Open("sqlite3", "cotacao.db")
	if err != nil {
		log.Println("Failed to open Database.")
		return err
	}
	defer db.Close()

	qry := "create table if not exists cotacao (id integer primary key autoincrement, val_dolar real not null, create_date text not null)"
	_, err = db.Exec(qry)
	if err != nil {
		log.Println("Failed to create Table.")
		return err
	}

	qry = `insert into cotacao(val_dolar, create_date) values(?,?)`
	//_, err = db.ExecContext(ctxDb, qry, cotacao.Dolar, cotacao.Date)
	_, err = db.Exec(qry, cotacao.Dolar, cotacao.Date)
	if err != nil {
		log.Println("Failed to insert register into Database.")
		return err
	}

	return nil
}

func returnResponse(w http.ResponseWriter, cotacao string, httpStatus int, error string) {
	log.Println("Initialize return JSON Response.")
	w.WriteHeader(httpStatus)
	response := Response{
		Dolar: cotacao,
		Error: error,
	}
	json.NewEncoder(w).Encode(response)
}
