package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/tarekseba/flight-scraper/internal/api/front"
	"github.com/tarekseba/flight-scraper/internal/logger"
)

func main() {
	var mux *http.ServeMux = http.NewServeMux()
	mux.HandleFunc("/api/", func(response http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(response, "Hello world")
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
		sig := <-sigChan

		logger.InfoLogger.Println(fmt.Sprintf("Signal received '%+v'", strings.ToUpper(sig.String())))
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
		defer cancel()
		server.Shutdown(ctx)
		// add services here later
		front.StopServices(ctx)
		logger.InfoLogger.Println("Server gracefully shutdown")
	}()

	err := server.ListenAndServe()

	if err != nil {
		logger.ErrorLogger.Println(err)
		os.Exit(1)
	}
}
