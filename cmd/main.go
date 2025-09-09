package cmd

import (
	"log"
	"sync"

	"github.com/BlueNyang/theday-theplace-cron/internal/config"
	"github.com/BlueNyang/theday-theplace-cron/internal/crawler"
	"github.com/BlueNyang/theday-theplace-cron/internal/database"
	"github.com/BlueNyang/theday-theplace-cron/internal/searcher"
	"github.com/supabase-community/postgrest-go"
)

func main() {
	cfg := config.LoadConfig()

	supabaseClient := database.InitSupabase(cfg.SupabaseURL, cfg.SupabaseServiceRoleKey)
	job(cfg, supabaseClient)
}

func job(cfg *config.Config, supabaseClient *postgrest.Client) {
	log.Println("[INFO] Starting job")

	//ctx := context.Background()

	query := ""
	urls, err := searcher.SearchGoogle(cfg.GoogleAPIKey, cfg.GoogleCX, query, "h12")
	if err != nil {
		log.Printf("[ERROR] Google search error: %v", err)
		return
	}
	if len(urls) == 0 {
		log.Printf("[INFO] No new results found. Ending job.")
		return
	}

	var wg sync.WaitGroup
	contentChan := make(chan crawler.PageContent, len(urls))

	for _, url := range urls {
		wg.Add(1)
		go crawler.CrawlPage(url, &wg, contentChan)
	}

	wg.Wait()
	close(contentChan)

	var pagesToProcess []crawler.PageContent
	for content := range contentChan {
		if content.Error != nil {
			log.Printf("[ERROR] Failed to process page: %v", content.Error)
		} else {
			pagesToProcess = append(pagesToProcess, content)
		}
	}

	// Gemini processing.
}
