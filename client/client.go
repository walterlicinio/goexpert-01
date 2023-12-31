package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Fatal("Erro criando request:", err)
	}
	log.Println("Sucesso criando request")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("Erro executando request:", err)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatal("Erro executando request:", resp.Status)
	}
	log.Println("Sucesso executando request")
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Erro lendo resposta:", err)
	}
	log.Println("Sucesso lendo resposta")

	value := string(body)
	content := fmt.Sprintf("Dólar: %s", value)
	log.Println(value)
	if err := os.WriteFile("./data/cotacao.txt", []byte(content), 0644); err != nil {
		log.Fatal("Erro escrevendo arquivo:", err)
	}
	log.Println("Sucesso escrevendo arquivo")
}
