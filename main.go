package main

import (
	"fmt"
	"os"

	"github.com/mercadolibre/minesweeper/server"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
		fmt.Printf("Fixed port to 5000")
	}
	fmt.Println("Listening port: " + port)
	server.New().Run(":" + port)
}
