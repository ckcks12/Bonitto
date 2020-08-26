package problems

import (
	"encoding/json"
	"fmt"
	"github.com/team-bonitto/bonitto/internal/model"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/rand"
	"net"
	"net/http"
	"strings"
	"time"
)

var P8User = model.Problem{
	No:    8,
	Title: "User",
	Content: `
Implement User API

Here are goals you should implement:
	- Create/Read/Update/Delete User
	- Invalid User Input Filtering
	- API Limit for Too Many Request

User Struct:
{
	id: string (len: 5~12) only allowed alphanumeric and _
	pw: string (hash: ~50)
	name: string (len: 2~50) only allowed alphanumeric
	email: string (len: ~50)
	phone: string (len: ~20) international phone number format
	visit: integer should be +1 every GET request
	deleted: boolean whether this user has been deleted
}
	

To sum up, These are endpoints:
	- GET /user?id=
		- input: 
		- output: {"id":"id","pw":"hashed","name":"name","email":"my@email.com","phone":"+821012341234","visit":34,"deleted":false}
	- POST /user
		- input: {"id":"id","pw":"pw","name":"name","email":"my@email.com","phone":"+821012341234"}
		- output: {"id":"id","pw":"hashed","name":"name","email":"my@email.com","phone":"+821012341234","visit":34,"deleted":false}
	- PUT /user
		- input: {"id":"id","name":"abc","email":"your@email.com","phone":"+12123412341"}
		- output: {"id":"id","name":"abc","email":"your@email.com","phone":"+12123412341","visit":34,"deleted":false}
	- DELETE /user
		- input: {"id":"id"}
		- output: 

	- Always output its full object
	- id, pw, deleted field cannot be modified

Note:
	- Server spec is 0.1 vCPU and 128MiB RAM.
	- For now, use memory as a database
	- Every output should be a JSON format and minified.
	- Every transactions should be done within 1 sec.
	- Use semantic status codes:
		- 200 OK
		- 400 Bad Request

Have fun :)
`,
	TestCases: [][]model.TestCase{
		{
			model.AccuracyHttpTestCase{
				Name:    "create a user",
				Method:  http.MethodPost,
				Path:    "/user",
				Input:   encodeJSON(fakeUser),
				Output:  encodeJSON(fakeUser),
				Status:  http.StatusOK,
				Timeout: 1000 * time.Millisecond,
				Client:  *http.DefaultClient,
			},
			model.AccuracyHttpTestCase{
				Name:    "get a user",
				Method:  http.MethodGet,
				Path:    "/user?id=" + fakeUser.ID,
				Input:   nil,
				Output:  encodeJSON(fakeUser),
				Status:  http.StatusOK,
				Timeout: 1000 * time.Millisecond,
				Client:  *http.DefaultClient,
			},
			model.AccuracyHttpTestCase{
				Name:    "update a user name",
				Method:  http.MethodPut,
				Path:    "/user",
				Input:   []byte(fmt.Sprintf(`{"id":"%s","name":"%s"}`, fakeUser.ID, fakeUser2.Name)),
				Output:  encodeJSON(fakeUser.WithName(fakeUser2.Name)),
				Status:  http.StatusOK,
				Timeout: 1000 * time.Millisecond,
				Client:  *http.DefaultClient,
			},
			model.AccuracyHttpTestCase{
				Name:    "delete a user",
				Method:  http.MethodDelete,
				Path:    "/user",
				Input:   []byte(fmt.Sprintf(`{"id":"%s"}`, fakeUser.ID)),
				Output:  nil,
				Status:  http.StatusOK,
				Timeout: 1000 * time.Millisecond,
				Client:  *http.DefaultClient,
			},
		},
		{
			model.AccuracyHttpTestCase{
				Name:    "get an unregistered user",
				Method:  http.MethodGet,
				Path:    "/user?id=abc123",
				Input:   nil,
				Output:  nil,
				Status:  http.StatusBadRequest,
				Timeout: 1000 * time.Millisecond,
				Client:  *http.DefaultClient,
			},
			model.AccuracyHttpTestCase{
				Name:    "create a user with incorrect id",
				Method:  http.MethodPost,
				Path:    "/user",
				Input:   encodeJSON(fakeUser.WithID(".")),
				Output:  nil,
				Status:  http.StatusBadRequest,
				Timeout: 1000 * time.Millisecond,
				Client:  *http.DefaultClient,
			},
			model.AccuracyHttpTestCase{
				Name:    "create a user with incorrect pw",
				Method:  http.MethodPost,
				Path:    "/user",
				Input:   encodeJSON(fakeUser.WithPW(newRandomAlphanumericString(51, 60))),
				Output:  nil,
				Status:  http.StatusBadRequest,
				Timeout: 1000 * time.Millisecond,
				Client:  *http.DefaultClient,
			},
			model.AccuracyHttpTestCase{
				Name:    "create a user with incorrect name",
				Method:  http.MethodPost,
				Path:    "/user",
				Input:   encodeJSON(fakeUser.WithName(newRandomAlphanumericString(51, 60))),
				Output:  nil,
				Status:  http.StatusBadRequest,
				Timeout: 1000 * time.Millisecond,
				Client:  *http.DefaultClient,
			},
			model.AccuracyHttpTestCase{
				Name:    "create a user with incorrect email",
				Method:  http.MethodPost,
				Path:    "/user",
				Input:   encodeJSON(fakeUser.WithEmail(newRandomAlphanumericString(10, 20))),
				Output:  nil,
				Status:  http.StatusBadRequest,
				Timeout: 1000 * time.Millisecond,
				Client:  *http.DefaultClient,
			},
			model.AccuracyHttpTestCase{
				Name:    "create a user with incorrect phone",
				Method:  http.MethodPost,
				Path:    "/user",
				Input:   encodeJSON(fakeUser.WithPhone(newRandomAlphanumericString(10, 20))),
				Output:  nil,
				Status:  http.StatusBadRequest,
				Timeout: 1000 * time.Millisecond,
				Client:  *http.DefaultClient,
			},
			model.AccuracyHttpTestCase{
				Name:    "update a user with incorrect name",
				Method:  http.MethodPut,
				Path:    "/user",
				Input:   encodeJSON(fakeUser.WithName(newRandomAlphanumericString(51, 60))),
				Output:  nil,
				Status:  http.StatusBadRequest,
				Timeout: 1000 * time.Millisecond,
				Client:  *http.DefaultClient,
			},
			model.AccuracyHttpTestCase{
				Name:    "update a user with incorrect email",
				Method:  http.MethodPut,
				Path:    "/user",
				Input:   encodeJSON(fakeUser.WithEmail(newRandomAlphanumericString(10, 20))),
				Output:  nil,
				Status:  http.StatusBadRequest,
				Timeout: 1000 * time.Millisecond,
				Client:  *http.DefaultClient,
			},
			model.AccuracyHttpTestCase{
				Name:    "update a user with incorrect phone",
				Method:  http.MethodPut,
				Path:    "/user",
				Input:   encodeJSON(fakeUser.WithPhone(newRandomAlphanumericString(10, 20))),
				Output:  nil,
				Status:  http.StatusBadRequest,
				Timeout: 1000 * time.Millisecond,
				Client:  *http.DefaultClient,
			},
			model.AccuracyHttpTestCase{
				Name:    "delete an unregistered user",
				Method:  http.MethodDelete,
				Path:    "/user",
				Input:   []byte(fmt.Sprintf(`{"id":"%s"}`, fakeUser.ID)),
				Output:  nil,
				Status:  http.StatusBadRequest,
				Timeout: 1000 * time.Millisecond,
				Client:  *http.DefaultClient,
			},
		},
	},
	Boilerplate: []model.Boilerplate{
		{
			Lang: model.Javascript,
			Code: P8BoilerplateJavascript,
		},
		{
			Lang: model.Go,
			Code: P8BoilerplateGo,
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

const P8BoilerplateJavascript = `const http = require('http');

const handler = (req, res) => {
    res.writeHead(200);
    res.end('hello, world');
}

const server = http.createServer(handler);
server.listen(80);`

const P8BoilerplateGo = `package main

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
}`

type P8UserStruct struct {
	ID      string `json:"id"`
	PW      string `json:"pw"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Visit   int    `json:"visit"`
	Deleted bool   `json:"deleted"`
}

var fakeUser = newFakeUser()
var fakeUser2 = newFakeUser()

func (u P8UserStruct) WithID(a string) P8UserStruct {
	u2 := u
	u2.ID = a
	return u2
}

func (u P8UserStruct) WithPW(a string) P8UserStruct {
	u2 := u
	u2.PW = a
	return u2
}

func (u P8UserStruct) WithName(a string) P8UserStruct {
	u2 := u
	u2.Name = a
	return u2
}

func (u P8UserStruct) WithEmail(a string) P8UserStruct {
	u2 := u
	u2.Email = a
	return u2
}

func (u P8UserStruct) WithPhone(a string) P8UserStruct {
	u2 := u
	u2.Phone = a
	return u2
}

func (u P8UserStruct) WithVisit(a int) P8UserStruct {
	u2 := u
	u2.Visit = a
	return u2
}

func (u P8UserStruct) WithDeleted(a bool) P8UserStruct {
	u2 := u
	u2.Deleted = a
	return u2
}

func newFakeUser() P8UserStruct {
	return P8UserStruct{
		ID:      newRandomAlphanumericString(5, 12),
		PW:      newRandomAlphanumericString(5, 30),
		Name:    newRandomAlphanumericString(4, 30),
		Email:   newRandomEmail(),
		Phone:   newRandomPhone(),
		Visit:   0,
		Deleted: false,
	}
}

func newRandomAlphanumericString(min, max int) string {
	len := rand.IntnRange(min, max)
	str := strings.Builder{}
	for i := 0; i < len; i++ {
		str.WriteRune(getRandomAlphanumericRune())
	}
	return str.String()
}

func getRandomAlphanumericRune() rune {
	sources := []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'j', 'k', 'l', 'm', 'n',
		'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '_'}
	len := len(sources)
	return sources[rand.Intn(len)]
}

func newRandomEmail() string {
	id := newRandomAlphanumericString(5, 15)
	domain := newRandomAlphanumericString(3, 10)
	tld := newRandomAlphanumericString(3, 5)
	return fmt.Sprintf("%s@%s.%s", id, domain, tld)
}

func newRandomPhone() string {
	country := rand.IntnRange(1, 20)
	area := rand.IntnRange(10, 1000)
	phone1 := rand.IntnRange(100, 10000)
	phone2 := rand.IntnRange(100, 10000)
	return fmt.Sprintf("+%d%d%d%d", country, area, phone1, phone2)
}

func encodeJSON(a interface{}) []byte {
	b, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}
	return b
}
