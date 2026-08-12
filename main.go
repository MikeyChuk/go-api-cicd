package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var db *sql.DB

func main() {
	var err error

	db, err = connectDB()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	http.HandleFunc("/health", healthHandler)

	log.Println("API listening on port 8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func connectDB() (*sql.DB, error) {
	host := getRequiredEnv("DB_HOST")
	port := getEnv("DB_PORT", "5432")
	user := getRequiredEnv("DB_USER")
	password := getRequiredEnv("DB_PASSWORD")
	database := getEnv("DB_NAME", "postgres")
	sslMode := getEnv("DB_SSLMODE", "require")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host,
		port,
		user,
		password,
		database,
		sslMode,
	)

	databaseHandle, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := databaseHandle.PingContext(ctx); err != nil {
		databaseHandle.Close()
		return nil, err
	}

	log.Println("Connected to PostgreSQL successfully")

	return databaseHandle, nil
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		http.Error(w, "database unavailable", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func getRequiredEnv(key string) string {
	value := os.Getenv(key)

	if value == "" {
		log.Fatalf("Required environment variable %s is missing", key)
	}

	return value
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// THIS CODE DOES NOT REQUIRE A CONNECTION TO A DATABASE. IT IS A SIMPLE API THAT RETURNS A HEALTH CHECK AND A LIST OF USERS.//
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// package main

// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"
// )

// type User struct {
// 	Name  string `json:"name"`
// 	Email string `json:"email"`
// }

// var users = []User{
// 	{
// 		Name:  "MichaelC",
// 		Email: "michael.Chidube@example.com",
// 	},
// }

// func main() {

// 	http.HandleFunc("/health", healthHandler)
// 	http.HandleFunc("/users", usersHandler)

// 	log.Println("API listening on port 8080")

// 	if err := http.ListenAndServe(":8080", nil); err != nil {
// 		log.Fatalf("Server failed: %v", err)
// 	}
// }

// func healthHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodGet {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	writeJSON(w, http.StatusOK, map[string]string{
// 		"status":  "ok",
// 		"service": "stacklaunch-go-api",
// 	})
// }

// func usersHandler(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case http.MethodGet:
// 		listUsersHandler(w)
// 	case http.MethodPost:
// 		createUserHandler(w, r)
// 	default:
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 	}
// }

// func listUsersHandler(w http.ResponseWriter) {
// 	writeJSON(w, http.StatusOK, users)
// }

// func createUserHandler(w http.ResponseWriter, r *http.Request) {
// 	var user User

// 	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
// 		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
// 		return
// 	}

// 	if user.Name == "" || user.Email == "" {
// 		http.Error(w, "Name and email are required", http.StatusBadRequest)
// 		return
// 	}

// 	users = append(users, user)

// 	writeJSON(w, http.StatusCreated, user)
// }

// func writeJSON(w http.ResponseWriter, status int, value any) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(status)

// 	if err := json.NewEncoder(w).Encode(value); err != nil {
// 		log.Printf("Unable to encode response: %v", err)
// 	}
// }
