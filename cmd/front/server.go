package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	tDB "github.com/tarekseba/flight-scraper/internal/api/db"
	"github.com/tarekseba/flight-scraper/internal/api/dto"
	"github.com/tarekseba/flight-scraper/internal/api/front"
	"github.com/tarekseba/flight-scraper/internal/logger"
	"github.com/tarekseba/flight-scraper/internal/scraper/types"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logger.ErrorLogger.Fatal(err)
	}

	db := tDB.InitDB()
	defer db.Close()

	var wg *sync.WaitGroup = new(sync.WaitGroup)
	queryService := front.NewQueryService(db, wg)

	var mux *http.ServeMux = http.NewServeMux()
	mux.HandleFunc("/api/", func(response http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(response, "Hello world")
	mux.HandleFunc("POST /api/query", func(response http.ResponseWriter, request *http.Request) {
		var query types.Query
		response.Header().Add("Content-Type", "application/json")
		err := json.NewDecoder(request.Body).Decode(&query)
		if err != nil {
			dto.HandleError(response, err)
			return
		}
		err = query.SanitizeQuery()
		if err != nil {
			dto.HandleError(response, err)
			return
		}
		_, err = queryService.DB.ExecContext(request.Context(), tDB.INSERT_QUERY, query.Departure, query.Destination, query.StayDuration, query.MonthHorizon, query.StringifyWeekdays())
		if err != nil {
			dto.HandleError(response, err)
			return
		}
		dto.HandleResponse(response, query, "Query added successfully")
	})

	var server http.Server = http.Server{
		Addr:                         ":8000",
		Handler:                      mux,
		DisableGeneralOptionsHandler: false,
		ErrorLog:                     logger.ErrorLogger,
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.ErrorLogger.Println(fmt.Sprintf("Server closed with error : %+v", err))
		}
	}()

	sig := <-sigChan
	logger.InfoLogger.Println(fmt.Sprintf("Signal received [%+v]", strings.ToUpper(sig.String())))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	err = server.Shutdown(ctx)
	if err != nil {
		logger.ErrorLogger.Println(err)
		os.Exit(1)
	}

	front.StopServices(ctx, &queryService)
	logger.InfoLogger.Println("Server gracefully shutdown")
}
