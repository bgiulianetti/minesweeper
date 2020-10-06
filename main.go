package main

import "github.com/mercadolibre/minesweeper/server"

func main() {
	server.New().Run(":5000")
}
