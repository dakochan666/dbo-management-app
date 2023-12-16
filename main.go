package main

import (
	"dbo-management-app/router"
	"fmt"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	r := router.StartApp()

	r.Run((":8080"))
}
