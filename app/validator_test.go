package main

import (
	"log"
	"testing"

	"github.com/valyala/fasthttp"
)

func BenchmarkValidator(b *testing.B) {

	list := make(map[string][]IPItem)
	var err error

	if list, err = readLines("bot.txt"); err != nil {
		log.Printf("[ERROR] file with ips error: %+v", err)
	}

	app := &application{
		HTTP: &fasthttp.Server{
			Name: "IP-CHECKER",
		},
		IPList: list,
	}

	for i := 0; i < b.N; i++ {
		app.validateIP("87.117.185.192")
		// perform the operation we're analyzing
	}
}
