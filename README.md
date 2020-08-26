# Bonitto

Live at : [bonitto.eunchan.com](http://bonitto.eunchan.com)
![image](https://imgur.com/0f53lw5.png)

Youtube : [https://www.youtube.com/watch?v=7c8_krY4jr4](https://www.youtube.com/watch?v=7c8_krY4jr4)
[![image](https://i.imgur.com/IdrLFEs.png)](https://www.youtube.com/watch?v=7c8_krY4jr4)

# Description
**실무 도메인 문제** 구현을 위한 코딩 채점 사이트 입니다.

[발표 PPT](https://drive.google.com/file/d/1p4-Cg-alPF9UmllmbkTNlXpB720f0BZn/view?usp=sharing)

실무 도메인 문제의 예시는 아래와 같습니다.
- User Authentication and Authorization
- Contents Management System
- Realtime Chat Server
- Image Resizing, Caching, Serving and Uploading

# Quck Start
이 프로젝트는 모노레포(단일저장소)로서 프론트와 백엔드 그리고 여러 마이크로서비스들을 포함하고 있습니다.

개발 실행과 배포 실행 두가지 종류가 있습니다.
1. 개발실행
    1. dockerizing backend & push to docker repo
        ```bash
       make docker-back
       docker push bonitto-back:lateset 도커저장소이미지주소
       ```
    2. setup kubernetes with backend image url you uploaded
        ```bash
       vi infra/kubernetes/configmap.yaml # example 따라 잘 만들어주시고
       kubectl apply -k infra/kubernetes
       ```
    3. run backend
        ```bash
        go run cmd/preparer/preparer.go
        go run cmd/recorder/recorder.go
        go run cmd/tester/tester.go
        go run cmd/web/web.go
       ```
    4. run frontend
        ```bash
        cd front
        yarn
        yarn start
        ```
    5. go to `http://localhost:3000`
2. 배포실행
    1. dockerizing backend, frontend
        ```bash
        make docker-back
        docker push bonitto-back:latest 도커저장소이미지주소
        make docker-front
        docker push bonitto-front:latest 도커저장소이미지주소
        ```
    2. setup kubernetes
        ```bash
       vi infra/kubernetes/configmap.yaml # example 따라 잘 만들어주시고
       kubectl apply -k infra/kubernetes
       ```
    3. go to service endpoint

# Flow
![https://i.imgur.com/mDa2o5K.png](https://i.imgur.com/mDa2o5K.png)

사용자의 소스를 실행시키고 평가하는게 주 업무입니다.

# Architecture
![https://i.imgur.com/x8aC7Vr.png](https://i.imgur.com/x8aC7Vr.png)

![https://i.imgur.com/haPtX3A.png](https://i.imgur.com/haPtX3A.png)

위와 같은 Event Driven 구조로 되어있습니다.

더 자세한 설명은 [PPT](https://drive.google.com/file/d/1p4-Cg-alPF9UmllmbkTNlXpB720f0BZn/view?usp=sharing) 를 참고해주세요.

# Installation
Project Directory Structure
```bash
.
├── cmd             # microservices' entry point
│   ├── preparer
│   ├── recorder
│   ├── tester
│   └── web
├── front           # CRA 앱
├── infra
│   └── kubernetes  # kubernetes kustomization 한방 인프라
└── internal        # microservice share codes
    ├── model
    ├── notifier
    ├── preparer
    ├── problems
    ├── queue
    └── recorder
```

## 1. Infrastructure (Kubernetes)
Directory Structure:
```
infra/kubernetes
├── cluster-role-binding.yaml
*** configmap.yaml              # 밑의 example 기반으로 직접 만들어야 합니다
├── configmap.yaml.example
├── deployment-preparer.yaml
├── deployment-recorder.yaml
├── deployment-web.yaml
├── kustomization.yaml
├── service-account.yaml
└── service.yaml
```
`configmap.yaml`: 직접 만들어주셔야 합니다
```yaml
# configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: bonitto
data:
  REDIS_URL:                포트번호까지 포함한 redis url
  WORKER_INTERVAL_MS:       마이크로서비스들 폴링 인터벌 (1000)
  KUBERNETES_SERVICE_HOST:  쿠버네티스 master url
  KUBERNETES_SERVICE_PORT:  쿠버네티스 master port
  IMAGE_TESTER: "ckcks12/bonitto-back:v0.17" # tester 이미지 입니다
  IMAGE_JAVASCRIPT: "node:14.8.0-alpine3.11"
  IMAGE_JAVA: "alpine3"
  IMAGE_GOLANG: "golang:1.15.0-alpine3.12"
  IMAGE_CPP: "alpine3A"
  COMMAND_JAVASCRIPT: 'echo "${CODE}" > /code.js && node /code.js& echo start; sleep 300; echo end'
  COMMAND_JAVA: "whoami"
  COMMAND_GOLANG: 'echo "${CODE}" > /code.go && go run /code.go& echo start; sleep 300; echo end'
  COMMAND_CPP: "whoami"
```
`kustomization.yaml`: 커스텀 이미지를 구워 사용할 경우 수정해주셔야 합니다.
```yaml
# kustomization.yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namespace: default
resources:
  - cluster-role-binding.yaml
  - deployment-web.yaml
  - deployment-preparer.yaml
  - deployment-recorder.yaml
  - service-account.yaml
  - configmap.yaml
  - service.yaml
images:
  - name: bonitto-front
    newName: ckcks12/bonitto-front  # 프론트 이미지 입니다.
    newTag: v0.11                   # 프론트 이미지 태그 입니다.
  - name: bonitto-back
    newName: ckcks12/bonitto-back   # 백엔드 이미지 입니다.
    newTag: v0.17                   # 백엔드 이미지 태그 입니다.
```
아래 명령어로 kustomization 을 apply 하여 쿠버네티스에 띄울 수 있습니다.
```bash
cd infra/kubernetes
kubectl apply -k .
```

## 2. Backend / Microservice Build
Directory Structure:
```bash
.
├── cmd             # microservices' entry point
│   ├── preparer    
│   ├── recorder
│   ├── tester
│   └── web
└── internal        # microservice share codes
    ├── model       # 기본 모델들 정의되어 있는 곳
    ├── notifier    # notifier producer/consumer 인터페이스 구현
    ├── preparer    # preparer producer/consumer 인터페이스 구현
    ├── problems    # 문제 정의(설명, 테스트케이스 등)
    ├── queue       # producer/consumer 인터페이스 정의 및 RedisProducer RedisConsumer 정의
    └── recorder    # recorder producer/consumer 인터페이스 구현
```
Build and Run:
```bash
# preparer
go run cmd/preparer/preparer.go

# recorder 
go run cmd/recorder/recorder.go

# tester 
go run cmd/tester/tester.go

# web 
go run cmd/web/web.go
```
Dockerizing:
```bash
make docker-back
docker run --rm -p 8080:8080 bonitto-back:latest /preparer
docker run --rm -p 8080:8080 bonitto-back:latest /recorder
docker run --rm -p 8080:8080 bonitto-back:latest /tester
docker run --rm -p 8080:8080 bonitto-back:latest /web
```

## 3. Front Build
Directory Structure:
```bash
front/src
├── App.tsx                     # router 설정 등 React App Main Entry Point
├── lib
│   ├── api.ts                  # 서버 통신하는 API 함수
│   └── type.ts                 # 모델 정의
└── page
    ├── MainPage.tsx            # 메인 페이지
    ├── ProblemListPage.tsx     # 문제 목록 페이지
    └── ProblemPage.tsx         # 문제 상세보기 페이지
```
Run:
```bash
cd front
yarn start
```
Build:
```bash
cd front
yarn build
```
Dockerizing: 

> 주의할 점: nginx.conf 에서 websocket proxy 주소가 localhost 입니다.
항상 세트로 front 와 back 이 배포된다 설계한 것 입니다.
따라서 loopback 이 다른 컨테이너에게 연결되게 해야합니다.
Linux 라면 `network_mode: host` 으로 간단히 해결 가능합니다.
Mac 과 Window 라면 차라리 nginx.conf 에서 proxy 주소를 바꾸신 뒤 이에 맞게 설정해주세요.
(예를 들어 docker-compose 라면 service 이름을 바꾸기)

가장 추천드리는 방법은 로컬에서는 위의 `yarn start` 로, 배포할때만 docker 를 사용하시는 것 입니다.
```bash
make docker-front
docker run --rm -p 80:80 bonitto-front
```

# Usage

Youtube : [https://www.youtube.com/watch?v=7c8_krY4jr4](https://www.youtube.com/watch?v=7c8_krY4jr4)

8번 문제 User 에 대한 정답 코드 golang 입니다.
```go
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type User struct {
	ID      string `json:"id"`
	PW      string `json:"pw"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Visit   int    `json:"visit"`
	Deleted bool   `json:"deleted"`
}

var DB = make([]User, 0)

func add(u User) {
	DB = append(DB, u)
}

func find(id string) (User, error) {
	for _, u := range DB {
		if u.ID == id {
			return u, nil
		}
	}
	return User{}, errors.New("not found")
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

func update(u User) (User, error) {
	u2, err := find(u.ID)
	if err != nil {
		return User{}, nil
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
		return User{}, err
	}
	add(u2)
	return u2, nil
}

func main() {
	server := &http.Server{Addr: ":80"}
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
	u := User{}
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
	u := User{}
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
	u := User{}
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
```

# Credit
Eunchan Lee(@ckcks12) - 3unch4n@gmail.com

# License
Beer License


