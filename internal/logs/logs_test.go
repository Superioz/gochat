package logs

import (
	"testing"
)

func TestCreateClient(t *testing.T) {
	log, err := CreateAndConnect("http://127.0.0.1:9200", "elastic", "changeme")
	if err != nil {
		t.Fatal(err)
	}

	if !log.Client.IsRunning() {
		t.Skip("couldn't connect to server")
	}
}

func TestChatLogger_AddEntry(t *testing.T) {
	log, err := CreateAndConnect("http://127.0.0.1:9200", "elastic", "changeme")
	if err != nil {
		t.Skip(err)
	}

	if !log.Client.IsRunning() {
		t.Skip("couldn't connect to server")
	}

	err = log.AddEntry("user", "a message")
	if err != nil {
		t.Skip(err)
	}
}
