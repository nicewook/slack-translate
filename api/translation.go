package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/nlopes/slack"
)

var slackSigningSecret string

func init() {
	slackSigningSecret = os.Getenv("SLACK_SIGNING_SECRET") // set SLACK_SIGNING_SECRET=<Signing Secret of your Slack App> in windows cli
	fmt.Printf("signing Secrect. type: %T, value: %s\n", slackSigningSecret, slackSigningSecret)
}

// TranslateEnglish2Korean is Translation function
func TranslateEnglish2Korean(w http.ResponseWriter, r *http.Request) {

	// 1. verify signed secret
	// verifier has signing secret and signature in r.Header
	verifier, err := slack.NewSecretsVerifier(r.Header, slackSigningSecret)
	if err != nil {
		log.Printf("fail to create verifier: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// verifier.Write(r.Body) and returns r.Body back (means not consuming r.Body)
	r.Body = ioutil.NopCloser(io.TeeReader(r.Body, &verifier))

	s, err := slack.SlashCommandParse(r)
	if err != nil {
		log.Printf("fail to parse slash command from Slack: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// now we have all ingredient for Verifing, so let's verify
	if err = verifier.Ensure(); err != nil {
		log.Printf("fail to authorize: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	json, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		log.Println(err)
	}
	slackPost := fmt.Sprintf("HTTP POST from Slack:\n--\n%s\n", string(json))
	fmt.Println(slackPost)
	fmt.Fprint(w, slackPost)
}
