package bot

import (
	"context"
	"fmt"
	"strings"
	"time"
	"tv-bot-go/internal/ai"
	"tv-bot-go/internal/analysis"
	"tv-bot-go/internal/market"

	"github.com/bwmarrin/discordgo"
)

// Bot represents the Discord bot application.
type Bot struct {
	Session         *discordgo.Session
	MarketService   *market.Service
	AnalysisService *analysis.Service
	AIService       *ai.Service
	GuildID         string
}

// NewBot creates a new Bot instance.
func NewBot(s *discordgo.Session, marketSvc *market.Service, analysisSvc *analysis.Service, aiSvc *ai.Service) *Bot {
	return &Bot{
		Session:         s,
		MarketService:   marketSvc,
		AnalysisService: analysisSvc,
		AIService:       aiSvc,
	}
}

// Start runs the bot and registers command handlers.
func (b *Bot) Start() error {
	b.Session.AddHandler(b.ready)
	b.Session.AddHandler(b.interactionCreate)

	return b.Session.Open()
}

// Stop gracefully shuts down the bot.
func (b *Bot) Stop() {
	// Clean up registered commands
	registeredCommands, _ := b.Session.ApplicationCommands(b.Session.State.User.ID, "")
	for _, v := range registeredCommands {
		b.Session.ApplicationCommandDelete(b.Session.State.User.ID, "", v.ID)
	}
	b.Session.Close()
}

func (b *Bot) ready(s *discordgo.Session, event *discordgo.Ready) {
	fmt.Printf("Logged in as: %v#%v\n", s.State.User.Username, s.State.User.Discriminator)
	// Register the /analyze command
	cmd := &discordgo.ApplicationCommand{
		Name:        "analyze",
		Description: "Analyze a crypto symbol",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "symbol",
				Description: "Crypto symbol (e.g., BTCUSDT)",
				Required:    true,
			},
		},
	}

	_, err := s.ApplicationCommandCreate(s.State.User.ID, b.GuildID, cmd)
	if err != nil {
		fmt.Printf("Cannot create command: %v\n", err)
	}
}

func (b *Bot) interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		if i.ApplicationCommandData().Name == "analyze" {
			b.handleAnalyzeCommand(s, i)
		}
	}
}

func (b *Bot) handleAnalyzeCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Acknowledge the command immediately
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	options := i.ApplicationCommandData().Options
	symbol := strings.ToUpper(options[0].Value.(string))

	// Default to USDT pair if not specified
	if !strings.Contains(symbol, "USDT") {
		symbol += "USDT"
	}

	// 1. Fetch Market Data
	marketData, err := b.MarketService.FetchMarketData(context.Background(), symbol)
	if err != nil {
		b.sendErrorResponse(s, i.Interaction, fmt.Sprintf("Error fetching market data for %s: %s", symbol, err))
		return
	}

	// 2. Perform Technical Analysis
	analysis1h := b.AnalysisService.AnalyzeKlines(marketData.Klines1h, "1h")
	analysis15m := b.AnalysisService.AnalyzeKlines(marketData.Klines15m, "15m")

	// 3. Generate AI Analysis
	aiSummary, err := b.AIService.GenerateAnalysis(context.Background(), symbol, analysis1h, analysis15m)
	if err != nil {
		b.sendErrorResponse(s, i.Interaction, fmt.Sprintf("Error generating AI analysis for %s: %s", symbol, err))
		return
	}

	// 4. Send the result
	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("Analysis for %s", symbol),
		Description: aiSummary,
		Color:       0x0099ff, // Blue
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
	})
}

func (b *Bot) sendErrorResponse(s *discordgo.Session, i *discordgo.Interaction, errorMsg string) {
	s.InteractionResponseEdit(i, &discordgo.WebhookEdit{
		Content: &errorMsg,
	})
}
