namespace go livelive.course

include "base.thrift"

struct CreateCourseReq{
    1:string classname
    2:string description
    3:i64 teacher_id

}

struct CreateCourseResp{

    255:base.BaseResp baseResp
}

struct JoinCourseReq{
    1:i64 student_id
    2:string classname
    3:string invitation_code
}

struct JoinCourseResp{

    255:base.BaseResp baseResp
}

struct CreateCourseInviteReq{
    1:string classname
    2:i64 max_usage
    3:i64 usage_count
    4:i64 expired_at
}

struct CreateCourseInviteResp{
    1:string invite_code

    255:base.BaseResp baseResp
}

service CourseService{
    CreateCourseResp CreateCourse(1:CreateCourseReq req)
    JoinCourseResp JoinCourse(1:JoinCourseReq req)
    CreateCourseInviteResp CreateCourseInvite(1:CreateCourseInviteReq req)
}