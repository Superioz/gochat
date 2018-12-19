package env

import (
	"github.com/superioz/gochat/internal/logs"
	"os"
)

const (
	Protocol    string = "GOCHAT_PROTOCOL"
	Logging     string = "GOCHAT_LOGGING"
	Host        string = "GOCHAT_SERVER_HOST"
	Port        string = "GOCHAT_SERVER_PORT"
	LoggingHost string = "GOCHAT_LOGGING_HOST"
	LoggingUser string = "GOCHAT_LOGGING_USER"
	LoggingPass string = "GOCHAT_LOGGING_PASS"
)

// returns the logging credentials fetched
// from the environmental variables
func GetLoggingCredentials() logs.LogCredentials {
	h := getOrDefault(LoggingHost, "127.0.0.1")
	u := getOrDefault(LoggingUser, "elastic")
	p := getOrDefault(LoggingPass, "changeme")

	return logs.LogCredentials{Host: h, User: u, Password: p}
}

// get the server port from environmental variables
func GetServerPort(defPort string) string {
	return getOrDefault(Port, defPort)
}

// get the server host + port from environmental variables
func GetServerIp(defPort string) string {
	host := getOrDefault(Host, "127.0.0.1")

	port := GetServerPort(defPort)
	return host + ":" + port
}

// get the protocol type from environmental variables
func GetProtocol() string {
	return getOrDefault(Protocol, "tcp")
}

// get if logging is enabled from environmental variables
func IsLoggingEnabled() bool {
	return getOrDefault(Logging, "false") == "true"
}

func getOrDefault(key string, def string) string {
	e, r := os.LookupEnv(key)
	if !r {
		return def
	}
	return e
}
