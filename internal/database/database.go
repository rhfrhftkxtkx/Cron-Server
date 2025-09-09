package database

import (
	"context"
	"log"

	"github.com/BlueNyang/theday-theplace-cron/internal/gemini"
	"github.com/supabase-community/postgrest-go"
)

func InitSupabase(url, key string) *postgrest.Client {
	client := postgrest.NewClient(url, "public", nil)
	client.SetAuthToken(key)

	log.Println("Supabase client created")
	return client
}

func SaveArticles(ctx context.Context, client *postgrest.Client, articles []*gemini.ProcessedData) error {
	savaedCount := 0
	skippedCount := 0

	for _, article := range articles {
		var result struct {
			Count int
		}
		_, err := client.From("articles").
			Select("count", "exact", false).
			Eq("url", article.URL).
			ExecuteTo(&result)

		if err != nil {
			log.Println("[Error] Failed to check existing article:", err)
			continue
		}

		if result.Count > 0 {
			skippedCount++
			continue
		}

		var insertResult []gemini.ProcessedData
		_, err = client.From("articles").
			Insert(article, false, "", "", "").
			ExecuteTo(&insertResult)
		if err != nil {
			log.Println("[Error] Failed to insert article:", err)
			continue
		}
		savaedCount++
	}

	log.Printf("Articles saved: %d, skipped: %d\n", savaedCount, skippedCount)
	return nil
}
