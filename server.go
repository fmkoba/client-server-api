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

type CurrencyResponse struct {
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

type Currency struct {
	ID        int    `json:"id"`
	Value     string `json:"value"`
	CreatedAt string `json:"created_at"`
}

type Response struct {
	Value string `json:"Dólar"`
}

func main() {
	log.Default().Println("Starting Server")
	sqliteDatabase, err := sql.Open("sqlite3", "./sqlite-database.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	createTable(sqliteDatabase)
	defer sqliteDatabase.Close()

	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		currency := getCurrency()
		AddCurrency(sqliteDatabase, Currency{Value: currency, CreatedAt: time.Now().String()})

		apiResponse := Response{Value: currency}
		JSONResponse, err := json.Marshal(apiResponse)
		if err != nil {
			panic(err)
		}
		log.Default().Println(string(JSONResponse))

		w.Write(JSONResponse)
	})
	http.ListenAndServe(":8080", nil)
}

func getCurrency() string {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var response CurrencyResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		panic(err)
	}
	return response.Usdbrl.Bid
}

func createTable(db *sql.DB) {
	createCurrencyTableSQL := `CREATE TABLE if NOT EXISTS currency (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"value" TEXT,
		"created_at" timestamp NOT NULL DEFAULT (DATETIME('now','localtime'))
	);`

	log.Println("Create currency table...")
	statement, err := db.Prepare(createCurrencyTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	log.Println("currency table created")
}

func AddCurrency(db *sql.DB, currency Currency) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	insertCurrencySQL := `INSERT INTO currency(value, created_at) VALUES (?, ?)`
	statement, err := db.Prepare(insertCurrencySQL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.ExecContext(ctx, currency.Value, currency.CreatedAt)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
