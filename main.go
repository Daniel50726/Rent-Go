package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

var db *sqlx.DB

type User struct {
	Username string `db:"username"`
	Password string `db:"password_hash"`
}

func main() {
	// Conectar a PostgreSQL
	dsn := "user=postgres password=mysecretpassword dbname=postgres sslmode=disable host=localhost port=5432"
	var err error
	db, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
        username TEXT PRIMARY KEY,
        password_hash TEXT
    )`)
	if err != nil {
		log.Fatalln(err)
	}

	r := mux.NewRouter()

	// Rutas para el registro y el login con prefijo /api
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/register", registerHandler).Methods("POST")
	api.HandleFunc("/login", loginHandler).Methods("POST")

	// Servir archivos estáticos
	staticDir := "./frontend/"
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(staticDir))))

	log.Println("Server listening on :8080")
	http.ListenAndServe(":8080", r)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Could not hash the password", http.StatusInternalServerError)
		return
	}

	query := "INSERT INTO users (username, password_hash) VALUES ($1, $2)"
	_, err = db.Exec(query, user.Username, string(hashedPassword))
	if err != nil {
		http.Error(w, "Could not register the user", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "User %s registered successfully", user.Username)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials User
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var user User
	err := db.Get(&user, "SELECT * FROM users WHERE username = $1", credentials.Username)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(w, "User %s logged in successfully", user.Username)
}
