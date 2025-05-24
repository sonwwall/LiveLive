package utils

import (
	"LiveLive/utils/md5"
	"fmt"
	"net/url"
	"time"
)

func GeneratePlayURL(baseURL, uri, secret string, Uid int64, ttl time.Duration) string {
	uid := fmt.Sprintf("%d", Uid)
	exp := time.Now().Add(ttl).Unix()
	sign := md5.GenerateSign(uri, uid, fmt.Sprintf("%d", exp), secret)

	// 创建 URL
	u := url.URL{
		Scheme: "http",
		Host:   baseURL, // 如 localhost:8080
		Path:   uri,     // 如 /live/movie.flv
	}

	// 添加查询参数
	q := u.Query()
	q.Set("exp", fmt.Sprintf("%d", exp))
	q.Set("sign", sign)
	q.Set("uid", uid)

	//q.Set("uid", uid)
	//q.Set("sign", sign)
	u.RawQuery = q.Encode()

	return u.String()
}
