package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Fatalf("Erro ao criar requisição: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Erro ao realizar requisição: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Erro ao ler resposta: %v", err)
	}

	err = ioutil.WriteFile("cotacao.txt", body, 0644)
	if err != nil {
		log.Fatalf("Erro ao salvar cotação: %v", err)
	}

	log.Println("Cotação salva com sucesso: ", string(body))
}
