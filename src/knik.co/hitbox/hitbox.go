package main

import (
	"net/http"
	"fmt"
	"net/url"
	"strings"
	"log"
)

func main() {
	config := LoadConfig("config.json")

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		query := strings.Split(r.URL.String(), "/ping?")[1]
		values, err := url.ParseQuery(query)

		if err != nil {
			fmt.Errorf("Error parsing query: %s", err.Error())
		}

		ip := getIP(r)
		locale := getLocale(config, ip)
		page := getPage(config, &values)
		referrer := getReferrer(config, &values)

		sendMessage(config, page, referrer, locale)

		log.Printf("%s%s%s%s", page, referrer, locale, ip)
	})

	log.Printf("Hitbox listening on port %s...", config.GetPort())
	http.ListenAndServe(config.GetPort(), nil)
}
