package middleware

import (
	"context"
	"expense_backend/database_sql"
	"expense_backend/services"
	models "expense_backend/types"
	"log"
	"net/http"
)

func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		// Use the AppContextValue instead of creating a new empty AppContext
		queries := database_sql.New(models.AppContextValue.GetDB())
		cookie, err := r.Cookie("token")

		if err != nil {
			http.Error(w, "Missing Auth token", http.StatusUnauthorized)
			return
		}

		tokenString := cookie.Value

		_, err = services.VerifyToken(ctx, queries, tokenString, w)
		if err != nil {
			log.Printf("Error ->", err.Error())
			http.Error(w, "Invalid Authentication", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}
