module main

go 1.12

require (
	github.com/lusis/slack-test v0.0.0-20190426140909-c40012f20018 // indirect
	github.com/nicewook/slack-translate/api v0.1.0
	github.com/nlopes/slack v0.5.0
	github.com/stretchr/testify v1.3.0 // indirect
)

replace github.com/nicewook/slack-translate/api v0.1.0 => ./api
