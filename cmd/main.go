package main

import (
	"context"
	"log"
	"net/http"
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

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	plausibleClients := make(map[string]*plausible.Client)
	for _, siteID := range siteIDs {
		plausibleClients[siteID] = &plausible.Client{HostAPIBase: plausibleHost, SiteID: siteID, Token: token}
	}

	metrics := prometheus.NewServer(siteIDs)

	updatePlausibleData := func() {
		for _, siteID := range siteIDs {
			data, err := plausibleClients[siteID].GetTimeseriesData()
			if err != nil {
				log.Printf("Refreshing data for site %s failed: %v", siteID, err)
				return
			}
			metrics.UpdateDataForSite(siteID, data)
			log.Printf("Data for site %s was refreshed from plausible", siteID)
		}
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
	if bearerAuthToken != "" {
		srv.SetBearerAuthToken(bearerAuthToken)
	}

	go func() {
		if err := srv.ListenAndServe(listenAddress); err != nil && err != http.ErrServerClosed {
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
