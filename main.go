package main

import (
	"3lab/handlers"
	"3lab/utils"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	db, err := utils.ConnectDB()
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}
	defer db.Close()

	reader := bufio.NewReader(os.Stdin)

	// Выбор между входом и регистрацией
	fmt.Println("Choose an option:")
	fmt.Println("1. Login")
	fmt.Println("2. Register")
	fmt.Print("Enter option number: ")
	optionInput, _ := reader.ReadString('\n')
	optionInput = strings.TrimSpace(optionInput)

	switch optionInput {
	case "1":
		// Аутентификация пользователя
		fmt.Print("Enter username: ")
		username, _ := reader.ReadString('\n')
		username = strings.TrimSpace(username)

		fmt.Print("Enter password: ")
		password, _ := reader.ReadString('\n')
		password = strings.TrimSpace(password)

		authenticated, err := handlers.Authenticate(db, username, password)
		if err != nil {
			fmt.Println("Error authenticating user:", err)
			return
		}

		if !authenticated {
			fmt.Println("Authentication failed.")
			return
		}

		// Получение user_id для аутентифицированного пользователя
		var userID int
		err = db.QueryRow("SELECT id FROM users WHERE username = $1", username).Scan(&userID)
		if err != nil {
			fmt.Println("Error retrieving user ID:", err)
			return
		}

		// Получение роли пользователя
		role, err := handlers.GetUserRole(db, userID)
		if err != nil {
			fmt.Println("Error retrieving user role:", err)
			return
		}

		if role == nil {
			fmt.Println("User has no role assigned.")
			return
		}

		fmt.Printf("User role: %s\n", role.Name)

		// Пример логики доступа на основе роли
		switch role.Name {
		case "admin":
			fmt.Println("Welcome, admin! You have full access.")
		case "teacher":
			fmt.Println("Welcome, teacher! You have standard access.")
		case "student":
			fmt.Println("Welcome, student! You have limited access.")
		default:
			fmt.Println("Unknown role.")
			return
		}

		// Проверка доступа к ресурсу
		fmt.Print("Enter resource: ")
		resource, _ := reader.ReadString('\n')
		resource = strings.TrimSpace(resource)

		fmt.Print("Enter permission (read/write): ")
		permission, _ := reader.ReadString('\n')
		permission = strings.TrimSpace(permission)

		hasAccess, err := handlers.CheckAccess(db, userID, resource, permission)
		if err != nil {
			fmt.Println("Error checking access:", err)
			return
		}

		if hasAccess {
			fmt.Println("Access granted!")
		} else {
			fmt.Println("Access denied.")
		}

	case "2":
		// Регистрация нового пользователя
		fmt.Print("Enter new username: ")
		newUsername, _ := reader.ReadString('\n')
		newUsername = strings.TrimSpace(newUsername)

		fmt.Print("Enter new password: ")
		newPassword, _ := reader.ReadString('\n')
		newPassword = strings.TrimSpace(newPassword)

		// Выбор роли
		fmt.Println("Choose a role:")
		fmt.Println("1. Student")
		fmt.Println("2. Teacher")
		fmt.Print("Enter role number: ")
		roleInput, _ := reader.ReadString('\n')
		roleInput = strings.TrimSpace(roleInput)

		var roleName string
		switch roleInput {
		case "1":
			roleName = "student"
		case "2":
			roleName = "teacher"
		default:
			fmt.Println("Invalid role choice.")
			return
		}

		err := handlers.Register(db, newUsername, newPassword, roleName)
		if err != nil {
			fmt.Println("Error registering user:", err)
			return
		}

		fmt.Println("Registration successful!")

	default:
		fmt.Println("Invalid option.")
	}
}
