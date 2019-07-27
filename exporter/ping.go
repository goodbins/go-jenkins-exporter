package exporter

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// Ping Replies with pong when called
func Ping(rw http.ResponseWriter, req *http.Request) {
	logrus.Debug("Got ping, replying with pong")
	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("pong"))
}
