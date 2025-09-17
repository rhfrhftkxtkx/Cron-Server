package database

import (
	"log"

	"github.com/BlueNyang/theday-theplace-cron/pkg/config"
	"github.com/BlueNyang/theday-theplace-cron/pkg/domain/common"
	"github.com/supabase-community/supabase-go"
)

type SupabaseClient struct {
	Client *supabase.Client
}

func InitSupabase(cfg *config.Config) SupabaseClient {
	log.Println("[INFO] (supabase.InitSupabase) Initializing Supabase Client")
	client, err := supabase.NewClient(cfg.SupabaseURL, cfg.SupabaseServiceKey, nil)
	if err != nil {
		log.Fatalf("[ERROR] (database.SaveExhibitions) Failed to create Supabase client: %v", err)
	}

	log.Println("[INFO] (supabase.InitSupabase) Supabase client created")
	return SupabaseClient{Client: client}
}

func (s *SupabaseClient) GetCrawlTargets() ([]common.CrawlTarget, error) {
	var targets []common.CrawlTarget

	_, err := s.Client.From("crawl_targets").
		Select("*", "", false).
		ExecuteTo(&targets)
	if err != nil {
		log.Println("[ERROR] (database.SaveExhibitions) Failed to fetch crawl targets:", err)
		return nil, err
	}

	return targets, nil
}

func SaveExhibitions(client *supabase.Client, exhibitions []*common.Exhibition) error {
	if exhibitions == nil {
		log.Println("[INFO] (database.SaveExhibitions) No exhibitions to save")
		return nil
	}

	execute, i, err := client.From("exhibitions").
		Upsert(exhibitions, "exhibition_id", "", "").
		Execute()
	if err != nil {
		log.Println("[ERROR] (database.SaveExhibitions) Failed to insert exhibitions:", err)
		return err
	}

	log.Printf("[INFO] (database.SaveExhibitions) Exhibitions saved: %v, %d rows affected\n", execute, i)
	return nil
}
