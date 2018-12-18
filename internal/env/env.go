package env

import "os"

const (
	envType    string = "GOCHAT_TYPE"
	envLogging string = "GOCHAT_LOGGING"
	envHost    string = "GOCHAT_SERVER_HOST"
	envPort    string = "GOCHAT_SERVER_PORT"
)

func GetServerPort(defPort string) string {
	port, r2 := os.LookupEnv(envPort)
	if !r2 {
		port = defPort
	}
	return port
}

func GetServerIp(defPort string) string {
	host, r := os.LookupEnv(envHost)
	if !r {
		host = "127.0.0.1"
	}

	port := GetServerPort(defPort)
	return host + ":" + port
}

func GetChatType() string {
	t, r := os.LookupEnv(envType)
	if !r {
		return "tcp"
	}
	return t
}

func IsLoggingEnabled() bool {
	t, r := os.LookupEnv(envLogging)

	if !r {
		return false
	}
	return t == "true"
}
