namespace go livelive.ai

include "base.thrift"

struct AnalyzeAudioReq{
    1:string wav_path
}

struct AnalyzeAudioResp{
    1:string summary

    255:base.BaseResp baseResp
}

struct ChatWithAIReq{
    1:string content
}

struct ChatWithAIResp{
    1:string content

    255:base.BaseResp baseResp
}

service AIService{
    AnalyzeAudioResp AnalyzeAudio(1:AnalyzeAudioReq req)
    ChatWithAIResp ChatWithAI(1:ChatWithAIReq req)
}