package logs

import (
	"context"
	"github.com/olivere/elastic"
	"github.com/pkg/errors"
	"time"
)

const (
	logIndex string = "chatlogs"
)

type LogCredentials struct {
	Host     string
	User     string
	Password string
}

type LogEntry struct {
	User      string `json:"user"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

type ChatLogger struct {
	Client *elastic.Client
}

func (l *ChatLogger) AddEntry(user string, message string) error {
	if !l.Client.IsRunning() {
		return errors.New("log client is not running")
	}

	entry := LogEntry{User: user, Message: message, Timestamp: time.Now().Format(time.RFC3339Nano)}
	_, err := l.Client.Index().Index(logIndex).Type("doc").BodyJson(entry).Do(context.Background())
	return err
}

func CreateAndConnect(cred LogCredentials) (ChatLogger, error) {
	// Create a new elastic search client
	client, err := elastic.NewClient(elastic.SetURL("http://"+cred.Host+":9200"), elastic.SetBasicAuth(cred.Host, cred.Password), elastic.SetSniff(false))
	if err != nil {
		return ChatLogger{}, err
	}

	// Create default index for our logs
	r, err := client.IndexExists(logIndex).Do(context.Background())
	if err != nil {
		return ChatLogger{}, err
	}

	if !r {
		_, err = client.CreateIndex(logIndex).Do(context.Background())
		if err != nil {
			return ChatLogger{}, err
		}
	}
	return ChatLogger{Client: client}, nil
}
