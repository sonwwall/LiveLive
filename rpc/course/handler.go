package main

import (
	course "LiveLive/kitex_gen/livelive/course"
	"context"
)

// CourseServiceImpl implements the last service interface defined in the IDL.
type CourseServiceImpl struct{}

// CreateCourse implements the CourseServiceImpl interface.
func (s *CourseServiceImpl) CreateCourse(ctx context.Context, req *course.CreateCourseReq) (resp *course.CreateCourseResp, err error) {
	// TODO: Your code here...
	return
}
