package main

import (
	"os"
	"encoding/json"
	"strings"
	"fmt"
	"io/ioutil"
)

type Config struct {
	SiteName   string `json:"site_name"`
	WebhookUrl string `json:"slack_webhook_url"`
	GeoIpUrl   string `json:"geoip_url_format"`
	IgnoreList []string `json:"ignore_referrer_domains"`
	Port       string `json:"port"`
}

func LoadConfig(filepath string) *Config {
	configFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Errorf("%s\n", err.Error())
		panic("Failed to load config.json!")
	}

	c := Config{}

	err = json.Unmarshal(configFile, &c)
	if err != nil {
		fmt.Errorf("%s\n", err.Error())
		panic("Failed to load config.json!")
	}

	return &c
}

func (c *Config) Site() string {
	return c.SiteName
}

func (c *Config) Webhook() string {
	return c.WebhookUrl
}

func (c *Config) GeoIp() string {
	return c.GeoIpUrl
}

func (c *Config) IsIgnored(domain string) bool {
	for _, ignored := range c.IgnoreList {

		// Because of sub-domains and the like, let's check
		// to make sure that the ignored domain isn't anywhere
		// inside the tested domain, not just that they're equal
		if strings.Index(domain, ignored) != -1 {
			return true
		}
	}

	return false
}

func (c *Config) GetPort() string {
	// First check if a port was defined in config.json,
	// then check for an environment variable
	if len(c.Port) > 0 {
		if strings.Index(c.Port, ":") != 0 {
			return fmt.Sprintf(":%s", c.Port)
		} else {
			return c.Port
		}
	} else if port := os.Getenv("PORT"); len(port) > 0 {
		if strings.Index(port, ":") != 0 {
			return fmt.Sprintf(":%s", port)
		} else {
			return port
		}
	}

	panic("No port defined through either config.json or $PORT!")
}
