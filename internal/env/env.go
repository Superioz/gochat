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

	defaultProtocol    string = "tcp"
	defaultLogging     string = "false"
	defaultHost        string = "127.0.0.1"
	defaultPort        string = "6000"
	defaultLoggingHost string = "127.0.0.1"
	defaultLoggingUser string = "elastic"
	defaultLoggingPass string = "changeme"
)

// set defaults for protocols other than
// `tcp`
func SetDefaults(prot string) {
	switch prot {
	case "amqp":
		_ = os.Setenv(Protocol, "amqp")
		_ = os.Setenv(Logging, "true")
		_ = os.Setenv(Host, "amqp://guest:guest@localhost")
		_ = os.Setenv(Port, "5672")
		break
	case "kafka":
		_ = os.Setenv(Protocol, "kafka")
		_ = os.Setenv(Logging, "true")
		_ = os.Setenv(Host, "localhost")
		_ = os.Setenv(Port, "9092")
		break
	case "tcp":
		// take default values
		break
	}
}

// returns the logging credentials fetched
// from the environmental variables
func GetLoggingCredentials() logs.LogCredentials {
	h := getOrDefault(LoggingHost, defaultLoggingHost)
	u := getOrDefault(LoggingUser, defaultLoggingUser)
	p := getOrDefault(LoggingPass, defaultLoggingPass)

	return logs.LogCredentials{Host: h, User: u, Password: p}
}

// get the server port from environmental variables
func GetServerPort(defPort string) string {
	return getOrDefault(Port, defPort)
}

// get the server host + port from environmental variables
func GetServerIp() string {
	host := getOrDefault(Host, defaultHost)

	port := GetServerPort(defaultPort)
	return host + ":" + port
}

// get the protocol type from environmental variables
func GetProtocol() string {
	return getOrDefault(Protocol, defaultProtocol)
}

// get if logging is enabled from environmental variables
func IsLoggingEnabled() bool {
	return getOrDefault(Logging, defaultLogging) == "true"
}

func getOrDefault(key string, def string) string {
	e, r := os.LookupEnv(key)
	if !r {
		return def
	}
	return e
}
