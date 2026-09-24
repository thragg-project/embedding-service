package heartbeat

import (
	"github.com/fr33dman/go-template/pkg/probes"
)

func NewProbes() *probes.Probes {
	appProbes := probes.NewProbes()
	appProbes.SetReadiness(true)
	return appProbes
}
