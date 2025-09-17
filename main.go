package main

import (
	"context"
	"log"
	"sync"

	"github.com/BlueNyang/theday-theplace-cron/pkg/config"
	"github.com/BlueNyang/theday-theplace-cron/pkg/database"
	"github.com/BlueNyang/theday-theplace-cron/pkg/domain/common"
	"github.com/BlueNyang/theday-theplace-cron/pkg/parser"
	"github.com/BlueNyang/theday-theplace-cron/pkg/parser/instances"
	"github.com/BlueNyang/theday-theplace-cron/pkg/worker"
)

func main() {
	cfg := config.LoadConfig()
	ctx := context.Background()

	FirstTierDataProcessing(cfg, ctx)
}

func FirstTierDataProcessing(cfg *config.Config, ctx context.Context) []*common.Exhibition {
	parser.InitializeParsers()

	jobsChan := make(chan instances.Job, 100)
	resultChan := make(chan *common.Exhibition, 100)
	var wg sync.WaitGroup

	log.Println("[INFO] Starting worker pool")
	numWorkers := 10
	for i := 0; i < numWorkers; i++ {
		go worker.Worker(ctx, cfg, &wg, jobsChan, resultChan)
	}

	supabase := database.InitSupabase(cfg)
	crawlTargets, err := supabase.GetCrawlTargets()
	if err != nil {
		log.Fatalf("[ERROR] Failed to get crawl targets: %v", err)
	}
	log.Printf("[INFO] Fetched %d crawl targets", len(crawlTargets))

	for _, target := range crawlTargets {
		log.Printf("[INFO] Enqueuing job for URL: %s", target.URL)
		wg.Add(1)

		// Code for Only Testing
		//if target.URL != "https://www.museum.go.kr/MUSEUM/contents/M0207000000.do" {
		//	wg.Done()
		//	continue
		//}

		jobsChan <- instances.Job{
			Url:   &target.URL,
			Depth: 1,
		}
	}

	go func() {
		wg.Wait()
		close(jobsChan)
	}()

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var allExhibitions []*common.Exhibition
	for resp := range resultChan {
		allExhibitions = append(allExhibitions, resp)
	}

	err = database.SaveExhibitions(supabase.Client, allExhibitions)
	if err != nil {
		log.Fatalf("[ERROR] Failed to save exhibitions: %v", err)
	}

	return allExhibitions
}
