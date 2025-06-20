// Code generated by Kitex v0.13.1. DO NOT EDIT.

package liveservice

import (
	live "LiveLive/kitex_gen/livelive/live"
	"context"
	client "github.com/cloudwego/kitex/client"
	callopt "github.com/cloudwego/kitex/client/callopt"
)

// Client is designed to provide IDL-compatible methods with call-option parameter for kitex framework.
type Client interface {
	GetStreamKey(ctx context.Context, req *live.GetStreamKeyReq, callOptions ...callopt.Option) (r *live.GetStreamKeyResp, err error)
	WatchLive(ctx context.Context, req *live.WatchLiveReq, callOptions ...callopt.Option) (r *live.WatchLiveResp, err error)
	PublishRegister(ctx context.Context, req *live.PublishRegisterReq, callOptions ...callopt.Option) (r *live.PublishRegisterResp, err error)
	StartRecording(ctx context.Context, req *live.StartRecordingReq, callOptions ...callopt.Option) (r *live.StartRecordingResp, err error)
	StopRecording(ctx context.Context, req *live.StopRecordingReq, callOptions ...callopt.Option) (r *live.StopRecordingResp, err error)
}

// NewClient creates a client for the service defined in IDL.
func NewClient(destService string, opts ...client.Option) (Client, error) {
	var options []client.Option
	options = append(options, client.WithDestService(destService))

	options = append(options, opts...)

	kc, err := client.NewClient(serviceInfoForClient(), options...)
	if err != nil {
		return nil, err
	}
	return &kLiveServiceClient{
		kClient: newServiceClient(kc),
	}, nil
}

// MustNewClient creates a client for the service defined in IDL. It panics if any error occurs.
func MustNewClient(destService string, opts ...client.Option) Client {
	kc, err := NewClient(destService, opts...)
	if err != nil {
		panic(err)
	}
	return kc
}

type kLiveServiceClient struct {
	*kClient
}

func (p *kLiveServiceClient) GetStreamKey(ctx context.Context, req *live.GetStreamKeyReq, callOptions ...callopt.Option) (r *live.GetStreamKeyResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetStreamKey(ctx, req)
}

func (p *kLiveServiceClient) WatchLive(ctx context.Context, req *live.WatchLiveReq, callOptions ...callopt.Option) (r *live.WatchLiveResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.WatchLive(ctx, req)
}

func (p *kLiveServiceClient) PublishRegister(ctx context.Context, req *live.PublishRegisterReq, callOptions ...callopt.Option) (r *live.PublishRegisterResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.PublishRegister(ctx, req)
}

func (p *kLiveServiceClient) StartRecording(ctx context.Context, req *live.StartRecordingReq, callOptions ...callopt.Option) (r *live.StartRecordingResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.StartRecording(ctx, req)
}

func (p *kLiveServiceClient) StopRecording(ctx context.Context, req *live.StopRecordingReq, callOptions ...callopt.Option) (r *live.StopRecordingResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.StopRecording(ctx, req)
}
