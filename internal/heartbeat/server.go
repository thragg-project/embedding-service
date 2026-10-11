package heartbeat

import (
	"net"
	"net/http"
	"strconv"

	"github.com/fr33dman/go-template/pkg/probes"
)

type ProbesServer struct {
	*http.Server
}

func NewProbesServer(host string, port int, appProbes *probes.Probes) ProbesServer {
	return ProbesServer{
		Server: &http.Server{
			Addr:    net.JoinHostPort(host, strconv.Itoa(port)),
			Handler: appProbes.Handler(),
		},
	}
}

func (s *ProbesServer) Start() error {
	return s.ListenAndServe()
}
