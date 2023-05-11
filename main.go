package main

import (
	"flag"
	"index-bot/logger"
	"index-bot/other"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	loadEnv()

	logger.InitLogger()
	logger.Info.Println("Starting Index bot application...")
	rest := other.CreateMarketStackRest(os.Getenv("MARKET_API_KEY"), 1000)

	logger.Info.Println("Rendering image...")
	err := other.DrawENDLeaderboard(
		&rest,
		[]string{"IXIC.INDX", "GSPC.INDX", "DJI.INDX"},
		"./assets/us-end.png",
		"./out/test-us-end.png",
	)
	if err != nil {
		logger.Error.Panic(err)
	}

	err = other.DrawYTDLeaderboard(
		&rest,
		[]string{"IXIC.INDX", "GSPC.INDX", "DJI.INDX"},
		"./assets/us-ytd.png",
		"./out/test-us-ytd.png",
	)
	if err != nil {
		logger.Error.Panic(err)
	}

	err = other.DrawYTDLeaderboard(
		&rest,
		[]string{"CDAXX.INDX", "UKX.INDX", "FCHI.INDX", "SSE180.INDX", "HSI.INDX", "N225.INDX"},
		"./assets/international-ytd.png",
		"./out/test-international-ytd.png",
	)
	if err != nil {
		logger.Error.Panic(err)
	}
}

func loadEnv() {
	var filename string
	flag.StringVar(&filename, "env", ".env", "switches which environmental file should be used (default: .env)")
	flag.Parse()

	if err := godotenv.Load(filename, ".env"); err != nil {
		logger.Error.Panic(err) // Unrecoverable!
	}

	if _, available := os.LookupEnv("MARKET_API_KEY"); !available {
		logger.Error.Panic("environmental file is missing credentials (invalid .env file)")
	}
}
