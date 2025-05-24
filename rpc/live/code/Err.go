package code

var (
	ErrDB              int64 = 40001 //数据库错误
	ErrCourseNotExist  int64 = 40002 //课程不存在
	ErrGetStreamKey    int64 = 40003 //从livego获取streamkey失败
	ErrTeacherNotExist int64 = 40004 //该老师不存在
)
