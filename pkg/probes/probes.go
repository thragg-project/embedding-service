package probes

import (
	"context"
	"sync/atomic"
	"time"
)

type (
	HeartBeat interface {
		Beat(ctx context.Context) error
	}
	Probes struct {
		isAlive, isReady atomic.Bool

		livenessHeartbeat, readinessHeartbeat                                  HeartBeat
		livenessPeriod, livenessThreshold, readinessPeriod, readinessThreshold time.Duration
	}
)

func NewProbes() *Probes {
	p := &Probes{
		livenessHeartbeat:  nil,
		readinessHeartbeat: nil,
	}
	p.isAlive.Store(true)
	p.isReady.Store(false)
	return p
}

func (p *Probes) SetLiveness(status bool) {
	p.isAlive.Store(status)
}

func (p *Probes) SetReadiness(status bool) {
	p.isReady.Store(status)
}

func (p *Probes) SetLivenessHeartbeat(heartbeat HeartBeat, period, threshold time.Duration) {
	p.livenessHeartbeat = heartbeat
	p.livenessPeriod = period
	p.livenessThreshold = threshold
}

func (p *Probes) SetReadinessHeartbeat(heartbeat HeartBeat, period, threshold time.Duration) {
	p.readinessHeartbeat = heartbeat
	p.readinessPeriod = period
	p.readinessThreshold = threshold
}

func (p *Probes) Heartbeat(ctx context.Context) {
	if p.livenessHeartbeat != nil {
		runHeartbeat(ctx, p.livenessHeartbeat, &p.isAlive, p.livenessPeriod, p.livenessThreshold)
	}
	if p.readinessHeartbeat != nil {
		runHeartbeat(ctx, p.readinessHeartbeat, &p.isReady, p.readinessPeriod, p.readinessThreshold)
	}
}

func runHeartbeat(
	ctx context.Context,
	heartbeat HeartBeat,
	probe *atomic.Bool,
	period, threshold time.Duration,
) {
	go func() {
		ticker := time.NewTicker(period)
		defer ticker.Stop()

		runBeat(ctx, heartbeat, probe, threshold)

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				runBeat(ctx, heartbeat, probe, threshold)
			}
		}
	}()
}

func runBeat(
	ctx context.Context,
	heartbeat HeartBeat,
	probe *atomic.Bool,
	threshold time.Duration,
) {
	beatCtx := ctx
	cancel := func() {}
	if threshold > 0 {
		beatCtx, cancel = context.WithTimeout(ctx, threshold)
	}
	defer cancel()

	err := heartbeat.Beat(beatCtx)
	probe.Store(err == nil)
}
