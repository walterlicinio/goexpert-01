package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Cotacao struct {
	Bid string `json:"bid"`
}

var db *gorm.DB

func main() {
	initDB()
	http.HandleFunc("/cotacao", getCotacao)
	http.ListenAndServe(":8080", nil)
}

func initDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("./data/database.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Erro conectando à database:", err)
	}
	log.Println("Sucesso em initDB()")
	db.AutoMigrate(&Cotacao{})
}

func getCotacao(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequest(http.MethodGet, "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		log.Println("Erro criando request:", err)
		http.Error(w, "Erro criando request", http.StatusInternalServerError)
		return
	}
	log.Println("Sucesso criando request")

	req = req.WithContext(ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Erro fetchando cotação:", err)
		http.Error(w, "Erro fetchando cotação", http.StatusInternalServerError)
		return
	}
	log.Println("Sucesso fetchando cotação")
	defer resp.Body.Close()

	var result map[string]Cotacao
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("Erro decodando resposta:", err)
		http.Error(w, "Erro decodando resposta", http.StatusInternalServerError)
		return
	}
	log.Println("Sucesso decodando resposta")

	bid := result["USDBRL"].Bid
	fmt.Fprint(w, bid)

	dbCtx, dbCancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer dbCancel()

	if err := db.WithContext(dbCtx).Create(&Cotacao{Bid: bid}).Error; err != nil {
		log.Println("Erro criando cotação:", err)
	} else {
		log.Println("Sucesso criando cotação")
	}
}
