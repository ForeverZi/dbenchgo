package task

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Task interface {
	SetUp(ctx context.Context) error
	Run(ctx context.Context) (statusCh <-chan Status, err error)
	CleanUp(ctx context.Context) error
	CollectResult(ctx context.Context) (result Result, err error)
}

type Status struct {
	Done      bool
	Completed int
	Total     int
	Error     error
}

type Result interface {
	Summary() string
}

type OpResult struct {
	mux      sync.Mutex
	Count    int
	Total    time.Duration
	Min      time.Duration
	Max      time.Duration
	ErrCount int
}

func NewOpResult() *OpResult {
	return &OpResult{
		Min: time.Hour * 100,
	}
}

func (result *OpResult) Record(cost time.Duration, success bool) {
	result.mux.Lock()
	defer result.mux.Unlock()
	result.Count++
	result.Total += cost
	if cost < result.Min {
		result.Min = cost
	}
	if cost > result.Max {
		result.Max = cost
	}
	if !success {
		result.ErrCount++
	}
}

func (result *OpResult) Summary() string {
	result.mux.Lock()
	defer result.mux.Unlock()
	summary := fmt.Sprintf(`================summary====================
	Execute Count: %v,
	Error Count: %v,
	Total Cost: %v,
	Avg Cost: %v,
	Min Cost: %v,
	Max Cost: %v
======================================================
`, result.Count, result.ErrCount, result.Total, result.Total/time.Duration(result.Count), result.Min, result.Max)
	return summary
}
