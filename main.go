package main

import (
	"3lab/handlers"
	"3lab/utils"
	"bufio"
	"fmt"
	"os"
)

func main() {
	db, err := utils.ConnectDB()
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}
	defer db.Close()

	reader := bufio.NewReader(os.Stdin)

	// Аутентификация пользователя
	fmt.Print("Enter username: ")
	username, _ := reader.ReadString('\n')
	username = username[:len(username)-1]

	fmt.Print("Enter password: ")
	password, _ := reader.ReadString('\n')
	password = password[:len(password)-1]

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

	// Проверка доступа
	fmt.Print("Enter resource: ")
	resource, _ := reader.ReadString('\n')
	resource = resource[:len(resource)-1]

	fmt.Print("Enter permission (read/write): ")
	permission, _ := reader.ReadString('\n')
	permission = permission[:len(permission)-1]

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
}
