package main

import (
	"fmt"
	consumer2 "github.com/team-bonitto/bonitto/internal/queue/consumer"
	recorder2 "github.com/team-bonitto/bonitto/internal/recorder"
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
	recorder, err := recorder2.New(consumer.RDB)
	if err != nil {
		panic(err)
	}
	for {
		<-time.After(time.Duration(interval) * time.Millisecond)
		start := time.Now()
		log := fmt.Sprintf("%s [recorder] : ", start.String())

		if err := consumer.Consume(recorder); err != nil {
			log += err.Error()
		}
		end := time.Now()
		log += fmt.Sprintf(" (%dms)", end.Unix()-start.Unix())
		fmt.Println(log)
	}
}
