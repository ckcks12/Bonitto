package problems

import (
	"github.com/team-bonitto/bonitto/internal/model"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"net"
	"net/http"
	"time"
)

var P1311TwitterAPI = model.Problem{
	No:    1311,
	Title: "Twitter API",
	Content: `
Implement Twitter API.

There are three big components:
	- User
	- Tweet
	- Like

Here are goals we should achieve:
	- Users create/read/delete tweets
	- Users retweet others' tweets
	- Users like others' tweets

To sum up, These are endpoints:
	- GET /user -> current user info
	- GET /user/:id -> user info by user id
	- POST /user -> create user
	- DELETE /user -> delete user
	- GET /tweets/:id -> tweets by user id
	- GET /tweets/:tid -> retweets by tweet id
	- GET /tweet/:tid -> tweet by tweet id
	- POST /tweet -> create tweet
	- DELETE /tweet/:tid -> delete tweet
	- POST /like/:tid -> like tweet
	- DELETE /like/:tid -> cancel like

Server spec is 0.5 vCPU and 256MiB RAM.
Server should be able to process more than 120req/s.
which means Concurrency is 60, Requests per Second is 2.

Note:
	Server should startup within 10 sec

Have fun :)
`,
	TestCases: [][]model.TestCase{
		{
			model.AccuracyHttpTestCase{
				Name:    "Get User",
				Method:  http.MethodGet,
				Path:    "/user",
				Input:   nil,
				Output:  nil,
				Status:  http.StatusOK,
				Timeout: 1000 * time.Millisecond,
				Client:  *http.DefaultClient,
			},
		},
	},
	Boilerplate: []model.Boilerplate{
		{
			Lang: model.Javascript,
			Code: `const http = require('http');

const handler = (req, res) => {
    res.writeHead(200);
    res.end('hello, world');
}

const server = http.createServer(handler);
server.listen(80);`,
		},
		{
			Lang: model.Go,
			Code: `package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	handler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "si señor\n")
	}

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":80", nil))
}`,
		},
	},
	Resource: v1.ResourceRequirements{
		Limits: v1.ResourceList{
			v1.ResourceCPU:    resource.MustParse("0.5"),
			v1.ResourceMemory: resource.MustParse("256Mi"),
		},
		Requests: v1.ResourceList{
			v1.ResourceCPU:    resource.MustParse("0.5"),
			v1.ResourceMemory: resource.MustParse("256Mi"),
		},
	},
	WaitForReady: func() <-chan bool {
		retry := 5
		timeout := 1 * time.Second
		c := make(chan bool)
		go func() {
			for i := 0; i < retry; i++ {
				<-time.After(timeout)
				conn, err := net.DialTimeout("tcp", "localhost:80", timeout)
				if err != nil {
					continue
				}
				_ = conn.Close()
				c <- true
				return
			}
			c <- false
		}()
		return c
	},
}
