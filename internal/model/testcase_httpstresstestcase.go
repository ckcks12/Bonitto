package model

import (
	"time"
)

type StressHttpTestCase struct {
	HttpTestCase
	Duration    time.Duration
	Concurrency int64
	RPS         int64
}

var _ TestCase = &StressHttpTestCase{}

func (tc StressHttpTestCase) GetName() string {
	return tc.Name
}

func (tc StressHttpTestCase) Run() TestResult {
	panic("implement me")
}
