package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var db *sqlx.DB

type User struct {
	Username string `db:"username"`
	Password string `db:"password_hash"`
}

func main() {
	// Reemplaza "your_postgres_container_ip" con la IP del contenedor de PostgreSQL
	dsn := "user=postgres password=mysecretpassword dbname=postgres sslmode=disable host=localhost port=5432"
	//dsn := "user=postgres password=mysecretpassword dbname=postgres sslmode=disable host=some-postgres port=5432"
	var err error
	db, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	// Crear la tabla de usuarios
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
        username TEXT PRIMARY KEY,
        password_hash TEXT
    )`)
	if err != nil {
		log.Fatalln(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/register", registerHandler).Methods("POST")
	r.HandleFunc("/login", loginHandler).Methods("POST")

	log.Println("Server listening on :8080")
	http.ListenAndServe(":8080", r)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Could not hash the password", http.StatusInternalServerError)
		return
	}

	query := "INSERT INTO users (username, password_hash) VALUES ($1, $2)"
	_, err = db.Exec(query, username, string(hashedPassword))
	if err != nil {
		http.Error(w, "Could not register the user", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "User %s registered successfully", username)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	var user User
	err := db.Get(&user, "SELECT * FROM users WHERE username = $1", username)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(w, "User %s logged in successfully", username)
}
