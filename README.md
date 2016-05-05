# slack-fm
Radio web app powered by Slack channel recommendations

# Installing & Running

Set up your PostgreSQL database using the instructions in setup.sql

Copy settings.go.example to settings.go, update your db settings (if necessary); you'll need to add a Slack API key, and SLACK_CHANNEL_ID is the channel ID of the Slack channel on which recommendations are being made.

: `go install && slack-fm`

Slack FM will now be running on port 8080. If you load /, all media posts in the channel will be loaded into your recommendations database.
