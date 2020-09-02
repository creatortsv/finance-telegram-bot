package env

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	token string
	port  string
)

func init() {
	flag.StringVar(&token, "token", "", "Telegram bot token")
	flag.StringVar(&port, "port", "", "Port of the Telegram bot")
}

// InitSettings ...
func InitSettings() {
	flag.Parse()

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}

	for name, value := range map[string]string{
		"TOKEN": token,
		"PORT":  port,
	} {
		if len(value) > 0 {
			os.Setenv(name, value)
		}
	}
}

// Get ...
func Get(key string, defaultValue ...string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	for _, val := range defaultValue {
		return val
	}

	return ""
}
