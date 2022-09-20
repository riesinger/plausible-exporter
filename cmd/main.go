package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/riesinger/plausible-exporter/plausible"
	"github.com/riesinger/plausible-exporter/prometheus"
	"github.com/riesinger/plausible-exporter/server"
)

func main() {
	if err := readConfig(); err != nil {
		log.Fatal(err)
	}

	log.Println("Starting plausible-exporter")

	ctx := context.Background()
	ctx, _ = signal.NotifyContext(ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	plausibleClient := &plausible.Client{
		HostAPIBase: plausibleHost,
		SiteID:      siteID,
		Token:       token,
	}

	metrics := prometheus.NewServer(siteID)

	updatePlausibleData := func() {
		data, err := plausibleClient.GetTimeseriesData()
		if err != nil {
			log.Printf("Refreshing data failed: %v", err)
			return
		}
		metrics.UpdateData(data)
		log.Println("Data was refreshed from plausible")
	}

	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		updatePlausibleData()
		for {
			select {
			case <-ticker.C:
				updatePlausibleData()
			case <-ctx.Done():
				log.Println("Stopping plausible refresh timer")
				return
			}
		}
	}()

	srv := server.New()
	go func() {
		if err := srv.ListenAndServe(listenAddress); err != nil {
			log.Fatalf("server: %v\n", err)
		}
	}()

	log.Println("Server started")
	<-ctx.Done()
	log.Println("Shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Graceful shutdown failed: %v\n", err)
	}
	log.Println("Bye")
}
