package code

var (
	ErrDB            int64 = 20001 //数据库错误
	ErrUserNotExists int64 = 20002 //用户不存在
	ErrUsernameExist int64 = 20003 //用户名已存在
	ErrEmailExist    int64 = 20004 //邮箱已存在
	ErrPhoneExist    int64 = 20005 //手机号已存在
)
