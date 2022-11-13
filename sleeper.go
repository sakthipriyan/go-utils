package main

import (
	"math/rand"
	"time"
)

const (
	DefaultMaxInterval             = 15 * time.Minute
	DefaultInitialInterval         = 500 * time.Millisecond
	DefaultRandomizationFactor     = 0.3
	DefaultMaxRandomization        = 2 * time.Minute
	DefaultUpMultiplier            = 1.5
	DefaultDownMultiplier          = 0.9
	DefaultDownMultiplierThreshold = 10
)

func NewAutoSleeper() *AutoSleeper {
	return &AutoSleeper{
		MaxInterval:             DefaultMaxInterval,
		InitialInterval:         DefaultInitialInterval,
		RandomizationFactor:     DefaultRandomizationFactor,
		MaxRandomization:        DefaultMaxRandomization,
		UpMultiplier:            DefaultUpMultiplier,
		DownMultiplier:          DefaultDownMultiplier,
		DownMultiplierThreshold: DefaultDownMultiplierThreshold,
	}
}

type AutoSleeperMetrics struct {
	TotalInvocation int
	TotalWentUp     int
	TotalWentDown   int
	TotalSlept      int
	TotalSleepTime  time.Duration
}

type AutoSleeper struct {
	InitialInterval         time.Duration
	MaxInterval             time.Duration
	MaxRandomization        time.Duration
	UpMultiplier            float64
	DownMultiplier          float64
	RandomizationFactor     float64
	DownMultiplierThreshold int

	metrics         AutoSleeperMetrics
	currentInterval time.Duration
	currentSuccess  int
}

func (s *AutoSleeper) GetMetrics() AutoSleeperMetrics {
	return s.metrics
}

func (s *AutoSleeper) SleepOnFailure() {
	s.metrics.TotalInvocation += 1
	s.goUp()
	s.sleep()
}

func (s *AutoSleeper) SleepOnSuccess() {
	s.metrics.TotalInvocation += 1
	if s.currentInterval == 0 {
		return
	}
	s.currentSuccess += 1
	if s.currentSuccess == s.DownMultiplierThreshold {
		s.goDown()
		s.currentSuccess = 0
	}
	s.sleep()
}

func (s *AutoSleeper) sleep() {
	s.metrics.TotalSleepTime += s.currentInterval
	s.metrics.TotalSlept += 1
	time.Sleep(s.currentInterval)
}

func (s *AutoSleeper) goDown() {
	s.metrics.TotalWentDown += 1
	interval := float64(s.currentInterval) * s.DownMultiplier
	random := getNextRandomInterval(interval, s.RandomizationFactor, float64(s.MaxRandomization))
	if random < float64(s.InitialInterval) {
		s.currentInterval = 0
		return
	}
	s.currentInterval = time.Duration(random)
}

func (s *AutoSleeper) goUp() {
	s.metrics.TotalWentUp += 1
	if s.currentInterval == 0 {
		s.currentInterval = s.InitialInterval
		return
	}
	interval := float64(s.currentInterval) * s.UpMultiplier
	random := getNextRandomInterval(interval, s.RandomizationFactor, float64(s.MaxRandomization))
	if random > float64(s.MaxInterval) {
		s.currentInterval = s.MaxInterval
		return
	}
	s.currentInterval = time.Duration(random)
}

func getNextRandomInterval(currentInterval, randomizationFactor, maxRandomization float64) float64 {
	if randomizationFactor == 0 {
		return currentInterval
	}
	delta := randomizationFactor * currentInterval
	if delta > maxRandomization {
		delta = maxRandomization
	}
	randomization := 2 * delta * rand.Float64()
	minInterval := currentInterval - delta
	return minInterval + randomization
}
