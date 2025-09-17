package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SupabaseURL        string
	SupabaseServiceKey string
	GoogleAPIKey       string
	GoogleCX           string
}

func LoadConfig() *Config {
	log.Println("[INFO] (config.LoadConfig) Loading configuration...")
	if os.Getenv("tdtp-cron") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Println("[WARNING] (config.LoadConfig) No .env file found")
		}
	}

	return &Config{
		SupabaseURL:        getEnv("SUPABASE_URL", ""),
		SupabaseServiceKey: getEnv("SUPABASE_SERVICE_ROLE_KEY", ""),
		GoogleAPIKey:       getEnv("GOOGLE_API_KEY", ""),
		GoogleCX:           getEnv("GOOGLE_CX", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	if fallback == "" {
		log.Fatalf("[ERROR] (config.getEnv) Environment variable %s not set and no fallback provided", key)
	}
	return fallback
}
