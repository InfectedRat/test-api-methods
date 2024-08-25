package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	_ "github.com/mattn/go-sqlite3"
)

// Структура для страны
type Country struct {
	AlfaTwo   string `json:"alfaTwo"`
	AlfaThree string `json:"alfaThree"`
	Name      string `json:"name"`
	NameBrief string `json:"nameBrief"`
}

// Структура для ответа
type ResponseCountry struct {
	Countries []Country `json:"countries"`
}

func getCountry() (*ResponseCountry, error) {

	url := "https://sandbox-invest-public-api.tinkoff.ru/rest/tinkoff.public.invest.api.contract.v1.InstrumentsService/GetCountries"

	err := godotenv.Load("/Users/maximbabichev/ProjectGolang/test-api-methods/configs/.env")
	if err != nil {
		log.Fatalf("Ошибка загрузки токена: %v", err)
	}

	token := os.Getenv("API_TOKEN")

	requestBody, err := json.Marshal(map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("Ошибка при создании тела запроса %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "application/json")

	client := http.Client{}

	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %v", err)
	}

	defer response.Body.Close()

	// Чтение ответа
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении ответа: %v", err)
	}

	var resp ResponseCountry
	err = json.Unmarshal(body, &resp)

	if err != nil {
		return nil, fmt.Errorf("ошибка при разборе ответа: %v", err)
	}

	return &resp, nil

}

func main() {

	http.HandleFunc("/country", getCountry())

}
