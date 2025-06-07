namespace go livelive.websocket

include "base.thrift"

struct BroadcastToCourseReq{
    i64 course_id
    binary data//[]byte类型
}

struct BroadcastToCourseResp{
     255:base.BaseResp baseResp
}

//答题统计
struct AggregateAnswersReq{
    i64 question_id
    i64 course_id
    i8 correct_answer
}


struct AggregateAnswersResp{
    double accuracy
}

//判断题答题统计
struct AggregateTrueOrFalseAnswersReq{
     i64 question_id
     i64 course_id
     i8 correct_answer
}

struct AggregateTrueOrFalseAnswersResp{
    double accuracy
}

//  统计签到结果
struct CountRegisterReq{
    i64 course_id

}

struct CountRegisterResp{


    255:base.BaseResp baseResp
}

service WebsocketService{
    BroadcastToCourseResp BroadcastToCourse(1:BroadcastToCourseReq req)
    AggregateAnswersResp AggregateAnswers(1:AggregateAnswersReq req)
    AggregateTrueOrFalseAnswersResp AggregateTrueOrFalseAnswers(1:AggregateTrueOrFalseAnswersReq req)
    CountRegisterResp CountRegister(1:CountRegisterReq req)
}