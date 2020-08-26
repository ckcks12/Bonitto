package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"github.com/gofiber/websocket"
	"github.com/team-bonitto/bonitto/internal/model"
	notifier2 "github.com/team-bonitto/bonitto/internal/notifier"
	"github.com/team-bonitto/bonitto/internal/preparer"
	"github.com/team-bonitto/bonitto/internal/problems"
	consumer2 "github.com/team-bonitto/bonitto/internal/queue/consumer"
	producer2 "github.com/team-bonitto/bonitto/internal/queue/producer"
	recorder2 "github.com/team-bonitto/bonitto/internal/recorder"
	"os"
	"strconv"
	"time"
)

var notiMap = make(map[string]chan string)

type SubmitInput struct {
	UserID string `json:"id"`
	Code   string `json:"code"`
	Lang   string `json:"lang"`
}

func main() {
	addr := os.Getenv("REDIS_URL")
	producer, err := producer2.New(addr)
	if err != nil {
		panic(err)
	}
	consumer, err := consumer2.New(addr)
	if err != nil {
		panic(err)
	}
	notifier, err := notifier2.New()
	if err != nil {
		panic(err)
	}
	recorder, err := recorder2.New(consumer.RDB)
	if err != nil {
		panic(err)
	}

	app := fiber.New()
	app.Use(middleware.Logger())

	app.Get("/live", func(c *fiber.Ctx) {
		c.SendStatus(200)
	})

	app.Get("/ready", func(c *fiber.Ctx) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		if err := consumer.RDB.Ping(ctx).Err(); err != nil {
			c.SendStatus(500)
			return
		}
		c.SendStatus(200)
	})

	app.Get("/api/problems", func(c *fiber.Ctx) {
		b, _ := json.Marshal(problems.Problems)
		c.SendBytes(b)
	})

	app.Get("/api/problem/:no", func(c *fiber.Ctx) {
		noStr := c.Params("no", "")
		if noStr == "" {
			c.SendStatus(400)
			return
		}
		no, _ := strconv.Atoi(noStr)
		var p *model.Problem
		for _, problem := range problems.Problems {
			if problem.No == int64(no) {
				p = &problem
				break
			}
		}
		if p == nil {
			c.SendStatus(400)
			return
		}
		b, _ := json.Marshal(p)
		c.SendBytes(b)
	})

	app.Post("/api/submit/:no", func(c *fiber.Ctx) {
		noStr := c.Params("no", "")
		if noStr == "" {
			c.SendStatus(400)
			return
		}
		no, _ := strconv.Atoi(noStr)
		submitInput := SubmitInput{}
		if err := c.BodyParser(&submitInput); err != nil {
			c.SendStatus(400)
			return
		}
		problem := model.Problem{}
		for _, p := range problems.Problems {
			if p.No == int64(no) {
				problem = p
			}
		}
		input := preparer.Input{
			UserID: submitInput.UserID,
			Codes: []model.Code{{
				Name:    "",
				Content: submitInput.Code,
			}},
			Problem: problem,
			Lang:    model.Language(submitInput.Lang),
		}
		if err := producer.Produce(input); err != nil {
			c.SendStatus(400)
			return
		}
		c.SendStatus(200)
	})

	app.Get("/api/result/:id/:no", func(c *fiber.Ctx) {
		id := c.Params("id", "")
		no := c.Params("no", "")
		if id == "" || no == "" {
			c.SendStatus(400)
			return
		}
		results, err := recorder.GetResults(id, no)
		if err != nil {
			c.SendStatus(400)
			return
		}
		c.JSON(results)
	})

	app.Get("/api/ws/:id", websocket.New(func(c *websocket.Conn) {
		id := c.Params("id")
		for {
			if _, _, err := c.ReadMessage(); err != nil {
				fmt.Println(err)
				break
			}
			ch, ok := notiMap[id]
			if !ok {
				continue
			}
			if len(ch) == 0 {
				continue
			}
			content := <-ch
			if err := c.WriteMessage(websocket.TextMessage, []byte(content)); err != nil {
				fmt.Println(err)
			}
		}
	}))

	go func() {
		for {
			data, err := consumer.ConsumeManually(notifier.GetQueueName())
			if err != nil {
				continue
			}
			noti := notifier2.Input{}
			if err := json.Unmarshal([]byte(data), &noti); err != nil {
				fmt.Println(err)
				continue
			}
			c, ok := notiMap[noti.UserID]
			if !ok {
				notiMap[noti.UserID] = make(chan string, 100)
				c = notiMap[noti.UserID]
			}
			c <- noti.Content
		}
	}()

	app.Listen(8080)
}
