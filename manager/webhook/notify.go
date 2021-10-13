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

type Notifier struct {
	Conf *viper.Viper
}

type Payload struct {
	Datetime  string   `json:datetime`
	Filenames []string `json:filenames`
}

func NewNotifier(conf *viper.Viper) *Notifier {
	return &Notifier{
		Conf: conf,
	}
}

func (n *Notifier) Notify(files *[]string) bool {
	if files == nil || len(*files) == 0 {
		log.Error("Error: files is nil\n")
		return false
	}
	urls := strings.Split(n.Conf.GetString("WEBHOOK_NOTIFY_URLS"), ",")
	payload := Payload{
		Datetime:  time.Now().String(),
		Filenames: *files,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		log.Errorf("Error: %s\n", err)
		return false
	}
	for _, url := range urls {
		log.Printf("Send POST: url=%s, body=%s\n", url, payload)
		_, err = http.Post(
			url,
			"application/json; charset=utf-8",
			bytes.NewBuffer(body),
		)
		if err != nil {
			log.Errorf("Error: %s\n", err)
			continue
		}
		log.Println("Notification sent!")
	}
	return true
}
