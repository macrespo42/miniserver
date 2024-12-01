package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/macrespo42/miniserver/internal/database"
)

func GetDbConnection() *database.Queries {
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Error when connection to database")
	}
	return database.New(db)
}

func main() {
	apiCfg := ApiConfig{}
	mux := http.NewServeMux()
	rootHandler := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", apiCfg.MiddlewareMetricsInc(rootHandler))

	mux.HandleFunc("GET /app/assets/logo.png", HandlerImage)

	mux.HandleFunc("GET /api/healthz", HandlerHealth)
	mux.HandleFunc("POST /api/validate_chirp", HandlerValidateChirp)
	mux.HandleFunc("POST /api/users", HandleCreateUser)

	mux.HandleFunc("POST /admin/reset", apiCfg.HandleReset)
	mux.HandleFunc("GET /admin/metrics", apiCfg.HandleServerHits)

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	log.Fatal(server.ListenAndServe())
}
