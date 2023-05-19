package main

import (
	"flag"
	"index-bot/commands"
	"index-bot/db"
	"index-bot/events"
	"index-bot/logger"
	"os"

	tempest "github.com/Amatsagu/Tempest"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()

	logger.InitLogger()
	logger.Info.Println("Starting Index bot application...")

	logger.Info.Println("Connecting with database...")
	db.InitConnectionWithDatabase()

	logger.Info.Println("Creating new MarketStack Rest API controller...")
	db.InitMarketStackRestAPIController()

	logger.Info.Println("Creating new Tempest app client...")
	client := tempest.CreateClient(tempest.ClientOptions{
		ApplicationID: tempest.StringToSnowflake(os.Getenv("DISCORD_APP_ID")),
		PublicKey:     os.Getenv("DISCORD_APP_PUBLIC_KEY"),
		Token:         "Bot " + os.Getenv("DISCORD_APP_TOKEN"),
	})

	logger.Info.Println("Registering slash commands...")
	client.RegisterCommand(commands.Ping)
	client.RegisterCommand(commands.RegisterChannel)
	client.RegisterCommand(commands.Blacklist)
	client.RegisterSubCommand(commands.BlacklistAdd, "blacklist")
	client.RegisterSubCommand(commands.BlacklistRemove, "blacklist")
	client.RegisterCommand(commands.SwitchCrons)

	// logger.Info.Println("Syncing local command cache with Discord API...")
	// err := client.SyncCommands(nil, nil, false)
	// if err != nil {
	// 	logger.Error.Panic(err) // Unrecoverable!
	// }

	logger.Info.Println("Starting all cron jobs...")
	go events.InitCronJobs(&client)

	logger.Info.Println("Launching application! Waiting for incoming interactions...")
	err := client.ListenAndServe("/", os.Getenv("APP_ADDRESS"))
	if err != nil {
		logger.Error.Panic(err) // Unrecoverable!
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
