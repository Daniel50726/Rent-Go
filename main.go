package main

import (
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
	Name         string `json:"name"`
	Biography    string `json:"biography"`
	Type         string `json:"type"`
	Color        string `json:"color"`
	Model        string `json:"model"`
	Brand        string `json:"brand"`
	Price        string `json:"price"`
	Assistance   bool   `json:"assistance"`
	Insurance    bool   `json:"insurance"`
	BabySeat     bool   `json:"babySeat"`
	Transmission string `json:"transmission"`
	LuxurySeat   bool   `json:"luxurySeat"`
	Fuel         string `json:"fuel"`
	Reservation  bool   `json:"reservation"`
	Images       Images `json:"images"`
}

type Images struct {
	CarImage string `json:"CarImage"`
}

var cars = []Car{
	{
		Name:         "Audi Rs 6",
		Biography:    "El Audi RS 6 es una poderosa combinación de lujo y rendimiento deportivo, diseñado para los amantes de la velocidad y el confort extremo.",
		Type:         "Sedan",
		Color:        "Blue",
		Model:        "2020",
		Brand:        "Audi",
		Price:        "$20000",
		Assistance:   true,
		Insurance:    true,
		BabySeat:     true,
		Transmission: "Automatic",
		LuxurySeat:   true,
		Fuel:         "Gas",
		Reservation:  false,
		Images: Images{
			CarImage: "https://images.pexels.com/photos/1035108/pexels-photo-1035108.jpeg?auto=compress&cs=tinysrgb&w=1260&h=750&dpr=1",
		},
	},
	{
		Name:         "Toyota Rav4",
		Biography:    "La Toyota RAV4 es una camioneta SUV compacta que combina estilo, versatilidad y rendimiento. Destaca por su diseño moderno, amplio espacio interior y eficiencia en consumo de combustible. Ideal para familias y aventuras urbanas o fuera de la carretera, la RAV4 ofrece un equilibrio entre confort y capacidad todoterreno.",
		Type:         "SUV",
		Color:        "White",
		Model:        "2019",
		Brand:        "Toyota",
		Price:        "$18000",
		Assistance:   false,
		Insurance:    true,
		BabySeat:     false,
		Transmission: "Manual",
		LuxurySeat:   false,
		Fuel:         "Gasoline",
		Reservation:  false,
		Images: Images{
			CarImage: "https://images.pexels.com/photos/2036544/pexels-photo-2036544.jpeg?auto=compress&cs=tinysrgb&w=1260&h=750&dpr=1",
		},
	},
	{
		Name:         "Camaro",
		Biography:    "El Chevrolet Camaro es un icónico muscle car estadounidense conocido por su potencia, estilo deportivo y emocionante rendimiento en carretera. Ideal para quienes buscan una experiencia de conducción dinámica y emocionante.",
		Type:         "Muscle Car",
		Color:        "Red",
		Model:        "2021",
		Brand:        "Chevrolet",
		Price:        "$22000",
		Assistance:   true,
		Insurance:    false,
		BabySeat:     true,
		Transmission: "Automatic",
		LuxurySeat:   true,
		Fuel:         "Electric",
		Reservation:  false,
		Images: Images{
			CarImage: "https://images.pexels.com/photos/3637981/pexels-photo-3637981.jpeg?auto=compress&cs=tinysrgb&w=1260&h=750&dpr=1",
		},
	},
	{
		Name:         "Bmw Serie 4",
		Biography:    "El BMW Serie 4 es una línea de automóviles deportivos de lujo que combina un diseño elegante con un rendimiento dinámico. Disponible en versiones coupé, convertible y Gran Coupé, destaca por su potente motor, confort interior de alta gama y tecnología avanzada, ofreciendo una experiencia de conducción emocionante y sofisticada.",
		Type:         "Coupe",
		Color:        "Blue",
		Model:        "2018",
		Brand:        "BMW",
		Price:        "$30000",
		Assistance:   false,
		Insurance:    true,
		BabySeat:     false,
		Transmission: "Manual",
		LuxurySeat:   true,
		Fuel:         "Hybrid",
		Reservation:  false,
		Images: Images{
			CarImage: "https://images.pexels.com/photos/898336/pexels-photo-898336.jpeg?auto=compress&cs=tinysrgb&w=1260&h=750&dpr=1",
		},
	},
	{
		Name:         "Chevrolet Tracker",
		Biography:    "Rentabiliza tu viaje con la Chevrolet Tracker, ideal para aventuras urbanas y escapadas.",
		Type:         "SUV",
		Color:        "Red",
		Model:        "2017",
		Brand:        "Chevrolet",
		Price:        "$25000",
		Assistance:   true,
		Insurance:    false,
		BabySeat:     true,
		Transmission: "Automatic",
		LuxurySeat:   true,
		Fuel:         "Hybrid",
		Reservation:  false,
		Images: Images{
			CarImage: "https://i.pinimg.com/originals/91/c0/26/91c026dd5f2233fc0d6b702603d143f7.jpg",
		},
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
	api.HandleFunc("/cars/search", searchCarsPostHandler).Methods("POST")
	api.HandleFunc("/cars/toggleReservation", toggleReservationHandler).Methods("POST")

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

func toggleReservationHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		CarName string `json:"carName"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	for i, car := range cars {
		if car.Name == request.CarName {
			cars[i].Reservation = !cars[i].Reservation
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(cars[i])
			return
		}
	}

	http.Error(w, "Car not found", http.StatusNotFound)
}
