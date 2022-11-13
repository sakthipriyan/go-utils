package main

import (
	"testing"
	"time"

	asserts "github.com/stretchr/testify/assert"
)

func TestAutoSleeper(t *testing.T) {
	sleeper := NewAutoSleeper()
	sleeper.RandomizationFactor = 0
	sleeper.UpMultiplier = 1.1
	sleeper.DownMultiplier = 0.9
	sleeper.InitialInterval = 10 * time.Microsecond
	sleeper.SleepOnSuccess()
	sleeper.SleepOnFailure()
	expected := AutoSleeperMetrics{
		TotalInvocation: 2,
		TotalWentUp:     1,
		TotalWentDown:   0,
		TotalSlept:      1,
		TotalSleepTime:  10 * time.Microsecond,
	}
	assert := asserts.New(t)
	assert.Equal(expected, sleeper.GetMetrics())
	sleeper.SleepOnSuccess()
	expected = AutoSleeperMetrics{
		TotalInvocation: 3,
		TotalWentUp:     1,
		TotalWentDown:   0,
		TotalSlept:      2,
		TotalSleepTime:  20 * time.Microsecond,
	}
	assert.Equal(expected, sleeper.GetMetrics())
	sleeper.SleepOnFailure()
	expected = AutoSleeperMetrics{
		TotalInvocation: 4,
		TotalWentUp:     2,
		TotalWentDown:   0,
		TotalSlept:      3,
		TotalSleepTime:  31 * time.Microsecond,
	}
	assert.Equal(expected, sleeper.GetMetrics())
	sleeper.SleepOnFailure()
	expected = AutoSleeperMetrics{
		TotalInvocation: 5,
		TotalWentUp:     3,
		TotalWentDown:   0,
		TotalSlept:      4,
		TotalSleepTime:  43100 * time.Nanosecond,
	}
	assert.Equal(expected, sleeper.GetMetrics())
}

func TestAutoSleeperWithGoDown(t *testing.T) {
	assert := asserts.New(t)
	sleeper := NewAutoSleeper()
	sleeper.RandomizationFactor = 0
	sleeper.UpMultiplier = 2
	sleeper.DownMultiplier = 0.5
	sleeper.DownMultiplierThreshold = 5
	sleeper.InitialInterval = 1 * time.Microsecond
	for i := 0; i < 5; i++ {
		sleeper.SleepOnFailure()
	}
	for i := 0; i < 30; i++ {
		sleeper.SleepOnSuccess()
	}
	expected := AutoSleeperMetrics{
		TotalInvocation: 35,
		TotalWentUp:     5,
		TotalWentDown:   5,
		TotalSlept:      30,
		TotalSleepTime:  (1 + 2 + 4 + 8 + 16*5 + 8*5 + 4*5 + 2*5 + 1*5) * time.Microsecond,
	}
	assert.Equal(expected, sleeper.GetMetrics())
}

func TestAutoSleeperWithRandomizationFactor(t *testing.T) {
	assert := asserts.New(t)
	sleeper := NewAutoSleeper()
	sleeper.RandomizationFactor = 0.2
	sleeper.UpMultiplier = 2
	sleeper.DownMultiplier = 0.5
	sleeper.DownMultiplierThreshold = 5
	sleeper.InitialInterval = 1 * time.Microsecond
	sleeper.SleepOnFailure()
	expectedMin := 800 * time.Nanosecond
	expectedMax := 1200 * time.Nanosecond
	actual := sleeper.GetMetrics().TotalSleepTime
	assert.True(actual >= expectedMin)
	assert.True(actual <= expectedMax)
}

func TestAutoSleeperForMaxInterval(t *testing.T) {
	assert := asserts.New(t)
	sleeper := NewAutoSleeper()
	sleeper.RandomizationFactor = 0
	sleeper.InitialInterval = 100 * time.Microsecond
	sleeper.MaxInterval = 150 * time.Microsecond
	sleeper.UpMultiplier = 2
	sleeper.DownMultiplier = 0.5
	sleeper.DownMultiplierThreshold = 1
	sleeper.SleepOnFailure()
	sleeper.SleepOnFailure()
	sleeper.SleepOnSuccess()
	sleeper.SleepOnSuccess()
	expected := AutoSleeperMetrics{
		TotalInvocation: 4,
		TotalWentUp:     2,
		TotalWentDown:   1,
		TotalSlept:      3,
		TotalSleepTime:  (100 + 150) * time.Microsecond,
	}
	assert.Equal(expected, sleeper.GetMetrics())
}

func TestAutoSleeperForMaxRandomization(t *testing.T) {
	assert := asserts.New(t)
	sleeper := NewAutoSleeper()
	sleeper.InitialInterval = 100 * time.Microsecond
	sleeper.UpMultiplier = 2
	sleeper.RandomizationFactor = 0.2
	sleeper.MaxRandomization = 10 * time.Microsecond
	sleeper.SleepOnFailure()
	sleeper.SleepOnFailure()
	expectedMin := (100 + 190) * time.Microsecond
	expectedMax := (100 + 210) * time.Microsecond
	actual := sleeper.GetMetrics().TotalSleepTime
	assert.True(expectedMin <= actual)
	assert.True(actual <= expectedMax)
}
