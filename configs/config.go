package configs

import (
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	defaultAddr             = ":8080"
	defaultDebug            = true
	defaultConnectionsLimit = 100
)

type Config struct {
	ListenAddr       string
	Debug            bool
	ConnectionsLimit int32
}

var conf *Config

func NewAppConfig() *Config {
	once := sync.Once{}
	once.Do(func() {
		conf = &Config{
			ListenAddr:       defaultAddr,
			Debug:            defaultDebug,
			ConnectionsLimit: defaultConnectionsLimit,
		}

		if addr, exists := os.LookupEnv("LISTEN_ADDR"); exists {
			conf.ListenAddr = addr
		}

		if debug, exists := os.LookupEnv("DEBUG"); exists && len(debug) != 0 { // might be DEBUG=
			conf.Debug = strings.ToUpper(debug) != "FALSE"
		}

		if limit, exists := os.LookupEnv("CONNECTIONS_LIMIT"); exists {
			lim, _ := strconv.ParseInt(limit, 10, 32)
			conf.ConnectionsLimit = int32(lim)
		}
	})
	return conf
}
