package worker

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/BlueNyang/theday-theplace-cron/pkg/config"
	"github.com/BlueNyang/theday-theplace-cron/pkg/domain/common"
	"github.com/BlueNyang/theday-theplace-cron/pkg/parser/instances"
)

func Worker(ctx context.Context, cfg *config.Config, wg *sync.WaitGroup, jobs chan instances.Job, result chan<- *common.Exhibition) {
	log.Println("[INFO] (worker.Worker) Worker started")
	// jobs 채널이 닫힐 때까지 계속해서 작업을 처리
	// 채널이 닫히면 for 루프가 종료되고 고루틴이 종료됨
	for job := range jobs {
		// URL에서 ?(쿼리문) 이전까지의 문자열을 키로 사용
		log.Printf("[INFO] (worker.Worker) Worker received job %+v\n", job)
		domain := strings.SplitN(*job.Url, "?", 2)[0]
		parser, err := instances.GetParser(domain)
		if err != nil {
			log.Printf("[ERROR] (worker.Worker) GetParser error: %+v\n", err)
			wg.Done()
			continue
		}

		parseResult, err := parser.Parsing(ctx, cfg, job)
		if err != nil {
			log.Printf("[ERROR] (worker.Worker) Parsing error: %+v\n", err)
			wg.Done()
			continue
		}

		fmt.Printf("[INFO] (worker.Worker) parseResult: %+v\n", parseResult)

		for _, newJob := range parseResult.DiscoveredJobs {
			wg.Add(1)
			jobs <- *newJob
		}

		for _, exhibition := range parseResult.FoundExhibitions {
			result <- exhibition
		}

		wg.Done()
	}
	log.Println("[INFO] (worker.Worker) Worker exiting")
}
