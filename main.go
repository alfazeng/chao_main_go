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

	r := mux.NewRouter()

	// Rutas públicas
	r.HandleFunc("/api/register", RegisterHandler).Methods("POST")
	r.HandleFunc("/api/login", LoginHandler).Methods("POST")

	// NUEVA RUTA PROTEGIDA
	// Usamos .Handle() en lugar de .HandleFunc() para poder envolverlo en el middleware
	r.Handle("/api/profile", AuthMiddleware(http.HandlerFunc(ProfileHandler))).Methods("GET")
	// Configuración de CORS para permitir peticiones desde tu frontend en React
// EN backend/main.go

c := cors.New(cors.Options{
    AllowedOrigins: []string{
        "http://localhost:3000", // La dejas para desarrollo local
        "https://chaotravelapp.com", // <- AÑADE TU URL DE NETLIFY AQUÍ
    },
    AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowedHeaders: []string{"Content-Type", "Authorization"},
})

	handler := c.Handler(r)

	log.Println("Servidor escuchando en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}