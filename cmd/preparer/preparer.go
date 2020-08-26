package main

import (
	"fmt"
	preparer2 "github.com/team-bonitto/bonitto/internal/preparer"
	consumer2 "github.com/team-bonitto/bonitto/internal/queue/consumer"
	"os"
	"strconv"
	"time"
)

func main() {
	addr := os.Getenv("REDIS_URL")
	interval, _ := strconv.Atoi(os.Getenv("WORKER_INTERVAL_MS"))
	consumer, err := consumer2.New(addr)
	if err != nil {
		panic(err)
	}
	preparer, err := preparer2.New(
		os.Getenv("IMAGE_TESTER"),
		os.Getenv("IMAGE_JAVASCRIPT"),
		os.Getenv("IMAGE_JAVA"),
		os.Getenv("IMAGE_GOLANG"),
		os.Getenv("IMAGE_CPP"),
		os.Getenv("COMMAND_JAVASCRIPT"),
		os.Getenv("COMMAND_JAVA"),
		os.Getenv("COMMAND_GOLANG"),
		os.Getenv("COMMAND_CPP"),
		addr,
	)
	if err != nil {
		panic(err)
	}
	for {
		<-time.After(time.Duration(interval) * time.Millisecond)
		start := time.Now()
		log := fmt.Sprintf("%s [preparer] : ", start.String())

		data, err := consumer.ConsumeManually(preparer.GetQueueName())
		if err != nil {
			log += err.Error()
		} else if err := preparer.Consume(data); err != nil {
			log += err.Error()
		} else {
			log += " done"
		}

		end := time.Now()
		log += fmt.Sprintf(" (%dms)", end.Unix()-start.Unix())
		fmt.Println(log)
	}
}
