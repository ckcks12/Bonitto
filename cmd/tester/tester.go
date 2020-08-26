package main

import (
	"fmt"
	"github.com/team-bonitto/bonitto/internal/model"
	"github.com/team-bonitto/bonitto/internal/notifier"
	"github.com/team-bonitto/bonitto/internal/problems"
	producer2 "github.com/team-bonitto/bonitto/internal/queue/producer"
	"github.com/team-bonitto/bonitto/internal/recorder"
	"os"
)

func main() {
	addr := os.Getenv("REDIS_URL")
	problemNo := os.Getenv("PROBLEM_NO")
	userID := os.Getenv("USER_ID")
	producer, err := producer2.New(addr)
	if err != nil {
		panic(err)
	}

	problem := model.Problem{}
	for _, p := range problems.Problems {
		no := fmt.Sprintf("%d", p.No)
		if no == problemNo {
			problem = p
			break
		}
	}
	if problem.No == 0 {
		panic("wrong problem no : " + problemNo)
	}

	if problem.WaitForReady != nil {
		notify(producer, userID, "🚪 Wait for program ready")
		res := <-problem.WaitForReady()
		if !res {
			notify(producer, userID, "😓 Program didn't get ready")
			os.Exit(1)
		}
	}

	results := make([][]model.TestResult, 0)
	for senarioIdx, scenario := range problem.TestCases {
		result := make([]model.TestResult, 0)
		for tcIdx, tc := range scenario {
			notify(producer, userID,
				fmt.Sprintf("🧪 Scenario #%d/%d Test #%d/%d", senarioIdx+1, len(problem.TestCases), tcIdx+1, len(scenario)))
			res := tc.Run()
			result = append(result, res)
		}
		results = append(results, result)
	}

	i := recorder.Input{
		ProblemNo:   problem.No,
		UserID:      userID,
		TestResults: results,
	}
	if err := producer.Produce(i); err != nil {
		panic(err)
	}

	notify(producer, userID, "🤩 All tests have been done!")
}

func notify(producer *producer2.RedisProducer, userID string, content string) {
	fmt.Println(content)
	noti := notifier.Input{
		UserID:  userID,
		Content: content,
	}
	if err := producer.Produce(noti); err != nil {
		fmt.Println(err)
	}
}
