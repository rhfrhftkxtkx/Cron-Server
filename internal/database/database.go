package database

import (
	"log"

	"github.com/BlueNyang/theday-theplace-cron/internal/domain/common"
	"github.com/supabase-community/postgrest-go"
)

func InitSupabase(url, key string) *postgrest.Client {
	client := postgrest.NewClient(url, "public", nil)
	client.SetAuthToken(key)

	log.Println("Supabase client created")
	return client
}

func SaveExhibitions(client *postgrest.Client, exhibitions []*common.Exhibition) error {
	if len(exhibitions) == 0 {
		log.Println("No exhibitions to save")
		return nil
	}

	var insertResult []common.Exhibition

	_, err := client.From("exhibitions").
		Upsert(exhibitions, "exhibition_id", "replace", "").
		ExecuteTo(&insertResult)
	if err != nil {
		log.Println("[Error] Failed to insert exhibitions:", err)
		return err
	}

	savedCount := len(insertResult)
	skippedCount := len(exhibitions) - savedCount

	log.Println("Exhibitions saved:", savedCount, "skipped:", skippedCount)
	return nil
}
