package conf

import (
	"github.com/sirupsen/logrus"
	"os"
)

type Config struct {
	botToken string `yaml:"bot_token" env:"BOT_TOKEN"`
	isDebug  bool   `yaml:"is_debug" env:"DEBUG" env-default:"false"`
	logger   *logrus.Logger
}

func (c *Config) GetBotToken() string {
	return c.botToken

}

func (c *Config) GetIsDebug() bool {
	return c.isDebug
}

func (c *Config) GetLogger() *logrus.Logger {
	return c.logger
}

func New() (conf *Config) {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		DisableColors:   true,
		TimestampFormat: "2006-01-02T15:04:05",
		FullTimestamp:   true,
	}
	logger.Out = os.Stdout
	logger.SetReportCaller(true)

	conf.botToken = os.Getenv("OWL_BOT_TOKEN")
	//conf.isDebug = os.Getenv("DEBUG")
	conf.logger = logger

	return conf
}
