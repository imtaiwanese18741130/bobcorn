package main

import (
	"net/url"
	"os"
	"strings"
)

type Config struct {
	AppEnv            string
	AppHostname       string
	AppPath           string
	AppPort           string
	AppTrustedProxies []string
	AppBaseURL        *url.URL
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		AppEnv:            os.Getenv("APP_ENV"),
		AppHostname:       os.Getenv("APP_HOSTNAME"),
		AppPath:           os.Getenv("APP_PATH"),
		AppPort:           os.Getenv("APP_PORT"),
		AppTrustedProxies: strings.Split(strings.ReplaceAll(os.Getenv("APP_TRUSTED_PROXIES"), " ", ""), ","),
	}

	if !strings.HasPrefix(cfg.AppPath, "/") {
		cfg.AppPath = "/" + cfg.AppPath
	}

	scheme := "https"
	if strings.HasPrefix(cfg.AppHostname, "localhost") {
		scheme = "http"
	}

	rootPath := ""
	if cfg.AppPath == "/" {
		rootPath = cfg.AppHostname
	} else {
		rootPath = cfg.AppHostname + cfg.AppPath
	}

	var err error

	cfg.AppBaseURL, err = url.ParseRequestURI(scheme + "://" + rootPath)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
