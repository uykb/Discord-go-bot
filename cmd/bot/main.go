package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"tv-bot-go/internal/ai"
	"tv-bot-go/internal/analysis"
	"tv-bot-go/internal/bot"
	"tv-bot-go/internal/binance"
	"tv-bot-go/internal/market"
	"tv-bot-go/pkg/config"

	"github.com/bwmarrin/discordgo"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		return
	}

	if cfg.DiscordToken == "" || cfg.AIAPIKey == "" || cfg.GuildID == "" {
		fmt.Println("Error: DISCORD_BOT_TOKEN, DEEPSEEK_API_KEY, and DISCORD_GUILD_ID must be set.")
		return
	}

	// Create a new Discord session
	dg, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		fmt.Println("Error creating Discord session:", err)
		return
	}

	// Initialize services
	binanceClient := binance.NewClient()
	marketService := market.NewService(binanceClient)
	analysisService := analysis.NewService()
	aiService := ai.NewService(cfg.AIAPIKey, cfg.AIEndpoint)

	// Create and start the bot
	app := bot.NewBot(dg, marketService, analysisService, aiService, cfg.GuildID)
	if err := app.Start(); err != nil {
		fmt.Println("Error starting bot:", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	fmt.Println("Shutting down bot...")
	app.Stop()
}
