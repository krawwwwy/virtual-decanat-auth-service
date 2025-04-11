package handlers

import (
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// Функция для регистрации нового пользователя
func Register(db *sql.DB, username, password, roleName string) error {
	// Хеширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Начало транзакции
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Вставка нового пользователя в базу данных
	query := "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id"
	var userID int
	err = tx.QueryRow(query, username, string(hashedPassword)).Scan(&userID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert user: %w", err)
	}

	// Получение ID роли
	var roleID int
	err = tx.QueryRow("SELECT id FROM roles WHERE name = $1", roleName).Scan(&roleID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get role ID: %w", err)
	}

	// Вставка записи в таблицу user_roles
	query = "INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)"
	_, err = tx.Exec(query, userID, roleID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert user role: %w", err)
	}

	// Фиксация транзакции
	return tx.Commit()
}

// Функция для аутентификации пользователя
func Authenticate(db *sql.DB, username, password string) (bool, error) {
	var storedPassword string
	query := "SELECT password FROM users WHERE username = $1"
	err := db.QueryRow(query, username).Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	if err != nil {
		return false, nil
	}
	return true, nil
}
