package scenarios

import (
	"context"
	"fmt"

	"github.com/tarekseba/flight-scraper/internal/logger"
	"github.com/tarekseba/flight-scraper/internal/scraper/types"
)

type ScenarioAction func(context.Context) error

func LogScenario(s types.Scenario) ScenarioAction {
	return func(ctx context.Context) error {
		logger.InfoLogger.Println(fmt.Sprintf("Scenario [%s] starting", s.Name()))
		err := s.Do(ctx)
		if err == nil {
			logger.InfoLogger.Println(fmt.Sprintf("Scenario [%s] done", s.Name()))
		}
		return err
	}
}
