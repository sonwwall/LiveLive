package code

var (
	ErrDB                 int64 = 40001 //数据库错误
	ErrCourseNotExist     int64 = 40002 //课程不存在
	ErrGetStreamKey       int64 = 40003 //从livego获取streamkey失败
	ErrTeacherNotExist    int64 = 40004 //该老师不存在
	ErrNoCoursePermission int64 = 40005 //该学生未加入课程
	ErrRedis              int64 = 40006 //Redis错误
	ErrStartRecording     int64 = 40007 //录制失败
	ErrStopRecording      int64 = 40008 //停止录制失败
)
