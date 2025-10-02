package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	
	"github.com/google/uuid" 
	"golang.org/x/crypto/bcrypt"
)

type RegisterPayload struct {
	FullName string `json:"fullName"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Country  string `json:"country"`
	Phone    string `json:"phone"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	user := User{
		FullName:     payload.FullName,
		Email:        payload.Email,
		PasswordHash: string(hashedPassword),
		Country:      payload.Country,
		Phone:        payload.Phone,
	}

	_, err = db.Exec(`
		INSERT INTO users (full_name, email, password_hash, country, phone)
		VALUES ($1, $2, $3, $4, $5)`,
		user.FullName, user.Email, user.PasswordHash, user.Country, user.Phone,
	)

	if err != nil {
		http.Error(w, "Email already exists", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var payload LoginPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var user User
	err := db.QueryRow("SELECT id, email, password_hash FROM users WHERE email = $1", payload.Email).Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(payload.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := GenerateJWT(user.ID)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}


// NUEVA FUNCIÓN: ProfileHandler
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Obtenemos el userID del contexto que el middleware añadió
	userID, ok := r.Context().Value("userID").(uuid.UUID)
	if !ok {
		http.Error(w, "Could not retrieve user ID", http.StatusInternalServerError)
		return
	}

	var user User
	// Buscamos al usuario en la BD usando el ID del token
	err := db.QueryRow(`
		SELECT id, full_name, email, country, phone, created_at
		FROM users WHERE id = $1`, userID).Scan(
		&user.ID, &user.FullName, &user.Email, &user.Country, &user.Phone, &user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}