namespace go livelive.live

include "base.thrift"

struct GetStreamKeyReq{
    1:i64 teacher_id
    2:string classname
}

struct GetStreamKeyResp{
    1:string rtmp_url
    2:string stream_key

    255:base.BaseResp baseResp
}

service LiveService{
    GetStreamKeyResp GetStreamKey(1:GetStreamKeyReq req)
}