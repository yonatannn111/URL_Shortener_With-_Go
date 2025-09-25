package main

import (
    "log"
    "net/http"
    "os"

    "github.com/joho/godotenv"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/cors"
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
    "github.com/go-redis/redis/v8"

    "github.com/yonatannn111/URL_Shortener_With_Go/internal/storage"
    "github.com/yonatannn111/URL_Shortener_With_Go/internal/api"
    "github.com/yonatannn111/URL_Shortener_With_Go/internal/worker"
)

func main() {
    // Load .env file
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using system env vars")
    }

    dbURL := os.Getenv("DATABASE_URL")
    redisAddr := os.Getenv("REDIS_ADDR")
    baseURL := os.Getenv("BASE_URL")

    if dbURL == "" || redisAddr == "" || baseURL == "" {
        log.Fatal("set DATABASE_URL, REDIS_ADDR, BASE_URL in env")
    }

    // Connect to Postgres
    db, err := sqlx.Connect("postgres", dbURL)
    if err != nil {
        log.Fatalf("db: %v", err)
    }

    // Connect to Redis
    _ = redis.NewClient(&redis.Options{Addr: redisAddr})

    // Initialize store and worker
    store := storage.NewStore(db)
    w := worker.NewWorker(store)
    w.Start()

    // Initialize API app
    app := &api.App{Store: store, Worker: w, BaseURL: baseURL}

    // Router
    r := chi.NewRouter()

    // === CORS Middleware ===
    r.Use(cors.Handler(cors.Options{
        AllowedOrigins:   []string{"*"}, // allow all origins for dev
        AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
        AllowCredentials: true,
        MaxAge:           300,
    }))

    // Handle preflight requests
    r.Options("/*", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    })

    // Routes
    r.Post("/shorten", app.ShortenHandler)
    r.Get("/{code}", app.RedirectHandler)
    r.Get("/{code}/stats", app.StatsHandler)

    log.Println("listening on :8080")
    http.ListenAndServe(":8080", r)
}
