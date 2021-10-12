package webhook

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"time"
	
	log "github.com/sirupsen/logrus"

	"github.com/ConsenSys/fc-latency-map/manager/config"
)

var conf = config.NewConfig()

type Payload struct {
    datetime 	string 		`json:datetime`
	filenames 	[]string 	`json:filenames`
}

func Notify(files *[]string) {
	urls := strings.Split(conf.GetString("WEBHOOK_NOTIFY_URLS"), ",")
	payload := Payload{
		datetime: time.Now().String(),
		filenames: *files,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		log.Errorf("Error: %s\n", err)
		return
	}
	for _, url := range urls {
		log.Printf("Notify... POST: url=%s, body=%s\n", url, payload)
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
}