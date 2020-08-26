package problems_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/team-bonitto/bonitto/internal/problems"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

var DB = make([]problems.P8UserStruct, 0)

func add(u problems.P8UserStruct) {
	DB = append(DB, u)
}

func find(id string) (problems.P8UserStruct, error) {
	for _, u := range DB {
		if u.ID == id {
			return u, nil
		}
	}
	return problems.P8UserStruct{}, errors.New("not found")
}

func del(id string) error {
	idx := -1
	for i, u := range DB {
		if u.ID == id {
			idx = i
			break
		}
	}
	if idx < 0 {
		return errors.New("not found")
	}
	l := len(DB)
	DB[idx] = DB[l-1]
	DB = DB[:l-1]
	return nil
}

func update(u problems.P8UserStruct) (problems.P8UserStruct, error) {
	u2, err := find(u.ID)
	if err != nil {
		return problems.P8UserStruct{}, nil
	}
	if u.Name != "" {
		u2.Name = u.Name
	}
	if u.Email != "" {
		u2.Email = u.Email
	}
	if u.Phone != "" {
		u2.Phone = u.Phone
	}
	if err := del(u.ID); err != nil {
		return problems.P8UserStruct{}, err
	}
	add(u2)
	return u2, nil
}

func Test(t *testing.T) {
	server := &http.Server{Addr: ":80"}
	go func() {
		http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				GetUser(w, r)
			case http.MethodPost:
				PostUser(w, r)
			case http.MethodPut:
				PutUser(w, r)
			case http.MethodDelete:
				DeleteUser(w, r)
			}
		})
		if err := http.ListenAndServe(":80", nil); err != nil {
			panic(err)
		}
	}()

	<-problems.P8User.WaitForReady()

	for _, scenario := range problems.P8User.TestCases {
		for _, tc := range scenario {
			fmt.Println(tc.GetName())
			res := tc.Run()
			fmt.Println(res.Passed, res.Result)
			if !res.Passed {
				t.FailNow()
			}
		}
	}

	server.Close()
}

func encodeJSON(a interface{}) []byte {
	b, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}
	return b
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	ids, ok := r.URL.Query()["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	id := ids[0]
	if !checkID(id) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	u, err := find(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	io.WriteString(w, string(encodeJSON(u)))
}

func PostUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	u := problems.P8UserStruct{}
	if err := json.Unmarshal(body, &u); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !(checkID(u.ID) &&
		checkPW(u.PW) &&
		checkName(u.Name) &&
		checkEmail(u.Email) &&
		checkPhone(u.Phone)) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	add(u)
	io.WriteString(w, string(encodeJSON(u)))
}

func PutUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	u := problems.P8UserStruct{}
	if err := json.Unmarshal(body, &u); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if u.Name != "" && !checkName(u.Name) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if u.Email != "" && !checkEmail(u.Email) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if u.Phone != "" && !checkPhone(u.Phone) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	u2, err := update(u)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	io.WriteString(w, string(encodeJSON(u2)))
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	u := problems.P8UserStruct{}
	if err := json.Unmarshal(body, &u); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := del(u.ID); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}

func checkOnlyAlphaNumeric(s string) bool {
	sources := []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'j', 'k', 'l', 'm', 'n',
		'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '_'}
	for _, a := range s {
		not := true
		for _, r := range sources {
			if a == r {
				not = false
				break
			}
		}
		if not {
			return false
		}
	}
	return true
}

func checkLength(s string, min, max int) bool {
	return min <= len(s) && len(s) <= max
}

func checkID(id string) bool {
	return checkLength(id, 5, 12) && checkOnlyAlphaNumeric(id)
}

func checkPW(pw string) bool {
	return checkLength(pw, 0, 50)
}

func checkName(name string) bool {
	return checkLength(name, 2, 50) && checkOnlyAlphaNumeric(name)
}

func checkEmail(email string) bool {
	return checkLength(email, 1, 50) && strings.Contains(email, "@")
}

func checkPhone(phone string) bool {
	return checkLength(phone, 1, 20) && strings.Contains(phone, "+")
}
