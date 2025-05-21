package utils

import "time"

// PtrToTimestamp 指针时间类型转换为int64时间戳
func PtrToTimestamp(t *time.Time) int64 {
	if t == nil {
		return 0
	}
	return t.Unix()
}

// TimestampToPtr int64时间戳转换为时间指针类型
func TimestampToPtr(ts int64) *time.Time {
	if ts == 0 {
		return nil
	}
	t := time.Unix(ts, 0)
	return &t
}
