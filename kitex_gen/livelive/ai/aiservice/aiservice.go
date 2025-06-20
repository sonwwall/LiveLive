// Code generated by Kitex v0.13.1. DO NOT EDIT.

package aiservice

import (
	ai "LiveLive/kitex_gen/livelive/ai"
	"context"
	"errors"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
)

var errInvalidMessageType = errors.New("invalid message type for service method handler")

var serviceMethods = map[string]kitex.MethodInfo{
	"AnalyzeAudio": kitex.NewMethodInfo(
		analyzeAudioHandler,
		newAIServiceAnalyzeAudioArgs,
		newAIServiceAnalyzeAudioResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"ChatWithAI": kitex.NewMethodInfo(
		chatWithAIHandler,
		newAIServiceChatWithAIArgs,
		newAIServiceChatWithAIResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
}

var (
	aIServiceServiceInfo                = NewServiceInfo()
	aIServiceServiceInfoForClient       = NewServiceInfoForClient()
	aIServiceServiceInfoForStreamClient = NewServiceInfoForStreamClient()
)

// for server
func serviceInfo() *kitex.ServiceInfo {
	return aIServiceServiceInfo
}

// for stream client
func serviceInfoForStreamClient() *kitex.ServiceInfo {
	return aIServiceServiceInfoForStreamClient
}

// for client
func serviceInfoForClient() *kitex.ServiceInfo {
	return aIServiceServiceInfoForClient
}

// NewServiceInfo creates a new ServiceInfo containing all methods
func NewServiceInfo() *kitex.ServiceInfo {
	return newServiceInfo(false, true, true)
}

// NewServiceInfo creates a new ServiceInfo containing non-streaming methods
func NewServiceInfoForClient() *kitex.ServiceInfo {
	return newServiceInfo(false, false, true)
}
func NewServiceInfoForStreamClient() *kitex.ServiceInfo {
	return newServiceInfo(true, true, false)
}

func newServiceInfo(hasStreaming bool, keepStreamingMethods bool, keepNonStreamingMethods bool) *kitex.ServiceInfo {
	serviceName := "AIService"
	handlerType := (*ai.AIService)(nil)
	methods := map[string]kitex.MethodInfo{}
	for name, m := range serviceMethods {
		if m.IsStreaming() && !keepStreamingMethods {
			continue
		}
		if !m.IsStreaming() && !keepNonStreamingMethods {
			continue
		}
		methods[name] = m
	}
	extra := map[string]interface{}{
		"PackageName": "ai",
	}
	if hasStreaming {
		extra["streaming"] = hasStreaming
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Thrift,
		KiteXGenVersion: "v0.13.1",
		Extra:           extra,
	}
	return svcInfo
}

func analyzeAudioHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*ai.AIServiceAnalyzeAudioArgs)
	realResult := result.(*ai.AIServiceAnalyzeAudioResult)
	success, err := handler.(ai.AIService).AnalyzeAudio(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newAIServiceAnalyzeAudioArgs() interface{} {
	return ai.NewAIServiceAnalyzeAudioArgs()
}

func newAIServiceAnalyzeAudioResult() interface{} {
	return ai.NewAIServiceAnalyzeAudioResult()
}

func chatWithAIHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*ai.AIServiceChatWithAIArgs)
	realResult := result.(*ai.AIServiceChatWithAIResult)
	success, err := handler.(ai.AIService).ChatWithAI(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newAIServiceChatWithAIArgs() interface{} {
	return ai.NewAIServiceChatWithAIArgs()
}

func newAIServiceChatWithAIResult() interface{} {
	return ai.NewAIServiceChatWithAIResult()
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) AnalyzeAudio(ctx context.Context, req *ai.AnalyzeAudioReq) (r *ai.AnalyzeAudioResp, err error) {
	var _args ai.AIServiceAnalyzeAudioArgs
	_args.Req = req
	var _result ai.AIServiceAnalyzeAudioResult
	if err = p.c.Call(ctx, "AnalyzeAudio", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) ChatWithAI(ctx context.Context, req *ai.ChatWithAIReq) (r *ai.ChatWithAIResp, err error) {
	var _args ai.AIServiceChatWithAIArgs
	_args.Req = req
	var _result ai.AIServiceChatWithAIResult
	if err = p.c.Call(ctx, "ChatWithAI", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}
