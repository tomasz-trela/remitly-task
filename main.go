package main

import (
	"fmt"
	"os"

	app "github.com/tomasz-trela/remitly-task/cmd/server"
	"github.com/tomasz-trela/remitly-task/internal/seeders"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a command: run or seed")
		return
	}
	switch os.Args[1] {
	case "run":
		app.Run()
	case "seed":
		seeders.SeedBanks()
	default:
		fmt.Println("Unknown command:", os.Args[1])
	}
}
