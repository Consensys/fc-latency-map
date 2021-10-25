package webhook

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type NotifierImpl struct {
	Conf *viper.Viper
}

type Payload struct {
	Datetime  string   `json:"datetime"`
	Filenames []string `json:"filenames"`
}

func NewNotifier(conf *viper.Viper) Notifier {
	return &NotifierImpl{
		Conf: conf,
	}
}

func (n *NotifierImpl) Notify(files *[]string) bool {
	if files == nil || len(*files) == 0 {
		log.Println("No updates to notify")
		return false
	}
	urls := strings.Split(n.Conf.GetString("WEBHOOK_NOTIFY_URLS"), ",")
	body, err := json.Marshal(Payload{
		Datetime:  time.Now().String(),
		Filenames: *files,
	})
	if err != nil {
		log.Errorf("Error: %s\n", err)
		return false
	}
	for _, url := range urls {
		log.Printf("Request POST: url=%s, body=%s\n", url, string(body))
		// nolint // G107: Potential HTTP request made with variable url
		resp, err := http.Post(
			url,
			"application/json; charset=utf-8",
			bytes.NewBuffer(body),
		)
		if err != nil {
			log.Errorf("Error: %s", err)
			continue
		}
		log.Println("Response status:", resp.Status)
		log.Println("Notification sent!")
	}
	return true
}
