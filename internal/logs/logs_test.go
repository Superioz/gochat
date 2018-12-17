package logs

import "testing"

func TestCreateClient(t *testing.T) {
	log, err := CreateAndConnect("http://127.0.0.1:9200", "elastic", "changeme")
	if err != nil {
		t.Fatal(err)
	}

	if !log.Client.IsRunning() {
		t.Fatal("cant connect to server")
	}
}
