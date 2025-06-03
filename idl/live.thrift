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

struct WatchLiveReq{
    1:string classname
    2:i64 student_id
    3:string teacher_name
}

struct WatchLiveResp{
    1:string Addr

    255:base.BaseResp baseResp
}

struct PublishRegisterReq{
    1:i64 teacher_id
    2:string classname
    3:string teacher_name
    4:i64 deadline
}

struct PublishRegisterResp{

    255:base.BaseResp baseResp
}

struct StartRecordingReq{
    1:i64 teacher_id
    2:i64 course_id
    3:string classname
}

struct StartRecordingResp{

    255:base.BaseResp baseResp
}

struct StopRecordingReq{
     1:i64 teacher_id
     2:i64 course_id
     3:string classname
}

struct StopRecordingResp{

    255:base.BaseResp baseResp
}

service LiveService{
    GetStreamKeyResp GetStreamKey(1:GetStreamKeyReq req)
    WatchLiveResp WatchLive(1:WatchLiveReq req)
    PublishRegisterResp PublishRegister(1:PublishRegisterReq req)
    StartRecordingResp StartRecording(1:StartRecordingReq req)
    StopRecordingResp StopRecording(1:StopRecordingReq req)
}