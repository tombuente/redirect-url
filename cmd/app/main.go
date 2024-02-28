package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tombuente/redirect-url/internal/redirect"
	sql_embed "github.com/tombuente/redirect-url/sql"
)

func main() {
	db := sqlx.MustConnect("sqlite3", "data.db")

	_, err := db.ExecContext(context.TODO(), sql_embed.RedirectSQLSchema)
	if err != nil {
		slog.Error("unable to load sql schema")
		return
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	redirectRepository := redirect.NewRepository(db)
	redirectService := redirect.NewRedirectService(redirectRepository)
	redirectHandler := redirect.NewRedirectHandler(redirectService)
	redirectAPIHandler := redirect.NewRedirectAPIHandler(redirectService)

	r.Mount("/", redirectHandler)
	r.Mount("/api", redirectAPIHandler)

	fmt.Println("Running...")
	http.ListenAndServe(":8080", r)
}
