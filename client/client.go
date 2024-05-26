package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func checkErr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 2000*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/cotacao", nil)
	checkErr(err)
	resp, err := http.DefaultClient.Do(req)
	checkErr(err)
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	checkErr(err)
	writeToFile(string(body))
}

func writeToFile(text string) {
	log.Default().Println("Writing to file")
	f, err := os.OpenFile("cotacao.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	checkErr(err)
	if _, err := f.Write([]byte(text + "\n")); err != nil {
		f.Close()
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
