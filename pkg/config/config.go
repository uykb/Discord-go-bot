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
	// Bind to environment variables. These keys match the `mapstructure` tags.
	viper.BindEnv("DISCORD_BOT_TOKEN")
	viper.BindEnv("DEEPSEEK_API_KEY")
	viper.SetDefault("AI_ENDPOINT", "https://api.deepseek.com/chat/completions")
	viper.BindEnv("AI_ENDPOINT")

	// If a path is provided (for local dev), also read from a config file.
	// Environment variables will take precedence.
	if path != "" {
		viper.AddConfigPath(path)
		viper.SetConfigName("app")
		viper.SetConfigType("env")

		if errRead := viper.ReadInConfig(); errRead != nil {
			if _, ok := errRead.(viper.ConfigFileNotFoundError); !ok {
				// A file was found but there was a different error.
				return config, errRead
			}
		}
	}

	// Unmarshal all settings into the Config struct.
	err = viper.Unmarshal(&config)
	return
}
