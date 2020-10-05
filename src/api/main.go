package main

import "github.com/mercadolibre/minesweeper/src/api/server"

func main() {
	server.New().Run(":8080")
}
