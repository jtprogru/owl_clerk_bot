package tg

type Config struct {
	BotToken string `yaml:"bot_token" env:"BOT_TOKEN"`
	IsDebug  bool   `yaml:"is_debug" env:"DEBUG" env-default:"false"`
}
