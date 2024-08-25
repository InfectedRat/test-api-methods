package logic

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
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

func ConnectDB() *sql.DB {

	db, err := sql.Open("sqlite3", "/Users/maximbabichev/Library/DBeaverData/workspace6/.metadata/sample-database-sqlite-1/Chinook.db")
	if err != nil {
		log.Fatalf("Не удалось подключиться к БД: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("База данных недоступна: %v", err)
	}
	return db

}

func CreateTables(db *sql.DB) {
	querys := map[string]string{"countries": `CREATE TABLE IF NOT EXISTS countries (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		alfaTwo TEXT NOT NULL,
		alfaThree TEXT NOT NULL,
		name TEXT NOT NULL,
		nameBrief TEXT NOT NULL);`,
		"accounts": `CREATE TABLE IF NOT EXISTS accounts (
		id TEXT PRIMARY KEY,
		type TEXT NOT NULL,
		name TEXT NOT NULL,
		status TEXT NOT NULL,
		openedDate TIMESTAMP NOT NULL,
		accessLevel TEXT NOT NULL
	);`}

	for name, query := range querys {
		_, err := db.Exec(query)
		if err != nil {
			log.Fatalf("Ошибка создания запроса: %v", err)
		}
		log.Printf("Таблица %s создана", name)
	}
}

func GetCountriesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}

	url := "https://sandbox-invest-public-api.tinkoff.ru/rest/tinkoff.public.invest.api.contract.v1.InstrumentsService/GetCountries"

	err := godotenv.Load("/Users/maximbabichev/ProjectGolang/test-api-methods/configs/.env")
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка загрузки токена: %v", err), http.StatusInternalServerError)
		return
	}

	token := os.Getenv("API_TOKEN")

	requestBody, err := json.Marshal(map[string]interface{}{})
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при создании тела запроса: %v", err), http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при создании запроса: %v", err), http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "application/json")

	client := http.Client{}

	response, err := client.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при выполнении запроса: %v", err), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при чтении ответа: %v", err), http.StatusInternalServerError)
		return
	}

	var resp ResponseCountry
	err = json.Unmarshal(body, &resp)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка при разборе ответа: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

}
