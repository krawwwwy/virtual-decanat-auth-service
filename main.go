package main

import (
	"3lab/handlers"
	"3lab/utils"
	"bufio"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"os"
)

func main() {
	// Пароль, который вы хотите захешировать
	passwords := "testpassword"

	// Генерация хеша пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwords), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error generating hash:", err)
		return
	}

	// Вывод сгенерированного хеша
	fmt.Println("Hashed Password:", string(hashedPassword))

	// Проверка хеша
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(passwords))
	if err != nil {
		fmt.Println("Password does not match")
	} else {
		fmt.Println("Password matches")
	}
	db, err := utils.ConnectDB()
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}
	defer db.Close()

	reader := bufio.NewReader(os.Stdin)
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

	if authenticated {
		fmt.Println("Authentication successful!")
	} else {
		fmt.Println("Authentication failed.")
	}
}
