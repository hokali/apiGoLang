package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"./routes"
	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"
)

func main() {
	// Get ENV and set variables
	e := godotenv.Load(os.ExpandEnv("./.env"))
	if e != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println(e)

	port := os.Getenv("PORT")

	// Handle routes
	r := routes.Handlers()
	http.Handle("/", r)

	// Serve
	log.Printf("Server up on port '%s'", port)
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(headers, methods, origins)(r)))
}
