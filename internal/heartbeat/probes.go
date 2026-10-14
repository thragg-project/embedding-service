package heartbeat

import (
	"github.com/thragg-project/embedding-service/pkg/probes"
)

func NewProbes() *probes.Probes {
	appProbes := probes.NewProbes()
	appProbes.SetReadiness(true)
	return appProbes
}
