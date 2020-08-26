package model

import (
	v1 "k8s.io/api/core/v1"
)

type Problem struct {
	No      int64  `json:"no"`
	Title   string `json:"title"`
	Content string `json:"content"`
	// []TestCase = same scenario (use same HttpClient)
	TestCases    [][]TestCase            `json:"-"`
	Boilerplate  []Boilerplate           `json:"boilerplate"`
	Resource     v1.ResourceRequirements `json:"-"`
	WaitForReady func() <-chan bool `json:"-"`
}

type Boilerplate struct {
	Lang Language `json:"lang"`
	Code string   `json:"code"`
}
