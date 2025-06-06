namespace go livelive.ai

include "base.thrift"

struct AnalyzeAudioReq{
    1:string wav_path
}

struct AnalyzeAudioResp{
    1:string summary

    255:base.BaseResp baseResp
}

service AIService{
    AnalyzeAudioResp AnalyzeAudio(1:AnalyzeAudioReq req)
}