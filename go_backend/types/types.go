package models

import (
	"context"
	"database/sql"
	"expense_backend/database_sql"
)


type Credentials struct {
	Id       uint   `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type AppContext struct {
    DB      *sql.DB
    Queries *database_sql.Queries
    Ctx     context.Context  // Changed from ctx to Ctx
}

func NewAppContext(db *sql.DB) *AppContext {
    queries := database_sql.New(db)
    return &AppContext{
        DB:      db,
        Queries: queries,
        Ctx:     context.Background(),
    }
}

// GetDB returns the database connection
func (a *AppContext) GetDB() *sql.DB {
    return a.DB
}

var AppContextValue *AppContext
