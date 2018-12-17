package logs

import (
	"context"
	"github.com/olivere/elastic"
)

const (
	logIndex string = "chatlogs"
)

type LogEntry struct {
	User      string `json:"user"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

type ChatLogger struct {
	Client *elastic.Client
}

func CreateAndConnect(url string, user string, pass string) (ChatLogger, error) {
	// Create a new elastic search client
	client, err := elastic.NewClient(elastic.SetURL(url), elastic.SetBasicAuth(user, pass), elastic.SetSniff(false))
	if err != nil {
		return ChatLogger{}, err
	}

	// Create default index for our logs
	_, err = client.CreateIndex(logIndex).Do(context.Background())
	if err != nil {
		return ChatLogger{}, err
	}

	return ChatLogger{Client: client}, nil
}
