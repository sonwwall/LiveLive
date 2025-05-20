namespace go livelive.course

include "base.thrift"

struct CreateCourseReq{
    1:string classname
    2:string description
    3:string teacher_id

}

struct CreateCourseResp{

    255:base.BaseResp baseResp
}

service CourseService{
    CreateCourseResp CreateCourse(1:CreateCourseReq req)
}