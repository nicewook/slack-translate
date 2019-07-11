package main

import (
	"log"
	"net/http"

	"github.com/nicewook/slack-translate/api"
)

func main() {
	log.Println("Server started...")
	http.HandleFunc("/", api.TranslateEnglish2Korean)
	http.ListenAndServe(":8080", nil)
}
