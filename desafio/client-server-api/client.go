package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Cotacao struct {
	Dollar string
}

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "get", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	//io.Copy(os.Stdout, res.Body)

	var cotacao Cotacao
	saveCotacao(&cotacao)

}

func saveCotacao(c *Cotacao) {
	file, err := os.Create("cotacao.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao criar arquivo: %w", err)
	}
	defer file.Close()

}
