# Slack Slash Command to Translate English to Korean

## Implementation Phase 1. Local Server which parses Slack's HTTP Post

1. Set environment variables (Windows)

  ```cmd
  set SLACK_SIGNING_SECRET=<Slack App's Signing Secret>
  ```

2. Locally, run the server

  ```cmd
  go run main.go
  ```

3. Port Forwarding, so Slack HTTP POST can reach the server

  ```cmd
  ssh -o ServerAliveInterval=60 -R 80:localhost:8080 serveo.net
  ```


## Implementation Phase 2. Local Server which receive Slack's HTTP Post and Translate it back
#### Using Google Cloud Translation API

1. Set environment variables (Windows) for Slack message verification

  ```cmd
  set SLACK_SIGNING_SECRET=<Slack App's Signing Secret>
  ```

2. Set environment variables (Windows) for Cloud Transalation API Authentication

  ```cmd
  set GOOGLE_APPLICATION_CREDENTIALS=<path to the .json file and .json file name>
  ```

3. Locally, run the server

  ```cmd
  go run main.go
  ```

4. Port Forwarding, so Slack HTTP POST can reach the server

  ```cmd
  ssh -o ServerAliveInterval=60 -R 80:localhost:8080 serveo.net
  ```
