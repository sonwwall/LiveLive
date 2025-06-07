namespace go livelive.quiz

include "base.thrift"

struct PublishChoiceQuestionReq{
    1:i64 teacher_id
    2:i64 course_id//直接从前端传入
    3:string title
    4:list<string> options
    5:i8 answer
    6:i64 deadline
}

struct PublishChoiceQuestionResp{

    255:base.BaseResp baseResp
}

struct PublishTrueOrFalseQuestionReq{
     1:i64 teacher_id
     2:i64 course_id//直接从前端传入
     3:string title
     5:i8 answer
     6:i64 deadline
}

struct PublishTrueOrFalseQuestionResp{

    255:base.BaseResp baseResp
}

service QuizService{
    PublishChoiceQuestionResp PublishChoiceQuestion(1:PublishChoiceQuestionReq req)
    PublishTrueOrFalseQuestionResp PublishTrueOrFalseQuestion(1:PublishTrueOrFalseQuestionReq req)

}

