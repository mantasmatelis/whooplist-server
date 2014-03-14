package whooplist

import (
	"bytes"
	"net/http"
	"strings"
)

var slackUrl = "https://whooplist.slack.com/services/hooks/incoming-webhook?token=GsURXK3zks9I7T3JhKBTOW0m"

func SlackPostError(msg string) {
	msg = strings.Replace(msg, "\n", "\\n", -1)
	_, _ = http.Post(slackUrl, "text/json", bytes.NewBufferString("{\"text\": \""+msg+"\"}"))
}
