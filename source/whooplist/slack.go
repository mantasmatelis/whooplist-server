package whooplist

import (
	"bytes"
	"log"
	"net/http"
	"strings"
)

var slackUrl = "https://whooplist.slack.com/services/hooks/incoming-webhook?token=GsURXK3zks9I7T3JhKBTOW0m"

func SlackPostError(msg string) {
	log.Print("{\"text\": \"" + msg + "\"}")
	msg = strings.Replace(msg, "\n", "\\n", -1)
	_, _ = http.Post(slackUrl, "text/json", bytes.NewBufferString("{\"text\": \""+msg+"\"}"))
}
