package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Inicializa la conexión con la base de datos
	InitDB()

	// Crea un nuevo enrutador
	r := mux.NewRouter()

	// Define los endpoints de la API
	r.HandleFunc("/api/register", RegisterHandler).Methods("POST")
	r.HandleFunc("/api/login", LoginHandler).Methods("POST")

	// Configuración de CORS para permitir peticiones desde tu frontend en React
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"}, // Cambia esto en producción
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	handler := c.Handler(r)

	log.Println("Servidor escuchando en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}