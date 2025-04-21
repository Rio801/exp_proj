package services

import (
	"context"
	"expense_backend/database_sql"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var secretKey = []byte("thesupersecretkey")
var ctx = context.Background()

func GenerateToken(queries *database_sql.Queries, uID uint, name string, w http.ResponseWriter) (string, error) {
	jti := uuid.New().String()

	expirationTime := time.Now().Add(1 * time.Hour)
	claims := jwt.MapClaims{
		"user_id":   uID,
		"username":  name,
		"expires":   expirationTime.Unix(),
		"jti":       jti,
		"issued_at": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	_, err := queries.StoreToken(ctx, database_sql.StoreTokenParams{UserID: int64(uID), Jti: jti, ExpiresAt: expirationTime})

	if err != nil {

		http.Error(w, "Invalid Token", http.StatusUnauthorized)
		return "", err
	}
	return token.SignedString(secretKey)
}

func VerifyToken(appCtx context.Context, queries *database_sql.Queries, tokenString string, w http.ResponseWriter) (jwt.MapClaims, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error with Signing Method")
		}

		return secretKey, nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Invalid Token", http.StatusUnauthorized)
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		jti, ok := claims["jti"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid jti in token claims")
		}
		isRevoked, err := queries.IsRevoked(appCtx, jti)
		if err != nil {
			return nil, err
		}
		if isRevoked.Valid && isRevoked.Bool {
			return nil, fmt.Errorf("token has been revoked")
		}

		// Check expiration
		if exp, ok := claims["expires"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return nil, fmt.Errorf("token has expired")
			}
		}

		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}

func RevokeToken(ctx context.Context, db database_sql.Queries, jti string) error {
	_, err := db.RevokeToken(ctx, jti)
	if err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}
	return nil
}
