// cmd/server/main.go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/go-chi/cors"
    "github.com/joho/godotenv"

    "github.com/PKR9759/LiftGo-backend/internal/db"
)

func main() {
    // load .env
    if err := godotenv.Load(); err != nil {
        log.Println("no .env file found, reading from environment")
    }

    ctx := context.Background()

    // connect to DB
    pool, err := db.Connect(ctx)
    if err != nil {
        log.Fatalf("db connect: %v", err)
    }
    defer pool.Close()
    log.Println("database connected")

    // run migrations
    if err := db.RunMigrations(ctx, pool); err != nil {
        log.Fatalf("migrations: %v", err)
    }
    log.Println("migrations complete")

    // router
    r := chi.NewRouter()

    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(middleware.RequestID)
    r.Use(cors.Handler(cors.Options{
        AllowedOrigins:   []string{"http://localhost:3000"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
        AllowCredentials: true,
    }))

    // health check
    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("ok"))
    })

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    fmt.Printf("LiftGo API running on :%s\n", port)
    log.Fatal(http.ListenAndServe(":"+port, r))
}