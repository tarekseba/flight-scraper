package front

import (
	"context"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/tarekseba/flight-scraper/internal/scraper/types"
)

type QueryService struct {
	DB   *sqlx.DB
	Sync *sync.WaitGroup
	DefaultStopper
}

func (s *QueryService) InsertQuery(con context.Context, query types.Query) error {
	s.Sync.Wait()
	return nil
}

func NewQueryService(db *sqlx.DB, sync *sync.WaitGroup) QueryService {
	return QueryService{
		DB:             db,
		Sync:           sync,
		DefaultStopper: DefaultStopper{},
	}
}
