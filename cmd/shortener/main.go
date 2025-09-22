package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/go-redis/redis/v8"
	"github.com/yonatannn111/URL_Shortener_With-_Go/storage"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	redisAddr := os.Getenv("REDIS_ADDR")
	baseURL := os.Getenv("BASE_URL")
	if dbURL == "" || baseURL == "" || redisAddr == "" {
		log.Fatal("set DATABASE_URL, REDIS_ADDR, BASE_URL in env")
	}
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil { log.Fatalf("db: %v", err) }
	rdb := redis.NewClient(&redis.Options{Addr: redisAddr})

	store := storage.NewStore(db, rdb)
	worker := worker.NewWorker(store)
	worker.Start()

	app := &api.App{Store: store, Worker: worker, BaseURL: baseURL}

	r := chi.NewRouter()
	r.Post("/shorten", app.ShortenHandler)
	r.Get("/{code}", app.RedirectHandler)
	r.Get("/{code}/stats", app.StatsHandler)

	log.Println("listening on :8080")
	http.ListenAndServe(":8080", r)
}
