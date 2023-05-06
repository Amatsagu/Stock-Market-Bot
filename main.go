package main

import (
	"flag"
	"index-bot/logger"
	"os"

	godotenv "github.com/joho/godotenv"
)

func main() {
	logger.InitLogger()
	logger.Info.Println("Starting Index bot application...")
}

func loadEnv() {
	var filename string
	flag.StringVar(&filename, "env", ".env", "switches which environmental file should be used (default: .env)")
	flag.Parse()

	if err := godotenv.Load(filename, ".env"); err != nil {
		logger.Error.Panic(err) // Unrecoverable!
	}

	if _, available := os.LookupEnv("APP_TOKEN"); !available {
		logger.Error.Panic("environmental file is missing credentials (invalid .env file)")
	}
}
