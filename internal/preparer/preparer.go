package preparer

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/team-bonitto/bonitto/internal/model"
	"github.com/team-bonitto/bonitto/internal/queue/consumer"
	"github.com/team-bonitto/bonitto/internal/queue/producer"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"strings"
)

const QueueName = "preparer"

var _ producer.Producer = Input{}
var _ consumer.Consumer = Preparer{}

type Input struct {
	UserID  string
	Codes   []model.Code
	Problem model.Problem
	Lang    model.Language
}

type Preparer struct {
	TesterImage               string
	ExecutorImageJavascript   string
	ExecutorImageJava         string
	ExecutorImageGo           string
	ExecutorImageCpp          string
	ExecutorCommandJavascript string
	ExecutorCommandJava       string
	ExecutorCommandGo         string
	ExecutorCommandCpp        string
	RedisAddr                 string
}

func New(testerImage string,
	executorImageJavascript string,
	executorImageJava string,
	executorImageGo string,
	executorImageCpp string,
	executorCommandJavascript string,
	executorCommandJava string,
	executorCommandGo string,
	executorCommandCpp string,
	redisAddr string,
) (*Preparer, error) {
	suspects := []string{
		testerImage,
		executorImageJavascript,
		executorImageJava,
		executorImageGo,
		executorImageCpp,
		executorCommandJavascript,
		executorCommandJava,
		executorCommandGo,
		executorCommandCpp,
		redisAddr,
	}
	for _, s := range suspects {
		if s == "" {
			return nil, errors.New(fmt.Sprintf("needs this : %s", s))
		}
	}
	p := &Preparer{
		TesterImage:               testerImage,
		ExecutorImageJavascript:   executorImageJavascript,
		ExecutorImageJava:         executorImageJava,
		ExecutorImageGo:           executorImageGo,
		ExecutorImageCpp:          executorImageCpp,
		ExecutorCommandJavascript: executorCommandJavascript,
		ExecutorCommandJava:       executorCommandJava,
		ExecutorCommandGo:         executorCommandGo,
		ExecutorCommandCpp:        executorCommandCpp,
		RedisAddr:                 redisAddr,
	}
	return p, nil
}

func (i Input) GetQueueName() string {
	return QueueName
}

func (p Preparer) GetQueueName() string {
	return QueueName
}

func (p Preparer) Consume(a string) error {
	input := Input{}
	if err := json.Unmarshal([]byte(a), &input); err != nil {
		return err
	}
	cli, err := p.newKubeClient()
	if err != nil {
		return err
	}
	pod, err := p.buildPod(input)
	if err != nil {
		return err
	}
	if _, err := cli.CoreV1().Pods("default").Create(pod); err != nil {
		return err
	}
	return nil
}

func (i Input) Marshal() []byte {
	b, _ := json.Marshal(i)
	return b
}

func (p Preparer) newKubeClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	cli, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return cli, nil
}

func (p Preparer) getExecutorImage(input Input) (string, error) {
	switch input.Lang {
	case model.Javascript:
		return p.ExecutorImageJavascript, nil
	case model.Java:
		return p.ExecutorImageJava, nil
	case model.Go:
		return p.ExecutorImageGo, nil
	case model.Cpp:
		return p.ExecutorImageCpp, nil
	default:
		return "", errors.New(fmt.Sprintf("no image for %s", input.Lang))
	}
}

func (p Preparer) getExecutorCommand(input Input) (string, error) {
	switch input.Lang {
	case model.Javascript:
		return p.ExecutorCommandJavascript, nil
	case model.Java:
		return p.ExecutorCommandJava, nil
	case model.Go:
		return p.ExecutorCommandGo, nil
	case model.Cpp:
		return p.ExecutorCommandCpp, nil
	default:
		return "", errors.New(fmt.Sprintf("no command for %s", input.Lang))
	}
}

func (p Preparer) buildContainerExecutor(input Input) (*v1.Container, error) {
	image, err := p.getExecutorImage(input)
	if err != nil {
		return nil, err
	}
	command, err := p.getExecutorCommand(input)
	if err != nil {
		return nil, err
	}
	con := &v1.Container{
		Name:       "executor",
		Image:      image,
		Command:    []string{"/bin/sh", "-c"},
		Args:       []string{command},
		WorkingDir: "/",
		Env: []v1.EnvVar{{
			Name:  "CODE",
			Value: input.Codes[0].Content,
		}},
		Resources:       input.Problem.Resource,
		ImagePullPolicy: v1.PullIfNotPresent,
	}
	return con, nil
}

func (p Preparer) buildContainerTester(input Input) (*v1.Container, error) {
	image := p.TesterImage
	con := &v1.Container{
		Name:       "tester",
		Image:      image,
		WorkingDir: "/",
		Command:    []string{"/tester"},
		Env: []v1.EnvVar{{
			Name:  "PROBLEM_NO",
			Value: fmt.Sprintf("%d", input.Problem.No),
		}, {
			Name: "REDIS_URL",
			Value: p.RedisAddr,
		}, {
			Name: "USER_ID",
			Value: input.UserID,
		}},
		ImagePullPolicy: v1.PullIfNotPresent,
	}
	return con, nil
}

func (p Preparer) buildPod(input Input) (*v1.Pod, error) {
	conExecutor, err := p.buildContainerExecutor(input)
	if err != nil {
		return nil, err
	}
	conTester, err := p.buildContainerTester(input)
	if err != nil {
		return nil, err
	}
	pod := &v1.Pod{
		TypeMeta: v12.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},
		ObjectMeta: v12.ObjectMeta{
			Namespace:    "default",
			GenerateName: fmt.Sprintf("%d-%s", input.Problem.No, strings.ToLower(input.UserID)),
		},
		Spec: v1.PodSpec{
			Containers:    []v1.Container{*conExecutor, *conTester},
			RestartPolicy: v1.RestartPolicyNever,
		},
	}
	return pod, nil
}
