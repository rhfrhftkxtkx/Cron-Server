package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SupabaseDBURL string
	GoogleAPIKey  string
	GoogleCX      string
	GeminiAPIKey  string
}

func LoadConfig() *Config {
	if os.Getenv("tdtp-cron") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Println("[Warning] No .env file found")
		}
	}

	return &Config{
		SupabaseDBURL: getEnv("SUPABASE_DB_URL", ""),
		GoogleAPIKey:  getEnv("GOOGLE_API_KEY", ""),
		GoogleCX:      getEnv("GOOGLE_CX", ""),
		GeminiAPIKey:  getEnv("GEMINI_API_KEY", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	if fallback == "" {
		log.Fatalf("[Error] Environment variable %s not set and no fallback provided", key)
	}
	return fallback
}
