package front

import (
	"context"
	"sync"

	"github.com/tarekseba/flight-scraper/internal/scraper/types"
)

type QueryService struct {
	Sync sync.WaitGroup
	DefaultStopper
}

func (s *QueryService) InsertQuery(con context.Context, query types.Query) error {
	s.Sync.Wait()
	return nil
}
