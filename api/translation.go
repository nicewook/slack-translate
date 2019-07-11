package api

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/nlopes/slack"
)

const (
	hSignature = "X-Slack-Signature"
	hTimestamp = "X-Slack-Request-Timestamp"
)

var slackSigningSecret string

func init() {
	// set SLACK_SIGNING_SECRET=<Signing Secret of your Slack App> in windows cli
	slackSigningSecret = os.Getenv("SLACK_SIGNING_SECRET")
}

// TranslateEnglish2Korean is Translation function
func TranslateEnglish2Korean(w http.ResponseWriter, r *http.Request) {

	// Verify Slack Request with Signing Secret
	if ok := verifySlackSignature(r, []byte(slackSigningSecret)); ok == false {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("failed on VerifyRequest()")
		return
	}

	// Parsing Slack Slash Command
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to parse Slack Slash Command: %v", err)
		return
	}

	// Print what Slack sent
	b, err := json.MarshalIndent(s, "", "--")
	if err != nil {
		log.Printf("failed to MarchalIndent: %v", err)
	}
	slackPost := fmt.Sprintf("HTTP POST from Slack\n%s\n", string(b))
	fmt.Println(slackPost)

	// Send back to Slack
	params := &slack.Msg{Text: slackPost}
	b, err = json.Marshal(params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

// checkMAC reports whether messageMAC is a valid HMAC tag for message.
func checkMAC(message, receivedMAC string, slackSigningSecret []byte) bool {
	mac := hmac.New(sha256.New, slackSigningSecret)
	if _, err := mac.Write([]byte(message)); err != nil {
		log.Printf("mac.Write(%v) failed\n", message)
		return false
	}
	calculatedMAC := "v0=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(calculatedMAC), []byte(receivedMAC))
}

// VerifySlackSignature verifies the request is coming from Slack
// Read https://api.slack.com/docs/verifying-requests-from-slack
func verifySlackSignature(r *http.Request, slackSigningSecret []byte) bool {
	if r.Body == nil {
		return false
	}

	// do not consume req.body
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	// prepare message for signing
	timestamp := r.Header.Get(hTimestamp)
	slackSignature := r.Header.Get(hSignature)
	message := "v0:" + timestamp + ":" + string(bodyBytes)

	return checkMAC(message, slackSignature, slackSigningSecret)
}
