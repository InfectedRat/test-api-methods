package main

import (
	"fmt"
	"log"
	"net/http"
	logic "test-api-methods/internal/app"
)

func main() {

	db := logic.ConnectDB()
	defer db.Close()

	logic.CreateTables(db)

	http.HandleFunc("/country", logic.GetCountriesHandler)

	fmt.Println("Сервер запущен на порту 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
