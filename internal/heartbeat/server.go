package heartbeat

import (
	"fmt"

	"github.com/labstack/echo/v4"

	"github.com/fr33dman/go-template/pkg/probes"
)

type ProbesServer struct {
	Host   string
	Port   int
	probes *probes.Probes
	*echo.Echo
}

func NewProbesServer(host string, port int, appProbes *probes.Probes) ProbesServer {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	return ProbesServer{
		Host:   host,
		Port:   port,
		probes: appProbes,
		Echo:   e,
	}
}

func (s *ProbesServer) Start() error {
	address := fmt.Sprintf("%s:%d", s.Host, s.Port)
	s.Add(echo.GET, "/", echo.WrapHandler(s.probes.Handler()))
	return s.Echo.Start(address)
}
