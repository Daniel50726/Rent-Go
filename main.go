package main

import (
	"encoding/base64"
	"encoding/json"
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

type Car struct {
	Name         string    `json:"name"`
	Biography    Biography `json:"biography"`
	Type         string    `json:"type"`
	Color        string    `json:"color"`
	Model        string    `json:"model"`
	Brand        string    `json:"brand"`
	Price        string    `json:"price"`
	Assistance   bool      `json:"assistance"`
	Insurance    bool      `json:"insurance"`
	BabySeat     bool      `json:"babySeat"`
	Transmission string    `json:"transmission"`
	LuxurySeat   bool      `json:"luxurySeat"`
	Fuel         string    `json:"fuel"`
	Reservation  bool      `json:"reservation"`
	Images       Images    `json:"images"`
}

type Biography struct {
	DescriptionCar string `json:"descriptionCar"`
}

type Images struct {
	CarImage string `json:"CarImage"`
}

var cars = []Car{
	{
		Name:         "Car 1",
		Biography:    Biography{DescriptionCar: "Description for Car 1"},
		Type:         "SUV",
		Color:        "Red",
		Model:        "2020",
		Brand:        "Brand A",
		Price:        "$20000",
		Assistance:   true,
		Insurance:    true,
		BabySeat:     true,
		Transmission: "Automatic",
		LuxurySeat:   true,
		Fuel:         "Gas",
		Reservation:  true,
		Images:       Images{CarImage: base64.StdEncoding.EncodeToString([]byte("imageData1"))},
	},
	{
		Name:         "Car 2",
		Biography:    Biography{DescriptionCar: "Description for Car 2"},
		Type:         "Sedan",
		Color:        "Blue",
		Model:        "2019",
		Brand:        "Brand B",
		Price:        "$18000",
		Assistance:   false,
		Insurance:    true,
		BabySeat:     false,
		Transmission: "Manual",
		LuxurySeat:   false,
		Fuel:         "Gasoline",
		Reservation:  false,
		Images:       Images{CarImage: base64.StdEncoding.EncodeToString([]byte("imageData2"))},
	},
	{
		Name:         "Car 3",
		Biography:    Biography{DescriptionCar: "Description for Car 3"},
		Type:         "Hatchback",
		Color:        "Green",
		Model:        "2021",
		Brand:        "Brand C",
		Price:        "$22000",
		Assistance:   true,
		Insurance:    false,
		BabySeat:     true,
		Transmission: "Automatic",
		LuxurySeat:   true,
		Fuel:         "Electric",
		Reservation:  true,
		Images:       Images{CarImage: base64.StdEncoding.EncodeToString([]byte("imageData3"))},
	},
	{
		Name:         "Car 4",
		Biography:    Biography{DescriptionCar: "Description for Car 4"},
		Type:         "Convertible",
		Color:        "Black",
		Model:        "2018",
		Brand:        "Brand D",
		Price:        "$30000",
		Assistance:   false,
		Insurance:    true,
		BabySeat:     false,
		Transmission: "Manual",
		LuxurySeat:   false,
		Fuel:         "Hybrid",
		Reservation:  false,
		Images:       Images{CarImage: base64.StdEncoding.EncodeToString([]byte("imageData4"))},
	},
	{
		Name:         "Car 5",
		Biography:    Biography{DescriptionCar: "Description for Car 5"},
		Type:         "Coupe",
		Color:        "White",
		Model:        "2022",
		Brand:        "Brand E",
		Price:        "$25000",
		Assistance:   true,
		Insurance:    false,
		BabySeat:     true,
		Transmission: "Automatic",
		LuxurySeat:   true,
		Fuel:         "Hybrid",
		Reservation:  true,
		Images:       Images{CarImage: base64.StdEncoding.EncodeToString([]byte("imageData5"))},
	},
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

	// Ruta para la búsqueda de vehículos usando POST
	api.HandleFunc("/cars/search", searchCarsPostHandler).Methods("POST")

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

func searchCarsPostHandler(w http.ResponseWriter, r *http.Request) {
	var searchCriteria map[string]string
	if err := json.NewDecoder(r.Body).Decode(&searchCriteria); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var results []Car
	for _, car := range cars {
		matches := true

		for key, value := range searchCriteria {
			switch key {
			case "type":
				if car.Type != value {
					matches = false
				}
			case "color":
				if car.Color != value {
					matches = false
				}
			case "model":
				if car.Model != value {
					matches = false
				}
			case "brand":
				if car.Brand != value {
					matches = false
				}
			case "transmission":
				if car.Transmission != value {
					matches = false
				}
			case "fuel":
				if car.Fuel != value {
					matches = false
				}
				// Add more search criteria as needed
			}
		}

		if matches {
			results = append(results, car)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
