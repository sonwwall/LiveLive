package rpc

import (
	"LiveLive/kitex_gen/livelive/live"
	"context"
)

func GetStreamKey(ctx context.Context, req *live.GetStreamKeyReq) (*live.GetStreamKeyResp, error) {
	return &live.GetStreamKeyResp{}, nil
}
