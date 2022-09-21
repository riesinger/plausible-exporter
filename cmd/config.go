package main

import (
	"fmt"
	"net/url"

	"github.com/spf13/viper"
)

var (
	listenAddress string
	plausibleHost *url.URL
	token         string
	siteIDs       []string
)

func readConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/plausible-exporter/")

	viper.SetDefault("listen_address", "0.0.0.0:8080")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("config: failed to read config file: %w", err)
		}
	}

	listenAddress = viper.GetString("listen_address")
	plausibleHostRaw := viper.GetString("plausible_host")
	if plausibleHostRaw == "" {
		return fmt.Errorf("config: no plausible host provided")
	}
	var err error
	plausibleHost, err = url.Parse(plausibleHostRaw)
	if err != nil {
		return fmt.Errorf("config: cannot parse plausible host as URL: %w", err)
	}
	token = viper.GetString("plausible_token")
	if token == "" {
		return fmt.Errorf("config: no plausible token provided")
	}
	siteIDs = viper.GetStringSlice("plausible_site_ids")
	if len(siteIDs) == 0 {
		return fmt.Errorf("config: no plausible site IDs provided")
	}

	return nil
}
