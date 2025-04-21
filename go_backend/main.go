package main

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"expense_backend/database_sql"
	"expense_backend/middleware"
	"expense_backend/services"
	models "expense_backend/types"
	"log"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var ddl string

func run() {
	return
}

func AuthMux() http.Handler {
	authMux := http.NewServeMux()

	authMux.HandleFunc("POST /login", login)

	return http.StripPrefix("/api/v1/auth", authMux)
}

func login(w http.ResponseWriter, r *http.Request) {
	var cred models.Credentials

	cred.Id = 12313213
	err := json.NewDecoder(r.Body).Decode(&cred)

	if err != nil || cred.Username != "user" || cred.Password != "password123" {
		http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	}
	queries := database_sql.New(models.AppContextValue.GetDB())
	token, err := services.GenerateToken(queries, cred.Id, cred.Username, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(1 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
	})

	json.NewEncoder(w).Encode("✅ Login Sucessful")
}

func hello(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Hello")
}
func protected(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("You have reached a protected route")
}

func main() {
	ctx := context.Background()

	db, err := sql.Open("sqlite", "./mydb.db")
	if err != nil {
		log.Fatal(err)
		return
	}

	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		log.Fatal(err)
		return
	}

	// Initialize the global AppContextValue
	appCtx := models.NewAppContext(db)
	models.AppContextValue = appCtx
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello", hello)
	mux.Handle("/api/v1/auth/", AuthMux())
	mux.Handle("/protected", middleware.Authenticate(protected))

	log.Fatal(http.ListenAndServe(":3000", mux))
}
