package front

import (
	"context"
	"sync"

	"github.com/tarekseba/flight-scraper/internal/logger"
)

type StoppableService interface {
	Stop(ctx context.Context)
}

type DefaultStopper struct {
	Sync sync.WaitGroup
}

func (s *DefaultStopper) Stop(ctx context.Context) {
	waitChannel := make(chan struct{})
	go func() {
		s.Sync.Wait()
		close(waitChannel)
	}()
	select {
	case <-ctx.Done():
		{
			logger.WarnLogger.Println("Context finished sooner then QueryService was done")
		}
	case <-waitChannel:
		{
			logger.InfoLogger.Println("QueryService gracefully shutdown")
		}
	}
}

func StopServices(ctx context.Context, s ...StoppableService) {
	for idx := range s {
		s[idx].Stop(ctx)
	}
}
