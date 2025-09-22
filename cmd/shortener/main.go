package main

import (
	"log"
	"net/http"
	"os"

	"github.com/yonatannn111/URL_Shortener_With_Go/internal/api"
	"github.com/yonatannn111/URL_Shortener_With_Go/internal/storage"
	"github.com/yonatannn111/URL_Shortener_With_Go/internal/worker"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/go-redis/redis/v8"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	redisAddr := os.Getenv("REDIS_ADDR")
	baseURL := os.Getenv("BASE_URL")
	if dbURL == "" || baseURL == "" || redisAddr == "" {
		log.Fatal("set DATABASE_URL, REDIS_ADDR, BASE_URL in env")
	}

	// connect to Postgres
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalf("db: %v", err)
	}

	// connect to Redis
	rdb := redis.NewClient(&redis.Options{Addr: redisAddr})

	// init store + worker
	store := storage.NewStore(db, rdb)
	worker := worker.NewWorker(store)
	worker.Start()

	// init app
	app := &api.App{Store: store, Worker: worker, BaseURL: baseURL}

	// routes
	r := chi.NewRouter()
	r.Post("/shorten", app.ShortenHandler)
	r.Get("/{code}", app.RedirectHandler)
	r.Get("/{code}/stats", app.StatsHandler)

	log.Println("listening on :8080")
	http.ListenAndServe(":8080", r)
}
