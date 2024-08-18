package main

import (
	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"log"
	"surge/cmd"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf(color.YellowString("Failed to load .env: %s"), err)
	}

	var rootCommand = cmd.BuildRootCommand()
	if err := rootCommand.Execute(); err != nil {
		log.Fatal(err)
	}
}
