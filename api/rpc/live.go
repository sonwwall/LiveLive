package rpc

import (
	"LiveLive/kitex_gen/livelive/live"
	"LiveLive/kitex_gen/livelive/live/liveservice"
	"context"
	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
)

var liveClient liveservice.Client

func InitLiveRPCClient() {
	r, err := etcd.NewEtcdResolver([]string{"http://127.0.0.1:2379"})
	if err != nil {
		panic(err)
	}

	c, err := liveservice.NewClient("livelive.live", client.WithResolver(r))
	if err != nil {
		panic(err)
	}

	liveClient = c
}

func GetStreamKey(ctx context.Context, req *live.GetStreamKeyReq) (*live.GetStreamKeyResp, error) {
	return liveClient.GetStreamKey(ctx, req)
}

func WatchLive(ctx context.Context, req *live.WatchLiveReq) (*live.WatchLiveResp, error) {
	return liveClient.WatchLive(ctx, req)
}
