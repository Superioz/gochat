package env

import (
	"github.com/superioz/gochat/internal/logs"
	"os"
)

const (
	envType        string = "GOCHAT_TYPE"
	envLogging     string = "GOCHAT_LOGGING"
	envHost        string = "GOCHAT_SERVER_HOST"
	envPort        string = "GOCHAT_SERVER_PORT"
	envLoggingHost string = "GOCHAT_LOGGING_HOST"
	envLoggingUser string = "GOCHAT_LOGGING_USER"
	envLoggingPass string = "GOCHAT_LOGGING_PASS"
)

func GetLoggingCredentials() logs.LogCredentials {
	h := getOrDefault(envLoggingHost, "127.0.0.1")
	u := getOrDefault(envLoggingUser, "elastic")
	p := getOrDefault(envLoggingPass, "changeme")

	return logs.LogCredentials{Host: h, User: u, Password: p}
}

func GetServerPort(defPort string) string {
	return getOrDefault(envPort, defPort)
}

func GetServerIp(defPort string) string {
	host := getOrDefault(envHost, "127.0.0.1")

	port := GetServerPort(defPort)
	return host + ":" + port
}

func GetChatType() string {
	return getOrDefault(envType, "tcp")
}

func IsLoggingEnabled() bool {
	return getOrDefault(envLogging, "false") == "true"
}

func getOrDefault(key string, def string) string {
	e, r := os.LookupEnv(key)
	if !r {
		return def
	}
	return e
}
