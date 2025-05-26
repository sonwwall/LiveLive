package rpc

import (
	"LiveLive/kitex_gen/livelive/quiz"
	"LiveLive/kitex_gen/livelive/quiz/quizservice"
	"context"
	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
)

var quizClient quizservice.Client

func InitQuizRPCClient() {
	r, err := etcd.NewEtcdResolver([]string{"http://127.0.0.1:2379"})
	if err != nil {
		panic(err)
	}

	c, err := quizservice.NewClient("livelive.quiz", client.WithResolver(r))
	if err != nil {
		panic(err)
	}

	quizClient = c
}

func PublishChoiceQuestion(ctx context.Context, req *quiz.PublishChoiceQuestionReq) (*quiz.PublishChoiceQuestionResp, error) {
	return quizClient.PublishChoiceQuestion(ctx, req)
}
