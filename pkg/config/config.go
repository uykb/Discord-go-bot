package config

import (
	"github.com/spf13/viper"
)

// Config stores all configuration for the application.
// The values are read by viper from a config file or environment variables.
type Config struct {
	DiscordToken string `mapstructure:"DISCORD_BOT_TOKEN"`
	GuildID      string `mapstructure:"DISCORD_GUILD_ID"`
	AIAPIKey     string `mapstructure:"DEEPSEEK_API_KEY"`
	AIEndpoint   string `mapstructure:"AI_ENDPOINT"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env") // or yaml, json etc.

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Config file was found but another error was produced
			return
		}
		// Config file not found; ignore error if desired
	}

	err = viper.Unmarshal(&config)
	return
}
