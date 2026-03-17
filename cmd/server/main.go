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

	"github.com/PKR9759/LiftGo-backend/internal/auth"
	"github.com/PKR9759/LiftGo-backend/internal/db"
	"github.com/PKR9759/LiftGo-backend/internal/user"
	"github.com/PKR9759/LiftGo-backend/internal/ride"
	"github.com/PKR9759/LiftGo-backend/internal/booking"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, reading from environment")
	}

	ctx := context.Background()

	pool, err := db.Connect(ctx)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer pool.Close()
	log.Println("database connected")

	if err := db.RunMigrations(ctx, pool); err != nil {
		log.Fatalf("migrations: %v", err)
	}
	log.Println("migrations complete")

	// auth
	authService := auth.NewService(pool)
	authHandler := auth.NewHandler(authService)

	//user
	userRepo := user.NewRepository(pool)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	//ride
	rideRepo    := ride.NewRepository(pool)
	rideService := ride.NewService(rideRepo)
	rideHandler := ride.NewHandler(rideService)

	// booking
	bookingRepo    := booking.NewRepository(pool)
	bookingService := booking.NewService(bookingRepo)
	bookingHandler := booking.NewHandler(bookingService)

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

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	// auth routes
	r.Route("/api/auth", func(r chi.Router) {// r inside callback func is sub router of main router, so all routes defined here will be prefixed with /api/auth , like api/auth/register and api/auth/login
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})

	// protected user routes
	r.Group(func(r chi.Router) {
		r.Use(auth.RequireAuth)

		r.Route("/api/users", func(r chi.Router) {
			r.Get("/me", userHandler.GetMe)
			r.Put("/me", userHandler.UpdateMe)
		})
	})

	// rides (mixed — some public, some protected)
	r.Route("/api/rides", func(r chi.Router) {
		// public
		r.Get("/",     rideHandler.Search)
		r.Get("/{id}", rideHandler.GetByID)

		// protected
		r.Group(func(r chi.Router) {
			r.Use(auth.RequireAuth)
			r.Post("/",       rideHandler.Create)
			r.Get("/mine",    rideHandler.GetMine)
			r.Put("/{id}",    rideHandler.Update)
			r.Delete("/{id}", rideHandler.Cancel)
		})
	})

	// bookings (all protected)
	r.Route("/api/bookings", func(r chi.Router) {
		r.Use(auth.RequireAuth)
		r.Post("/",              bookingHandler.Create)
		r.Get("/mine",           bookingHandler.GetMine)
		r.Get("/incoming",       bookingHandler.GetIncoming)
		r.Get("/{id}",           bookingHandler.GetByID)
		r.Put("/{id}/confirm",   bookingHandler.Confirm)
		r.Put("/{id}/cancel",    bookingHandler.Cancel)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("LiftGo API running on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}