package model

import (
	"net/http"
	"time"
)

type HttpTestCase struct {
	Name         string
	Method       string
	Path         string
	Input        []byte
	Output       []byte
	Status       int
	Client       http.Client
	Timeout      time.Duration
}
