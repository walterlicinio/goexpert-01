package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://server:8080/cotacao", nil)
	if err != nil {
		fmt.Println("Erro criando request:", err)
		return
	}
	fmt.Println("Sucesso criando request")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Erro executando request:", err)
		return
	}
	fmt.Println("Sucesso executando request")
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erro lendo resposta:", err)
		return
	}
	fmt.Println("Sucesso lendo resposta")

	value := string(body)
	content := fmt.Sprintf("Dólar: %s", value)
	fmt.Println(value)
	if err := os.WriteFile("/data/cotacao.txt", []byte(content), 0644); err != nil {
		fmt.Println("Erro escrevendo arquivo:", err)
	}
	fmt.Println("Sucesso escrevendo arquivo")
}