package main

import (
	live "LiveLive/kitex_gen/livelive/live"
	"context"
)

// LiveServiceImpl implements the last service interface defined in the IDL.
type LiveServiceImpl struct{}

// GetStreamKey implements the LiveServiceImpl interface.
func (s *LiveServiceImpl) GetStreamKey(ctx context.Context, req *live.GetStreamKeyReq) (resp *live.GetStreamKeyResp, err error) {
	// TODO: Your code here...
	return
}
