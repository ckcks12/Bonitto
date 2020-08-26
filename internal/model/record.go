package model

type Record struct {
	ProblemNo   int64 `json:"no"`
	UserID      string `json:"id"`
	TestResults [][]TestResult `json:"results"`
}
