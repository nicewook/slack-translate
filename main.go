package main

import (
	"net/http"

	"github.com/nicewook/slack-translate/api"
)

func main() {
	http.HandleFunc("/", api.TranslateEnglish2Korean)
	http.ListenAndServe(":8080", nil)
}
