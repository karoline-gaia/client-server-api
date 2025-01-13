package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Cotacao struct {
	Bid string `json:"bid"`
}

func main() {
	db, err := sql.Open("sqlite3", "cotacoes.db")
	if err != nil {
		log.Fatalf("Erro ao abrir banco de dados: %v", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS cotacoes (id INTEGER PRIMARY KEY, bid TEXT, timestamp DATETIME DEFAULT CURRENT_TIMESTAMP)")
	if err != nil {
		log.Fatalf("Erro ao criar tabela: %v", err)
	}

	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
		defer cancel()

		cotacao, err := buscarCotacao(ctx)
		if err != nil {
			log.Printf("Erro ao buscar cotação: %v", err)
			http.Error(w, "Erro ao buscar cotação", http.StatusInternalServerError)
			return
		}

		ctx, cancel = context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		err = salvarCotacao(ctx, db, cotacao.Bid)
		if err != nil {
			log.Printf("Erro ao salvar cotação: %v", err)
			http.Error(w, "Erro ao salvar cotação", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cotacao)
	})

	log.Println("Servidor iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func buscarCotacao(ctx context.Context) (Cotacao, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return Cotacao{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Cotacao{}, err
	}
	defer resp.Body.Close()

	var dados map[string]Cotacao
	if err := json.NewDecoder(resp.Body).Decode(&dados); err != nil {
		return Cotacao{}, err
	}

	return dados["USDBRL"], nil
}

func salvarCotacao(ctx context.Context, db *sql.DB, bid string) error {
	query := "INSERT INTO cotacoes (bid) VALUES (?)"
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, bid)
	return err
}
