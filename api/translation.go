package api

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/translate"
	"github.com/nlopes/slack"
	"golang.org/x/text/language"
)

const (
	hSignature = "X-Slack-Signature"
	hTimestamp = "X-Slack-Request-Timestamp"
)

// Hangul related constants.
const (
	HangulBase = 0xAC00
	HangulEnd  = 0xD7A4
)

var slackSigningSecret string

func isHangul(s string) bool {

	for _, char := range s {
		if HangulBase <= char && char < HangulEnd {
			return true
		}
	}
	return false
}

func init() {
	// set SLACK_SIGNING_SECRET=<Signing Secret of your Slack App> in windows cli
	slackSigningSecret = os.Getenv("SLACK_SIGNING_SECRET")
}

// TranslateEnglish2Korean is Translation function
func TranslateEnglish2Korean(w http.ResponseWriter, r *http.Request) {

	// Verify Slack Request with Signing Secret, and Timeout check
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

	// Creates a client.
	ctx := context.Background()
	client, err := translate.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// check if it's Korean
	var targetLanguage string
	if isHangul(s.Text) {
		targetLanguage = "en"
	} else {
		targetLanguage = "ko"
	}
	// Sets the target language.
	target, err := language.Parse(targetLanguage)
	if err != nil {
		log.Fatalf("failed to parse target language: %v", err)
	}

	// Translates the text into en <-> ko
	translations, err := client.Translate(ctx, []string{s.Text}, target, nil)
	if err != nil {
		log.Fatalf("failed to translate text: %v", err)
	}

	srcText := s.Text
	tgtText := translations[0].Text
	fmt.Printf("source: %v\n", srcText)
	fmt.Printf("target: %v\n", tgtText)

	// Send back to Slack
	slactPost := fmt.Sprintf("`source`: %s\n`target`: %s\n", srcText, tgtText)
	params := &slack.Msg{
		Type: "mrkdwn",
		Text: slactPost,
	}

	b, err := json.Marshal(params)
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

	// Timeout check
	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		log.Printf("failed strconv.ParseInt%v\n", err)
		return false
	}

	tSince := time.Since(time.Unix(ts, 0))
	diff := time.Duration(abs64(int64(tSince)))
	if diff > 5*time.Minute {
		log.Println("timed out")
		return false
	}

	// Not timeouted, then check Mac
	return checkMAC(message, slackSignature, slackSigningSecret)
}

func abs64(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}
