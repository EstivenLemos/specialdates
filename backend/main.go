package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"specialdates-backend/internal/db"
	"specialdates-backend/internal/handlers"
	"specialdates-backend/internal/repository"
	"specialdates-backend/internal/service"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN no definida. Ej: user:pass@tcp(localhost:3306)/specialdates?parseTime=true")
	}

	sqlDB, err := db.Connect(dsn)
	if err != nil {
		log.Fatalf("error conectando DB: %v", err)
	}
	defer sqlDB.Close()

	repo := repository.NewMySQLRepository(sqlDB)
	svc := service.NewDateService(repo)
	h := handlers.NewHandler(svc)

	r := mux.NewRouter()

	// CORS middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if req.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, req)
		})
	})

	api := r.PathPrefix("/api").Subrouter()
	h.RegisterRoutes(api)

	// start server
	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Servidor backend escuchando en :8080")
	log.Fatal(srv.ListenAndServe())
}
