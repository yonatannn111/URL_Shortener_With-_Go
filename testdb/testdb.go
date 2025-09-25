package main

import (
    "fmt"
    "os"
    "github.com/joho/godotenv"
    "database/sql"
    _ "github.com/lib/pq"
)

func main() {
    // Load .env in the current folder
    if err := godotenv.Load(); err != nil {
        fmt.Println("⚠️ No .env file found, using system env vars")
    }

    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        fmt.Println("❌ DATABASE_URL not set in environment")
        return
    }

    db, err := sql.Open("postgres", dsn)
    if err != nil {
        panic(err)
    }
    defer db.Close()

    if err := db.Ping(); err != nil {
        panic(err)
    }

    fmt.Println("✅ Connected to Supabase Postgres!")
}
