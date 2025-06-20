// Code generated by Kitex v0.13.1. DO NOT EDIT.

package websocketservice

import (
	websocket "LiveLive/kitex_gen/livelive/websocket"
	"context"
	"errors"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
)

var errInvalidMessageType = errors.New("invalid message type for service method handler")

var serviceMethods = map[string]kitex.MethodInfo{
	"BroadcastToCourse": kitex.NewMethodInfo(
		broadcastToCourseHandler,
		newWebsocketServiceBroadcastToCourseArgs,
		newWebsocketServiceBroadcastToCourseResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"AggregateAnswers": kitex.NewMethodInfo(
		aggregateAnswersHandler,
		newWebsocketServiceAggregateAnswersArgs,
		newWebsocketServiceAggregateAnswersResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"AggregateTrueOrFalseAnswers": kitex.NewMethodInfo(
		aggregateTrueOrFalseAnswersHandler,
		newWebsocketServiceAggregateTrueOrFalseAnswersArgs,
		newWebsocketServiceAggregateTrueOrFalseAnswersResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"CountRegister": kitex.NewMethodInfo(
		countRegisterHandler,
		newWebsocketServiceCountRegisterArgs,
		newWebsocketServiceCountRegisterResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
}

var (
	websocketServiceServiceInfo                = NewServiceInfo()
	websocketServiceServiceInfoForClient       = NewServiceInfoForClient()
	websocketServiceServiceInfoForStreamClient = NewServiceInfoForStreamClient()
)

// for server
func serviceInfo() *kitex.ServiceInfo {
	return websocketServiceServiceInfo
}

// for stream client
func serviceInfoForStreamClient() *kitex.ServiceInfo {
	return websocketServiceServiceInfoForStreamClient
}

// for client
func serviceInfoForClient() *kitex.ServiceInfo {
	return websocketServiceServiceInfoForClient
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
	serviceName := "WebsocketService"
	handlerType := (*websocket.WebsocketService)(nil)
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
		"PackageName": "websocket",
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

func broadcastToCourseHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*websocket.WebsocketServiceBroadcastToCourseArgs)
	realResult := result.(*websocket.WebsocketServiceBroadcastToCourseResult)
	success, err := handler.(websocket.WebsocketService).BroadcastToCourse(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newWebsocketServiceBroadcastToCourseArgs() interface{} {
	return websocket.NewWebsocketServiceBroadcastToCourseArgs()
}

func newWebsocketServiceBroadcastToCourseResult() interface{} {
	return websocket.NewWebsocketServiceBroadcastToCourseResult()
}

func aggregateAnswersHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*websocket.WebsocketServiceAggregateAnswersArgs)
	realResult := result.(*websocket.WebsocketServiceAggregateAnswersResult)
	success, err := handler.(websocket.WebsocketService).AggregateAnswers(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newWebsocketServiceAggregateAnswersArgs() interface{} {
	return websocket.NewWebsocketServiceAggregateAnswersArgs()
}

func newWebsocketServiceAggregateAnswersResult() interface{} {
	return websocket.NewWebsocketServiceAggregateAnswersResult()
}

func aggregateTrueOrFalseAnswersHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*websocket.WebsocketServiceAggregateTrueOrFalseAnswersArgs)
	realResult := result.(*websocket.WebsocketServiceAggregateTrueOrFalseAnswersResult)
	success, err := handler.(websocket.WebsocketService).AggregateTrueOrFalseAnswers(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newWebsocketServiceAggregateTrueOrFalseAnswersArgs() interface{} {
	return websocket.NewWebsocketServiceAggregateTrueOrFalseAnswersArgs()
}

func newWebsocketServiceAggregateTrueOrFalseAnswersResult() interface{} {
	return websocket.NewWebsocketServiceAggregateTrueOrFalseAnswersResult()
}

func countRegisterHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*websocket.WebsocketServiceCountRegisterArgs)
	realResult := result.(*websocket.WebsocketServiceCountRegisterResult)
	success, err := handler.(websocket.WebsocketService).CountRegister(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newWebsocketServiceCountRegisterArgs() interface{} {
	return websocket.NewWebsocketServiceCountRegisterArgs()
}

func newWebsocketServiceCountRegisterResult() interface{} {
	return websocket.NewWebsocketServiceCountRegisterResult()
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) BroadcastToCourse(ctx context.Context, req *websocket.BroadcastToCourseReq) (r *websocket.BroadcastToCourseResp, err error) {
	var _args websocket.WebsocketServiceBroadcastToCourseArgs
	_args.Req = req
	var _result websocket.WebsocketServiceBroadcastToCourseResult
	if err = p.c.Call(ctx, "BroadcastToCourse", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) AggregateAnswers(ctx context.Context, req *websocket.AggregateAnswersReq) (r *websocket.AggregateAnswersResp, err error) {
	var _args websocket.WebsocketServiceAggregateAnswersArgs
	_args.Req = req
	var _result websocket.WebsocketServiceAggregateAnswersResult
	if err = p.c.Call(ctx, "AggregateAnswers", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) AggregateTrueOrFalseAnswers(ctx context.Context, req *websocket.AggregateTrueOrFalseAnswersReq) (r *websocket.AggregateTrueOrFalseAnswersResp, err error) {
	var _args websocket.WebsocketServiceAggregateTrueOrFalseAnswersArgs
	_args.Req = req
	var _result websocket.WebsocketServiceAggregateTrueOrFalseAnswersResult
	if err = p.c.Call(ctx, "AggregateTrueOrFalseAnswers", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) CountRegister(ctx context.Context, req *websocket.CountRegisterReq) (r *websocket.CountRegisterResp, err error) {
	var _args websocket.WebsocketServiceCountRegisterArgs
	_args.Req = req
	var _result websocket.WebsocketServiceCountRegisterResult
	if err = p.c.Call(ctx, "CountRegister", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}
