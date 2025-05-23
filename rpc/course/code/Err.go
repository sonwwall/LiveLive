package code

var (
	ErrCourseExist       int64 = 30001 //课程已存在
	ErrDB                int64 = 30002 //数据库错误
	ErrCourseNotExist    int64 = 30003 //课程不存在
	ErrInviteCodeInvalid int64 = 30004 //验证码无效
	ErrTeacherNotExist   int64 = 30005 //教师不存在
)
