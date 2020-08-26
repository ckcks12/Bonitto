package problems

import (
	"github.com/team-bonitto/bonitto/internal/model"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"net"
	"net/http"
	"time"
)

var P7HelloWorld = model.Problem{
	No:    7,
	Title: "Hello Word",
	Content: `
Just say Hello World to me.

- GET /say -> "Hello World"

Have fun :)
`,
	TestCases: [][]model.TestCase{
		{
			model.AccuracyHttpTestCase{
				Name:    "Hello World",
				Method:  http.MethodGet,
				Path:    "/say",
				Input:   nil,
				Output:  []byte("Hello World"),
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
	sayHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "si señor\n")
	}

	http.HandleFunc("/say", sayHandler)
	log.Fatal(http.ListenAndServe(":80", nil))
}`,
		},
	},
	Resource: v1.ResourceRequirements{
		Limits: v1.ResourceList{
			v1.ResourceCPU:    resource.MustParse("1"),
			v1.ResourceMemory: resource.MustParse("1Gi"),
		},
		Requests: v1.ResourceList{
			v1.ResourceCPU:    resource.MustParse("1"),
			v1.ResourceMemory: resource.MustParse("1Gi"),
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
