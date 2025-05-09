package handlers

import (
	"auth-service/models"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func HandleRegister(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			log.Println("Error parsing form:", err)
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		user := models.User{
			Username: r.FormValue("username"),
			Password: r.FormValue("password"),
			Role:     r.FormValue("role"),
		}

		log.Printf("Register request: %+v\n", user)

		err = Register(db, user.Username, user.Password, user.Role)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Registration successful!"))
	}
}

func HandleLogin(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			log.Println("Error parsing form:", err)
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		user := models.User{
			Username: r.FormValue("username"),
			Password: r.FormValue("password"),
			Role:     "-",
		}

		log.Printf("Login request: %+v\n", user)

		authenticated, err := Authenticate(db, user.Username, user.Password)
		if err != nil {
			log.Println("Error authenticating:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if authenticated {
			w.Write([]byte("Login successful!"))
		} else {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		}
	}
}

func Register(db *sql.DB, username, password, roleName string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	var userID int
	err = tx.QueryRow("INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id", username, string(hashedPassword)).Scan(&userID)
	if err != nil {
		log.Println("Error inserting user:", err)
		tx.Rollback()
		return err
	}

	var roleID int
	err = tx.QueryRow("SELECT id FROM roles WHERE name = $1", roleName).Scan(&roleID)
	if err != nil {
		log.Println("Error selecting role:", err)
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)", userID, roleID)
	if err != nil {
		log.Println("Error inserting user role:", err)
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func Authenticate(db *sql.DB, username, password string) (bool, error) {
	var storedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = $1", username).Scan(&storedPassword)
	if err != nil {
		log.Printf("Cant get data to login as %v:", username, err)
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	return err == nil, nil
}
