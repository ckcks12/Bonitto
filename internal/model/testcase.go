package model

import (
	"time"
)

type TestCase interface {
	GetName() string
	Run() TestResult
}

type TestResult struct {
	Passed bool      `json:"passed"`
	Result string    `json:"result"`
	Date   time.Time `json:"date"`
}

func NewTestResult(passed bool, result string) TestResult {
	return TestResult{
		Passed: passed,
		Result: result,
		Date:   time.Now(),
	}
}
