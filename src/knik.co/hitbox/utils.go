package main

import (
	"net/url"
	"fmt"
	"strings"
	"net/http"
	"net"
	"bytes"
	"io/ioutil"
	"encoding/json"
)

func getPage(config *Config, values *url.Values) string {
	if _, exists := (*values)["p"]; exists {
		page := (*values)["p"][0]

		return fmt.Sprintf("%s%s ", config.Site(), page)
	}

	// Return an empty string if no page was provided
	return ""
}

func getReferrer(config *Config, values *url.Values) string {
	if _, exists := (*values)["r"]; exists {
		s := (*values)["r"][0]
		if strings.Index(s, "//") >= 0 {
			s = strings.Split(s, "//")[1]
		}

		if strings.Index(s, "/") >= 0 {
			s = strings.Split(s, "/")[0]
		}

		// Check if the referrer has been explicitly ignored
		if !config.IsIgnored(s) {
			return fmt.Sprintf("%s ", s)
		}
	}

	// Return an empty string if an empty referrer was provided
	return ""
}

func getIP(r *http.Request) string {
	if ipProxy := r.Header.Get("X-FORWARDED-FOR"); len(ipProxy) > 0 {
		return ipProxy
	}

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	return ip
}

func getLocale(config *Config, ip string) string {
	resp, err := http.Get(fmt.Sprintf(config.GeoIp(), ip))
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Errorf("Error querying GeoIP (%s): %s", config.GeoIp(), err.Error())
	}

	localeInfo := map[string]string{}

	json.Unmarshal(body, &localeInfo)

	country := localeInfo["country_name"]
	region := localeInfo["region_name"]
	city := localeInfo["city"]

	var locale string

	if len(country) > 0 {
		locale = fmt.Sprintf("%s", country)
		if len(region) > 0 {
			locale = fmt.Sprintf("%s, %s", region, locale)
		}

		if len(city) > 0 {
			locale = fmt.Sprintf("%s, %s", city, locale)
		}
	}

	return fmt.Sprintf("%s ", locale)
}

func sendMessage(config *Config, page, referrer, locale string) {
	message := []byte(fmt.Sprintf("{\"text\":\"%s%s%s\"}", page, referrer, locale))
	req, err := http.NewRequest("POST", config.Webhook(), bytes.NewBuffer(message))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		fmt.Errorf("Failed to POST message to Slack: %s", err.Error())
	}
}
